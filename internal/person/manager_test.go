package person

import (
	"context"
	"errors"
	"github.com/muhomorfus/lab1-template/internal/models"
	"github.com/muhomorfus/lab1-template/internal/person/mocks"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestManager_Create(t *testing.T) {
	t.Run("positive: created person", func(t *testing.T) {
		ctx := context.Background()

		toCreate := models.Person{
			ID:   0,
			Name: "Vasya",
		}

		repoMock := mocks.NewPersonRepository(t)
		repoMock.EXPECT().Create(ctx, toCreate).Return(1, nil)

		m := New(repoMock)

		got, err := m.Create(ctx, toCreate)
		require.NoError(t, err)
		assert.Equal(t, 1, got)
	})

	t.Run("negative: cant create in repo", func(t *testing.T) {
		ctx := context.Background()

		toCreate := models.Person{
			ID:   0,
			Name: "Vasya",
		}

		repoMock := mocks.NewPersonRepository(t)
		repoMock.EXPECT().Create(ctx, toCreate).Return(0, errors.New("error"))

		m := New(repoMock)

		got, err := m.Create(ctx, toCreate)
		require.Error(t, err)
		assert.Zero(t, got)
	})

	t.Run("negative: invalid person", func(t *testing.T) {
		ctx := context.Background()

		toCreate := models.Person{
			ID:   0,
			Name: "",
		}

		repoMock := mocks.NewPersonRepository(t)

		m := New(repoMock)

		got, err := m.Create(ctx, toCreate)
		require.Error(t, err)
		assert.Zero(t, got)
	})
}

func TestManager_Get(t *testing.T) {
	t.Run("positive: created person", func(t *testing.T) {
		ctx := context.Background()

		person := &models.Person{
			ID:   228,
			Name: "Vasya",
		}

		repoMock := mocks.NewPersonRepository(t)
		repoMock.EXPECT().Get(ctx, 228).Return(person, nil)

		m := New(repoMock)

		got, err := m.Get(ctx, 228)
		require.NoError(t, err)
		assert.Equal(t, person, got)
	})

	t.Run("positive: cant get user from repo", func(t *testing.T) {
		ctx := context.Background()

		repoMock := mocks.NewPersonRepository(t)
		repoMock.EXPECT().Get(ctx, 228).Return(nil, errors.New("error"))

		m := New(repoMock)

		got, err := m.Get(ctx, 228)
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestManager_Update(t *testing.T) {
	t.Run("positive: created person", func(t *testing.T) {
		ctx := context.Background()

		person := models.Person{
			ID:   228,
			Name: "Petya",
		}

		patch := models.Person{
			ID:   228,
			Name: "Vasya",
			Age:  lo.ToPtr(18),
		}

		toUpdate := models.Person{
			ID:   228,
			Name: "Vasya",
			Age:  lo.ToPtr(18),
		}

		repoMock := mocks.NewPersonRepository(t)
		repoMock.EXPECT().Get(ctx, 228).Return(&person, nil)
		repoMock.EXPECT().Update(ctx, toUpdate).Return(nil)

		m := New(repoMock)

		got, err := m.Update(ctx, patch)
		require.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("positive: cant update user", func(t *testing.T) {
		ctx := context.Background()

		person := models.Person{
			ID:   228,
			Name: "Petya",
		}

		patch := models.Person{
			ID:   228,
			Name: "Vasya",
			Age:  lo.ToPtr(18),
		}

		toUpdate := models.Person{
			ID:   228,
			Name: "Vasya",
			Age:  lo.ToPtr(18),
		}

		repoMock := mocks.NewPersonRepository(t)
		repoMock.EXPECT().Get(ctx, 228).Return(&person, nil)
		repoMock.EXPECT().Update(ctx, toUpdate).Return(errors.New("error"))

		m := New(repoMock)

		got, err := m.Update(ctx, patch)
		require.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("positive: cant get user", func(t *testing.T) {
		ctx := context.Background()

		patch := models.Person{
			ID:   228,
			Name: "Vasya",
			Age:  lo.ToPtr(18),
		}

		repoMock := mocks.NewPersonRepository(t)
		repoMock.EXPECT().Get(ctx, 228).Return(nil, errors.New("error"))

		m := New(repoMock)

		got, err := m.Update(ctx, patch)
		require.Error(t, err)
		assert.Nil(t, got)
	})
}
