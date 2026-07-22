package service_test

import (
	"context"
	"testing"

	"github.com/alexperezortuno/streaming/internal/config"
	"github.com/alexperezortuno/streaming/internal/model"
	"github.com/alexperezortuno/streaming/internal/repository/mock"
	"github.com/alexperezortuno/streaming/internal/service"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login_Success(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	userRepo := &mock.UserRepository{
		FindByUsernameFunc: func(_ context.Context, username string) (*model.User, error) {
			return &model.User{
				ID:       "user-1",
				Username: username,
				Password: string(hashed),
				Role:     model.RoleAdmin,
			}, nil
		},
	}

	cfg := &config.Config{JWTSecret: "test-secret", JWTExpiration: 3600000000000}
	svc := service.NewAuthService(userRepo, cfg)

	resp, err := svc.Login(context.Background(), model.LoginRequest{
		Username: "admin",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected a token")
	}
	if resp.User.Username != "admin" {
		t.Fatalf("expected username admin, got: %s", resp.User.Username)
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	userRepo := &mock.UserRepository{
		FindByUsernameFunc: func(_ context.Context, username string) (*model.User, error) {
			return &model.User{
				ID:       "user-1",
				Username: username,
				Password: string(hashed),
				Role:     model.RoleUser,
			}, nil
		},
	}

	cfg := &config.Config{JWTSecret: "test-secret"}
	svc := service.NewAuthService(userRepo, cfg)

	_, err := svc.Login(context.Background(), model.LoginRequest{
		Username: "admin",
		Password: "wrong-password",
	})
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	userRepo := &mock.UserRepository{
		FindByUsernameFunc: func(_ context.Context, _ string) (*model.User, error) {
			return nil, nil
		},
	}

	cfg := &config.Config{JWTSecret: "test-secret"}
	svc := service.NewAuthService(userRepo, cfg)

	_, err := svc.Login(context.Background(), model.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
}

func TestAuthService_Register_Success(t *testing.T) {
	userRepo := &mock.UserRepository{
		FindByUsernameFunc: func(_ context.Context, _ string) (*model.User, error) {
			return nil, nil
		},
		CreateFunc: func(_ context.Context, req model.RegisterRequest, hashedPassword string) (*model.User, error) {
			return &model.User{
				ID:       "new-user",
				Username: req.Username,
				Role:     model.RoleUser,
			}, nil
		},
	}

	cfg := &config.Config{JWTSecret: "test-secret"}
	svc := service.NewAuthService(userRepo, cfg)

	user, err := svc.Register(context.Background(), model.RegisterRequest{
		Username: "newuser",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if user.Username != "newuser" {
		t.Fatalf("expected username newuser, got: %s", user.Username)
	}
}

func TestAuthService_Register_Duplicate(t *testing.T) {
	userRepo := &mock.UserRepository{
		FindByUsernameFunc: func(_ context.Context, username string) (*model.User, error) {
			return &model.User{ID: "existing", Username: username}, nil
		},
	}

	cfg := &config.Config{JWTSecret: "test-secret"}
	svc := service.NewAuthService(userRepo, cfg)

	_, err := svc.Register(context.Background(), model.RegisterRequest{
		Username: "existing",
		Password: "password123",
	})
	if err != model.ErrDuplicate {
		t.Fatalf("expected ErrDuplicate, got: %v", err)
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	userRepo := &mock.UserRepository{
		FindByUsernameFunc: func(_ context.Context, username string) (*model.User, error) {
			return &model.User{ID: "u1", Username: username, Password: string(hashed), Role: model.RoleAdmin}, nil
		},
	}

	cfg := &config.Config{JWTSecret: "test-secret", JWTExpiration: 3600000000000}
	svc := service.NewAuthService(userRepo, cfg)

	resp, err := svc.Login(context.Background(), model.LoginRequest{Username: "admin", Password: "password"})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	user, err := svc.ValidateToken(resp.Token)
	if err != nil {
		t.Fatalf("validate token failed: %v", err)
	}
	if user.Username != "admin" {
		t.Fatalf("expected username admin, got: %s", user.Username)
	}
}

func TestAuthService_ValidateToken_Invalid(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	svc := service.NewAuthService(&mock.UserRepository{}, cfg)

	_, err := svc.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}
