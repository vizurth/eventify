package wsserver

import (
	"context"
	"encoding/json"
	"eventify/common/logger"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type WsServer interface {
	Start() error
	BroadcastMessage(msg NotificationMessage) error
	GetConnectedClientsCount() int
}

type wsServer struct {
	mux     *http.ServeMux
	srv     *http.Server
	wsUpg   *websocket.Upgrader
	log     *logger.Logger
	clients map[string]*ClientConnection
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewWsServer(addr string, log *logger.Logger) WsServer {
	ctx, cancel := context.WithCancel(context.Background())
	mux := http.NewServeMux()
	return &wsServer{
		mux: mux,
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		wsUpg: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Разрешаем все origin для разработки
			},
		},
		log:     log,
		clients: make(map[string]*ClientConnection),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (ws *wsServer) Start() error {
	ws.mux.HandleFunc("/ws", ws.wsHandler)

	// Запускаем горутину для очистки неактивных соединений
	go ws.cleanupInactiveConnections()

	ws.log.Info(ws.ctx, "websocket server starting", zap.String("addr", ws.srv.Addr))
	return ws.srv.ListenAndServe()
}

func (ws *wsServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conn, err := ws.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		ws.log.Error(ctx, "websocket upgrade failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Создаем уникальный ID для клиента
	clientID := conn.RemoteAddr().String()

	// Создаем клиентское соединение
	client := &ClientConnection{
		ID:       clientID,
		Send:     make(chan []byte, 256),
		Conn:     conn,
		LastPing: time.Now(),
	}

	// Добавляем клиента в список
	ws.mu.Lock()
	ws.clients[clientID] = client
	ws.mu.Unlock()

	ws.log.Info(ctx, "websocket connection established",
		zap.String("client_id", clientID),
		zap.Int("total_clients", len(ws.clients)))

	// Запускаем горутины для чтения и записи
	//go ws.handleClientRead(client)
	go ws.handleClientWrite(client)
}

//func (ws *wsServer) handleClientRead(client *ClientConnection) {
//	defer func() {
//		ws.removeClient(client.ID)
//		client.Conn.Close()
//	}()
//
//	client.Conn.SetReadLimit(512)
//	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
//	client.Conn.SetPongHandler(func(string) error {
//		client.LastPing = time.Now()
//		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
//		return nil
//	})
//
//	for {
//		select {
//		case <-ws.ctx.Done():
//			return
//		default:
//			_, message, err := client.Conn.ReadMessage()
//			if err != nil {
//				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
//					ws.log.Error(ws.ctx, "websocket read error", zap.Error(err))
//				}
//				return
//			}
//
//			// Обрабатываем ping/pong
//			if string(message) == "ping" {
//				client.Send <- []byte("pong")
//			}
//		}
//	}
//}

func (ws *wsServer) handleClientWrite(client *ClientConnection) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-ws.ctx.Done():
			return
		}
	}
}

func (ws *wsServer) BroadcastMessage(msg NotificationMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	ws.mu.RLock()
	defer ws.mu.RUnlock()

	// Отправляем сообщение всем подключенным клиентам
	for _, client := range ws.clients {
		select {
		case client.Send <- data:
		default:
			// Если канал полный, удаляем клиента
			ws.log.Warn(ws.ctx, "client channel full, removing client", zap.String("client_id", client.ID))
			go ws.removeClient(client.ID)
		}
	}

	ws.log.Info(ws.ctx, "broadcasted message to clients",
		zap.String("message_type", msg.Type),
		zap.Int("clients_count", len(ws.clients)))

	return nil
}

func (ws *wsServer) removeClient(clientID string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if client, exists := ws.clients[clientID]; exists {
		close(client.Send)
		delete(ws.clients, clientID)
		ws.log.Info(ws.ctx, "client disconnected",
			zap.String("client_id", clientID),
			zap.Int("remaining_clients", len(ws.clients)))
	}
}

func (ws *wsServer) cleanupInactiveConnections() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ws.mu.Lock()
			now := time.Now()
			for clientID, client := range ws.clients {
				if now.Sub(client.LastPing) > 70*time.Second {
					ws.log.Info(ws.ctx, "removing inactive client", zap.String("client_id", clientID))
					close(client.Send)
					delete(ws.clients, clientID)
				}
			}
			ws.mu.Unlock()
		case <-ws.ctx.Done():
			return
		}
	}
}

func (ws *wsServer) GetConnectedClientsCount() int {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return len(ws.clients)
}
