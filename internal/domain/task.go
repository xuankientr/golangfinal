package domain

import (
	"fmt"
)

// Task model
type Task struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type ValidateTask struct {
	Title string
	Done  bool
}

func (t *Task) Validate() error {
	if len(t.Title) == 0 {
		return fmt.Errorf("title is required")
	}

	if len(t.Title) < 3 {
		return fmt.Errorf("title must be at least 3 characters long")
	}
	return nil
}
