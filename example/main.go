package main

import (
	"bytes"
	"encoding/json"
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
			buf, _ := json.Marshal(info)
			fmt.Println(string(buf))
		},
	})

	d.POST("/user/login", loginTracker, func(c *elton.Context) (err error) {
		c.SetHeader(elton.HeaderContentType, elton.MIMEApplicationJSON)
		c.BodyBuffer = bytes.NewBuffer(c.RequestBody)
		return
	})

	err := d.ListenAndServe(":3000")
	if err != nil {
		panic(err)
	}
}
