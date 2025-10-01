package store

import (
    "encoding/json"
    "os"
    "sync"

    "github.com/yourname/download-service/internal/model"
)

type Store struct {
    mu    sync.Mutex
    tasks map[string]*model.Task
    file  string
}

func NewStore(file string) *Store {
    s := &Store{
        tasks: make(map[string]*model.Task),
        file:  file,
    }
    s.load()
    return s
}

func (s *Store) SaveTask(task *model.Task) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.tasks[task.ID] = task
    s.persist()
}

func (s *Store) GetTask(id string) *model.Task {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.tasks[id]
}

func (s *Store) GetPendingTasks() []*model.Task {
    s.mu.Lock()
    defer s.mu.Unlock()
    var res []*model.Task
    for _, t := range s.tasks {
        if t.Status == "pending" || t.Status == "in_progress" {
            res = append(res, t)
        }
    }
    return res
}

func (s *Store) persist() {
    f, _ := os.Create(s.file)
    defer f.Close()
    json.NewEncoder(f).Encode(s.tasks)
}

func (s *Store) load() {
    f, err := os.Open(s.file)
    if err != nil {
        return
    }
    defer f.Close()
    json.NewDecoder(f).Decode(&s.tasks)
}
