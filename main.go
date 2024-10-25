package main

import (
	"encoding/json"
	"flag"
	"os"
	"strings"

	"github.com/5HT2C/chrome-bookmarks-converter/parse"
	"github.com/5HT2C/chrome-bookmarks-converter/utils"
	"github.com/virtualtam/netscape-go/v2"
)

var (
	flagProd = flag.Bool("prod", false, "Enables only test mode")
	flagSafe = flag.Bool("unsafe", false, "Ignores errors and attempts to continue")
)

func main() {
	flag.Parse()
	utils.IsSafe = *flagSafe

	entries, _ := os.ReadDir(".")

	for _, f := range entries {
		if *flagProd || (strings.HasPrefix(f.Name(), "test_") && strings.HasSuffix(f.Name(), ".json")) {
			if b, err := os.ReadFile(f.Name()); err == nil {
				var bookmarks *parse.Gen

				if err := json.Unmarshal(b, &bookmarks); err != nil {
					panic(err)
				}

				m, err := netscape.Marshal(bookmarks.ToNetscape())
				if err != nil {
					panic(err)
				}

				//fmt.Printf("TODO: %s\n", f.Name())
				//fmt.Print(string(m))

				if err := os.WriteFile(strings.TrimSuffix(f.Name(), ".json")+".html", m, 0644); err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		}
	}
}
