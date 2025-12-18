package qrCode

import (
	"booking/internal/app/models"
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(data string) ([]byte, error) {
	return qrcode.Encode(data, qrcode.Medium, 256)
}

func GenerateTicketPDF(
	ticket models.Ticket,
	screening models.Screening,
	seat models.SeatScreening,
	film models.Film,
	hall models.Hall,
	cinema models.Cinema,
	qr []byte,
) ([]byte, error) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	pdf.AddUTF8Font("DejaVu", "", "assets/fonts/DejaVuSans.ttf")
	pdf.AddUTF8Font("DejaVu", "B", "assets/fonts/DejaVuSans-Bold.ttf")

	pdf.SetFont("DejaVu", "B", 18)
	pdf.Cell(0, 12, "Детали билета")
	pdf.Ln(14)

	pdf.SetFont("DejaVu", "", 12)

	pdf.Cell(0, 8, fmt.Sprintf("Кинотеатр: %s", cinema.Name))
	pdf.Ln(6)

	pdf.Cell(0, 8, fmt.Sprintf("Фильм: %s", film.Name))
	pdf.Ln(6)

	pdf.Cell(0, 8, fmt.Sprintf(
		"Дата и время: %s",
		screening.Start_At.Format("02.01.2006 15:04"),
	))
	pdf.Ln(6)

	pdf.Cell(0, 8, fmt.Sprintf("Зал: %s", hall.Name))
	pdf.Ln(6)

	pdf.Cell(0, 8, fmt.Sprintf(
		"Место: ряд %d, место %d",
		seat.Row,
		seat.Number,
	))
	pdf.Ln(10)

	opt := gofpdf.ImageOptions{
		ImageType: "PNG",
		ReadDpi:   true,
	}

	pdf.RegisterImageOptionsReader(
		"qr",
		opt,
		bytes.NewReader(qr),
	)

	pdf.Image("qr", 150, 40, 40, 40, false, "", 0, "")

	pdf.Ln(60)
	pdf.SetFont("DejaVu", "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("Номер билета: %d", ticket.ID))
	pdf.Ln(5)
	pdf.Cell(0, 6, fmt.Sprintf("Статус: %s", ticket.Status))

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
