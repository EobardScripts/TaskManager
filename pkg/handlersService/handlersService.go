package handlersService

import (
	"TaskManager/pkg/logger"
	"TaskManager/pkg/storage"
	"TaskManager/pkg/utilities"
	"context"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type HandlersService struct {
	storage *storage.Storage
}

func New(storage *storage.Storage) *HandlersService {
	return &HandlersService{storage: storage}
}

func (h *HandlersService) PreloadRoutes() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	r := mux.NewRouter()
	//Задачи
	{
		//Эндпоинт всех задач
		r.HandleFunc("/alltasks", h.AllTasks).Methods(http.MethodGet, http.MethodOptions)
		//Одна задача по ID
		r.HandleFunc("/gettask", h.TaskById).Queries("id", "{id}").Methods(http.MethodGet, http.MethodOptions)
		//Создание новой задачи
		r.HandleFunc("/createtask", h.CreateTask).Methods(http.MethodPost, http.MethodOptions)
		//Обновление задачи
		r.HandleFunc("/updatetask", h.UpdateTask).Methods(http.MethodPut, http.MethodOptions)
		//Удаление задачи
		r.HandleFunc("/deletetask", h.DeleteTask).Queries("id", "{id}").Methods(http.MethodGet, http.MethodOptions)
		//Поиск задачи по taskID и authorID
		r.HandleFunc("/taskby", h.TaskBy).Queries("tid", "{tid}").Queries("aid", "{aid}").Methods(http.MethodGet,
			http.MethodOptions)
		//Создание массива задачи
		r.HandleFunc("/createtasks", h.CreateTasks).Methods(http.MethodPost, http.MethodOptions)
		//Поиск задач по автору authorID
		r.HandleFunc("/taskbyauthor", h.TaskByAuthor).Queries("id", "{id}").Methods(http.MethodGet,
			http.MethodOptions)
		//Поиск задач по метке и labelID
		r.HandleFunc("/taskbylabel", h.TaskByLabel).Queries("id", "{id}").Methods(http.MethodGet,
			http.MethodOptions)
	}

	//Пользователи
	{
		//Эндпоинт всех юзеров
		r.HandleFunc("/allusers", h.AllUsers).Methods(http.MethodGet, http.MethodOptions)
		//Юзер по ID
		r.HandleFunc("/getuser", h.UserById).Queries("id", "{id}").Methods(http.MethodGet, http.MethodOptions)
		//Создание нового юзера
		r.HandleFunc("/createuser", h.CreateUser).Methods(http.MethodPost, http.MethodOptions)
		//Обновление юзера
		r.HandleFunc("/updateuser", h.UpdateUser).Methods(http.MethodPut, http.MethodOptions)
		//Удаление юзера
		r.HandleFunc("/deleteuser", h.DeleteUser).Queries("id", "{id}").Methods(http.MethodGet, http.MethodOptions)
	}

	//Метки
	{
		//Эндпоинт всех метов
		r.HandleFunc("/alllabels", h.AllLabels).Methods(http.MethodGet, http.MethodOptions)
		//Метка по ID
		r.HandleFunc("/getlabel", h.LabelById).Queries("id", "{id}").Methods(http.MethodGet, http.MethodOptions)
		//Создание новой метки
		r.HandleFunc("/createlabel", h.CreateLabel).Methods(http.MethodPost, http.MethodOptions)
		//Обновление метки
		r.HandleFunc("/updatelabel", h.UpdateLabel).Methods(http.MethodPut, http.MethodOptions)
		//Удаление метки
		r.HandleFunc("/deletelabel", h.DeleteLabel).Queries("id", "{id}").Methods(http.MethodGet, http.MethodOptions)
	}

	r.Use(cors.Default().Handler, mux.CORSMethodMiddleware(r))
	// CORS обработчик
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})
	handler := crs.Handler(r)

	srv := &http.Server{
		Addr: "0.0.0.0:8010",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      cors.AllowAll().Handler(handler), // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		logger.Error("%s", err.Error())
	}
	log.Println("shutting down")
	os.Exit(0)
}

//----------------------------------Метки-------------------------------------------------------------

// AllLabels - эндпоинт /alllabels, возвращает массив меток в JSON или 204 код при отсутствии данных
func (h *HandlersService) AllLabels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allLabels, err := h.storage.AllLabels()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("%s", err.Error())
		return
	}

	str := utilities.ToJSON(allLabels)
	if str == "null" {
		http.Error(w, "Метки отсутствуют", http.StatusNoContent)
		logger.Warn("Пустой массив")
		return
	}
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// LabelById - эндпоинт /getlabel?id={id}, возвращает метку в JSON или 204 код при отсутствии данных
func (h *HandlersService) LabelById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	labelID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("%s", err.Error())
		return
	}
	label, err := h.storage.LabelById(labelID)
	if err != nil && err.Error() != "no rows in result set" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("%s", err.Error())
		return
	}

	str := utilities.ToJSON(label)
	if str == "null" || label.ID == 0 {
		http.Error(w, "Метка отсутствуют", http.StatusNoContent)
		logger.Warn("Задача отсутствуют")
		return
	}
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// CreateLabel - эндпоинт /createlabel, возвращает новую метку в JSON или ошибку
func (h HandlersService) CreateLabel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	newLabel := &storage.Label{}

	if err := json.NewDecoder(r.Body).Decode(newLabel); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		logger.Error("Ошибка при декодировании тела запроса: %s", err.Error())
		return
	}

	err := h.storage.NewLabel(newLabel)
	if err != nil {
		logger.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	str := utilities.ToJSON(newLabel)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// UpdateLabel - эндпоинт /updatelabel, возвращает обновленную метку в JSON или ошибку
func (h HandlersService) UpdateLabel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	updateLabel := &storage.Label{}

	if err := json.NewDecoder(r.Body).Decode(updateLabel); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		logger.Error("Ошибка при декодировании тела запроса: %s", err.Error())
		return
	}

	logger.Info("update label: %s", utilities.ToJSON(updateLabel))

	err := h.storage.UpdateLabel(updateLabel)
	if err != nil {
		logger.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	str := utilities.ToJSON(updateLabel)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// DeleteLabel - эндпоинт /deletelabel?id={id}, возвращает удаленную метку в JSON или ошибку
func (h HandlersService) DeleteLabel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		logger.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deletedLabel, err := h.storage.DeleteLabel(id)
	if err != nil {
		logger.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	str := utilities.ToJSON(deletedLabel)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

//----------------------------------Пользователи-------------------------------------------------------------

// AllUsers - эндпоинт /allusers, возвращает массив пользователей в JSON или 204 код при отсутствии данных
func (h *HandlersService) AllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allUsers, err := h.storage.AllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("%s", err.Error())
		return
	}

	str := utilities.ToJSON(allUsers)
	if str == "null" {
		http.Error(w, "Пользователи отсутствуют", http.StatusNoContent)
		logger.Warn("Пустой массив")
		return
	}
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// UserById - эндпоинт /getuser?id={id}, возвращает юзера в JSON или 204 код при отсутствии данных
func (h *HandlersService) UserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("%s", err.Error())
		return
	}
	user, err := h.storage.UserById(userID)
	if err != nil && err.Error() != "no rows in result set" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("%s", err.Error())
		return
	}

	str := utilities.ToJSON(user)
	if str == "null" || user.ID == 0 {
		http.Error(w, "Задача отсутствуют", http.StatusNoContent)
		logger.Warn("Задача отсутствуют")
		return
	}
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// CreateUser - эндпоинт /createuser, возвращает новго юзера в JSON или ошибку
func (h HandlersService) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	newUser := &storage.User{}

	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		logger.Error("Ошибка при декодировании тела запроса: %s", err.Error())
		return
	}

	err := h.storage.NewUser(newUser)
	if err != nil {
		logger.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	str := utilities.ToJSON(newUser)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// UpdateUser - эндпоинт /updateuser, возвращает обновленного юзера в JSON или ошибку
func (h HandlersService) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	updateUser := &storage.User{}

	if err := json.NewDecoder(r.Body).Decode(updateUser); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		logger.Error("Ошибка при декодировании тела запроса: %s", err.Error())
		return
	}

	logger.Info("update user: %s", utilities.ToJSON(updateUser))

	err := h.storage.UpdateUser(updateUser)
	if err != nil {
		logger.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	str := utilities.ToJSON(updateUser)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// DeleteUser - эндпоинт /deleteuser?id={id}, возвращает удаленного юзера в JSON или ошибку
func (h HandlersService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		logger.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deletedUser, err := h.storage.DeleteUser(id)
	if err != nil {
		logger.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	str := utilities.ToJSON(deletedUser)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

//----------------------------------Задачи-----------------------------------------------------------------

// TaskByLabel - эндпоинт /taskbylabel?id={id}, возвращает задачу в JSON или 204 код при отсутствии данных
func (h *HandlersService) TaskByLabel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	labelID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("%s", err.Error())
		return
	}

	task, err := h.storage.TasksByLabel(labelID)
	if err != nil && err.Error() != "no rows in result set" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("%s", err.Error())
		return
	}

	str := utilities.ToJSON(task)
	if str == "null" {
		http.Error(w, "Задача отсутствуют", http.StatusNoContent)
		logger.Warn("Задача отсутствуют")
		return
	}
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// TaskByAuthor - эндпоинт /taskbyauthor?id={id}, возвращает задачу в JSON или 204 код при отсутствии данных
func (h *HandlersService) TaskByAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	authorID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("%s", err.Error())
		return
	}

	task, err := h.storage.TasksByAuthor(authorID)
	if err != nil && err.Error() != "no rows in result set" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("%s", err.Error())
		return
	}

	str := utilities.ToJSON(task)
	if str == "null" {
		http.Error(w, "Задача отсутствуют", http.StatusNoContent)
		logger.Warn("Задача отсутствуют")
		return
	}
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// TaskBy - эндпоинт /taskby?tid={tid}&aid={aid}, возвращает задачу в JSON или 204 код при отсутствии данных
func (h *HandlersService) TaskBy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["tid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("%s", err.Error())
		return
	}
	authorID, err := strconv.Atoi(vars["aid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("%s", err.Error())
		return
	}

	task, err := h.storage.Tasks(taskID, authorID)
	if err != nil && err.Error() != "no rows in result set" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("%s", err.Error())
		return
	}

	str := utilities.ToJSON(task)
	if str == "null" {
		http.Error(w, "Задача отсутствуют", http.StatusNoContent)
		logger.Warn("Задача отсутствуют")
		return
	}
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// TaskById - эндпоинт /gettask?id={id}, возвращает задачу в JSON или 204 код при отсутствии данных
func (h *HandlersService) TaskById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	taskID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("%s", err.Error())
		return
	}
	task, err := h.storage.TaskById(int(taskID))
	if err != nil && err.Error() != "no rows in result set" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("%s", err.Error())
		return
	}

	str := utilities.ToJSON(task)
	if str == "null" || task.ID == 0 {
		http.Error(w, "Задача отсутствуют", http.StatusNoContent)
		logger.Warn("Задача отсутствуют")
		return
	}
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// AllTasks - эндпоинт /alltasks, возвращает массив задач в JSON или 204 код при отсутствии данных
func (h *HandlersService) AllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allTasks, err := h.storage.AllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("%s", err.Error())
		return
	}

	str := utilities.ToJSON(allTasks)
	if str == "null" {
		http.Error(w, "Задачи отсутствуют", http.StatusNoContent)
		logger.Warn("Пустой массив")
		return
	}
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// CreateTask - эндпоинт /CreateTask, возвращает созданную задач в JSON или ошибку
func (h HandlersService) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	newTask := &storage.Task{}

	if err := json.NewDecoder(r.Body).Decode(newTask); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		logger.Error("Ошибка при декодировании тела запроса: %s", err.Error())
		return
	}

	err := h.storage.NewTask(newTask)
	if err != nil {
		logger.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	str := utilities.ToJSON(newTask)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// CreateTasks - эндпоинт /createtasks, возвращает созданные задачи в JSON или ошибку
func (h HandlersService) CreateTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	newTasks := []*storage.Task{}

	if err := json.NewDecoder(r.Body).Decode(&newTasks); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		logger.Error("Ошибка при декодировании тела запроса: %s", err.Error())
		return
	}

	logger.Info("Массив задач: %s", utilities.ToJSON(newTasks))

	err := h.storage.NewTasks(newTasks)
	if err != nil {
		logger.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, task := range newTasks {
		newTasks[i], err = h.storage.TaskById(task.ID)
	}

	str := utilities.ToJSON(newTasks)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// UpdateTask - эндпоинт /updatetask, возвращает обновленную задачу в JSON или ошибку
func (h HandlersService) UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	updateTask := &storage.Task{}

	if err := json.NewDecoder(r.Body).Decode(updateTask); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		logger.Error("Ошибка при декодировании тела запроса: %s", err.Error())
		return
	}

	logger.Info("update task: %s", utilities.ToJSON(updateTask))

	err := h.storage.UpdateTask(updateTask)
	if err != nil {
		logger.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	str := utilities.ToJSON(updateTask)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// DeleteTask - эндпоинт /deletetask?id={id}, возвращает удаленную задачу в JSON или ошибку
func (h HandlersService) DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		logger.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deletedTask, err := h.storage.DeleteTask(int(id))
	if err != nil {
		logger.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	str := utilities.ToJSON(deletedTask)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}
