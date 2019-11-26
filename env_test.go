package withenv

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/spf13/afero"
)

const someFile = "/file"

func setupEnv(envs map[string]string) *Env {
	env := New()

	if len(envs) > 0 {
		for key, val := range envs {
			env.Set(key, val)
		}
	}

	return env
}

func TestEnv_GetOrDefault(t *testing.T) {
	tests := []struct {
		name  string
		envs  map[string]string
		key   string
		deflt string
		want  string
	}{
		{"empty env", nil, "x", "DEF", "DEF"},
		{"not existing key", map[string]string{"a": "b"}, "x", "DEF", "DEF"},
		{"existing key", map[string]string{"a": "b"}, "a", "DEF", "b"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			env := setupEnv(tt.envs)
			if got := env.GetOrDefault(tt.key, tt.deflt); got != tt.want {
				t.Errorf("env.GetOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_Get(t *testing.T) {
	tests := []struct {
		name string
		envs map[string]string
		key  string
		ok   bool
		want string
	}{
		{"empty env", nil, "x", false, ""},
		{"not existing key", map[string]string{"a": "b"}, "x", false, ""},
		{"existing key", map[string]string{"a": "b"}, "a", true, "b"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			env := setupEnv(tt.envs)

			got, gotOK := env.Get(tt.key)
			if gotOK != tt.ok {
				t.Errorf("env.Get() ok = %v, want %v", gotOK, tt.ok)
			} else if got != tt.want {
				t.Errorf("env.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_Expand(t *testing.T) {
	tests := []struct {
		name string
		envs map[string]string
		key  string
		want string
	}{
		{"empty env", nil, "$x", ""},
		{"not existing key", map[string]string{"a": "b"}, "$x", ""},
		{"existing key", map[string]string{"a": "b"}, "$a", "b"},
		{"expand key", map[string]string{"a": "b"}, "a $a c", "a b c"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			env := setupEnv(tt.envs)

			if got := env.Expand(tt.key); got != tt.want {
				t.Errorf("env.Expand() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_EnvRead(t *testing.T) {
	tests := []struct {
		name    string
		hasFile bool
		cont    []byte
		envs    map[string]string
		want    map[string]string
		wantErr bool
	}{
		{"no file", false, []byte(""), nil, map[string]string{}, true},
		{"empty file", true, []byte(""), nil, map[string]string{}, false},
		{"set variable", true, []byte("a=b"), nil, map[string]string{"a": "b"}, false},
		{
			"add variable",
			true,
			[]byte("a=b"),
			map[string]string{"b": "c"},
			map[string]string{"a": "b", "b": "c"},
			false,
		},
		{
			"update variable",
			true,
			[]byte("a=b"),
			map[string]string{"a": "c"},
			map[string]string{"a": "b"},
			false,
		},
		{
			"expand variable",
			true,
			[]byte("a=a.$a.$b.$c"),
			map[string]string{"a": "b", "c": "d"},
			map[string]string{"a": "a.b..d", "c": "d"},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			oldFS := envFS
			envFS = setupfs(tt.hasFile, tt.cont)

			env := setupEnv(tt.envs)
			err := env.Read(someFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("readEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(env.Vars, tt.want); diff != nil {
				t.Errorf("Read() Vars %v", diff)
			}
			envFS = oldFS
		})
	}
}

func setupfs(create bool, content []byte) afero.Fs {
	fs := afero.NewMemMapFs()

	if create {
		_ = afero.WriteFile(fs, someFile, content, 0o644)
	}

	return fs
}
