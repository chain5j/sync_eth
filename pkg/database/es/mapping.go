// description: sync_eth 
// 
// @author: xwc1125
// @date: 2020/10/05
package es

import (
	"encoding/json"
	"fmt"
	"github.com/chain5j/sync_eth/pkg/util/reflectutil"
	"strings"
	"unicode"
)

// 获取m对象的mapping
// 需要注意的是：结构体的json字符串需为下划线的形式（UnderScoreCase）
func Mapping(m interface{}) string {
	properties := make(map[string]Property)
	fieldInfo := reflectutil.GetTagName(m, "es")

	for _, f := range fieldInfo {
		tagValueAll := f.TagValue
		tagValues := strings.Split(tagValueAll, ",")
		typeValue := getDataTypeValue(tagValues)
		if typeValue == "" {
			typeValue = getDataType(f.FieldType)
		}
		maps := arrayToMap(tagValues)
		_, store := maps["store"]
		_, fielddata := maps["fielddata"]

		key := Camel2Case(f.FieldName)
		properties[key] = Property{
			Type:      typeValue,
			Store:     store,
			Fielddata: fielddata,
		}
	}
	esMappings := EsMappings{
		Settings: &Settings{
			NumberOfShards:   1,
			NumberOfReplicas: 0,
		},
		Mappings: &Mappings{
			Properties: properties,
		},
	}
	mapData := esMappings.String()
	fmt.Println("mapData:" + mapData)
	return mapData
}

func Camel2Case(name string) string {
	buffer := NewBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}
	return buffer.String()
}

type EsMappings struct {
	Settings *Settings `json:"settings,omitempty"`
	Mappings *Mappings `json:"mappings,omitempty"`
}

func (e EsMappings) String() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}

type Settings struct {
	NumberOfShards   int64 `json:"number_of_shards,omitempty"`
	NumberOfReplicas int64 `json:"number_of_replicas,omitempty"`
}

type Mappings struct {
	Properties map[string]Property `json:"properties"`
}

type Property struct {
	Type      string `json:"type"`
	Store     bool   `json:"store,omitempty"`
	Fielddata bool   `json:"fielddata,omitempty"`
}

func getDataType(t string) string {
	has := hasEsDataKey(t)
	if has {
		return t
	}
	switch t {
	// ES type
	// Golang type
	case "Time":
		return "date"
	case "string":
		return "text"
	case "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "int":
		return "integer"
	case "float32", "float64", "complex64", "complex128", "rune", "uintptr":
		return "double"
	default:
		return ""
	}
}

func getDataTypeValue(s []string) string {
	for _, v := range s {
		dataType := getDataType(v)
		if dataType != "" {
			return dataType
		}
	}
	return ""
}

var (
	esDataKey = map[string]string{
		"text": "", "keyword": "",
		"integer": "", "long": "", "short": "", "byte": "",
		"double": "", "float": "", "half_float": "", "scaled_float": "",
		"boolean": "", "date": "", "range": "", "binary": "",
		"array": "", "object": "", "nested": "", "geo_point": "", "geo_shape": "",
		"ip": "", "completion": "", "token_count": "", "attachment": "", "percolator": "",
	}
	esSysKey = map[string]string{
		"store": "", "fielddata": "",
	}
)

func hasEsDataKey(k string) bool {
	_, ok := esDataKey[k]
	return ok
}

func hasEsSysKey(k string) bool {
	_, ok := esSysKey[k]
	return ok
}

func arrayToMap(strs []string) map[string]string {
	maps := make(map[string]string)
	if strs == nil || len(strs) == 0 {
		return maps
	}
	for _, s := range strs {
		maps[s] = ""
	}
	return maps
}
