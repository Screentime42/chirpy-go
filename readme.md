# 🐦 chirpy-go  
A small Go backend built as a learning project. It implements a simple “chirp” (micro‑post) API along with user authentication, token refresh/revocation, and basic admin/metrics endpoints. The goal is to practice idiomatic Go server design, routing, validation, and state management.



## 🚀 Features

- User registration & login  
- Access token + refresh token flow  
- CRUD operations for “chirps”  
- Basic validation (e.g., banned words)  
- Admin endpoints for metrics & resets  
- Readiness/health check endpoint  
- Simple in‑memory or file‑based persistence (depending on implementation)



## 📁 Project Structure

```.
├── main.go              # Server entrypoint & route registration
├── metrics.go           # Metrics counter + /admin/metrics
├── reset.go             # /admin/reset handler
├── readiness.go         # /api/healthz handler
├── validate.go          # Input validation helpers
├── users.go             # User creation, login, update
├── chirps.go            # Chirp creation, listing, retrieval, deletion
└── ...
```



## 🧩 API Endpoints

### **Health**
- **GET `/api/healthz`**  
  Readiness probe. Returns a simple success response.

---

### **Admin**
- **POST `/admin/reset`**  
  Resets internal metrics counters.

- **GET `/admin/metrics`**  
  Returns current metrics (e.g., file server hit count).

---

### **Users**
- **POST `/api/users`**  
  Create a new user.

- **POST `/api/login`**  
  Authenticate a user and return tokens.

- **POST `/api/refresh`**  
  Exchange a refresh token for a new access token.

- **POST `/api/revoke`**  
  Revoke a refresh token.

- **PUT `/api/users`**  
  Update the authenticated user’s credentials.

- **POST `/api/polka/webhooks`**  
  Handle webhook events that upgrade a user (e.g., premium status).

---

### **Chirps**
- **POST `/api/chirps`**  
  Create a new chirp (with validation).

- **GET `/api/chirps`**  
  Retrieve all chirps.

- **GET `/api/chirps/{id}`**  
  Retrieve a single chirp by ID.

- **DELETE `/api/chirps/{chirp_id}`**  
  Delete a chirp (typically requires authorization).



## 🏃 Running the Project

No special setup required — just Go.

```bash
go run main.go
```
The server will start on the configured port (often :8080 unless changed in code).



## 🎯 Purpose of This Project

This repository exists purely as a learning exercise to practice:

Idiomatic Go HTTP server patterns

Handler struct organization

Token-based authentication

Basic validation pipelines

Route design

State and metrics management

It is not intended for production use.



## 🧪 Example Requests

Create a user
```bash
curl -X POST http://localhost:8080/api/users \
  -d '{"email":"test@example.com","password":"secret"}'
```
Create a chirp
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer <token>" \
  -d '{"body":"Hello world!"}'
  ```