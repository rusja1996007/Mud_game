package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewHTTPServer(hh *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: hh,
	}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	router.Path("/tasks").Methods("POST").HandlerFunc(s.httpHandlers.HandleCreateTask)
	router.Path("/tasks/{title}").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetTask)
	router.Path("/tasks").Methods("GET").Queries("completed", "true").HandlerFunc(s.httpHandlers.HandleGetAllUncompletedTasks)
	router.Path("/tasks").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetAllTask)
	router.Path("/tasks/{title}").Methods("PATCH").HandlerFunc(s.httpHandlers.HandleCompleteTask)
	router.Path("/tasks/{title}").Methods("DELETE").HandlerFunc(s.httpHandlers.HandleDeleteTask)

	if err := http.ListenAndServe(":228", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil

}
