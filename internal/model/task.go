package model

type Task struct {
    ID     string   `json:"id"`
    Name   string   `json:"name"`
    URLs   []string `json:"urls"`
    Status string   `json:"status"`
}
