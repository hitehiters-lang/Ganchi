package blog

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"

	"ganchi_app/additional"
	"ganchi_app/connection"
	"ganchi_app/localisation"
	"ganchi_app/models"

	"github.com/uptrace/bun"
)

// @Tags Блог
// @Summary Получить список статей
// @Description Возвращает список статей с пагинацией, фильтрацией по категориям и переводом
// @ID get_blog_posts
// @Accept json
// @Produce json
// @Param page query int false "Страница" default(1)
// @Param perPage query int false "На странице" default(10)
// @Param category query []string false "Слаги категорий (можно несколько)" collectionFormat(multi)
// @Param lang query string false "Язык: ru, en, kk, zh" default(ru)
// @Success 200 {object} models.BlogPostsResponse
// @Router /api/blog/posts [get]
func GetPosts(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	ctx := r.Context()
	db := connection.GetDatabase()

	// Язык
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "ru"
	}
	allowedLangs := map[string]bool{"ru": true, "en": true, "kk": true, "zh": true}
	if !allowedLangs[lang] {
		lang = "ru"
	}

	// Пагинация
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	// Получаем все категории: поддерживаем ?category=a&category=b и ?category=a,b,c
	categorySlugs := []string{}
	for _, val := range r.URL.Query()["category"] {
		if val != "" {
			// Если передано через запятую — разбиваем
			for _, slug := range strings.Split(val, ",") {
				slug = strings.TrimSpace(slug)
				if slug != "" {
					categorySlugs = append(categorySlugs, slug)
				}
			}
		}
	}

	posts := []models.BlogPost{}
	query := db.NewSelect().Model(&posts).Order("published_at DESC")

	// Фильтр по категориям
	if len(categorySlugs) > 0 {
		query = query.
			Join("JOIN blog_categories c ON c.id = blog_post.category_id").
			Where("c.slug IN (?)", bun.In(categorySlugs))
	}

	total, err := query.Count(ctx)
	if err != nil {
		additional.PrintError(path, err)
		http.Error(w, `{"error": "failed to count posts"}`, http.StatusInternalServerError)
		return
	}

	err = query.Limit(perPage).Offset((page - 1) * perPage).Scan(ctx)
	if err != nil {
		additional.PrintError(path, err)
		http.Error(w, `{"error": "failed to fetch posts"}`, http.StatusInternalServerError)
		return
	}

	// Загружаем категории и авторов отдельно + переводим
	for i := range posts {
		if posts[i].CategoryID != nil {
			var cat models.BlogCategory
			if db.NewSelect().Model(&cat).Where("id = ?", *posts[i].CategoryID).Scan(ctx) == nil {
				localisation.TranslateCategory(ctx, &cat, lang)
				posts[i].Category = &cat
			}
		}
		if posts[i].AuthorID != nil {
			var auth models.BlogAuthor
			if db.NewSelect().Model(&auth).Where("id = ?", *posts[i].AuthorID).Scan(ctx) == nil {
				localisation.TranslateAuthor(ctx, &auth, lang)
				posts[i].Author = &auth
			}
		}
		localisation.TranslatePost(ctx, &posts[i], lang)
	}

	// Формируем ответ
	items := make([]models.BlogPostListItem, 0, len(posts))
	for _, p := range posts {
		items = append(items, models.BlogPostListItem{
			ID: p.ID, Slug: p.Slug, Title: p.Title, Excerpt: p.Excerpt,
			CoverImage: p.CoverImage, PublishedAt: p.PublishedAt,
			Category: p.Category, Author: p.Author,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	response := models.BlogPostsResponse{
		Data: items,
		Meta: models.PaginationMeta{Page: page, PerPage: perPage, Total: total, TotalPages: totalPages},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	additional.PrintSuccess(path, "Статьи получены на языке "+lang)
}
