package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"goschool/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

const testSecret = "01234567890123456789012345678901"

// InMemoryDB — заглушка хранилища для тестов.
// Методы на указателе: иначе мьютекс и append работали бы с копией структуры.
type InMemoryDB struct {
	mu    sync.RWMutex
	users []models.User
}

func (s *InMemoryDB) AddUser(_ context.Context, u *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	u.Id = int64(len(s.users) + 1)
	s.users = append(s.users, *u)
	return nil
}

func (s *InMemoryDB) GetUserById(_ context.Context, id int64) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.users {
		if s.users[i].Id == id {
			u := s.users[i]
			return &u, nil
		}
	}
	return nil, pgx.ErrNoRows
}

func (s *InMemoryDB) GetAllUsers(_ context.Context) ([]models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]models.User, len(s.users))
	copy(users, s.users)
	return users, nil
}

func (s *InMemoryDB) DeleteUserById(_ context.Context, _ int64) error {
	return nil
}

func (s *InMemoryDB) GetUserByEmail(_ context.Context, email string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.users {
		if s.users[i].Email == email {
			u := s.users[i]
			return &u, nil
		}
	}
	return nil, pgx.ErrNoRows
}

func newTestRouter(store models.Storage) *gin.Engine {
	gin.SetMode(gin.TestMode)

	hs := NewStorageHandler(store, []byte(testSecret))

	r := gin.New()
	r.POST("/users/login", hs.Login)

	admin := r.Group("/admin", AuthMiddleware([]byte(testSecret)), AdminMiddleware())
	admin.GET("/User", hs.GetAllUsers)
	admin.POST("/User", hs.AddUser)

	return r
}

func makeUser(t *testing.T, id int64, email, password, role string) models.User {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	return models.User{
		Id:       id,
		Name:     "Test User",
		Email:    email,
		Password: string(hash),
		Role:     role,
	}
}

func doRequest(r http.Handler, method, path, token, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func login(t *testing.T, r http.Handler, email, password string) string {
	t.Helper()

	body := fmt.Sprintf(`{"email":%q,"password":%q}`, email, password)
	w := doRequest(r, http.MethodPost, "/users/login", "", body)
	if w.Code != http.StatusOK {
		t.Fatalf("login: want 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("login: decode response: %v", err)
	}
	if resp.AccessToken == "" {
		t.Fatal("login: empty access_token")
	}

	return resp.AccessToken
}

// Главный тест: логин -> токен -> доступ в админку.
func TestLoginThenAccessAdmin(t *testing.T) {
	store := &InMemoryDB{users: []models.User{
		makeUser(t, 1, "admin@example.com", "secret123", models.RoleAdmin),
	}}
	r := newTestRouter(store)

	token := login(t, r, "admin@example.com", "secret123")

	w := doRequest(r, http.MethodGet, "/admin/User", token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("admin: want 200, got %d: %s", w.Code, w.Body.String())
	}
	if strings.Contains(w.Body.String(), "password") {
		t.Errorf("admin: response leaks password field: %s", w.Body.String())
	}
}

func TestAdminWithoutToken(t *testing.T) {
	r := newTestRouter(&InMemoryDB{})

	w := doRequest(r, http.MethodGet, "/admin/User", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", w.Code)
	}
}

func TestAdminWithStudentToken(t *testing.T) {
	store := &InMemoryDB{users: []models.User{
		makeUser(t, 1, "student@example.com", "secret123", models.RoleStudent),
	}}
	r := newTestRouter(store)

	token := login(t, r, "student@example.com", "secret123")

	w := doRequest(r, http.MethodGet, "/admin/User", token, "")
	if w.Code != http.StatusForbidden {
		t.Fatalf("want 403, got %d", w.Code)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	store := &InMemoryDB{users: []models.User{
		makeUser(t, 1, "admin@example.com", "secret123", models.RoleAdmin),
	}}
	r := newTestRouter(store)

	body := `{"email":"admin@example.com","password":"wrong"}`
	w := doRequest(r, http.MethodPost, "/users/login", "", body)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", w.Code)
	}
}

func TestLoginUnknownEmail(t *testing.T) {
	r := newTestRouter(&InMemoryDB{})

	body := `{"email":"nobody@example.com","password":"secret123"}`
	w := doRequest(r, http.MethodPost, "/users/login", "", body)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", w.Code)
	}
}

// Логин не должен требовать email-формат: в проекте админ входит по "roro".
func TestLoginWithNonEmailIdentifier(t *testing.T) {
	store := &InMemoryDB{users: []models.User{
		makeUser(t, 1, "roro", "roro", models.RoleAdmin),
	}}
	r := newTestRouter(store)

	token := login(t, r, "roro", "roro")

	w := doRequest(r, http.MethodGet, "/admin/User", token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("admin: want 200, got %d", w.Code)
	}
}

// Токен, подписанный другим секретом, не должен приниматься.
func TestTokenFromAnotherSecret(t *testing.T) {
	otherSecret := []byte(strings.Repeat("x", 32))
	token, err := GenerateJWT(otherSecret, 1, models.RoleAdmin)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	r := newTestRouter(&InMemoryDB{})
	w := doRequest(r, http.MethodGet, "/admin/User", token, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", w.Code)
	}
}
