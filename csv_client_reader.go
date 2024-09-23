package csv

import "encoding/csv"

type CsvClientReader struct {
	r *csv.Reader
}
