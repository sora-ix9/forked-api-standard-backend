package handlers_test

import (
	"fdlp-standard-api/internal/handler"
	mock_services "fdlp-standard-api/internal/tests/mock/services"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWebSocketHandler_WebSocketInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock_services.NewMockWebSocketService(ctrl)
	realUpgrader := websocket.Upgrader{}
	h := handler.NewWebsocketHandler(mockSvc, realUpgrader)
	e := echo.New()

	// 1. Upgrade Failure
	t.Run("UpgradeError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ws", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.WebSocketInit(c)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Could not upgrade")
	})

	// 2. Success Flow
	t.Run("Success", func(t *testing.T) {
		// We use a channel to signal completion of handler
		done := make(chan bool)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := e.NewContext(r, w)

			// Expectation: Register is called
			mockSvc.EXPECT().RegisterClient(gomock.Any())
			// Expectation: Unregister is called when connection closes
			mockSvc.EXPECT().UnregisterClient(gomock.Any())

			err := h.WebSocketInit(c)
			assert.NoError(t, err)
			done <- true
		}))
		defer server.Close()

		u := "ws" + strings.TrimPrefix(server.URL, "http")
		conn, _, err := websocket.DefaultDialer.Dial(u, nil)
		assert.NoError(t, err)

		// Send message
		err = conn.WriteMessage(websocket.TextMessage, []byte("test"))
		assert.NoError(t, err)

		// Close connection to trigger handler exit
		conn.Close()

		// Wait for server handler to finish
		select {
		case <-done:
			// OK
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for handler to return")
		}
	})

	// 3. Read Error (Client disconnects immediately)
	t.Run("ReadError", func(t *testing.T) {
		done := make(chan bool)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := e.NewContext(r, w)

			mockSvc.EXPECT().RegisterClient(gomock.Any())
			mockSvc.EXPECT().UnregisterClient(gomock.Any())

			err := h.WebSocketInit(c)
			assert.NoError(t, err)
			done <- true
		}))
		defer server.Close()

		u := "ws" + strings.TrimPrefix(server.URL, "http")
		conn, _, err := websocket.DefaultDialer.Dial(u, nil)
		assert.NoError(t, err)

		// Immediately close without writing
		conn.Close()

		select {
		case <-done:
			// OK
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for handler to return")
		}
	})
}
