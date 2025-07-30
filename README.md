# GoXML_JSON - Сервис обработки XML в JSON

Микросервис для преобразования XML данных пользователей в JSON формат с отправкой на внешний сервер.

## 🏗️ Архитектура проекта

```
GoXML_JSON/
├── cmd/                   # Точки входа приложений
│   ├── server/            # Основной сервер (порт 8080)
│   │   └── main.go
│   └── test_server/       # Тестовый сервер (порт 8081)
│       └── main.go
├── internal/              # Внутренняя логика приложения
│   ├── client/            # HTTP клиент для отправки данных
│   │   ├── client.go      # Клиент встроен в хэндлер сервера 8080 и обращается к серверу 8081
│   │   └── send_users.go  # Метод клиента для отправки JSON юзеров
│   ├── converter/         # Конвертация данных
│   │   ├── converter.go   # Структура реализующая методы
│   │   ├── xml_parser.go  # Парси XML
│   │   ├── json_converter.go # Конвертирует записи асинхронно
│   │   └── age_groups.go  # Определяет группу пользователей
│   ├── handler/           # HTTP обработчики
│   │   ├── handler.go
│   │   └── postUsers.go
│   ├── middleware/       # Промежуточное ПО
│   │   └── auth.go       # Для аутентификации пользователя по ключу
│   └── models/           # Модели данных
│       ├── user.go
│       └── errors.go
├── pkg/                 # Переиспользуемые пакеты
│   └── logger/          # Логирование
│       └── logger.go
├── settings/            # Конфигурация
│   └── settings.go
├── test_users.xml       # Тестовые данные
├── provider.log         # Лог файл
├── go.mod               # Зависимости Go
└── README.md            # Документация
```

## 🎯 Основные компоненты

### 📡 Основной сервер (порт 8080)
- **Функция**: Принимает XML данные, конвертирует в JSON, отправляет на внешний сервер
- **Эндпоинты**:
  - `POST /users` - обработка XML пользователей
  - `GET /health` - проверка состояния

### 🔧 Тестовый сервер (порт 8081)
- **Функция**: Принимает JSON данные и возвращает подтверждение
- **Эндпоинты**:
  - `POST /users` - обработка JSON пользователей
  - `GET /health` - проверка состояния

### 🔐 Авторизация
- **Тип**: Bearer Token
- **Ключ**: `1234567890`
- **Header**: `Authorization: Bearer 1234567890`

### 📊 Обработка данных
- **XML → JSON**: Конвертация с возрастными группами
- **Возрастные группы**:
  - `до 25` - молодые
  - `от 25 до 35` - средние
  - `старше 35` - старшие
- **Асинхронность**: Обработка пользователей в goroutines

## 🚀 Запуск приложения

### 1. Запуск основного сервера
```bash
go run ./cmd/server/main.go
```
Сервер запустится на `http://localhost:8080`

### 2. Запуск тестового сервера
```bash
go run ./cmd/test_server/main.go
```
Сервер запустится на `http://localhost:8081`

### 3. Проверка состояния серверов
```bash
# Основной сервер
curl http://localhost:8080/health

# Тестовый сервер
curl http://localhost:8081/health
```

## 🧪 Тестирование с помощью curl

### Тестовые данные (test_users.xml)
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

### 1. Полный цикл обработки данных
```bash
curl -v -X POST http://localhost:8080/users \
  -H "Content-Type: application/xml" \
  -H "Authorization: Bearer 1234567890" \
  -d @test_users.xml
```

**Ожидаемый ответ:**
```json
{
  "data": {
    "status": "success",
    "message": "Успешно получено и обработано 2 пользователей",
    "received_at": "2025-07-30 19:55:57",
    "processed_by": "Test Server 8081",
    "user_count": 2,
    "users": [
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
  },
  "usersProcessed": 2
}
```

### 2. Тестирование авторизации
```bash
# Без токена (должен вернуть 401)
curl -v -X POST http://localhost:8080/users \
  -H "Content-Type: application/xml" \
  -d @test_users.xml

# Неправильный токен (должен вернуть 401)
curl -v -X POST http://localhost:8080/users \
  -H "Content-Type: application/xml" \
  -H "Authorization: Bearer wrong_token" \
  -d @test_users.xml
```

### 3. Тестирование обработки ошибок
```bash
# Некорректный XML (должен вернуть 400)
curl -v -X POST http://localhost:8080/users \
  -H "Content-Type: application/xml" \
  -H "Authorization: Bearer 1234567890" \
  -d "некорректный XML"

# Пустое тело запроса (должен вернуть 400)
curl -v -X POST http://localhost:8080/users \
  -H "Content-Type: application/xml" \
  -H "Authorization: Bearer 1234567890"
```

### 4. Прямое тестирование тестового сервера
```bash
curl -v -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '[{"id":"1","full_name":"Тест","email":"test@example.com","age_group":"от 25 до 35"}]'
```

## 📋 Логирование

Все операции логируются в файл `provider.log`:

```
2025/07/30 19:55:57 ✅ авторизованный доступ: remote=127.0.0.1:46366, path=/users
2025/07/30 19:55:57 🙏 Users: начало обработки запроса
2025/07/30 19:55:57 ✅ Users: тело запроса успешно прочитано, размер: 292 байт
2025/07/30 19:55:57 ✅ Users: XML успешно пропарсен: &{{ users} [{1 Иван Иванов ivan@example.com 30} {2 Мария Петрова maria@example.com 25}]}
2025/07/30 19:55:57 ✅ Сконвертировано 2 валидных пользователей. Начинаем отправку...
2025/07/30 19:55:57 🙏 Users: Отправляем пользователей на сервер
2025/07/30 19:55:57 ✅ Пользователи успешно отправлены на сервер
```

## 🔧 Конфигурация

Основные настройки в `settings/settings.go`:

```go
const (
    AuthKey = "1234567890"           // Ключ авторизации
    OutputFile = "provider.log"       // Файл логов
    ClientTimeout = 10 * time.Second  // Таймаут HTTP запросов
    ServerHost = "localhost"          // Хост сервера
    ServerPort = "8080"              // Порт сервера
    ClientURL = "http://localhost:8081/users"  // URL внешнего сервера
)
```

## 🏛️ Архитектурные принципы

### Clean Architecture
- **Handlers**: HTTP обработчики (транспортный слой)
- **Services**: Бизнес-логика (конвертация, валидация)
- **Repositories**: Доступ к данным (HTTP клиент)
- **Models**: Структуры данных

### Dependency Injection
- Все зависимости инжектируются через конструкторы
- Интерфейсы для тестируемости
- Разделение ответственности

### Error Handling
- Кастомные типы ошибок
- Детальное логирование
- Graceful degradation

### Concurrency
- Асинхронная обработка пользователей
- Goroutines для параллельной работы
- Thread-safe логирование

## 🧪 Тестирование

### Unit тесты
```bash
# Запуск всех тестов
go test ./...

# Запуск с покрытием
go test -cover ./...

# Запуск бенчмарков
go test -bench=. ./...
```

### Интеграционные тесты
- Тестирование полного цикла обработки
- Проверка взаимодействия между серверами
- Тестирование обработки ошибок

## 📊 Мониторинг

### Health Checks
```bash
# Основной сервер
curl http://localhost:8080/health

# Тестовый сервер
curl http://localhost:8081/health
```

### Логи
- Все операции логируются в `provider.log`
- Структурированные логи с временными метками
- Детальная информация об ошибках

## 🚀 Производительность

- **Асинхронная обработка**: Пользователи обрабатываются параллельно
- **Оптимизированная память**: Минимальные аллокации
- **Эффективная сеть**: Переиспользование HTTP соединений
- **Быстрая сериализация**: Нативная поддержка XML/JSON

## 🔒 Безопасность

- **Авторизация**: Bearer token аутентификация
- **Валидация**: Проверка входных данных
- **Санитизация**: Очистка от лишних пробелов
- **Логирование**: Аудит всех операций
