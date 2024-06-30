package views

import (
	"fmt"
	"go-interview/models"
)

func ShowPost(post models.Posts) {
    fmt.Printf("Post ID: %s\nTitle: %s\nContent: %s\nStatus: %s\nTags: %s\n", 
        post.ID, post.Title, post.Content, post.Status, post.Tags)
    if post.PublishDate.Valid {
        fmt.Printf("Publish Date: %s\n", post.PublishDate.Time.Format("2006-01-02"))
    }
}

func ShowError(err error) {
	fmt.Printf("Error: %v\n", err)
}

func ShowSuccess(message string) {
	fmt.Println(message)
}