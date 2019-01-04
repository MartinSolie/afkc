package main

type Task struct {
	ID          string `json:"task_id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Estimate    int    `json:"estimate,omitempty"`
	Spent       int    `json:"spent,omitempty"`
}

