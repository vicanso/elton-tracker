# cod-tracker

[![Build Status](https://img.shields.io/travis/vicanso/cod-tracker.svg?label=linux+build)](https://travis-ci.org/vicanso/cod-tracker)

Tracker middleware for cod, it can track route params, include query, params, form and handle result.

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/vicanso/cod"

	bodyparser "github.com/vicanso/cod-body-parser"
	tracker "github.com/vicanso/cod-tracker"
)

func main() {
	d := cod.New()

	d.Use(bodyparser.NewDefault())

	loginTracker := tracker.New(tracker.Config{
		OnTrack: func(info *tracker.Info, _ *cod.Context) {
			fmt.Println(info)
		},
	})

	d.POST("/user/login", loginTracker, func(c *cod.Context) (err error) {
		c.SetHeader(cod.HeaderContentType, cod.MIMEApplicationJSON)
		c.BodyBuffer = bytes.NewBuffer(c.RequestBody)
		return
	})

	d.ListenAndServe(":7001")
}

```