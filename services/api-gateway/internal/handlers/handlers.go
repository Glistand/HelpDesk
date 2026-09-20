package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"

	assignmentv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/assignment/v1"
	auditv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/audit/v1"
	authv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/auth/v1"
	searchv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/search/v1"
	slav1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/sla/v1"
	ticketv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/ticket/v1"
	"github.com/Glistand/HelpDesk/services/api-gateway/internal/clients"
	"github.com/Glistand/HelpDesk/services/api-gateway/internal/middleware"
)

type API struct {
	c       *clients.Clients
	siteKey string
}

func New(c *clients.Clients, siteKey string) *API {
	if siteKey == "" {
		siteKey = "demo-site"
	}
	return &API{c: c, siteKey: siteKey}
}

func (a *API) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	var body loginBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}
	resp, err := a.c.Auth.Login(r.Context(), &authv1.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		writeErr(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"access_token": resp.GetAccessToken(),
		"expires_at":   resp.GetExpiresAtUnix(),
		"user":         userJSON(resp.GetUser()),
	})
}

type createTicketBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Category    string `json:"category"`
	Requester   string `json:"requester"`
}

func (a *API) CreateTicket(w http.ResponseWriter, r *http.Request) {
	var body createTicketBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.Requester == "" {
		if u := middleware.UserFromContext(r.Context()); u != nil {
			body.Requester = u.GetName()
		}
	}
	resp, err := a.c.Ticket.CreateTicket(r.Context(), &ticketv1.CreateTicketRequest{
		Title:       body.Title,
		Description: body.Description,
		Priority:    parsePriority(body.Priority),
		Category:    body.Category,
		Requester:   body.Requester,
	})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, ticketJSON(resp.GetTicket()))
}

func (a *API) ListTickets(w http.ResponseWriter, r *http.Request) {
	statusQ := r.URL.Query().Get("status")
	assignee := r.URL.Query().Get("assignee_id")
	resp, err := a.c.Ticket.ListTickets(r.Context(), &ticketv1.ListTicketsRequest{
		Status:     parseStatus(statusQ),
		AssigneeId: assignee,
		PageSize:   50,
	})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	items := make([]map[string]any, 0, len(resp.GetTickets()))
	for _, t := range resp.GetTickets() {
		items = append(items, ticketJSON(t))
	}
	writeJSON(w, http.StatusOK, map[string]any{"tickets": items})
}

func (a *API) GetTicket(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErr(w, http.StatusBadRequest, "missing id")
		return
	}
	resp, err := a.c.Ticket.GetTicket(r.Context(), &ticketv1.GetTicketRequest{Id: id})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ticketJSON(resp.GetTicket()))
}

func (a *API) GetTimeline(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErr(w, http.StatusBadRequest, "missing id")
		return
	}
	resp, err := a.c.Audit.GetTimeline(r.Context(), &auditv1.GetTimelineRequest{TicketId: id})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	events := make([]map[string]any, 0, len(resp.GetEvents()))
	for _, e := range resp.GetEvents() {
		events = append(events, map[string]any{
			"id":          e.GetId(),
			"ticket_id":   e.GetTicketId(),
			"event_id":    e.GetEventId(),
			"event_type":  e.GetEventType(),
			"title":       e.GetTitle(),
			"detail":      e.GetDetail(),
			"actor":       e.GetActor(),
			"occurred_at": e.GetOccurredAt(),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"events": events})
}

func (a *API) GetSLA(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErr(w, http.StatusBadRequest, "missing id")
		return
	}
	resp, err := a.c.SLA.GetSLA(r.Context(), &slav1.GetSLARequest{TicketId: id})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, slaJSON(resp.GetSla()))
}

// GetTicketCard aggregates ticket + assignment + sla + timeline for the UI.
func (a *API) GetTicketCard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErr(w, http.StatusBadRequest, "missing id")
		return
	}

	ctx := r.Context()
	var (
		ticketResp *ticketv1.GetTicketResponse
		assignResp *assignmentv1.GetAssignmentResponse
		slaResp    *slav1.GetSLAResponse
		tlResp     *auditv1.GetTimelineResponse
		ticketErr, assignErr, slaErr, tlErr error
		wg         sync.WaitGroup
	)
	wg.Add(4)
	go func() {
		defer wg.Done()
		ticketResp, ticketErr = a.c.Ticket.GetTicket(ctx, &ticketv1.GetTicketRequest{Id: id})
	}()
	go func() {
		defer wg.Done()
		assignResp, assignErr = a.c.Assignment.GetAssignment(ctx, &assignmentv1.GetAssignmentRequest{TicketId: id})
	}()
	go func() {
		defer wg.Done()
		slaResp, slaErr = a.c.SLA.GetSLA(ctx, &slav1.GetSLARequest{TicketId: id})
	}()
	go func() {
		defer wg.Done()
		tlResp, tlErr = a.c.Audit.GetTimeline(ctx, &auditv1.GetTimelineRequest{TicketId: id})
	}()
	wg.Wait()

	if ticketErr != nil {
		writeGRPCErr(w, ticketErr)
		return
	}

	card := map[string]any{
		"ticket": ticketJSON(ticketResp.GetTicket()),
	}
	if assignErr == nil && assignResp.GetAssignment() != nil {
		as := assignResp.GetAssignment()
		card["assignment"] = map[string]any{
			"ticket_id":     as.GetTicketId(),
			"assignee_id":   as.GetAssigneeId(),
			"assignee_name": as.GetAssigneeName(),
			"assigned_at":   as.GetAssignedAt(),
		}
	} else {
		card["assignment"] = nil
	}
	if slaErr == nil && slaResp.GetSla() != nil {
		card["sla"] = slaJSON(slaResp.GetSla())
	} else {
		card["sla"] = nil
	}
	events := make([]map[string]any, 0)
	if tlErr == nil {
		for _, e := range tlResp.GetEvents() {
			events = append(events, map[string]any{
				"id":          e.GetId(),
				"ticket_id":   e.GetTicketId(),
				"event_id":    e.GetEventId(),
				"event_type":  e.GetEventType(),
				"title":       e.GetTitle(),
				"detail":      e.GetDetail(),
				"actor":       e.GetActor(),
				"occurred_at": e.GetOccurredAt(),
			})
		}
	}
	card["timeline"] = events
	writeJSON(w, http.StatusOK, card)
}

func (a *API) SearchTickets(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	limit := int32(20)
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = int32(n)
		}
	}
	resp, err := a.c.Search.SearchTickets(r.Context(), &searchv1.SearchTicketsRequest{
		Query: q,
		Limit: limit,
	})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	hits := make([]map[string]any, 0, len(resp.GetHits()))
	for _, h := range resp.GetHits() {
		hits = append(hits, map[string]any{
			"id":          h.GetId(),
			"title":       h.GetTitle(),
			"description": h.GetDescription(),
			"status":      h.GetStatus(),
			"priority":    h.GetPriority(),
			"category":    h.GetCategory(),
			"requester":   h.GetRequester(),
			"assignee_id": h.GetAssigneeId(),
			"updated_at":  h.GetUpdatedAt(),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"hits":            hits,
		"estimated_total": resp.GetEstimatedTotal(),
		"query":           q,
	})
}

func slaJSON(s *slav1.SLA) map[string]any {
	return map[string]any{
		"ticket_id":          s.GetTicketId(),
		"state":              slaStateString(s.GetState()),
		"first_response_due": s.GetFirstResponseDue(),
		"resolve_due":        s.GetResolveDue(),
		"warned_at":          s.GetWarnedAt(),
		"breached_at":        s.GetBreachedAt(),
		"policy":             s.GetPolicy(),
	}
}

func slaStateString(s slav1.SLAState) string {
	switch s {
	case slav1.SLAState_SLA_STATE_OK:
		return "ok"
	case slav1.SLAState_SLA_STATE_WARNING:
		return "warning"
	case slav1.SLAState_SLA_STATE_BREACHED:
		return "breached"
	case slav1.SLAState_SLA_STATE_CANCELLED:
		return "cancelled"
	default:
		return "unspecified"
	}
}

type updateStatusBody struct {
	Status string `json:"status"`
}

func (a *API) UpdateTicketStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeErr(w, http.StatusBadRequest, "missing id")
		return
	}
	var body updateStatusBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}
	resp, err := a.c.Ticket.UpdateTicketStatus(r.Context(), &ticketv1.UpdateTicketStatusRequest{
		Id:     id,
		Status: parseStatus(body.Status),
	})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ticketJSON(resp.GetTicket()))
}

func ticketJSON(t *ticketv1.Ticket) map[string]any {
	return map[string]any{
		"id":          t.GetId(),
		"title":       t.GetTitle(),
		"description": t.GetDescription(),
		"status":      statusString(t.GetStatus()),
		"priority":    priorityString(t.GetPriority()),
		"category":    t.GetCategory(),
		"requester":   t.GetRequester(),
		"assignee_id": t.GetAssigneeId(),
		"created_at":  t.GetCreatedAt(),
		"updated_at":  t.GetUpdatedAt(),
	}
}

func userJSON(u *authv1.User) map[string]any {
	return map[string]any{
		"id":    u.GetId(),
		"email": u.GetEmail(),
		"name":  u.GetName(),
		"role":  roleString(u.GetRole()),
	}
}

func parsePriority(s string) ticketv1.TicketPriority {
	switch strings.ToLower(s) {
	case "low":
		return ticketv1.TicketPriority_TICKET_PRIORITY_LOW
	case "high":
		return ticketv1.TicketPriority_TICKET_PRIORITY_HIGH
	case "urgent":
		return ticketv1.TicketPriority_TICKET_PRIORITY_URGENT
	default:
		return ticketv1.TicketPriority_TICKET_PRIORITY_NORMAL
	}
}

func parseStatus(s string) ticketv1.TicketStatus {
	switch strings.ToLower(s) {
	case "new":
		return ticketv1.TicketStatus_TICKET_STATUS_NEW
	case "open":
		return ticketv1.TicketStatus_TICKET_STATUS_OPEN
	case "pending":
		return ticketv1.TicketStatus_TICKET_STATUS_PENDING
	case "resolved":
		return ticketv1.TicketStatus_TICKET_STATUS_RESOLVED
	case "closed":
		return ticketv1.TicketStatus_TICKET_STATUS_CLOSED
	default:
		return ticketv1.TicketStatus_TICKET_STATUS_UNSPECIFIED
	}
}

func statusString(s ticketv1.TicketStatus) string {
	switch s {
	case ticketv1.TicketStatus_TICKET_STATUS_NEW:
		return "new"
	case ticketv1.TicketStatus_TICKET_STATUS_OPEN:
		return "open"
	case ticketv1.TicketStatus_TICKET_STATUS_PENDING:
		return "pending"
	case ticketv1.TicketStatus_TICKET_STATUS_RESOLVED:
		return "resolved"
	case ticketv1.TicketStatus_TICKET_STATUS_CLOSED:
		return "closed"
	default:
		return "unspecified"
	}
}

func priorityString(p ticketv1.TicketPriority) string {
	switch p {
	case ticketv1.TicketPriority_TICKET_PRIORITY_LOW:
		return "low"
	case ticketv1.TicketPriority_TICKET_PRIORITY_HIGH:
		return "high"
	case ticketv1.TicketPriority_TICKET_PRIORITY_URGENT:
		return "urgent"
	default:
		return "normal"
	}
}

func roleString(r authv1.Role) string {
	switch r {
	case authv1.Role_ROLE_AGENT:
		return "agent"
	case authv1.Role_ROLE_ADMIN:
		return "admin"
	case authv1.Role_ROLE_REQUESTER:
		return "requester"
	default:
		return "unspecified"
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func writeGRPCErr(w http.ResponseWriter, err error) {
	msg := err.Error()
	code := http.StatusBadGateway
	if strings.Contains(msg, "visitor mismatch") || strings.Contains(msg, "site_key mismatch") {
		code = http.StatusForbidden
	} else if strings.Contains(msg, "NotFound") || strings.Contains(msg, "not found") {
		code = http.StatusNotFound
	} else if strings.Contains(msg, "InvalidArgument") || strings.Contains(msg, "required") {
		code = http.StatusBadRequest
	}
	writeErr(w, code, msg)
}
