package commands

import (
	"fmt"
	"testing"
)

func Test_makeDynamic(t *testing.T) {
	for i := 0; i < 50; i++ {
		fmt.Println(i, makeDynamic())
	}
}
