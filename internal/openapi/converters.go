package openapi

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/muhomorfus/lab1-template/internal/generated"
	"github.com/muhomorfus/lab1-template/internal/models"
	"github.com/samber/lo"
)

func toPerson(p generated.PersonRequest) models.Person {
	return models.Person{
		Name:    p.Name,
		Address: p.Address,
		Age:     toIntPtr(p.Age),
		Work:    p.Work,
	}
}

func fromPerson(p models.Person) generated.PersonResponse {
	return generated.PersonResponse{
		Address: p.Address,
		Age:     toInt32Ptr(p.Age),
		Id:      int32(p.ID),
		Name:    p.Name,
		Work:    p.Work,
	}
}

func toIntPtr(ptr *int32) *int {
	if ptr == nil {
		return nil
	}

	return lo.ToPtr(int(*ptr))
}

func toInt32Ptr(ptr *int) *int32 {
	if ptr == nil {
		return nil
	}

	return lo.ToPtr(int32(*ptr))
}

func errorResponse(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, models.ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(generated.ErrorResponse{Message: lo.ToPtr(err.Error())})
	case errors.Is(err, models.ErrInvalidData):
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			m := make(map[string]string)
			for _, e := range validationErrors {
				m[e.Field()] = e.Error()
			}

			return c.Status(fiber.StatusBadRequest).JSON(generated.ValidationErrorResponse{
				Errors:  lo.ToPtr(m),
				Message: lo.ToPtr(err.Error()),
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(generated.ValidationErrorResponse{
			Message: lo.ToPtr(err.Error()),
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(generated.ErrorResponse{Message: lo.ToPtr(err.Error())})
	}
}
