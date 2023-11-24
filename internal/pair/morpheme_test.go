package pair

import (
	"fmt"
	"testing"
)

func TestT(t *testing.T) {
	s := "おっきい絵。"
	r := ma.parse(s)
	for _, l := range r.Result {
		fmt.Println(l)
	}
}
