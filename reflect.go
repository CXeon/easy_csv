package easy_csv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// marshal a structure object or a pointer to two-dimensional slice
// you can set column name in csv file by using struct tag,if not , struct field name will be column name
// for example :
//
//	type bean struct{
//		Email string  `csv:"邮箱,email_desensitization"`
//		Phone string  `csv:"手机号,phone_desensitization"`
//	}
//
// 'email_desensitization' will desensitize the email ,and  'phone_desensitization' will desensitize the phone number
//
// bean interface{}:  a structure or a pointer of structure，every field of struct should be convertible to type string
//
// setTitle bool:  true: the index 0 in result will be set column name
func marshalStructure(bean interface{}, setTitle bool) ([][]string, error) {

	if bean == nil {
		return [][]string{}, nil
	}

	reflectValue := reflect.ValueOf(bean)
	reflectType := reflectValue.Type()

	if reflectType.Kind() != reflect.Pointer && reflectType.Kind() != reflect.Struct {
		return [][]string{}, errors.New("bean must be a struct or a pointer to struct")
	}

	if reflectType.Kind() == reflect.Pointer {
		if reflectType.Elem().Kind() != reflect.Struct {
			return [][]string{}, errors.New("bean must be a pointer to struct")
		}
	}

	//strType := reflect.TypeOf("")
	if reflectValue.Kind() == reflect.Pointer {
		reflectValue = reflectValue.Elem()
		reflectType = reflectType.Elem()
	}

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
		//	return [][]string{}, errors.New("All structure field should be convertible to type string")
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

// marshal a list to two-dimensional slice,the item of list should be a structure or a pointer of structure
//
// list interface{}: any item of list should be a structure or a pointer of structure，every field of struct should be convertible to type string
//
// setTitle bool : true: the index 0 in result will be set column name
func marshalList(list interface{}, setTitle bool) ([][]string, error) {
	if list == nil {
		return [][]string{}, nil
	}

	var itemTypeKind reflect.Kind

	result := make([][]string, 0)

	reflectListValue := reflect.ValueOf(list)

	if reflectListValue.Kind() != reflect.Slice && reflectListValue.Kind() != reflect.Pointer {
		return [][]string{}, errors.New("list must be a slice or a pointer to a slice")
	}

	if reflectListValue.Kind() == reflect.Pointer {
		if reflectListValue.Elem().Kind() != reflect.Slice {
			return [][]string{}, errors.New("list must be a slice or a pointer to a slice")
		}
	}

	if reflectListValue.Kind() == reflect.Pointer {
		reflectListValue = reflectListValue.Elem()
	}

	for i := 0; i < reflectListValue.Len(); i++ {
		item := reflectListValue.Index(i)
		if item.Kind() == reflect.Pointer {
			item = item.Elem()
		}
		bean := item.Interface()

		reflectBeanValue := reflect.ValueOf(bean)
		reflectBeanType := reflectBeanValue.Type()

		if i == 0 {
			itemTypeKind = reflectBeanType.Kind()
		}

		if i > 0 && reflectBeanType.Kind() != itemTypeKind {
			return [][]string{}, errors.New("all list item type must is same")
		}

		if reflectBeanType.Kind() != reflect.Struct {
			return [][]string{}, errors.New("item must be structure or structure pointer")
		}

		if i == 0 {
			rows, err := marshalStructure(bean, setTitle)
			if err != nil {
				return [][]string{}, errors.New(fmt.Sprintf("err:%+v ,invalid row data: %v", err, bean))
			}

			result = append(result, rows...)
		}

		if i > 0 {
			rows, err := marshalStructure(bean, false)
			if err != nil {
				return [][]string{}, errors.New(fmt.Sprintf("err:%+v ,invalid row data: %v", err, bean))
			}

			result = append(result, rows...)
		}
	}

	return result, nil
}

// unmarshal a one-dimensional slice to a pointer of structure.
// The data of slice will be parsed in the order in which the structure is stored.
//
// source []string: a one-dimensional slice
//
// target interface{}: a pointer of structure
func unmarshalOneDSlice(source []string, target interface{}) error {

	if target == nil {
		return errors.New("target cannot be nil")
	}
	if len(source) == 0 {
		return nil
	}

	reflectValue := reflect.ValueOf(target)
	if reflectValue.Kind() != reflect.Pointer {
		return errors.New("target must be a pointer of structure")
	}
	sourceLen := len(source)
	reflectValue = reflectValue.Elem()
	if reflectValue.Kind() != reflect.Struct {
		return errors.New("target must be a pointer of structure")
	}
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

// unmarshal a two-dimensional slice to a list,the list must be a pointer of list,and the item of list must be a structure or a pointer of a structure.
//
// source [][]string:a two-dimensional slice pointer  of a list
//
// target interface{}: a pointer of a list
func unmarshalTwoDSlice(source [][]string, target interface{}) error {
	if target == nil {
		return errors.New("target can't is nil")
	}

	if len(source) == 0 {
		return nil
	}

	reflectPtrValue := reflect.ValueOf(target)
	if reflectPtrValue.Kind() != reflect.Pointer {
		return errors.New("target must be a structure slice ptr")
	}

	reflectSliValue := reflectPtrValue.Elem()

	reflectSliType := reflectSliValue.Type()

	if reflectSliType.Kind() != reflect.Slice {
		return errors.New("target must be a structure slice ptr")
	}

	itemReflectType := reflectSliType.Elem()

	for _, row := range source {

		var subTarget reflect.Value
		if itemReflectType.Kind() == reflect.Pointer {
			subTarget = reflect.New(itemReflectType.Elem())
		}
		if itemReflectType.Kind() == reflect.Struct {
			subTarget = reflect.New(itemReflectType)
		}

		if subTarget.Kind() == reflect.Invalid {
			return errors.New("target must be a structure slice ptr")
		}

		var err error

		err = unmarshalOneDSlice(row, subTarget.Interface())

		if err != nil {
			return err
		}

		//判断target的每个元素是结构体指针还是结构体
		if itemReflectType.Kind() == reflect.Ptr {
			reflectSliValue = reflect.Append(reflectSliValue, subTarget)
		}
		if itemReflectType.Kind() == reflect.Struct {
			reflectSliValue = reflect.Append(reflectSliValue, subTarget.Elem())
		}
	}
	reflectPtrValue.Elem().Set(reflectSliValue)

	return nil
}

// unmarshal a one-dimensional slice to a pointer of structure,and names is used to specify the order.
// The data of slice will be parsed in the order in which the structure is stored.
//
// names []string: field names of structure ,the order of source
//
// source []string: a one-dimensional slice
//
// target interface{}: a pointer of structure
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

// unmarshal a two-dimensional slice to a list,the list must be a pointer of list,and the item of list must be a structure or a pointer of a structure.
//
// names []string: field names of structure ,the order of source
//
// source [][]string:a two-dimensional slice pointer  of a list
//
// target interface{}: a pointer of a list
func unmarshalTwoDSliceWithNames(names []string, source [][]string, target interface{}) error {

	if target == nil {
		return errors.New("target can't is nil")
	}

	if len(source) == 0 {
		return errors.New("source must have a value")
	}

	reflectPtrValue := reflect.ValueOf(target)
	if reflectPtrValue.Kind() != reflect.Ptr {
		return errors.New("target must be a structure slice ptr")
	}

	reflectSliValue := reflectPtrValue.Elem()

	reflectSliType := reflectSliValue.Type()
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

		var err error
		err = unmarshalOneDSliceWithNames(names, row, subTarget.Interface())

		if err != nil {
			return err
		}
		//判断target的每个元素是结构体指针还是结构体
		if itemReflectType.Kind() == reflect.Ptr {
			reflectSliValue = reflect.Append(reflectSliValue, subTarget)
		}
		if itemReflectType.Kind() == reflect.Struct {
			reflectSliValue = reflect.Append(reflectSliValue, subTarget.Elem())
		}
	}
	reflectPtrValue.Elem().Set(reflectSliValue)
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
