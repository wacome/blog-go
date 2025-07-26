package models

import "time"

// User 用户模型
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // 不输出到JSON
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Role      string    `json:"role"` // admin, user
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Post 文章模型
type Post struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Excerpt     string    `json:"excerpt"`
	CoverImage  string    `json:"coverImage"`
	Published   bool      `json:"published"`
	AuthorID    int       `json:"authorId"`
	Tags        []Tag     `json:"tags,omitempty"`
	Views       int       `json:"views"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	PublishedAt time.Time `json:"publishedAt,omitempty"`
}

// Tag 标签模型
type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Count     int       `json:"count"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Comment 评论模型
type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	PostID    int       `json:"postId"`
	UserID    int       `json:"userId,omitempty"`
	Author    string    `json:"author"`
	Email     string    `json:"email"`
	Website   string    `json:"website,omitempty"`
	Approved  bool      `json:"approved"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Friend 友链模型
type Friend struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Avatar    string    `json:"avatar"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Book 图书模型
type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Cover       string    `json:"cover"`
	Publisher   string    `json:"publisher"`
	PublishDate string    `json:"publishDate"`
	ISBN        string    `json:"isbn"`
	Pages       int       `json:"pages"`
	Rating      float64   `json:"rating"`
	Status      string    `json:"status"` // reading, finished, want
	Review      string    `json:"review"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Hitokoto 一言模型
type Hitokoto struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Image 图片模型
type Image struct {
	ID         int       `json:"id"`
	Filename   string    `json:"filename"`
	URL        string    `json:"url"`
	Size       int64     `json:"size"`
	Type       string    `json:"type"`
	UploadedBy int       `json:"uploadedBy"`
	CreatedAt  time.Time `json:"createdAt"`
}
