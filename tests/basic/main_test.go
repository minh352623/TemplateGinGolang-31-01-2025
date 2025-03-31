package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddOne(t *testing.T) {
	// var (
	// 	input = 1
	// 	want  = 3
	// )

	// got := AddOne(input)
	// if got != want {
	// 	t.Errorf("AddOne(%d) = %d; want %d", input, got, want)
	// }
	assert.Equal(t, AddOne(2), 3, "AddOne(2) should be 3")
}