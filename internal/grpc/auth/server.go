package auth

import (
	"context"

	ssov1 "github.com/VariableSan/gia-protos/gen/go/sso"
	"github.com/VariableSan/gia-sso/pkg/validator"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validator.ValidateLoginRequest(req); err != nil {
		return nil, err
	}

	// TODO: implement login via auth service

	return &ssov1.LoginResponse{
		Token: "token",
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if err := validator.ValidateRegisterRequest(req); err != nil {
		return nil, err
	}

	// TODO: implement registration logic

	return &ssov1.RegisterResponse{
		UserId: 123, // Replace with actual user ID
	}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	if err := validator.ValidateIsAdminRequest(req); err != nil {
		return nil, err
	}

	// TODO: implement admin check logic

	return &ssov1.IsAdminResponse{
		IsAdmin: false, // Replace with actual check
	}, nil
}

func (s *serverAPI) Logout(
	ctx context.Context,
	req *ssov1.LogoutRequest,
) (*ssov1.LogoutResponse, error) {
	if err := validator.ValidateLogoutRequest(req); err != nil {
		return nil, err
	}

	// TODO: implement logout logic

	return &ssov1.LogoutResponse{
		Success: true,
	}, nil
}
