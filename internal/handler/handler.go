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

func NewHandler(logger *logger.Logger, converter *converter.Converter, client *client.Client) *Handler {
	return &Handler{
		logger:    logger,
		converter: converter,
		client:    client,
	}
}
