// Simple time servant

package main

import (
	"time"
	"github.com/fixermark/belphanior_go/servant"
)

func main() {
	var s servant.Servant
	s.SetRoleUrl("http://belphanior.net/roles/time/v1")
	s.RegisterHandler(servant.Handler{
		Name: "get time string",
		Method: "GET",
		Path: "/time",
		Handler: func() string {
			t := time.Now()
			return t.Format("Monday, January 2, 2006, 03:04 PM")
		},
	})

	s.Run()
}

