package golangdotenv

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Env is a struct that contains the
//   system environment variables
type Env struct {
	params map[string]string
}

// Count returns the number of Params in e
func (e Env) Count() int {
	return len(e.params)
}

// Get retrieves the value of the given
//   key in e
func (e Env) Get(key string) string {
	return e.params[key]
}

// Set adds the given key-value pair to e
func (e Env) Set(key, val string) {
	e.params[key] = val
}

// Keys returns a slice of strings of the keys in e
func (e Env) Keys() []string {
	paramKeys := make([]string, e.Count())

	for key := range e.params {
		paramKeys = append(paramKeys, key)
	}

	return paramKeys
}

// Values returns a slice of strings of the values in e
func (e Env) Values() []string {
	paramVals := make([]string, e.Count())

	for _, val := range e.params {
		paramVals = append(paramVals, val)
	}

	return paramVals
}

// Merge modifies returns a struct Env containing params
//   from both e1 and e2.
//   Keys in e1 that exist in e2 will be overwritten,
//   by values in e2, depending on the overwrite boolean.
func Merge(e1, e2 Env, overwrite bool) Env {
	params := make(map[string]string)

	for key, val := range e1.params {
		params[key] = val
	}

	for key, val := range e2.params {
		isNil := e1.Get(key) == ""

		if isNil || !isNil && overwrite {
			params[key] = val
		}
	}

	env := Env{params}
	return env
}

// Load returns a struct Env which contains system
//   parameters in key "params" that have been
//   imported from the specified file.
// NOTE: Does not yet support comments in .env file
func Load(filename string) (env Env, err error) {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		fmt.Printf("%s\n", err)
		return env, err
	}

	fileBytes, _ := ioutil.ReadAll(file)
	lines := strings.Split(string(fileBytes), "\n")

	params := make(map[string]string)

	for i := 0; i < len(lines); i++ {
		currParam := strings.Split(lines[i], "=")

		if len(currParam) < 2 {
			continue
		}

		key := strings.TrimSpace(currParam[0])
		val := strings.TrimSpace(currParam[1])

		if key == "GO_LOAD" {
			tmpEnv, _ := Load(val)

			tmpEnv = Merge(Env{params}, tmpEnv, true)
			params = tmpEnv.params
		}

		params[key] = val
	}

	return Env{params: params}, err
}
