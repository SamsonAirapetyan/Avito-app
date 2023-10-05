package handler

import (
	"SergeyProject/internal/domain"
	"SergeyProject/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"net/http"
)

type CounterHandler struct {
	logger         hclog.Logger
	counterUsecase domain.ICounterUsecases
}

func NewCounterHandler(counterUsecase domain.ICounterUsecases) *CounterHandler {
	return &CounterHandler{logger: logger.GetLogger(), counterUsecase: counterUsecase}
}

func (ch *CounterHandler) Register(router *mux.Router) {
	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/set/{val:[0-9]+}", ch.handlerSetCounter).Queries("name", "{[a-z]+}")
	getRouter.HandleFunc("/set/{val:[0-9]+}", ch.handlerSetCounter)
	getRouter.HandleFunc("/inc/{val:[0-9]+}", ch.handlerIncreaseCounter).Queries("name", "{[a-z]+}")
}
