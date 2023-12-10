package emailclient

import (
	"context"

	pb "github.com/a-dev-mobile/kidneysmart-auth/proto"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

type EmailClient struct {
	client pb.EmailSenderApiClient
	Logger *slog.Logger
}

func NewEmailClient(conn *grpc.ClientConn, lg *slog.Logger) *EmailClient {
	return &EmailClient{
		client: pb.NewEmailSenderApiClient(conn),
		Logger: lg,
	}
}

func (s *EmailClient) SendEmail(recipient, subject, fromName, fromEmail, body string) error {
	req := &pb.EmailSenderRequest{
		RecipientEmail: recipient,
		Subject:        subject,
		FromName:       fromName,
		FromEmail:      fromEmail,
		Body:           body,
	}

	_, err := s.client.SendEmail(context.Background(), req)
	if err != nil {
		s.Logger.Warn("could not send email:","error", err)
		return err
	}

	return nil
}
