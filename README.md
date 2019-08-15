# elton-tracker

[![Build Status](https://img.shields.io/travis/vicanso/elton-tracker.svg?label=linux+build)](https://travis-ci.org/vicanso/elton-tracker)

Tracker middleware for elton, it can track route params, include query, params, form and handle result.

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/vicanso/elton"

	bodyparser "github.com/vicanso/elton-body-parser"
	tracker "github.com/vicanso/elton-tracker"
)

func main() {
	d := elton.New()

	d.Use(bodyparser.NewDefault())

	loginTracker := tracker.New(tracker.Config{
		OnTrack: func(info *tracker.Info, _ *elton.Context) {
			fmt.Println(info)
		},
	})

	d.POST("/user/login", loginTracker, func(c *elton.Context) (err error) {
		c.SetHeader(elton.HeaderContentType, elton.MIMEApplicationJSON)
		c.BodyBuffer = bytes.NewBuffer(c.RequestBody)
		return
	})

	d.ListenAndServe(":7001")
}

```