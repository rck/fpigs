package main

import (
	"fmt"
	"testing"
)

func TestFlagUnitLen(t *testing.T) {
	lu, lU := len(units), len(Units)

	if len(units) != len(Units) {
		t.Errorf("Expected len(units) == len(Units), but %d != %d", lu, lU)
	}
}

func TestFlagUnitString(t *testing.T) {
	for _, us := range units {
		uv := Units[us]
		flagName := fmt.Sprintf("TestFlagUnitString%s", us)
		fu := UnitFlag(flagName, uv, "test")
		if fmt.Sprintf("%s", *fu) != us {
			t.Errorf("Expected %s for unit %s", us, us)
		}
	}
}

func TestFlagUnitValue(t *testing.T) {
	for _, us := range units {
		uv := Units[us]
		flagName := fmt.Sprintf("TestFlagUnitValue%s", us)
		fu := UnitFlag(flagName, uv, "test")
		if *fu != uv {
			t.Errorf("Expected %d for unit %s", uv, us)
		}
	}
}

func TestFlagIgnore(t *testing.T) {
	regexs := []string{"x", "^x", "x$"}
	var i ignoreFlag
	for _, r := range regexs {
		i.Set(r)
	}

	lr, li := len(regexs), len(i.Ignores)
	if lr != li {
		t.Errorf("Expected that number of input regex (%d) are set, but found %d", lr, li)
	}
}
