package entity

type Comment struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
	BlogID  int    `json:"blog_id"`
}
