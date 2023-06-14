package user

// ticketDate
// departureOn
// travelTime
// seat
// lokasiAwal
// lokasiAkhir
// qrcode
type UserTicket struct {
	TicketDate  string `json:"ticketDate"`
	DepartureOn string `json:"departureOn"`
	TravelTime  int    `json:"travelTime"`
	Seat        string `json:"seat"`
	From        int    `json:"lokasiAwal"`
	To          int    `json:"lokasiAkhir"`
	Qrcode      string `json:"qrcode"`
}

type CreateUserTicket struct {
	TicketDate  string `json:"ticketDate" validate:"required"`
	DepartureOn string `json:"departureOn" validate:"required"`
	TravelTime  int    `json:"travelTime" validate:"required"`
	Seat        string `json:"seat" validate:"required"`
	From        int    `json:"lokasiAwal" validate:"required"`
	To          int    `json:"lokasiAkhir" validate:"required"`
}
