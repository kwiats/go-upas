# Go User Profile Auth System (Go UPAS)

## Introduction
Go UPAS is a Go-based authentication and user profile management system. This project is designed to provide a robust and efficient solution for handling user authentication, profile storage, and access control in Go applications. Leveraging the concurrency and efficiency of Go, Go UPAS offers a scalable approach to managing user data and securing access.

## Features
- User Authentication: Securely manage user logins with password hashing and verification.
- Profile Management: Create and update user profiles with customizable fields.
- Access Control: Define and manage user roles and permissions.
- Scalable Architecture: Built to handle concurrent requests efficiently in Go.

## Getting Started
These instructions will help you set up Go UPAS on your local machine for development and testing purposes.

### Prerequisites
What you need to install:

```bash
go get -u [dependencies, if any]
```

## Installing
Follow these steps to get a development environment running:

1. **Clone the repository:**
   ```bash
   git clone https://github.com/kwiats/go-upas.git
    ```

2. **Navigate to the project directory:**

    ```bash
    cd go-upas
    ```

3. **Run the application:**
    ```bash
    make run
    ```


## Running the Tests
To run the automated tests:

 ```bash
    make test
```