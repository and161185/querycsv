# QueryCSV

QueryCSV is a command-line tool written in Go that allows users to execute SQL-like queries on CSV files. It enables quick and efficient data filtering, aggregation, and transformation without requiring a dedicated database.

## Features

- Read CSV files and execute SQL-like queries.
- Support for basic SQL operations: `SELECT`, `WHERE`, `ORDER BY`, `GROUP BY`.
- Fast data processing using Go's concurrency capabilities.
- Simple command-line interface.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/and161185/querycsv.git
   cd querycsv
   ```
2. Build the application:
   ```bash
   go build -o querycsv .
   ```

## Usage

Run the tool with a CSV file and query as arguments:
```bash
./querycsv -file data.csv -query "SELECT name, age FROM data WHERE age > 30 ORDER BY age DESC"
```

### Command-line Flags
- `-file`: Path to the CSV file.
- `-query`: SQL-like query string.
- `-delimiter`: Specify the CSV delimiter (default: `,`).
- `-header`: Indicate whether the CSV has a header row (default: `true`).

## Configuration

Configuration options are defined in `config/config.json`. An example configuration:
```json
{
    "delimiter": ",",
    "header": true
}
```

## Testing

Run unit tests with:
```bash
go test ./...
```

## Dependencies

- [Go](https://golang.org/)
- [github.com/xwb1989/sqlparser](https://github.com/xwb1989/sqlparser) for SQL query parsing.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

