# 📡 WiFi Payment & Employee Attendance API

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)
![Gin](https://img.shields.io/badge/Gin-Gonic-00ADD8?style=for-the-badge&logo=go)
![GORM](https://img.shields.io/badge/GORM-DB-red?style=for-the-badge)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql)

Modern Backend API built with **Golang** and **Gin-Gonic** to serve as the core engine for a Field Employee Attendance system and a WiFi Subscriber Billing & Invoicing system.

---

## ✨ Features

- **🛡️ Role-Based Access Control (RBAC):** Supports 3 distinct roles: `Admin` (Business Owner), `Employee` (Field Worker), and `Customer` (WiFi Subscriber).
- **📍 GPS-Based Attendance:** Employees can check-in and check-out daily with GPS coordinates and photo evidence.
- **💰 Smart Billing System:** Generate monthly invoices for WiFi subscribers.
- **✅ Payment Verification:** Customers can upload transfer proofs, and Admins can verify them directly.
- **📱 WhatsApp Integration (Planned):** Automated invoice and payment reminders via WhatsApp.

## 🚀 Tech Stack

- **Language:** [Go (Golang)](https://go.dev/)
- **Framework:** [Gin](https://gin-gonic.com/) (Fast HTTP web framework)
- **Database:** PostgreSQL
- **ORM:** [GORM](https://gorm.io/) (The fantastic ORM library for Golang)
- **Configuration:** Godotenv

## 📁 Project Structure

```text
Backend-System/
├── config/           # Database and Environment Configurations
├── controllers/      # HTTP Request Handlers (API Endpoints)
├── middlewares/      # Interceptors (JWT Auth, Error Handling, etc.)
├── models/           # Database Entities / Structs (GORM)
├── repositories/     # Database Queries (Data Access Layer)
├── routes/           # API Routing Definitions
├── services/         # Core Business Logic
└── main.go           # Application Entry Point
```

## 🛠️ Quick Start

### 1. Prerequisites
- Go 1.20 or higher installed.
- PostgreSQL running locally or on a server.

### 2. Environment Setup
Create a `.env` file in the root of the project with the following keys:
```ini
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_database_name
DB_PORT=5432
DB_SSLMODE=disable
```

### 3. Run the Application
1. Download dependencies:
   ```bash
   go mod tidy
   ```
2. Start the development server:
   ```bash
   go run main.go
   ```
3. The API will be available at `http://localhost:8080/api/health`

---
*Built for Modern ISP & Workforce Management.*
