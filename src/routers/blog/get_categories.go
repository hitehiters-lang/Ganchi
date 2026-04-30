package blog

import (
	"encoding/json"
	"net/http"

	"ganchi_app/additional"
	"ganchi_app/connection"
	"ganchi_app/localisation"
	"ganchi_app/models"
)

// @Tags Блог
// @Summary Получить список категорий
// @Description Возвращает все категории блога с переводом
// @ID get_blog_categories
// @Accept json
// @Produce json
// @Param lang query string false "Язык: ru, en, kk, zh" default(ru)
// @Success 200 {array} models.BlogCategory
// @Router /api/blog/categories [get]
func GetCategories(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	ctx := r.Context()
	db := connection.GetDatabase()

	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "ru"
	}
	allowedLangs := map[string]bool{"ru": true, "en": true, "kk": true, "zh": true}
	if !allowedLangs[lang] {
		lang = "ru"
	}

	var categories []models.BlogCategory
	err := db.NewSelect().Model(&categories).Order("slug ASC").Scan(ctx)
	if err != nil {
		additional.PrintError(path, err)
		http.Error(w, `{"error": "failed to fetch categories"}`, http.StatusInternalServerError)
		return
	}

	// Переводим названия категорий
	for i := range categories {
		localisation.TranslateCategory(ctx, &categories[i], lang)
	}

	response := models.BlogCategoriesResponse{Data: categories}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	additional.PrintSuccess(path, "Категории получены на языке "+lang)
}
