package store

import (
	"errors"
	"sync"

	authv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/auth/v1"
	"golang.org/x/crypto/bcrypt"
)

var ErrNotFound = errors.New("user not found")
var ErrInvalidCredentials = errors.New("invalid credentials")

type User struct {
	ID           string
	Email        string
	Name         string
	Role         authv1.Role
	PasswordHash string
}

type Store struct {
	mu    sync.RWMutex
	byID  map[string]User
	email map[string]string
}

func NewWithSeed() *Store {
	s := &Store{
		byID:  make(map[string]User),
		email: make(map[string]string),
	}
	// Dev seed: agent@helpdesk.local / password
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	u := User{
		ID:           "a-1",
		Email:        "agent@helpdesk.local",
		Name:         "Алексей К.",
		Role:         authv1.Role_ROLE_AGENT,
		PasswordHash: string(hash),
	}
	s.byID[u.ID] = u
	s.email[u.Email] = u.ID

	adminHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	admin := User{
		ID:           "a-admin",
		Email:        "admin@helpdesk.local",
		Name:         "Admin",
		Role:         authv1.Role_ROLE_ADMIN,
		PasswordHash: string(adminHash),
	}
	s.byID[admin.ID] = admin
	s.email[admin.Email] = admin.ID
	return s
}

func (s *Store) Authenticate(email, password string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.email[email]
	if !ok {
		return User{}, ErrInvalidCredentials
	}
	u := s.byID[id]
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return User{}, ErrInvalidCredentials
	}
	return u, nil
}

func (s *Store) Get(id string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.byID[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

func (u User) Proto() *authv1.User {
	return &authv1.User{
		Id:    u.ID,
		Email: u.Email,
		Name:  u.Name,
		Role:  u.Role,
	}
}
