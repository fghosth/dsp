package filter

import (
	"fmt"
	"testing"
	"time"
)

func TestDev(t *testing.T) {
	a := time.Now().UTC()

	fmt.Println(int(a.Weekday()))
}
