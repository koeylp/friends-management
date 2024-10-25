
# Friends Management

A RESTful API for managing friend relationships, built using Go and PostgreSQL. This application allows users to add, manage, and interact with their friend lists.

## Table of Contents
- [Features](#features)
- [Getting Started](#getting-started)
- [Using Docker](#using-docker)
- [Success Cases](#success-cases)
- [Error Cases](#error-cases)

## Features
- Create a friend connection
- Retrieve the friends list for an email address
- Retrieve the common friend list between two email addresses
- Subscribe to updates from an email address
- Block updates from an email address
- Retrieve all updatable email addresses

## Getting Started

### Prerequisites
- Go 1.18 or later
- PostgreSQL
- Docker (optional, for running PostgreSQL in a container)

## Running the Application

### Using Docker



  **Start the application:**
   Run the following command in the root directory of your project:
   ```bash
   docker-compose up --build -d
   ```

   This will build the Docker images and start the PostgreSQL database and your application. The application will be available at `http://localhost:8080`.

### Stopping the Application
To stop the application, run:
```bash
docker-compose down
```

## Success Cases

### Example Request
- **Endpoint:** `POST /api/v1/friends`
- **Request Body:**
  ```json
  {
      "friends": [
          "user2@example.com",
          "user6@example.com"
      ]
  }
  ```
### Retrieve Friends List
- **Endpoint:** `POST /api/friends/list`
- **Example Response:**
  ```json
  {
     "email": "user6@example.com"
  }
  ```
### Retrieve the common friend list
- **Endpoint:** `POST /api/friends/common-list`
- **Example Response:**
  ```json
  {
      "email": "user6@example.com"
  }
  ```
### Subscribe updates 
- **Endpoint:** `POST /api/subcription`
- **Example Response:**
  ```json
  {
      "requestor": "user@example.com",
      "target": "user@example.com"
  }
  ```
### Block updates
- **Endpoint:** `POST /api/block`
- **Example Response:**
  ```json
  {
      "requestor": "user@example.com",
      "target": "user@example.com"
  }
  ```
### Retrieve all updatable email addresses
- **Endpoint:** `POST /api/subcription/recipients`
- **Example Response:**
  ```json
  {
   "sender": "user9@example.com",
   "text": "Hello World! user7@example.com, user8@example.com, user9@example.com"
  }
  ```
  
### Example Response
- **Status Code:** 201 Created
- **Response Body:**
  ```json
  {
      "success": true
  }
  ```
### Example Response
- **Status Code:** 200 Ok
- **Response Body:**
  ```json
  {
    "count": 2,
    "recipients": [
        "user7@example.com",
        "user8@example.com"
    ],
    "success": true
  }
  ```

## Error Cases

### Example Error Response
- **Error Case:** User not found
- **Endpoint:** `POST /api/friends`
- **Request Body:**
  ```json
  {
    "Message": "user not found with email user222@example.com",
    "Time": "2024-10-24T17:14:11.389023506+07:00"
  }
  ```

### Example Error Response
- **Status Code:** 400 Bad request error
- **Response Body:**
  ```json
  {
    "Message": "Invalid request payload",
    "Time": "2024-10-24T17:22:18.835078843+07:00"
  }
  ```


### Database Error
- **Error Case:** Database connection issue
- **Example Response:**
- **Status Code:** 500 Internal Server Error
- **Response Body:**
  ```json
  {
      "error": "Database connection error."
  }
  ```



