package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Collection struct {
	ent.Schema
}

func (Collection) Fields() []ent.Field {
	return []ent.Field{
		field.String("type").NotEmpty(),
		field.String("title").NotEmpty(),
		field.String("author").Optional(),
		field.String("cover").Optional(),
		field.String("date").Optional(),
		field.String("link").Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
