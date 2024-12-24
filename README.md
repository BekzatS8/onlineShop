OnlineShop – E-commerce Platform
OnlineShop is a simple e-commerce platform designed for online shopping. It allows users to browse products, add items to the cart, and place orders via a RESTful API while providing a basic frontend interface for interaction.

🌟 Project Overview
Purpose
OnlineShop demonstrates the creation of an e-commerce platform using Go, RESTful APIs, database migrations, and a simple frontend.

Core Features
Supports GET and POST methods for product management and order handling.
Secure data storage in a SQL database.
Basic HTML/CSS frontend interface for user interactions.
Target Audience
Developers learning Go (Golang) for backend development.
Students exploring e-commerce application development.
Teams building scalable web-based e-commerce applications.
👥 Team Members
Sapargali Bekzat
Dinmukhammed Tauasar

🚀 How to Run the Project
Prerequisites
Install Go (version 1.16 or higher): Download Go.
Set up a SQL database (e.g., PostgreSQL or MySQL).
Install golang-migrate for managing database migrations.
Steps to Run OnlineShop
1. Clone the Repository
bash
Копировать код
```bash
   git https://github.com/BekzatS8/onlineShop
   cd online-shop
```
2. Configure Environment Variables
Create a .env file in the project root with the following content:

makefile
Копировать код
DB_HOST=localhost
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=online_shop_db
DB_PORT=5432
3. Apply Database Migrations
Run the following command to apply migrations:

bash
Копировать код
migrate -path ./db/migrations -database "postgres://username:password@localhost:5432/online_shop_db?sslmode=disable" up
4. Start the Server
Run the Go server:

bash
Копировать код
go run main.go
5. Access the Web Interface
Open the following URL in your browser:

bash
Копировать код
http://localhost:8080
📡 API Endpoints
Method	Endpoint	Description
GET	/products	Retrieve all available products
POST	/order	Place a new order in JSON format
🛠️ Tools and Technologies
Go: Backend development.
gorilla/mux: HTTP routing.
SQL: Database for storing products and orders.
golang-migrate: Manages database migrations.
HTML/CSS: Provides a simple frontend interface.
💡 Future Enhancements
Implement user authentication (login/logout).
Integrate a payment gateway for processing orders.
Build a dynamic frontend using JavaScript for a better user experience.