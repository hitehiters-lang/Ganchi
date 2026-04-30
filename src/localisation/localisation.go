package localisation

import (
	"context"
	"ganchi_app/connection"
	"ganchi_app/models"
)

// GetTranslation возвращает перевод для ключа и языка
func GetTranslation(ctx context.Context, locKey, lang string) string {
	if locKey == "" {
		return ""
	}
	db := connection.GetDatabase()
	var loc models.Localisation
	err := db.NewSelect().
		Model(&loc).
		Where("loc = ?", locKey).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return locKey // fallback на ключ, если перевода нет
	}
	switch lang {
	case "en":
		return loc.En
	case "kk":
		return loc.Kk
	case "zh":
		return loc.Zh
	default:
		return loc.Ru
	}
}

// TranslateCategory заполняет поле Name переводом
func TranslateCategory(ctx context.Context, cat *models.BlogCategory, lang string) {
	if cat != nil && cat.NameLoc != "" {
		cat.Name = GetTranslation(ctx, cat.NameLoc, lang)
	}
}

// TranslateAuthor заполняет поле Name переводом
func TranslateAuthor(ctx context.Context, author *models.BlogAuthor, lang string) {
	if author != nil && author.NameLoc != "" {
		author.Name = GetTranslation(ctx, author.NameLoc, lang)
	}
}

// TranslateTag заполняет поле Name переводом
func TranslateTag(ctx context.Context, tag *models.BlogTag, lang string) {
	if tag != nil && tag.NameLoc != "" {
		tag.Name = GetTranslation(ctx, tag.NameLoc, lang)
	}
}

// TranslatePost заполняет текстовые поля поста переводами
func TranslatePost(ctx context.Context, post *models.BlogPost, lang string) {
	if post == nil {
		return
	}
	post.Title = GetTranslation(ctx, post.TitleLoc, lang)
	post.Excerpt = GetTranslation(ctx, post.ExcerptLoc, lang)
	post.Content = GetTranslation(ctx, post.ContentLoc, lang)
	post.MetaTitle = GetTranslation(ctx, post.MetaTitleLoc, lang)
	post.MetaDescription = GetTranslation(ctx, post.MetaDescriptionLoc, lang)
}
