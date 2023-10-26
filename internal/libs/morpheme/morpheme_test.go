package morpheme

import (
	"fmt"
	"testing"
)

func TestT(t *testing.T) {
	ma, err := UseMorphologicalAnalytics()
	if err != nil {
		t.Error(err)
	}

	s := "おっきい絵。"
	r := ma.Parse(s)
	for _, l := range r.Result {
		fmt.Println(l)
	}
}
