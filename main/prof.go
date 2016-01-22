package main

import (
	"github.com/cagnosolutions/bpt"
	"github.com/pkg/profile"
)

func main() {
	//prof := profile.CPUProfile
	prof := profile.MemProfile
	defer profile.Start(prof, profile.ProfilePath(".")).Stop()

	t := bpt.NewTree()
	for i := 0; i < 1000000; i++ {
		x := bpt.UUID()
		t.Set(x, x)
	}
	t.Close()
}
