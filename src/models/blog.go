package models

import "time"

// BlogCategory
type BlogCategory struct {
	ID      int    `json:"id" bun:"id,pk,autoincrement"`
	NameLoc string `json:"-" bun:"name_loc"` // ключ локализации
	Slug    string `json:"slug" bun:"slug"`
	Name    string `json:"name" bun:"-"` // заполняется переводом
}

// BlogAuthor
type BlogAuthor struct {
	ID      int    `json:"id" bun:"id,pk,autoincrement"`
	NameLoc string `json:"-" bun:"name_loc"`
	Avatar  string `json:"avatar" bun:"avatar"`
	Name    string `json:"name" bun:"-"`
}

// BlogTag
type BlogTag struct {
	ID      int    `json:"id" bun:"id,pk,autoincrement"`
	NameLoc string `json:"-" bun:"name_loc"`
	Slug    string `json:"slug" bun:"slug"`
	Name    string `json:"name" bun:"-"`
}

// BlogPost
type BlogPost struct {
	ID                 int       `json:"id" bun:"id,pk,autoincrement"`
	Slug               string    `json:"slug" bun:"slug"`
	TitleLoc           string    `json:"-" bun:"title_loc"`
	ExcerptLoc         string    `json:"-" bun:"excerpt_loc"`
	ContentLoc         string    `json:"-" bun:"content_loc"`
	CoverImage         string    `json:"coverImage" bun:"cover_image"`
	PublishedAt        time.Time `json:"publishedAt" bun:"published_at"`
	CategoryID         *int      `json:"-" bun:"category_id"`
	AuthorID           *int      `json:"-" bun:"author_id"`
	MetaTitleLoc       string    `json:"-" bun:"meta_title_loc"`
	MetaDescriptionLoc string    `json:"-" bun:"meta_description_loc"`

	// Заполняются переводом
	Title           string `json:"title" bun:"-"`
	Excerpt         string `json:"excerpt" bun:"-"`
	Content         string `json:"content" bun:"-"`
	MetaTitle       string `json:"-" bun:"-"`
	MetaDescription string `json:"-" bun:"-"`

	Category *BlogCategory `json:"category,omitempty" bun:"rel:belongs-to,join:category_id=id"`
	Author   *BlogAuthor   `json:"author,omitempty" bun:"rel:belongs-to,join:author_id=id"`
}

// BlogPostListItem — для списка
type BlogPostListItem struct {
	ID          int           `json:"id"`
	Slug        string        `json:"slug"`
	Title       string        `json:"title"`
	Excerpt     string        `json:"excerpt"`
	CoverImage  string        `json:"coverImage"`
	PublishedAt time.Time     `json:"publishedAt"`
	Category    *BlogCategory `json:"category"`
	Author      *BlogAuthor   `json:"author"`
}

// BlogPostDetail — полная статья
type BlogPostDetail struct {
	ID          int           `json:"id"`
	Slug        string        `json:"slug"`
	Title       string        `json:"title"`
	Content     string        `json:"content"`
	CoverImage  string        `json:"coverImage"`
	PublishedAt time.Time     `json:"publishedAt"`
	Category    *BlogCategory `json:"category"`
	Author      *BlogAuthor   `json:"author"`
	Tags        []BlogTag     `json:"tags"`
	SEO         BlogSEO       `json:"seo"`
}

type BlogSEO struct {
	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

type BlogPostsResponse struct {
	Data []BlogPostListItem `json:"data"`
	Meta PaginationMeta     `json:"meta"`
}

type BlogPostResponse struct {
	Data BlogPostDetail `json:"data"`
}

type BlogCategoriesResponse struct {
	Data []BlogCategory `json:"data"`
}
