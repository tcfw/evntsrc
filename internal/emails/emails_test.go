package emails

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/tcfw/go-queue"

	pb "github.com/tcfw/evntsrc/internal/emails/protos"
)

type mockEmailProcessor struct {
	msgs []*mail.SGMailV3
}

func (p *mockEmailProcessor) Handle(job interface{}) {
	msg := job.(*mail.SGMailV3)
	p.msgs = append(p.msgs, msg)
}

func TestBasicSend(t *testing.T) {
	//Spin up worker
	processor := &mockEmailProcessor{}
	worker = queue.NewDispatcher(processor)
	worker.MaxWorkers = 1
	worker.Run()

	msg := &pb.Email{
		Subject:   "Test email",
		To:        []string{"tester@evntsrc.io"},
		Html:      "<b>This is a test</b>",
		PlainText: "This is a test",
	}

	s := NewServer()
	_, err := s.Send(context.Background(), msg)

	if assert.NoError(t, err) {
		//Wait for worker to process
		for {
			if len(processor.msgs) > 0 {
				break
			}
		}

		assert.Len(t, processor.msgs, 1)
		assert.Equal(t, msg.Subject, processor.msgs[0].Subject)
	}
}

func TestBasicAttachment(t *testing.T) {
	//Spin up worker
	processor := &mockEmailProcessor{}
	worker = queue.NewDispatcher(processor)
	worker.MaxWorkers = 1
	worker.Run()

	msg := &pb.Email{
		Subject:   "Test email",
		To:        []string{"tester@evntsrc.io"},
		Html:      "<b>This is a test</b>",
		PlainText: "This is a test",
		Attachments: []*pb.Attachment{
			{
				Filename: "test.txt",
				Type: &pb.Attachment_Data{
					Data: []byte("this is a test"),
				},
			},
		},
	}

	s := NewServer()
	_, err := s.Send(context.Background(), msg)

	if assert.NoError(t, err) {
		//Wait for worker to process
		for {
			if len(processor.msgs) > 0 {
				break
			}
		}

		assert.Len(t, processor.msgs, 1)
		assert.Len(t, processor.msgs[0].Attachments, 1)
	}
}

func TestRemoteAttachment(t *testing.T) {
	//Test will require internet to complete successfully

	//Spin up worker
	processor := &mockEmailProcessor{}
	worker = queue.NewDispatcher(processor)
	worker.MaxWorkers = 1
	worker.Run()

	msg := &pb.Email{
		Subject:   "Test email",
		To:        []string{"tester@evntsrc.io"},
		Html:      "<b>This is a test</b>",
		PlainText: "This is a test",
		Attachments: []*pb.Attachment{
			{
				Filename: "test.txt",
				Type: &pb.Attachment_Uri{
					Uri: "https://api.staging.evntsrc.io/",
				},
			},
		},
	}

	s := NewServer()
	_, err := s.Send(context.Background(), msg)

	if assert.NoError(t, err) {
		//Wait for worker to process
		for {
			if len(processor.msgs) > 0 {
				break
			}
		}

		assert.Len(t, processor.msgs, 1)
		assert.Len(t, processor.msgs[0].Attachments, 1)
	}
}
