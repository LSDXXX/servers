package conversation

type AskReq struct {
	ConversationId string `json:"conversation_id"`
	Content        string `json:"content" binding:"required"`
}
