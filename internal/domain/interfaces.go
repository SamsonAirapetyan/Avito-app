package domain

import (
	"SergeyProject/internal/domain/entity"
	"SergeyProject/internal/dto"
	"context"
)

type (
	IPrivilegeUsecases interface {
		GetRecordByTitle(context.Context, *dto.PrivilegeDTO) (*dto.PrivilegeResponseDTO, error)
	}

	IPrivilegeRepository interface {
		GetRecordByTitle(context.Context, string) (*entity.Privilege, error)
	}
)
