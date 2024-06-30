// models/tags.go
package models

import (
	"database/sql"
    "fmt"
	"github.com/google/uuid"
    "github.com/lib/pq"
)

type Tags struct {
    ID          uuid.UUID    `json:"id"`
    Label       string       `json:"label"`
}

func (t *Tags) CreateTag(db *sql.DB) (map[string]interface{}, error) {
    tx, err := db.Begin()
    if err != nil {
        return nil,err
    }
    defer tx.Rollback()

    query := `INSERT INTO tags (id, label) VALUES ($1, $2) RETURNING id`
    t.ID = uuid.New()
    err = tx.QueryRow(query, t.ID, t.Label).Scan(&t.ID)
    if err != nil {
        return nil,err
    }

    err = tx.Commit()

    if err != nil {
        return nil, err
    }

    result := map[string]interface{}{
        "id":      t.ID,
        "label":   t.Label,
    }
    
    return result, nil
}

func (t *Tags) UpdateTag(db *sql.DB) (map[string]interface{}, error) {
    tx, err := db.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()
    var updatedId uuid.UUID
    var updatedLabel string

    query := `UPDATE tags SET label = $2 WHERE id = $1 RETURNING id, label`
    err = tx.QueryRow(query, t.ID, t.Label).Scan(&updatedId, &updatedLabel)

    if err != nil {
        return nil, err
    }

    err = tx.Commit();

    if err != nil {
        return nil, err
    }

    result := map[string]interface{}{
        "id": updatedId,
        "title": updatedLabel,
    }
    
    return result, nil
}

func GetTags(db *sql.DB) ([]map[string]interface{}, error) {
    query := `
    SELECT t.id, t.label, 
    COALESCE(
        ARRAY_AGG(p.title) FILTER (WHERE p.title IS NOT NULL), ARRAY[]::text[]
    ) AS posts 
    FROM tags t
    LEFT JOIN post_tags pt ON pt.tag_id = t.id
    LEFT JOIN posts p ON p.id = pt.post_id
    GROUP BY t.id
    `

    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tags []map[string]interface{}

    for rows.Next() {
        var (
            id    uuid.UUID
            label string
            posts []sql.NullString
        )
        err := rows.Scan(&id, &label, pq.Array(&posts))
        if err != nil {
            return nil, fmt.Errorf("error scanning row: %v", err)
        }

        var postTitles []string
        for _, post := range posts {
            if post.Valid {
                postTitles = append(postTitles, post.String)
            }
        }

        tag := map[string]interface{}{
            "id":    id,
            "label": label,
            "posts": postTitles,
        }

        tags = append(tags, tag)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return tags, nil
}

func (t *Tags) GetTag(db *sql.DB) (map[string]interface{}, error) {

    tx, err := db.Begin()

    if err != nil {
        return nil,err
    }
    defer tx.Rollback()

    
    var updatedId uuid.UUID
    var updatedLabel string
    var updatedPosts []sql.NullString
    
    query := `
        SELECT t.id, t.label, 
        COALESCE(
            ARRAY_AGG(p.title) FILTER (WHERE p.title IS NOT NULL), ARRAY[]::text[]
        ) AS posts 
        FROM tags t
        LEFT JOIN post_tags pt ON pt.tag_id = t.id
        LEFT JOIN posts p ON p.id = pt.post_id
        WHERE t.id = $1
        GROUP BY t.id
        `
    
    err = tx.QueryRow(query, t.ID).Scan(&updatedId, &updatedLabel, pq.Array(&updatedPosts))

    var postTitles []string
    for _, post := range updatedPosts {
        if post.Valid {
            postTitles = append(postTitles, post.String)
        }
    }

    if err != nil {
        return nil, err
    }

    err = tx.Commit()

    if err != nil {
        return nil, err
    }

    result := map[string]interface{}{
        "id": updatedId,
        "label": updatedLabel,
        "posts": postTitles,
    }
    return result, nil
}

func (t *Tags) DeleteTag(db *sql.DB) (map[string]interface{}, error) {

    tx, err := db.Begin()

    if err != nil {
        return nil,err
    }
    defer tx.Rollback()

    queryDeleteLink := `DELETE FROM post_tags WHERE tag_id = $1`
    _, err = tx.Exec(queryDeleteLink, t.ID)
    if err != nil {
        return nil, fmt.Errorf("error delete link: %v", err)
    }

    queryDeleteEntry := `DELETE FROM tags WHERE id = $1`
    result, err := tx.Exec(queryDeleteEntry, t.ID)
    if err != nil {
        return nil, fmt.Errorf("error delete entry: %v", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return nil, fmt.Errorf("error rows sync: %v", err)
    }
    if rowsAffected == 0 {
        return nil, fmt.Errorf("wrong id inserted %v", t.ID)
    }

    err = tx.Commit()
    if err != nil {
        return nil, fmt.Errorf("error in delete process: %v", err)
    }

    return map[string]interface{}{
        "id": t.ID,
    }, nil
}
