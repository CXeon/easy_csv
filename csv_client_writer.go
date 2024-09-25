package csv

import (
	"encoding/csv"
	"io"
)

type ClientWriter struct {
	w *csv.Writer
}

type ClientWriterOption struct {
	Comma   rune // Field delimiter (cloud set to ',')
	UseCRLF bool // True to use \r\n as the line terminator
}

type ClientWriterOptionFunc func(*ClientWriterOption)

func NewCsvClientWriter(writer io.Writer, opts ...ClientWriterOptionFunc) *ClientWriter {
	option := &ClientWriterOption{}
	for _, opt := range opts {
		opt(option)
	}

	w := csv.NewWriter(writer)
	if option.Comma != 0 {
		w.Comma = option.Comma
	}

	w.UseCRLF = option.UseCRLF

	return &ClientWriter{w: w}
}

func WithWriterComma(comma rune) ClientWriterOptionFunc {
	return func(opt *ClientWriterOption) {
		opt.Comma = comma
	}
}

func WithWriterUseCRLF(useCRLF bool) ClientWriterOptionFunc {
	return func(opt *ClientWriterOption) {
		opt.UseCRLF = useCRLF
	}
}

// WriteRow2File 将一行数据写入文件中
func (writer *ClientWriter) WriteRow2File(data interface{}, setTitle ...bool) error {
	flag := false
	if len(setTitle) > 0 {
		flag = setTitle[0]
	}

	records, err := marshalStructure(data, flag)
	if err != nil {
		return err
	}
	err = writer.WriteString2File(records)
	if err != nil {
		return err
	}
	return nil
}

// WriteRows2File 将多行数据写入文件中
func (writer *ClientWriter) WriteRows2File(list interface{}, setTitle ...bool) error {
	flag := false
	if len(setTitle) > 0 {
		flag = setTitle[0]
	}

	records, err := marshalList(list, flag)
	if err != nil {
		return err
	}
	err = writer.WriteString2File(records)
	if err != nil {
		return err
	}
	return nil
}

// WriteString2File 向CSV文件中写入文本数据
//
// data 切片每个元素代表一行，每行元素还是一个切片，其中每个元素代表一列
func (writer *ClientWriter) WriteString2File(data [][]string) error {
	err := writer.w.WriteAll(data)
	if err != nil {
		return err
	}
	writer.w.Flush()
	return nil
}
