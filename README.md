# jsonmodifier

`jsonmodifier` is a Go library designed to provide fine-grained control over JSON serialization of objects and arrays. It allows users to specify which fields should be included or excluded when generating JSON, with support for handling nested structures using dot notation to specify field paths.

## Features

- **Selective Inclusion/Exclusion**: Choose exactly which fields to include or exclude when serializing to JSON.
- **Nested Structure Support**: Easily handle multi-level data structures by specifying field paths with dot notation.
- **Customizable**: Integrate seamlessly with existing Go structs and JSON tags.

## Installation

To install `jsonmodifier`, use `go get`:

```sh
go get github.com/puper/jsonmodifier
```

## Usage

Here's a simple example of how to use `jsonmodifier` to include or exclude fields from a JSON serialization:

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/puper/jsonmodifier"
)

type User struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Password string  `json:"password"` // We don't want to include this in the JSON
	Profile  Profile `json:"profile"`
}

type Profile struct {
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func main() {
	user := User{
		ID:       1,
		Username: "johndoe",
		Password: "supersecret",
		Profile: Profile{
			Email: "johndoe@example.com",
			Age:   30,
		},
	}
	obj := jsonmodifier.JsonModify(user)
	// Specify fields to include. Fields not listed here will be excluded.
	includeFields := []string{"id", "username", "profile.email"}
	obj.Only(includeFields...)
	b, _ := json.Marshal(obj)
	fmt.Printf("%v\n", string(b))
	// Alternatively, specify fields to exclude. Fields listed here will be omitted.
	excludeFields := []string{"password", "profile.email"}

	obj.Except(excludeFields...)
	b, _ = json.Marshal(obj)
	fmt.Printf("%v\n", string(b))
	users := [][]User{
		{
			user,
		},
	}
	obj = jsonmodifier.JsonModify(users)
	obj.Only(includeFields...)
	b, _ = json.Marshal(obj)
	fmt.Printf("%v\n", string(b))

	obj.Except(excludeFields...)
	b, _ = json.Marshal(obj)
	fmt.Printf("%v\n", string(b))
}
```

The output will be:

```
{"id":1,"profile":{"email":"johndoe@example.com"},"username":"johndoe"}
{"id":1,"profile":{"age":30},"username":"johndoe"}
[[{"id":1,"profile":{"email":"johndoe@example.com"},"username":"johndoe"}]]
[[{"id":1,"profile":{"age":30},"username":"johndoe"}]]
```