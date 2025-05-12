package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Hitokoto holds the schema definition for the Hitokoto entity.
type Hitokoto struct {
	ent.Schema
}

// Fields of the Hitokoto.
func (Hitokoto) Fields() []ent.Field {
	return []ent.Field{
		field.Text("content").NotEmpty(),
		field.String("source").Optional(),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}
