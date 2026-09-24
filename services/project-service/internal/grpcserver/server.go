package grpcserver

import (
	"context"
	"strings"
	"time"

	projectv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/project/v1"
	"github.com/Glistand/HelpDesk/libs/grpckit/statuserr"
	"github.com/Glistand/HelpDesk/services/project-service/internal/repository"
)

type Server struct {
	projectv1.UnimplementedProjectServiceServer
	repo *repository.Repo
}

func New(repo *repository.Repo) *Server { return &Server{repo: repo} }

func (s *Server) GetProfile(ctx context.Context, _ *projectv1.GetProfileRequest) (*projectv1.GetProfileResponse, error) {
	profile, err := s.repo.Get(ctx)
	if err != nil {
		return nil, statuserr.Internal(err.Error())
	}
	return &projectv1.GetProfileResponse{Profile: toProto(profile)}, nil
}

func (s *Server) UpdateProfile(ctx context.Context, req *projectv1.UpdateProfileRequest) (*projectv1.UpdateProfileResponse, error) {
	if req.GetProfile() == nil || strings.TrimSpace(req.GetProfile().GetName()) == "" || strings.TrimSpace(req.GetProfile().GetSiteKey()) == "" {
		return nil, statuserr.InvalidArgument("profile name and site_key are required")
	}
	profile, err := s.repo.Update(ctx, fromProto(req.GetProfile()))
	if err != nil {
		return nil, statuserr.Internal(err.Error())
	}
	return &projectv1.UpdateProfileResponse{Profile: toProto(profile)}, nil
}

func toProto(profile repository.Profile) *projectv1.ProjectProfile {
	return &projectv1.ProjectProfile{
		Name:             profile.Name,
		Description:      profile.Description,
		Domain:           profile.Domain,
		WebsiteUrl:       profile.WebsiteURL,
		SupportEmail:     profile.SupportEmail,
		SupportPhone:     profile.SupportPhone,
		Timezone:         profile.Timezone,
		SiteKey:          profile.SiteKey,
		BotInstructions:  profile.BotInstructions,
		AutoCreateTicket: profile.AutoCreateTicket,
		UpdatedAt:        profile.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func fromProto(profile *projectv1.ProjectProfile) repository.Profile {
	return repository.Profile{
		Name:             strings.TrimSpace(profile.GetName()),
		Description:      strings.TrimSpace(profile.GetDescription()),
		Domain:           strings.TrimSpace(profile.GetDomain()),
		WebsiteURL:       strings.TrimSpace(profile.GetWebsiteUrl()),
		SupportEmail:     strings.TrimSpace(profile.GetSupportEmail()),
		SupportPhone:     strings.TrimSpace(profile.GetSupportPhone()),
		Timezone:         strings.TrimSpace(profile.GetTimezone()),
		SiteKey:          strings.TrimSpace(profile.GetSiteKey()),
		BotInstructions:  strings.TrimSpace(profile.GetBotInstructions()),
		AutoCreateTicket: profile.GetAutoCreateTicket(),
	}
}
