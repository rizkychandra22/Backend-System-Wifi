package controllers

import (
	"backend-wifi/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func CreatePayment(c *gin.Context) {
	var input struct {
		CustomerID    uint `json:"customer_id" binding:"required"`
		WifiPackageID uint `json:"wifi_package_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, appErr := services.CreatePayment(input.CustomerID, input.WifiPackageID)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Payment recorded successfully", "data": payment})
}

func GetCustomerPayments(c *gin.Context) {
	payments, appErr := services.GetCustomerPayments(c.Param("customer_id"))
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payments})
}

func GeneratePaymentPDF(c *gin.Context) {
	paymentID := c.Param("id")
	payment, appErr := services.GetPaymentByID(paymentID)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Invoice Pembayaran Layanan WiFi")

	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("ID Pembayaran: %d", payment.ID))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Tanggal: %s", payment.CreatedAt.Format("02 Jan 2006 15:04")))
	
	pdf.Ln(12)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Detail Customer")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 12)
	if payment.Customer != nil {
		pdf.Cell(40, 10, fmt.Sprintf("Nama: %s", payment.Customer.Name))
		pdf.Ln(8)
		pdf.Cell(40, 10, fmt.Sprintf("No HP: %s", payment.Customer.Phone))
	} else {
		pdf.Cell(40, 10, "Nama: -")
	}

	pdf.Ln(12)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Detail Layanan")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 12)
	if payment.WifiPackage != nil {
		pdf.Cell(40, 10, fmt.Sprintf("Paket: %s", payment.WifiPackage.Name))
	}

	pdf.Ln(12)
	pdf.Cell(40, 10, fmt.Sprintf("Harga Paket: Rp %.0f", payment.PackagePrice))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("PPN (11%%): Rp %.0f", payment.PPN))
	pdf.Ln(8)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Total Bayar: Rp %.0f", payment.TotalAmount))

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=invoice-%d.pdf", payment.ID))
	
	err := pdf.Output(c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}
}
