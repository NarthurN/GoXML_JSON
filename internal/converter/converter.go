package converter

import (
	"github.com/NarthurN/GoXML_JSON/pkg/logger"
)

type Converter struct {
	logger *logger.Logger
}

func NewConverter(logger *logger.Logger) *Converter {
	return &Converter{
		logger: logger,
	}
}
