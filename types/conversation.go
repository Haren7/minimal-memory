package types

type RegisterConversationInput struct {
	Agent string
	User  string
}

type RegisterConversationOutput struct {
	ConversationID string
}
