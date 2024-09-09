package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"unicode"
)
type AssetReq struct {
	CoinName string `json:"coinName,optional" path:"coinName,optional"`
	Ip       string `json:"ip,optional"`
}
// determineTypeProto 根据 Go 类型确定 Protobuf 类型
func determineTypeProto(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int32:
		return "int32"
	case reflect.Int64:
		return "int64"
	case reflect.Float32:
		return "float"
	case reflect.Float64:
		return "double"
	case reflect.Bool:
		return "bool"
	case reflect.Slice:
		return "repeated " + determineTypeProto(t.Elem())
	case reflect.Struct:
		return t.Name() // 直接使用结构体的名字
	default:
		return "bytes" // 其他未知类型默认使用 bytes
	}
}

func toProtobufName(name string) string {
	// Protobuf 字段名通常使用小写字母和下划线
	var result string
	for i, ch := range name {
		if i > 0 && unicode.IsUpper(ch) {
			result += "_"
		}
		result += string(unicode.ToLower(ch))
	}
	return result
}


func toProtobufMessageName(name string) string {
	// Protobuf 消息名通常使用驼峰式命名
	return strings.Title(name)
}

func toProtobufFieldName(name string) string {
	// Protobuf 字段名不能以数字开头，因此确保首字符是字母
	if len(name) > 0 && unicode.IsDigit(rune(name[0])) {
		name = "field" + name
	}

	// 将字段名转换为小写加下划线
	return toProtobufName(name)
}

// GetStruct 生成结构体的 Protobuf 格式描述
func GetStruct(s interface{}) string {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	st := make([]string, 0)

	for i := 0; i < t.NumField(); i++ {
		fieldName := toProtobufFieldName(t.Field(i).Name)
		fieldType := determineTypeProto(t.Field(i).Type)

		// 构建 Protobuf 字段描述
		fieldDesc := fmt.Sprintf("%s %s = %d;", fieldType, fieldName, i+1)
		st = append(st, fieldDesc)
	}

	// 将字段描述连接成一个完整的 message
	return fmt.Sprintf("message %s {\n%s\n}", toProtobufMessageName(t.Name()), "\n"+strings.Join(st, "\n"))
}

func main() {
	
	//类型定义占位
	st := make([]string, 0)

st = append(st ,	GetStruct(new(AssetReq)))

	create, err := os.Create("./output.proto")

	if err != nil {
		fmt.Println(err)
	}

	defer create.Close()

	_, err = create.Write([]byte(strings.Join(st, "\n")))

	if err != nil {
		fmt.Println(err)
	}

}
