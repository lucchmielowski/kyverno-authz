package yaml_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/kyverno/kyverno-authz/pkg/cel/impl"
	yamlcel "github.com/kyverno/kyverno-authz/pkg/cel/libs/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYamlLibrary(t *testing.T) {
	env, err := cel.NewEnv(
		yamlcel.Lib(&impl.YamlImpl{}),
	)
	require.NoError(t, err)

	tests := []struct {
		name       string
		expression string
		setup      func(t *testing.T) string // returns cleanup path if needed
		want       any
		wantErr    bool
	}{
		{
			name:       "unmarshal simple yaml string",
			expression: `yaml.Unmarshal("name: test\nvalue: 42")`,
			want:       map[string]any{"name": "test", "value": float64(42)},
		},
		{
			name:       "unmarshal yaml with nested structure",
			expression: `yaml.Unmarshal("person:\n  name: John\n  age: 30")`,
			want:       map[string]any{"person": map[string]any{"name": "John", "age": float64(30)}},
		},
		{
			name:       "unmarshal yaml array",
			expression: `yaml.Unmarshal("- apple\n- banana\n- cherry")`,
			want:       []any{"apple", "banana", "cherry"},
		},
		{
			name: "unmarshal from file",
			setup: func(t *testing.T) string {
				tmpFile := filepath.Join(t.TempDir(), "test.yaml")
				content := "database:\n  host: localhost\n  port: 5432"
				err := os.WriteFile(tmpFile, []byte(content), 0644)
				require.NoError(t, err)
				return tmpFile
			},
			expression: `yaml.UnmarshalFromFile(%q)`,
			want:       map[string]any{"database": map[string]any{"host": "localhost", "port": float64(5432)}},
		},
		{
			name:       "unmarshal invalid yaml",
			expression: `yaml.Unmarshal("invalid: yaml: content: [")`,
			wantErr:    true,
		},
		{
			name:       "unmarshal from non-existent file",
			expression: `yaml.UnmarshalFromFile("/nonexistent/path/file.yaml")`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expression := tt.expression
			if tt.setup != nil {
				path := tt.setup(t)
				expression = formatExpression(tt.expression, path)
			}

			ast, issues := env.Compile(expression)
			require.NoError(t, issues.Err())

			prg, err := env.Program(ast)
			require.NoError(t, err)

			result, _, err := prg.Eval(map[string]any{})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, result.Value())
		})
	}
}

func formatExpression(expr, arg string) string {
	// Handle format strings in expressions using fmt.Sprintf
	if len(arg) > 0 {
		return fmt.Sprintf(expr, arg)
	}
	return expr
}
