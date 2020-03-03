package dotenv_test

import (
	"reflect"
	"testing"

	"github.com/navaz-alani/dotenv"
)

var params = map[string]string{
	"key1": "test1",
	"key2": "test2",
	"key3": "test3",
	"key4": "test4",
}

func TestLoadNoComments(t *testing.T) {
	e, err := dotenv.Load("tests/noComments.env", false)
	if err != nil {
		t.Fatal("Expected no error; got ", err)
	}

	for k, v := range params {
		if eVal := e.Get(k); eVal != v {
			t.Errorf("Expected %s; got %s", v, eVal)
		}
	}

	paramCount := 4
	if ct := e.Count(); ct != paramCount {
		t.Errorf("Expected count %d; got %d", paramCount, ct)
	}
}

func TestLoadComments(t *testing.T) {
	e, err := dotenv.Load("tests/commented.env", false)
	if err != nil {
		t.Fatal("Expected no error; got ", err)
	}

	for k, v := range params {
		if eVal := e.Get(k); eVal != v {
			t.Errorf("Expected %s; got %s", v, eVal)
		}
	}

	if k5 := e.Get("key5"); k5 != "" {
		t.Errorf("Expected 'key5' to be ''; got %s", k5)
	}

	if k6 := e.Get("key6"); k6 == "" {
		t.Errorf("Expected 'key6' to be non-empty")
	}

	paramCount := 5
	if ct := e.Count(); ct != paramCount {
		t.Errorf("Expected count %d; got %d", paramCount, ct)
	}
}

func TestLoadChainedNoOverwrite(t *testing.T) {
	e, err := dotenv.Load("tests/chaining.env", false)
	if err != nil {
		t.Fatal("Expected no error; got ", err)
	}

	for k, v := range params {
		if eVal := e.Get(k); eVal != v {
			t.Errorf("Expected %s; got %s", v, eVal)
		}
	}
}

func TestLoadChainedOverwrite(t *testing.T) {
	e, err := dotenv.Load("tests/chaining.env", true)
	if err != nil {
		t.Fatal("Expected no error; got ", err)
	}

	for k, v := range params {
		if eVal := e.Get(k); eVal != v {
			if k == "key1" && eVal == "test1-overwrite" {
				continue
			} else {
				t.Errorf("Expected %s; got %s", v, eVal)
			}
		}
	}
}

func TestLoadInvalidChain(t *testing.T) {
	_, err := dotenv.Load("tests/invalidChain.env", true)
	if err == nil {
		t.Fatal("Expected load error; got none", err)
	}
}

func TestEnv_CheckRequired(t *testing.T) {
	e, err := dotenv.Load("tests/undef.env", true)
	if err != nil {
		t.Fatal("Expected no error; got ", err)
	}

	req := []string{"key1", "key2", "key3", "key4"}
	reqUndef := []string{"key1", "key2", "key4"}

	undef := e.CheckRequired(req)
	if len(undef) == 0 {
		t.Errorf("Expected undef to have: %v; got []", req)
	}

	if !reflect.DeepEqual(reqUndef, undef) {
		t.Errorf("Expected undef to be: %v; got %v", reqUndef, undef)
	}
}
