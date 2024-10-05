package openapi

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muhomorfus/lab1-template/internal/generated"
	"github.com/muhomorfus/lab1-template/internal/models"
	"github.com/samber/lo"
)

type Server struct {
	mgr personManager
}

func New(mgr personManager) *Server {
	return &Server{mgr: mgr}
}

func (s *Server) ListPersons(c *fiber.Ctx) error {
	persons, err := s.mgr.List(c.Context())
	if err != nil {
		return errorResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(lo.Map(persons, func(p models.Person, _ int) generated.PersonResponse {
		return fromPerson(p)
	}))
}

func (s *Server) CreatePerson(c *fiber.Ctx) error {
	var req generated.PersonRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, err)
	}

	id, err := s.mgr.Create(c.Context(), toPerson(req))
	if err != nil {
		return errorResponse(c, err)
	}

	c.Set("Location", fmt.Sprintf("/api/v1/persons/%d", id))
	return c.SendStatus(fiber.StatusCreated)
}

func (s *Server) DeletePerson(c *fiber.Ctx, id int32) error {
	if err := s.mgr.Delete(c.Context(), int(id)); err != nil {
		return errorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (s *Server) GetPerson(c *fiber.Ctx, id int32) error {
	person, err := s.mgr.Get(c.Context(), int(id))
	if err != nil {
		return errorResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fromPerson(*person))
}

func (s *Server) EditPerson(c *fiber.Ctx, id int32) error {
	var req generated.PersonRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, err)
	}

	person := toPerson(req)
	person.ID = int(id)

	updated, err := s.mgr.Update(c.Context(), person)
	if err != nil {
		return errorResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fromPerson(*updated))
}

type personManager interface {
	Create(ctx context.Context, person models.Person) (int, error)
	Get(ctx context.Context, id int) (*models.Person, error)
	Update(ctx context.Context, person models.Person) (*models.Person, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]models.Person, error)
}
