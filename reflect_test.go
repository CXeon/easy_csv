package csv

import (
	"fmt"
	"testing"
)

type testBean struct {
	Name  string `csv:"name"`
	Age   int
	Grade string
	Score float64 `csv:"分数"`
	Email string  `csv:"邮箱,email_desensitization"`
	Phone string  `csv:"手机号,phone_desensitization"`
}

func TestParseInterface(t *testing.T) {

	tb := testBean{
		Name:  "王五",
		Age:   12,
		Grade: "6年级",
		Score: 123.01,
		Email: "wangwu@qq.com",
		Phone: "1331111",
	}

	data, err := marshalInterface(tb, true)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("有表头：%v\n", data)

	data, err = marshalInterface(tb, false)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("无表头：%v\n", data)
}

func TestParseInterfaceList(t *testing.T) {
	list := []testBean{
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

	data, err := marshalInterfaceList(list, true)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("有表头：%v\n", data)

	data, err = marshalInterfaceList(list, false)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("无表头：%v\n", data)
}

func TestUnmarshal2Interface(t *testing.T) {
	tb := testBean{}
	source := []string{
		"王五",
		"12",
		"6年级",
		"bool",
		"wangwu@qq.com",
		"1331111",
	}

	err := unmarshalOneDSlice(source, &tb)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("interface: %+v\n", tb)
	return

}

func TestUnmarshal2InterfaceList(t *testing.T) {
	list := []testBean{}
	source := [][]string{
		{
			"王五",
			"12",
			"6年级",
			"101.1",
			"wangwu@qq.com",
			"1331111",
		},
		{
			"李四",
			"13",
			"4年级",
			"105.3",
			"lisi@qq.com",
			"22222",
		},
	}
	fmt.Printf("%p\n", &list)
	err := unmarshalTwoDSlice(source, &list)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%p\n", &list)
	return
}

func TestUnmarshalOneDSliceWithNames(t *testing.T) {

	names := []string{
		"Name",
		"Phone",
		"Email",
		"Age",
		"Grade",
		"Score",
	}

	source := []string{
		"王五",
		"1331111",
		"wangwu@qq.com",
		"12",
		"6年级",
		"101.1",
	}

	tb := testBean{}

	err := unmarshalOneDSliceWithNames(names, source, &tb)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("interface: %+v\n", tb)
	return
}

func TestUnmarshalTwoDSliceWithNames(t *testing.T) {
	list := []testBean{}
	source := [][]string{
		{
			"王五",
			"1331111",
			"wangwu@qq.com",
			"12",
			"6年级",
			"101.1",
		},
		{
			"李四",
			"22222",
			"lisi@qq.com",
			"13",
			"4年级",
			"105.3",
		},
	}

	names := []string{
		"Name",
		"Phone",
		"Email",
		"Age",
		"Grade",
		"Score",
	}

	fmt.Printf("%p\n", &list)
	err := unmarshalTwoDSliceWithNames(names, source, &list)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%p\n", &list)
	return
}
