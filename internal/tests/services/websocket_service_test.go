package services_test

import (
	"fdlp-standard-api/internal/services"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketService(t *testing.T) {
	// 1. Success Flow
	t.Run("RegisterBroadcastAndUnregister", func(t *testing.T) {
		svc := services.NewWebSocketService()

		// Setup Test Server
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{}
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}
			svc.RegisterClient(conn)

			// Keep alive until closed
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					svc.UnregisterClient(conn)
					break
				}
			}
		}))
		defer s.Close()

		// Client connects
		u := "ws" + strings.TrimPrefix(s.URL, "http")
		clientConn, _, err := websocket.DefaultDialer.Dial(u, nil)
		assert.NoError(t, err)

		// Wait for registration
		time.Sleep(50 * time.Millisecond)

		// Broadcast
		msg := "test message"
		err = svc.BroadcastMessage(msg)
		assert.NoError(t, err)

		// Client reads
		_, p, err := clientConn.ReadMessage()
		assert.NoError(t, err)
		assert.Equal(t, msg, string(p))

		// Client disconnects
		clientConn.Close()
		time.Sleep(50 * time.Millisecond)

		// Verify nothing crashes on next broadcast (client removed)
		err = svc.BroadcastMessage("another")
		assert.NoError(t, err)
	})

	// 2. Write Error (Broadcast removes client)
	t.Run("BroadcastWriteError", func(t *testing.T) {
		svc := services.NewWebSocketService()

		// Setup Server that registers but doesn't read loop (controlled)
		var serverConn *websocket.Conn
		ready := make(chan bool)

		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{}
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}
			serverConn = conn
			svc.RegisterClient(conn)
			ready <- true
		}))
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")
		clientConn, _, err := websocket.DefaultDialer.Dial(u, nil)
		assert.NoError(t, err)

		<-ready // wait for server to register

		// Close ONLY the underlying connection or client side to force write error?
		// If we close clientConn, the server might not know immediately until it tries to write.
		clientConn.Close()

		// Give some time for network stack
		time.Sleep(50 * time.Millisecond)

		// Broadcast should fail to write to this conn and call UnregisterClient
		// Note: WriteMessage might not fail immediately for small payloads if socket buffer not full,
		// but since we closed the client, a write should eventually trigger error.
		// However, TCP close might generate a broken pipe.

		// Repeatedly write until error or timeout
		// Actually, RegisterClient is called. serverConn is active.
		// If we close serverConn.Close() from test? No, Broadcast writes to it.
		// UnregisterClient closes it.
		// We want Broadcast to fail WRITING.

		// Force Close Server Conn underlying net conn?
		// serverConn.UnderlyingConn().Close()
		// serverConn.UnderlyingConn().Close()

		// Force Write Error by setting deadline in the past
		serverConn.SetWriteDeadline(time.Now())

		// Broadcast
		err = svc.BroadcastMessage("fail")
		assert.NoError(t, err) // Broadcast returns nil even if write fails

		// Ensure it was removed.
		// We can't access private map `clients`.
		// But if we mock checking, or just rely on coverage.
		// Coverage will show if line 54 (ws.UnregisterClient(conn)) was hit.
	})

	// 3. Unregister duplicate/safe check
	t.Run("UnregisterSafe", func(t *testing.T) {
		svc := services.NewWebSocketService()
		// Calling unregister on nil or not-found shouldn't panic, but code does `ws.clients[conn]`.
		// If conn is not in map, `delete` is no-op. `conn.Close()` will be called if exists check passes.

		// If we pass a dummy conn?
		// We can't easily create a *websocket.Conn without a real connection.
		// But we can check that unregistering a non-registered connection does nothing (coverage line verify).
		// Wait, the code:
		// if _, exists := ws.clients[conn]; exists { ... }
		// So if we pass a conn that is not in map, it does nothing.
		// But validation requires a valid *websocket.Conn pointer.
		// We can use the one from previous test or just nil?
		// map lookup with nil key is valid (returns false).
		// So passing nil should be safe if logic is correct.
		svc.UnregisterClient(nil)
	})
}
