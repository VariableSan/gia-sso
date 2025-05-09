package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/VariableSan/gia-sso/internal/domain/models"
	"github.com/VariableSan/gia-sso/internal/storage"
	"github.com/VariableSan/gia-sso/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userProvider UserProvider
	tokenTTL     time.Duration
}

type UserProvider interface {
	User(ctx context.Context, email string) (models models.User, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (int64, error)
}

type Provider interface {
	UserProvider
}

func New(
	log *slog.Logger,
	provider Provider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userProvider: provider,
		tokenTTL:     tokenTTL,
	}
}

func (auth *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	jwtSecret string,
) (string, error) {
	const operation = "auth.Login"

	log := auth.log.With(
		slog.String("operation", operation),
	)

	log.Info("attempting to login user")

	user, err := auth.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			auth.log.Warn("user not found")
			return "", fmt.Errorf("%s: %w", operation, storage.ErrInvalidCredentials)
		}
		auth.log.Error("failed to get user")
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		auth.log.Error("invalid credentials")
		return "", fmt.Errorf("%s: %w", operation, storage.ErrInvalidCredentials)
	}

	token, err := jwt.NewToken(user, jwtSecret, auth.tokenTTL)
	if err != nil {
		auth.log.Error("failed to generate token")
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("user logged in successfully")

	return token, nil
}

func (auth *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (userID int64, err error) {
	const operation = "auth.RegisterNewUser"

	log := auth.log.With(
		slog.String("operation", operation),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash")
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	id, err := auth.userProvider.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists")
			return 0, fmt.Errorf("%s: %w", operation, storage.ErrUserExists)
		}

		log.Error("failed to save user")
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("user registered")

	return id, nil
}

func (auth *Auth) IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	const operation = "auth.IsAdmin"

	log := auth.log.With(
		slog.String("operation", operation),
	)

	log.Info("checking if user is admin")

	isAdmin, err := auth.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
