# SnipNet - Code Snippet Management Application
SnipNet is a versatile code snippet management application. Currently, it provides functionality through a terminal interface and is planned to be expanded into a web application in the future.


## Features
- User authentication (signup, signin, signout)
- CRUD operations on snippets
- User management (retrieve, update, delete users)
- Secure endpoints with authentication middleware

## Future Web Interface (planned):
- Intuitive web-based user interface
- Enhanced collaboration features
- Integration with other development tools

## Installation
To install SnipNet, follow these steps:
- Clone this repository:
```bash
git clone https://github.com/tonie-ng/nest-backend.git snipnet
```
- Navigate to the project directory:
```bash
cd snipnet
```
- Install dependencies
```
go mod tidy
```
- Set up environment variables
```md
PORT=address:port
DB_USER=
DB_NAME=
DB_PASSWORD=
DB_PORT=
DB_HOST=
DB_CONN_STRING="postgres://postgres:password@localhost:5432/testdb?sslmode=disable"
REDIS_URL=
```
- Optionally, setup the database (optional because you might choose another database config)
```bash
./setup-db.sh
```

- Start the backend application
```bash
cd cmd
go run main.go
```

## Usage
### Terminal Interface
To use SnipNet via the terminal interface, execute the following command:
```
./snipnet.sh
```

### Web Interface (planned)
The web interface for SnipNet is currently under development. Stay tuned for updates on how to access and use the web application.

Thank you.
