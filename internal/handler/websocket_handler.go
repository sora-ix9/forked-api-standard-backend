package handler

import (
	"fdlp-standard-api/internal/services"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler interface {
	WebSocketInit(c echo.Context) error
}

type webSocketHandler struct {
	services services.WebSocketService
	upgrader websocket.Upgrader
}

func NewWebsocketHandler(services services.WebSocketService, upgrader websocket.Upgrader) WebSocketHandler {
	return &webSocketHandler{services: services, upgrader: upgrader}
}

func (h *webSocketHandler) WebSocketInit(c echo.Context) error {
	conn, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not upgrade to websocket")
	}
	defer conn.Close()

	h.services.RegisterClient(conn)
	defer h.services.UnregisterClient(conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading websocket message: %v", err)
			break // Optionally handle different types of errors differently
		}
		log.Printf("Received message: %s", string(msg))
		// Further processing of the message
	}
	return nil
}
