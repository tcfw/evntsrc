package emails

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
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
		log.Printf("Failed to send: %v ~ %e", msg, err)
	} else {
		if resp.StatusCode == 400 {
			log.Printf("Failed to send email: %v :: %v", msg, resp)
			return
		}
		log.Printf("Sent email %d (%s)\n", resp.StatusCode, resp.Headers["X-Message-Id"])
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

//Send triggers an email with templating
//Plain text is not modified by the template
func (s *Server) Send(ctx context.Context, request *pb.Email) (*pb.EmailResponse, error) {
	tmpl, err := template.New("email").Parse(globalTemplate)
	if err != nil {
		return nil, fmt.Errorf("Failed to init global template: %s", err)
	}

	buf := &bytes.Buffer{}
	if err = tmpl.Execute(buf, template.HTML(request.Html)); err != nil {
		return nil, fmt.Errorf("Failed to execute global template: %s", err)
	}
	request.Html = buf.String()

	return s.SendRaw(ctx, request)
}

//SendRaw triggers a raw email without templating
func (s *Server) SendRaw(ctx context.Context, request *pb.Email) (*pb.EmailResponse, error) {
	from := mail.NewEmail("EvntSrc.io", "no-reply@evntsrc.io")

	if len(request.To) == 0 {
		return nil, fmt.Errorf("Can't send to no one")
	}

	if request.PlainText == "" {
		request.PlainText = stripHTML(request.Html)
	}

	for _, recipt := range request.To {
		to := mail.NewEmail(recipt.Name, recipt.Email)
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
