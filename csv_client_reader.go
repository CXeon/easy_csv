package easycsv

import (
	"encoding/csv"
	"io"
)

// ClientReader a reader client is used to read and unmarshal file of csv
type ClientReader struct {
	r *csv.Reader
}

type ClientReaderOption struct {
	// Comma is the field delimiter.
	// It is set to comma (',') by NewReader.
	// Comma must be a valid rune and must not be \r, \n,
	// or the Unicode replacement character (0xFFFD).
	Comma rune

	// Comment, if not 0, is the comment character. Lines beginning with the
	// Comment character without preceding whitespace are ignored.
	// With leading whitespace the Comment character becomes part of the
	// field, even if TrimLeadingSpace is true.
	// Comment must be a valid rune and must not be \r, \n,
	// or the Unicode replacement character (0xFFFD).
	// It must also not be equal to Comma.
	Comment rune

	// FieldsPerRecord is the number of expected fields per record.
	// If FieldsPerRecord is positive, Read requires each record to
	// have the given number of fields. If FieldsPerRecord is 0, Read sets it to
	// the number of fields in the first record, so that future records must
	// have the same field count. If FieldsPerRecord is negative, no check is
	// made and records may have a variable number of fields.
	FieldsPerRecord int

	// If LazyQuotes is true, a quote may appear in an unquoted field and a
	// non-doubled quote may appear in a quoted field.
	LazyQuotes bool

	// If TrimLeadingSpace is true, leading white space in a field is ignored.
	// This is done even if the field delimiter, Comma, is white space.
	TrimLeadingSpace bool

	// ReuseRecord controls whether calls to Read may return a slice sharing
	// the backing array of the previous call's returned slice for performance.
	// By default, each call to Read returns newly allocated memory owned by the caller.
	ReuseRecord bool
}

type ClientReaderOptionFunc func(opt *ClientReaderOption)

func NewClientReader(reader io.Reader, opts ...ClientReaderOptionFunc) *ClientReader {

	option := &ClientReaderOption{}
	for _, o := range opts {
		o(option)
	}

	r := csv.NewReader(reader)
	if option.Comma != 0 {
		r.Comma = option.Comma
	}
	if option.Comment != 0 {
		r.Comment = option.Comment
	}

	r.FieldsPerRecord = option.FieldsPerRecord
	r.LazyQuotes = option.LazyQuotes
	r.TrimLeadingSpace = option.TrimLeadingSpace
	r.ReuseRecord = option.ReuseRecord

	return &ClientReader{
		r: r,
	}
}

func WithReaderComma(comma rune) ClientReaderOptionFunc {
	return func(opt *ClientReaderOption) {
		opt.Comma = comma
	}
}

func WithReaderComment(comment rune) ClientReaderOptionFunc {
	return func(opt *ClientReaderOption) {
		opt.Comment = comment
	}
}

func WithReaderFieldsPerRecord(fieldsPerRecord int) ClientReaderOptionFunc {
	return func(opt *ClientReaderOption) {
		opt.FieldsPerRecord = fieldsPerRecord
	}
}

func WithReaderLazyQuotes(lazyQuotes bool) ClientReaderOptionFunc {
	return func(opt *ClientReaderOption) {
		opt.LazyQuotes = lazyQuotes
	}
}

func WithTrimLeadingSpace(trimLeadingSpace bool) ClientReaderOptionFunc {
	return func(opt *ClientReaderOption) {
		opt.TrimLeadingSpace = trimLeadingSpace
	}
}

func WithReuseRecord(reuse bool) ClientReaderOptionFunc {
	return func(opt *ClientReaderOption) {
		opt.ReuseRecord = reuse
	}
}

// Read Read one line at a time
func (reader *ClientReader) Read() ([]string, error) {
	return reader.r.Read()
}

// ReadAll Read all the remaining lines in sequence
func (reader *ClientReader) ReadAll() ([][]string, error) {
	return reader.r.ReadAll()
}

// ReadRowFromFile Read a row of lines and parse it into each field of the structure in the order of columns.
//
// structure: The parameter structure is a structure pointer
func (reader *ClientReader) ReadRowFromFile(structure interface{}) error {
	row, err := reader.Read()
	if err != nil {
		return err
	}

	err = unmarshalOneDSlice(row, structure)
	if err != nil {
		return err
	}
	return nil
}

// ReadRowFromFileWithNames Read a row of lines and parse the data into the corresponding field name of the structure in the order specified by names.
//
// names: The parameter names is field names of structure in order to specify column of file to structure
//
// structure: The parameter structure is a structure pointer
func (reader *ClientReader) ReadRowFromFileWithNames(names []string, structure interface{}) error {
	row, err := reader.Read()
	if err != nil {
		return err
	}

	err = unmarshalOneDSliceWithNames(names, row, structure)
	if err != nil {
		return err
	}
	return nil
}

// ReadRowsFromFile Read rows of remaining lines and parse it into each field of the structure in the order of columns.
//
// list: The parameter list is a list pointer,the item of list must be a structure or a structure pointer
func (reader *ClientReader) ReadRowsFromFile(list interface{}) error {
	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}
	err = unmarshalTwoDSlice(rows, list)
	if err != nil {
		return err
	}

	return nil
}

// ReadRowsFromFileWithNames Read rows of lines and parse the data into the corresponding field name of the structure in the order specified by names.
//
// names: The parameter names is field names of structure in order to specify column of file to structure
//
// list: The parameter list is a list pointer,the item of list must be a structure or a structure pointer
func (reader *ClientReader) ReadRowsFromFileWithNames(names []string, list interface{}) error {
	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}
	err = unmarshalTwoDSliceWithNames(names, rows, list)
	if err != nil {
		return err
	}

	return nil
}
