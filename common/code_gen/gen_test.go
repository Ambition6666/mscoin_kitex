package code_gen

import (
	"testing"
)

func TestYamlToStruct(t *testing.T) {
	n := NewYaml()
	err := n.GetStruct("./test.yaml")
	if err != nil {
		t.Error(err)
	}
}

func TestYamlToJson(t *testing.T) {
	n := NewYaml()
	err := n.ToJSON("./test.yaml", "./output.json")
	if err != nil {
		t.Error(err)
	}
}

func TestNewProto(t *testing.T) {
	p := NewProto()
	p.WriteProto("test.go", "output.go")
}
