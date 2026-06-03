package chat

import (
	"net/http"
	"noosphere/backend-api/internal/pkg/httpio" 

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ChatController struct {
	service  ChatService
	validate *validator.Validate
}

func NewChatController(service ChatService) *ChatController {
	return &ChatController{
		service:  service,
		validate: validator.New(),
	}
}

func (c *ChatController) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	var req SendMessageRequest

	if err := httpio.DecodeAndValidate(w, r, &req, c.validate); err != nil {
		return
	}

	sessionUUID := uuid.MustParse(req.SessionID)

	msg, err := c.service.SendMessage(r.Context(), sessionUUID, req.Role, req.Content)
	if err != nil {
		httpio.RespondWithError(w, http.StatusInternalServerError, "Failed to process chat mutation: "+err.Error())
		return
	}

	httpio.RespondWithJSON(w, http.StatusCreated, MapMessageToResponse(msg))
}