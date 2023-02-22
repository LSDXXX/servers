package bot

type askMessage struct {
	Id      string `json:"id"`
	Role    string `json:"role"`
	Content struct {
		ContentType string   `json:"content_type"`
		Parts       []string `json:"parts"`
	} `json:"content"`
}

type askReq struct {
	Action          string       `json:"action"`
	Messages        []askMessage `json:"messages"`
	ConversationId  string       `json:"conversation_id,omitempty"`
	ParentMessageId string       `json:"parent_message_id"`
	Model           string       `json:"model"`
}

type ResponseMessage struct {
	Message        Message     `json:"message"`
	ConversationID string      `json:"conversation_id"`
	Error          interface{} `json:"error"`
}

type Content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}
type FinishDetails struct {
	Type string `json:"type"`
	Stop string `json:"stop"`
}

type Metadata struct {
	MessageType   string        `json:"message_type"`
	ModelSlug     string        `json:"model_slug"`
	FinishDetails FinishDetails `json:"finish_details"`
}
type Message struct {
	ID         string      `json:"id"`
	Role       string      `json:"role"`
	User       interface{} `json:"user"`
	CreateTime interface{} `json:"create_time"`
	UpdateTime interface{} `json:"update_time"`
	Content    Content     `json:"content"`
	EndTurn    interface{} `json:"end_turn"`
	Weight     float64     `json:"weight"`
	Metadata   Metadata    `json:"metadata"`
	Recipient  string      `json:"recipient"`
}

type ConversationReq struct {
	ConversationId string `json:"conversation_id"`
	Content        string `json:"content" binding:"required"`
}

type HttpResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type User struct {
	Id        int
	UserName  string
	FirstName string
	LastName  string
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
