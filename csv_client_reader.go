package csv

import "encoding/csv"

type ClientReader struct {
	r *csv.Reader
}
