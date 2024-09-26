package easy_csv

import (
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

func TestMarshalStructure(t *testing.T) {

	tb := testBean{
		Name:  "王五",
		Age:   12,
		Grade: "6年级",
		Score: 123.01,
		Email: "wangwu@qq.com",
		Phone: "1331111",
	}

	data, err := marshalStructure(&tb, true)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("有表头：%v\n", data)

	data, err = marshalStructure(tb, false)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("无表头：%v\n", data)
}

func TestMarshalList(t *testing.T) {
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

	data, err := marshalList(list, true)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("有表头：%v\n", data)

	data, err = marshalList(&list, false)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("无表头：%v\n", data)

	tb := testBean{
		Name:  "李四",
		Age:   11,
		Grade: "6年级",
		Score: 100.45,
		Email: "lisi@qq.com",
		Phone: "1120000000",
	}
	list2 := make([]*testBean, 1)
	list2[0] = &tb

	data, err = marshalList(list2, true)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("有表头：%v\n", data)

	data, err = marshalList(&list2, false)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("无表头：%v\n", data)
}

func TestUnmarshalOneDSlice(t *testing.T) {
	tb := testBean{}
	source := []string{
		"王五",
		"12",
		"6年级",
		"111.20",
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

func TestUnmarshalTwoDSlice(t *testing.T) {
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
	t.Logf("%p\n", &list)
	err := unmarshalTwoDSlice(source, &list)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v\n", list)
	t.Logf("%p\n", &list)
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

	t.Logf("%p\n", &list)
	err := unmarshalTwoDSliceWithNames(names, source, &list)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v\n", list)
	t.Logf("%p\n", &list)

	list2 := make([]*testBean, 0)

	err = unmarshalTwoDSliceWithNames(names, source, &list2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v\n", list)

	return
}
