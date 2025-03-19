package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type ErrHandler struct {
	Error []DetailErr `json:"error"`
}

type DetailErr struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewError(message ...string) *ErrHandler {
	err := &ErrHandler{
		Error: nil,
	}
	if len(message) > 0 {
		err.AddDetailErr("err", message[0])
	}
	return err
}

func (e *ErrHandler) AddDetailErr(key string, value string) {
	e.Error = append(e.Error, DetailErr{
		Key:   key,
		Value: value,
	})
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}
	res := NewError(err.Error())
	return ctx.Status(code).JSON(res)
}
