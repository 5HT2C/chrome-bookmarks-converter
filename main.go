package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/5HT2C/chrome-bookmarks-converter/parse"
	"github.com/5HT2C/chrome-bookmarks-converter/util"
	"github.com/virtualtam/netscape-go/v2"
)

var (
	flagProd = flag.Bool("prod", false, "Enables only test mode")
	flagSegf = flag.Bool("unsafe", false, "Ignores errors and attempts to continue")
	flagInfo = flag.Bool("quiet", false, "Disables debug logging")
)

func main() {
	flag.Parse()
	util.LoggerPanic = !*flagSegf
	util.LoggerQuiet = *flagInfo

	bkFilePrefix := "exported-bookmarks"
	bkTimeSuffix := fmt.Sprintf("%v", time.Now().UnixMilli())
	bkFileExport := func(s, ext string) string {
		parts := []string{bkFilePrefix}
		if !util.StringEmpty(s) {
			parts = append(parts, s)
		}
		parts = append(parts, bkTimeSuffix)

		return fmt.Sprintf(
			"%s.%s",
			strings.Join(parts, "-"),
			ext,
		)
	}

	bkFolded, bkUnique, bkExists, bkPrefer :=
		make(parse.GenChildren, 0),
		make(parse.GenChildren, 0),
		make(parse.GenChildren, 0),
		make(parse.GenChildren, 0)

	entries, _ := os.ReadDir(".")
	for _, f := range entries {
		util.Log(util.LogInfo, "main() got dir child", f.Name())

		if *flagProd || (strings.HasPrefix(f.Name(), "test_") && strings.HasSuffix(f.Name(), ".json")) {
			if f.IsDir() {
				continue
			}

			if b, err := os.ReadFile(f.Name()); err == nil {
				var bkGenRoot *parse.Gen

				if err := json.Unmarshal(b, &bkGenRoot); err != nil {
					util.Log(util.LogFuck, "main() json.Unmarshal()", err, string(b))
				}

				bkParsed := bkGenRoot.ToNetscape()
				bkFolded = append(bkFolded, bkGenRoot.CollectChildren()...)
				writeExport(bkParsed, bkFolded, strings.TrimSuffix(f.Name(), ".json")+".html")
			} else {
				util.Log(util.LogWarn, "main() os.ReadFile()", err, f.Name())
			}
		} else {
			util.Log(util.LogInfo, "main() skipped", f.Name())
		}
	}

	bkUnique, bkExists, bkPrefer = (&bkFolded).CollectUnique()
	writeChild(bkUnique, bkFileExport("unique", "json"))
	writeChild(bkExists, bkFileExport("exists", "json"))
	writeChild(bkPrefer, bkFileExport("prefer", "json"))

	bkFolded = make(parse.GenChildren, 0)
	bkFolded = append(bkFolded, bkUnique...)
	bkFolded = append(bkFolded, bkPrefer...)
	bkParsed := (&bkFolded).ToNetscapeUnique()

	writeExport(bkParsed, bkFolded, bkFileExport("", "html"))
}

func writeExport(bkParsed *netscape.Document, bkFolded parse.GenChildren, bkExport string) {
	if len(bkFolded) == 0 {
		util.Log(util.LogWarn, "main() writeExport() skipped", bkExport, bkParsed)
		return
	}

	m, err := netscape.Marshal(bkParsed)
	if err != nil {
		util.Log(util.LogFuck, "main() writeExport() netscape.Marshal()", err, bkParsed)
	} else {
		util.Log(util.LogWarn, "main() writeExport() marshalled", bkExport, len(bkFolded))

		if err := os.WriteFile(bkExport, m, 0644); err != nil {
			util.Log(util.LogWarn, "main() writeExport() write failed", err, bkExport)
		} else {
			util.Log(util.LogInfo, "main() writeExport() finished", bkExport, len(bkFolded))
		}
	}
}

func writeChild(bkFolded parse.GenChildren, bkExport string) {
	if len(bkFolded) == 0 {
		util.Log(util.LogWarn, "main() writeChild() skipped", bkExport)
		return
	}

	if j, err := json.MarshalIndent(&bkFolded, "", "    "); err != nil {
		util.Log(util.LogFuck, "main() writeChild() json.MarshalIndent()", err, &bkFolded)
	} else {
		util.Log(util.LogWarn, "main() writeChild() marshalled", bkExport, len(bkFolded))

		if err := os.WriteFile(bkExport, j, 0644); err != nil {
			util.Log(util.LogWarn, "main() writeChild() write failed", err, bkExport)
		} else {
			util.Log(util.LogInfo, "main() writeChild() finished", bkExport, len(bkFolded))
		}
	}
}
