package verifysvc

import (
	"context"
	"fmt"
	"time"

	"github.com/resend/resend-go/v2"
)

type SendParams struct {
	Email  string
	UserId string
}

func (s *Service) Send(ctx context.Context, sendparams SendParams) error {
	apiKey := s.cfg.ResendApiKey
	client := resend.NewClient(apiKey)

	verify, err := s.db.User.Verify.Create(ctx, sendparams.UserId)
	if err != nil {
		return err
	}
	verifyid := verify.ID

	verifyemail := s.cfg.Url + "/auth/verification/" + verifyid

	params := &resend.SendEmailRequest{
		From:    "no-reply@dooleyonline.net",
		To:      []string{sendparams.Email},
		Subject: "Verify your email",
		Html:    "<p>Please verify your email by clicking below:</p><a href=\"" + verifyemail + "\">Verify Email</a>",
	}
	_, err = client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) IsValid(ctx context.Context, id string) (bool, error) {
	verify, err := s.db.User.Verify.Get(ctx, id)
	if err != nil {
		return false, err
	}

	if verify.ExpiredAt.Before(time.Now()) {
		return false, nil
	}
	return true, nil
}

func (s *Service) VerifyUserEmail(ctx context.Context, id string) error {
	verify, err := s.IsValid(ctx, id)
	if err != nil {
		return err
	}
	if !verify {
		return fmt.Errorf("verification email is not valid")
	}
	v, err := s.db.User.Verify.Get(ctx, id)
	if err != nil {
		return err
	}

	err = s.db.User.User.Verify(ctx, v.UserID)
	if err != nil {
		return err
	}
	return nil

}
