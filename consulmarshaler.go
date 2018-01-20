package consulmarshaler

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	consul "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

// Client is the consul marshaler
type Client struct {
	kv *consul.KV
}

// New returns a new consul marshaler
func New(conf *consul.Config) (*Client, error) {
	c, err := consul.NewClient(conf)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to consul")
	}
	kv := c.KV()
	client := &Client{kv}
	return client, nil
}

// Unmarshal will pull KV pairs from a prefix in consul and unmarshal them into a struct
func (c *Client) Unmarshal(prefix string, v interface{}) error {

	if prefix == "" {
		return errors.New("prefix required")
	}
	status := false // if at least one field was read successfully, status will be true

	// getter gets values from consul
	getter := func(name string) string {
		pair, _, err := c.kv.Get(name, nil)
		if err != nil || pair == nil {
			log.Println(errors.Wrap(err, "could not read from consul"))
			return ""
		}
		status = true
		return string(pair.Value[:])
	}

	walkRead(getter, prefix, reflect.TypeOf(v).Elem(), reflect.ValueOf(v).Elem())
	if !status {
		return errors.New("no fields found")
	}
	return nil
}

// Marshal will save the struct into consul
func (c *Client) Marshal(prefix string, v interface{}) error {
	if prefix == "" {
		return errors.New("prefix required")
	}
	// setter sets values to consul
	setter := func(name string, value []byte) {
		kv := &consul.KVPair{
			Key:   name,
			Value: value,
		}
		_, err := c.kv.Put(kv, &consul.WriteOptions{})
		if err != nil {
			log.Println(errors.Wrap(err, "could not write to consul"))
		}
	}
	walkWrite(setter, prefix, reflect.TypeOf(v).Elem(), reflect.ValueOf(v).Elem())
	return nil
}

func walkWrite(fn func(string, []byte), prefix string, st reflect.Type, v reflect.Value) {
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		val := v.Field(i)
		switch field.Type.Kind() {
		case reflect.Struct:
			walkWrite(fn, prefix+"/"+field.Name, field.Type, val)
		case reflect.String:
			key := fmt.Sprintf("%s/%s", prefix, field.Name)
			fn(key, []byte(val.String()))
		case reflect.Int:
			key := fmt.Sprintf("%s/%s", prefix, field.Name)
			fn(key, []byte(strconv.Itoa(int(val.Int()))))
		case reflect.Bool:
			key := fmt.Sprintf("%s/%s", prefix, field.Name)
			boolVal := val.Bool()
			if boolVal {
				fn(key, []byte("true"))
			} else {
				fn(key, []byte("false"))
			}
		default:
			log.Println("unsupported type:", field.Type.Kind())
		}
	}
}

func walkRead(fn func(string) string, prefix string, st reflect.Type, v reflect.Value) {
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		val := v.Field(i)
		switch field.Type.Kind() {
		case reflect.Struct:
			walkRead(fn, prefix+"/"+field.Name, field.Type, val)
		case reflect.String:
			key := fmt.Sprintf("%s/%s", prefix, field.Name)
			if cval := fn(key); cval != "" {
				val.SetString(cval)
			}
		case reflect.Int:
			key := fmt.Sprintf("%s/%s", prefix, field.Name)
			if cval := fn(key); cval != "" {
				i, err := strconv.ParseInt(cval, 10, 64)
				if err == nil {
					val.SetInt(i)
				}
			}
		case reflect.Bool:
			key := fmt.Sprintf("%s/%s", prefix, field.Name)
			if cval := fn(key); cval != "" {
				if cval == "true" {
					val.SetBool(true)
				}
				if cval == "false" {
					val.SetBool(false)
				}
			}
		default:
			log.Println("unsupported type:", field.Type.Kind())
		}
	}
}
