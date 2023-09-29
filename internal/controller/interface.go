package controller

import "github.com/gorilla/mux"

type IHandler interface {
	Register(router *mux.Router)
}
