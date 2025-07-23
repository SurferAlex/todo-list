package db

import (
	"fmt"
	"testi/internal/entity"
)

func InsertPost(post entity.Post) error {
	_, err := db.Exec("INSERT INTO posts (user_id, content) VALUES ($1, $2)",
		post.UserID, post.Content)
	return err
}

func GetAllPosts() ([]entity.Post, error) {
	rows, err := db.Query(`
		SELECT posts.id, posts.user_id, posts.content, posts.created_at, users.username
		FROM posts
		JOIN users ON posts.user_id = users.id
		ORDER BY posts.created_at DESC
	`)
	if err != nil {
		fmt.Printf("Ошибка SQL запроса: %v\n", err) // добавь для отладки
		return nil, err
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post
		var username string
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt, &username)
		if err != nil {
			fmt.Printf("Ошибка сканирования строки: %v\n", err) // добавь для отладки
			return nil, err
		}
		post.Username = username

		// Загружаем изображения для поста
		images, err := GetImagesByPostID(post.ID)
		if err != nil {
			return nil, err
		}
		post.Images = images

		posts = append(posts, post)
	}

	fmt.Printf("Найдено постов: %d\n", len(posts)) // добавь для отладки
	return posts, nil
}

func CreatePostsTable() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	return err
}

// Удаление поста
func DeletePostByID(postID int, userID int) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = $1 AND user_id = $2", postID, userID)
	return err
}

func InsertPostImage(image entity.PostImage) error {
	_, err := db.Exec(`
		INSERT INTO post_images (post_id, filename, original_name, file_path) 
		VALUES ($1, $2, $3, $4)
	`, image.PostID, image.Filename, image.OriginalName, image.FilePath)
	return err
}

func GetImagesByPostID(postID int) ([]entity.PostImage, error) {
	rows, err := db.Query(`
		SELECT id, post_id, filename, original_name, file_path, created_at
		FROM post_images WHERE post_id = $1 ORDER BY created_at
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []entity.PostImage
	for rows.Next() {
		var img entity.PostImage
		err := rows.Scan(&img.ID, &img.PostID, &img.Filename, &img.OriginalName, &img.FilePath, &img.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

func GetLastPostID(userID int) (int, error) {
	var postID int
	err := db.QueryRow("SELECT id FROM posts WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1", userID).Scan(&postID)
	return postID, err
}
