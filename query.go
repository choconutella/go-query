package query

import (
	"database/sql"
	"fmt"
)

type Recordset struct {
	Rows *sql.Rows
}

func NewRecordset(rows *sql.Rows) *Recordset {
	return &Recordset{
		Rows: rows,
	}
}

func (r *Recordset) Query() ([][]interface{}, error) {
	if r.Rows == nil {
		return nil, fmt.Errorf("rows is nil")
	}

	columns, err := r.Rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get column names: %v", err)
	}

	colTypes, err := r.Rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("failed to get column types: %v", err)
	}

	// Initialize slice with capacity to avoid frequent reallocations
	dataFrame := make([][]interface{}, 0, 10)

	for r.Rows.Next() {
		rowData := make([]interface{}, len(columns))
		rowPointers := make([]interface{}, len(columns))

		for i := range rowData {

			switch colTypes[i].DatabaseTypeName() {
			case "VARCHAR", "TEXT", "CHAR", "UUID":
				rowData[i] = new(sql.NullString)
			case "BOOL":
				rowData[i] = new(sql.NullBool)
			case "INT", "BIGINT", "SMALLINT":
				rowData[i] = new(sql.NullInt64)
			case "FLOAT", "DOUBLE", "DECIMAL":
				rowData[i] = new(sql.NullFloat64)
			case "TIMESTAMP", "DATETIME", "DATE":
				rowData[i] = new(sql.NullTime)
			default:
				rowData[i] = new(interface{})
			}
			rowPointers[i] = rowData[i]
		}

		if err := r.Rows.Scan(rowPointers...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		// Create a new row slice for actual values
		actualRow := make([]interface{}, len(columns))
		for i, data := range rowPointers {
			// Extract actual value from sql.Null* types
			switch v := data.(type) {
			case *sql.NullString:
				if v.Valid {
					actualRow[i] = v.String
				}
			case *sql.NullInt64:
				if v.Valid {
					actualRow[i] = v.Int64
				}
			case *sql.NullFloat64:
				if v.Valid {
					actualRow[i] = v.Float64
				}
			case *sql.NullBool:
				if v.Valid {
					actualRow[i] = v.Bool
				}
			case *sql.NullTime:
				if v.Valid {
					actualRow[i] = v.Time
				}
			default:
				actualRow[i] = data
			}
		}

		dataFrame = append(dataFrame, actualRow)
	}

	// Check for errors that occurred during iteration
	if err = r.Rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return dataFrame, nil
}

func (r *Recordset) QueryAsMap() ([]map[string]interface{}, error) {
	columns, err := r.Rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get column names: %v", err)
	}

	// Create a slice to hold the row data
	rowData := make([]interface{}, len(columns))
	rowPointers := make([]interface{}, len(columns))
	for i := range rowData {
		rowPointers[i] = &rowData[i]
	}

	// Create slice to hold all rows
	var result []map[string]interface{}

	// Iterate through rows
	for r.Rows.Next() {
		if err := r.Rows.Scan(rowPointers...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		// Create a map for current row
		row := make(map[string]interface{})

		// Fill the map with column names and values
		for i, col := range columns {
			if rowData[i] != nil {
				row[col] = rowData[i]
			}
		}

		result = append(result, row)
	}

	if err = r.Rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return result, nil
}
