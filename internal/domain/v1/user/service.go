package user

import (
	"context"
	"strings"

	"github.com/oapi-codegen/runtime/types"

	"mrb-service/internal/db"
	v1GenAPI "mrb-service/internal/generated/v1"
	"mrb-service/internal/password"
)

type UserService interface {
	Registration(ctx context.Context, body v1GenAPI.PostUserRegisterJSONRequestBody) (*v1GenAPI.UserResponse, error)
}

type userService struct {
	repo repository
}

func NewService(repo repository) *userService {
	return &userService{
		repo: repo,
	}
}

func (u *userService) validateRegistration(
	username, email, passwordRow, confirmation string,
) error {
	var err error

	switch {
	case username == "":
		err = ErrEmptyUsername

	case len([]rune(username)) > 25:
		err = ErrLongUsername

	case email == "":
		err = ErrEmptyEmail

	case len([]rune(email)) > 254:
		err = ErrLongEmail

	case passwordRow == "":
		err = ErrEmptyPassword

	case len([]rune(passwordRow)) < 8:
		err = ErrShortPassword

	case len([]rune(passwordRow)) > 128:
		err = ErrLongPassword

	case passwordRow != confirmation:
		err = ErrPasswordsDoNotMatch

	default:
		err = nil
	}

	return err
}

func (u *userService) Registration(
	ctx context.Context,
	body v1GenAPI.PostUserRegisterJSONRequestBody,
) (*v1GenAPI.UserResponse, error) {
	username := strings.TrimSpace(body.Username)
	email := strings.TrimSpace(strings.ToLower(string(body.Email)))
	passwordRow := body.Password
	confirmation := body.Confirmation

	if err := u.validateRegistration(
		username,
		email,
		passwordRow,
		confirmation,
	); err != nil {
		return nil, err
	}

	passwordHash, err := password.Hash(passwordRow)
	if err != nil {
		return nil, err
	}

	createdUser, err := u.repo.CreateUser(
		ctx,
		db.CreateUserParams{
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		},
	)
	if err != nil {
		return nil, mapCreateUserError(err)
	}

	return &v1GenAPI.UserResponse{
		Id:       createdUser.ID,
		Username: createdUser.Username,
		Email:    types.Email(createdUser.Email),
		Role:     v1GenAPI.UserResponseRole(createdUser.Role),
	}, nil
}
