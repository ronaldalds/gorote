package middlewares

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ronaldalds/base-go-api/internal/config/handlers"
)

func (m *Middleware) ValidationMiddleware(requestStruct any, inputType string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		switch inputType {
		case "query":
			if err := ctx.QueryParser(requestStruct); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid query parameters: %s", err.Error()))
			}
		case "json":
			if err := ctx.BodyParser(requestStruct); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid body: %s", err.Error()))
			}
		case "params":
			if err := ctx.ParamsParser(requestStruct); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid URL parameters: %s", err.Error()))
			}
		default:
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid validation type"))
		}

		// Valide os dados usando o validator
		if err := validateStruct(requestStruct); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(err)
		}

		// Armazene os dados validados no contexto
		ctx.Locals("validatedData", requestStruct)

		// Prossiga para o próximo middleware ou handler
		return ctx.Next()
	}
}

func validateStruct(data any) *handlers.ErrHandler {
	var validate = validator.New()
	// Verifica se o objeto possui erros de validação
	err := validate.Struct(data)
	if err != nil {
		// Converte o erro para ValidationErrors, se aplicável
		errors := handlers.NewError()
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, err := range validationErrors {
				fieldName := err.Field()
				tag := err.Tag()
				errors.AddDetailErr(fieldName, fmt.Sprintf("invalid field: %s", tag))
			}
			return errors
		}
		// Retorna erro genérico se não for ValidationErrors
		errors.AddDetailErr("error", "validate Structure Failed!")
		return errors
	}
	// Nenhum erro encontrado
	return nil
}
