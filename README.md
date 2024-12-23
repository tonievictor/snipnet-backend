# Snipnet (**backend**)

## Prerequisites
Before you begin, ensure your environment has the following tools installed:
- Docker (for containerization)
- Docker Compose (for orchestrating multi-container Docker applications)
- Git (for version control)
- [Migrate](https://github.com/golang-migrate/migrate) (for database migrations)

## Setup
To set up the backend application, follow these steps:

1. Clone the repository to your local machine:
```bash
git clone https://github.com/tonievictor/snipnet-backend.git
```
2. Navigate into the project directory
```bash
cd snipnet-backend
```
3. Copy the `.env.example` file to `.env` and provide the necessary values for your local setup:
```bash
cp .env.example .env
```
> Note: When using Docker Compose, set the database host to the service name defined in `compose.yml`.

4. To apply database migrations, use the following command:
```bash
migrate -database ${POSTGRESQL_URL} -path migrations up
```
> Note: If you're using Docker Compose to run the database, it will be available on localhost:5432

5. To launch the application, simply run:
```bash
docker compose up
```

## Documentation
All endpoints are documented using Swagger. Once the application is running, you can access the Swagger UI at:
```bash
http://localhost:<PORT>/swagger
```
