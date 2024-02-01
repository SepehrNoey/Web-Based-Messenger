package handler

import (
	"fmt"
	"net/http"

	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/model"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/repository/accountrepo"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/repository/chatrepo"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/domain/repository/messagerepo"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/infra/http/auth"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/infra/http/clientdto"
	"github.com/SepehrNoey/Web-Based-Messenger/Back-end/internal/infra/http/request"
	"github.com/labstack/echo/v4"
)

var lastRegisteredChatID = 0

type ChatHandler struct {
	accRepo   accountrepo.Repository
	chatRepo  chatrepo.Repository
	msgRepo   messagerepo.Repository
	jwtConfig auth.JWTConfig
}

func NewChatHandler(accRepo accountrepo.Repository, chatRepo chatrepo.Repository, msgRepo messagerepo.Repository, jwtConfig auth.JWTConfig) *ChatHandler {
	return &ChatHandler{
		accRepo:   accRepo,
		chatRepo:  chatRepo,
		msgRepo:   msgRepo,
		jwtConfig: jwtConfig,
	}
}

func (ch *ChatHandler) BelongsTo(memberID *uint64, chat *model.Chat) bool {
	for _, id := range chat.Members {
		if *memberID == id {
			return true
		}
	}

	return false
}

func (ch *ChatHandler) Create(c echo.Context) error {
	var req request.ChatCreate

	if err := request.BindT(&req, c); err != nil {
		return echo.ErrBadRequest
	}
	if err := request.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var claims map[string]interface{}
	var err error
	claims, err = ch.jwtConfig.ValidateToken(*req.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	userID, _ := claims["id"].(uint64)
	contactAccs := ch.accRepo.Get(c.Request().Context(), accountrepo.GetCommand{
		ID: req.ContactID,
	})
	if len(contactAccs) > 1 {
		return echo.ErrInternalServerError
	}
	if len(contactAccs) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, model.ErrUserNotFound.Error())
	}

	mems := []uint64{userID, *req.ContactID}
	chatsWithContact := ch.chatRepo.Get(c.Request().Context(), chatrepo.GetCommand{
		Members: &mems,
	})
	if len(chatsWithContact) > 1 {
		return echo.ErrInternalServerError
	}
	if len(chatsWithContact) == 1 {
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrChatDuplicate.Error())
	}

	if err := ch.chatRepo.Create(c.Request().Context(), model.Chat{
		ID:      uint64(lastRegisteredChatID + 1),
		Members: []uint64{userID, *req.ContactID},
	}); err != nil {
		return echo.ErrInternalServerError
	}

	lastRegisteredChatID++
	return c.JSON(http.StatusCreated, fmt.Sprintf("chat_id: %v", lastRegisteredChatID))
}

func (ch *ChatHandler) Get(c echo.Context) error {
	var req request.TokenOnly

	if err := request.BindT(&req, c); err != nil {
		return err
	}
	if err := request.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var claims map[string]interface{}
	var err error
	claims, err = ch.jwtConfig.ValidateToken(*req.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	userID, _ := claims["id"].(uint64)
	mems := []uint64{userID}
	chats := ch.chatRepo.Get(c.Request().Context(), chatrepo.GetCommand{
		Members: &mems,
	})

	return c.JSON(http.StatusOK, chats)
}

func (ch *ChatHandler) GetByID(c echo.Context) error {
	var req request.TokenAndChatID

	if err := request.BindT(&req, c); err != nil {
		return err
	}
	if err := request.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var claims map[string]interface{}
	var err error
	claims, err = ch.jwtConfig.ValidateToken(*req.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	userID, _ := claims["id"].(uint64)
	chats := ch.chatRepo.Get(c.Request().Context(), chatrepo.GetCommand{
		ID: req.ID,
	})

	if len(chats) > 1 {
		return echo.ErrInternalServerError
	} else if len(chats) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, model.ErrChatNotFound.Error())
	}

	chat := chats[0]
	if !ch.BelongsTo(&userID, &chat) {
		return echo.NewHTTPError(http.StatusForbidden, model.ErrAccessForbidden.Error())
	}

	msgs := ch.msgRepo.Get(c.Request().Context(), messagerepo.GetCommand{
		ChatID: req.ID,
	})
	dto := clientdto.ChatWithContentDTO{
		ID:        req.ID,
		Members:   &chat.Members,
		CreatedAt: &chat.CreatedAt,
		UpdatedAt: &chat.UpdatedAt,
		Messages:  &msgs,
	}

	return c.JSON(http.StatusOK, dto)

}

func (ch *ChatHandler) Delete(c echo.Context) error {
	var req request.TokenAndChatID

	if err := request.BindT(&req, c); err != nil {
		return err
	}
	if err := request.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var claims map[string]interface{}
	var err error
	claims, err = ch.jwtConfig.ValidateToken(*req.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	userID, _ := claims["id"].(uint64)
	chats := ch.chatRepo.Get(c.Request().Context(), chatrepo.GetCommand{
		ID: req.ID,
	})

	if len(chats) > 1 {
		return echo.ErrInternalServerError
	} else if len(chats) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, model.ErrChatNotFound.Error())
	}

	chat := chats[0]
	if !ch.BelongsTo(&userID, &chat) {
		return echo.NewHTTPError(http.StatusForbidden, model.ErrAccessForbidden.Error())
	}

	if err := ch.chatRepo.Delete(c.Request().Context(), chatrepo.GetCommand{
		ID: &chat.ID,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (ch *ChatHandler) DeleteMsg(c echo.Context) error {
	var req request.ChatDeleteMsg

	if err := request.BindT(&req, c); err != nil {
		return err
	}
	if err := request.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var claims map[string]interface{}
	var err error
	claims, err = ch.jwtConfig.ValidateToken(*req.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	userID, _ := claims["id"].(uint64)
	chats := ch.chatRepo.Get(c.Request().Context(), chatrepo.GetCommand{
		ID: req.ChatID,
	})

	if len(chats) > 1 {
		return echo.ErrInternalServerError
	} else if len(chats) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, model.ErrChatNotFound.Error())
	}

	chat := chats[0]
	if !ch.BelongsTo(&userID, &chat) {
		return echo.NewHTTPError(http.StatusForbidden, model.ErrAccessForbidden.Error())
	}

	msgs := ch.msgRepo.Get(c.Request().Context(), messagerepo.GetCommand{
		ID:     req.MsgID,
		ChatID: req.ChatID,
	})
	if len(msgs) > 1 {
		return echo.ErrInternalServerError
	} else if len(msgs) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, model.ErrMessageNotFound.Error())
	}

	if err := ch.msgRepo.Delete(c.Request().Context(), messagerepo.GetCommand{
		ID:     req.MsgID,
		ChatID: req.ChatID,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	// handle updating the deletion of msg in client app of the other clients !
	return c.NoContent(http.StatusOK)
}

func (ch *ChatHandler) RegisterMethods(g *echo.Group) {
	g.POST("chats", ch.Create)
	g.GET("chats", ch.Get)
	g.GET("chats/:chat_id", ch.GetByID)
	g.DELETE("chats/:chat_id", ch.Delete)
	g.DELETE("chats/:chat_id/messages/:message_id", ch.DeleteMsg)
}
