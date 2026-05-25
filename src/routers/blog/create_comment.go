package blog

import (
	"encoding/json"
	"fmt"
	"ganchi_app/additional"
	"ganchi_app/connection"
	"ganchi_app/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

// @Tags Блог
// @Summary Создать комментарий к статье
// @Description Добавляет новый комментарий (требует модерации)
// @ID create_post_comment
// @Accept json
// @Produce json
// @Param slug path string true "Слаг статьи"
// @Param comment body models.CommentCreateRequest true "Данные комментария"
// @Success 201 {object} models.Comment
// @Failure 400 {object} map[string]string
// @Router /api/blog/posts/{slug}/comments [post]
func CreateComment(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	ctx := r.Context()
	db := connection.GetDatabase()

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, `{"error": "slug is required"}`, http.StatusBadRequest)
		return
	}

	var post models.BlogPost
	err := db.NewSelect().Model(&post).Where("slug = ?", slug).Scan(ctx)
	if err != nil {
		http.Error(w, `{"error": "post not found"}`, http.StatusNotFound)
		return
	}

	var req models.CommentCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := validateComment(&req); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	if req.ParentID != nil {
		var parent models.Comment
		err = db.NewSelect().
			Model(&parent).
			Where("id = ?", *req.ParentID).
			Where("post_slug = ?", slug).
			Scan(ctx)
		if err != nil {
			http.Error(w, `{"error": "parent comment not found"}`, http.StatusBadRequest)
			return
		}
	}

	comment := models.Comment{
		PostSlug:    slug,
		ParentID:    req.ParentID,
		AuthorName:  req.Name,
		AuthorEmail: req.Email,
		Content:     req.Content,
		IsApproved:  true,
		CreatedAt:   time.Now(),
	}

	_, err = db.NewInsert().Model(&comment).Exec(ctx)
	if err != nil {
		additional.PrintError(path, err)
		http.Error(w, `{"error": "failed to create comment"}`, http.StatusInternalServerError)
		return
	}

	response := models.Comment{
		ID:         comment.ID,
		ParentID:   comment.ParentID,
		AuthorName: comment.AuthorName,
		Content:    comment.Content,
		CreatedAt:  comment.CreatedAt,
		Replies:    []models.Comment{},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	additional.PrintSuccess(path, "Создан комментарий #"+strconv.Itoa(comment.ID))
}

func validateComment(req *models.CommentCreateRequest) error {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Content = strings.TrimSpace(req.Content)

	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(req.Name) > 255 {
		return fmt.Errorf("name too long")
	}
	// Простая валидация email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return fmt.Errorf("invalid email format")
	}
	if req.Content == "" {
		return fmt.Errorf("content is required")
	}
	if len(req.Content) < 3 {
		return fmt.Errorf("content too short")
	}
	if len(req.Content) > 5000 {
		return fmt.Errorf("content too long")
	}
	return nil
}
