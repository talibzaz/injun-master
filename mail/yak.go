package mail

import (
	"github.com/domodwyer/mailyak"
	"github.com/spf13/viper"
	"sync"
	"net/smtp"
	"bytes"
	"html/template"
	"github.com/GeertJohan/go.rice"
	log "github.com/sirupsen/logrus"
	"github.com/graphicweave/injun/mail/model"
)

var (
	once, templateOnce sync.Once
	mail               *mailyak.MailYak

	abandonedTemplate            *template.Template
	approvalTemplate             *template.Template
	attendeeConfirmationTemplate *template.Template
	confirmationTemplate         *template.Template
	creationTemplate             *template.Template
	forgotPasswordTemplate       *template.Template
	publishedTemplate            *template.Template
	reminderTemplate             *template.Template
	ticketTemplate               *template.Template
	visitorTemplate              *template.Template
	welcomeTemplate              *template.Template
	sponsorEnquiryTemplate		*template.Template
	exhibitorEnquiryTemplate		*template.Template
	brochureRequestTemplate		*template.Template
)

type Mailer interface {
	SendWelcomeEmail(w model.Welcome) error
	SendCreationEmail(c model.Creation) error
	SendApprovalEmail(a model.Approval) error
	SendConfirmationEmail(c model.Confirmation) error
	SendAbandonedEmail(a model.Abandoned) error
	SendPublishedEmail(p model.Published) error
	SendReminderEmail(r model.Reminder) error
	SendVisitorEmail(v model.Visitor) error
	SendForgotPasswordEmail(f model.ForgotPassword) error
	SendTicketEmail(t model.Ticket) error
	SendAttendeeConfirmationEmail(a model.AttendeeConfirmation) error
	SendSponsorEnquiryEmail( f model.SponsorEnquiry) error
	SendExhibitorEnquiryEmail(q model.ExhibitorEnquiry) error
	SendBrochureRequestEmail (a model.BrochureRequest) error
}

type Mail struct {
	to, from, subject, body, attachmentName string
	attachment                              []byte
}

func NewMail(to, from string) *Mail {

	mail := &Mail{
		to:   to,
		from: from,
	}

	templateOnce.Do(func() {
		compileTemplates()
	})

	return mail
}

func compileTemplates() {

	box, err := rice.FindBox("mail-templates")

	if err != nil {
		log.Errorln("failed to compile email templates")
		log.Errorln("error:", err)
	}

	compiler := Compiler{}

	welcomeTemplate = compiler.MustCompile(box, "welcome.html", "welcome")
	attendeeConfirmationTemplate = compiler.MustCompile(box, "attendee-confirmation.html", "attendee-confirmation")
	creationTemplate = compiler.MustCompile(box, "creation.html", "creation")
	approvalTemplate = compiler.MustCompile(box, "approval.html", "approval")
	confirmationTemplate = compiler.MustCompile(box, "confirmation.html", "confirmation")
	ticketTemplate = compiler.MustCompile(box, "ticket.html", "ticket")
	forgotPasswordTemplate = compiler.MustCompile(box, "forgot-password.html", "forgotPassword")
	publishedTemplate = compiler.MustCompile(box, "published.html", "published")
	reminderTemplate = compiler.MustCompile(box, "reminder.html", "reminder")
	visitorTemplate = compiler.MustCompile(box, "visitor.html", "visitor")
	abandonedTemplate = compiler.MustCompile(box, "abandoned.html", "abandoned")
	sponsorEnquiryTemplate = compiler.MustCompile(box, "sponsor-enquiry.html", "sponsorEnquiry")
	exhibitorEnquiryTemplate = compiler.MustCompile(box, "exhibitor-enquiry.html", "exhibitorEnquiry")
	brochureRequestTemplate = compiler.MustCompile(box, "brochure.html", "brochureRequest")
}

func (m *Mail) SendWelcomeEmail(w model.Welcome) error {

	sWriter := ""
	writer := bytes.NewBufferString(sWriter)

	welcomeTemplate.Execute(writer, w)

	m.subject = "Welcome to Eventackle"
	m.body = writer.String()

	return m.send()
}

func (m *Mail) SendCreationEmail(c model.Creation) error {

	sWriter := ""
	writer := bytes.NewBufferString(sWriter)

	creationTemplate.Execute(writer, c)

	m.subject = "Your organizer profile has been created"
	m.body = writer.String()

	return m.send()
}

func (m *Mail) SendApprovalEmail(a model.Approval) error {
	sWriter := ""
	writer := bytes.NewBufferString(sWriter)

	approvalTemplate.Execute(writer, a)

	m.subject = "Your organizer profile has been approved"
	m.body = writer.String()

	return m.send()
}

func (m *Mail) SendConfirmationEmail(c model.Confirmation) error {

	sWriter := ""
	writer := bytes.NewBufferString(sWriter)

	confirmationTemplate.Execute(writer, c)

	m.subject = "Order confirmation"
	m.body = writer.String()

	return m.send()
}

func (m *Mail) SendForgotPasswordEmail(f model.ForgotPassword) error {

	sWriter := ""
	writer := bytes.NewBufferString(sWriter)

	forgotPasswordTemplate.Execute(writer, f)

	m.subject = "Forgot password"
	m.body = writer.String()

	return m.send()
}

func (m *Mail) SendAbandonedEmail(a model.Abandoned) error {
	sWriter := ""
	writer := bytes.NewBufferString(sWriter)

	abandonedTemplate.Execute(writer, a)

	m.subject = "Forgot something?"
	m.body = writer.String()

	return m.send()
}

func (m *Mail) SendTicketEmail(t model.Ticket) error {
	sWriter := ""
	writer := bytes.NewBufferString(sWriter)

	ticketTemplate.Execute(writer, t)

	m.subject = "Your ticket"
	m.body = writer.String()

	return m.send()
}

func (m *Mail) SendPublishedEmail(p model.Published) error {
	sWriter := ""
	writer := bytes.NewBufferString(sWriter)

	publishedTemplate.Execute(writer, p)

	m.subject = "Event Published"
	m.body = writer.String()

	return m.send()
}

func (m *Mail) SendReminderEmail(r model.Reminder) error {
	sWriter := ""
	writer := bytes.NewBufferString(sWriter)

	reminderTemplate.Execute(writer, r)

	m.subject = "Event reminder"
	m.body = writer.String()

	return m.send()
}

func (m *Mail) SendVisitorEmail(v model.Visitor) error {
	sWriter := ""
	writer := bytes.NewBufferString(sWriter)

	visitorTemplate.Execute(writer, v)

	m.subject = "Event reminder"
	m.body = writer.String()

	return m.send()
}

func (m *Mail) SendAttendeeConfirmationEmail(a model.AttendeeConfirmation) error {
	sWriter := ""
	writer := bytes.NewBufferString(sWriter)

	attendeeConfirmationTemplate.Execute(writer, a)

	m.subject = "Confirm your attendance"
	m.body = writer.String()

	return m.send()
}

func (m *Mail) SendSponsorEnquiryEmail( f model.SponsorEnquiry) error {

	sWriter := ""
	writer := bytes.NewBufferString(sWriter)

	sponsorEnquiryTemplate.Execute(writer, f)

	m.subject = f.EventName
	m.body = writer.String()

	return m.send()
}

func (m *Mail) SendExhibitorEnquiryEmail(q model.ExhibitorEnquiry) error {

	sWriter := ""
	writer := bytes.NewBufferString(sWriter)

	exhibitorEnquiryTemplate.Execute(writer, q)

	m.subject = q.EventName
	m.body = writer.String()

	return m.send()
}

func (m *Mail) SendBrochureRequestEmail (a model.BrochureRequest) error {

	sWriter := ""
	writer := bytes.NewBufferString(sWriter)

	brochureRequestTemplate.Execute(writer, a)

	m.subject = a.EventName
	m.body = writer.String()

	return m.send()
}

func (m *Mail) SetAttachment(attachment []byte) {
	m.attachment = attachment
}

func (m *Mail) send() error {

	once.Do(func() {
		mailHost := viper.GetString("MAIL_HOST")
		mail = mailyak.New(mailHost+":"+viper.GetString("MAIL_PORT"), smtp.PlainAuth("", viper.GetString("MAIL_USERNAME"), viper.GetString("MAIL_PASSWORD"), mailHost))
	})

	mail.To(m.to)
	mail.From(m.from)
	mail.Subject(m.subject)
	mail.HTML().Set(m.body)

	if m.attachment != nil && len(m.attachmentName) > 0 {
		mail.Attach(m.attachmentName, bytes.NewBuffer(m.attachment))
	}

	return mail.Send()
}
