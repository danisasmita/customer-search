# Customer Search - Setup Guide

## 📥 Clone Backend Repository

Clone the backend repository from GitHub:

```bash
git clone https://github.com/danisasmita/customer-search.git
```

Navigate to the project directory:

```bash
cd customer-search
```

## ⚙️ Configure `.env`

Copy the `.env.example` file to `.env`:

```bash
cp .env.example .env
```

Open the `.env` file and update the PostgreSQL database configuration:

```ini
DB_HOST=localhost  
DB_PORT=5432  
DB_USER=your_username  
DB_PASSWORD=your_password  
DB_NAME=your_database  
```

Replace `your_username`, `your_password`, and `your_database` with your actual database credentials.

## 📚 Install Dependencies

Ensure Go is installed, then run:

```bash
go mod tidy
```

## 🗂️ Run Database Migration & Seeding

If the project requires database migration and initial data seeding, run:

```bash
go run cmd/main.go --migrate --seed
```

After successful execution, you will see an output like this:


Data Seeded  
Seeded Data. Here is the customer data seeded into the database:
  

- John Doe - john@example.com  
- Jane Smith - jane@example.com  
- Robert Johnson - robert@example.com  
- Emily Davis - emily@example.com  
- Michael Wilson - michael@example.com  
- Sarah Brown - sarah@example.com  
- David Lee - david@example.com  
- Jennifer Taylor - jennifer@example.com  
- Kevin Martinez - kevin@example.com  
- Lisa Anderson - lisa@example.com  
- Thomas Wright - thomas@example.com  
```

```
## 🚀 Start the Backend Server

Run the following command to start the backend application:

```bash
go run cmd/main.go
```

---
## 🗂️ Run Unit Test

```bash
go test -v ./...
```

## 🎨 Clone Frontend Repository

To register, log in, and search for customers, clone and run the frontend project:

```bash
git clone https://github.com/danisasmita/customer-search-app.git
```

