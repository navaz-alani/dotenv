/*
Package dotenv provides functionality, similar to NodeJs, for
handling environment variables. It is independent of the underlying
operating system and allows more functionality with .env files.

This includes the ability to link together environment variable
source files using the special "__GO_LOAD" key, specifying required
environment variables and more.

This package has been written with concurrency in mind and is therefore
safe for concurrent access by multiple go routines.
*/
package dotenv

import (
	"io/ioutil"
	"regexp"
	"strings"
	"sync"
)

/*
loadKey is used in an env source file to specify a path to
another env source file.
*/
const loadKey = "__GO_LOAD"

/*
These are regular expressions which match lines in the env
variable source files.
*/
const (
	// comment is regexp to match a whole-line comment
	comment = `^[ \t]+#.*?$`
	/*
		entry is a regular expression which matches a line with
		a valid definition of an environment variable.
		Please note that this regular expression has been written
		to ignore entries with an empty string ("") as the value.
	*/
	entry = `^[ \t]*[^=#]+[ \t]*=[ \t]*"[^"]+"[ \t]*(#?.*)?$`
	/*
		kvEntry is used to pick out the key-value pair from a
		line which could possibly include an inline comment.
	*/
	kvEntry = `^[ \t]*[^=#]+[ \t]*=[ \t]*"[^"]+"[ \t]*`
)

/*
Env is a type which defines a collection of environment variables.
*/
type Env struct {
	/*
		mu is a mutex which guards access to
		the structure.
	*/
	mu sync.Mutex
	/*
		vars is the underlying map which stores
		environment variables.
	*/
	vars map[string]string
}

/*
Load reads environment variables from the given source
file. It returns a pointer to an Env type which holds the
environment variables in that file.
If all goes well, err will be nil.

Values must be enclosed within quotes, but quotes within the
value are not permitted (feature to be added). This requirement
comes as a result of allowing comments in the env source files.
Also, note that an entry should exist on ONE line only.

Comments can begin a line/start in the middle and continue
until the end of the line.
*/
func Load(source string, overWrite bool) (e *Env, err error) {
	e = &Env{
		mu:   sync.Mutex{},
		vars: make(map[string]string),
	}

	raw, err := ioutil.ReadFile(source)
	if err != nil {
		return nil, err
	}

	file := string(raw)
	lines := strings.Split(file, "\n")

	validEntry := regexp.MustCompile(entry)
	commentLine := regexp.MustCompile(comment)
	keyValEntry := regexp.MustCompile(kvEntry)

	for _, line := range lines {
		if commentLine.MatchString(line) ||
			!validEntry.MatchString(line) ||
			strings.TrimSpace(line) == "" {
			continue
		}

		// split line to ignore comment
		entry := keyValEntry.FindString(line)
		kvPair := strings.Split(entry, "=")

		key, val := strings.TrimSpace(kvPair[0]), strings.TrimSpace(kvPair[1])
		// remove quotes
		val = val[1 : len(val)-1]

		// recursive load if loadKey encountered
		if key == loadKey {
			subEnv, err := Load(val, overWrite)
			if err != nil {
				return nil, err
			}

			e.Merge(subEnv, overWrite)
		}

		e.vars[key] = val
	}

	return e, nil
}

/*
Merge adds all of the keys and values in the given env
into e. The overWrite parameter specifies whether keys in e
that are also in  env should be overWritten with their values
in env.
*/
func (e *Env) Merge(env *Env, overWrite bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for k, v := range env.vars {
		if _, ok := e.vars[k]; ok && !overWrite {
			continue
		}
		e.vars[k] = v
	}
}

/*
Get retrieves the value of the given key in the env variables
held within e.
*/
func (e *Env) Get(key string) string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.vars[key]
}

/*
Count returns the number of key-value pairs of environment
variables stored in e.
*/
func (e *Env) Count() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.vars)
}
