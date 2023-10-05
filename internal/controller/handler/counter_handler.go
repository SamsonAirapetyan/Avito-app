package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (ch *CounterHandler) handlerSetCounter(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	val, _ := strconv.Atoi(vars["val"])
	name := r.URL.Query().Get("name")
	_ = ch.counterUsecase.SetValue(name, val)
	rw.WriteHeader(http.StatusOK)
}

func (ch *CounterHandler) handlerIncreaseCounter(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	val, _ := strconv.Atoi(vars["val"])
	name := r.URL.Query().Get("name")
	currentValue := ch.counterUsecase.IncreaseValue(name, val)
	if currentValue == -1 {
		ch.logger.Warn("MaxInt ceiling has been hit")
		return
	} else if currentValue == -2 {
		ch.logger.Warn("Default value is set to zero")
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (ch *CounterHandler) handlerDecreaseCounter(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	val, _ := strconv.Atoi(vars["val"])
	name := r.URL.Query().Get("name")

	currentValue := ch.counterUsecase.DecreaseValue(name, val)
	if currentValue == -1 {
		ch.logger.Warn("Current value is zero, set the new one to continue decrementing")
		return
	}

	rw.WriteHeader(http.StatusOK)
}
