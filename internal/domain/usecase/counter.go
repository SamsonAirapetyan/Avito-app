package usecase

import (
	"SergeyProject/internal/domain"
	"SergeyProject/pkg/logger"
	"github.com/hashicorp/go-hclog"
	"math"
)

type CounterUsecase struct {
	logger hclog.Logger
	repo   domain.ICounterRepository
}

func NewCounterUsecase(repo domain.ICounterUsecases) *CounterUsecase {
	return &CounterUsecase{logger: logger.GetLogger(), repo: repo}
}

func (cu *CounterUsecase) SetValue(name string, val int) int {
	return cu.repo.SetValue(name, val)
}
func (cu *CounterUsecase) IncreaseValue(name string, val int) int {
	storage := cu.repo.GetStorage()
	if _, ok := storage[name]; ok {
		if storage[name]+val <= math.MaxInt {
			return cu.repo.IncreaseValue(name, val)
		} else {
			return -1
		}
	} else {
		storage[name] = 0
		return -2
	}
}

func (cu *CounterUsecase) DecreaseValue(name string, val int) int {
	storage := cu.repo.GetStorage()
	if _, ok := storage[name]; ok {
		if storage[name]-val >= 0 {
			return cu.repo.DecreaseValue(name, val)
		}
	}
	return -1
}

func (cu *CounterUsecase) GetStorage() map[string]int {
	return cu.repo.GetStorage()
}
