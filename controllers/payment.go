package controllers

import (
	"backend-wifi/models"
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

	userIDClaim, _ := c.Get("userID")
	userID := uint(userIDClaim.(float64))

	payment, appErr := services.CreatePayment(input.CustomerID, input.WifiPackageID, userID)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pembayaran berhasil dicatat", "data": payment})
}

func GetCustomerPayments(c *gin.Context) {
	userRoleClaim, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userRole := userRoleClaim.(string)
	
	if userRole == string(models.RoleCustomer) {
		userIDClaim, _ := c.Get("userID")
		userID := uint(userIDClaim.(float64))
		
		if fmt.Sprintf("%d", userID) != c.Param("customer_id") {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own payments"})
			return
		}
	}

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

	invoiceNum := payment.InvoiceNumber
	if invoiceNum == "" {
		invoiceNum = fmt.Sprintf("INV-%04d", payment.ID)
	}

	customerName := "Customer"
	if payment.Customer != nil && payment.Customer.Name != "" {
		customerName = payment.Customer.Name
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Invoice Pembayaran Layanan WiFi")

	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("ID Pembayaran: %s", invoiceNum))
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

	filename := fmt.Sprintf("%s - %s.pdf", customerName, invoiceNum)

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	
	err := pdf.Output(c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}
}

func GetAllPayments(c *gin.Context) {
	payments, appErr := services.GetAllPayments()
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payments})
}

func UpdatePayment(c *gin.Context) {
	var input struct {
		CustomerID    uint `json:"customer_id" binding:"required"`
		WifiPackageID uint `json:"wifi_package_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentID := c.Param("id")
	payment, appErr := services.UpdatePayment(paymentID, input.CustomerID, input.WifiPackageID)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran berhasil diperbarui", "data": payment})
}

func DeletePayment(c *gin.Context) {
	paymentID := c.Param("id")
	appErr := services.DeletePayment(paymentID)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran berhasil dihapus"})
}
