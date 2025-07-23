package entity

import "time"

type Post struct {
	ID        int
	UserID    int
	Content   string
	CreatedAt time.Time
	Username  string
	Images    []PostImage
}

type PostImage struct {
	ID           int
	PostID       int
	Filename     string
	OriginalName string
	FilePath     string
	CreatedAt    time.Time
}
