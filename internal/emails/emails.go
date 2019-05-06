package emails

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	pb "github.com/tcfw/evntsrc/internal/emails/protos"
	"github.com/tcfw/go-queue"
)

//Server core struct
type Server struct {
	mu sync.Mutex
}

type emailProcessor struct{}

func (p *emailProcessor) Handle(job interface{}) {
	msg := job.(*mail.SGMailV3)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	resp, err := client.Send(msg)
	if err != nil {
		fmt.Printf("Failed to send: %v ~ %e", msg, err)
	} else {
		fmt.Printf("Sent email (%s)\n", resp.Headers["X-Message-Id"])
	}
}

var worker *queue.Dispatcher

func startWorker() {
	worker = queue.NewDispatcher(&emailProcessor{})
	worker.Run()
}

//NewServer creates a new struct to interface the streams server
func NewServer() *Server {
	return &Server{}
}

//Send turns a call into an email
func (s *Server) Send(ctx context.Context, request *pb.Email) (*pb.EmailResponse, error) {

	from := mail.NewEmail("EvntSrc.io", "no-reply@evntsrc.io")

	if len(request.To) == 0 {
		return nil, fmt.Errorf("Can't send to no one")
	}

	for _, email := range request.To {
		to := mail.NewEmail("", email)
		message := mail.NewSingleEmail(from, request.Subject, to, request.PlainText, request.Html)

		addHeaders(message, request.Headers)

		if _, err := addAttachments(message, request.Attachments); err != nil {
			return nil, err
		}

		worker.Queue(message)
	}

	return &pb.EmailResponse{}, nil
}

func addHeaders(message *mail.SGMailV3, headers map[string]string) *mail.SGMailV3 {
	for k, v := range headers {
		message.SetHeader(k, v)
	}

	return message
}

func addAttachments(message *mail.SGMailV3, attachments []*pb.Attachment) (*mail.SGMailV3, error) {
	for _, attachment := range attachments {
		att, err := fetchAttachment(attachment)
		if err != nil {
			return nil, err
		}
		message.AddAttachment(att)
	}

	return message, nil
}

func fetchAttachment(attachment *pb.Attachment) (*mail.Attachment, error) {

	switch mtype := attachment.GetType().(type) {
	case *pb.Attachment_Uri:
		content, err := getURIContent(attachment.GetUri())
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch attachment: %s", err.Error())
		}

		return &mail.Attachment{
			Content:     base64.StdEncoding.EncodeToString(*content),
			Filename:    attachment.Filename,
			Type:        mimeFromFilename(attachment.Filename),
			Disposition: "attachment",
		}, nil

	case *pb.Attachment_Data:
		return &mail.Attachment{
			Content:     base64.StdEncoding.EncodeToString(attachment.GetData()),
			Filename:    attachment.Filename,
			Type:        mimeFromFilename(attachment.Filename),
			Disposition: "attachment",
		}, nil

	default:
		return nil, fmt.Errorf("Unknown attachment type %T", mtype)
	}

}

func mimeFromFilename(filename string) string {
	ext := filepath.Ext(filename)
	if mime := mime.TypeByExtension(ext); mime != "" {
		return mime
	}

	return "application/octet-stream"
}

func getURIContent(uri string) (*[]byte, error) {
	resp, err := http.Get(uri)
	// handle the error if there is one
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &content, nil
}
