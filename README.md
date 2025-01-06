# tribe-payments-wallet-golang-interview-assignment


---

# Setup Instructions

Use the terminal to navigate to your playground. Example:

```bash
cd E:\
```

Cool now you have to make a copy of my repo.

```bash
git clone https://github.com/evelinkolev/tribe-payments-wallet-golang-interview-assignment
```

Open the folder in IDE of your choice. I have used VS Code with Go extension. Initialise the go.mod (only if does not exist).

```bash
go mod init tribe-payments-wallet-golang-interview-assignment
```

Install dependencies.

```bash
go get github.com/microsoft/go-mssqldb
go get -u github.com/golang-migrate/migrate/v4
go get -u github.com/golang-migrate/migrate/v4/database/sqlserver
go get -u github.com/golang-migrate/migrate/v4/source/file
go install -tags 'sqlserver' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Open SSMS and create database with name wallet_db. Run the migrate command. (only If you have issue with migrate like in my case, check https://github.com/golang-migrate/migrate and build the .exe file on your machine).

```bash
migrate -database "sqlserver://localhost:1433?database=wallet_db&integratedSecurity=true&trustServerCertificate=true" -path ./migrations up
```

Start the application

```bash
go run . api
```

The API will be available at: [http://localhost:8080](http://localhost:8080)

---

# API Documentation

## API Endpoints

- **Create Wallet:** `POST /v1/wallets`
- **Get Wallet:** `GET /v1/wallets/{id}`
- **Deposit Funds:** `POST /v1/wallets/{id}/deposit`
- **Withdraw Funds:** `POST /v1/wallets/{id}/withdraw`

---

## Create Wallet

**Request:**
```powershell
curl.exe -X POST http://localhost:8080/v1/wallets -H "Content-Type: application/json" -d '{\"currency\": \"USD\"}'
```

**Response:**
```json
{
  "id": "34fde074-262c-4ba4-8104-ec09e7a39e12",
  "balance": 0,
  "currency": "USD",
  "created_at": "2025-01-06T08:42:10+02:00",
  "updated_at": "2025-01-06T08:42:10+02:00"
}
```

## Get Wallet

**Request:**
```powershell
curl.exe -X GET http://localhost:8080/v1/wallets/34fde074-262c-4ba4-8104-ec09e7a39e12
```

**Response:**
```json
{
  "id": "34fde074-262c-4ba4-8104-ec09e7a39e12",
  "balance": 0,
  "currency": "USD",
  "created_at": "2025-01-06T08:42:10Z",
  "updated_at": "2025-01-06T08:42:10Z"
}
```

## Deposit Funds

**Request:**
```powershell
curl.exe -X POST "http://localhost:8080/v1/wallets/34fde074-262c-4ba4-8104-ec09e7a39e12/deposit" `
-H "Content-Type: application/json" `
-d '{\"balance\": 100.50}'
```

**Response:**
```http
HTTP/1.1 204 No Content
```

## Withdraw Funds

**Request:**
```powershell
curl.exe -X POST "http://localhost:8080/v1/wallets/34fde074-262c-4ba4-8104-ec09e7a39e12/withdraw" `
-H "Content-Type: application/json" `
-d '{\"balance\": 100.50}'
```

**Response:**
```http
HTTP/1.1 204 No Content
```


---

## Troubleshooting

- Ensure the database is running and accessible.
- Check logs for errors during startup or requests.


