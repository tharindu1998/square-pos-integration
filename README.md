# square-pos-integration
A multi-tenant backend REST API written in Go, designed to integrate with Square POS. This system enables restaurants to manage orders and process payments securely and efficiently.

# Prerequisites

- Go 1.19 or higher
- MySQL 5.7 or higher
- Square Developer Account
- Square Application ID and credentials

# Installation

1. Clone the repository:

~~~bash  
git clone https://github.com/tharindu1998/square-pos-integration.git
cd square-pos-integration
~~~

2. Install dependencies:

~~~bash  
go mod download
~~~

3. Set up your database:

~~~bash  
CREATE DATABASE square_pos_db;
~~~

4. Configure environment variables (see Configuration section below)
Run the application:

~~~bash  
go run main.go
~~~

# Configuration

Create a .env file in the root directory with the following configuration:

~~~bash  
# Database Configuration
DB_DSN=username:password@tcp(localhost:3306)/square_pos_db?charset=utf8mb4&parseTime=True&loc=Local

# Server Configuration
PORT=8080

# JWT Secret Key (Change this in production)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Square API Configuration
SQUARE_APPLICATION_ID=your_square_app_id
SQUARE_ENVIRONMENT=sandbox # or production

# Logging
LOG_LEVEL=info
~~~

# Square Setup

1. Create a Square Developer Account:

Visit Square Developer Dashboard
Create a new application or use an existing one


2. Get Your Credentials:

Navigate to your application in the Square Developer Dashboard
Copy your Application ID from the Credentials tab
Generate and copy your Access Token


3. Configure Webhooks (Optional):

Set up webhook endpoints for real-time payment updates
Configure webhook signature verification

# API Endpoints

1. Authentication (Public)
- POST /api/v1/register-restaurant – Register a new restaurant and admin user

- POST /api/v1/login – Authenticate a user and return a JWT token

2. Profile (Protected)
- GET /api/v1/profile – Retrieve the authenticated user's profile

3. Orders (Protected)
- POST /api/v1/orders – Create a new order

- GET /api/v1/orders/table/:table_number – Get an order by table number

- GET /api/v1/orders/:id – Get an order by order ID

4. Payments (Protected)
- POST /api/v1/payment/:id/payment-intent – Create a payment intent for an order

- POST /api/v1/payment/complete – Complete a payment

5. Admin (Protected - Admin Role Only)
- POST /api/v1/admin/users – Create a new user (Admin only)

# License
This project is licensed under the MIT License - see the LICENSE file for details.
# Support
For support and questions:

- Create an issue in the GitHub repository
- Review Square's developer documentation
- Check the Square Developer Community

# Acknowledgments

- Square Developer Platform for comprehensive APIs
- Gin Web Framework for HTTP router
- GORM for database ORM
- Go community for excellent libraries and tools



