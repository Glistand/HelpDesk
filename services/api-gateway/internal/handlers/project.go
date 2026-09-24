package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	authv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/auth/v1"
	projectv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/project/v1"
	"github.com/Glistand/HelpDesk/services/api-gateway/internal/middleware"
)

func (a *API) GetProject(w http.ResponseWriter, r *http.Request) {
	resp, err := a.c.Project.GetProfile(r.Context(), &projectv1.GetProfileRequest{})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, projectJSON(resp.GetProfile()))
}

func (a *API) UpdateProject(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || user.GetRole() != authv1.Role_ROLE_ADMIN {
		writeErr(w, http.StatusForbidden, "administrator role required")
		return
	}
	var input projectInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}
	resp, err := a.c.Project.UpdateProfile(r.Context(), &projectv1.UpdateProfileRequest{Profile: input.proto()})
	if err != nil {
		writeGRPCErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, projectJSON(resp.GetProfile()))
}

type projectInput struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	Domain           string `json:"domain"`
	WebsiteURL       string `json:"website_url"`
	SupportEmail     string `json:"support_email"`
	SupportPhone     string `json:"support_phone"`
	Timezone         string `json:"timezone"`
	SiteKey          string `json:"site_key"`
	BotInstructions  string `json:"bot_instructions"`
	AutoCreateTicket bool   `json:"auto_create_ticket"`
}

func (p projectInput) proto() *projectv1.ProjectProfile {
	return &projectv1.ProjectProfile{
		Name:             strings.TrimSpace(p.Name),
		Description:      strings.TrimSpace(p.Description),
		Domain:           strings.TrimSpace(p.Domain),
		WebsiteUrl:       strings.TrimSpace(p.WebsiteURL),
		SupportEmail:     strings.TrimSpace(p.SupportEmail),
		SupportPhone:     strings.TrimSpace(p.SupportPhone),
		Timezone:         strings.TrimSpace(p.Timezone),
		SiteKey:          strings.TrimSpace(p.SiteKey),
		BotInstructions:  strings.TrimSpace(p.BotInstructions),
		AutoCreateTicket: p.AutoCreateTicket,
	}
}

func projectJSON(profile *projectv1.ProjectProfile) map[string]any {
	if profile == nil {
		return nil
	}
	return map[string]any{
		"name":               profile.GetName(),
		"description":        profile.GetDescription(),
		"domain":             profile.GetDomain(),
		"website_url":        profile.GetWebsiteUrl(),
		"support_email":      profile.GetSupportEmail(),
		"support_phone":      profile.GetSupportPhone(),
		"timezone":           profile.GetTimezone(),
		"site_key":           profile.GetSiteKey(),
		"bot_instructions":   profile.GetBotInstructions(),
		"auto_create_ticket": profile.GetAutoCreateTicket(),
		"updated_at":         profile.GetUpdatedAt(),
		"server_time":        time.Now().UTC().Format(time.RFC3339),
	}
}
