# UFMarketplace

## Group 35

## Problem Statement:

The University of Florida currently lacks a university-specific online marketplace where students and staff can post, browse, and connect with each other to buy or sell products. Existing solutions like social media groups and bulletin boards are unstructured, inefficient, and often pose trust and security challenges. A dedicated platform is needed to streamline these transactions while fostering a trusted community environment where users can connect easily and securely.

## About:

UFMarketplace is a web-based platform exclusively designed for the University of Florida community, enabling students and staff to post and browse product listings. Built with Golang for the backend and React for the frontend, the platform provides a simple and secure way for users to showcase items for sale and connect with potential buyers or sellers.

## Contributors:

- Manikumar Honnenahalli Lakshminarayana Swamy (front-end)
- Rushang Sunil Chiplunkar (back-end)
- Nischith Bairannanavara Omprakash (back-end)
- Mayur Sai Yaram (front-end)

--

# Steps to intsall and run

## Backend

---

### Database Setup (PostgreSQL)

1. **Install PostgreSQL** on your system (from [https://www.postgresql.org/download](https://www.postgresql.org/download) or using your OS package manager).
2. **Create a new database** (e.g., `ufmarketplace`).
3. **Create a user** (e.g., `ufmarketplace`) and grant it admin or full write permissions on the database.
4. **Get your PostgreSQL connection string** in the following format: postgres://<username>:<password>@<host>:<port>/<dbname>?sslmode=disable
5. Set the POSTGRES_CONN Environment Variable

This connection string must be set as an environment variable so the application can connect to the database.

### Linux/macOS (Temporary)

```bash
export POSTGRES_CONN="postgres://ufmarketplace:yourpassword@localhost:5432/ufmarketplace?sslmode=disable"
```

### Windows CMD (Temporary)

```cmd
set POSTGRES_CONN=postgres://ufmarketplace:yourpassword@localhost:5432/ufmarketplace?sslmode=disable
```

### Windows PowerShell (Temporary)

```powershell
$env:POSTGRES_CONN="postgres://ufmarketplace:yourpassword@localhost:5432/ufmarketplace?sslmode=disable"
```

---

## Install Go

1. Download Go from the official site:  
   https://go.dev/dl/

2. Run the installer and follow the instructions for your OS.

3. After installation, verify it works:

```bash
go version
```

You should see the installed version printed.

---

## Clone the Repository

Clone the project to your local machine using Git:

```bash
git clone https://github.com/MANIKUMAR98/UFMarketPlace.git
cd UFMarketPlace
```

## Build and Run the Application

1. Navigate to the backend directory:

```bash
cd Backend
```

2. Build the Go application:

```bash
go build
```

4. Run Unit Tests:

```bash
go test -v

```

5. Run the application:

```bash
go run .

```

---
