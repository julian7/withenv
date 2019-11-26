package withenv

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/afero"
)

var (
	envFS        = afero.NewOsFs()
	ErrDeclError = errors.New("env declaration error")
)

// Env contains environment variables for command
type Env struct {
	Vars map[string]string
	sync.RWMutex
}

// New returns a new, empty Env
func New() *Env {
	return &Env{Vars: map[string]string{}}
}

func (env *Env) SetEnviron(envs []string) error {
	for _, item := range envs {
		if err := env.setenv(item); err != nil {
			return fmt.Errorf("setting envinonment %q: %v", item, err)
		}
	}

	return nil
}

// GetOrDefault retrieves an env. value by key, if exists.
// It returns default value otherwise.
func (env *Env) GetOrDefault(key, deflt string) string {
	val, ok := env.Get(key)
	if !ok {
		return deflt
	}

	return val
}

// Get retrieves an environment value by key
func (env *Env) Get(key string) (string, bool) {
	env.RLock()
	defer env.RUnlock()

	val, ok := env.Vars[key]

	return val, ok
}

// Set sets an environment key with a value
func (env *Env) Set(key, val string) {
	val = env.Expand(val)

	env.Lock()
	defer env.Unlock()

	env.Vars[key] = val
}

// setenv sets a key=value pair into environment
func (env *Env) setenv(keyval string) error {
	items := strings.SplitN(keyval, "=", 2)
	if len(items) != 2 {
		return ErrDeclError
	}
	env.Set(items[0], items[1])
	return nil
}

// Expand does variable expansion for strings provided by replacing `${var}`
// and `$var` with values set in Env. Undefined variables replaced by an
// empty string.
func (env *Env) Expand(s string) string {
	return os.Expand(s, func(key string) string {
		return env.GetOrDefault(key, "")
	})
}

// Read loads Env data from file
func (env *Env) Read(fn string) error {
	fd, err := envFS.Open(fn)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	scanner := bufio.NewScanner(bufio.NewReader(fd))

	lineno := 0
	for scanner.Scan() {
		lineno++
		if err := env.setenv(scanner.Text()); err != nil {
			return fmt.Errorf("reading file %q line %d: %w", fn, lineno, err)
		}
	}

	return nil
}

// Environ returns a slice of strings with environment
// declarations
func (env *Env) Environ() []string {
	list := []string{}

	env.RLock()
	defer env.RUnlock()

	for key, val := range env.Vars {
		list = append(list, key+"="+val)
	}

	return list
}
