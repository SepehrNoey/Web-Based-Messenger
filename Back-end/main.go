package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/model"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/infra/http/auth"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/infra/http/handler"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/infra/repository/account/accountsql"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/infra/repository/chat/chatsql"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/infra/repository/contact/contactsql"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/infra/repository/message/messagesql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=0150188511 dbname=Web-Based-Messenger port=5432 sslmode=disable TimeZone=Asia/Tehran"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("error connecting to database: %v", err)
		return
	}

	if err = db.AutoMigrate(&accountsql.AccountDTO{}, &contactsql.ContactDTO{},
		&chatsql.ChatDTO{}, &messagesql.MessageDTO{}); err != nil {
		fmt.Printf("failed to automigrate: %v", err)
	}

	app := echo.New()
	accountRepo := accountsql.New(db)
	contactRepo := contactsql.New(db)
	chatRepo := chatsql.New(db)
	messageRepo := messagesql.New(db)
	upgrader := websocket.Upgrader{
		ReadBufferSize:  int(model.MaxMsgSize),
		WriteBufferSize: int(model.MaxMsgSize),
	}
	broadcast := make(chan model.Message)
	openConns := make(map[uint64]map[uint64]*websocket.Conn)
	connsMutex := sync.Mutex{}

	secretKey := []byte("secret-key-of-web-based-messenger-for-jwt-authentication")
	expDur := time.Minute * 5
	signMethod := jwt.SigningMethodHS256
	iss := "messenger-server"
	aud := make([]string, 0)
	aud = append(aud, "messenger-client")
	jwtConfig := auth.NewJWTConfig(secretKey, expDur, signMethod, iss, aud)

	accHnd := handler.NewAccountHandler(accountRepo, *jwtConfig)
	conHnd := handler.NewContactHandler(contactRepo, accountRepo, *jwtConfig)
	chHnd := handler.NewChatHandler(accountRepo, chatRepo, messageRepo, *jwtConfig)
	msgHnd := handler.NewMessageHandler(messageRepo, chatRepo, jwtConfig, &upgrader, &broadcast, &openConns, &connsMutex)

	accHnd.RegisterMethods(app.Group("api/"))
	conHnd.RegisterMethods(app.Group("api/users/"))
	chHnd.RegisterMethods(app.Group("api/chats/"))
	msgHnd.RegisterMethods(app.Group("api/chats/"))

	go msgHnd.HandleMessages(&broadcast)
	if err := app.Start("0.0.0.0:2024"); err != nil {
		log.Fatalf("server failed to start %v", err)
	}
}
