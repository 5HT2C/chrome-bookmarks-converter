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

	entries, _ := os.ReadDir(".")
	bkFolded := make([]netscape.Bookmark, 0)
	bkExport := fmt.Sprintf("exported-bookmarks-%v.json", time.Now().UnixMilli())

	for _, f := range entries {
		util.Log(util.LogInfo, "main() got dir child", f)

		if *flagProd || (strings.HasPrefix(f.Name(), "test_") && strings.HasSuffix(f.Name(), ".json")) {
			if f.IsDir() {
				continue
			}

			if b, err := os.ReadFile(f.Name()); err == nil {
				var bkParsed *parse.Gen

				if err := json.Unmarshal(b, &bkParsed); err != nil {
					util.Log(util.LogFuck, "main() json.Unmarshal()", err, string(b))
				}

				m, err := netscape.Marshal(bkParsed.ToNetscape())
				if err != nil {
					util.Log(util.LogFuck, "main() netscape.Marshal()", err, bkParsed)
				}

				bkFolded = append(bkFolded, bkParsed.CollectBookmarks()...)

				util.Log(util.LogWarn, "main() marshalled file", f)

				if err := os.WriteFile(strings.TrimSuffix(f.Name(), ".json")+".html", m, 0644); err != nil {
					util.Log(util.LogWarn, "main() os.WriteFile()", err, f.Name())
				}
			} else {
				util.Log(util.LogWarn, "main() os.ReadFile()", err, f.Name())
			}
		} else {
			util.Log(util.LogInfo, "main() skipped", f.Name())
		}
	}

	if j, err := json.MarshalIndent(bkFolded, "", "    "); err != nil {
		util.Log(util.LogFuck, "main() json.Unmarshal() bkFolded", err, bkFolded)
	} else {
		if err := os.WriteFile(bkExport, j, 0644); err != nil {
			util.Log(util.LogWarn, "main() os.WriteFile()", err, bkExport)
		} else {

		}
	}
}
