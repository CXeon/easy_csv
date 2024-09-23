package csv

import "testing"

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

	data, err := parseInterface(tb, true)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("有表头：%v\n", data)

	data, err = parseInterface(tb, false)
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

	data, err := parseInterfaceList(list, true)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("有表头：%v\n", data)

	data, err = parseInterfaceList(list, false)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("无表头：%v\n", data)
}
