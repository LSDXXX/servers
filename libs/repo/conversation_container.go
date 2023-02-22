package repo

type ConversationContainer interface {
	GetParentId(convId string) (string, bool)

	SetParentId(convId string, parentId string)
}
