package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Comment holds the schema definition for the Comment entity.
type Comment struct {
	ent.Schema
}

// Fields of the Comment.
func (Comment) Fields() []ent.Field {
	return []ent.Field{
		field.Text("content").NotEmpty(),
		field.String("author").NotEmpty(),
		field.String("email").NotEmpty(),
		field.String("website").Optional(),
		field.Bool("approved").Default(false),
		field.Time("created_at"),
		field.Time("updated_at"),
		field.String("avatar").Optional().Default("/images/default-avatar.png"),
		field.Int("parent_id").Optional().Nillable(),
	}
}

// Edges of the Comment.
func (Comment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("post", Post.Type).
			Ref("comments").
			Unique().
			Required(),
		edge.From("user", User.Type).
			Ref("comments").
			Unique(),
		edge.To("children", Comment.Type).
			From("parent").
			Unique().
			Field("parent_id"),
	}
}
