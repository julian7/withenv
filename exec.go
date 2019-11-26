package withenv

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// Exec runs a given command, replacing the current executable.
// It returns only if exec is not possible.
func (env *Env) Exec(cmd string, args ...string) error {
	var err error

	cmd, err = exec.LookPath(cmd)
	if err != nil {
		return err
	}

	return syscall.Exec(cmd, args, env.Environ())
}

// Child runs a given command as a child process, and returns its
// return value
func (env *Env) Child(out, err io.Writer, cmd string, args ...string) error {
	for idx := range args {
		args[idx] = env.Expand(args[idx])
	}

	c := exec.Command(cmd, args...)
	c.Env = env.Environ()
	c.Stdin = os.Stdin
	c.Stdout = out
	c.Stderr = err

	return c.Run()
}

// Run runs a given command, waiting for execution, then returns
// the subprocess' return value.
func (env *Env) Run(cmd string, args ...string) error {
	fmt.Printf("Running %s %s\n", cmd, strings.Join(args, " "))
	return env.Child(os.Stdout, os.Stderr, cmd, args...)
}
