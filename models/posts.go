// models/posts.go
package models

import (
	"database/sql"
    "fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Posts struct {
    ID          uuid.UUID    `json:"id"`
    Title       string       `json:"title"`
    Content     string       `json:"content"`
    Status      string       `json:"status"`
    PublishDate sql.NullTime `json:"publish_date"`
    Tags        []string     `json:"tags"`
}

func (p *Posts) CreatePost(db *sql.DB) (map[string]interface{}, error) {
    tx, err := db.Begin()
    if err != nil {
        return nil,err
    }
    defer tx.Rollback()

    query := `INSERT INTO posts (id, title, content, status, publish_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`
    p.ID = uuid.New()
    err = tx.QueryRow(query, p.ID, p.Title, p.Content, p.Status, p.PublishDate).Scan(&p.ID)
    if err != nil {
        return nil,err
    }

    var tagIDs []uuid.UUID
    var tagLabels []string

    for _, tag := range p.Tags {
        var tagID uuid.UUID
        err = tx.QueryRow("INSERT INTO tags (label) VALUES ($1) ON CONFLICT (label) DO UPDATE SET label = EXCLUDED.label RETURNING id", tag).Scan(&tagID)
        if err != nil {
            return nil,err
        }

        _, err = tx.Exec("INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2)", p.ID, tagID)
        if err != nil {
            return nil, err
        }

        _, err = tx.Exec("UPDATE tags SET posts = array_append(posts, $1) WHERE id = $2", p.ID, tagID)
        if err != nil {
            return nil, err
        }

        tagLabels = append(tagLabels, tag)
        tagIDs = append(tagIDs, tagID)
    }

    _, err = tx.Exec("UPDATE posts SET tags = $1 WHERE id = $2", pq.Array(tagIDs), p.ID)
    if err != nil {
        return nil,err
    }

    err = tx.Commit()

    if err != nil {
        return nil, err
    }

    result := map[string]interface{}{
        "id":      p.ID,
        "title":   p.Title,
        "content": p.Content,
        "tags":    tagLabels,
    }
    
    return result, nil
}

func (p *Posts) UpdatePost(db *sql.DB) (map[string]interface{}, error) {
    tx, err := db.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()
    var updatedId uuid.UUID
    var updatedTitle, updatedContent string

    query := `UPDATE posts SET title = $2, content = $3 WHERE id = $1 RETURNING id, title, content`
    err = tx.QueryRow(query, p.ID, p.Title, p.Content).Scan(&updatedId, &updatedTitle, &updatedContent)

    if err != nil {
        return nil, err
    }

    _, err = tx.Exec("UPDATE tags SET posts = array_remove(posts, $1)", p.ID)
    if err != nil {
        return nil, err
    }

    _, err = tx.Exec("DELETE FROM post_tags WHERE post_id = $1", p.ID)
    if err != nil {
        return nil, err
    }

    var tagLabels []string

    for _, tag := range p.Tags {
        var tagID uuid.UUID
        err = tx.QueryRow("INSERT INTO tags (label) VALUES ($1) ON CONFLICT (label) DO UPDATE SET label = EXCLUDED.label RETURNING id", tag).Scan(&tagID)
        if err != nil {
            return nil, err
        }

        _, err = tx.Exec("INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2)", p.ID, tagID)
        if err != nil {
            return nil, err
        }
        tagLabels = append(tagLabels, tag)
    }


    err = tx.Commit();

    if err != nil {
        return nil, err
    }

    result := map[string]interface{}{
        "id": updatedId,
        "title": updatedTitle,
        "content": updatedContent,
        "tags": tagLabels,
    }
    
    return result, nil
}

func GetPosts(db *sql.DB, tagQuery string) ([]map[string]interface{}, error) {
    var query string
    var args []interface{}
    if tagQuery != "" {
        query = `
        SELECT DISTINCT p.id, p.title, p.content,
        (
            SELECT ARRAY_AGG(t2.label)
            FROM post_tags pt2
            JOIN tags t2 ON pt2.tag_id = t2.id
            WHERE pt2.post_id = p.id
        ) AS tags
        FROM posts p
        JOIN post_tags pt ON p.id = pt.post_id
        JOIN tags t ON pt.tag_id = t.id
        WHERE t.label = $1
        GROUP BY p.id`
    args = append(args, tagQuery)
    } else {
        query = `
        SELECT p.id, p.title, p.content,
        COALESCE(
            ARRAY_AGG(t.label) FILTER (WHERE t.label IS NOT NULL), ARRAY[]::text[]
        ) AS tags
        FROM posts p
        LEFT JOIN post_tags pt ON p.id = pt.post_id
        LEFT JOIN tags t ON pt.tag_id = t.id
        GROUP BY p.id`
    }

    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []map[string]interface{}

    for rows.Next() {
        var (
            id          uuid.UUID
            title       string
            content     string
            tags        []sql.NullString
        )

        err := rows.Scan(&id, &title, &content, pq.Array(&tags))
        if err != nil {
            return nil, err
        }

        var tagLabels []string
        for _, tag := range tags {
            if tag.Valid {
                tagLabels = append(tagLabels, tag.String)
            }
        }

        post := map[string]interface{}{
            "id":      id,
            "title":   title,
            "content": content,
            "tags":    tagLabels,
        }

        posts = append(posts, post)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return posts, nil
}

func (p *Posts) GetPost(db *sql.DB) (map[string]interface{}, error) {

    tx, err := db.Begin()

    if err != nil {
        return nil,err
    }
    defer tx.Rollback()

    var updatedId uuid.UUID
    var updatedTitle, updatedContent string
    var updatedTags []sql.NullString
    query := `
        SELECT p.id, p.title, p.content,
        COALESCE(
            ARRAY_AGG(t.label) FILTER (WHERE t.label IS NOT NULL), ARRAY[]::text[]
        ) AS tags        
        FROM posts p
        LEFT JOIN post_tags pt ON pt.post_id = p.id 
        LEFT JOIN tags t ON t.id= pt.tag_id 
        WHERE p.id = $1
        GROUP BY p.id
    `
    err = tx.QueryRow(query, p.ID).Scan(&updatedId, &updatedTitle, &updatedContent, pq.Array(&updatedTags))

    var tagLabels []string
    for _, tag := range updatedTags {
        if tag.Valid {
            tagLabels = append(tagLabels, tag.String)
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
        "title": updatedTitle,
        "content": updatedContent,
        "tags": tagLabels,
    }
    return result, nil
}

func (p *Posts) DeletePost(db *sql.DB) (map[string]interface{}, error) {

    tx, err := db.Begin()

    if err != nil {
        return nil,err
    }
    defer tx.Rollback()

    queryDeleteLink := `DELETE FROM post_tags WHERE post_id = $1`
    _, err = tx.Exec(queryDeleteLink, p.ID)
    if err != nil {
        return nil, fmt.Errorf("error delete link: %v", err)
    }

    queryDeleteEntry := `DELETE FROM posts WHERE id = $1`
    result, err := tx.Exec(queryDeleteEntry, p.ID)
    if err != nil {
        return nil, fmt.Errorf("error delete entry: %v", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return nil, fmt.Errorf("error rows sync: %v", err)
    }
    if rowsAffected == 0 {
        return nil, fmt.Errorf("wrong id inserted %v", p.ID)
    }

    err = tx.Commit()
    if err != nil {
        return nil, fmt.Errorf("error in delete process: %v", err)
    }

    return map[string]interface{}{
        "id": p.ID,
    }, nil
}
