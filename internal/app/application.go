package app

import (
	"fmt"
	"strconv"
	"time"

	"go_final_project/internal/config"
)

type Application struct {
	storage Storage
}

func CreateApplication(storage Storage) *Application {
	return &Application{storage: storage}
}

// Принимает текущее время now, предыдущую установленную дату задачи date, правило повторения repeat и возвращает новую дату
func (app Application) NextDate(now, date, repeat string) (string, error) {
	t, err := time.Parse(config.DBDateFormat, now)
	if err != nil {
		return "", fmt.Errorf("Application.NextDate: %w", err)
	}

	next, err := NextDate(t, date, repeat)
	if err != nil {
		return next, fmt.Errorf("Application.NextDate: %w", err)
	}
	return next, nil
}

// Проверка полученной задачи на соответствие всем необходимым условиям
func (app Application) CheckTask(task Task) (Task, error) {
	if len(task.Title) == 0 {
		return task, fmt.Errorf("Application.CheckTask: Error, Task.Title is empty ")
	}

	now := time.Now()

	if len(task.Date) == 0 {
		task.Date = now.Format(config.DBDateFormat)
	}

	date, err := time.Parse(config.DBDateFormat, task.Date)
	if err != nil {
		return task, fmt.Errorf("Application.CheckTask: Task.Date has Invalid format ")
	}

	var nextDate string
	if len(task.Repeat) > 0 {
		nextDate, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return task, fmt.Errorf("Application.CheckTask: %w ", err)
		}
	}

	if !date.After(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(config.DBDateFormat)
		} else {
			task.Date = nextDate
		}
	}
	return task, nil
}

// Добавляет новую задачу и возвращает её id
func (app Application) AddTask(task Task) (int64, error) {
	task, err := app.CheckTask(task)
	if err != nil {
		return -1, err
	}

	return app.storage.AddTask(task)
}

// Возвращает задачу по её id
func (app Application) GetTask(id string) (Task, error) {
	if len(id) == 0 {
		return Task{}, fmt.Errorf("Application.GetTask: id not setted")
	}
	_, err := strconv.Atoi(id)
	if err != nil {
		return Task{}, fmt.Errorf("Application.GetTask: invalid id : %v", err)
	}
	return app.storage.GetTaskByID(id)
}

// Меняет содержимое задачи по id, указанному в переданной структуре
func (app Application) UpdateTask(task Task) error {
	_, err := strconv.Atoi(task.ID)
	if err != nil {
		return fmt.Errorf("Application.UpdateTask : invalid task.ID=%s ", task.ID)
	}

	task, err = app.CheckTask(task)
	if err != nil {
		return err
	}
	return app.storage.UpdateTask(task)
}

// Удаляет задачу по id
func (app Application) RemoveTask(id string) error {
	_, err := app.storage.GetTaskByID(id)
	if err != nil {
		return fmt.Errorf("Application.TaskDone : %v", err)
	}

	return app.storage.RemoveTask(id)
}

// Отмечает задачу по её id как завершенную (удаляет при отсутствии правила повторения или переносит при наличии такого правила
func (app Application) FinishTask(id string) error {
	task, err := app.storage.GetTaskByID(id)
	if err != nil {
		return fmt.Errorf("Application.FinishTask : %v", err)
	}

	if len(task.Repeat) == 0 {
		app.storage.RemoveTask(id)
		return nil
	}

	task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return fmt.Errorf("Application.FinishTask : %v", err)
	}

	err = app.storage.UpdateTask(task)
	if err != nil {
		return fmt.Errorf("Application.FinishTask : %v", err)
	}
	return nil
}

// Возвращает слайс задач максимальной длиной maxLen, удовлетворяющих по названию, комментарию или дате фильтру searchString.
func (app Application) GetTaskList(searchString string, maxLen int64) ([]Task, error) {
	date, err := time.Parse(config.WebDateFormat, searchString)
	if err == nil {
		searchString = date.Format(config.DBDateFormat)
	}

	tasks, err := app.storage.GetTaskList(searchString, maxLen)
	if err != nil {
		return nil, fmt.Errorf("Application.GetTaskList: %v", err)
	}

	if len(tasks) == 0 {
		return make([]Task, 0), nil
	}
	return tasks, nil
}
