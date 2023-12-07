# Microservice Readme

## Description

This Golang microservice is designed for creating and applying promo codes to users. Interaction with this service involves working with another service responsible for storing and managing users.

## Technologies

1. **Programming Language**: Golang
2. **Framework**: Gin (for handling HTTP requests)
3. **ORM**: Gorm (for interacting with PostgreSQL database)
4. **Database**: PostgreSQL
5. **Logging**: Output to standard input stream (stdin)

## Available API of the External User Service

### Get available users for applying a promo code

**Endpoint**: `GET api/users/promo_type/:type/available`

**Response**:

```json
{
  "user_id": [1, 5, 7, 9]
}
```

### Get users to whom the promo code is applied

**Endpoint**: `GET api/users/promo/:id/applied`

**Response**:

```json
{
  "users_id": [1, 2, 3, 4]
}
```

### Apply promo code to users

**Endpoint**: `POST api/users/promo/:id/apply`

**Request**:

```json
{
  "users_id": [1, 2, 6, 9]
}
```

**Response status**: "ok"

## Microservice API Contract

### Check service's availability

**Endpoint**: `GET /check_alive`

**Response status**: 200

### Create a promo code

**Endpoint**: `POST /promo/create`

**Request**:

```json
{
  "name": "New Promo",
  "slug": "new-promo",
  "url": "https://www.example.com/promo",
  "description": "Try our new promo code. You'll love it.",
  "type": 1
}
```

### Add promo codes to a promo

**Endpoint**: `POST /promo/:id/codes`

**Request**:

```json
{
  "codes": ["123", "456", "789"]
}
```

### Get promo with promo codes

**Endpoint**: `GET /promo/:id`

### Apply promo to all eligible users

**Endpoint**: `POST /promo/:id/apply_all`

### Apply promo to specified users

**Endpoint**: `POST /promo/:id/apply_users`

**Request**:

```json
{
  "users_id": [1, 4, 8]
}
```

### Get users for whom the promo is applied

**Endpoint**: `GET /promo/:id/users`

**Response**:

```json
{
  "users_id": [1, 3, 6, 8]
}
```

## Running the Microservice

To run the microservice, follow these steps:

1. **Clone the repository:**

   ```bash
   git clone https://github.com/iamrosada/microservice-goland
   ```

2. **Navigate to the microservice directory:**

   ```bash
   cd microservice-goland
   ```

3. **Run the user service:**
   The user service will be available on port 8000.

   ```bash
   cd user-service
   go run main.go
   ```

4. **Open a new terminal and run the promo service:**
   The promo service will be available on port 8080.

   ```bash
   cd /promo-service
   go run main.go
   ```

After completing these steps, the microservice will be accessible through the specified endpoints. Make sure you have Go installed on your system before running these commands.
