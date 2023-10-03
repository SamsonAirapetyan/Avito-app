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

/*
GetRecordByTitle

Функция для получения данных о привилегии по её названию
*/
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

/*
GetRecordByID

Функция для получения наименования привиегии по её id
*/
func (ps *PrivilegeUsecase) GetRecordByID(ctx context.Context, priv_id int) (string, error) {
	privilege, err := ps.repo.GetRecordByID(ctx, priv_id)
	if err != nil {
		return "", err
	}
	return privilege, nil
}

/*
CreatePrivilege

Функция для создания новой привилегии. Сначала проводится валидация данных, затем в Бд добавляется новая привилегия
*/
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

/*
DeletePrivilege

Функция для удаления привилении
*/
func (ps *PrivilegeUsecase) DeletePrivilege(ctx context.Context, id int) error {
	_, err := ps.GetRecordByID(ctx, id)
	if err != nil {
		if err == errors.ErrNoRecordFound {
			return errors.ErrNoRecordFound
		}
		return err
	}
	if err = ps.repo.DeletePrivilege(ctx, id); err != nil {
		return err
	}

	return nil
}

/*
AddPrivilegeToUser на вход поступает структура ID user и название привилегий которые надо добавить
(1) Сответственно мы для начала проверяем валидность введенных данных
(2) Затем проходимся по каждой введеной привелегии чтобы определить из БД id такой привелегии
(3) Далее получаем данные всех имеющихся привилегий у пользователя с введенным id (получаем все id привилегий)
(4) Проходимся по каждому id привилегий чтобы из БД получить название этой привилегии
(5) Если такая привилегия уже ест у пользователя, то выдается ошибка
(6) Если все в порядке, то вызывается функция repository которая добовляет привилегию пользователю
*/
func (ps *PrivilegeUsecase) AddPrivilegeToUser(ctx context.Context, req *dto.PrivilegedUserCreateDTO) (string, error) {
	//(1)
	validator := utils.NewValidator()

	if err := validator.Struct(req); err != nil {
		validation_error := utils.ValidatorErrors(err)
		for _, er := range validation_error {
			ps.logger.Error("Validation error", "error", er)
		}
		return "", err
	}
	//(2)
	UserID := req.UserID
	for _, privilege := range req.PrivilegeList {
		//(2)
		entity, err := ps.repo.GetRecordByTitle(ctx, privilege)
		if err != nil {
			return "", err
		}
		//(3)
		privileged_ids, err := ps.repo.GetUserPrivilegesByID(ctx, req.UserID)
		if err != nil {
			return "", nil
		}
		//(4)
		for _, privilege_id := range privileged_ids {
			record_privilege, err := ps.repo.GetRecordByID(ctx, privilege_id)
			if err != nil {
				return "", nil
			}
			//(5)
			if record_privilege == privilege {
				return privilege, errors.ErrRecordAlreadyExists
			}
		}
		//(6)
		if err = ps.repo.AddPrivilegeToUser(ctx, UserID, entity.ID); err != nil {
			return "", err
		}
	}
	return "", nil
}

/*
GetAllUsers

Функция для получения всех пользователей
*/
func (ps *PrivilegeUsecase) GetAllUsers(ctx context.Context) ([]*dto.PrivilegedUserDTO, error) {
	records := []*dto.PrivilegedUserDTO{}

	entitys, err := ps.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	for _, entity := range entitys {
		priv_title, err := ps.repo.GetRecordByID(ctx, entity.PrivilegeID)
		if err != nil {
			return nil, err
		}
		record := &dto.PrivilegedUserDTO{UserID: entity.UserID, Privilege: priv_title}
		records = append(records, record)

	}
	return records, nil
}

/*
RemoveUserPrivilege

Функция для удаления привилегий у пользователя. Для начала проводится валидация данных.
Затем по каждому введенной привилегии ищется id привилегии. После чего определются привилегии принадлежащие данному пользователю.
Нужные привилегии удаляются
*/
func (ps *PrivilegeUsecase) RemoveUserPrivilege(ctx context.Context, req *dto.PrivilegedUserDeleteDTO) (string, error) {
	validator := utils.NewValidator()

	if err := validator.Struct(req); err != nil {
		validation_error := utils.ValidatorErrors(err)
		for _, er := range validation_error {
			ps.logger.Error("Validation error", "error", er)
		}
		return "", err
	}
	User_ID := req.UserID
	flag := false
	for _, privilege := range req.PrivilegeList {
		entity, err := ps.repo.GetRecordByTitle(ctx, privilege)
		if err != nil {
			return "", err
		}

		privilege_list, err := ps.repo.GetUserPrivilegesByID(ctx, User_ID)
		if err != nil {
			return "", err
		}
		for _, priv := range privilege_list {
			record_privilege, err := ps.repo.GetRecordByID(ctx, priv)
			if err != nil {
				return "", err
			}
			if record_privilege == privilege {
				flag = true
				break
			}
		}
		if flag {
			if err = ps.repo.RemoveUserPrivilege(ctx, User_ID, entity.ID); err != nil {
				return "", err
			} else {
				return privilege, errors.ErrNoRecordFound
			}
		}
	}
	return "", nil
}

/*
DeletePrivilegedUser
Удаление пользователя.
Проверка что пользователь существует, а затем уже вызов функции удаляющей пользователя.
*/
func (ps *PrivilegeUsecase) DeletePrivilegedUser(ctx context.Context, user_id int) error {
	_, err := ps.repo.GetUserByID(ctx, user_id)
	if err != nil {
		return err
	}

	if err = ps.repo.DeletePrivilegeUser(ctx, user_id); err != nil {
		return err
	}
	return nil

}
