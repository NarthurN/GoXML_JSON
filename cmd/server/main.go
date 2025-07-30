package main

import (
	"log"

	"github.com/NarthurN/GoXML_JSON/pkg/logger"
)

func main() {
	logg, err := logger.New()
	if err != nil {
		log.Fatalf("❌ не удалось создать логгер: %v", err)
	}
	defer logg.Close()

	logg.Log("✅ логер инциализирован ")


}
