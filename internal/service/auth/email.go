package authsvc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dooleyonline/backend/internal/model"
	"github.com/resend/resend-go/v2"
)

type SendParams struct {
	Email  string `json:"email"`
	UserID string `json:"userID"`
}

func (s *Service) CreateVerification(ctx context.Context, sendparams SendParams) (string, error) {
	apiKey := s.cfg.ResendApiKey
	client := resend.NewClient(apiKey)

	ver, err := s.db.User.Verify.Create(ctx, sendparams.UserID)
	if err != nil {
		return "", err
	}

	link := fmt.Sprintf("%s%s?id=%s&token=%s", s.cfg.FrontendUrl, "/auth/verify", ver.ID, ver.Token)

	params := &resend.SendEmailRequest{
		From:    "hello@dooleyonline.net",
		To:      []string{sendparams.Email},
		Subject: "Verify your email",
		Html: fmt.Sprintf(`
			<p>Please verify your email by clicking below:</p>
			<a href="%s">Verify Email</a>
		`, link),
	}
	_, err = client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to send verification email: %w", err)
	}
	return ver.ID, nil
}

func (s *Service) GetVerification(ctx context.Context, verficicationID string) (*model.Verification, error) {
	verification, err := s.db.User.Verify.Get(ctx, verficicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get verification: %w", err)
	}
	user, err := s.db.User.User.GetByID(ctx, verification.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.Verified {
		return nil, errors.New("user is already verified")
	}

	if verification.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("email verification has expired")
	}

	return &verification, nil
}

func (s *Service) VerifyUserEmail(ctx context.Context, verificationID string, token string) error {
	verification, err := s.GetVerification(ctx, verificationID)
	if err != nil {
		return fmt.Errorf("failed to get verification: %w", err)
	}

	if token != verification.Token {
		return fmt.Errorf("invalid token")
	}

	err = s.db.User.User.Verify(ctx, verification.UserID)
	if err != nil {
		return fmt.Errorf("failed to set user verification to true: %w", err)
	}

	return nil
}
