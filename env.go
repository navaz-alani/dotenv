/*
Package golangdotenv implements, to some extent, the dotenv
functionality provided by other development environments such
as NodeJS.

An additional feature that is provided is the chaining of
environment variable files. In any environment file, the Load
function can be reciursively called to read other files.
This functionality is achieved by using the load key, "__GO_LOAD".
Once encountered, this load key makes the Load function to
recursively call itself on the value of the load key, with the same
overwriteFlag that was used in the initial call to the Load function.
*/
package golangdotenv

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

/*
Env is a struct that contains environment variables stored
as key-value pairs under the field'params'

Keys and values are both stored as strings. Therefore if a user
needs to store an integer in the environment variables, it
would have to be converted from a string to an integer when needed.
*/
type Env struct {
	params map[string]string
}

/*
Count returns the number of Params in e.
*/
func (e Env) Count() int {
	return len(e.params)
}

/*
Get retrieves the value of the given
key in e.
*/
func (e Env) Get(key string) string {
	return e.params[key]
}

/*
Set adds the given key-value pair to the
environment e.
*/
func (e Env) Set(key, val string) {
	e.params[key] = val
}

// Keys returns a slice of strings of the keys in e.
func (e Env) Keys() []string {
	paramKeys := make([]string, e.Count())

	for key := range e.params {
		paramKeys = append(paramKeys, key)
	}

	return paramKeys
}

// Values returns a slice of strings of the values in e.
func (e Env) Values() []string {
	paramVals := make([]string, e.Count())

	for _, val := range e.params {
		paramVals = append(paramVals, val)
	}

	return paramVals
}

/*
Merge adds all of other environment's parameters into e.

Keys in e that exist in other will be overwritten,
by values in other, depending on the overwrite boolean.
*/
func (e Env) Merge(other Env, overwrite bool) Env {
	params := make(map[string]string)

	for key, val := range e.params {
		params[key] = val
	}

	for key, val := range other.params {
		isNil := e.Get(key) == ""

		if isNil || !isNil && overwrite {
			params[key] = val
		}
	}

	env := Env{params}
	return env
}

/*
Load returns a struct Env which contains system
parameters in key "params" that have been
imported from the specified file.

There is one optional boolean parameter supported. All other
parameters will be ignored. The first optional boolean
argument specifies whether to overwrite keys when loading
chained environment files. By default, the overwrite flag is
set to be true.

NOTE: Comments in .env file NOT supported yet!
*/
func Load(filename string, overwrite ...bool) (env Env, err error) {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		fmt.Printf("%s\n", err)
		return env, err
	}

	var overwriteFlag bool = true

	if len(overwrite) != 0 {
		overwriteFlag = overwrite[0]
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

		// GO_LOAD environment variable chaining functionality
		if key == "__GO_LOAD" {
			// Recursively Load and over
			tmpEnv, _ := Load(val, overwriteFlag)

			tmpEnv = Env{params}.Merge(tmpEnv, overwriteFlag)
			params = tmpEnv.params
		}

		params[key] = val
	}

	return Env{params: params}, err
}
