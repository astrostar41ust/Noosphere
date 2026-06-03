package chat

type SendMessageRequest struct {
	SessionID string `json:"session_id" validate:"required,uuid"`
	Role      string `json:"role"       validate:"required,oneof=user assistant system"`
	Content   string `json:"content"    validate:"required,min=1"`
}

type ChatMessageResponse struct {
	ID        string `json:"id"`
	SessionID string `json:"session_id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func MapMessageToResponse(msg *Message) ChatMessageResponse {
	return ChatMessageResponse{
		ID:        msg.ID.String(),
		SessionID: msg.SessionID.String(),
		Role:      msg.Role,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}