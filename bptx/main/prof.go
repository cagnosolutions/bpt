package main

import (
	"github.com/cagnosolutions/bpt/bptx"
	"github.com/pkg/profile"
)

func main() {
	prof := profile.CPUProfile
	//prof := profile.MemProfile
	defer profile.Start(prof, profile.ProfilePath(".")).Stop()

	t := bptx.NewTree()
	for i := 0; i < 10000; i++ {
		x := bptx.UUID()
		t.Set(x, x)
	}
	t.Close()
}
