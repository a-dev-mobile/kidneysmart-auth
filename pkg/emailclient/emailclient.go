package emailclient

import (
	"context"
	"log"

	pb "github.com/a-dev-mobile/kidneysmart-auth/proto"
	"google.golang.org/grpc"
)

type EmailClient struct {
	client pb.EmailSenderApiClient
}

func NewEmailClient(conn *grpc.ClientConn) *EmailClient {
	return &EmailClient{
		client: pb.NewEmailSenderApiClient(conn),
	}
}

func (e *EmailClient) SendEmail(recipient, subject, fromName, fromEmail, body string) error {
	req := &pb.EmailSenderRequest{
		RecipientEmail: recipient,
		Subject:        subject,
		FromName:       fromName,
		FromEmail:      fromEmail,
		Body:           body,
	}

	_, err := e.client.SendEmail(context.Background(), req)
	if err != nil {
		log.Printf("could not send email: %v", err)
		return err
	}

	return nil
}
