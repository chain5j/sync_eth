// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/3/11
package reflectutil

import (
	"log"
	"reflect"
)

func ToPointer(w interface{}) interface{} {
	typeOf := reflect.TypeOf(w)
	if typeOf.Kind() == reflect.Ptr {
		return w
	}
	//reflect.PtrTo(typeOf)
	valueOf := reflect.ValueOf(w)
	pv := reflect.New(typeOf)
	pv.Elem().Set(valueOf)
	return pv.Interface()
}

func DelPointer(w interface{}) interface{} {
	typeOf := reflect.TypeOf(w)
	if typeOf.Kind() != reflect.Ptr {
		return w
	}
	valueOf := reflect.ValueOf(w)
	return valueOf.Elem().Interface()
}

func GetFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

func GetValueByFieldName(structName interface{}, fieldName string) interface{} {
	t := reflect.ValueOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldByName := t.FieldByName(fieldName)
	return fieldByName.Interface()
}

func GetTagName(structName interface{}, tagKey string) []FieldInfo {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	fieldInfos := make([]FieldInfo, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		fieldName := t.Field(i).Name
		fieldType := t.Field(i).Type.Name()

		fieldInfo := FieldInfo{
			FieldName: fieldName,
			FieldType: fieldType,
		}
		//
		tagValue := t.Field(i).Tag.Get(tagKey)
		fieldInfo.TagValue = tagValue
		fieldInfos = append(fieldInfos, fieldInfo)
	}
	return fieldInfos
}

type FieldInfo struct {
	FieldName string
	FieldType string
	TagValue  string
}
