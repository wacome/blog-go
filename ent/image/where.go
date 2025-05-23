// Code generated by ent, DO NOT EDIT.

package image

import (
	"blog-go/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Image {
	return predicate.Image(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Image {
	return predicate.Image(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Image {
	return predicate.Image(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Image {
	return predicate.Image(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Image {
	return predicate.Image(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Image {
	return predicate.Image(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Image {
	return predicate.Image(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Image {
	return predicate.Image(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Image {
	return predicate.Image(sql.FieldLTE(FieldID, id))
}

// Filename applies equality check predicate on the "filename" field. It's identical to FilenameEQ.
func Filename(v string) predicate.Image {
	return predicate.Image(sql.FieldEQ(FieldFilename, v))
}

// URL applies equality check predicate on the "url" field. It's identical to URLEQ.
func URL(v string) predicate.Image {
	return predicate.Image(sql.FieldEQ(FieldURL, v))
}

// Size applies equality check predicate on the "size" field. It's identical to SizeEQ.
func Size(v int64) predicate.Image {
	return predicate.Image(sql.FieldEQ(FieldSize, v))
}

// Type applies equality check predicate on the "type" field. It's identical to TypeEQ.
func Type(v string) predicate.Image {
	return predicate.Image(sql.FieldEQ(FieldType, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Image {
	return predicate.Image(sql.FieldEQ(FieldCreatedAt, v))
}

// Width applies equality check predicate on the "width" field. It's identical to WidthEQ.
func Width(v int) predicate.Image {
	return predicate.Image(sql.FieldEQ(FieldWidth, v))
}

// Height applies equality check predicate on the "height" field. It's identical to HeightEQ.
func Height(v int) predicate.Image {
	return predicate.Image(sql.FieldEQ(FieldHeight, v))
}

// FilenameEQ applies the EQ predicate on the "filename" field.
func FilenameEQ(v string) predicate.Image {
	return predicate.Image(sql.FieldEQ(FieldFilename, v))
}

// FilenameNEQ applies the NEQ predicate on the "filename" field.
func FilenameNEQ(v string) predicate.Image {
	return predicate.Image(sql.FieldNEQ(FieldFilename, v))
}

// FilenameIn applies the In predicate on the "filename" field.
func FilenameIn(vs ...string) predicate.Image {
	return predicate.Image(sql.FieldIn(FieldFilename, vs...))
}

// FilenameNotIn applies the NotIn predicate on the "filename" field.
func FilenameNotIn(vs ...string) predicate.Image {
	return predicate.Image(sql.FieldNotIn(FieldFilename, vs...))
}

// FilenameGT applies the GT predicate on the "filename" field.
func FilenameGT(v string) predicate.Image {
	return predicate.Image(sql.FieldGT(FieldFilename, v))
}

// FilenameGTE applies the GTE predicate on the "filename" field.
func FilenameGTE(v string) predicate.Image {
	return predicate.Image(sql.FieldGTE(FieldFilename, v))
}

// FilenameLT applies the LT predicate on the "filename" field.
func FilenameLT(v string) predicate.Image {
	return predicate.Image(sql.FieldLT(FieldFilename, v))
}

// FilenameLTE applies the LTE predicate on the "filename" field.
func FilenameLTE(v string) predicate.Image {
	return predicate.Image(sql.FieldLTE(FieldFilename, v))
}

// FilenameContains applies the Contains predicate on the "filename" field.
func FilenameContains(v string) predicate.Image {
	return predicate.Image(sql.FieldContains(FieldFilename, v))
}

// FilenameHasPrefix applies the HasPrefix predicate on the "filename" field.
func FilenameHasPrefix(v string) predicate.Image {
	return predicate.Image(sql.FieldHasPrefix(FieldFilename, v))
}

// FilenameHasSuffix applies the HasSuffix predicate on the "filename" field.
func FilenameHasSuffix(v string) predicate.Image {
	return predicate.Image(sql.FieldHasSuffix(FieldFilename, v))
}

// FilenameEqualFold applies the EqualFold predicate on the "filename" field.
func FilenameEqualFold(v string) predicate.Image {
	return predicate.Image(sql.FieldEqualFold(FieldFilename, v))
}

// FilenameContainsFold applies the ContainsFold predicate on the "filename" field.
func FilenameContainsFold(v string) predicate.Image {
	return predicate.Image(sql.FieldContainsFold(FieldFilename, v))
}

// URLEQ applies the EQ predicate on the "url" field.
func URLEQ(v string) predicate.Image {
	return predicate.Image(sql.FieldEQ(FieldURL, v))
}

// URLNEQ applies the NEQ predicate on the "url" field.
func URLNEQ(v string) predicate.Image {
	return predicate.Image(sql.FieldNEQ(FieldURL, v))
}

// URLIn applies the In predicate on the "url" field.
func URLIn(vs ...string) predicate.Image {
	return predicate.Image(sql.FieldIn(FieldURL, vs...))
}

// URLNotIn applies the NotIn predicate on the "url" field.
func URLNotIn(vs ...string) predicate.Image {
	return predicate.Image(sql.FieldNotIn(FieldURL, vs...))
}

// URLGT applies the GT predicate on the "url" field.
func URLGT(v string) predicate.Image {
	return predicate.Image(sql.FieldGT(FieldURL, v))
}

// URLGTE applies the GTE predicate on the "url" field.
func URLGTE(v string) predicate.Image {
	return predicate.Image(sql.FieldGTE(FieldURL, v))
}

// URLLT applies the LT predicate on the "url" field.
func URLLT(v string) predicate.Image {
	return predicate.Image(sql.FieldLT(FieldURL, v))
}

// URLLTE applies the LTE predicate on the "url" field.
func URLLTE(v string) predicate.Image {
	return predicate.Image(sql.FieldLTE(FieldURL, v))
}

// URLContains applies the Contains predicate on the "url" field.
func URLContains(v string) predicate.Image {
	return predicate.Image(sql.FieldContains(FieldURL, v))
}

// URLHasPrefix applies the HasPrefix predicate on the "url" field.
func URLHasPrefix(v string) predicate.Image {
	return predicate.Image(sql.FieldHasPrefix(FieldURL, v))
}

// URLHasSuffix applies the HasSuffix predicate on the "url" field.
func URLHasSuffix(v string) predicate.Image {
	return predicate.Image(sql.FieldHasSuffix(FieldURL, v))
}

// URLEqualFold applies the EqualFold predicate on the "url" field.
func URLEqualFold(v string) predicate.Image {
	return predicate.Image(sql.FieldEqualFold(FieldURL, v))
}

// URLContainsFold applies the ContainsFold predicate on the "url" field.
func URLContainsFold(v string) predicate.Image {
	return predicate.Image(sql.FieldContainsFold(FieldURL, v))
}

// SizeEQ applies the EQ predicate on the "size" field.
func SizeEQ(v int64) predicate.Image {
	return predicate.Image(sql.FieldEQ(FieldSize, v))
}

// SizeNEQ applies the NEQ predicate on the "size" field.
func SizeNEQ(v int64) predicate.Image {
	return predicate.Image(sql.FieldNEQ(FieldSize, v))
}

// SizeIn applies the In predicate on the "size" field.
func SizeIn(vs ...int64) predicate.Image {
	return predicate.Image(sql.FieldIn(FieldSize, vs...))
}

// SizeNotIn applies the NotIn predicate on the "size" field.
func SizeNotIn(vs ...int64) predicate.Image {
	return predicate.Image(sql.FieldNotIn(FieldSize, vs...))
}

// SizeGT applies the GT predicate on the "size" field.
func SizeGT(v int64) predicate.Image {
	return predicate.Image(sql.FieldGT(FieldSize, v))
}

// SizeGTE applies the GTE predicate on the "size" field.
func SizeGTE(v int64) predicate.Image {
	return predicate.Image(sql.FieldGTE(FieldSize, v))
}

// SizeLT applies the LT predicate on the "size" field.
func SizeLT(v int64) predicate.Image {
	return predicate.Image(sql.FieldLT(FieldSize, v))
}

// SizeLTE applies the LTE predicate on the "size" field.
func SizeLTE(v int64) predicate.Image {
	return predicate.Image(sql.FieldLTE(FieldSize, v))
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v string) predicate.Image {
	return predicate.Image(sql.FieldEQ(FieldType, v))
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v string) predicate.Image {
	return predicate.Image(sql.FieldNEQ(FieldType, v))
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...string) predicate.Image {
	return predicate.Image(sql.FieldIn(FieldType, vs...))
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...string) predicate.Image {
	return predicate.Image(sql.FieldNotIn(FieldType, vs...))
}

// TypeGT applies the GT predicate on the "type" field.
func TypeGT(v string) predicate.Image {
	return predicate.Image(sql.FieldGT(FieldType, v))
}

// TypeGTE applies the GTE predicate on the "type" field.
func TypeGTE(v string) predicate.Image {
	return predicate.Image(sql.FieldGTE(FieldType, v))
}

// TypeLT applies the LT predicate on the "type" field.
func TypeLT(v string) predicate.Image {
	return predicate.Image(sql.FieldLT(FieldType, v))
}

// TypeLTE applies the LTE predicate on the "type" field.
func TypeLTE(v string) predicate.Image {
	return predicate.Image(sql.FieldLTE(FieldType, v))
}

// TypeContains applies the Contains predicate on the "type" field.
func TypeContains(v string) predicate.Image {
	return predicate.Image(sql.FieldContains(FieldType, v))
}

// TypeHasPrefix applies the HasPrefix predicate on the "type" field.
func TypeHasPrefix(v string) predicate.Image {
	return predicate.Image(sql.FieldHasPrefix(FieldType, v))
}

// TypeHasSuffix applies the HasSuffix predicate on the "type" field.
func TypeHasSuffix(v string) predicate.Image {
	return predicate.Image(sql.FieldHasSuffix(FieldType, v))
}

// TypeEqualFold applies the EqualFold predicate on the "type" field.
func TypeEqualFold(v string) predicate.Image {
	return predicate.Image(sql.FieldEqualFold(FieldType, v))
}

// TypeContainsFold applies the ContainsFold predicate on the "type" field.
func TypeContainsFold(v string) predicate.Image {
	return predicate.Image(sql.FieldContainsFold(FieldType, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Image {
	return predicate.Image(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Image {
	return predicate.Image(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Image {
	return predicate.Image(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Image {
	return predicate.Image(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Image {
	return predicate.Image(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Image {
	return predicate.Image(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Image {
	return predicate.Image(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Image {
	return predicate.Image(sql.FieldLTE(FieldCreatedAt, v))
}

// WidthEQ applies the EQ predicate on the "width" field.
func WidthEQ(v int) predicate.Image {
	return predicate.Image(sql.FieldEQ(FieldWidth, v))
}

// WidthNEQ applies the NEQ predicate on the "width" field.
func WidthNEQ(v int) predicate.Image {
	return predicate.Image(sql.FieldNEQ(FieldWidth, v))
}

// WidthIn applies the In predicate on the "width" field.
func WidthIn(vs ...int) predicate.Image {
	return predicate.Image(sql.FieldIn(FieldWidth, vs...))
}

// WidthNotIn applies the NotIn predicate on the "width" field.
func WidthNotIn(vs ...int) predicate.Image {
	return predicate.Image(sql.FieldNotIn(FieldWidth, vs...))
}

// WidthGT applies the GT predicate on the "width" field.
func WidthGT(v int) predicate.Image {
	return predicate.Image(sql.FieldGT(FieldWidth, v))
}

// WidthGTE applies the GTE predicate on the "width" field.
func WidthGTE(v int) predicate.Image {
	return predicate.Image(sql.FieldGTE(FieldWidth, v))
}

// WidthLT applies the LT predicate on the "width" field.
func WidthLT(v int) predicate.Image {
	return predicate.Image(sql.FieldLT(FieldWidth, v))
}

// WidthLTE applies the LTE predicate on the "width" field.
func WidthLTE(v int) predicate.Image {
	return predicate.Image(sql.FieldLTE(FieldWidth, v))
}

// HeightEQ applies the EQ predicate on the "height" field.
func HeightEQ(v int) predicate.Image {
	return predicate.Image(sql.FieldEQ(FieldHeight, v))
}

// HeightNEQ applies the NEQ predicate on the "height" field.
func HeightNEQ(v int) predicate.Image {
	return predicate.Image(sql.FieldNEQ(FieldHeight, v))
}

// HeightIn applies the In predicate on the "height" field.
func HeightIn(vs ...int) predicate.Image {
	return predicate.Image(sql.FieldIn(FieldHeight, vs...))
}

// HeightNotIn applies the NotIn predicate on the "height" field.
func HeightNotIn(vs ...int) predicate.Image {
	return predicate.Image(sql.FieldNotIn(FieldHeight, vs...))
}

// HeightGT applies the GT predicate on the "height" field.
func HeightGT(v int) predicate.Image {
	return predicate.Image(sql.FieldGT(FieldHeight, v))
}

// HeightGTE applies the GTE predicate on the "height" field.
func HeightGTE(v int) predicate.Image {
	return predicate.Image(sql.FieldGTE(FieldHeight, v))
}

// HeightLT applies the LT predicate on the "height" field.
func HeightLT(v int) predicate.Image {
	return predicate.Image(sql.FieldLT(FieldHeight, v))
}

// HeightLTE applies the LTE predicate on the "height" field.
func HeightLTE(v int) predicate.Image {
	return predicate.Image(sql.FieldLTE(FieldHeight, v))
}

// HasUploadedBy applies the HasEdge predicate on the "uploaded_by" edge.
func HasUploadedBy() predicate.Image {
	return predicate.Image(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, UploadedByTable, UploadedByColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUploadedByWith applies the HasEdge predicate on the "uploaded_by" edge with a given conditions (other predicates).
func HasUploadedByWith(preds ...predicate.User) predicate.Image {
	return predicate.Image(func(s *sql.Selector) {
		step := newUploadedByStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasBooks applies the HasEdge predicate on the "books" edge.
func HasBooks() predicate.Image {
	return predicate.Image(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, BooksTable, BooksColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasBooksWith applies the HasEdge predicate on the "books" edge with a given conditions (other predicates).
func HasBooksWith(preds ...predicate.Book) predicate.Image {
	return predicate.Image(func(s *sql.Selector) {
		step := newBooksStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Image) predicate.Image {
	return predicate.Image(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Image) predicate.Image {
	return predicate.Image(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Image) predicate.Image {
	return predicate.Image(sql.NotPredicates(p))
}
