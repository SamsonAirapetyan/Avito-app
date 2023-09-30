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
	}

	IPrivilegeRepository interface {
		GetRecordByTitle(context.Context, string) (*entity.Privilege, error)
		GetRecordByID(context.Context, int) (string, error)
		CreatePrivilege(context.Context, *entity.Privilege) error
		DeletePrivilege(context.Context, int) error
	}
)
