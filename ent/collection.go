// Code generated by ent, DO NOT EDIT.

package ent

import (
	"blog-go/ent/collection"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// Collection is the model entity for the Collection schema.
type Collection struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty"`
	// Title holds the value of the "title" field.
	Title string `json:"title,omitempty"`
	// Author holds the value of the "author" field.
	Author string `json:"author,omitempty"`
	// Cover holds the value of the "cover" field.
	Cover string `json:"cover,omitempty"`
	// Date holds the value of the "date" field.
	Date string `json:"date,omitempty"`
	// Link holds the value of the "link" field.
	Link string `json:"link,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Collection) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case collection.FieldID:
			values[i] = new(sql.NullInt64)
		case collection.FieldType, collection.FieldTitle, collection.FieldAuthor, collection.FieldCover, collection.FieldDate, collection.FieldLink:
			values[i] = new(sql.NullString)
		case collection.FieldCreatedAt, collection.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Collection fields.
func (c *Collection) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case collection.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			c.ID = int(value.Int64)
		case collection.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				c.Type = value.String
			}
		case collection.FieldTitle:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field title", values[i])
			} else if value.Valid {
				c.Title = value.String
			}
		case collection.FieldAuthor:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field author", values[i])
			} else if value.Valid {
				c.Author = value.String
			}
		case collection.FieldCover:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field cover", values[i])
			} else if value.Valid {
				c.Cover = value.String
			}
		case collection.FieldDate:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field date", values[i])
			} else if value.Valid {
				c.Date = value.String
			}
		case collection.FieldLink:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field link", values[i])
			} else if value.Valid {
				c.Link = value.String
			}
		case collection.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				c.CreatedAt = value.Time
			}
		case collection.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				c.UpdatedAt = value.Time
			}
		default:
			c.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Collection.
// This includes values selected through modifiers, order, etc.
func (c *Collection) Value(name string) (ent.Value, error) {
	return c.selectValues.Get(name)
}

// Update returns a builder for updating this Collection.
// Note that you need to call Collection.Unwrap() before calling this method if this Collection
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Collection) Update() *CollectionUpdateOne {
	return NewCollectionClient(c.config).UpdateOne(c)
}

// Unwrap unwraps the Collection entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Collection) Unwrap() *Collection {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Collection is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Collection) String() string {
	var builder strings.Builder
	builder.WriteString("Collection(")
	builder.WriteString(fmt.Sprintf("id=%v, ", c.ID))
	builder.WriteString("type=")
	builder.WriteString(c.Type)
	builder.WriteString(", ")
	builder.WriteString("title=")
	builder.WriteString(c.Title)
	builder.WriteString(", ")
	builder.WriteString("author=")
	builder.WriteString(c.Author)
	builder.WriteString(", ")
	builder.WriteString("cover=")
	builder.WriteString(c.Cover)
	builder.WriteString(", ")
	builder.WriteString("date=")
	builder.WriteString(c.Date)
	builder.WriteString(", ")
	builder.WriteString("link=")
	builder.WriteString(c.Link)
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(c.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(c.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// Collections is a parsable slice of Collection.
type Collections []*Collection
