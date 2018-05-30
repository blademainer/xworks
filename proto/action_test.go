package proto

import (
	"testing"
	"github.com/golang/protobuf/proto"
	"log"
)

func Test_Marshal(t *testing.T) {
	test := &Action{
		Name: "hello",
		Type:  17,
	}
	data, err := proto.Marshal(test)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	newTest := &Action{}
	err = proto.Unmarshal(data, newTest)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	// Now test and newTest contain the same data.
	if test.GetName() != newTest.GetName() {
		log.Fatalf("data mismatch %q != %q", test.GetName(), newTest.GetName())
	}
	// etc.
}
