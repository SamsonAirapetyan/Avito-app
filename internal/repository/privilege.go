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

/*
GetRecordByTitle

Функция для обращения к БД для получения дынных о привилегии
*/
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

/*
GetRecordByID

Функция, которая обращается к БД и возвращает название привилегии по её id
*/
func (pr *PrivilegeRepository) GetRecordByID(ctx context.Context, priv_id int) (string, error) {
	query := `SELECT * FROM privileges WHERE id = $1 `
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	entity := &entity.Privilege{}
	conn, err := pr.storage.GetPgConnPool().Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	if err = tx.QueryRow(ctx, query, priv_id).Scan(&entity.ID, &entity.PrivilegeTitle, &entity.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return "", errors.ErrNoRecordFound
		}
		return "", err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}
	return entity.PrivilegeTitle, nil
}

/*
CreatePrivilege

функци которая обращается к БД для создания новой привилегии
*/
func (pr *PrivilegeRepository) CreatePrivilege(ctx context.Context, req *entity.Privilege) error {
	query := `INSERT INTO privileges(privilege_title, created_at) VALUES ($1,$2)`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := pr.storage.GetPgConnPool().Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, query, req.PrivilegeTitle, req.CreatedAt)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

/*
DeletePrivilege

Функция, которая обращается к БД для удаления привилегии.
*/
func (pr *PrivilegeRepository) DeletePrivilege(ctx context.Context, id int) error {
	query := `DELETE FROM privileges WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := pr.storage.GetPgConnPool().Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, query, id); err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

/*
GetUserPrivilegesByID

Функция, которая обращается к БД для получения всех имеющихся привилегий у пользователя.
*/
func (pr *PrivilegeRepository) GetUserPrivilegesByID(ctx context.Context, user_id int) ([]int, error) {
	query := `SELECT * FROM privileged_users WHERE user_id = $1`
	privileges := []int{}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := pr.storage.GetPgConnPool().Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		privilege := &entity.PrivilegedUser{}
		err = rows.Scan(&privilege.UserID, &privilege.PrivilegeID, &privilege.AssignedAt)
		if err != nil {
			return nil, err
		}
		privileges = append(privileges, privilege.PrivilegeID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return privileges, nil
}

/*
AddPrivilegeToUser

Функция, которая обращается к БД для добавления пользователю новой привилегии.
*/
func (pr *PrivilegeRepository) AddPrivilegeToUser(ctx context.Context, user_id int, privilege_id int) error {
	query := `INSERT INTO privileged_users (user_id, privilege_id, assigned_at) VALUES ($1,$2,$3)`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := pr.storage.GetPgConnPool().Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, query, user_id, privilege_id, time.Now())
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

/*
RemoveUserPrivilege

Функция, которая обращается к БД для удаления у пользователя привилегии
*/
func (pr *PrivilegeRepository) RemoveUserPrivilege(ctx context.Context, user_id int, priv_id int) error {
	query := `DELETE FROM privileged_users WHERE user_id = $1 AND privilege_id = $2`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := pr.storage.GetPgConnPool().Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query, user_id, priv_id)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

/*
GetAllUsers

Функция, которая обращается к БД для получения всех пользователей
*/
func (pr *PrivilegeRepository) GetAllUsers(ctx context.Context) ([]*entity.PrivilegedUser, error) {
	query := `SELECT * FROM privileged_users`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	entities := []*entity.PrivilegedUser{}
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
	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		priv := &entity.PrivilegedUser{}
		err = rows.Scan(&priv.UserID, &priv.PrivilegeID, &priv.AssignedAt)
		if err != nil {
			return nil, err
		}
		entities = append(entities, priv)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

/*
GetUserByID

Функция для проверки существования пользователя по получаемому на вход id
*/
func (pr *PrivilegeRepository) GetUserByID(ctx context.Context, user_id int) (int, error) {
	query := `SELECT * FROM privileged_users WHERE user_id = $1`
	entuty := &entity.PrivilegedUser{}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := pr.storage.GetPgConnPool().Acquire(ctx)
	if err != nil {
		return -1, err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback(ctx)

	if err := tx.QueryRow(ctx, query, user_id).Scan(&entuty.UserID, &entuty.PrivilegeID, &entuty.AssignedAt); err != nil {
		if err == pgx.ErrNoRows {
			return -1, errors.ErrNoRecordFound
		}
		return entuty.UserID, nil
	}
	err = tx.Commit(ctx)
	if err != nil {
		return -1, err
	}
	return entuty.UserID, nil
}

/*
DeletePrivilegeUser

Функция, которая обращается к БД для удаления пользователя по его id
*/
func (pr *PrivilegeRepository) DeletePrivilegeUser(ctx context.Context, user_id int) error {
	query := `DELETE FROM privileged_users WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := pr.storage.GetPgConnPool().Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, query, user_id)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
