package authsvc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/resend/resend-go/v2"
)

type SendParams struct {
	Email  string
	UserId string
}

func (s *Service) CreateVerification(ctx context.Context, sendparams SendParams) (string, error) {
	apiKey := s.cfg.ResendApiKey
	client := resend.NewClient(apiKey)

	ver, err := s.db.User.Verify.Create(ctx, sendparams.UserId)
	if err != nil {
		return "", err
	}

	link := s.cfg.FrontendUrl + "/auth/verification/" + ver.ID

	params := &resend.SendEmailRequest{
		From:    "no-reply@dooleyonline.net",
		To:      []string{sendparams.Email},
		Subject: "Verify your email",
		Html:    "<p>Please verify your email by clicking below:</p><a href=\"" + link + "\">Verify Email</a>",
	}
	_, err = client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to create verification: %w", err)
	}
	return ver.ID, nil
}

func (s *Service) IsValid(ctx context.Context, id string) (string, error) {
	verify, err := s.db.User.Verify.Get(ctx, id)
	if err != nil {
		return "", err
	}

	if verify.ExpiredAt.Before(time.Now()) {
		return "", errors.New("email verification has expired")
	}
	return verify.UserID, nil
}

func (s *Service) VerifyUserEmail(ctx context.Context, id string) error {
	id, err := s.IsValid(ctx, id)
	if err != nil {
		return fmt.Errorf("verification email is not valid")
	}

	err = s.db.User.User.Verify(ctx, id)
	if err != nil {
		return fmt.Errorf("couldn't set user verification to true: %w", err)
	}
	return nil

}
