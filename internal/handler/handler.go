package handler

import (
	"github.com/NarthurN/GoXML_JSON/internal/client"
	"github.com/NarthurN/GoXML_JSON/internal/converter"
	"github.com/NarthurN/GoXML_JSON/pkg/logger"
)

type Handler struct {
	logger    *logger.Logger
	converter *converter.Converter
	client    *client.Client
}

func NewHandler(logger *logger.Logger) *Handler {
	return &Handler{
		logger:    logger,
		converter: converter.NewConverter(logger),
		client:    client.NewClient(logger),
	}
}
