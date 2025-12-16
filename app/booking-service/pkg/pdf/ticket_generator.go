package pdf

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"booking-service/internal/domain"

	"github.com/jung-kurt/gofpdf"
)

type TicketGenerator interface {
	GenerateTicket(ctx context.Context, ticket *domain.Ticket, details *domain.TicketWithDetails, qrCodeBase64 string) ([]byte, error)
}

type ticketGenerator struct {
	logoPath string
}

func NewTicketGenerator(logoPath string) TicketGenerator {
	return &ticketGenerator{
		logoPath: logoPath,
	}
}

func (g *ticketGenerator) GenerateTicket(ctx context.Context, ticket *domain.Ticket, details *domain.TicketWithDetails, qrCodeBase64 string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 24)

	pdf.SetFillColor(41, 128, 185)
	pdf.Rect(0, 0, 210, 40, "F")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 15, "", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 15, "MOVIE TICKET", "", 1, "C", false, 0, "")

	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Ticket #%d", ticket.ID))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(41, 128, 185)
	pdf.Cell(0, 10, details.Film.Name)
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 11)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(0, 7, details.Cinema.Name)
	pdf.Ln(5)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 6, details.Cinema.Address)
	pdf.Ln(10)

	pdf.SetFillColor(245, 245, 245)
	pdf.Rect(20, pdf.GetY(), 170, 45, "F")
	pdf.SetFont("Arial", "B", 11)
	
	yPos := pdf.GetY() + 7

	pdf.SetXY(25, yPos)
	pdf.Cell(0, 6, "DATE & TIME")
	pdf.SetFont("Arial", "", 11)
	pdf.SetXY(25, yPos+5)
	pdf.Cell(0, 6, details.Screening.StartAt.Format("Mon, Jan 02, 2006 - 3:04 PM"))

	pdf.SetFont("Arial", "B", 11)
	pdf.SetXY(25, yPos+15)
	pdf.Cell(0, 6, fmt.Sprintf("FORMAT: %s | LANGUAGE: %s", details.Screening.Format, details.Screening.Language))

	pdf.SetFont("Arial", "B", 11)
	pdf.SetXY(25, yPos+25)
	pdf.Cell(0, 6, "SEAT")
	pdf.SetFont("Arial", "", 11)
	pdf.SetXY(25, yPos+30)
	pdf.Cell(0, 6, fmt.Sprintf("Row %s, Seat %d", details.Seat.Row, details.Seat.Number))

	pdf.Ln(50)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, fmt.Sprintf("Price: $%.2f", ticket.Price))
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "SCAN QR CODE FOR CHECK-IN")
	pdf.Ln(5)

	if qrCodeBase64 != "" {
		qrImagePath := fmt.Sprintf("/tmp/qr_%d_%d.png", ticket.ID, time.Now().Unix())
		if err := saveBase64Image(qrCodeBase64, qrImagePath); err == nil {
			pdf.ImageOptions(qrImagePath, 70, pdf.GetY(), 70, 70, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
			pdf.Ln(75)
		}
	}

	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 5, fmt.Sprintf("QR Code: %s", ticket.QRCode))
	pdf.Ln(10)

	pdf.SetFillColor(255, 243, 205)
	pdf.Rect(20, pdf.GetY(), 170, 35, "F")
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(25, pdf.GetY()+5)
	pdf.Cell(0, 5, "IMPORTANT INFORMATION")
	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(25, pdf.GetY()+5)
	pdf.MultiCell(160, 4, "• Please arrive 30 minutes before the screening\n• This ticket is non-refundable\n• Present this QR code at the entrance for check-in\n• Outside food and beverages are not permitted", "", "", false)

	pdf.Ln(10)

	pdf.SetY(-30)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.CellFormat(0, 5, "Thank you for choosing our cinema!", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Issued: %s", time.Now().Format("2006-01-02 15:04:05")), "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

func saveBase64Image(base64Data, filepath string) error {
	//This is a placeholder - implement base64 to file conversion
	//returning nil to not break the flow
	return nil
}