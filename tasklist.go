package main

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"os"
	"time"

	"tasklist/internal"
)

type taskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *taskServer {
	store := taskstore.New()
	return &taskServer{store: store}
}

func (ts *taskServer) createTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling task create at %s\n", req.URL.Path)

	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	type ResponseId struct {
		Id int `json:"id"`
	}

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rt RequestTask
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	js, err := json.Marshal(ResponseId{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) getAllTasksHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handing get all tasks at %s\n", req.URL.Path)

	allTasks := ts.store.GetAllTasks()
	js, err := json.Marshal(allTasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// func (ts *taskServer) getTaskHandler(w http.ResponseWriter, req *http.Request) {
// 	log.Printf("handling get task as %s\n", req.URL.Path)
//
// }

func main() {
	mux := http.NewServeMux()
	server := NewTaskServer()

	mux.HandleFunc("POST /task/", server.createTaskHandler)
	mux.HandleFunc("GET /task/", server.getAllTasksHandler)

	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVEPORT"), mux))
}
