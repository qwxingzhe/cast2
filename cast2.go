package cast2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-module/carbon"
	"github.com/spf13/cast"
	"log"
	"reflect"
	"sort"
)

// ToListMap 将list转换成指定key为下标的map
func ToListMap[TKey comparable, T any](list []T, keyName string) map[TKey]T {
	data := make(map[TKey]T)
	for _, item := range list {
		kv := StructValue(item, keyName).(TKey)
		data[kv] = item
	}
	return data
}

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
func CopyStruct[T1 any, T2 any](original T1, aim T2) T2 {
	return CopyStructAdv[T1, T2](original, aim, FieldConversionConfig{})
}

type FieldConversionConfig struct {
	PartialConversionFields []string          //需要进行部分转换的字段，未设置则进行全部转换（源字段）
	ReplaceField            map[string]string //需要进行的字段替换配置（源字段:新字段）
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

// CopyStructAdv 将目标字段映射赋值
func CopyStructAdv[T1 any, T2 any](original T1, aim T2, c FieldConversionConfig) T2 {

	mapConfig := c.ReplaceField

	aimTypeOf := reflect.TypeOf(aim)
	aimElem := reflect.ValueOf(&aim).Elem()

	t := reflect.TypeOf(original)
	v := reflect.ValueOf(original)
	var discardField []string

	partialConversionFieldsCount := len(c.PartialConversionFields)
	sort.Strings(c.PartialConversionFields)

	for k := 0; k < t.NumField(); k++ {
		sourceFieldName := t.Field(k).Name
		//fmt.Println(" ---||| ", sourceFieldName, partialConversionFieldsCount, InStringsSorted(sourceFieldName, c.PartialConversionFields))
		if partialConversionFieldsCount > 0 && !InStringsSorted(sourceFieldName, c.PartialConversionFields) {
			continue
		}
		sourceFieldKind := t.Field(k).Type.Kind()
		fieldValue := v.Field(k)

		aimFieldName := sourceFieldName
		// 字段映射
		if fv, ok := mapConfig[aimFieldName]; ok {
			aimFieldName = fv
		}

		// 判断 aim 中是否包含该字段
		if _, ok := aimTypeOf.FieldByName(aimFieldName); ok {
			// 判断字段类型是否一致，是否需要转换跳过
			aimFieldKind := aimElem.FieldByName(aimFieldName).Kind()
			aimFieldType := aimElem.FieldByName(aimFieldName).Type().String()
			if aimFieldType == "time.Time" && sourceFieldKind == reflect.String { // 时间转换
				fieldValue = reflect.ValueOf(carbon.Parse(fieldValue.Interface().(string)).Carbon2Time())
			} else if aimFieldKind == reflect.String && (sourceFieldKind == reflect.Struct || sourceFieldKind == reflect.Array || sourceFieldKind == reflect.Slice) {
				//fmt.Println("******| ", aimFieldName)
				// original 是结构体 或者 数组
				// aim 是字符串，自动将 original 的值转成json 存 aim
				fvb, err := json.Marshal(fieldValue.Interface())
				if err != nil {
					//fmt.Println("discard 1 ------| ", aimFieldName, "source:", sourceFieldKind, "aim:", aimFieldKind)
					discardField = append(discardField, sourceFieldName)
					continue
				}
				fieldValue = reflect.ValueOf(string(fvb))
			} else if sourceFieldKind != aimFieldKind {
				fv, err := reflectConversion(fieldValue, aimFieldKind)
				if err != nil {
					//fmt.Println("discard 2 ------| ", aimFieldName, "source:", sourceFieldKind, "aim:", aimFieldKind)
					discardField = append(discardField, sourceFieldName)
					continue
				} else {

					fieldValue = fv
				}
			}
			//fmt.Println("赋值 ------| ", aimFieldName, "source:", sourceFieldKind, "aim:", aimFieldKind, "fieldValue:", fieldValue)
			// 赋值
			aimElem.FieldByName(aimFieldName).Set(fieldValue)
		} else {
			fmt.Println("目标字段不存在 ------| ", aimFieldName, "source:", sourceFieldKind, "fieldValue:", fieldValue)
		}
	}

	if len(discardField) > 0 {
		log.Println("")
		log.Println("discardField : ", discardField)
		log.Println("")
	}

	return aim
}

func To[T any](sourceValue any) (aim T, err error) {
	vo := reflect.ValueOf(sourceValue)
	v, err := reflectConversion(vo, reflect.TypeOf(aim).Kind())
	if err != nil {
		return
	}
	aim = v.Interface().(T)
	return
}

func GetColumn[TField any, TList any](list []TList, key string) []TField {
	var vl []TField
	for _, i2 := range list {
		vl = append(vl, GetStructValue(i2, key).(TField))
	}
	return vl
}

// CreateList 使用一个list创建另外一个list
func CreateList[T1 any, T2 any](sourceList []T2, c FieldConversionConfig) []T1 {
	var aimList []T1
	for _, item := range sourceList {
		var aim T1
		v := CopyStructAdv(item, aim, c)
		aimList = append(aimList, v)
	}
	return aimList
}

func LoadList[T1 any, T2 any](baseList []T1, baseListKey string, newList []T2, newListKey string, c FieldConversionConfig) []T1 {

	// 将 interviewPersonalList 转换成以 InterviewId 为键的map
	var mapNewList = map[any]T2{}
	for _, item := range newList {
		kv := GetStructValue(item, newListKey)
		mapNewList[kv] = item
	}

	// 再循环 list 将匹配到的写入
	for i, item := range baseList {
		kv := GetStructValue(item, baseListKey)
		if v, ok := mapNewList[kv]; ok {
			// 此处需要做字段赋值控制处理
			item = CopyStructAdv(v, item, c)
			baseList[i] = item
		}
	}

	return baseList
}

func reflectConversion(sourceFieldValue reflect.Value, aimFieldKind reflect.Kind) (reflect.Value, error) {
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
	case float32, float64:
		if s, err := cast.ToFloat64E(source); err == nil {
			valueOf.SetFloat(s)
		}
	case uint, uint8, uint16, uint32, uint64:
		if s, err := cast.ToUint64E(source); err == nil {
			valueOf.SetUint(s)
		}
	case int, int8, int16, int32, int64:
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

func GetStructValue[T any](data T, field string) any {
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
