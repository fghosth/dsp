package util

import (
	"fmt"
	"testing"
)

var city []string

func init() {
	city = []string{"shanghai", "benjin", "guangzhou"}
}

func aTestRoaringB(t *testing.T) {

	for i := 0; i < len(city); i++ {
		code := Hashcode(city[i])
		fmt.Println(code)
	}

}

func TestHashcode(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := Hashcode(tt.args.str); got != tt.want {
			t.Errorf("%q. Hashcode() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
