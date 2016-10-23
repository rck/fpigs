package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
)

type Unit int64

var units = []string{"B", "KiB", "MiB", "GiB", "TiB", "KB", "MB", "GB", "TB"}

func allUnits() string {
	return strings.Join(units, ", ")
}

var Units = map[string]Unit{
	units[0]: 1,
	units[1]: 1 << 10,
	units[2]: 1 << 20,
	units[3]: 1 << 30,
	units[4]: 1 << 40,
	units[5]: 1000,
	units[6]: 1000 * 1000,
	units[7]: 1000 * 1000 * 1000,
	units[8]: 1000 * 1000 * 1000 * 1000,
}

func (u Unit) String() string {
	for k, v := range Units {
		if v == u {
			return k
		}
	}
	return ""
}

type unitFlag struct{ Unit }

func (u *unitFlag) Set(s string) error {
	if v, ok := Units[s]; ok {
		u.Unit = v
		return nil
	}

	return fmt.Errorf("invalid unit %q, choices: %s", s, allUnits())
}

func UnitFlag(name string, value Unit, usage string) *Unit {
	f := unitFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Unit
}

type Ignores []regexp.Regexp

func (i Ignores) String() string {
	var strs []string
	for _, e := range i {
		strs = append(strs, e.String())
	}
	return fmt.Sprintf("%s", strings.Join(strs, ", "))
}

type ignoreFlag struct{ Ignores }

func (u *ignoreFlag) Set(s string) error {
	r, err := regexp.Compile(s)
	if err != nil {
		return err
	}

	u.Ignores = append(u.Ignores, *r)
	return nil
}

func IgnoreFlag(name string, value Ignores, usage string) *Ignores {
	f := ignoreFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Ignores
}
