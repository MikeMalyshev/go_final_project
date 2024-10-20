package db

import (
	"database/sql"
	"fmt"
	"strings"

	"go_final_project/internal/app"
)

func (storage *DBStorage) AddTask(task app.Task) (int64, error) {
	res, err := storage.db.Exec(
		`
		INSERT
			INTO scheduler
			(date, title, comment, repeat)
			VALUES (:date, :title, :comment, :repeat)
		`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if num, _ := res.RowsAffected(); num != 1 {
		return 0, fmt.Errorf("DBStorage.AddTask: task already exists")
	}

	if err != nil {
		return 0, fmt.Errorf("DBStorage.AddTask: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return id, fmt.Errorf("DBStorage.AddTask: %v", err)
	}
	return id, nil
}

func (storage *DBStorage) UpdateTask(task app.Task) error {
	res, err := storage.db.Exec(
		`
		UPDATE scheduler
			SET date = :date, title = :title, comment = :comment, repeat = :repeat
			WHERE id = :id
		`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))

	r, err := res.RowsAffected()
	if r == 0 {
		return fmt.Errorf("DBStorage.UpdateTask: coudn't find task.id=%s", task.ID)
	} else if r > 1 {
		panic("DBStorage.UpdateTask number of rowsAffected > 1 !!!")
	}

	if err != nil {
		return fmt.Errorf("DBStorage.UpdateTask: %v", err)
	}

	return nil
}

func (storage *DBStorage) RemoveTask(id string) error {
	_, err := storage.db.Exec(
		`
		DELETE
			FROM scheduler
			WHERE id = :id
		`,
		sql.Named("id", id))

	if err != nil {
		return fmt.Errorf("DBStorage.RemoveTask: %v", err)
	}

	return nil
}

func (storage *DBStorage) GetTaskByID(id string) (app.Task, error) {
	row := storage.db.QueryRow(
		`
		SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE id = :id
		`,
		sql.Named("id", id))

	var task app.Task

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return app.Task{}, fmt.Errorf("DBStorage.GetTask: %v", err)
	}
	return task, nil
}

func (storage *DBStorage) GetTaskList(searchString string, maxLen int64) ([]app.Task, error) {
	var tasks []app.Task
	dateString := searchString
	caseInsensitiveRegExpr := ""
	for _, c := range searchString {
		caseInsensitiveRegExpr += "[" + strings.ToUpper(string(c)) + strings.ToLower(string(c)) + "]"
	}

	rows, err := storage.db.Query(
		`
		SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE
				title REGEXP :search OR
				comment REGEXP :search OR
				date = :date
			ORDER BY date
			LIMIT :limit
		`,
		sql.Named("search", caseInsensitiveRegExpr),
		sql.Named("date", dateString),
		sql.Named("limit", maxLen))

	defer rows.Close()

	if err != nil {
		return []app.Task{}, fmt.Errorf("DBStorage.AddTask: %v", err)
	}

	for rows.Next() {
		task := app.Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return []app.Task{}, fmt.Errorf("DBStorage.AddTask: %v", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []app.Task{}, fmt.Errorf("DBStorage.AddTask: %v", err)
	}

	return tasks, nil
}

func (storage *DBStorage) FindTask(title, date string) (string, error) {
	var task app.Task

	db, err := sql.Open("sqlite3_ext", storage.cfg.DBPath())
	if err != nil {
		return "", fmt.Errorf("DBStorage.GetTasks: %v", err)
	}
	defer db.Close()

	row := db.QueryRow(
		`
		SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE
				title = :title AND
				date = :date
		`,
		sql.Named("title", title),
		sql.Named("date", date),
	)
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err == sql.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("DBStorage.AddTask: %v", err)
	}
	return task.ID, nil
}
