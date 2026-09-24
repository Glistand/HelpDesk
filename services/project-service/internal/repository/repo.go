package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Glistand/HelpDesk/services/project-service/internal/config"
)

type Profile struct {
	Name             string
	Description      string
	Domain           string
	WebsiteURL       string
	SupportEmail     string
	SupportPhone     string
	Timezone         string
	SiteKey          string
	BotInstructions  string
	AutoCreateTicket bool
	UpdatedAt        time.Time
}

type Repo struct{ db *sql.DB }

func New(db *sql.DB) *Repo { return &Repo{db: db} }

func (r *Repo) Bootstrap(ctx context.Context, seed config.BootstrapProfile) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO project_profile (
			singleton, name, description, domain, website_url, support_email, support_phone,
			timezone, site_key, bot_instructions, auto_create_ticket
		) VALUES (TRUE, $1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		ON CONFLICT (singleton) DO NOTHING`,
		seed.Name, seed.Description, seed.Domain, seed.WebsiteURL, seed.SupportEmail,
		seed.SupportPhone, seed.Timezone, seed.SiteKey, seed.BotInstructions, seed.AutoCreateTicket,
	)
	return err
}

func (r *Repo) Get(ctx context.Context) (Profile, error) {
	var profile Profile
	err := r.db.QueryRowContext(ctx, `
		SELECT name, description, domain, website_url, support_email, support_phone,
			timezone, site_key, bot_instructions, auto_create_ticket, updated_at
		FROM project_profile WHERE singleton=TRUE`,
	).Scan(
		&profile.Name, &profile.Description, &profile.Domain, &profile.WebsiteURL,
		&profile.SupportEmail, &profile.SupportPhone, &profile.Timezone, &profile.SiteKey,
		&profile.BotInstructions, &profile.AutoCreateTicket, &profile.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return Profile{}, fmt.Errorf("profile not initialized")
	}
	return profile, err
}

func (r *Repo) Update(ctx context.Context, profile Profile) (Profile, error) {
	err := r.db.QueryRowContext(ctx, `
		UPDATE project_profile SET
			name=$1, description=$2, domain=$3, website_url=$4, support_email=$5,
			support_phone=$6, timezone=$7, site_key=$8, bot_instructions=$9,
			auto_create_ticket=$10, updated_at=now()
		WHERE singleton=TRUE
		RETURNING name, description, domain, website_url, support_email, support_phone,
			timezone, site_key, bot_instructions, auto_create_ticket, updated_at`,
		profile.Name, profile.Description, profile.Domain, profile.WebsiteURL,
		profile.SupportEmail, profile.SupportPhone, profile.Timezone, profile.SiteKey,
		profile.BotInstructions, profile.AutoCreateTicket,
	).Scan(
		&profile.Name, &profile.Description, &profile.Domain, &profile.WebsiteURL,
		&profile.SupportEmail, &profile.SupportPhone, &profile.Timezone, &profile.SiteKey,
		&profile.BotInstructions, &profile.AutoCreateTicket, &profile.UpdatedAt,
	)
	return profile, err
}
