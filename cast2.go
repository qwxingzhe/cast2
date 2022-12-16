package cast2

import "reflect"

func ToMap(obj interface{}) map[string]interface{} {
	return ToMapByTag(obj, "")
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
