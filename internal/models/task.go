package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateTaskPayload struct {
	Name string `json:"name"`
}

type Task struct {
	UUID      uuid.UUID  `json:"uuid"`
	UserUUID  uuid.UUID  `json:"userUuid"`
	Name      string     `json:"name"`
	StartTime time.Time  `json:"startTime"`
	EndTime   *time.Time `json:"endTime,omitempty"`
}
