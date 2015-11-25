Library for constructing Belphanior servants with Go.

For more information, visit http://belphanior.net

This library serves as glue between the Go standard net/http handler and the Belphanior protocol. Handler functions accept positional arguments from the path and data argument substitution variables, which are sent as strings to the handling function. The handler can return 

* nothing
* a string (which will be converted into the response to the servant's caller with a 200 HTTP status)
* a string and an error (if error is not nil, it will be sent as the response with a 500 HTTP status; otherwise, the string is sent with 200 status)

Example of usage:

```
package main

import (
        "fmt"
        "github.com/fixermark/belphanior_go/servant"
)

func main() {
        var s servant.Servant
        s.SetRoleUrl("http://belphanior.net/roles/speech/v1")
        s.RegisterHandler(servant.Handler{
                Name: "output",
                Method: "POST",
                Path: "/output/$(output)",
                Handler: func(output string) {
                        fmt.Println("Output:", output)
                },
        })
        s.Run()
}
```
