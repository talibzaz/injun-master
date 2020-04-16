package mail

import (
	"golang.org/x/net/context"
	"github.com/graphicweave/injun/mail"
	"github.com/spf13/viper"
	"github.com/graphicweave/injun/mail/model"
	"time"
	"github.com/sirupsen/logrus"
	"crypto/sha256"
	"fmt"
	"github.com/graphicweave/injun/database"
	"github.com/graphicweave/injun/service"
)

/**
 * Created by  ♅ Salfi Farooq on 23/06/18.
 */

type MailService struct {
}

func (m MailService) ConfirmAttendee(ctx context.Context, attendee *mail.AttendeeConfirmationDetail) (*mail.AnyResponse, error) {

	mail_ := mail.NewMail(attendee.EmailId, viper.GetString("SENDER"))

	attendeeInfo := model.AttendeeConfirmation{
		Name:                  attendee.Name,
		EventName:             attendee.EventName,
		EventCoverImage:       attendee.EventCoverImage,
		EventDateTime:         attendee.EventDateTime,
		EventFullVenue:        attendee.EventFullVenue,
		EventOrganizerCompany: attendee.EventOrganizerCompany,
		EventURL:              attendee.EventURL,
		ConfirmationURL:       attendee.ConfirmationURL,
		Meta:                  getMeta(),
	}

	err := mail_.SendAttendeeConfirmationEmail(attendeeInfo)
	if err != nil {
		logrus.Error("couldn't send confirmation email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}
	logrus.Info("Confirmation Email Sent.")
	return &mail.AnyResponse{Status: "OK", Message: "Email Sent."}, nil
}

func (m MailService) ForgotPassword(ctx context.Context, req *mail.ForgotPasswordRequest) (*mail.ForgotPasswordResponse, error) {
	mail_ := mail.NewMail(req.EmailId, viper.GetString("SENDER"))

	resetToken := getToken(req.EmailId+req.Id)
	fmt.Println("resu",resetToken)
   err :=  mail_.SendForgotPasswordEmail(model.ForgotPassword{
   		Name: req.Name,
   		Email: req.EmailId,
    	PasswordURL: viper.GetString("PASSWORD_URL")+resetToken,
    	Meta :getMeta(),
  })

	if err != nil {
		logrus.Error("couldn't send reset password email :", err)
		return &mail.ForgotPasswordResponse{Status: "ERR", ResetToken: ""}, err
	}
	logrus.Info("Reset Password Email Sent.")
	return &mail.ForgotPasswordResponse{Status: "OK", ResetToken: resetToken}, nil

}

func (m MailService) WelcomeEmail(ctx context.Context, req *mail.Email) (*mail.AnyResponse, error) {
	mail_ := mail.NewMail(req.EmailId, viper.GetString("SENDER"))
	err := mail_.SendWelcomeEmail(model.Welcome{
		AccountSettings: viper.GetString("ACCOUNT_SETTINGS_URL"),
		OrganizerProfile: viper.GetString("ORGANIZER_PROFILE_URL"),
		EventSearch: viper.GetString("EVENTACKLE"),
		Meta: getMeta(),
	})
	if err != nil {
		logrus.Error("couldn't send welcome email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}
	logrus.Info("Welcome Email Sent.")
	return &mail.AnyResponse{Status: "OK", Message: "Email Sent."}, nil
}

func (m MailService) VisitorEmail(ctx context.Context, req *mail.VisitorDetail) (*mail.AnyResponse, error) {
	mail_ := mail.NewMail(req.EmailId, viper.GetString("SENDER"))
	err := mail_.SendVisitorEmail(model.Visitor{
		Name: req.Name,
		EventURL: req.EventURL,
		EventName: req.EventName,
		EventCoverImage: req.EventCoverImage,
		EventDateTime: req.EventDateTime,
		EventFullVenue: req.EventFullVenue,
		EventOrganizerCompany: req.EventOrganizerCompany,
		Coordinates: req.Coordinates,
		Meta :getMeta(),
	})
	if err != nil {
		logrus.Error("couldn't send visitor email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}
	logrus.Info("visitor Email Sent.")
	return &mail.AnyResponse{Status: "OK", Message: "Email Sent."}, nil
}

func (m MailService) ReminderEmail(ctx context.Context, req *mail.ReminderRequestDetail) (*mail.AnyResponse, error) {
	mail_ := mail.NewMail(req.EmailId, viper.GetString("SENDER"))
	err := mail_.SendReminderEmail(model.Reminder{
		Name: req.Name,
		EventURL: req.EventURL,
		EventCoverImage: fmt.Sprintf(viper.GetString("MINIO") + "/uploads/" + req.EventCoverImage),
		EventDateTime: req.EventDateTime,
		EventFullVenue: req.EventFullVenue,
		EventOrganizerCompany: req.EventOrganizerCompany,
		PrintTicketURL: req.PrintTicketURL,
		DaysToStart: req.DaysToStart,
		Coordinates: req.Coordinates,
		Meta :getMeta(),
	})
	if err != nil {
		logrus.Error("couldn't send reminder email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}
	logrus.Info("reminder Email Sent.")
	return &mail.AnyResponse{Status: "OK", Message: "Email Sent."}, nil
}

func (m MailService) PublishedEmail(ctx context.Context, req *mail.PublishedRequestDetail) (*mail.AnyResponse, error) {
	EmailId,err := database.GetOrganizerEmail(req.OrganizerId)
	if err != nil {
		logrus.Error("couldn't send brochure request email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}
	mail_ := mail.NewMail(EmailId, viper.GetString("SENDER"))
	err = mail_.SendPublishedEmail(model.Published{
		Name: req.Name ,
		EventURL: req.EventURL   ,
		EventName: req.EventName ,
		EventCoverImage: fmt.Sprintf(viper.GetString("MINIO") + "/uploads/" + req.EventCoverImage ) ,
		EventDateTime: req.EventDateTime ,
		EventFullVenue: req.EventFullVenue ,
		EventPromotionGuide: viper.GetString("EVENT_PROMOTION_GUIDE") ,
		Meta: getMeta(),
	})
	if err != nil {
		logrus.Error("couldn't send published email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}
	logrus.Info("published Email Sent.")
	return &mail.AnyResponse{Status: "OK", Message: "Email Sent."}, nil
}

func (m MailService) TicketEmail(ctx context.Context, req *mail.TicketRequestDetail) (*mail.AnyResponse, error) {
	logrus.Infof("got request to send ticket with userId = %v and ticketId = %v", req.UserId, req.TicketId)

	ctx = context.Background()
	db, err := database.NewArangoDB(ctx)
	if err != nil {
		logrus.Error("couldn't send ticket email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}
	Details, err := db.GetEventDetails(req.TicketId, req.UserId)
	if err != nil {
		logrus.Error("couldn't send ticket email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}
	var userDetails = Details["user"].([]interface{})
	var eventDetails = Details["event"].([]interface{})

	mail_ := mail.NewMail(userDetails[0].(map[string]interface{})["email"].(string), viper.GetString("SENDER"))
	err = mail_.SendTicketEmail(model.Ticket{
		Name: userDetails[0].(map[string]interface{})["name"].(string) ,
		EventURL: fmt.Sprint(viper.GetString("EVENTACKLE") + "/search/view?=" + eventDetails[0].(map[string]interface{})["eventId"].(string)),
		EventName: eventDetails[0].(map[string]interface{})["eventName"].(string),
		EventCoverImage: fmt.Sprintf(viper.GetString("MINIO") + "/uploads/" + eventDetails[0].(map[string]interface{})["eventCoverImage"].(string)),
		EventDateTime: eventDetails[0].(map[string]interface{})["eventDateTime"].(string),
		EventFullVenue: eventDetails[0].(map[string]interface{})["eventFullVenue"].(string),
		EventOrganizerCompany: eventDetails[0].(map[string]interface{})["eventOrganizerCompany"].(string),
		PrintTicketURL: fmt.Sprint(viper.GetString("EVENTACKLE") + "/print-tickets/" + req.UserId + "/" + req.TicketId),
		Meta: getMeta(),
	})
	if err != nil {
		logrus.Error("couldn't send ticket email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}
	logrus.Info("ticket Email Sent.")
	return &mail.AnyResponse{Status: "OK", Message: "Email Sent."}, nil
}

func (m MailService) AbandonedEmail(ctx context.Context, req *mail.AbandonedRequestDetail) (*mail.AnyResponse, error) {
	mail_ := mail.NewMail(req.EmailId, viper.GetString("SENDER"))
	err := mail_.SendAbandonedEmail(model.Abandoned{
		Name: req.Name,
		EventName: req.EventName,
		EventDateTime: req.EventDateTime,
		EventFullVenue: req.EventFullVenue,
		EventOrganizerCompany: req.EventOrganizerCompany,
		EventCoverImage: fmt.Sprintf(viper.GetString("MINIO") + "/uploads/" + req.EventCoverImage),
		CheckoutURL: req.CheckoutURL ,
		TicketAmount: req.TicketAmount,
		Coordinates: req.Coordinates,
		Meta: getMeta(),
	})
	if err != nil {
		logrus.Error("couldn't send abandoned email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}
	logrus.Info("abandoned Email Sent.")
	return &mail.AnyResponse{Status: "OK", Message: "Email Sent."}, nil
}

func (m MailService) ConfirmationEmail(ctx context.Context, req *mail.ConfirmationRequestDetail) (*mail.AnyResponse, error) {

	ctx = context.Background()
	db, err := database.NewArangoDB(ctx)
	if err != nil {
		logrus.Error("couldn't send ticket email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}
	Details, err := db.GetEventDetails(req.TicketNumber, "")
	if err != nil {
		logrus.Error("couldn't send ticket email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}

	var eventDetails = Details["event"].([]interface{})

	mail_ := mail.NewMail(req.EmailId, viper.GetString("SENDER"))
	err = mail_.SendConfirmationEmail(model.Confirmation{
		Name: req.Name ,
		EventURL: fmt.Sprint(viper.GetString("EVENTACKLE") + "/search/view?id=" + eventDetails[0].(map[string]interface{})["eventId"].(string)),
		EventName: eventDetails[0].(map[string]interface{})["eventName"].(string) ,
		EventDateTime: eventDetails[0].(map[string]interface{})["eventDateTime"].(string),
		EventFullVenue: eventDetails[0].(map[string]interface{})["eventFullVenue"].(string),
		EventOrganizerCompany: eventDetails[0].(map[string]interface{})["eventOrganizerCompany"].(string),
		TicketNumber: req.TicketNumber,
		EventCoverImage: fmt.Sprintf(viper.GetString("MINIO") + "/uploads/" + eventDetails[0].(map[string]interface{})["eventCoverImage"].(string)),
		TicketPurchaseDate: req.TicketPurchaseDate,
		TicketAmount: req.TicketAmount,
		Coordinates: eventDetails[0].(map[string]interface{})["coordinates"],
		PrintTicketURL: fmt.Sprint(viper.GetString("EVENTACKLE") + "/print-tickets/" + req.TicketNumber),
		Meta: getMeta(),
	})
	if err != nil {
		logrus.Error("couldn't send confirmation email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}
	logrus.Info("confirmation Email Sent.")
	return &mail.AnyResponse{Status: "OK", Message: "Email Sent."}, nil
}

func (m MailService) ApprovalEmail(ctx context.Context, req *mail.CreationRequestDetail) (*mail.AnyResponse, error) {
	mail_ := mail.NewMail(req.EmailId, viper.GetString("SENDER"))
	err := mail_.SendApprovalEmail(model.Approval{
		Name: req.OrganizerName,
		CreationGuide: viper.GetString("EVENT_CREATION_GUIDE"),
		Meta: getMeta(),
	})
	if err != nil {
		logrus.Error("couldn't send approval email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}
	logrus.Info("approval Email Sent.")
	return &mail.AnyResponse{Status: "OK", Message: "Email Sent."}, nil
}

func (m MailService) CreationEmail(ctx context.Context, req *mail.CreationRequestDetail) (*mail.AnyResponse, error) {
	mail_ := mail.NewMail(req.EmailId, viper.GetString("SENDER"))
	err := mail_.SendCreationEmail(model.Creation{
		OrganizerName: req.OrganizerName,
		EventCreationGuide: viper.GetString("EVENT_CREATION_GUIDE"),
		Meta: getMeta(),
	})
	if err != nil {
		logrus.Error("couldn't send creation email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}
	logrus.Info("creation Email Sent.")
	return &mail.AnyResponse{Status: "OK", Message: "Email Sent."}, nil
}

func (m MailService) SponsorEnquiryEmail (ctx context.Context, req *mail.EnquiryRequest) (*mail.AnyResponse, error) {

	logrus.Info("got sponsor enquiry request: ", req)

	organizerEmail,err := database.GetOrganizerEmail(req.OrganizerId)

	if err != nil {
		logrus.Error("couldn't send sponsor enquiry email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}

	mail_ := mail.NewMail(organizerEmail, viper.GetString("SENDER"))
	err = mail_.SendSponsorEnquiryEmail(model.SponsorEnquiry{
		Name: req.Name,
		Phone: req.Phone,
		Company: req.Company,
		JobTitle: req.JobTitle,
		Email: req.Email,
		CompanyWebsite: req.CompanyWebsite,
		Comments: req.Comments,
		EventId: req.EventId,
		EventName: req.EventName,
		Meta: getMeta(),
	})

	if err != nil {
		logrus.Error("couldn't send sponsor enquiry email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}


	arangoDb, err := service.NewArangoDB(context.Background())
	if err != nil {
		logrus.Error("failed to get arangodb client")
		logrus.Errorln("Error ", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}

	req.EnquiryType = "sponsor"

	err = arangoDb.InsertEnquiryDetails(req)
	if err != nil {
		logrus.Error("couldn't save sponsor enquiry request detail :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}

	logrus.Info("sponsor enquiry Email Sent to: ",organizerEmail)
	return &mail.AnyResponse{Status: "OK", Message: "Email Sent."}, nil
}

func (m MailService) ExhibitorEnquiryEmail (ctx context.Context, req *mail.EnquiryRequest) (*mail.AnyResponse, error) {
	logrus.Info("got exhibitor enquiry request: ", req)
	organizerEmail,err := database.GetOrganizerEmail(req.OrganizerId)
	if err != nil {
		logrus.Error("couldn't send exhibitor enquiry email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}
	mail_ := mail.NewMail( organizerEmail, viper.GetString("SENDER"))
	err = mail_.SendExhibitorEnquiryEmail(model.ExhibitorEnquiry{
		Name:req.Name,
		Phone: req.Phone,
		Company: req.Company	,
		JobTitle: req.JobTitle,
		Email: req.Email,
		CompanyWebsite: req.CompanyWebsite,
		Comments: req.Comments,
		EventId: req.EventId,
		EventName: req.EventName,
		Meta: getMeta(),

	})
	if err != nil {
		logrus.Error("couldn't send exhibitor enquiry email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}

	arangoDb, err := service.NewArangoDB(context.Background())
	if err != nil {
		logrus.Error("failed to get arangodb client")
		logrus.Errorln("Error ", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}

	req.EnquiryType = "exhibitor"

	err = arangoDb.InsertEnquiryDetails(req)
	if err != nil {
		logrus.Error("couldn't save exhibitor enquiry request detail :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}


	logrus.Info("exhibitor enquiry Email Sent to: ",organizerEmail)
	return &mail.AnyResponse{Status: "OK", Message: "Email Sent."}, nil
}

func (m MailService) BrochureRequestEmail(ctx context.Context, req *mail.BrochureRequest) (*mail.AnyResponse, error) {
	logrus.Info("got sponsor request: ", req)
	organizerEmail,err := database.GetOrganizerEmail(req.OrganizerId)
	if err != nil {
		logrus.Error("couldn't send brochure request email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}
	mail_ := mail.NewMail(organizerEmail, viper.GetString("SENDER"))
	err = mail_.SendBrochureRequestEmail(model.BrochureRequest{
		Name: req.Name,
		Phone: req.Phone ,
		Company: req.Company	,
		Email: req.Email,
		CompanyWebsite: req.CompanyWebsite,
		Comments: req.Comments,
		Address1: req.Address1,
		Address2: req.Address2,
		City: req.City,
		Country: req.Country,
		EventId: req.EventId,
		EventName: req.EventName,
		Meta: getMeta(),
	})
	if err != nil {
		logrus.Error("couldn't send brochure request email :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}

	arangoDb, err := service.NewArangoDB(context.Background())
	if err != nil {
		logrus.Error("failed to get arangodb client")
		logrus.Errorln("Error ", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}

	if err := arangoDb.BrochureRequest(req); err != nil {
		logrus.Error("couldn't save brochure request detail :", err)
		return &mail.AnyResponse{Status: "ERR", Message: "Internal Server Error"}, err
	}

	logrus.Info("brochure request Email Sent to: ",organizerEmail)
	return &mail.AnyResponse{Status: "OK", Message: "Email Sent."}, nil
}


func getMeta() (model.Meta) {
	return model.Meta{
		Address1:        viper.GetString("ADDRESS1"),
		Address2:        viper.GetString("ADDRESS2"),
		Address3:        viper.GetString("ADDRESS3"),
		Phone:           viper.GetString("PHONE"),
		EventackleEmail: viper.GetString("EMAIL"),
		Year:            time.Now().Year(),
	}
}

func getToken(message string)string{
	h := sha256.New()
	h.Write([]byte(message))
	token := fmt.Sprintf("%x",h.Sum(nil))
	return token
}
