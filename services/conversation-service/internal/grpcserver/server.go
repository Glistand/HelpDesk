package grpcserver

import (
	"context"
	"log/slog"
	"strings"
	"time"

	conversationv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/conversation/v1"
	"github.com/Glistand/HelpDesk/libs/eventkit/envelope"
	"github.com/Glistand/HelpDesk/libs/eventkit/natsx"
	"github.com/Glistand/HelpDesk/libs/eventkit/subjects"
	"github.com/Glistand/HelpDesk/libs/grpckit/statuserr"
	"github.com/Glistand/HelpDesk/services/conversation-service/internal/openrouter"
	"github.com/Glistand/HelpDesk/services/conversation-service/internal/repository"
)

type Server struct {
	conversationv1.UnimplementedConversationServiceServer
	repo                *repository.Repo
	bot                 *openrouter.Client
	projectInstructions string
	pub                 *natsx.Publisher
	logger              *slog.Logger
}

func New(repo *repository.Repo, bot *openrouter.Client, projectInstructions string, pub *natsx.Publisher, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	return &Server{repo: repo, bot: bot, projectInstructions: projectInstructions, pub: pub, logger: logger}
}

func (s *Server) CreateConversation(ctx context.Context, req *conversationv1.CreateConversationRequest) (*conversationv1.CreateConversationResponse, error) {
	site := strings.TrimSpace(req.GetSiteKey())
	visitor := strings.TrimSpace(req.GetVisitorId())
	if site == "" || visitor == "" {
		return nil, statuserr.InvalidArgument("site_key and visitor_id required")
	}
	c, welcome, err := s.repo.CreateConversation(ctx, site, visitor)
	if err != nil {
		return nil, statuserr.Internal(err.Error())
	}
	s.publish(ctx, subjects.ConversationCreated, c.ID, map[string]any{
		"conversation_id": c.ID, "site_key": c.SiteKey, "visitor_id": c.VisitorID,
	})
	return &conversationv1.CreateConversationResponse{
		Conversation: toProtoConv(c),
		Messages:     []*conversationv1.Message{toProtoMsg(welcome)},
	}, nil
}

func (s *Server) GetConversation(ctx context.Context, req *conversationv1.GetConversationRequest) (*conversationv1.GetConversationResponse, error) {
	c, err := s.repo.GetConversation(ctx, req.GetId())
	if err != nil {
		return nil, statuserr.NotFound("conversation not found")
	}
	return &conversationv1.GetConversationResponse{Conversation: toProtoConv(c)}, nil
}

func (s *Server) ListConversations(ctx context.Context, req *conversationv1.ListConversationsRequest) (*conversationv1.ListConversationsResponse, error) {
	status := statusString(req.GetStatus())
	list, err := s.repo.ListConversations(ctx, status, int(req.GetPageSize()))
	if err != nil {
		return nil, statuserr.Internal(err.Error())
	}
	out := make([]*conversationv1.Conversation, 0, len(list))
	for _, c := range list {
		out = append(out, toProtoConv(c))
	}
	return &conversationv1.ListConversationsResponse{Conversations: out}, nil
}

func (s *Server) PostMessage(ctx context.Context, req *conversationv1.PostMessageRequest) (*conversationv1.PostMessageResponse, error) {
	body := strings.TrimSpace(req.GetBody())
	if req.GetConversationId() == "" || body == "" {
		return nil, statuserr.InvalidArgument("conversation_id and body required")
	}
	c, err := s.repo.GetConversation(ctx, req.GetConversationId())
	if err != nil {
		return nil, statuserr.NotFound("conversation not found")
	}
	role := roleString(req.GetRole())
	if role == "" || role == "unspecified" {
		return nil, statuserr.InvalidArgument("role required")
	}
	author := strings.TrimSpace(req.GetAuthorId())
	msg, err := s.repo.AddMessage(ctx, c.ID, role, body, author)
	if err != nil {
		return nil, statuserr.Internal(err.Error())
	}
	s.publish(ctx, subjects.ConversationMessage, c.ID, map[string]any{
		"conversation_id": c.ID, "message_id": msg.ID, "role": role,
	})

	resp := &conversationv1.PostMessageResponse{Message: toProtoMsg(msg)}

	if req.GetInvokeBot() && role == "visitor" && (c.Status == "bot" || c.Status == "open") {
		botBody := s.generateBotReply(ctx, c.ID)
		botMsg, err := s.repo.AddMessage(ctx, c.ID, "bot", botBody, "bot")
		if err != nil {
			return nil, statuserr.Internal(err.Error())
		}
		s.publish(ctx, subjects.ConversationMessage, c.ID, map[string]any{
			"conversation_id": c.ID, "message_id": botMsg.ID, "role": "bot",
		})
		resp.BotReply = toProtoMsg(botMsg)
	}

	if role == "agent" && c.Status == "waiting_agent" {
		c, _ = s.repo.SetStatus(ctx, c.ID, "open", author)
	} else {
		c, _ = s.repo.GetConversation(ctx, c.ID)
	}
	resp.Conversation = toProtoConv(c)
	return resp, nil
}

func (s *Server) ListMessages(ctx context.Context, req *conversationv1.ListMessagesRequest) (*conversationv1.ListMessagesResponse, error) {
	list, err := s.repo.ListMessages(ctx, req.GetConversationId(), req.GetAfterId())
	if err != nil {
		return nil, statuserr.Internal(err.Error())
	}
	out := make([]*conversationv1.Message, 0, len(list))
	for _, m := range list {
		out = append(out, toProtoMsg(m))
	}
	return &conversationv1.ListMessagesResponse{Messages: out}, nil
}

func (s *Server) RequestHuman(ctx context.Context, req *conversationv1.RequestHumanRequest) (*conversationv1.RequestHumanResponse, error) {
	c, err := s.repo.GetConversation(ctx, req.GetConversationId())
	if err != nil {
		return nil, statuserr.NotFound("conversation not found")
	}
	c, err = s.repo.SetStatus(ctx, c.ID, "waiting_agent", "")
	if err != nil {
		return nil, statuserr.Internal(err.Error())
	}
	sys, err := s.repo.AddMessage(ctx, c.ID, "system", "Запрос передан агенту. Обычно отвечаем в течение нескольких минут.", "system")
	if err != nil {
		return nil, statuserr.Internal(err.Error())
	}
	s.publish(ctx, subjects.ConversationHandoff, c.ID, map[string]any{"conversation_id": c.ID})
	return &conversationv1.RequestHumanResponse{
		Conversation:  toProtoConv(c),
		SystemMessage: toProtoMsg(sys),
	}, nil
}

func (s *Server) ResolveConversation(ctx context.Context, req *conversationv1.ResolveConversationRequest) (*conversationv1.ResolveConversationResponse, error) {
	c, err := s.repo.SetStatus(ctx, req.GetId(), "resolved", "")
	if err != nil {
		return nil, statuserr.NotFound("conversation not found")
	}
	return &conversationv1.ResolveConversationResponse{Conversation: toProtoConv(c)}, nil
}

func (s *Server) generateBotReply(ctx context.Context, conversationID string) string {
	fallback := "Сейчас не удалось получить ответ бота. Нажмите «Нужен человек», чтобы связаться с агентом."
	if s.bot == nil || s.bot.APIKey == "" {
		return fallback
	}
	hist, err := s.repo.RecentForBot(ctx, conversationID, 16)
	if err != nil {
		s.logger.Warn("bot history failed", "error", err)
		return fallback
	}
	msgs := make([]openrouter.ChatMessage, 0, len(hist))
	for _, m := range hist {
		if m.Role == "system" {
			continue
		}
		msgs = append(msgs, openrouter.HistoryFromRepo(m.Role, m.Body))
	}
	botCtx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()
	reply, err := s.bot.ReplyWithInstructions(botCtx, msgs, s.projectInstructions)
	if err != nil {
		s.logger.Warn("openrouter failed", "error", err)
		return fallback
	}
	return reply
}

func (s *Server) publish(ctx context.Context, subject, aggregateID string, payload map[string]any) {
	if s.pub == nil {
		return
	}
	ev, err := envelope.New(strings.TrimPrefix(subject, "helpdesk."), aggregateID, "", payload)
	if err != nil {
		return
	}
	_ = s.pub.Publish(ctx, subject, ev)
}

func toProtoConv(c repository.Conversation) *conversationv1.Conversation {
	return &conversationv1.Conversation{
		Id:         c.ID,
		SiteKey:    c.SiteKey,
		VisitorId:  c.VisitorID,
		Status:     toProtoStatus(c.Status),
		AssigneeId: c.AssigneeID,
		Preview:    c.Preview,
		CreatedAt:  c.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:  c.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func toProtoMsg(m repository.Message) *conversationv1.Message {
	return &conversationv1.Message{
		Id:             m.ID,
		ConversationId: m.ConversationID,
		Role:           toProtoRole(m.Role),
		Body:           m.Body,
		AuthorId:       m.AuthorID,
		CreatedAt:      m.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func toProtoStatus(s string) conversationv1.ConversationStatus {
	switch s {
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

func statusString(s conversationv1.ConversationStatus) string {
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
		return ""
	}
}

func toProtoRole(r string) conversationv1.MessageRole {
	switch r {
	case "visitor":
		return conversationv1.MessageRole_MESSAGE_ROLE_VISITOR
	case "bot":
		return conversationv1.MessageRole_MESSAGE_ROLE_BOT
	case "agent":
		return conversationv1.MessageRole_MESSAGE_ROLE_AGENT
	case "system":
		return conversationv1.MessageRole_MESSAGE_ROLE_SYSTEM
	default:
		return conversationv1.MessageRole_MESSAGE_ROLE_UNSPECIFIED
	}
}

func roleString(r conversationv1.MessageRole) string {
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
		return ""
	}
}
