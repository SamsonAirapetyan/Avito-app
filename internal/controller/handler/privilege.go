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
	getRouter.HandleFunc("/priv", ph.handlePrivilegeGetByTitle) //+
	getRouter.HandleFunc("/priv/user", ph.handlerGetAllUsers)   //+

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/priv", ph.handlePrivilegeCreate)                    //+
	postRouter.HandleFunc("/priv/user/add", ph.handlerAttachPrivilegeToUser)    //+
	postRouter.HandleFunc("/priv/user/remove", ph.handlerRemovePrivilegeToUser) //-

	deleteRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/priv/{id:[0-9]+}", ph.handlerPrivilegeDelete)          //+
	deleteRouter.HandleFunc("/priv/user/{id:[0-9]+}", ph.handlerPrivilegeUserDelete) //+
}
