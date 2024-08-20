package code_gen

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Yaml struct {
	Struct []string
}

func NewYaml() *Yaml {
	return &Yaml{}
}

// GetStruct 读取 YAML 文件并生成相应的 Go 结构体
func (y *Yaml) GetStruct(filepath string) (err error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return
	}
	var temp map[string]interface{}
	err = yaml.Unmarshal(data, &temp)
	if err != nil {
		return
	}

	// 生成结构体
	y.Struct = append(y.Struct, y.generateStructCode("config", temp))

	create, err := os.Create("./struct_gen.go")
	if err != nil {
		return err
	}
	defer create.Close()

	_, err = create.Write([]byte("package code_gen\n\n" + strings.Join(y.Struct, "\n")))

	y.Struct = []string{}

	return
}

// generateStructCode 递归生成结构体代码
func (y *Yaml) generateStructCode(structName string, yamlMap map[string]interface{}) string {
	var structFields []string

	for key, value := range yamlMap {
		fieldName := toCamelCase(key)
		fieldType := determineType(value)

		if nestedMap, ok := value.(map[string]interface{}); ok {
			nestedStructName := fieldName
			nestedStruct := y.generateStructCode(nestedStructName, nestedMap)
			y.Struct = append(y.Struct, nestedStruct)
			fieldType = nestedStructName
		} else if array, ok := value.([]interface{}); ok {
			// Determine the element type of the array
			if len(array) > 0 {
				elemType := determineType(array[0])
				fieldType = "[]" + elemType
			} else {
				fieldType = "[]interface{}"
			}
		}

		structFields = append(structFields, fmt.Sprintf("\t%s %s `yaml:\"%s\"`", fieldName, fieldType, key))
	}

	return fmt.Sprintf("type %s struct {\n%s\n}\n", structName, strings.Join(structFields, "\n"))

}

// determineType 根据 YAML 数据类型确定 Go 类型
func determineType(value interface{}) string {
	switch value.(type) {
	case string:
		return "string"
	case int, int64, int32:
		return "int"
	case float64, float32:
		return "float64"
	case bool:
		return "bool"
	case []interface{}:
		return "[]interface{}"
	case map[interface{}]interface{}:
		return "struct"
	default:
		return "interface{}"
	}
}

// toCamelCase 将 YAML 字段名称转换为驼峰式
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

func (y *Yaml) ToJSON(inputPath string, outputPath string) (err error) {
	// 读取 YAML 文件
	yamlFile, err := os.ReadFile(inputPath)
	if err != nil {
		return
	}

	// 解析 YAML 数据
	var yamlData map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &yamlData)
	if err != nil {
		return
	}

	// 转换为 JSON
	jsonData, err := json.MarshalIndent(yamlData, "", "  ")
	if err != nil {
		return
	}

	// 将 JSON 数据写入文件
	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		return
	}

	return
}
