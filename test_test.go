package goapidoc

import (
	"fmt"
	"testing"
)

func failNow(t *testing.T, msg string) {
	fmt.Println(msg)
	t.FailNow()
}

func testPanic(t *testing.T, want bool, fn func(), fnName string) {
	didPanic := false
	var msg interface{}
	func() {
		defer func() {
			if msg = recover(); msg != nil {
				didPanic = true
			}
		}()
		fn()
	}()

	if didPanic && !want {
		failNow(t, fmt.Sprintf("Test case for '%s' want no panic but panic with '%s'", fnName, msg))
	} else if !didPanic && want {
		failNow(t, fmt.Sprintf("Test case for '%s' want panic but no panic happened", fnName))
	}
}

func testMatchElements(t *testing.T, s1, s2 []string, s1Name, s2Name string) {
	if len(s1) != len(s2) {
		failNow(t, fmt.Sprintf("Two slice ('%s' and '%s')'s lengths is not same", s1Name, s2Name))
	}
	for _, i1 := range s1 {
		contained := false
		for _, i2 := range s2 {
			if i1 == i2 {
				contained = true
				break
			}
		}
		if !contained {
			failNow(t, fmt.Sprintf("There are some items in '%s' that are not found in '%s'", s1Name, s2Name))
		}
	}
	for _, i2 := range s2 {
		contained := false
		for _, i1 := range s1 {
			if i1 == i2 {
				contained = true
				break
			}
		}
		if !contained {
			failNow(t, fmt.Sprintf("There are some items in '%s' that are not found in '%s'", s2Name, s1Name))
		}
	}
}
