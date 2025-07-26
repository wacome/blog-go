package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").NotEmpty().Unique(),
		field.String("email").NotEmpty().Unique(),
		field.String("password").NotEmpty(),
		field.String("role").Default("user"),
		field.String("avatar").Optional(),
		field.String("nickname").Optional(),
		field.String("bio").Optional(),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", Post.Type),
		edge.To("comments", Comment.Type),
		edge.To("tags", Tag.Type),
		edge.To("images", Image.Type),
		edge.To("books", Book.Type),
	}
}
