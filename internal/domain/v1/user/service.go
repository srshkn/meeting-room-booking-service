package user

import (
	"context"
	"errors"
	"strings"

	"github.com/oapi-codegen/runtime/types"

	"mrb-service/internal/db"
	v1GenAPI "mrb-service/internal/generated/v1"
	"mrb-service/internal/hash"
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

func (u *userService) validateCreateUser(
	ctx context.Context,
	username, rowEmail string,
) error {
	emailUser := strings.ToLower(rowEmail)

	user, _ := u.repo.CheckUserExists(
		ctx,
		db.CheckUserExistsParams{
			Username: username,
			Email:    emailUser,
		},
	)

	if emailUser == user.Email {
		return errors.New("a user with this email already exists")
	}

	if username == user.Username {
		return errors.New("a user with that username already exists")
	}

	return nil
}

func (u *userService) Registration(
	ctx context.Context,
	body v1GenAPI.PostUserRegisterJSONRequestBody,
) (*v1GenAPI.UserResponse, error) {
	if err := u.validateCreateUser(ctx, body.Username, string(body.Email)); err != nil {
		return &v1GenAPI.UserResponse{}, err
	}

	hash, err := hash.Hash(body.Password)
	if err != nil {
		return &v1GenAPI.UserResponse{}, err
	}

	createUser, err := u.repo.CreateUser(
		ctx,
		db.CreateUserParams{
			Username:     body.Username,
			Email:        strings.ToLower(string(body.Email)),
			PasswordHash: hash,
		},
	)
	if err != nil {
		return &v1GenAPI.UserResponse{}, err
	}

	return &v1GenAPI.UserResponse{
		Id:       createUser.ID,
		Username: createUser.Username,
		Email:    types.Email(createUser.Email),
		Role:     v1GenAPI.UserResponseRole(createUser.Role),
	}, nil
}
