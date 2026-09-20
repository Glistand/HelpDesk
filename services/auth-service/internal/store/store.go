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
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	seeds := []User{
		{ID: "a-1", Email: "agent@helpdesk.local", Name: "Алексей К.", Role: authv1.Role_ROLE_AGENT},
		{ID: "a-2", Email: "maria@helpdesk.local", Name: "Марина С.", Role: authv1.Role_ROLE_AGENT},
		{ID: "a-3", Email: "ivan@helpdesk.local", Name: "Денис В.", Role: authv1.Role_ROLE_AGENT},
		{ID: "a-admin", Email: "admin@helpdesk.local", Name: "Admin", Role: authv1.Role_ROLE_ADMIN},
		{ID: "r-1", Email: "anna@company.local", Name: "Анна П.", Role: authv1.Role_ROLE_REQUESTER},
		{ID: "r-2", Email: "peter@company.local", Name: "Пётр В.", Role: authv1.Role_ROLE_REQUESTER},
	}
	for _, u := range seeds {
		u.PasswordHash = string(hash)
		s.byID[u.ID] = u
		s.email[u.Email] = u.ID
	}
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
