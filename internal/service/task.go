package service

import (
	"context"
	"errors"
	"fmt"
	db "time-tracker/internal/db/sqlc"
	"time-tracker/internal/models"
	"time-tracker/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrTaskAlreadyExists = errors.New("task already exists")
	ErrTaskNotFound      = errors.New("task not found")
)

type TaskService struct {
	repository db.Querier
}

func NewTaskService(repository db.Querier) *TaskService {
	return &TaskService{
		repository: repository,
	}
}

func (ts *TaskService) CreateTask(ctx context.Context, userUUID uuid.UUID, payload *models.CreateTaskPayload) (*models.Task, error) {
	params := db.CreateTaskParams{
		UserUuid: pgtype.UUID{Bytes: userUUID, Valid: true},
		Name:     payload.Name,
	}

	taskRaw, err := ts.repository.CreateTask(ctx, params)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23503":
				return nil, ErrForeignKeyViolation
			case "23505":
				return nil, ErrTaskAlreadyExists
			}
		}
		return nil, err
	}

	task, err := utils.ConvertDBTaskToModelsTask(taskRaw)
	if err != nil {
		return nil, fmt.Errorf("error converting user: %v", err)
	}

	return task, nil
}

func (ts *TaskService) FinishTask(ctx context.Context, userUUID uuid.UUID) (*models.CompletedTask, error) {
	userPgUUID := pgtype.UUID{Bytes: userUUID, Valid: true}

	_, err := ts.repository.GetUserByUUID(ctx, userPgUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	taskRaw, err := ts.repository.UpdateTaskEndTime(ctx, userPgUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	params := db.CreateTaskHistoryParams{
		UserUuid:  userPgUUID,
		Name:      taskRaw.Name,
		StartTime: taskRaw.StartTime,
		EndTime:   taskRaw.EndTime,
	}

	taskHistoryRaw, err := ts.repository.CreateTaskHistory(ctx, params)
	if err != nil {
		return nil, err
	}

	err = ts.repository.DeleteTask(ctx, userPgUUID)
	if err != nil {
		return nil, err
	}

	return &models.CompletedTask{
		Name:     taskHistoryRaw.Name,
		Duration: taskHistoryRaw.Duration.(string),
	}, nil
}

func (ts *TaskService) GetTasksResult(ctx context.Context, userUUID uuid.UUID, days int) (*models.TasksResult, error) {
	userPgUUID := pgtype.UUID{Bytes: userUUID, Valid: true}
	_, err := ts.repository.GetUserByUUID(ctx, userPgUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	params := db.GetTasksResultByPeriodParams{
		UserUuid: userPgUUID,
		Column1:  pgtype.Interval{Days: int32(days), Valid: true},
	}

	taskResultByPeriodRows, err := ts.repository.GetTasksResultByPeriod(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(taskResultByPeriodRows) == 0 {
		return nil, ErrTaskNotFound
	}

	var completedTasks = make([]models.CompletedTask, len(taskResultByPeriodRows))
	for i, task := range taskResultByPeriodRows {
		completedTasks[i] = models.CompletedTask{
			Name:     task.TaskName,
			Duration: task.Duration.(string),
		}
	}

	return &models.TasksResult{
		CompletedTask: completedTasks,
		TotalDuration: taskResultByPeriodRows[0].TotalDuration.(string),
	}, nil
}
