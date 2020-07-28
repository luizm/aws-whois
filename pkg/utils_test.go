package whois

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveElementOfSliceString(t *testing.T) {
	a := []string{"a", "b", "c", "d", "e"}
	b := []string{"c", "d"}
	result := []string{"a", "b", "e"}

	diff := RemoveElementOfSliceString(a, b)
	fmt.Println(diff)
	assert.Equal(t, diff, result)
}
