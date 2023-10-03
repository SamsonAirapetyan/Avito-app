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

/*
handlePrivilegeGetByTitle

	Ручка для получения данных о привиллегии
*/
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

/*
handlePrivilegeCreate

	Ручка для создания привилегии
*/
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

/*
handlerAttachPrivilegeToUser

Ручка для добавления пользователю привилегий
*/
func (ph *PrivilegeHandler) handlerAttachPrivilegeToUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()
	req := &dto.PrivilegedUserCreateDTO{}
	if err := utils.StructDecode(r, req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	priv, err := ph.privilegeUsecases.AddPrivilegeToUser(ctx, req)
	if err != nil {
		if err == errors.ErrNoRecordFound {
			ph.logger.Error("No privilege record with such title exists", "error", err)
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		} else if err == errors.ErrRecordAlreadyExists {
			ph.logger.Error("Such privilege is already assigned to the user", "privilege", priv, "error", err)
			http.Error(rw, fmt.Sprintf("%s: %s", err.Error(), priv), http.StatusBadRequest)
			return
		}

		ph.logger.Error("Internal error", "message", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(`{"message": "Records have been created"}`))
}

/*
handlerRemovePrivilegeToUser

Ручка для удаления привилегии у пользователя
*/
func (ph *PrivilegeHandler) handlerRemovePrivilegeToUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()
	req := &dto.PrivilegedUserDeleteDTO{}
	if err := utils.StructDecode(r, req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	priv, err := ph.privilegeUsecases.RemoveUserPrivilege(ctx, req)
	if err != nil {
		if err == errors.ErrNoRecordFound {
			ph.logger.Error("No privilege record with such title attached to the user", "error", err)
			http.Error(rw, fmt.Sprintf("%s: %s", err.Error(), priv), http.StatusNotFound)
			return
		}

		ph.logger.Error("Internal error", "message", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{"message": "Privileges have been deleted"}`))
}

/*
handlerPrivilegeDelete

Ручка для удаления определенной привилегии
*/
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

/*
handlerGetAllUsers

Ручка для получения всех пользователей
*/
func (ph *PrivilegeHandler) handlerGetAllUsers(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()
	records, err := ph.privilegeUsecases.GetAllUsers(ctx)
	if err != nil {
		ph.logger.Error("Couldn't get records from privilege table", "error", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = utils.ToJSON(records, rw); err != nil {
		ph.logger.Error("JSON sezialisation didn't complete successfuly", "error", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

/*
handlerPrivilegeUserDelete

Ручка для удаления пользователя с привилегиями
*/
func (ph *PrivilegeHandler) handlerPrivilegeUserDelete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := ph.privilegeUsecases.DeletePrivilegedUser(ctx, id)
	if err != nil {
		if err == errors.ErrNoRecordFound {
			ph.logger.Error("No privileged user record with such id has been found", "error", err)
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		ph.logger.Error("Internal error", "message", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(fmt.Sprintf(`{"message": "Record has been deleted", "deleted privileged user id": %d}`, id)))

}
