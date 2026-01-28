package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"restapi/todo"
	"time"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	todoList *todo.List
}

func NewHTTPHandlers(todoList *todo.List) *HTTPHandlers {
	return &HTTPHandlers{
		todoList: todoList,
	}
}

/*
	если по РЕСТ АПИ:

патерн - /tasks
Метод - POST
принимаемая инфа откуда и как - json в теле входящего http запроса

успех:

	статус код - 201 Created
	тело ответа - JSON с ответом данных задачи

ошибка:
статус кода - 400,409,500...
тело ответа - Json с ошибкой
*/
func (h *HTTPHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	//читаем входящий запрос
	var taskDTO TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	//валидируем
	if err := taskDTO.ValidateForCreate(); err != nil {

		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return

	}
	//добваляем задачи в туду лист
	todoTask := todo.NewTask(taskDTO.Title, taskDTO.Opisanie)
	if err := h.todoList.AddTask(todoTask); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, todo.ErrTaskAlreadyExists) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)

		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return

	} //добавляем ответ об добавлении задачи
	b, err := json.MarshalIndent(todoTask, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

/*
патерн - /tasks(title)
Метод - GET
Доп инфа и так будет в патерну(title)

успех:
статус - 200 ok
тело ответа - json с задачей

ошибка:
статус код - 400,404,500...
тело ответа - Json с ошибкой + время
*/
func (h HTTPHandlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]
	task, err := h.todoList.GetTask(title)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)

		}
		return
	}
	b, err := json.MarshalIndent(task, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("ошибка отправки http ответа")
		return
	}
}

/*
патерн - /tasks
Метод - GET
Доп инфа -

успех:
статус - 200 ok
тело ответа - json с задачами

ошибка:
статус код - 400,500...
тело ответа - Json с ошибкой + время


*/

func (h HTTPHandlers) HandleGetAllTask(w http.ResponseWriter, r *http.Request) {
	tasks := h.todoList.ListTasks()
	b, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("ошибка отправеки аштитипи ответа", err)
		return
	}
}

/*
патерн - /tasks?complited=true
Метод - GET
Доп инфа - в квери параметрах

успех:
статус - 200 ok
тело ответа - json с задачами не выполнеными

ошибка:
статус код - 400,500...
тело ответа - Json с ошибкой + время


*/

func (h HTTPHandlers) HandleGetAllUncompletedTasks(w http.ResponseWriter, r *http.Request) {
	uncompletedTasks := h.todoList.ListUncompletedTasks()
	b, err := json.MarshalIndent(uncompletedTasks, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("ошибка отправеки аштитипи ответа", err)
		return
	}
}

/*
патерн - /tasks/{title}-идентификатор
Метод - PATCH
Доп инфа - в патерне+Json в теле

успех:
статус - 200 ok
тело ответа - json с изменеными задачами

ошибка:
статус код - 400,409,500...
тело ответа - Json с ошибкой + время


*/

func (h HTTPHandlers) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
	var completeDTO CompleteTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	title := mux.Vars(r)["title"]

	var (
		changedTask todo.Task
		err         error
	)
	if completeDTO.Complete {
		changedTask, err = h.todoList.CompleteTask(title)
	} else {
		changedTask, err = h.todoList.UncompleteTask(title)
	}
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)

		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}
	x, err := json.MarshalIndent(changedTask, "", "    ")
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(x); err != nil {
		fmt.Println("failed write http response", err)
		return
	}

}

/*
патерн - /tasks/{title}-идентификатор
Метод - DELETE
Доп инфа - в патерне+Json в теле

успех:
статус - 204 No content
тело ответа - пустое

ошибка:
статус код - 400,404,500...
тело ответа - Json с ошибкой + время наступления ошибки


*/

func (h HTTPHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]
	if err := h.todoList.DeleteTask(title); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)

		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
