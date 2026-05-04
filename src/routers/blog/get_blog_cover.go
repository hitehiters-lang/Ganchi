package blog

import (
	"database/sql"
	"errors"
	"ganchi_app/additional"
	"ganchi_app/connection"
	"ganchi_app/models"
	"net/http"
	"os"
	"runtime"

	"github.com/go-chi/chi/v5"
)

// @Tags Блог
// @Summary Получить фото обложки статьи по слагy
// @Description Отправляет изображение обложки поста с указанным слагом
// @ID get_blog_cover_image
// @Produce image/png,image/jpeg,image/webp
// @Param slug path string true "Slug"
// @Success 200 {file} binary
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/blog/posts/{slug}/cover [get]
func GetBlogCoverImage(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	ctx := r.Context()
	db := connection.GetDatabase()

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		additional.PrintError(path, errors.New("slug is empty"))
		http.Error(w, `{"error":"slug is required"}`, http.StatusBadRequest)
		return
	}

	var post models.BlogPost
	err := db.NewSelect().Model(&post).Where("slug = ?", slug).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		additional.PrintError(path, err)
		http.Error(w, `{"error":"post not found"}`, http.StatusNotFound)
		return
	}
	if err != nil {
		additional.PrintError(path, err)
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}

	// Выбираем путь в зависимости от ОС
	var imagePath string
	switch runtime.GOOS {
	case "windows":
		imagePath = post.BlogPicturePathWindows
	default:
		imagePath = post.BlogPicturePathLinux
	}

	// Если путь пустой — возвращаем 404
	if imagePath == "" {
		additional.PrintError(path, errors.New("image path is empty"))
		http.Error(w, `{"error":"image path not configured"}`, http.StatusNotFound)
		return
	}

	// Проверяем существование файла
	_, err = os.Stat(imagePath)
	if errors.Is(err, os.ErrNotExist) {
		additional.PrintError(path, err)
		http.Error(w, `{"error":"image file not found"}`, http.StatusNotFound)
		return
	}
	if err != nil {
		additional.PrintError(path, err)
		http.Error(w, `{"error":"file system error"}`, http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, imagePath)
	additional.PrintSuccess(path, "Обложка поста '"+slug+"' отправлена")
}
