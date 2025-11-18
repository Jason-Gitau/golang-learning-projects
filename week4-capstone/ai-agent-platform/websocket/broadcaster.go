package websocket

// HubBroadcaster adapts the Hub to work with the agent runner
type HubBroadcaster struct {
	hub *Hub
}

// NewHubBroadcaster creates a new hub broadcaster
func NewHubBroadcaster(hub *Hub) *HubBroadcaster {
	return &HubBroadcaster{hub: hub}
}

// BroadcastChunk broadcasts a response chunk
func (b *HubBroadcaster) BroadcastChunk(conversationID, messageID, content string) error {
	msg := NewMessage(MessageTypeAgentResponseChunk, content)
	msg.ConversationID = conversationID
	msg.MessageID = messageID

	b.hub.BroadcastToConversation(conversationID, msg)
	return nil
}

// BroadcastToolCall broadcasts a tool call
func (b *HubBroadcaster) BroadcastToolCall(conversationID, tool string, params map[string]interface{}) error {
	msg := NewToolCallMessage(tool, params)
	msg.ConversationID = conversationID

	b.hub.BroadcastToConversation(conversationID, msg)
	return nil
}

// BroadcastToolResult broadcasts a tool result
func (b *HubBroadcaster) BroadcastToolResult(conversationID, tool string, result interface{}) error {
	msg := NewToolResultMessage(tool, result)
	msg.ConversationID = conversationID

	b.hub.BroadcastToConversation(conversationID, msg)
	return nil
}

// BroadcastComplete broadcasts a completion message
func (b *HubBroadcaster) BroadcastComplete(conversationID, messageID string) error {
	msg := NewMessage(MessageTypeAgentResponseComplete, "")
	msg.ConversationID = conversationID
	msg.MessageID = messageID

	b.hub.BroadcastToConversation(conversationID, msg)
	return nil
}

// BroadcastError broadcasts an error message
func (b *HubBroadcaster) BroadcastError(conversationID, errorMsg string) error {
	msg := NewErrorMessage(errorMsg)
	msg.ConversationID = conversationID

	b.hub.BroadcastToConversation(conversationID, msg)
	return nil
}
