package csv

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// 解析结构体的数据到切片
// bean interface{} 要求必须是结构体类型，不能是指针，每一个field都可以转换成字符串类型
// setTitle bool 是否在返回的数据中第一行添加表头数据
func parseInterface(bean interface{}, setTitle bool) ([][]string, error) {

	if bean == nil {
		return [][]string{}, nil
	}

	reflectValue := reflect.ValueOf(bean)
	reflectType := reflectValue.Type()

	//确认bean是结构体类型
	if reflectType.Kind() != reflect.Struct {
		return [][]string{}, errors.New("param type interface  must is a structure ")
	}

	//strType := reflect.TypeOf("")

	numField := reflectValue.NumField()

	//行数据
	rowData := make([]string, numField)
	//表头
	title := make([]string, numField)

	//结构体每一个参数必须可以转换成字符串
	for i := 0; i < numField; i++ {
		field := reflectValue.Field(i)
		fieldType := reflectType.Field(i)
		//fmt.Println(field.String())
		//fmt.Println(strType)
		//if !field.CanConvert(strType) {
		//	return [][]string{}, errors.New("all structure field must can convert to type string")
		//}

		rowData[i] = fmt.Sprint(field.Interface())

		tagStr := fieldType.Tag.Get("csv")
		if len(tagStr) == 0 {
			if setTitle {
				title[i] = reflectType.Field(i).Name
			}
			continue
		}

		tagStrSpl := strings.Split(tagStr, ",")
		if setTitle {
			title[i] = tagStrSpl[0]
		}

		if len(tagStrSpl) > 1 {
			format := tagStrSpl[1]

			if format == "phone_desensitization" {
				//手机号脱敏
				phone := []rune(rowData[i])

				if len(phone) > 6 {
					rowData[i] = string(phone[0:3]) + "****" + string(phone[7:])
				}
			}

			if format == "email_desensitization" {
				//邮箱脱敏
				email := rowData[i]

				emailSpl := strings.Split(email, "@")
				if len(emailSpl) == 2 {
					user := []rune(emailSpl[0])
					domain := emailSpl[1]

					if len(user) > 3 {
						userStr := string(user[0:2]) + "***" + string(user[len(user)-1:])
						rowData[i] = userStr + "@" + domain
					}
				}
			}
		}

	}

	if !setTitle {
		return [][]string{rowData}, nil
	}

	return [][]string{title, rowData}, nil
}

// 解析结构体列表到二维切片
// list interface{} 要求必须是同类型结构体列表，数组里面可以放结构体指针，结构体每一个field都可以转换成字符串类型
// setTitle bool 是否在返回的数据中第一行添加表头数据
func parseInterfaceList(list interface{}, setTitle bool) ([][]string, error) {
	if list == nil {
		return [][]string{}, nil
	}

	var itemTypeKind reflect.Kind

	result := make([][]string, 0)

	reflectListValue := reflect.ValueOf(list)
	//	reflectListType := reflectListValue.Type()

	if reflectListValue.Kind() != reflect.Slice {
		return [][]string{}, errors.New("list  type must is slice")
	}

	for i := 0; i < reflectListValue.Len(); i++ {
		bean := reflectListValue.Index(i).Interface()

		reflectBeanValue := reflect.ValueOf(bean)
		reflectBeanType := reflectBeanValue.Type()

		if i == 0 {
			itemTypeKind = reflectBeanType.Kind()
		}

		if i > 0 && reflectBeanType.Kind() != itemTypeKind {
			return [][]string{}, errors.New("all list item type must is same")
		}

		if reflectBeanType.Kind() != reflect.Struct && reflectBeanType.Kind() != reflect.Pointer {
			return [][]string{}, errors.New("item type must is structure or structure pointer")
		}

		if i == 0 {
			rows, err := parseInterface(bean, setTitle)
			if err != nil {
				return [][]string{}, errors.New(fmt.Sprintf("err:%+v ,invalid row data: %v", err, bean))
			}

			result = append(result, rows...)
		}

		if i > 0 {
			rows, err := parseInterface(bean, false)
			if err != nil {
				return [][]string{}, errors.New(fmt.Sprintf("err:%+v ,invalid row data: %v", err, bean))
			}

			result = append(result, rows...)
		}
	}

	return result, nil
}
