package validator

import (
	"regexp"

	ssov1 "github.com/VariableSan/gia-protos/gen/go/sso"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validate *validator.Validate
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func init() {
	validate = validator.New()
}

// Register custom validators
func customValidationExample(validate *validator.Validate) {
	/*
		type StructExampleUsage struct {
			Email    string `validate:"required,validEmail"`
			Password string `validate:"required,min=6"`
			AppID    int32  `validate:"required,gt=0"`
		}
	*/
	_ = validate.RegisterValidation(
		"validEmail",
		func(fl validator.FieldLevel) bool {
			return emailRegex.MatchString(fl.Field().String())
		},
	)
}

// Validate validates any struct and returns a gRPC error if validation fails
func Validate(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		if len(validationErrors) > 0 {
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", validationErrors[0])
		}
	}
	return nil
}

// LoginRequestValidator validates LoginRequest
type LoginRequestValidator struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
	AppID    int32  `validate:"required,gt=0"`
}

// RegisterRequestValidator validates RegisterRequest
type RegisterRequestValidator struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
}

// LogoutRequestValidator validates LogoutRequest
type LogoutRequestValidator struct {
	Token string `validate:"required"`
}

// IsAdminRequestValidator validates IsAdminRequest
type IsAdminRequestValidator struct {
	UserID int64 `validate:"required,gt=0"`
}

// ValidateLoginRequest validates LoginRequest fields
func ValidateLoginRequest(req *ssov1.LoginRequest) error {
	return Validate(LoginRequestValidator{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		AppID:    req.GetAppId(),
	})
}

// ValidateRegisterRequest validates RegisterRequest fields
func ValidateRegisterRequest(req *ssov1.RegisterRequest) error {
	return Validate(RegisterRequestValidator{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
}

// ValidateLogoutRequest validates LogoutRequest fields
func ValidateLogoutRequest(req *ssov1.LogoutRequest) error {
	return Validate(LogoutRequestValidator{
		Token: req.GetToken(),
	})
}

// ValidateIsAdminRequest validates IsAdminRequest fields
func ValidateIsAdminRequest(req *ssov1.IsAdminRequest) error {
	return Validate(IsAdminRequestValidator{
		UserID: req.GetUserId(),
	})
}
