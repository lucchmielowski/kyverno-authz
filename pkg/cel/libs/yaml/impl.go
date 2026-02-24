package yaml

import (
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/kyverno/kyverno-authz/pkg/cel/utils"
)

type impl struct {
	types.Adapter
}

// unmarshal parses YAML data from a string
func (i *impl) unmarshal(receiver ref.Val, data ref.Val) ref.Val {
	yamlObj, err := utils.ConvertToNative[Yaml](receiver)
	if err != nil {
		return types.WrapErr(err)
	}

	dataBytes, err := utils.ConvertToNative[[]byte](data)
	if err != nil {
		// Try converting from string
		dataStr, err := utils.ConvertToNative[string](data)
		if err != nil {
			return types.WrapErr(err)
		}
		dataBytes = []byte(dataStr)
	}

	result, err := yamlObj.Unmarshal(dataBytes)
	if err != nil {
		return types.WrapErr(err)
	}

	return i.NativeToValue(result)
}

// unmarshalFromFile reads and parses a YAML file from the local filesystem
func (i *impl) unmarshalFromFile(receiver ref.Val, path ref.Val) ref.Val {
	yamlObj, err := utils.ConvertToNative[Yaml](receiver)
	if err != nil {
		return types.WrapErr(err)
	}

	pathStr, err := utils.ConvertToNative[string](path)
	if err != nil {
		return types.WrapErr(err)
	}

	result, err := yamlObj.UnmarshalFromFile(pathStr)
	if err != nil {
		return types.WrapErr(err)
	}

	return i.NativeToValue(result)
}

// unmarshalFromURL fetches and parses a YAML file from a URL
func (i *impl) unmarshalFromURL(receiver ref.Val, url ref.Val) ref.Val {
	yamlObj, err := utils.ConvertToNative[Yaml](receiver)
	if err != nil {
		return types.WrapErr(err)
	}

	urlStr, err := utils.ConvertToNative[string](url)
	if err != nil {
		return types.WrapErr(err)
	}

	result, err := yamlObj.UnmarshalFromURL(urlStr)
	if err != nil {
		return types.WrapErr(err)
	}

	return i.NativeToValue(result)
}
