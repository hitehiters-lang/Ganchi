package blog

import (
	"encoding/json"
	"ganchi_app/additional"
	"ganchi_app/connection"
	"ganchi_app/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// @Tags Блог
// @Summary Получить комментарии к статье
// @Description Возвращает плоский список всех комментариев с parentId
// @ID get_post_comments
// @Accept json
// @Produce json
// @Param slug path string true "Слаг статьи"
// @Success 200 {object} models.CommentResponse
// @Router /api/blog/posts/{slug}/comments [get]
func GetComments(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	ctx := r.Context()
	db := connection.GetDatabase()

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, `{"error": "slug is required"}`, http.StatusBadRequest)
		return
	}

	var comments []models.Comment
	err := db.NewSelect().
		Model(&comments).
		Where("post_slug = ? AND is_approved = ?", slug, true).
		Order("created_at ASC").
		Scan(ctx)

	if err != nil {
		additional.PrintError(path, err)
		http.Error(w, `{"error": "failed to fetch comments"}`, http.StatusInternalServerError)
		return
	}

	for i := range comments {
		comments[i].Replies = nil
	}

	response := models.CommentResponse{
		Comments: comments,
		Total:    len(comments),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	additional.PrintSuccess(path, "Получено комментариев: "+strconv.Itoa(len(comments)))
}
