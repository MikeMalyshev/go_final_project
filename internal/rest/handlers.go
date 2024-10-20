package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go_final_project/internal/app"
	"go_final_project/internal/config"
)

// Хэндлер обращений к `/api/task`
func (mux Mux) TaskHandler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		mux.TaskPostHandler(resp, req)

	case http.MethodDelete:
		mux.TaskDeleteHandler(resp, req)

	case http.MethodPut:
		mux.TaskPutHandler(resp, req)

	case http.MethodGet:
		mux.TaskGetHandler(resp, req)

	default:
		mux.makeErrorJsonResponse("TaskHandler: invalid request", resp)
	}
}

// Хэндлер POST обращений к `/api/task`
func (mux Mux) TaskPostHandler(resp http.ResponseWriter, req *http.Request) {
	var task app.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		mux.makeErrorJsonResponse(err.Error(), resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		mux.makeErrorJsonResponse(err.Error(), resp)
		return
	}

	id, err := mux.app.AddTask(task)
	if err != nil {
		fmt.Println("trying to make json responce")
		mux.makeErrorJsonResponse(err.Error(), resp)
		return
	}

	mux.makeJsonResponse(fmt.Sprintf(`{"id":"%d"}`, id), resp)
}

// Хэндлер DELETE обращений к `/api/task`
func (mux Mux) TaskDeleteHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
	id := req.URL.Query().Get("id")
	err := mux.app.RemoveTask(id)

	if err != nil {
		mux.makeErrorJsonResponse(err.Error(), resp)
		return
	}
	mux.makeEmptyJsonResponse(resp)
}

// Хэндлер PUT обращений к `/api/task`
func (mux Mux) TaskPutHandler(resp http.ResponseWriter, req *http.Request) {
	var task app.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		mux.makeErrorJsonResponse(err.Error(), resp)
		return
	}

	err = mux.app.UpdateTask(task)
	// Фронтэнд игнорирует ошибку
	if err != nil {
		mux.makeErrorJsonResponse(err.Error(), resp)
		return
	}
	mux.makeEmptyJsonResponse(resp)
}

// Хэндлер GET обращений к `/api/task`
func (mux Mux) TaskGetHandler(resp http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	task, err := mux.app.GetTask(id)
	if err != nil {
		mux.makeErrorJsonResponse(err.Error(), resp)
		return
	}

	jsonResponse, err := json.Marshal(task)
	if err != nil {
		mux.makeErrorJsonResponse(err.Error(), resp)
		return
	}

	mux.makeJsonResponse(string(jsonResponse), resp)
}

// Хэндлер POST обращений к `/api/task/done`
func (mux Mux) TaskDoneHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(resp, fmt.Errorf("Invalid Request").Error(), http.StatusBadRequest)
		return
	}

	id := req.URL.Query().Get("id")
	err := mux.app.FinishTask(id)

	if err != nil {
		mux.makeErrorJsonResponse(err.Error(), resp)
		return
	}
	mux.makeEmptyJsonResponse(resp)
}

// Хэндлер GET обращений к `/api/tasks`
func (mux Mux) TasksHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(resp, fmt.Errorf("Invalid Request").Error(), http.StatusBadRequest)
		return
	}

	searchString := req.URL.Query().Get("search")

	tasks, err := mux.app.GetTaskList(searchString, config.TaskReturnLimit)
	if err != nil {
		mux.makeErrorJsonResponse(err.Error(), resp)
		return
	}

	taskListBytes, err := json.Marshal(app.TaskList{List: tasks})
	if err != nil {
		mux.makeErrorJsonResponse(err.Error(), resp)
		return
	}

	mux.makeJsonResponse(string(taskListBytes), resp)
}

// Хэндлер GET обращений по `/api/nextdate`
func (mux Mux) NextDateHandler(resp http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	now := query.Get("now")
	date := query.Get("date")
	repeat := query.Get("repeat")

	nextDate, err := mux.app.NextDate(now, date, repeat)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
	}

	_, err = resp.Write([]byte(nextDate))
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
	}
}

// Хэндлер POST обращений к `/api/signin`
func (mux Mux) SignupHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(resp, fmt.Errorf("Invalid Request").Error(), http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		mux.makeErrorJsonResponse(fmt.Errorf("Mux.SignupHandler: %v", err).Error(), resp)
		return
	}

	passStruct := struct {
		Password string `json:"password"`
	}{}
	err = json.Unmarshal(buf.Bytes(), &passStruct)
	if err != nil {
		mux.makeErrorJsonResponse(fmt.Errorf("Mux.SignupHandler: %v", err).Error(), resp)
		return
	}

	valid := mux.auth.VerifyPassword(passStruct.Password)
	if !valid {
		mux.makeErrorJsonResponse("Неверный пароль", resp)
		return
	}
	tokenString, err := mux.auth.CreateToken()
	if err != nil {
		mux.makeErrorJsonResponse(fmt.Errorf("Mux.SignupHandler: %v", err).Error(), resp)
		return
	}
	mux.makeJsonResponse(fmt.Sprintf(`{"token":"%s"}`, tokenString), resp)
}
