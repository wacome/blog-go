package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.Text("content").NotEmpty(),
		field.String("excerpt").NotEmpty(),
		field.String("cover_image").Default("/images/post-cover.jpg"),
		field.Bool("published").Default(false),
		field.Time("created_at"),
		field.Time("updated_at"),
		field.Time("published_at").Optional().Nillable(),
		field.Int("views").Default(0),
		field.Enum("author_type").Values("original", "repost").Default("original"),
		field.String("author").NotEmpty(),
	}
}

// Edges of the Post.
func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("comments", Comment.Type),
		edge.To("tags", Tag.Type),
	}
}
