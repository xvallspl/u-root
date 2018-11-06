package pacman

import (
	"testing"
)

func TestPackage(t *testing.T) {
	p, err := Package("extra", "x86_64", "dhcp")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("package: %v", p.info)
}
