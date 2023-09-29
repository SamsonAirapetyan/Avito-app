package repository

import (
	"SergeyProject/internal/domain/entity"
	"SergeyProject/internal/errors"
	"SergeyProject/pkg/logger"
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v5"
	"time"
)

type PrivilegeRepository struct {
	logger  hclog.Logger
	storage *Storage //poolConnection
}

func NewPrivilegeRepository(storage *Storage) *PrivilegeRepository {
	return &PrivilegeRepository{logger: logger.GetLogger(), storage: storage}
}

func (pr *PrivilegeRepository) GetRecordByTitle(ctx context.Context, title string) (*entity.Privilege, error) {
	query := `SELECT * FROM privileges WHERE privilege_title = $1`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	entity := &entity.Privilege{}
	conn, err := pr.storage.GetPgConnPool().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err = tx.QueryRow(ctx, query, title).Scan(&entity.ID, &entity.PrivilegeTitle, &entity.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrNoRecordFound
		}
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return entity, nil
}
