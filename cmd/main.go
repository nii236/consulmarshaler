package main

import (
	"log"

	"github.com/hashicorp/consul/api"
	"github.com/nii236/consulmarshaler"
)

type KeyValue struct {
	First  string
	Second int
	Third  bool
}

func main() {
	// Marshal into consul
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
