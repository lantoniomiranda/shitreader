# ShitReader

A Go application for importing Portuguese electricity sector data from Excel files into a PostgreSQL database.

## Features

- Imports data from multiple Excel files containing regulatory tables
- Processes geographic data (countries, districts, municipalities, parishes)
- Handles CAE (economic activity classification) data
- Progress tracking with real-time statistics
- Automatic database migrations
- Environment-based configuration

## Requirements

- Go 1.25.7+
- PostgreSQL 18+
- Docker (optional, for running database)

## Setup

1. Clone the repository

2. Create a `.env` file with your database configuration:
```env
DB_HOST=localhost
DB_USER=shitreader
DB_PASS=shitreader
DB_NAME=shitreader
DB_PORT=5432
DB_SSL_MODE=disable
```

3. Start the database (using Docker):
```bash
docker-compose up -d
```

4. Place your Excel files in the `files/` directory:
   - `tabelas-dados.xlsx`
   - `cae.xlsx`
   - `paises.xlsx`
   - `distritos.xlsx`
   - `concelhos.xlsx`
   - `freguesias.xlsx`
   - `ine-zonas.xlsx`

## Running

```bash
go run cmd/main.go
```

The application will:
1. Connect to the database
2. Run migrations automatically
3. Import all Excel files sequentially
4. Display progress and statistics

## Project Structure

```
.
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── app/                 # Application setup
│   ├── services/            # Business logic
│   ├── store/               # Database layer
│   └── types/               # Data types and mappings
├── migrations/              # SQL migrations
├── files/                   # Excel data files
└── docker-compose.yml       # PostgreSQL setup
```

## License

MIT
