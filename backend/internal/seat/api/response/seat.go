package response

type SeatResponse struct {
	ID       string  `json:"id"`
	SeatNo   string  `json:"seat_no"`
	ZoneID   string  `json:"zone_id"`
	ZoneName *string `json:"zone_name,omitempty"`
	Price    *int    `json:"price,omitempty"`
	Status   *string `json:"status,omitempty"`
}
