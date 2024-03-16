package process

import (
	"bytes"
	"text/template"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gojek/work"
	"gopkg.in/gomail.v2"
)

type SendMail struct {
	jobName string
	pool    *work.WorkerPool
}

func NewSendMailProcess(jobName string, pool *work.WorkerPool) Processor {
	return &SendMail{
		jobName: jobName,
		pool:    pool,
	}
}

func (s *SendMail) log(job *work.Job, next work.NextMiddlewareFunc) error {
	log.Info("Starting job: ", job.Name)
	return next()
}

func (s *SendMail) execute(job *work.Job) error {
	// topic := job.ArgString("topic")
	// message := job.Args["message"]
	receiver := job.ArgString("receiver")
	subject := job.ArgString("subject")
	payload := job.Args["payload"]

	if err := job.ArgError(); err != nil {
		return err
	}
	// Define the email template
	// TODO: change this to a template dynamically
	const emailTemplate = `
	<html>
	<body>
			<h1>Hello, {{.Name}}!</h1>
			<p>This is a templated email sent from Go.</p>
	</body>
	</html>
`

	// Create a new template and parse the email template
	tmpl := template.Must(template.New("email").Parse(emailTemplate))

	// Execute the template with the provided data
	var emailBody bytes.Buffer
	err := tmpl.Execute(&emailBody, payload)
	if err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", "app.notif@gmf-aeroasia.co.id")
	mail.SetHeader("To", receiver)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", emailBody.String())

	// Set up the email client
	dialer := gomail.NewDialer("smtp.office365.com", 587, "app.notif@gmf-aeroasia.co.id", "@Ap.notif")

	// Send the email
	if err := dialer.DialAndSend(mail); err != nil {
		log.Error("Mail send error: ", err)
		return err
	}

	log.Infof("Job completed successfully for %s:%s", job.Name, job.ID)
	return nil
}

// GetWorkerPool implements Processor.
func (s *SendMail) GetWorkerPool() *work.WorkerPool {
	return s.pool
}

func (s *SendMail) Do() {
	s.pool.Middleware(s.log)
	s.pool.Job(s.jobName, s.execute)
}

func (s *SendMail) ProcessName() string {
	return s.jobName
}
