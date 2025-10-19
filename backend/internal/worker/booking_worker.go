package worker

import (
	"fmt"
	"time"

	"github.com/baimhons/stadiumhub/internal/booking"
)

type BookingWorker struct {
	bookingService booking.BookingService
}

func NewBookingWorker(bookingService booking.BookingService) *BookingWorker {
	return &BookingWorker{bookingService: bookingService}
}

func (w *BookingWorker) Start() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("[Worker] Checking for expired bookings...")
		expiredCount, err := w.bookingService.CancelExpiredBookings(30 * time.Minute)
		if err != nil {
			fmt.Println("[Worker] Error:", err)
		} else {
			fmt.Printf("[Worker] Cancelled %d expired bookings\n", expiredCount)
		}
	}
}
