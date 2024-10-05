package person

import (
	"context"
	"fmt"
	"github.com/muhomorfus/lab1-template/internal/models"
)

type Manager struct {
	repo personRepository
}

func New(repo personRepository) *Manager {
	return &Manager{repo: repo}
}

func (m *Manager) Create(ctx context.Context, person models.Person) (int, error) {
	if err := person.Validate(); err != nil {
		return 0, fmt.Errorf("invalid person to create: %w", err)
	}

	id, err := m.repo.Create(ctx, person)
	if err != nil {
		return 0, fmt.Errorf("create user by repo: %w", err)
	}

	return id, nil
}

func (m *Manager) Get(ctx context.Context, id int) (*models.Person, error) {
	person, err := m.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return person, nil
}

func (m *Manager) Update(ctx context.Context, patch models.Person) (*models.Person, error) {
	person, err := m.repo.Get(ctx, patch.ID)
	if err != nil {
		return nil, fmt.Errorf("get user from repo: %w", err)
	}

	person.Merge(patch)

	if err := person.Validate(); err != nil {
		return nil, fmt.Errorf("invalid person to update: %w", err)
	}

	if err := m.repo.Update(ctx, *person); err != nil {
		return nil, fmt.Errorf("update user in repo: %w", err)
	}

	return person, nil
}

func (m *Manager) Delete(ctx context.Context, id int) error {
	if err := m.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete user by id: %w", err)
	}

	return nil
}

func (m *Manager) List(ctx context.Context) ([]models.Person, error) {
	persons, err := m.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	return persons, nil
}

//go:generate mockery --all --with-expecter --exported --output mocks/

type personRepository interface {
	Create(ctx context.Context, person models.Person) (int, error)
	Get(ctx context.Context, id int) (*models.Person, error)
	Update(ctx context.Context, person models.Person) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]models.Person, error)
}
