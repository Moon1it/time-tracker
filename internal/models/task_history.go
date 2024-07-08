package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskHistory struct {
	Uuid      uuid.UUID `json:"uuid"`
	TaskUuid  uuid.UUID `json:"taskUuid"`
	UserUuid  uuid.UUID `json:"userUuid"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endtime"`
}

type CompletedTask struct {
	Name     string `json:"name"`
	Duration string `json:"duration"`
}

type TasksResult struct {
	TotalDuration string          `json:"totalDuration"`
	CompletedTask []CompletedTask `json:"CompletedTask"`
}
