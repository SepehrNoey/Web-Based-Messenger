package handler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/model"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/repository/chatrepo"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/repository/messagerepo"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/infra/http/auth"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/infra/http/request"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var lastRegisteredMessageID = 0

type MessageHandler struct {
	msgRepo       messagerepo.Repository
	chatrepo      chatrepo.Repository
	jwtConfig     *auth.JWTConfig
	upgrader      *websocket.Upgrader
	broadcast     *chan model.Message                    // a channel that every receiving messages to server goes to
	userOpenConns *map[uint64]map[uint64]*websocket.Conn // userID -> a map of (chatID -> *websocket.Conn)
	connsMutex    *sync.Mutex                            // a mutex to change userOpenConns in a consistent way
}

func NewMessageHandler(msgRepo messagerepo.Repository, chatRepo chatrepo.Repository,
	jwtConfig *auth.JWTConfig, upgrader *websocket.Upgrader, broadcast *chan model.Message,
	userOpenConns *map[uint64]map[uint64]*websocket.Conn, connsMutex *sync.Mutex) *MessageHandler {

	return &MessageHandler{
		msgRepo:       msgRepo,
		chatrepo:      chatRepo,
		jwtConfig:     jwtConfig,
		upgrader:      upgrader,
		broadcast:     broadcast,
		userOpenConns: userOpenConns,
		connsMutex:    connsMutex,
	}
}

func (mh *MessageHandler) Subscribe(c echo.Context) error {
	var req request.MessageSubscribe

	if err := request.BindT(&req, c); err != nil {
		return echo.ErrBadRequest
	}
	if err := request.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var claims map[string]interface{}
	var err error
	claims, err = mh.jwtConfig.ValidateToken(*req.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	userID, _ := claims["id"].(uint64)
	username, _ := claims["username"].(string)
	if userID != *req.UserID {
		return echo.NewHTTPError(http.StatusForbidden, model.ErrAccessForbidden.Error())
	}

	mh.connsMutex.Lock()
	_, ok := (*mh.userOpenConns)[userID]
	if !ok {
		(*mh.userOpenConns)[userID] = make(map[uint64]*websocket.Conn)
	}

	chatConnMap := (*mh.userOpenConns)[userID]
	_, ok = chatConnMap[*req.ChatID]
	if ok { // if a connection for that user in that chat exists, close it and later make another one
		(*mh.userOpenConns)[userID][*req.ChatID].Close()
	}

	// create a new websocket
	ws, err := mh.upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return echo.ErrBadRequest
	}

	(*mh.userOpenConns)[userID][*req.ChatID] = ws
	mh.connsMutex.Unlock()

	ws.SetReadLimit(int64(model.MaxMsgSize))

readLoop:
	for {
		msgType, msgStr, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
				break
			}

			fmt.Printf("error while receiving msg from user: [%v] , chatID: [%v]\n", username, *req.ChatID)
		}

		switch msgType {
		case websocket.CloseMessage:
			break readLoop
		case websocket.TextMessage:
			lastRegisteredMessageID++
			msg := model.Message{}
			msg.ID = uint64(lastRegisteredMessageID)
			msg.ChatID = *req.ChatID
			msg.SenderID = *req.UserID
			msg.Content = string(msgStr)
			*mh.broadcast <- msg

		default:
			fmt.Printf("message type not supported, from username: [%v], chatID: [%v]\n", username, *req.ChatID)
			continue
		}

	}

	ws.Close()
	fmt.Printf("closed socket for user: [%v], chatID: [%v]\n", username, *req.ChatID)
	return nil
}

func (mh *MessageHandler) HandleMessages(broadcast *chan model.Message) {
	for {
		msg := <-*broadcast
		for userID, connsMap := range *mh.userOpenConns {
			if msg.SenderID != userID { // we don't send the message to the sender himself
				for chatID, ws := range connsMap {
					if msg.ChatID == chatID {
						if err := ws.WriteMessage(websocket.TextMessage, []byte(msg.Content)); err != nil {
							fmt.Printf("error while sending msg to socket, locAddr: [%v], rmtAddr: [%v]", ws.UnderlyingConn().LocalAddr(), ws.UnderlyingConn().RemoteAddr())
						}
					}
				}
			}
		}
	}
}

func (mh *MessageHandler) RegisterMethods(g *echo.Group) {
	g.POST("chats/:chat_id/subscribe/:id", mh.Subscribe)
}
