package domain

import (
	"SergeyProject/internal/domain/entity"
	"SergeyProject/internal/dto"
	"context"
)

type (
	IPrivilegeUsecases interface {
		GetRecordByTitle(context.Context, *dto.PrivilegeDTO) (*dto.PrivilegeResponseDTO, error)
		GetRecordByID(context.Context, int) (string, error)
		CreatePrivilege(context.Context, *dto.PrivilegeDTO) error
		DeletePrivilege(context.Context, int) error

		AddPrivilegeToUser(context.Context, *dto.PrivilegedUserCreateDTO) (string, error)
		GetAllUsers(context.Context) ([]*dto.PrivilegedUserDTO, error)
		RemoveUserPrivilege(context.Context, *dto.PrivilegedUserDeleteDTO) (string, error)
		DeletePrivilegedUser(context.Context, int) error
	}

	IPrivilegeRepository interface {
		GetRecordByTitle(context.Context, string) (*entity.Privilege, error)
		GetRecordByID(context.Context, int) (string, error)
		CreatePrivilege(context.Context, *entity.Privilege) error
		DeletePrivilege(context.Context, int) error

		GetUserPrivilegesByID(context.Context, int) ([]int, error)
		AddPrivilegeToUser(context.Context, int, int) error
		GetAllUsers(context.Context) ([]*entity.PrivilegedUser, error)
		GetUserByID(context.Context, int) (int, error)
		RemoveUserPrivilege(context.Context, int, int) error
		DeletePrivilegeUser(context.Context, int) error
	}

	ICounterUsecases interface {
		SetValue(string, int) int
		IncreaseValue(string, int) int
		DecreaseValue(string, int) int
		GetStorage() map[string]int
	}

	ICounterRepository interface {
		GetStorage() map[string]int
		SetValue(string, int) int
		IncreaseValue(string, int) int
		DecreaseValue(string, int) int
	}
)
