package utils

import (
	"time"
	db "time-tracker/internal/db/sqlc"
	"time-tracker/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertDBUserToModelsUser(user db.User) (*models.User, error) {
	var uuid uuid.UUID
	err := uuid.UnmarshalBinary(user.Uuid.Bytes[:])
	if err != nil {
		return nil, err
	}

	var patronymic *string
	if user.Patronymic.Valid {
		patronymic = &user.Patronymic.String
	}

	return &models.User{
		UUID:           uuid,
		PassportNumber: user.PassportNumber,
		Surname:        user.Surname,
		Name:           user.Name,
		Patronymic:     patronymic,
		Address:        user.Address,
		CreatedAt:      user.CreatedAt.Time,
		UpdatedAt:      user.UpdatedAt.Time,
	}, nil
}

func ConvertDBTaskToModelsTask(dbTask db.Task) (*models.Task, error) {
	var taskUUID uuid.UUID
	err := taskUUID.UnmarshalBinary(dbTask.Uuid.Bytes[:])
	if err != nil {
		return nil, err
	}

	var userUUID uuid.UUID
	err = userUUID.UnmarshalBinary(dbTask.UserUuid.Bytes[:])
	if err != nil {
		return nil, err
	}

	var endTime *time.Time
	if dbTask.EndTime.Valid {
		endTime = &dbTask.EndTime.Time

	}

	return &models.Task{
		UUID:      taskUUID,
		UserUUID:  userUUID,
		Name:      dbTask.Name,
		StartTime: dbTask.StartTime.Time,
		EndTime:   endTime,
	}, nil
}

func ToPgText(s *string) pgtype.Text {
	if s != nil {
		return pgtype.Text{String: *s, Valid: true}
	}
	return pgtype.Text{Valid: false}
}

func ConvertDBTaskHistoryToModelsTaskHistory(dbTask db.TaskHistory) (*models.TaskHistory, error) {
	var modelsTask models.TaskHistory

	if dbTask.Uuid.Valid {
		modelsTask.Uuid, _ = uuid.FromBytes(dbTask.Uuid.Bytes[:])
	}
	if dbTask.UserUuid.Valid {
		modelsTask.UserUuid, _ = uuid.FromBytes(dbTask.UserUuid.Bytes[:])
	}
	if dbTask.StartTime.Valid {
		modelsTask.StartTime = dbTask.StartTime.Time
	}
	if dbTask.EndTime.Valid {
		modelsTask.EndTime = dbTask.EndTime.Time
	}

	modelsTask.Name = dbTask.Name

	return &modelsTask, nil
}
