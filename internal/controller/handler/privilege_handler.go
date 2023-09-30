package handler

import (
	"SergeyProject/internal/dto"
	"SergeyProject/internal/errors"
	"SergeyProject/internal/utils"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (ph *PrivilegeHandler) handlePrivilegeGetByTitle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()

	req := &dto.PrivilegeDTO{}
	if err := utils.StructDecode(r, req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := ph.privilegeUsecases.GetRecordByTitle(ctx, req)
	if err != nil {
		if err == errors.ErrNoRecordFound {
			ph.logger.Error("No privilege record has been found", "filter title", req.PrivilegeTitle)
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}
		ph.logger.Error("Internal error", "message", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = utils.ToJSON(record, rw); err != nil {
		ph.logger.Error("JSON sezialisation didn't complete successfuly", "error", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (ph *PrivilegeHandler) handlePrivilegeCreate(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()
	req := &dto.PrivilegeDTO{}
	if err := utils.StructDecode(r, req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	err := ph.privilegeUsecases.CreatePrivilege(ctx, req)
	if err != nil {
		if err == errors.ErrRecordAlreadyExists {
			ph.logger.Error("Cannot create a record because record with such name already exists", "error", err)
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		ph.logger.Error("Internal error", "error", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte("REcord has been created.\n"))
}

func (ph *PrivilegeHandler) handlerPrivilegeDelete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := ph.privilegeUsecases.DeletePrivilege(ctx, id)
	if err != nil {
		if err == errors.ErrNoRecordFound {
			ph.logger.Error("No such privilege record has been found")
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		ph.logger.Error("Internal error", "message", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(fmt.Sprintf(`{"message": "Record has been deleted", "privilege_id": %d}}.`, id)))
}
