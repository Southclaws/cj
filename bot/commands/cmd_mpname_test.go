package commands

import (
	"fmt"
	"testing"
)

func Test_mpname(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(i, mpname())
	}
}
