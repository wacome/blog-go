package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Friend holds the schema definition for the Friend entity.
type Friend struct {
	ent.Schema
}

// Fields of the Friend.
func (Friend) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("url").NotEmpty(),
		field.String("avatar").Default("/images/default-avatar.png"),
		field.String("desc").Optional(),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}
