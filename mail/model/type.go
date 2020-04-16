package model

type Meta struct {
	Address1        string
	Address2        string
	Address3        string
	EventackleEmail string
	Phone           string
	Year int
}

type Welcome struct {
	AccountSettings  string
	OrganizerProfile string
	EventSearch      string
	Meta
}

type Creation struct {
	OrganizerName      string
	EventCreationGuide string
	Meta
}

type Confirmation struct {
	Name                  string
	EventURL 			  string
	EventName             string
	EventDateTime         string
	EventFullVenue        string
	EventOrganizerCompany string
	TicketNumber          string
	EventCoverImage       string
	TicketPurchaseDate 		string
	TicketAmount	 		int32
	PrintTicketURL 		string
	Coordinates        interface{}
	Meta
}

type Approval struct {
	Name 	string
	CreationGuide string
	Meta
}

type ForgotPassword struct {
	Name string
	Email string
	PasswordURL string
	Meta
}

type Abandoned struct {
	Name                  string
	EventName             string
	EventDateTime         string
	EventFullVenue        string
	EventOrganizerCompany string
	EventCoverImage       string
	CheckoutURL           string
	TicketAmount 			int32
	Coordinates        []float64
	Meta
}

type Ticket struct {
	Name                  string
	EventURL              string
	EventName             string
	EventCoverImage       string
	EventDateTime         string
	EventFullVenue        string
	EventOrganizerCompany string
	PrintTicketURL        string
	Coordinates []float64
	Meta
}

type Published struct {
	Name                  string
	EventURL              string
	EventName             string
	EventCoverImage       string
	EventDateTime         string
	EventFullVenue        string
	EventOrganizerCompany string
	EventPromotionGuide   string
	Meta
}

type Reminder struct {
	Name                  string
	EventURL              string
	EventName             string
	EventCoverImage       string
	EventDateTime         string
	EventFullVenue        string
	EventOrganizerCompany string
	PrintTicketURL        string
	DaysToStart int32
	Coordinates []float64
	Meta
}

type AttendeeConfirmation struct {
	Name                  string
	EventURL              string
	EventName             string
	EventCoverImage       string
	EventDateTime         string
	EventFullVenue        string
	EventOrganizerCompany string
	ConfirmationURL       string
	Meta
}

type Visitor struct {
	Name                  string
	EventURL              string
	EventName             string
	EventCoverImage       string
	EventDateTime         string
	EventFullVenue        string
	EventOrganizerCompany string
	Coordinates           []float64
	Meta
}

type SponsorEnquiry struct {
	Name		string
	Phone 		string
	Company		string
	JobTitle	string
	Email		string
	CompanyWebsite	string
	Comments	string
	EventId 	string
	EventName 	string
	Meta
}

type ExhibitorEnquiry struct {
	Name		string
	Phone 		string
	Company		string
	JobTitle	string
	Email		string
	CompanyWebsite	string
	Comments	string
	EventId 	string
	EventName 	string
	Meta
}

type BrochureRequest struct {
	Name		string
	Phone 		string
	Company		string
	Email		string
	CompanyWebsite	string
	Comments	string
	Address1	string
	Address2	string
	City		string
	Country		string
	EventId 	string
	EventName 	string
	Meta
}
