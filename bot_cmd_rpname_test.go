package main

import (
	"fmt"
	"testing"
)

func Test_rpname(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(i, rpname())
	}
}
