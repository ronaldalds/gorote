package example

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	Service *Service
}

func NewController() *Controller {
	return &Controller{
		Service: NewService(),
	}
}