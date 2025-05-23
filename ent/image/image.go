// Code generated by ent, DO NOT EDIT.

package image

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the image type in the database.
	Label = "image"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldFilename holds the string denoting the filename field in the database.
	FieldFilename = "filename"
	// FieldURL holds the string denoting the url field in the database.
	FieldURL = "url"
	// FieldSize holds the string denoting the size field in the database.
	FieldSize = "size"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldWidth holds the string denoting the width field in the database.
	FieldWidth = "width"
	// FieldHeight holds the string denoting the height field in the database.
	FieldHeight = "height"
	// EdgeUploadedBy holds the string denoting the uploaded_by edge name in mutations.
	EdgeUploadedBy = "uploaded_by"
	// EdgeBooks holds the string denoting the books edge name in mutations.
	EdgeBooks = "books"
	// Table holds the table name of the image in the database.
	Table = "images"
	// UploadedByTable is the table that holds the uploaded_by relation/edge.
	UploadedByTable = "images"
	// UploadedByInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UploadedByInverseTable = "users"
	// UploadedByColumn is the table column denoting the uploaded_by relation/edge.
	UploadedByColumn = "user_images"
	// BooksTable is the table that holds the books relation/edge.
	BooksTable = "books"
	// BooksInverseTable is the table name for the Book entity.
	// It exists in this package in order to avoid circular dependency with the "book" package.
	BooksInverseTable = "books"
	// BooksColumn is the table column denoting the books relation/edge.
	BooksColumn = "image_books"
)

// Columns holds all SQL columns for image fields.
var Columns = []string{
	FieldID,
	FieldFilename,
	FieldURL,
	FieldSize,
	FieldType,
	FieldCreatedAt,
	FieldWidth,
	FieldHeight,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "images"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"user_images",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// FilenameValidator is a validator for the "filename" field. It is called by the builders before save.
	FilenameValidator func(string) error
	// URLValidator is a validator for the "url" field. It is called by the builders before save.
	URLValidator func(string) error
	// SizeValidator is a validator for the "size" field. It is called by the builders before save.
	SizeValidator func(int64) error
	// TypeValidator is a validator for the "type" field. It is called by the builders before save.
	TypeValidator func(string) error
	// DefaultWidth holds the default value on creation for the "width" field.
	DefaultWidth int
	// DefaultHeight holds the default value on creation for the "height" field.
	DefaultHeight int
)

// OrderOption defines the ordering options for the Image queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByFilename orders the results by the filename field.
func ByFilename(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFilename, opts...).ToFunc()
}

// ByURL orders the results by the url field.
func ByURL(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldURL, opts...).ToFunc()
}

// BySize orders the results by the size field.
func BySize(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSize, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByWidth orders the results by the width field.
func ByWidth(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldWidth, opts...).ToFunc()
}

// ByHeight orders the results by the height field.
func ByHeight(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHeight, opts...).ToFunc()
}

// ByUploadedByField orders the results by uploaded_by field.
func ByUploadedByField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newUploadedByStep(), sql.OrderByField(field, opts...))
	}
}

// ByBooksCount orders the results by books count.
func ByBooksCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newBooksStep(), opts...)
	}
}

// ByBooks orders the results by books terms.
func ByBooks(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newBooksStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newUploadedByStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(UploadedByInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, UploadedByTable, UploadedByColumn),
	)
}
func newBooksStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(BooksInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, BooksTable, BooksColumn),
	)
}
