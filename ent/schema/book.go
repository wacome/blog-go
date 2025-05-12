package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Book holds the schema definition for the Book entity.
type Book struct {
	ent.Schema
}

// Fields of the Book.
func (Book) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.String("author").NotEmpty(),
		field.Text("desc").Optional(),
		field.String("cover").Default("/images/default-book-cover.jpg"),
		field.String("publisher").Optional(),
		field.String("publish_date").Optional(),
		field.String("isbn").Optional(),
		field.Int("pages").Optional(),
		field.Float("rating").Default(0),
		field.Enum("status").Values("reading", "finished", "want").Default("want"),
		field.Text("review").Optional(),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Edges of the Book.
func (Book) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("cover_image", Image.Type).
			Ref("books").
			Unique(),
		edge.From("owner", User.Type).
			Ref("books").
			Unique().
			Required(),
	}
}
