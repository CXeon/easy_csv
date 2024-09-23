package csv

import (
	"encoding/csv"
	"os"
)

type CsvClientWriter struct {
	w *csv.Writer
}

type CsvClientOption struct {
	Comma   rune // Field delimiter (cloud set to ',')
	UseCRLF bool // True to use \r\n as the line terminator
}

type CsvClientOptionFunc func(*CsvClientOption)

func NewCsvClientWriter(file *os.File, opts ...CsvClientOptionFunc) *CsvClientWriter {
	option := &CsvClientOption{}
	for _, opt := range opts {
		opt(option)
	}

	w := csv.NewWriter(file)
	if option.Comma != 0 {
		w.Comma = option.Comma
	}

	w.UseCRLF = option.UseCRLF

	return &CsvClientWriter{w: w}
}

func WithComma(comma rune) CsvClientOptionFunc {
	return func(opt *CsvClientOption) {
		opt.Comma = comma
	}
}

func WithUseCRLF(useCRLF bool) CsvClientOptionFunc {
	return func(opt *CsvClientOption) {
		opt.UseCRLF = useCRLF
	}
}

// WriteRow2File 将一行数据写入文件中
func (w *CsvClientWriter) WriteRow2File(data interface{}, setTitle ...bool) error {
	flag := false
	if len(setTitle) > 0 {
		flag = setTitle[0]
	}

	records, err := parseInterface(data, flag)
	if err != nil {
		return err
	}
	err = w.WriteString2File(records)
	if err != nil {
		return err
	}
	return nil
}

// WriteRows2File 将多行数据写入文件中
func (w *CsvClientWriter) WriteRows2File(list interface{}, setTitle ...bool) error {
	flag := false
	if len(setTitle) > 0 {
		flag = setTitle[0]
	}

	records, err := parseInterfaceList(list, flag)
	if err != nil {
		return err
	}
	err = w.WriteString2File(records)
	if err != nil {
		return err
	}
	return nil
}

// WriteString2File 向CSV文件中写入文本数据
//
// data 切片每个元素代表一行，每行元素还是一个切片，其中每个元素代表一列
func (w *CsvClientWriter) WriteString2File(data [][]string) error {
	err := w.w.WriteAll(data)
	if err != nil {
		return err
	}
	w.w.Flush()
	return nil
}
