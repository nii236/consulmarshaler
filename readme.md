[![](https://godoc.org/github.com/nii236/consulmarshaler?status.svg)](http://godoc.org/github.com/nii236/consulmarshaler)
[![Go Report Card](https://goreportcard.com/badge/github.com/nii236/consulmarshaler)](https://goreportcard.com/report/github.com/nii236/consulmarshaler)

# Golang Consul Marshaler

So you can use consul for its unintended purpose, a storage backend!

# Usage

Run Consul.

```
consul agent -dev
```

Use the package.

```go
package main

import (
	"github.com/hashicorp/consul/api"
	"github.com/nii236/consulmarshaler"
)

type KeyValue struct {
	First  string
	Second int
	Third  bool
}

func main() {

	// Marshal to consul
	KVCase := &KeyValue{
		First:  "Hello",
		Second: 1,
		Third:  true,
	}
	m, err := consulmarshaler.New(api.DefaultConfig())
	if err != nil {
		log.Println(err)
	}

	err = m.Marshal("testmarshal", KVCase)
	if err != nil {
		log.Println(err)
	}

	// Unmarshal from consul
	unmarshalled := &KeyValue{}
	err = m.Unmarshal("testmarshal", unmarshalled)
	if err != nil {
		log.Println(err)
	}

	log.Printf("%+v", unmarshalled)
}
```

Fiddle with your values in Consul.

![](/images/consul.png)

# Todo

- Support more data types (struct pointers most useful for now)

# Contributions

- Make pull request