package main

import (
	"encoding/json"
	"fmt"

	"github.com/krhubert/null/v1"
)

type OtherStruct struct {
	S string
}

type Person struct {
	Age    null.Null[uint8]
	Weight null.Null[int]
	Other  null.Null[OtherStruct]
}

func main() {
	fmt.Printf("%#v\n", Person{})
	p := Person{Age: null.New[uint8](8)}
	b, err := json.Marshal(p)
	fmt.Println(string(b), err)

	var p1 Person
	json.Unmarshal(b, &p1)
	fmt.Println(p1)

	var p2 Person
	json.Unmarshal([]byte(`{"Age":null,"Other":{"S":":)"}}`), &p2)
	fmt.Println(p2)
}
