//go:build !windows

package main

// Command-line harness used to test the patch engine on Linux.
//   patcher status DIR
//   patcher apply DIR id,id,...   (empty list = original)
//   patcher intro DIR on|off

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

var lastLayers = map[string]string{}

func setColor256(path string) error {
	lastLayers[path] = mergeLayers(lastLayers[path])
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: status DIR | apply DIR ids | intro DIR on|off")
		os.Exit(2)
	}
	dir := os.Args[2]
	switch os.Args[1] {
	case "status":
		for _, f := range Files {
			st := Analyze(dir, f)
			var on []string
			for k := range st.On {
				on = append(on, k)
			}
			sort.Strings(on)
			fmt.Printf("%s present=%v known=%v backup=%v canRestore=%v on=%s\n", f.Name, st.Present, st.Known,
				st.FromBackup, st.Original != nil, strings.Join(on, ","))
		}
		fmt.Println("intro", IntroState(dir))
	case "apply":
		want := map[string]bool{}
		if len(os.Args) > 3 {
			for _, id := range strings.Split(os.Args[3], ",") {
				if id != "" {
					want[id] = true
				}
			}
		}
		for _, f := range Files {
			st := Analyze(dir, f)
			ch, err := ApplyFile(dir, st, want)
			fmt.Println(f.Name, "changed", ch, "err", err)
			if err != nil {
				os.Exit(1)
			}
		}
	case "intro":
		err := SetIntro(dir, os.Args[3] == "on")
		fmt.Println("intro", IntroState(dir), err, lastLayers)
	}
}
