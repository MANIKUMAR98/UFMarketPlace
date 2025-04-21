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

## Frontend

### Prerequisites

- **Node.js and npm:** Make sure Node.js and npm are installed. Download the latest LTS version from [https://nodejs.org/en/download](https://nodejs.org/en/download).

### Install and Run the Frontend

1. **Navigate to the Frontend directory:**

   ```bash
   cd UFMarketPlace/Frontend
   ```

2. **Install the dependencies:**

   ```bash
   npm install
   ```

3. **Start the React development server:**

   ```bash
   npm run dev
   ```

4. **Open your browser and go to:**
   ```
   http://localhost:5173
   ```

### Running Test Cases

- **Unit Tests:** To run the unit tests, use the following command:

  ```bash
  npm test
  ```

- **Cypress Tests:** For end-to-end testing with Cypress:
  - Ensure both the frontend and backend applications are running.
  - Open the Cypress test runner with:
    ```bash
    npx cypress open
    ```

### Notes

- The frontend is built using React and Vite.
- Make sure the backend server is running before using the frontend. It is expected to be available at `http://localhost:8080`.
- If your backend is on a different URL or port, update the API base URL in the frontend code.

---

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
