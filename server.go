// server.go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "go-interview/controllers"
    "go-interview/configs"
)

type Server struct {
    pc *controllers.PostsController
    tc *controllers.TagsController
}

func NewServer() (*Server, error) {
    db, err := configs.InitDB()
    if err != nil {
        return nil, fmt.Errorf("init db failed %v", err)
    }
    pc := &controllers.PostsController{DB: db}
    tc := &controllers.TagsController{DB: db}
    return &Server{
        pc: pc,
        tc: tc,
    }, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if strings.HasPrefix(r.URL.Path, "/api/posts") {
        s.handlePostsRoutes(w, r)
    } else if strings.HasPrefix(r.URL.Path, "/api/tags") {
        s.handleTagsRoutes(w, r)
    } else {
        http.NotFound(w, r)
    }
}

func (s *Server) handlePostsRoutes(w http.ResponseWriter, r *http.Request) {
    switch {
    case r.Method == "GET" && r.URL.Path == "/api/posts":
        s.handleGetAllPosts(w, r)
    case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/api/posts/"):
        s.handleGetPost(w, r)
    case r.Method == "POST" && r.URL.Path == "/api/posts":
        s.handleCreatePost(w, r)
    case r.Method == "PUT" && strings.HasPrefix(r.URL.Path, "/api/posts/"):
        s.handleUpdatePost(w, r)
    case r.Method == "DELETE" && strings.HasPrefix(r.URL.Path, "/api/posts/"):
        s.handleDeletePost(w, r)
    default:
        http.NotFound(w, r)
    }
}

func (s *Server) handleTagsRoutes(w http.ResponseWriter, r *http.Request) {
    switch {
    case r.Method == "POST" && r.URL.Path == "/api/tags":
        s.handleCreateTag(w, r)
    case r.Method == "PUT" && strings.HasPrefix(r.URL.Path, "/api/tags/"):
        s.handleUpdateTag(w, r)
    case r.Method == "GET" && r.URL.Path == "/api/tags":
        s.handleGetAllTags(w, r)
    case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/api/tags/"):
        s.handleGetTag(w, r)
    case r.Method == "DELETE" && strings.HasPrefix(r.URL.Path, "/api/tags/"):
        s.handleDeleteTag(w, r)
    default:
        http.NotFound(w, r)
    }
}

////////////////////////////// POST SECTION ////////////////////////////////////
func (s *Server) handleCreatePost(w http.ResponseWriter, r *http.Request) {
    var post struct {
        Title   string   `json:"title"`
        Content string   `json:"content"`
        Tags    []string `json:"tags"`
    }

    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result, err := s.pc.CreatePost(post.Title, post.Content, post.Tags)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func (s *Server) handleUpdatePost(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/api/posts/");
    var post struct {
        Title       string   `json:"title"`
        Content     string   `json:"content"`
        Status      string   `json:"status"`
        PublishDate string   `json:"publish_date"`
        Tags        []string `json:"tags"`
    }

    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result, err := s.pc.UpdatePost(id, post.Title, post.Content, post.Tags)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func (s *Server) handleGetAllPosts(w http.ResponseWriter, r *http.Request) {
    tagQuery := r.URL.Query().Get("tags")
    posts, err := s.pc.GetAllPosts(tagQuery)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}

func (s *Server) handleGetPost(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/api/posts/")
    post, err := s.pc.GetPost(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}

func (s *Server) handleDeletePost(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/api/posts/")
    post, err := s.pc.DeletePost(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}

////////////////////////////// TAG SECTION ////////////////////////////////////
func (s *Server) handleCreateTag(w http.ResponseWriter, r *http.Request) {
    var tags struct {
        ID   string   `json:"id"`
        Label string   `json:"label"`
    }

    if err := json.NewDecoder(r.Body).Decode(&tags); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result, err := s.tc.CreateTag(tags.Label)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func (s *Server) handleUpdateTag(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/api/tags/");
    var tags struct {
        ID   string   `json:"id"`
        Label string   `json:"label"`
    }

    if err := json.NewDecoder(r.Body).Decode(&tags); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result, err := s.tc.UpdateTag(id, tags.Label)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func (s *Server) handleGetAllTags(w http.ResponseWriter, r *http.Request) {
    tags, err := s.tc.GetAllTags()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tags)
}

func (s *Server) handleGetTag(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/api/tags/")
    tags, err := s.tc.GetTag(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tags)
}

func (s *Server) handleDeleteTag(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/api/tags/")
    tags, err := s.tc.DeleteTag(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tags)
}
