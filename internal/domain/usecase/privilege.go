package usecase

import (
	"SergeyProject/internal/domain"
	entity2 "SergeyProject/internal/domain/entity"
	"SergeyProject/internal/dto"
	"SergeyProject/internal/errors"
	"SergeyProject/internal/utils"
	"SergeyProject/pkg/logger"
	"context"
	"github.com/hashicorp/go-hclog"
	"time"
)

type PrivilegeUsecase struct {
	logger hclog.Logger
	repo   domain.IPrivilegeRepository
}

func NewPrivilegeUsecase(repo domain.IPrivilegeRepository) *PrivilegeUsecase {
	return &PrivilegeUsecase{logger: logger.GetLogger(), repo: repo}
}

func (ps *PrivilegeUsecase) GetRecordByTitle(ctx context.Context, req *dto.PrivilegeDTO) (*dto.PrivilegeResponseDTO, error) {
	record, err := ps.repo.GetRecordByTitle(ctx, req.PrivilegeTitle)
	if err != nil {
		return nil, err
	}
	resp := &dto.PrivilegeResponseDTO{
		ID:             record.ID,
		PrivilegeTitle: record.PrivilegeTitle,
	}
	return resp, nil
}

func (ps *PrivilegeUsecase) CreatePrivilege(ctx context.Context, req *dto.PrivilegeDTO) error {
	validate := utils.NewValidator()
	if err := validate.Struct(req); err != nil {
		validation_error := utils.ValidatorErrors(err)
		for _, er := range validation_error {
			ps.logger.Error("Validation error", "error", er)
		}
		return err
	}
	entity := &entity2.Privilege{
		PrivilegeTitle: req.PrivilegeTitle,
		CreatedAt:      time.Now(),
	}
	_, err := ps.repo.GetRecordByTitle(ctx, req.PrivilegeTitle)
	if err == nil {
		return errors.ErrRecordAlreadyExists
	} else {
		if err != errors.ErrNoRecordFound {
			return err
		}
	}

	if err = ps.repo.CreatePrivilege(ctx, entity); err != nil {
		return err
	}

	return nil
}
