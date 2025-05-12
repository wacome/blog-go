package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Image holds the schema definition for the Image entity.
type Image struct {
	ent.Schema
}

// Fields of the Image.
func (Image) Fields() []ent.Field {
	return []ent.Field{
		field.String("filename").NotEmpty(),
		field.String("url").NotEmpty(),
		field.Int64("size").Positive(),
		field.String("type").NotEmpty(),
		field.Time("created_at"),
		field.Int("width").Default(0),
		field.Int("height").Default(0),
	}
}

// Edges of the Image.
func (Image) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("uploaded_by", User.Type).
			Ref("images").
			Unique().
			Required(),
		edge.To("books", Book.Type),
	}
}
