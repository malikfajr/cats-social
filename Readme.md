# Cats Social

[Cats Social](https://openidea-projectsprint.notion.site/Cats-Social-9e7639a6a68748c38c67f81d9ab3c769) is an API for a social networking application where cat owners can match their cats with other cats. 

## 📜Table of Contents

- [Cats Social](#cats-social)
  - [📜Table of Contents](#table-of-contents)
  - [🔍Requirements](#requirements)
  - [🛠️Installation](#installation)
  - [🌟Features](#features)
  - [🚀Usage](#usage)
  - [⚙️Configuration](#configuration)
  - [💾Database Migration](#database-migration)
  - [🔨Build](#build)
  - [🧪Testing](#testing)

## 🔍Requirements

This application requires the following:

- [Golang](https://golang.org/dl/)
- [PostgreSQL](https://www.postgresql.org/)
- [Golang Migrate](https://github.com/golang-migrate/migrate)

## 🛠️Installation

To install the boilerplate, follow these steps:

1. Make sure you have [Golang](https://golang.org/dl/),  [PostgreSQL](https://www.postgresql.org/), and [Golang Migrate](https://github.com/golang-migrate/migrate) installed and configured on your system.

2. Clone this repository:

   ```bash
   git clone https://github.com/malikfajr/cats-social.git
   ```

3. Navigate to the project directory:

   ```bash
   cd cats-social
   ```

4. Download the required dependencies:

   ```bash
   go mod download
   ```


## 🌟Features

Cats Social offers the following features:

- **Authentication**:
  - User registration
  - User login
- **Cat Management (CRUD)**:
  - Create new cat profiles
  - View existing cat profiles
  - Update cat profiles
  - Delete cat profiles
- **Matching**:
  - Match your cat with other cats
  - View matching cats
  - Approve or reject matches
  - Delete matches

## 🚀Usage

1. **Setting Up Environment Variables**

   Before starting the application, you need to set up the following environment variables:

   - `DB_HOST`: Hostname or IP address of your PostgreSQL server
   - `DB_PORT`: Port of your PostgreSQL database (default: 5432)
   - `DB_USERNAME`: Username for your PostgreSQL database
   - `DB_PASSWORD`: Password for your PostgreSQL database
   - `DB_NAME`: Name of your PostgreSQL database
   - `DB_PARAMS` : Additional connection parameters for PostgreSQL (e.g., sslmode=disable)
   - `JWT_SECRET`: Secret key used for generating JSON Web Tokens (JWT)
   - `BCRYPT_SALT`: Salt for password hashing (use a higher value than 8 in production!)

2. **Database Migrations**

   Cats Social uses Golang Migrate to manage database schema changes. Here's how to work with migrations:

   - Apply migrations to the database:

     ```bash
     migrate -database "postgres://username:password@host:port/dbname?sslmode=disable" -path db/migrations up
     ```

3. **Running the Application**

   Once you have set up the environment variables, you can start the application by running:

   ```bash
   go run main.go
   ```

   This will start the Cats Social application on the default port (usually 8080). You can access the application in your web browser at `http://localhost:8080`.

## ⚙️Configuration

The application uses environment variables for configuration. You can configure the database connection, JWT secret, and bcrypt salt by setting the following environment variables:

- Refer to the [Usage](#usage) section for a detailed explanation of each environment variable.

## 💾Database Migration

Database migration must use golang-migrate as a tool to manage database migration

1. Direct your terminal to your project folder first

2. Initiate folder

   ```bash
   mkdir db/migrations
   ```

3. Create migration

   ```bash
   migrate create -ext sql -dir db/migrations add_user_table
   ```

   This command will create two new files named `add_user_table.up.sql` and `add_user_table.down.sql` inside the `db/migrations` folder

   - `.up.sql` can be filled with database queries to create / delete / change the table
   - `.down.sql` can be filled with database queries to perform a `rollback` or return to the state before the table from `.up.sql` was created

4. Execute migration

   ```bash
   migrate -database "postgres://username:password@host:port/dbname?sslmode=disable" -path db/migrations up
   ```

5. Rollback migration

   ```bash
   migrate -database "postgres://username:password@host:port/dbname?sslmode=disable" -path db/migrations down
   ```

6. View the current migration state

   ```bash
   migrate -database "postgres://username:password@host:port/dbname?sslmode=disable" version 
   ```

## 🔨Build

To build the Cats Social API for different operating systems and architectures, you can use the following commands:

1. **Windows (amd64)**:

    ```bash
    GOOS=windows GOARCH=amd64 go build -o cats-social main.go
    ```

2. **Linux (amd64)**:

    ```bash
    GOOS=linux GOARCH=amd64 go build -o cats-social main.go
    ```

3. **macOS (amd64)**:

    ```bash
    GOOS=darwin GOARCH=amd64 go build -o cats-social main.go
    ```

4. **Linux (ARM)**:
    ```bash
    GOOS=linux GOARCH=arm go build -o cats-social main.go
    ```

## 🧪Testing

To test the Cats Social API, you can use the testing [repository](https://github.com/nandanugg/ProjectSprintBatch2Week1TestCases) provided.