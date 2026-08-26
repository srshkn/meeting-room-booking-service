package auth

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"mrb-service/internal/db"
	"mrb-service/internal/jwt"
	"mrb-service/internal/password"

	v1GenAPI "mrb-service/internal/generated/v1"
)

var (
	dummyUserID  = uuid.MustParse("01a02084-0c80-7a11-8a11-111111111111")
	dummyAdminID = uuid.MustParse("01a02084-0c80-7a22-8a22-222222222222")
)

type AuthService interface {
	Login(
		ctx context.Context,
		body v1GenAPI.PostAuthLoginJSONRequestBody,
	) (*v1GenAPI.Token, jwt.RefreshToken, error)
	LoginDummy(
		ctx context.Context,
		body v1GenAPI.PostAuthDummyLoginJSONRequestBody,
	) (*v1GenAPI.Token, jwt.RefreshToken, error)
}

type authService struct {
	repo repository
	jwt  jwt.Manager
}

func NewService(repo repository, jwt jwt.Manager) *authService {
	return &authService{
		repo: repo,
		jwt:  jwt,
	}
}

func (s *authService) validateLogin(
	email, passwordRow string,
) error {
	var err error

	switch {
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

	default:
		err = nil
	}

	return err
}

func (s *authService) validateLoginDummy(
	email, passwordRow string,
) error {
	var err error

	switch {
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

	default:
		err = nil
	}

	return err
}

func (s *authService) Login(
	ctx context.Context,
	body v1GenAPI.PostAuthLoginJSONRequestBody,
) (*v1GenAPI.Token, jwt.RefreshToken, error) {
	var refreshToken jwt.RefreshToken

	email := strings.TrimSpace(string(body.Email))
	passwordRow := body.Password

	if err := s.validateLogin(email, passwordRow); err != nil {
		return nil, refreshToken, err
	}

	user, err := s.repo.GetUserByLogin(ctx, email)
	if err != nil {
		return nil, refreshToken, err
	}

	if ok, err := password.Compare(passwordRow, user.PasswordHash); !ok {
		return nil, refreshToken, ErrIncorrectPassword
	} else if err != nil {
		return nil, refreshToken, err
	}

	return s.issueTokens(ctx, user.ID, user.Role, user.Permissions)
}

func (s *authService) LoginDummy(
	ctx context.Context,
	body v1GenAPI.PostAuthDummyLoginJSONRequestBody,
) (*v1GenAPI.Token, jwt.RefreshToken, error) {
	var refreshToken jwt.RefreshToken
	var id uuid.UUID

	role := body.Role

	if !role.Valid() {
		return nil, refreshToken, ErrIncorrectPassword //
	}

	if role == v1GenAPI.PostAuthDummyLoginJSONBodyRoleAdmin {
		id = dummyAdminID
	} else {
		id = dummyUserID
	}

	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, refreshToken, err
	}

	return s.issueTokens(ctx, user.ID, user.Role, user.Permissions)
}

func (s *authService) issueTokens(
	ctx context.Context,
	userID uuid.UUID,
	role string,
	permissions []string,
) (
	*v1GenAPI.Token,
	jwt.RefreshToken,
	error,
) {
	var response v1GenAPI.Token
	var refreshToken jwt.RefreshToken
	var err error

	response.Token, err = s.jwt.CreateAccessToken(userID.String(), role, permissions)
	if err != nil {
		return nil, refreshToken, err
	}

	refreshToken.Token, err = s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, refreshToken, err
	}

	hashRefreshToken := s.jwt.HashToken(refreshToken.Token)
	refreshToken.ExpiresAt = s.jwt.RefreshTokenExpiresAt()

	_, err = s.repo.CreateRefreshToken(
		ctx, db.CreateRefreshTokenParams{
			UserID:    userID,
			TokenHash: hashRefreshToken,
			ExpiresAt: refreshToken.ExpiresAt,
		},
	)
	if err != nil {
		return nil, refreshToken, err
	}

	return &response, refreshToken, nil
}
