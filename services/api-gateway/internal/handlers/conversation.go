package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	conversationv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/conversation/v1"
	"github.com/Glistand/HelpDesk/services/api-gateway/internal/middleware"
	"github.com/google/uuid"
)

type widgetSessionBody struct {
	VisitorID string `json:"visitor_id"`
	SiteKey   string `json:"site_key"`
}

func (a *API) WidgetSession(w http.ResponseWriter, r *http.Request) {
	var body widgetSessionBody
	_ = json.NewDecoder(r.Body).Decode(&body)

	siteKey := strings.TrimSpace(body.SiteKey)
	if siteKey == "" {
		siteKey = strings.TrimSpace(r.Header.Get("X-Site-Key"))
	}
	if siteKey == "" {
		siteKey = a.siteKey
	}
	if siteKey != a.siteKey {
		writeErr(w, http.StatusForbidden, "invalid site_key")
		return
	}

	visitorID := strings.TrimSpace(body.VisitorID)
	if visitorID == "" {
		visitorID = strings.TrimSpace(r.Header.Get("X-Visitor-Id"))
	}
	if visitorID == "" {
		visitorID = uuid.NewString()
	}

	resp, err := a.c.Conversation.CreateConversation(r.Context(), &conversationv1.CreateConversationRequest{
		SiteKey:   siteKey,
		VisitorId: visitorID,
	})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"visitor_id":   visitorID,
		"conversation": conversationJSON(resp.GetConversation()),
		"messages":     messagesJSON(resp.GetMessages()),
	})
}

type postMessageBody struct {
	Body string `json:"body"`
}

func (a *API) WidgetPostMessage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	visitorID := strings.TrimSpace(r.Header.Get("X-Visitor-Id"))
	if visitorID == "" {
		writeErr(w, http.StatusBadRequest, "X-Visitor-Id required")
		return
	}
	if err := a.assertVisitorOwns(r, id, visitorID); err != nil {
		writeGRPCErr(w, err)
		return
	}

	var body postMessageBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Body) == "" {
		writeErr(w, http.StatusBadRequest, "body required")
		return
	}

	resp, err := a.c.Conversation.PostMessage(r.Context(), &conversationv1.PostMessageRequest{
		ConversationId: id,
		Role:           conversationv1.MessageRole_MESSAGE_ROLE_VISITOR,
		Body:           body.Body,
		AuthorId:       visitorID,
		InvokeBot:      true,
	})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	out := map[string]any{
		"message":      messageJSON(resp.GetMessage()),
		"conversation": conversationJSON(resp.GetConversation()),
	}
	if resp.GetBotReply() != nil {
		out["bot_reply"] = messageJSON(resp.GetBotReply())
	}
	writeJSON(w, http.StatusCreated, out)
}

func (a *API) WidgetListMessages(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	visitorID := strings.TrimSpace(r.Header.Get("X-Visitor-Id"))
	if visitorID == "" {
		writeErr(w, http.StatusBadRequest, "X-Visitor-Id required")
		return
	}
	if err := a.assertVisitorOwns(r, id, visitorID); err != nil {
		writeGRPCErr(w, err)
		return
	}
	resp, err := a.c.Conversation.ListMessages(r.Context(), &conversationv1.ListMessagesRequest{
		ConversationId: id,
		AfterId:        r.URL.Query().Get("after"),
	})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"messages": messagesJSON(resp.GetMessages())})
}

func (a *API) WidgetHandoff(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	visitorID := strings.TrimSpace(r.Header.Get("X-Visitor-Id"))
	if visitorID == "" {
		writeErr(w, http.StatusBadRequest, "X-Visitor-Id required")
		return
	}
	if err := a.assertVisitorOwns(r, id, visitorID); err != nil {
		writeGRPCErr(w, err)
		return
	}
	resp, err := a.c.Conversation.RequestHuman(r.Context(), &conversationv1.RequestHumanRequest{
		ConversationId: id,
	})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"conversation":   conversationJSON(resp.GetConversation()),
		"system_message": messageJSON(resp.GetSystemMessage()),
	})
}

func (a *API) ListConversations(w http.ResponseWriter, r *http.Request) {
	status := parseConversationStatus(r.URL.Query().Get("status"))
	resp, err := a.c.Conversation.ListConversations(r.Context(), &conversationv1.ListConversationsRequest{
		Status:   status,
		PageSize: 50,
	})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	items := make([]map[string]any, 0, len(resp.GetConversations()))
	for _, c := range resp.GetConversations() {
		items = append(items, conversationJSON(c))
	}
	writeJSON(w, http.StatusOK, map[string]any{"conversations": items})
}

func (a *API) GetConversation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	conv, err := a.c.Conversation.GetConversation(r.Context(), &conversationv1.GetConversationRequest{Id: id})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	msgs, err := a.c.Conversation.ListMessages(r.Context(), &conversationv1.ListMessagesRequest{ConversationId: id})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"conversation": conversationJSON(conv.GetConversation()),
		"messages":     messagesJSON(msgs.GetMessages()),
	})
}

func (a *API) AgentPostMessage(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id := r.PathValue("id")
	var body postMessageBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Body) == "" {
		writeErr(w, http.StatusBadRequest, "body required")
		return
	}
	resp, err := a.c.Conversation.PostMessage(r.Context(), &conversationv1.PostMessageRequest{
		ConversationId: id,
		Role:           conversationv1.MessageRole_MESSAGE_ROLE_AGENT,
		Body:           body.Body,
		AuthorId:       user.GetId(),
		InvokeBot:      false,
	})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"message":      messageJSON(resp.GetMessage()),
		"conversation": conversationJSON(resp.GetConversation()),
	})
}

func (a *API) ResolveConversation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	resp, err := a.c.Conversation.ResolveConversation(r.Context(), &conversationv1.ResolveConversationRequest{Id: id})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"conversation": conversationJSON(resp.GetConversation())})
}

func (a *API) assertVisitorOwns(r *http.Request, conversationID, visitorID string) error {
	resp, err := a.c.Conversation.GetConversation(r.Context(), &conversationv1.GetConversationRequest{Id: conversationID})
	if err != nil {
		return err
	}
	c := resp.GetConversation()
	if c.GetVisitorId() != visitorID {
		return errForbidden("visitor mismatch")
	}
	if c.GetSiteKey() != a.siteKey {
		return errForbidden("site_key mismatch")
	}
	return nil
}

type simpleErr string

func (e simpleErr) Error() string { return string(e) }

func errForbidden(msg string) error { return simpleErr(msg) }

func conversationJSON(c *conversationv1.Conversation) map[string]any {
	if c == nil {
		return nil
	}
	return map[string]any{
		"id":          c.GetId(),
		"site_key":    c.GetSiteKey(),
		"visitor_id":  c.GetVisitorId(),
		"status":      conversationStatusString(c.GetStatus()),
		"assignee_id": c.GetAssigneeId(),
		"preview":     c.GetPreview(),
		"created_at":  c.GetCreatedAt(),
		"updated_at":  c.GetUpdatedAt(),
	}
}

func messageJSON(m *conversationv1.Message) map[string]any {
	if m == nil {
		return nil
	}
	return map[string]any{
		"id":              m.GetId(),
		"conversation_id": m.GetConversationId(),
		"role":            messageRoleString(m.GetRole()),
		"body":            m.GetBody(),
		"author_id":       m.GetAuthorId(),
		"created_at":      m.GetCreatedAt(),
	}
}

func messagesJSON(msgs []*conversationv1.Message) []map[string]any {
	out := make([]map[string]any, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, messageJSON(m))
	}
	return out
}

func conversationStatusString(s conversationv1.ConversationStatus) string {
	switch s {
	case conversationv1.ConversationStatus_CONVERSATION_STATUS_BOT:
		return "bot"
	case conversationv1.ConversationStatus_CONVERSATION_STATUS_WAITING_AGENT:
		return "waiting_agent"
	case conversationv1.ConversationStatus_CONVERSATION_STATUS_OPEN:
		return "open"
	case conversationv1.ConversationStatus_CONVERSATION_STATUS_RESOLVED:
		return "resolved"
	default:
		return "unspecified"
	}
}

func messageRoleString(r conversationv1.MessageRole) string {
	switch r {
	case conversationv1.MessageRole_MESSAGE_ROLE_VISITOR:
		return "visitor"
	case conversationv1.MessageRole_MESSAGE_ROLE_BOT:
		return "bot"
	case conversationv1.MessageRole_MESSAGE_ROLE_AGENT:
		return "agent"
	case conversationv1.MessageRole_MESSAGE_ROLE_SYSTEM:
		return "system"
	default:
		return "unspecified"
	}
}

func parseConversationStatus(s string) conversationv1.ConversationStatus {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "bot":
		return conversationv1.ConversationStatus_CONVERSATION_STATUS_BOT
	case "waiting_agent":
		return conversationv1.ConversationStatus_CONVERSATION_STATUS_WAITING_AGENT
	case "open":
		return conversationv1.ConversationStatus_CONVERSATION_STATUS_OPEN
	case "resolved":
		return conversationv1.ConversationStatus_CONVERSATION_STATUS_RESOLVED
	default:
		return conversationv1.ConversationStatus_CONVERSATION_STATUS_UNSPECIFIED
	}
}
