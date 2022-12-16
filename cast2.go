package cast2

import (
	"encoding/json"
	"errors"
	"github.com/golang-module/carbon"
	"github.com/spf13/cast"
	"log"
	"reflect"
)

func ToMap(obj interface{}) map[string]interface{} {
	return ToMapByTag(obj, "")
}

func ToMapByTagJson(obj interface{}) map[string]interface{} {
	return ToMapByTag(obj, "json")
}

func ToMapByTag(obj interface{}, tagName string) map[string]interface{} {
	var maps = make(map[string]interface{}, 0)

	typeof := reflect.TypeOf(obj)
	values := reflect.ValueOf(obj)
	for i := 0; i < typeof.NumField(); i++ {
		fieldName := typeof.Field(i).Name
		keyName := fieldName
		if tagName != "" {
			field := typeof.Field(i)
			keyName = field.Tag.Get(tagName)
		}

		fv := values.FieldByName(fieldName)
		if fv.CanInterface() && keyName != "" {
			fieldValue := values.FieldByName(fieldName).Interface()
			maps[keyName] = fieldValue
		}
	}

	return maps
}

// SetStructValue 按属性key给结构体赋值
func SetStructValue[T any, T2 any](data T, key string, value T2) T {
	elem := reflect.ValueOf(&data).Elem()
	fieldKind := elem.FieldByName(key).Kind()
	if fieldKind == reflect.Invalid {
		return data
	}
	fieldValue := reflect.ValueOf(value)
	elem.FieldByName(key).Set(fieldValue)
	return data
}

func GetStructKeyKind[T any](data T, key string) reflect.Kind {
	elem := reflect.ValueOf(&data).Elem()
	return elem.FieldByName(key).Kind()
}

// Unmarshal 将原字符串值转对象，不考虑错误情况
func Unmarshal[T any](oldData string) T {
	var data T
	json.Unmarshal([]byte(oldData), &data)
	return data
}

// ToString 将对象转字符串，不考虑异常情况
func ToString(v any) string {
	bt, err := json.Marshal(v)
	if err != nil {
		log.Println("ToString error:", err.Error())
		return ""
	}
	return string(bt)
}

// CopyStruct 匹配到及赋值
func CopyStruct[T1 any, T2 any](a1 T1, a2 T2) T2 {
	return CopyStructMapping[T1, T2](a1, a2, map[string]string{})
}

// CopyStructMapping 将目标字段映射赋值
func CopyStructMapping[T1 any, T2 any](a1 T1, a2 T2, mapConfig map[string]string) T2 {
	aimTypeOf := reflect.TypeOf(a2)
	aimElem := reflect.ValueOf(&a2).Elem()

	t := reflect.TypeOf(a1)
	v := reflect.ValueOf(a1)
	var discardField []string

	for k := 0; k < t.NumField(); k++ {
		sourceFieldName := t.Field(k).Name
		sourceFieldKind := t.Field(k).Type.Kind()
		fieldValue := v.Field(k)

		aimFieldName := sourceFieldName
		// 字段映射
		if fv, ok := mapConfig[aimFieldName]; ok {
			aimFieldName = fv
		}

		// 判断 a2 中是否包含该字段
		if _, ok := aimTypeOf.FieldByName(aimFieldName); ok {
			// 判断字段类型是否一致，是否需要转换跳过
			aimFieldKind := aimElem.FieldByName(aimFieldName).Kind()
			aimFieldType := aimElem.FieldByName(aimFieldName).Type().String()
			if aimFieldType == "time.Time" && sourceFieldKind == reflect.String { // 时间转换
				fieldValue = reflect.ValueOf(carbon.Parse(fieldValue.Interface().(string)).Carbon2Time())
			} else if aimFieldKind == reflect.String && (sourceFieldKind == reflect.Struct || sourceFieldKind == reflect.Array || sourceFieldKind == reflect.Slice) {
				//fmt.Println("******| ", aimFieldName)
				// a1 是结构体 或者 数组
				// a2 是字符串，自动将 a1 的值转成json 存 a2
				fvb, err := json.Marshal(fieldValue.Interface())
				if err != nil {
					//fmt.Println("discard 1 ------| ", aimFieldName, "source:", sourceFieldKind, "aim:", aimFieldKind)
					discardField = append(discardField, sourceFieldName)
					continue
				}
				fieldValue = reflect.ValueOf(string(fvb))
			} else if sourceFieldKind != aimFieldKind {
				fv, err := ReflectConversion(fieldValue, aimFieldKind)
				if err != nil {
					//fmt.Println("discard 2 ------| ", aimFieldName, "source:", sourceFieldKind, "aim:", aimFieldKind)
					discardField = append(discardField, sourceFieldName)
					continue
				} else {
					fieldValue = fv
				}
			}

			// 赋值
			aimElem.FieldByName(aimFieldName).Set(fieldValue)
		}
	}

	if len(discardField) > 0 {
		log.Println("")
		log.Println("discardField : ", discardField)
		log.Println("")
	}

	return a2
}

func ReflectConversion(sourceFieldValue reflect.Value, aimFieldKind reflect.Kind) (reflect.Value, error) {
	switch aimFieldKind {
	case reflect.Float64:
		return doReflectConversion[float64](sourceFieldValue)
	case reflect.Float32:
		return doReflectConversion[float32](sourceFieldValue)
	case reflect.String:
		return doReflectConversion[string](sourceFieldValue)
	case reflect.Int:
		return doReflectConversion[int](sourceFieldValue)
	case reflect.Int8:
		return doReflectConversion[int8](sourceFieldValue)
	case reflect.Int16:
		return doReflectConversion[int16](sourceFieldValue)
	case reflect.Int32:
		return doReflectConversion[int32](sourceFieldValue)
	case reflect.Int64:
		return doReflectConversion[int64](sourceFieldValue)
	case reflect.Uint:
		return doReflectConversion[uint](sourceFieldValue)
	case reflect.Uint8:
		return doReflectConversion[uint8](sourceFieldValue)
	case reflect.Uint16:
		return doReflectConversion[uint16](sourceFieldValue)
	case reflect.Uint32:
		return doReflectConversion[uint32](sourceFieldValue)
	case reflect.Uint64:
		return doReflectConversion[uint64](sourceFieldValue)
	}
	return sourceFieldValue, errors.New("not support! ")
}

func doReflectConversion[T any](fieldValue reflect.Value) (reflect.Value, error) {
	var v T
	typeConversion(fieldValue.Interface(), &v)
	return reflect.ValueOf(v), nil
}

// typeConversion TODO 将数值转换成目标数据类型
func typeConversion(source any, output interface{}) (err error) {

	// 获取变量a的反射值对象(a的地址)
	valueOf := reflect.ValueOf(output)
	// 取出a地址的元素(a的值)
	valueOf = valueOf.Elem()

	switch indirect(output).(type) {
	case string:
		if s, err := cast.ToStringE(source); err == nil {
			valueOf.SetString(s)
		}
	case float64:
		if s, err := cast.ToFloat64E(source); err == nil {
			valueOf.SetFloat(s)
		}
	case uint:
		if s, err := cast.ToUint64E(source); err == nil {
			valueOf.SetUint(s)
		}
	case int:
		if s, err := cast.ToInt64E(source); err == nil {
			valueOf.SetInt(s)
		}
	default:
		err = errors.New("typeConversion 类型转换缺失:" + indirect(output).(string))
	}
	return
}

func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

func ErgodicObj(obj interface{}, fn func(fieldName string)) {
	typeof := reflect.TypeOf(obj)
	for k := 0; k < typeof.NumField(); k++ {
		fn(typeof.Field(k).Name)
	}
}

func StructValue[T any](data T, field string) any {
	v := reflect.ValueOf(data)
	fv := v.FieldByName(field)
	if fv.Kind() == reflect.Invalid {
		return nil
	}
	if fv.CanInterface() {
		return fv.Interface()
	}
	return nil
}

func StructHaveField[T any](data T, field string) bool {
	elem := reflect.ValueOf(&data).Elem()
	kind := elem.FieldByName(field).Kind()
	if kind == reflect.Invalid {
		return false
	}
	return true
}
