package yaml

import (
	"reflect"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

var YamlType = types.NewOpaqueType("yaml.Yaml")

// YamlImpl is the interface that concrete implementations must satisfy
type YamlImpl interface {
	// Unmarshal parses YAML data and returns the result
	Unmarshal(data []byte) (any, error)
	// UnmarshalFromFile reads and parses a YAML file from the local filesystem
	UnmarshalFromFile(path string) (any, error)
	// UnmarshalFromURL fetches and parses a YAML file from a URL
	UnmarshalFromURL(url string) (any, error)
}

// Yaml is the CEL-visible type that wraps the implementation
type Yaml struct {
	YamlImpl
}

// ConvertToNative implements ref.Val.ConvertToNative.
func (d Yaml) ConvertToNative(typeDesc reflect.Type) (any, error) {
	switch typeDesc {
	case reflect.TypeFor[Yaml]():
		return d, nil
	case reflect.TypeFor[YamlImpl]():
		return d.YamlImpl, nil
	}
	return d.YamlImpl, nil
}

// ConvertToType implements ref.Val.ConvertToType.
func (d Yaml) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case YamlType:
		return d
	case types.TypeType:
		return YamlType
	}
	return types.NewErr("type conversion error from '%s' to '%s'", YamlType, typeVal)
}

// Equal implements ref.Val.Equal.
func (d Yaml) Equal(other ref.Val) ref.Val {
	o, ok := other.(Yaml)
	if !ok {
		return types.MaybeNoSuchOverloadErr(other)
	}
	return types.Bool(d.YamlImpl == o.YamlImpl)
}

// Type implements ref.Val.Type.
func (Yaml) Type() ref.Type {
	return YamlType
}

// Value implements ref.Val.Value.
func (d Yaml) Value() any {
	return d.YamlImpl
}
