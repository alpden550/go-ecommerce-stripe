package main

import (
	"fmt"
	"net/http"

	"github.com/alpden550/go-ecommerce-stripe/internal/helpers"
	"github.com/alpden550/go-ecommerce-stripe/internal/mailer"
	"github.com/alpden550/go-ecommerce-stripe/internal/models"
	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
)

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func CreateAndSendInvoice(writer http.ResponseWriter, request *http.Request) {
	var invoice models.Invoice

	err := helpers.ReadJSON(writer, request, &invoice)
	if err != nil {
		helpers.BadRequest(writer, request, err)
		return
	}

	err = createInvoicePDF(invoice)
	if err != nil {
		helpers.BadRequest(writer, request, err)
		return
	}

	attachments := []string{
		fmt.Sprintf("./invoices/%d.pdf", invoice.ID),
	}
	var data interface{}
	if err = mailer.SendEmail(
		invoicer,
		invoicer.Config.SMTP.EmailFrom,
		invoice.Email,
		"Your Invoice",
		"invoices/invoice.plain.tmpl",
		"invoices/invoice.html.tmpl",
		attachments,
		data,
	); err != nil {
		invoicer.ErrorLog.Printf("%w", fmt.Errorf("%e", err))
		_ = helpers.BadRequest(writer, request, err)
	}

	response := Response{
		Error:   false,
		Message: fmt.Sprintf("Invoice id %d created and sent to %s", invoice.ID, invoice.Email),
	}
	err = helpers.WriteJSON(writer, http.StatusCreated, response)
	if err != nil {
		helpers.BadRequest(writer, request, err)
		return
	}
}

func createInvoicePDF(invoice models.Invoice) error {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(10, 13, 10)
	pdf.SetAutoPageBreak(true, 0)

	importer := gofpdi.NewImporter()

	t := importer.ImportPage(pdf, "./pdf-templates/invoice.pdf", 1, "/MediaBox")

	pdf.AddPage()
	importer.UseImportedTemplate(pdf, t, 0, 0, 215.9, 0)

	// write info
	pdf.SetY(50)
	pdf.SetX(10)
	pdf.SetFont("Times", "", 11)
	pdf.CellFormat(97, 8, fmt.Sprintf("Attention: %s %s", invoice.FirstName, invoice.LastName), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(97, 8, invoice.Email, "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(97, 8, invoice.CreatedAt.Format("2006-01-02"), "", 0, "L", false, 0, "")

	pdf.SetX(58)
	pdf.SetY(93)
	pdf.CellFormat(155, 8, invoice.Product, "", 0, "L", false, 0, "")
	pdf.SetX(166)
	pdf.CellFormat(20, 8, fmt.Sprintf("%d", invoice.Quantity), "", 0, "C", false, 0, "")

	pdf.SetX(185)
	pdf.CellFormat(20, 8, fmt.Sprintf("$%.2f", float32(invoice.Amount/100.0)), "", 0, "R", false, 0, "")

	invoicePath := fmt.Sprintf("./invoices/%d.pdf", invoice.ID)
	err := pdf.OutputFileAndClose(invoicePath)
	if err != nil {
		return err
	}

	return nil
}
