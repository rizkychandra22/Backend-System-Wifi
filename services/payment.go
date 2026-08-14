package services

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"backend-wifi/utils"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
)

func CreatePayment(customerID, WifiPackageID uint) (*models.Payment, *utils.AppError) {
	var customer models.User
	if err := config.DB.First(&customer, customerID).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Customer not found")
	}
	if customer.Role != models.RoleCustomer {
		return nil, utils.NewAppError(http.StatusBadRequest, "User is not a customer")
	}

	var ws models.WifiPackage
	if err := config.DB.First(&ws, WifiPackageID).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Wifi service not found")
	}

	totalAmount := ws.Price
	ppn := totalAmount * 0.11
	packagePrice := totalAmount - ppn

	// Generate unique random invoice number INV-XXXX-XXXX
	n1, _ := rand.Int(rand.Reader, big.NewInt(10000))
	n2, _ := rand.Int(rand.Reader, big.NewInt(10000))
	invoiceNumber := fmt.Sprintf("INV-%04d-%04d", n1.Int64(), n2.Int64())

	payment := models.Payment{
		CustomerID:    customerID,
		WifiPackageID: WifiPackageID,
		PackagePrice:  packagePrice,
		PPN:           ppn,
		TotalAmount:   totalAmount,
		Status:        "paid",
		InvoiceNumber: invoiceNumber,
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Failed to create payment")
	}

	return &payment, nil
}

func GetCustomerPayments(customerID string) ([]models.Payment, *utils.AppError) {
	var payments []models.Payment
	if err := config.DB.Preload("WifiPackage").Where("customer_id = ?", customerID).Order("created_at desc").Find(&payments).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Failed to retrieve payments")
	}
	return payments, nil
}

func GetPaymentByID(id string) (*models.Payment, *utils.AppError) {
	var payment models.Payment
	if err := config.DB.Preload("Customer").Preload("WifiPackage").First(&payment, id).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Payment not found")
	}
	return &payment, nil
}

func GetAllPayments() ([]models.Payment, *utils.AppError) {
	var payments []models.Payment
	if err := config.DB.Preload("Customer").Preload("WifiPackage").Order("created_at desc").Find(&payments).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Failed to retrieve all payments")
	}
	return payments, nil
}

func UpdatePayment(id string, customerID, wifiPackageID uint) (*models.Payment, *utils.AppError) {
	var payment models.Payment
	if err := config.DB.First(&payment, id).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Payment not found")
	}

	var customer models.User
	if err := config.DB.First(&customer, customerID).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Customer not found")
	}
	if customer.Role != models.RoleCustomer {
		return nil, utils.NewAppError(http.StatusBadRequest, "User is not a customer")
	}

	var ws models.WifiPackage
	if err := config.DB.First(&ws, wifiPackageID).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Wifi service not found")
	}

	totalAmount := ws.Price
	ppn := totalAmount * 0.11
	packagePrice := totalAmount - ppn

	payment.CustomerID = customerID
	payment.WifiPackageID = wifiPackageID
	payment.PackagePrice = packagePrice
	payment.PPN = ppn
	payment.TotalAmount = totalAmount

	if err := config.DB.Save(&payment).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Failed to update payment")
	}

	return &payment, nil
}

func DeletePayment(id string) *utils.AppError {
	var payment models.Payment
	if err := config.DB.First(&payment, id).Error; err != nil {
		return utils.NewAppError(http.StatusNotFound, "Payment not found")
	}

	if err := config.DB.Delete(&payment).Error; err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "Failed to delete payment")
	}

	return nil
}
