// controllers/posts_controller.go
package controllers

import (
	"database/sql"
	"go-interview/models"
	"time"
	"github.com/google/uuid"
)

type PostsController struct {
    DB *sql.DB
}

func (pc *PostsController) CreatePost(title, content string, tags []string) (map[string]interface{}, error) {
    post := models.Posts{
        Title:   title,
        Content: content,
        Tags:    tags,
        Status:  "draft",
        PublishDate: sql.NullTime{Time: time.Now(), Valid: true},
    }

    result, err := post.CreatePost(pc.DB)
    if err != nil {
        return nil, err
    }

    return result, nil
}

func (pc *PostsController) UpdatePost(idStr, title, content string, tags []string) (map[string]interface{}, error) {
    id, err := uuid.Parse(idStr)
    post := models.Posts{
        ID: id, 
        Title: title, 
        Content: content, 
        Tags: tags,
    }

    result,err := post.UpdatePost(pc.DB)
    if err != nil {
        return nil, err
    }

    return result, nil
}

func (pc *PostsController) GetAllPosts(tagQuery string) ([]map[string]interface{}, error) {
    posts, err := models.GetPosts(pc.DB, tagQuery)
    if err != nil {
        return nil, err
    }
    return posts, nil
}

func (pc *PostsController) GetPost(idStr string) (map[string]interface{}, error) {
    id, err := uuid.Parse(idStr)
    post := models.Posts{
        ID: id,
    }
    result, err := post.GetPost(pc.DB)
    if err != nil {
        return nil, err
    }

    return result, nil

}

func (pc *PostsController) DeletePost(idStr string) (map[string]interface{}, error) {
    id, err := uuid.Parse(idStr)
    post := models.Posts{
        ID: id,
    }
    result, err := post.DeletePost(pc.DB)
    if err != nil {
        return nil, err
    }

    return result, nil

}