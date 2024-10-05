package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/muhomorfus/lab1-template/internal/models"
	"github.com/samber/lo"
)

type Repository struct {
	db db
}

func New(db db) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, p models.Person) (int, error) {
	query := `insert into persons(name, age, address, work) values ($1, $2, $3, $4) returning id`

	var id int
	err := r.db.GetContext(ctx, &id, query, p.Name, p.Age, p.Address, p.Work)
	if err != nil {
		return 0, fmt.Errorf("inserting person: %w", err)
	}

	return id, nil
}

func (r *Repository) Get(ctx context.Context, id int) (*models.Person, error) {
	query := `select * from persons where id = $1`

	var p person
	err := r.db.GetContext(ctx, &p, query, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, fmt.Errorf("get person from db: %w", models.ErrNotFound)
	case err != nil:
		return nil, fmt.Errorf("get person from db: %w", err)
	default:
		return lo.ToPtr(toPerson(p)), nil
	}
}

func (r *Repository) Update(ctx context.Context, p models.Person) error {
	query := `select * from persons where id = $1`
	if err := r.db.GetContext(ctx, &person{}, query, p.ID); errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("person not found: %w", models.ErrNotFound)
	}

	query = `update persons set name = $2, age = $3, address = $4, work = $5 where id = $1`
	if _, err := r.db.ExecContext(ctx, query, p.ID, p.Name, p.Age, p.Address, p.Work); err != nil {
		return fmt.Errorf("update person: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `delete from persons where id = $1`
	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("delete person: %w", err)
	}

	return nil
}

func (r *Repository) List(ctx context.Context) ([]models.Person, error) {
	query := `select * from persons`

	var persons []person
	if err := r.db.SelectContext(ctx, &persons, query); err != nil {
		return nil, fmt.Errorf("list persons: %w", err)
	}

	return lo.Map(persons, func(p person, _ int) models.Person {
		return toPerson(p)
	}), nil
}

type db interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type person struct {
	ID      int     `db:"id"`
	Name    string  `validate:"required" db:"name"`
	Address *string `db:"address"`
	Age     *int    `db:"age"`
	Work    *string `db:"work"`
}

func toPerson(p person) models.Person {
	return models.Person{
		ID:      p.ID,
		Name:    p.Name,
		Address: p.Address,
		Age:     p.Age,
		Work:    p.Work,
	}
}
