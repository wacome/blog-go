// Code generated by ent, DO NOT EDIT.

package ent

import (
	"blog-go/ent/book"
	"blog-go/ent/image"
	"blog-go/ent/user"
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// BookCreate is the builder for creating a Book entity.
type BookCreate struct {
	config
	mutation *BookMutation
	hooks    []Hook
}

// SetTitle sets the "title" field.
func (bc *BookCreate) SetTitle(s string) *BookCreate {
	bc.mutation.SetTitle(s)
	return bc
}

// SetAuthor sets the "author" field.
func (bc *BookCreate) SetAuthor(s string) *BookCreate {
	bc.mutation.SetAuthor(s)
	return bc
}

// SetDesc sets the "desc" field.
func (bc *BookCreate) SetDesc(s string) *BookCreate {
	bc.mutation.SetDesc(s)
	return bc
}

// SetNillableDesc sets the "desc" field if the given value is not nil.
func (bc *BookCreate) SetNillableDesc(s *string) *BookCreate {
	if s != nil {
		bc.SetDesc(*s)
	}
	return bc
}

// SetCover sets the "cover" field.
func (bc *BookCreate) SetCover(s string) *BookCreate {
	bc.mutation.SetCover(s)
	return bc
}

// SetNillableCover sets the "cover" field if the given value is not nil.
func (bc *BookCreate) SetNillableCover(s *string) *BookCreate {
	if s != nil {
		bc.SetCover(*s)
	}
	return bc
}

// SetPublisher sets the "publisher" field.
func (bc *BookCreate) SetPublisher(s string) *BookCreate {
	bc.mutation.SetPublisher(s)
	return bc
}

// SetNillablePublisher sets the "publisher" field if the given value is not nil.
func (bc *BookCreate) SetNillablePublisher(s *string) *BookCreate {
	if s != nil {
		bc.SetPublisher(*s)
	}
	return bc
}

// SetPublishDate sets the "publish_date" field.
func (bc *BookCreate) SetPublishDate(s string) *BookCreate {
	bc.mutation.SetPublishDate(s)
	return bc
}

// SetNillablePublishDate sets the "publish_date" field if the given value is not nil.
func (bc *BookCreate) SetNillablePublishDate(s *string) *BookCreate {
	if s != nil {
		bc.SetPublishDate(*s)
	}
	return bc
}

// SetIsbn sets the "isbn" field.
func (bc *BookCreate) SetIsbn(s string) *BookCreate {
	bc.mutation.SetIsbn(s)
	return bc
}

// SetNillableIsbn sets the "isbn" field if the given value is not nil.
func (bc *BookCreate) SetNillableIsbn(s *string) *BookCreate {
	if s != nil {
		bc.SetIsbn(*s)
	}
	return bc
}

// SetPages sets the "pages" field.
func (bc *BookCreate) SetPages(i int) *BookCreate {
	bc.mutation.SetPages(i)
	return bc
}

// SetNillablePages sets the "pages" field if the given value is not nil.
func (bc *BookCreate) SetNillablePages(i *int) *BookCreate {
	if i != nil {
		bc.SetPages(*i)
	}
	return bc
}

// SetRating sets the "rating" field.
func (bc *BookCreate) SetRating(f float64) *BookCreate {
	bc.mutation.SetRating(f)
	return bc
}

// SetNillableRating sets the "rating" field if the given value is not nil.
func (bc *BookCreate) SetNillableRating(f *float64) *BookCreate {
	if f != nil {
		bc.SetRating(*f)
	}
	return bc
}

// SetStatus sets the "status" field.
func (bc *BookCreate) SetStatus(b book.Status) *BookCreate {
	bc.mutation.SetStatus(b)
	return bc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (bc *BookCreate) SetNillableStatus(b *book.Status) *BookCreate {
	if b != nil {
		bc.SetStatus(*b)
	}
	return bc
}

// SetReview sets the "review" field.
func (bc *BookCreate) SetReview(s string) *BookCreate {
	bc.mutation.SetReview(s)
	return bc
}

// SetNillableReview sets the "review" field if the given value is not nil.
func (bc *BookCreate) SetNillableReview(s *string) *BookCreate {
	if s != nil {
		bc.SetReview(*s)
	}
	return bc
}

// SetCreatedAt sets the "created_at" field.
func (bc *BookCreate) SetCreatedAt(t time.Time) *BookCreate {
	bc.mutation.SetCreatedAt(t)
	return bc
}

// SetUpdatedAt sets the "updated_at" field.
func (bc *BookCreate) SetUpdatedAt(t time.Time) *BookCreate {
	bc.mutation.SetUpdatedAt(t)
	return bc
}

// SetCoverImageID sets the "cover_image" edge to the Image entity by ID.
func (bc *BookCreate) SetCoverImageID(id int) *BookCreate {
	bc.mutation.SetCoverImageID(id)
	return bc
}

// SetNillableCoverImageID sets the "cover_image" edge to the Image entity by ID if the given value is not nil.
func (bc *BookCreate) SetNillableCoverImageID(id *int) *BookCreate {
	if id != nil {
		bc = bc.SetCoverImageID(*id)
	}
	return bc
}

// SetCoverImage sets the "cover_image" edge to the Image entity.
func (bc *BookCreate) SetCoverImage(i *Image) *BookCreate {
	return bc.SetCoverImageID(i.ID)
}

// SetOwnerID sets the "owner" edge to the User entity by ID.
func (bc *BookCreate) SetOwnerID(id int) *BookCreate {
	bc.mutation.SetOwnerID(id)
	return bc
}

// SetOwner sets the "owner" edge to the User entity.
func (bc *BookCreate) SetOwner(u *User) *BookCreate {
	return bc.SetOwnerID(u.ID)
}

// Mutation returns the BookMutation object of the builder.
func (bc *BookCreate) Mutation() *BookMutation {
	return bc.mutation
}

// Save creates the Book in the database.
func (bc *BookCreate) Save(ctx context.Context) (*Book, error) {
	bc.defaults()
	return withHooks(ctx, bc.sqlSave, bc.mutation, bc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (bc *BookCreate) SaveX(ctx context.Context) *Book {
	v, err := bc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (bc *BookCreate) Exec(ctx context.Context) error {
	_, err := bc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bc *BookCreate) ExecX(ctx context.Context) {
	if err := bc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (bc *BookCreate) defaults() {
	if _, ok := bc.mutation.Cover(); !ok {
		v := book.DefaultCover
		bc.mutation.SetCover(v)
	}
	if _, ok := bc.mutation.Rating(); !ok {
		v := book.DefaultRating
		bc.mutation.SetRating(v)
	}
	if _, ok := bc.mutation.Status(); !ok {
		v := book.DefaultStatus
		bc.mutation.SetStatus(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (bc *BookCreate) check() error {
	if _, ok := bc.mutation.Title(); !ok {
		return &ValidationError{Name: "title", err: errors.New(`ent: missing required field "Book.title"`)}
	}
	if v, ok := bc.mutation.Title(); ok {
		if err := book.TitleValidator(v); err != nil {
			return &ValidationError{Name: "title", err: fmt.Errorf(`ent: validator failed for field "Book.title": %w`, err)}
		}
	}
	if _, ok := bc.mutation.Author(); !ok {
		return &ValidationError{Name: "author", err: errors.New(`ent: missing required field "Book.author"`)}
	}
	if v, ok := bc.mutation.Author(); ok {
		if err := book.AuthorValidator(v); err != nil {
			return &ValidationError{Name: "author", err: fmt.Errorf(`ent: validator failed for field "Book.author": %w`, err)}
		}
	}
	if _, ok := bc.mutation.Cover(); !ok {
		return &ValidationError{Name: "cover", err: errors.New(`ent: missing required field "Book.cover"`)}
	}
	if _, ok := bc.mutation.Rating(); !ok {
		return &ValidationError{Name: "rating", err: errors.New(`ent: missing required field "Book.rating"`)}
	}
	if _, ok := bc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "Book.status"`)}
	}
	if v, ok := bc.mutation.Status(); ok {
		if err := book.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Book.status": %w`, err)}
		}
	}
	if _, ok := bc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Book.created_at"`)}
	}
	if _, ok := bc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Book.updated_at"`)}
	}
	if _, ok := bc.mutation.OwnerID(); !ok {
		return &ValidationError{Name: "owner", err: errors.New(`ent: missing required edge "Book.owner"`)}
	}
	return nil
}

func (bc *BookCreate) sqlSave(ctx context.Context) (*Book, error) {
	if err := bc.check(); err != nil {
		return nil, err
	}
	_node, _spec := bc.createSpec()
	if err := sqlgraph.CreateNode(ctx, bc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	bc.mutation.id = &_node.ID
	bc.mutation.done = true
	return _node, nil
}

func (bc *BookCreate) createSpec() (*Book, *sqlgraph.CreateSpec) {
	var (
		_node = &Book{config: bc.config}
		_spec = sqlgraph.NewCreateSpec(book.Table, sqlgraph.NewFieldSpec(book.FieldID, field.TypeInt))
	)
	if value, ok := bc.mutation.Title(); ok {
		_spec.SetField(book.FieldTitle, field.TypeString, value)
		_node.Title = value
	}
	if value, ok := bc.mutation.Author(); ok {
		_spec.SetField(book.FieldAuthor, field.TypeString, value)
		_node.Author = value
	}
	if value, ok := bc.mutation.Desc(); ok {
		_spec.SetField(book.FieldDesc, field.TypeString, value)
		_node.Desc = value
	}
	if value, ok := bc.mutation.Cover(); ok {
		_spec.SetField(book.FieldCover, field.TypeString, value)
		_node.Cover = value
	}
	if value, ok := bc.mutation.Publisher(); ok {
		_spec.SetField(book.FieldPublisher, field.TypeString, value)
		_node.Publisher = value
	}
	if value, ok := bc.mutation.PublishDate(); ok {
		_spec.SetField(book.FieldPublishDate, field.TypeString, value)
		_node.PublishDate = value
	}
	if value, ok := bc.mutation.Isbn(); ok {
		_spec.SetField(book.FieldIsbn, field.TypeString, value)
		_node.Isbn = value
	}
	if value, ok := bc.mutation.Pages(); ok {
		_spec.SetField(book.FieldPages, field.TypeInt, value)
		_node.Pages = value
	}
	if value, ok := bc.mutation.Rating(); ok {
		_spec.SetField(book.FieldRating, field.TypeFloat64, value)
		_node.Rating = value
	}
	if value, ok := bc.mutation.Status(); ok {
		_spec.SetField(book.FieldStatus, field.TypeEnum, value)
		_node.Status = value
	}
	if value, ok := bc.mutation.Review(); ok {
		_spec.SetField(book.FieldReview, field.TypeString, value)
		_node.Review = value
	}
	if value, ok := bc.mutation.CreatedAt(); ok {
		_spec.SetField(book.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := bc.mutation.UpdatedAt(); ok {
		_spec.SetField(book.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if nodes := bc.mutation.CoverImageIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   book.CoverImageTable,
			Columns: []string{book.CoverImageColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(image.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.image_books = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := bc.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   book.OwnerTable,
			Columns: []string{book.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.user_books = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// BookCreateBulk is the builder for creating many Book entities in bulk.
type BookCreateBulk struct {
	config
	err      error
	builders []*BookCreate
}

// Save creates the Book entities in the database.
func (bcb *BookCreateBulk) Save(ctx context.Context) ([]*Book, error) {
	if bcb.err != nil {
		return nil, bcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(bcb.builders))
	nodes := make([]*Book, len(bcb.builders))
	mutators := make([]Mutator, len(bcb.builders))
	for i := range bcb.builders {
		func(i int, root context.Context) {
			builder := bcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*BookMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, bcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, bcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, bcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (bcb *BookCreateBulk) SaveX(ctx context.Context) []*Book {
	v, err := bcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (bcb *BookCreateBulk) Exec(ctx context.Context) error {
	_, err := bcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bcb *BookCreateBulk) ExecX(ctx context.Context) {
	if err := bcb.Exec(ctx); err != nil {
		panic(err)
	}
}
