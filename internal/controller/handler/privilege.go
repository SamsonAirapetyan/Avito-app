package handler

import (
	"SergeyProject/internal/controller"
	"SergeyProject/internal/domain"
	"SergeyProject/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"net/http"
)

type PrivilegeHandler struct {
	logger            hclog.Logger
	privilegeUsecases domain.IPrivilegeUsecases
}

func NewPrivilegeHandler(privilegeUsecases domain.IPrivilegeUsecases) controller.IHandler {
	return &PrivilegeHandler{logger: logger.GetLogger(), privilegeUsecases: privilegeUsecases}
}

func (ph *PrivilegeHandler) Register(router *mux.Router) {
	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/priv", ph.handlePrivilegeGetByTitle)
	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/priv", ph.handlePrivilegeCreate)
}
