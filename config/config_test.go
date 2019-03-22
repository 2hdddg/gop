package config

import (
	"testing"
)

func TestPackageFromPath(t *testing.T) {
	cases := []struct {
		c *Config
		i string
		p string // Expected package name
		o bool   // Expected result
	}{
		// One level package matching system path
		{&Config{SystemPath: "/x/y/z"}, "/x/y/z/p/x.go", "p", true},
		// One level package matching workspace path
		{&Config{WorkspacePath: "/x/y/z"}, "/x/y/z/p/x.go", "p", true},
		// Two level package matching system path
		{&Config{WorkspacePath: "/x/y"}, "/x/y/z/p/x.go", "z/p", true},
		// No match
		{&Config{WorkspacePath: "/x/y"}, "/y/z/p/x.go", "", false},
	}

	for _, c := range cases {
		p, o := c.c.PackageFromPath(c.i)
		if o != c.o {
			t.Errorf("Expected success indication to be %v but was %v",
				c.o, o)
		}
		if p != c.p {
			t.Errorf("Expected package to be %v but was %v", c.p, p)
		}
	}
}
