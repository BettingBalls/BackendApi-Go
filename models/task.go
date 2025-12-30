package models

type Task struct {
	ID 		int64  `json:"id"`
	UserID 	int64  `json:"user_id"`
	Title 	string `json:"title"`
	Description string `json:"description"`
	Status 	string `json:"status"`
	Deadline string `json:"deadline"`
	CreatedAt string `json:"created_at,omitempty"`
}