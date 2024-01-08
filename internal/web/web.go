package web

import (
	"auth-backend/internal/config"
	"auth-backend/internal/web/controllers"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
	"log"
	"reflect"
)

type WebServer struct {
	cfg    *config.WebServerConfig
	client *fiber.App
}

func NewWebServer(cfg *config.WebServerConfig) *WebServer {
	return &WebServer{
		cfg:    cfg,
		client: fiber.New(fiber.Config{}),
	}
}

func (w *WebServer) RegisterRoutes(routes []controllers.Controller) {

	for _, route := range routes {

		group := w.client.Group(route.GetGroup())

		for _, handler := range route.GetHandlers() {
			switch handler.GetMethod() {
			case "GET":
				group.Get(handler.GetPath(), handler.GetHandler())
			case "POST":
				group.Post(handler.GetPath(), handler.GetHandler())
			case "PUT":
				group.Put(handler.GetPath(), handler.GetHandler())
			case "PATCH":
				group.Patch(handler.GetPath(), handler.GetHandler())
			case "DELETE":
				group.Delete(handler.GetPath(), handler.GetHandler())
			default:
				fmt.Printf("unsupported method: %s, path: %s controller: %s", handler.GetMethod(), handler.GetPath(), reflect.TypeOf(handler).Elem().Name())
			}
		}
	}

}

func (w *WebServer) Run() {
	err := w.client.Listen(fmt.Sprintf(":%d", w.cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
}

func (w *WebServer) LogMiddleware(channel *amqp.Channel) {
	w.client.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return err
		}
		err = channel.Publish(
			"",       // exchange
			"logger", // routing key
			false,    // mandatory
			false,    // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        c.Response().Body(),
			})
		return nil
	})
}
