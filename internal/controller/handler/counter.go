package handler

import (
	"SergeyProject/internal/domain"
	"SergeyProject/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type CounterHandler struct {
	logger         hclog.Logger
	counterUsecase domain.ICounterUsecases
}

func NewCounterHandler(counterUsecase domain.ICounterUsecases) *CounterHandler {
	return &CounterHandler{logger: logger.GetLogger(), counterUsecase: counterUsecase}
}

func (ch *CounterHandler) Register(router *mux.Router) {
}
