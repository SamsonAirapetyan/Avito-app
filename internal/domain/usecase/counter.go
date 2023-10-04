package usecase

import (
	"SergeyProject/internal/domain"
	"SergeyProject/pkg/logger"
	"github.com/hashicorp/go-hclog"
)

type CounterUsecase struct {
	logger hclog.Logger
	repo   domain.ICounterUsecases
}

func NewCounterUsecase(repo domain.ICounterUsecases) *CounterUsecase {
	return &CounterUsecase{logger: logger.GetLogger(), repo: repo}
}
