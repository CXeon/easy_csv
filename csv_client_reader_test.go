package csv

import (
	"os"
	"testing"
)

type testStudentInfo2 struct {
	Score float64 `csv:"分数"`
	Email string  `csv:"邮箱,email_desensitization"`
	Phone string  `csv:"手机号,phone_desensitization"`
	Name  string  `csv:"name"`
	Age   int
	Grade string
}

func TestClientReader_ReadRowFromFile(t *testing.T) {
	file, err := os.Open("./TestCsvReader.csv")
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer file.Close()

	clientReader := NewClientReader(file)

	clientReader.Read() //第一行表头不能处理成结构体，读取第一行

	data := testStudentInfo{}
	err = clientReader.ReadRowFromFile(&data) //继续读取才是数据
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("data:%+v\n", data)
}

func TestClientReader_ReadRowFromFileWithNames(t *testing.T) {
	file, err := os.Open("./TestCsvReader2.csv")
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer file.Close()

	clientReader := NewClientReader(file)

	clientReader.Read() //第一行表头不能处理成结构体，读取第一行

	data := testStudentInfo{}
	names := []string{
		"Score",
		"Email",
		"Phone",
		"Name",
		"Age",
		"Grade",
	}
	err = clientReader.ReadRowFromFileWithNames(names, &data) //继续读取才是数据
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("data:%+v\n", data)
}

func TestClientReader_ReadRowsFromFile(t *testing.T) {
	file, err := os.Open("./TestCsvReader.csv")
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer file.Close()

	clientReader := NewClientReader(file)

	clientReader.Read() //第一行表头不能处理成结构体，读取第一行

	var list []testStudentInfo
	err = clientReader.ReadRowsFromFile(&list) //继续读取才是数据
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("data:%+v\n", list)
}

func TestClientReader_ReadRowsFromFileWithNames(t *testing.T) {
	file, err := os.Open("./TestCsvReader2.csv")
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer file.Close()

	clientReader := NewClientReader(file)

	clientReader.Read() //第一行表头不能处理成结构体，读取第一行

	names := []string{
		"Score",
		"Email",
		"Phone",
		"Name",
		"Age",
		"Grade",
	}
	var list []testStudentInfo2
	err = clientReader.ReadRowsFromFileWithNames(names, &list) //继续读取才是数据
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("data:%+v\n", list)
}
