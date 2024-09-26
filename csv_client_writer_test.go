package easy_csv

import (
	"os"
	"testing"
)

type testStudentInfo struct {
	Name  string `csv:"name"`
	Age   int
	Grade string
	Score float64 `csv:"分数"`
	Email string  `csv:"邮箱,email_desensitization"`
	Phone string  `csv:"手机号,phone_desensitization"`
}

func TestCsvClientWriter(t *testing.T) {

	//打开文件
	fileName := "./TestCsvWriter.csv"

	csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Error(err)
		return
	}
	defer csvFile.Close()

	//多行数据
	list := []testStudentInfo{
		{
			Name:  "王五",
			Age:   12,
			Grade: "6年级",
			Score: 123.01,
			Email: "wangwu@qq.com",
			Phone: "133111",
		},
		{
			Name:  "张三",
			Age:   11,
			Grade: "5年级",
			Score: 123.02,
			Email: "zh@qq.com",
			Phone: "13322225559",
		},
	}

	//一行数据
	row := testStudentInfo{
		Name:  "李四",
		Age:   10,
		Grade: "4年级",
		Score: 99.01,
		Email: "lisi@qq.com",
		Phone: "13322226666",
	}

	writer := NewClientWriter(csvFile)
	err = writer.WriteRows2File(list, true)
	if err != nil {
		t.Error(err)
		return
	}

	err = writer.WriteRow2File(row)
	if err != nil {
		t.Error(err)
		return
	}

	return
}
