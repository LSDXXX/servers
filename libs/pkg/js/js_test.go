package js

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestScript(t *testing.T) {
	s := New()
	code, err := s.Compile(`
	function main(msg) {
		msg.payload = 'test'
		return msg
	}
	msg = main(msg)
	
	`)
	if err != nil {
		t.Fatal(err)
	}
	s.Start()
	val, err := s.RunCode(context.Background(), code, map[string]interface{}{
		"msg": map[string]interface{}{
			"payload": "oyo",
		},
	}, nil, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	v := ToGoValue(val)
	fmt.Println(v)
}

func BenchmarkScript(b *testing.B) {
	s := New()
	code, err := s.Compile(`
	function main(msg) {
		msg.payload = 'test'
		return msg
	}
	msg = main(msg)
	
	`)
	if err != nil {
		b.Fatal(err)
	}
	s.Start()
	for i := 0; i < b.N; i++ {
		_, err := s.RunCode(context.Background(), code, map[string]interface{}{
			"msg": map[string]interface{}{
				"payload": "oyo",
			},
		}, nil, time.Second)
		if err != nil {
			b.Fatal(err)
		}
	}
}
