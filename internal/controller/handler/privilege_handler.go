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

// @Summary GetPrivilege
// @Tags Privilege
// @Description Get All Privileges
// @ID GetPrivilege
// @Accept  json
// @Produce  json
// @Param input body dto.PrivilegeDTO true "account info"
// @Success 200 {object} dto.PrivilegeResponseDTO
// @Failure 400,404
// @Failure 500
// @Failure default
// @Router /priv [get]
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

// @Summary CreatePrivilege
// @Tags Privilege
// @Description Create Privileges
// @ID CreatePrivilege
// @Accept  json
// @Produce  json
// @Param input body dto.PrivilegeDTO true "account info"
// @Success 201 {object} []byte
// @Failure 400,404
// @Failure 500
// @Failure default
// @Router /priv [post]
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
	rw.Write([]byte("Record has been created.\n"))
}

// @Summary AddPrivilegeToUser
// @Tags User
// @Description Add Privileges to User
// @ID AddPrivileges
// @Accept  json
// @Produce  json
// @Param input body dto.PrivilegedUserCreateDTO true "account info"
// @Success 201 {object} []byte
// @Failure 400,404
// @Failure 500
// @Failure default
// @Router /priv/user/add [post]
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

// @Summary RemovePrivilegeToUser
// @Tags User
// @Description Remove Privileges to User
// @ID RemovePrivileges
// @Accept  json
// @Produce  json
// @Param input body dto.PrivilegedUserDeleteDTO true "account info"
// @Success 200 {object} []byte
// @Failure 400,404
// @Failure 500
// @Failure default
// @Router /priv/user/remove [post]
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

// @Summary DeletePrivilege
// @Tags Privilege
// @Description Delete Privilege by id
// @ID DeletePrivilege
// @Produce  json
// @Success 200 {object} []byte
// @Failure 400,404
// @Failure 500
// @Failure default
// @Router /priv/:id [delete]
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
	rw.Write([]byte(fmt.Sprintf(`{"message": "Record has been deleted", "privilege_id": %d}`, id)))
}

// @Summary GetAllUsers
// @Tags User
// @Description Get All Users
// @ID GetUsers
// @Produce  json
// @Success 200 {object} []dto.PrivilegedUserDTO
// @Failure 400,404
// @Failure 500
// @Failure default
// @Router /priv/user [get]
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

// @Summary DeleteUser
// @Tags User
// @Description Delete User by id
// @ID DeleteUser
// @Produce  json
// @Success 200 {object} []byte
// @Failure 400,404
// @Failure 500
// @Failure default
// @Router /priv/user/:id [delete]
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
