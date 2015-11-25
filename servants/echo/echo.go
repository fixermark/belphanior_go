// Echoing test

package main

import (
	"fmt"
	"github.com/fixermark/belphanior_go/servant"
)

func main() {
	var s servant.Servant
	s.SetRoleUrl("http://belphanior.net/roles/speech/v1")
	s.RegisterHandler(servant.Handler{
		Name:   "output",
		Method: "POST",
		Path:   "/say",
		Data:   "$(output)",
		Handler: func(msg string) {
			fmt.Printf("%s\n", msg)
		},
	})

	s.Run()
}
