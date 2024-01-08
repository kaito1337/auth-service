package controllers

import "github.com/gofiber/fiber/v2"

type Controller interface {
	GetGroup() string // return path
	GetHandlers() []ControllerHandler
}

type ControllerHandler interface {
	GetMethod() string
	GetPath() string
	GetHandler() func(c *fiber.Ctx) error
}

type Handler struct {
	Method  string
	Path    string
	Handler func(c *fiber.Ctx) error
}

func (h *Handler) GetPath() string {
	return h.Path
}

func (h *Handler) GetMethod() string {
	return h.Method
}
func (h *Handler) GetHandler() func(c *fiber.Ctx) error {
	return h.Handler
}
