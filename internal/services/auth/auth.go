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
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models models.User, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (auth *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
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

	app, err := auth.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, app, auth.tokenTTL)
	if err != nil {
		auth.log.Error("failed to generate token")
		return "", fmt.Errorf("%s: %w", operation, err)
	}

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

	id, err := auth.userSaver.SaveUser(ctx, email, passHash)
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
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("user not found")
			return false, fmt.Errorf("%s: %w", operation, storage.ErrInvalidAppID)
		}

		return false, fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
