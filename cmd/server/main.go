package main

import (
    "log"
    "net/http"

    "github.com/yourname/download-service/internal/api"
    "github.com/yourname/download-service/internal/store"
    "github.com/yourname/download-service/internal/downloader"
    "github.com/yourname/download-service/pkg/logger"
)

func main() {
    log.Println("starting server...")

    st := store.NewStore("tasks.json")
    dl := downloader.NewWorker(st)

    api.RegisterHandlers(st, dl)

    log.Println("server started at :8080")
    logger.LogFatal(http.ListenAndServe(":8080", nil))
}
