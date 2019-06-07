package tracker

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/cod"
)

func TestNoTrackPanic(t *testing.T) {
	assert := assert.New(t)
	done := false
	defer func() {
		r := recover()
		assert.Equal(r.(error), errNoTrackFunction)
		done = true
	}()

	New(Config{})
	assert.True(done)
}

func TestConverMap(t *testing.T) {
	assert := assert.New(t)
	assert.Nil(convertMap(nil, nil))
	assert.Equal(convertMap(map[string]string{
		"password": "123",
		"foo":      "bar",
	}, defaultMaskFields), map[string]string{
		"foo":      "bar",
		"password": "***",
	})
}

func TestTracker(t *testing.T) {
	assert := assert.New(t)
	customErr := errors.New("abcd")
	done := false
	fn := New(Config{
		OnTrack: func(info *Info, _ *cod.Context) {
			assert.Equal(info.Result, HandleFail)
			assert.Equal(info.Query["type"], "1")
			assert.Equal(info.Query["passwordType"], "***")
			assert.Equal(info.Params["category"], "login")
			assert.Equal(info.Form["password"], "***")
			done = true
		},
	})
	req := httptest.NewRequest("POST", "/users/login?type=1&passwordType=2", nil)
	c := cod.NewContext(nil, req)
	c.RequestBody = []byte(`{
		"account": "tree.xie",
		"password": "password"
	}`)
	c.Params = map[string]string{
		"category": "login",
	}
	c.Next = func() error {
		return customErr
	}
	err := fn(c)
	assert.Equal(err, customErr)
	assert.True(done, "tracker middleware fail")
}

// https://stackoverflow.com/questions/50120427/fail-unit-tests-if-coverage-is-below-certain-percentage
func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	rc := m.Run()

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 0.9 {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}
	os.Exit(rc)
}
