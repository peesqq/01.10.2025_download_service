package downloader

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yourname/download-service/internal/model"
	"github.com/yourname/download-service/internal/store"
)

type Worker struct {
	store *store.Store
	queue chan *model.Task
}

func NewWorker(st *store.Store) *Worker {
	w := &Worker{
		store: st,
		queue: make(chan *model.Task, 100),
	}
	go w.run()
	// Восстанавливаем pending задачи
	for _, t := range st.GetPendingTasks() {
		w.Enqueue(t)
	}
	return w
}

func (w *Worker) Enqueue(task *model.Task) {
	w.queue <- task
}

func (w *Worker) run() {
	for task := range w.queue {
		log.Printf("[INFO] Задача начата: %s (urls=%d)", task.Name, len(task.URLs))
		task.Status = "in_progress"
		w.store.SaveTask(task)

		w.download(task)

		task.Status = "done"
		w.store.SaveTask(task)
		log.Printf("[INFO] Задача завершена: %s", task.Name)
	}
}

func (w *Worker) download(task *model.Task) {
	safeName := strings.ReplaceAll(task.Name, " ", "_")
	dir := filepath.Join("downloads", safeName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("[ERROR] Не удалось создать директорию %s: %v", dir, err)
		return
	}

	for _, url := range task.URLs {
		log.Printf("[INFO] Начато скачивание файла: %s (task=%s)", url, task.Name)

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("[ERROR] Ошибка при скачивании %s: %v", url, err)
			continue
		}
		defer resp.Body.Close()

		fname := filepath.Base(url)
		fpath := filepath.Join(dir, fname)

		f, err := os.Create(fpath)
		if err != nil {
			log.Printf("[ERROR] Ошибка при создании файла %s: %v", fpath, err)
			continue
		}

		if _, err := io.Copy(f, resp.Body); err != nil {
			log.Printf("[ERROR] Ошибка при записи файла %s: %v", fpath, err)
		}
		f.Close()

		log.Printf("[INFO] Завершено скачивание файла: %s (task=%s)", fpath, task.Name)
	}
}
