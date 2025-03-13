1. Clone Repository

Clone the repository from GitHub using the following command:

git clone https://github.com/danisasmita/customer-search.git

2. Navigate to the Project Directory

cd customer-search

3. Create and Configure the .env File

Copy the .env.example file to .env:

cp .env.example .env

Then, open the .env file and configure the PostgreSQL database settings:

DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database

Make sure to replace your_username, your_password, and your_database with the appropriate values.

4. Install Dependencies

Ensure you have Go installed, then run the following command to download the project dependencies:

go mod tidy

5. Run Database Migration (Optional)

If the project requires database migration, run the following command:

go run cmd/main.go migrate

Or use the specific migration command for the project.

6. Start the Application

Use the following command to run the application:

go run cmd/main.go

Or if using air for hot reload:

air

