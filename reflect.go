package csv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// 解析结构体的数据到切片
// bean interface{} 要求必须是结构体类型，不能是指针，每一个field都可以转换成字符串类型
// setTitle bool 是否在返回的数据中第一行添加表头数据
func marshalInterface(bean interface{}, setTitle bool) ([][]string, error) {

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

		//if !field.CanConvert(strType) {
		//	return [][]string{}, errors.New("All structure field should be convertible to strings")
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
// list interface{} 要求必须是同类型结构体列表，结构体每一个field都可以转换成字符串类型
// setTitle bool 是否在返回的数据中第一行添加表头数据
func marshalInterfaceList(list interface{}, setTitle bool) ([][]string, error) {
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
			rows, err := marshalInterface(bean, setTitle)
			if err != nil {
				return [][]string{}, errors.New(fmt.Sprintf("err:%+v ,invalid row data: %v", err, bean))
			}

			result = append(result, rows...)
		}

		if i > 0 {
			rows, err := marshalInterface(bean, false)
			if err != nil {
				return [][]string{}, errors.New(fmt.Sprintf("err:%+v ,invalid row data: %v", err, bean))
			}

			result = append(result, rows...)
		}
	}

	return result, nil
}

// 将字符串数组转换成结构体
//
// source 需要处理的数据
// target 结构体对象指针
func unmarshalOneDSlice(source []string, target interface{}) error {

	reflectValue := reflect.ValueOf(target)
	fmt.Println(reflectValue.Kind().String())
	if reflectValue.Kind() != reflect.Ptr {
		return errors.New("target must be a structure ptr")
	}
	sourceLen := len(source)
	reflectValue = reflectValue.Elem()
	fmt.Println(reflectValue)
	fmt.Println(reflectValue.Kind().String())
	fieldNum := reflectValue.NumField()

	for i := 0; i < fieldNum; i++ {
		field := reflectValue.Field(i)

		if i < sourceLen {
			s := source[i]
			err := setFieldValue(field, s)
			if err != nil {
				return err
			}
		} else {
			break
		}
	}

	return nil
}

// 将多行数解析到结构体列表
//
// source 多行数据
// target list指针 list的元素可以是结构体，也可以是结构体指针
func unmarshalTwoDSlice(source [][]string, target interface{}) error {
	if target == nil {
		return errors.New("target can't is nil")
	}

	reflectPtrValue := reflect.ValueOf(target)
	fmt.Println(reflectPtrValue.Kind().String())
	if reflectPtrValue.Kind() != reflect.Ptr {
		return errors.New("target must be a structure slice ptr")
	}

	reflectSliValue := reflectPtrValue.Elem()

	reflectSliType := reflectSliValue.Type()
	fmt.Println(reflectSliType.Kind().String())
	if reflectSliType.Kind() != reflect.Slice {
		return errors.New("target must be a structure slice ptr")
	}

	fmt.Println(reflectSliType.Elem())
	fmt.Println(reflectSliType.Elem().Kind())
	fmt.Println(reflectSliType.Elem().Kind().String())
	itemReflectType := reflectSliType.Elem()
	for _, row := range source {

		var subTarget reflect.Value
		if itemReflectType.Kind() == reflect.Ptr {
			subTarget = reflect.New(itemReflectType.Elem())
		}
		if itemReflectType.Kind() == reflect.Struct {
			subTarget = reflect.New(itemReflectType)
		}

		if subTarget.Kind() == reflect.Invalid {
			return errors.New("target must be a structure slice ptr")
		}
		fmt.Println(subTarget.Kind())
		var err error

		err = unmarshalOneDSlice(row, subTarget.Interface())

		if err != nil {
			return err
		}
		fmt.Println(subTarget.Interface())
		//判断target的每个元素是结构体指针还是结构体
		if itemReflectType.Kind() == reflect.Ptr {
			reflectSliValue = reflect.Append(reflectSliValue, subTarget)
		}
		if itemReflectType.Kind() == reflect.Struct {
			reflectSliValue = reflect.Append(reflectSliValue, subTarget.Elem())
		}
	}
	reflectPtrValue.Elem().Set(reflectSliValue)
	fmt.Printf("%+v\n", reflectSliValue.Interface())
	return nil
}

// 将字符串数组数据解析到结构体
// names 结构体field名称
// source 需要处理的数据
// target 结构体对象指针
func unmarshalOneDSliceWithNames(names []string, source []string, target interface{}) error {

	if target == nil {
		return errors.New("target can't is nil")
	}

	if len(names) == 0 || len(source) == 0 {
		return errors.New("titles and source must have a value")
	}
	if len(names) != len(source) {
		return errors.New("titles and source must have same length")
	}

	reflectValue := reflect.ValueOf(target)

	if reflectValue.Kind() != reflect.Ptr {
		return errors.New("target must be a structure slice ptr")
	}

	if reflectValue.Elem().Kind() != reflect.Struct {
		return errors.New("target must be a structure slice ptr")
	}

	reflectValue = reflectValue.Elem()

	for i := 0; i < len(names); i++ {
		t := names[i]
		s := source[i]

		field := reflectValue.FieldByName(t)
		if field.Kind() != reflect.Invalid {
			err := setFieldValue(field, s)
			if err != nil {
				return err
			}
		}
	}
	return nil

}

// 将多行数据解析到结构体列表
//
// names 结构体field名称
// source 多行数据
// target list指针 list的元素可以是结构体，也可以是结构体指针
func unmarshalTwoDSliceWithNames(names []string, source [][]string, target interface{}) error {

	if target == nil {
		return errors.New("target can't is nil")
	}

	if len(source) == 0 {
		return errors.New("source must have a value")
	}

	reflectPtrValue := reflect.ValueOf(target)
	fmt.Println(reflectPtrValue.Kind().String())
	if reflectPtrValue.Kind() != reflect.Ptr {
		return errors.New("target must be a structure slice ptr")
	}

	reflectSliValue := reflectPtrValue.Elem()

	reflectSliType := reflectSliValue.Type()
	fmt.Println(reflectSliType.Kind().String())
	if reflectSliType.Kind() != reflect.Slice {
		return errors.New("target must be a structure slice ptr")
	}

	itemReflectType := reflectSliType.Elem()
	for _, row := range source {
		var subTarget reflect.Value
		if itemReflectType.Kind() == reflect.Ptr {
			subTarget = reflect.New(itemReflectType.Elem())
		}
		if itemReflectType.Kind() == reflect.Struct {
			subTarget = reflect.New(itemReflectType)
		}

		if subTarget.Kind() == reflect.Invalid {
			return errors.New("target must be a structure slice ptr")
		}

		fmt.Println(subTarget.Kind())
		var err error

		err = unmarshalOneDSliceWithNames(names, row, subTarget.Interface())

		if err != nil {
			return err
		}
		fmt.Println(subTarget.Interface())
		//判断target的每个元素是结构体指针还是结构体
		if itemReflectType.Kind() == reflect.Ptr {
			reflectSliValue = reflect.Append(reflectSliValue, subTarget)
		}
		if itemReflectType.Kind() == reflect.Struct {
			reflectSliValue = reflect.Append(reflectSliValue, subTarget.Elem())
		}
	}
	reflectPtrValue.Elem().Set(reflectSliValue)
	fmt.Printf("%+v\n", reflectSliValue.Interface())
	return nil
}

func setFieldValue(field reflect.Value, str string) error {
	t := field.Type()
	switch t.Kind() {
	case reflect.String:
		field.SetString(str)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		field.SetFloat(v)
	case reflect.Bool:
		switch strings.ToLower(str) {
		case "true":
			field.SetBool(true)
		default:
			field.SetBool(false)
		}
	default:
		return errors.New(fmt.Sprintf("unsupported type %s to convert %s", t.Kind().String(), str))
	}
	return nil
}
