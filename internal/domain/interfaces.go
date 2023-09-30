package domain

import (
	"SergeyProject/internal/domain/entity"
	"SergeyProject/internal/dto"
	"context"
)

type (
	IPrivilegeUsecases interface {
		GetRecordByTitle(context.Context, *dto.PrivilegeDTO) (*dto.PrivilegeResponseDTO, error)
		CreatePrivilege(context.Context, *dto.PrivilegeDTO) error
	}

	IPrivilegeRepository interface {
		GetRecordByTitle(context.Context, string) (*entity.Privilege, error)
		CreatePrivilege(context.Context, *entity.Privilege) error
	}
)
