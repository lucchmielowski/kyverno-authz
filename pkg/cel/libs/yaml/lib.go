package yaml

import (
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/ext"
)

type lib struct {
	yaml Yaml
}

// Lib returns a cel.EnvOption for registering the YAML library
func Lib(yaml YamlImpl) cel.EnvOption {
	return cel.Lib(&lib{
		yaml: Yaml{yaml},
	})
}

// LibraryName implements cel.Lib.LibraryName.
func (*lib) LibraryName() string {
	return "kyverno.yaml"
}

// CompileOptions implements cel.Lib.CompileOptions.
func (l *lib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		ext.NativeTypes(reflect.TypeFor[Yaml]()),
		cel.Variable("yaml", YamlType),
		l.extendEnv,
	}
}

// ProgramOptions implements cel.Lib.ProgramOptions.
func (l *lib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{
		cel.Globals(map[string]any{
			"yaml": l.yaml,
		}),
	}
}

func (*lib) extendEnv(env *cel.Env) (*cel.Env, error) {
	impl := impl{Adapter: env.CELTypeAdapter()}

	libraryDecls := map[string][]cel.FunctionOpt{
		// yaml.Unmarshal(string) -> dyn
		// Parses YAML data from a string and returns the parsed structure
		"Unmarshal": {
			cel.MemberOverload(
				"yaml_unmarshal_string",
				[]*cel.Type{YamlType, types.StringType},
				types.DynType,
				cel.BinaryBinding(impl.unmarshal),
			),
			cel.MemberOverload(
				"yaml_unmarshal_bytes",
				[]*cel.Type{YamlType, types.BytesType},
				types.DynType,
				cel.BinaryBinding(impl.unmarshal),
			),
		},
		// yaml.UnmarshalFromFile(string) -> dyn
		// Reads and parses a YAML file from the local filesystem
		"UnmarshalFromFile": {
			cel.MemberOverload(
				"yaml_unmarshal_from_file",
				[]*cel.Type{YamlType, types.StringType},
				types.DynType,
				cel.BinaryBinding(impl.unmarshalFromFile),
			),
		},
		// yaml.UnmarshalFromURL(string) -> dyn
		// Fetches and parses a YAML file from a URL
		"UnmarshalFromURL": {
			cel.MemberOverload(
				"yaml_unmarshal_from_url",
				[]*cel.Type{YamlType, types.StringType},
				types.DynType,
				cel.BinaryBinding(impl.unmarshalFromURL),
			),
		},
	}

	options := []cel.EnvOption{}
	for name, overloads := range libraryDecls {
		options = append(options, cel.Function(name, overloads...))
	}
	return env.Extend(options...)
}
