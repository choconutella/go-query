# Go SQL Recordset

A lightweight and efficient Go package that provides a flexible wrapper for SQL query results, offering both array and map-based data access patterns.

## Features

- Type-safe scanning of SQL query results
- Support for NULL value handling using `sql.Null*` types
- Multiple output formats:
  - Array-based access (`[][]interface{}`)
  - Map-based access with column names as keys (`[]map[string]interface{}`)
- Automatic type detection and conversion
- Memory-efficient data handling
- Comprehensive error handling

## Installation

```bash
go get github.com/choconutella/recordset
```

## Usage

### Basic Example

```go
package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq" // or your preferred database driver
    "github.com/choconutella/recordset"
)

func main() {
    // Open database connection
    db, err := sql.Open("postgres", "postgres://username:password@localhost/dbname?sslmode=disable")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Execute query
    rows, err := db.Query("SELECT id, name, age FROM users")
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    // Create new recordset
    rs := recordset.NewRecordset(rows)

    // Get results as array
    results, err := rs.Query()
    if err != nil {
        panic(err)
    }

    // Process results
    for _, row := range results {
        fmt.Printf("Row: %v\n", row)
    }
}
```

### Map-Based Access

```go
// Get results as map
mapResults, err := rs.QueryAsMap()
if err != nil {
    panic(err)
}

// Access data using column names
for _, row := range mapResults {
    fmt.Printf("Name: %v, Age: %v\n", row["name"], row["age"])
}
```

## Supported Data Types

The package automatically handles the following SQL data types:

- String types: VARCHAR, TEXT, CHAR, UUID
- Numeric types: INT, BIGINT, SMALLINT, FLOAT, DOUBLE, DECIMAL
- Boolean: BOOL
- Date/Time: TIMESTAMP, DATETIME, DATE
- Other types are handled as generic interfaces

## NULL Value Handling

The package properly handles NULL values using Go's `sql.Null*` types:

- `sql.NullString` for string types
- `sql.NullInt64` for integer types
- `sql.NullFloat64` for floating-point types
- `sql.NullBool` for boolean types
- `sql.NullTime` for date/time types

NULL values are converted to nil in the output data structures.

## Memory Efficiency

The package implements several optimizations for memory efficiency:

- Pre-allocated slices with initial capacity
- New row data allocation for each iteration to prevent data overwrites
- Efficient type conversions
- Proper memory management for large result sets

## Error Handling

Comprehensive error handling is implemented throughout the package:

- Nil checks for database rows
- Column name and type retrieval errors
- Row scanning errors
- Row iteration errors

All errors include detailed messages to help with debugging.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

Adhil Novandri

## Support

If you encounter any problems or have suggestions, please open an issue in the GitHub repository.
