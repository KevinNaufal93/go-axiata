// controllers/tags_controller.go
package controllers

import (
	"database/sql"
	"go-interview/models"
	"github.com/google/uuid"
)

type TagsController struct {
    DB *sql.DB
}

func (tc *TagsController) CreateTag(label string) (map[string]interface{}, error) {
    tags := models.Tags{Label: label}

    result, err := tags.CreateTag(tc.DB)

    if err != nil {
        return nil, err
    }

    return result, nil
}

func (tc *TagsController) UpdateTag(idStr, label string) (map[string]interface{}, error) {
    id, err := uuid.Parse(idStr)

    tags := models.Tags{
        ID: id, 
        Label: label, 
    }

    result,err := tags.UpdateTag(tc.DB)
    if err != nil {
        return nil, err
    }

    return result, nil
}

func (tc *TagsController) GetAllTags() ([]map[string]interface{}, error) {
    tags, err := models.GetTags(tc.DB)
    if err != nil {
        return nil, err
    }
    return tags, nil
}

func (tc *TagsController) GetTag(idStr string) (map[string]interface{}, error) {
    id, err := uuid.Parse(idStr)
    tags := models.Tags{
        ID: id,
    }
    result, err := tags.GetTag(tc.DB)
    if err != nil {
        return nil, err
    }

    return result, nil

}

func (tc *TagsController) DeleteTag(idStr string) (map[string]interface{}, error) {
    id, err := uuid.Parse(idStr)
    tags := models.Tags{
        ID: id,
    }
    result, err := tags.DeleteTag(tc.DB)
    if err != nil {
        return nil, err
    }

    return result, nil

}