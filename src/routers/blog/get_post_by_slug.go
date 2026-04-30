package blog

import (
	"encoding/json"
	"net/http"

	"ganchi_app/additional"
	"ganchi_app/connection"
	"ganchi_app/localisation"
	"ganchi_app/models"

	"github.com/go-chi/chi/v5"
)

// @Tags Блог
// @Summary Получить статью по слагy
// @Description Возвращает полную статью с тегами, переводом и SEO
// @ID get_blog_post_by_slug
// @Accept json
// @Produce json
// @Param slug path string true "Слаг статьи"
// @Param lang query string false "Язык: ru, en, kk, zh" default(ru)
// @Success 200 {object} models.BlogPostResponse
// @Failure 404 {object} map[string]string
// @Router /api/blog/posts/{slug} [get]
func GetPostBySlug(w http.ResponseWriter, r *http.Request) {
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

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, `{"error": "slug is required"}`, http.StatusBadRequest)
		return
	}

	var post models.BlogPost
	err := db.NewSelect().Model(&post).Where("slug = ?", slug).Scan(ctx)
	if err != nil {
		additional.PrintError(path, err)
		http.Error(w, `{"error": "post not found"}`, http.StatusNotFound)
		return
	}

	// Категория
	if post.CategoryID != nil {
		var cat models.BlogCategory
		if db.NewSelect().Model(&cat).Where("id = ?", *post.CategoryID).Scan(ctx) == nil {
			localisation.TranslateCategory(ctx, &cat, lang)
			post.Category = &cat
		}
	}

	// Автор
	if post.AuthorID != nil {
		var auth models.BlogAuthor
		if db.NewSelect().Model(&auth).Where("id = ?", *post.AuthorID).Scan(ctx) == nil {
			localisation.TranslateAuthor(ctx, &auth, lang)
			post.Author = &auth
		}
	}

	// Теги + перевод
	var tags []models.BlogTag
	err = db.NewSelect().Model(&tags).
		Join("JOIN blog_post_tags ppt ON ppt.tag_id = tags.id").
		Where("ppt.post_id = ?", post.ID).Scan(ctx)
	if err == nil {
		for i := range tags {
			localisation.TranslateTag(ctx, &tags[i], lang)
		}
	}

	// Перевод поста
	localisation.TranslatePost(ctx, &post, lang)

	response := models.BlogPostResponse{
		Data: models.BlogPostDetail{
			ID: post.ID, Slug: post.Slug, Title: post.Title, Content: post.Content,
			CoverImage: post.CoverImage, PublishedAt: post.PublishedAt,
			Category: post.Category, Author: post.Author, Tags: tags,
			SEO: models.BlogSEO{
				MetaTitle:       post.MetaTitle,
				MetaDescription: post.MetaDescription,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	additional.PrintSuccess(path, "Статья "+slug+" получена на языке "+lang)
}
