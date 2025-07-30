# Архитектура TCP-сервера на Go для обработки XML/JSON

## Описание проекта

TCP-сервер на порту 8080, который принимает XML-документы через POST-запросы, преобразует их в JSON и отправляет на внешний сервис. Включает авторизацию, асинхронную обработку и логирование.

## Структура проекта

```
tcp-xml-processor/
├── main.go                 # Точка входа приложения
├── internal/
│   ├── handler/
│   │   ├── handler.go      # HTTP handler с бизнес-логикой
│   │   └── auth.go         # Логика авторизации
│   ├── models/
│   │   ├── user.go         # Структуры данных User, Users
│   │   └── json_user.go    # Структура UserJSON
│   ├── converter/
│   │   ├── xml_parser.go   # Парсинг XML
│   │   ├── json_converter.go # Конвертация в JSON
│   │   └── age_groups.go   # Логика возрастных групп
│   ├── client/
│   │   └── external_client.go # HTTP-клиент для внешнего сервиса
│   ├── logger/
│   │   └── logger.go       # Логирование в файл
│   └── processor/
│       └── async_processor.go # Асинхронная обработка с goroutines
├── config/
│   └── config.go           # Конфигурация приложения
├── logs/
│   └── provider.log        # Файл логов (создается автоматически)
├── examples/
│   ├── sample_request.xml  # Пример входящего XML
│   └── sample_response.json # Пример выходного JSON
├── go.mod
├── go.sum
└── README.md
```

## Архитектура решения

### Диаграмма компонентов

```
┌─────────────────────────────────────────────────────────────┐
│                      TCP Server :8080                       │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐    ┌──────────────┐    ┌─────────────────┐ │
│  │   Logger    │    │   Handler    │    │ External Client │ │
│  │             │◄───┤              ├────┤                 │ │
│  └─────────────┘    │ - logger     │    │ HTTP Client     │ │
│                     │ - httpClient │    └─────────────────┘ │
│  ┌─────────────┐    │ - authKey    │                       │
│  │XML Parser   │◄───┤              │                       │
│  └─────────────┘    └──────────────┘                       │
│                             │                              │
│  ┌─────────────┐    ┌──────────────┐    ┌─────────────────┐ │
│  │JSON         │◄───┤ Async        │────┤ Age Groups      │ │
│  │Converter    │    │ Processor    │    │ Logic           │ │
│  └─────────────┘    └──────────────┘    └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                             │
                             ▼
                ┌─────────────────────────┐
                │  External Service       │
                │  localhost:8081/users   │
                └─────────────────────────┘
```

## Подробный план решения

### Этап 1: Структуры данных

#### Модели данных (internal/models/)

**user.go**
```go
// User - структура для парсинга XML
type User struct {
    ID    string `xml:"id,attr"`
    Name  string `xml:"name"`
    Email string `xml:"email"`
    Age   int    `xml:"age"`
}

// Users - корневая структура XML документа
type Users struct {
    XMLName xml.Name `xml:"users"`
    Users   []User   `xml:"user"`
}
```

**json_user.go**
```go
// UserJSON - структура для JSON ответа
type UserJSON struct {
    ID       string `json:"id"`
    FullName string `json:"full_name"`
    Email    string `json:"email"`
    AgeGroup string `json:"age_group"`
}
```

### Этап 2: Логирование (internal/logger/)

**logger.go**
```go
type Logger struct {
    file *os.File
    mu   sync.Mutex
}

// Функции:
func NewLogger(filename string) (*Logger, error)    // Создание логгера
func (l *Logger) Log(level, message string)         // Запись в лог
func (l *Logger) LogError(err error, context string) // Логирование ошибок
func (l *Logger) LogInfo(message string)            // Информационные сообщения
func (l *Logger) Close() error                      // Закрытие файла
```

**Особенности:**
- Thread-safe запись с использованием mutex
- Форматирование с временными метками
- Различные уровни логирования (INFO, ERROR, DEBUG)

### Этап 3: HTTP Handler (internal/handler/)

**handler.go**
```go
type Handler struct {
    logger     *Logger
    httpClient *http.Client
    authKey    string
    processor  *processor.AsyncProcessor
}

// Функции:
func NewHandler(logger *Logger, authKey string) *Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request)
func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request)
```

**auth.go**
```go
// Функции авторизации:
func (h *Handler) validateAuth(r *http.Request) bool
func (h *Handler) extractAuthKey(r *http.Request) string
```

### Этап 4: Парсинг и конвертация (internal/converter/)

**xml_parser.go**
```go
// Функции:
func ParseXML(data []byte) (*models.Users, error)
func ValidateXML(data []byte) error
```

**json_converter.go**
```go
// Функции:
func ConvertToJSON(user models.User) models.UserJSON
func ConvertUsersSlice(users []models.User) []models.UserJSON
```

**age_groups.go**
```go
// Функции:
func GetAgeGroup(age int) string
func ValidateAge(age int) bool

// Константы возрастных групп:
const (
    AgeGroupYoung  = "до 25"
    AgeGroupMiddle = "от 25 до 35"
    AgeGroupSenior = "старше 35"
)
```

### Этап 5: Асинхронная обработка (internal/processor/)

**async_processor.go**
```go
type AsyncProcessor struct {
    logger      *Logger
    workerCount int
}

// Функции:
func NewAsyncProcessor(logger *Logger, workerCount int) *AsyncProcessor
func (p *AsyncProcessor) ProcessUsers(users []models.User) ([]models.UserJSON, error)
func (p *AsyncProcessor) processUserWorker(userChan <-chan models.User, resultChan chan<- models.UserJSON, errorChan chan<- error)
```

**Архитектура обработки:**
```
Input Users → Worker Pool → JSON Results → Aggregation
     │              │              │            │
     └─── User1 ────┼─── Worker1 ───┼─── JSON1 ──┤
     └─── User2 ────┼─── Worker2 ───┼─── JSON2 ──┼─── Final Result
     └─── User3 ────┼─── Worker3 ───┼─── JSON3 ──┤
     └─── UserN ────┼─── WorkerN ───┼─── JSONN ──┘
```

### Этап 6: Внешний клиент (internal/client/)

**external_client.go**
```go
type ExternalClient struct {
    httpClient *http.Client
    baseURL    string
    logger     *Logger
}

// Функции:
func NewExternalClient(baseURL string, logger *Logger) *ExternalClient
func (c *ExternalClient) SendUsers(users []models.UserJSON) (*http.Response, error)
func (c *ExternalClient) handleResponse(resp *http.Response) error
```

### Этап 7: Конфигурация (config/)

**config.go**
```go
type Config struct {
    ServerPort      string
    AuthKey         string
    ExternalURL     string
    LogFile         string
    WorkerCount     int
    RequestTimeout  time.Duration
}

// Функции:
func LoadConfig() *Config
func (c *Config) Validate() error
```

## Поток выполнения

### 1. Инициализация (main.go)
```go
func main() {
    // 1. Загрузка конфигурации
    cfg := config.LoadConfig()
    
    // 2. Создание логгера
    logger, err := logger.NewLogger(cfg.LogFile)
    if err != nil {
        log.Fatal("Failed to create logger:", err)
    }
    defer logger.Close()
    
    // 3. Создание обработчика с встроенными компонентами
    handler := handler.NewHandler(logger, cfg.AuthKey)
    
    // 4. Настройка и запуск HTTP-сервера
    server := &http.Server{
        Addr:    ":" + cfg.ServerPort,
        Handler: handler,
    }
    
    logger.LogInfo("Server starting on port " + cfg.ServerPort)
    log.Fatal(server.ListenAndServe())
}
```

### 2. Обработка запроса
```
HTTP POST Request
        │
        ▼
┌─────────────────┐
│ Validate Method │ ──── (не POST) ────► HTTP 405
└─────────────────┘
        │ (POST)
        ▼
┌─────────────────┐
│ Check Auth Key  │ ──── (неверный) ───► HTTP 401
└─────────────────┘
        │ (валидный)
        ▼
┌─────────────────┐
│ Parse XML Body  │ ──── (ошибка) ─────► HTTP 400
└─────────────────┘
        │ (успешно)
        ▼
┌─────────────────┐
│ Async Processing│ ──── (ошибка) ─────► HTTP 500
└─────────────────┘
        │ (успешно)
        ▼
┌─────────────────┐
│ Send to External│ ──── (ошибка) ─────► HTTP 502
└─────────────────┘
        │ (успешно)
        ▼
┌─────────────────┐
│ Return Response │ ────────────────────► HTTP 200
└─────────────────┘
```

### 3. Асинхронная обработка пользователей
```go
// Создание каналов
userChan := make(chan models.User, len(users))
resultChan := make(chan models.UserJSON, len(users))
errorChan := make(chan error, len(users))

// Запуск worker pool
for i := 0; i < workerCount; i++ {
    go p.processUserWorker(userChan, resultChan, errorChan)
}

// Отправка задач
for _, user := range users {
    userChan <- user
}
close(userChan)

// Сбор результатов
var results []models.UserJSON
var errors []error

for i := 0; i < len(users); i++ {
    select {
    case result := <-resultChan:
        results = append(results, result)
    case err := <-errorChan:
        errors = append(errors, err)
    }
}
```

## Обработка ошибок

### Уровни ошибок:
1. **HTTP уровень**: неверный метод, авторизация
2. **Парсинг уровень**: некорректный XML
3. **Обработка уровень**: ошибки конвертации
4. **Сеть уровень**: недоступность внешнего сервиса

### Стратегия логирования:
```go
// Каждая ошибка логируется с контекстом
logger.LogError(err, "Failed to parse XML request")
logger.LogInfo("Successfully processed 10 users")
logger.LogError(networkErr, "External service unavailable")
```

## Синхронизация и конкурентность

### Механизмы синхронизации:
- **Channels** - для передачи данных между goroutines
- **sync.Mutex** - для защиты записи в лог-файл
- **Worker Pool** - для ограничения количества concurrent операций
- **Context** - для отмены длительных операций

### Преимущества архитектуры:
1. **Модульность** - четкое разделение ответственности
2. **Тестируемость** - каждый компонент можно тестировать отдельно
3. **Масштабируемость** - легко изменить количество workers
4. **Надежность** - comprehensive error handling
5. **Производительность** - асинхронная обработка

## Пример использования

### Входной XML:
```xml
<users>
    <user id="1">
        <name>Иван Иванов</name>
        <email>ivan@example.com</email>
        <age>30</age>
    </user>
    <user id="2">
        <name>Мария Петрова</name>
        <email>maria@example.com</email>
        <age>25</age>
    </user>
</users>
```

### Выходной JSON:
```json
[
    {
        "id": "1",
        "full_name": "Иван Иванов",
        "email": "ivan@example.com",
        "age_group": "от 25 до 35"
    },
    {
        "id": "2",
        "full_name": "Мария Петрова",
        "email": "maria@example.com",
        "age_group": "от 25 до 35"
    }
]
```

### Команда запуска:
```bash
# Установка зависимостей
go mod tidy

# Запуск сервера
go run main.go

# Тестовый запрос
curl -X POST http://localhost:8080/ \
  -H "Authorization: your-secret-key" \
  -H "Content-Type: application/xml" \
  -d @examples/sample_request.xml
```

Данная архитектура обеспечивает высокую производительность, надежность и легкость сопровожден