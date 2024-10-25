package parse

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/5HT2C/chrome-bookmarks-converter/util"
	"github.com/virtualtam/netscape-go/v2"
)

var (
	rootName   = fmt.Sprintf("root-folder-%v", time.Now().UnixNano())
	rootFolder = GenChild{
		Name: rootName,
		Guid: fmt.Sprintf("%s", sha256.Sum256([]byte(rootName))),
	}
)

func (g *Gen) ToNetscape() *netscape.Document {
	g.Origin = g.PopulateOrigin()

	genTitle := fmt.Sprintf("Bookmarks-%s-%s", g.Origin, g.Checksum)
	util.Log(util.LogInfo, "Gen.ToNetscape() title", genTitle)

	return &netscape.Document{
		Title: genTitle,
		Root: netscape.Folder{
			Name:        genTitle,
			Description: fmt.Sprintf("%s - %s", "root", g.Checksum),
			Subfolders:  g.Roots.ToNetscape(g.PopulateRoots()),
		},
	}
}

func (g *Gen) CollectBookmarks() []netscape.Bookmark {
	bookmarks := make([]netscape.Bookmark, 0)

	for _, folder := range *g.PopulateRoots() {
		folder.Children = append(make([]GenChild, 0), folder.Children...) // unnecessary imo
		bookmarks = append(bookmarks, folder.ToNetscapeBookmarks()...)    // we should avoid calling this twice, ideally, this is not good
	}

	return bookmarks
}

func (g GenOrigin) String() string {
	switch g {
	case GenOriginUnknown:
		fallthrough
	case GenOriginChrome:
		return "chrome"
	case GenOriginEdge:
		return "edge"
	case GenOriginOther:
		return "other"
	default:
		return "unknown"
	}
}

func (g *GenChild) Description() string {
	sep := " - "
	if len(g.Type) == 0 && len(g.Guid) == 0 {
		sep = ""
	}

	return fmt.Sprintf("%s%s%s", g.Type, sep, g.Guid)
}

func (g *GenChild) makeBookmark() netscape.Bookmark {
	b := netscape.Bookmark{
		Title:       g.Name,
		URL:         g.Url,
		Description: g.Description(),
		Private:     false,
		Tags:        nil,
		Attributes:  g.AttrStr(nil),
	}

	util.Log(util.LogInfo, "GenFolder.ToNetscapeBookmarks() made bookmark", g.Name, g.Guid)
	return b
}

func (g *GenChild) ToNetscapeBookmarks() []netscape.Bookmark {
	bookmarks := make([]netscape.Bookmark, 0)
	childUrls := make([]string, 0)

	if !g.IsFolder() {
		bookmarks = append(bookmarks, g.makeBookmark())
	}

	for _, child := range g.Children {
		if child.IsFolder() {
			continue
		}

		bookmarks = append(bookmarks, child.makeBookmark())
		childUrls = append(childUrls, child.Url)
	}

	util.Log(util.LogInfo, "GenFolder.ToNetscapeBookmarks() found", g.Name, g.Guid, len(bookmarks), childUrls)
	return bookmarks
}

func (g *GenChild) makeFolder(caller string) netscape.Folder {
	f := netscape.Folder{
		Description: g.Description(),
		Name:        g.Name,
		Attributes:  g.AttrStr(nil),
		Bookmarks:   g.ToNetscapeBookmarks(),  // we fucking love recurse around here
		Subfolders:  g.ToNetscapeSubfolder(g), // whoooo yeah!!!! recursive loops!!!!!
	}

	util.Log(util.LogInfo, caller+" made folder", g.Name, g.Guid)
	return f
}

func (g *GenChild) ToNetscapeSubfolder(parent *GenChild) []netscape.Folder {
	folders := make([]netscape.Folder, 0)
	subDirs := make([]string, 0)

	if g.IsFolder() && util.StringOrDefault(g.Guid, g.Name) != util.StringOrDefault(parent.Guid, parent.Name) {
		folders = append(folders, g.makeFolder("GenFolder.ToNetscapeSubfolder()"))
	}

	for _, child := range g.Children {
		if !child.IsFolder() {
			util.Log(util.LogInfo, "GenFolder.ToNetscapeSubfolder() skipped", g.Name+"/"+child.Name, g.Guid+"/"+child.Guid, len(child.Children))
			continue
		}

		folders = append(folders, child.makeFolder("GenFolder.ToNetscapeSubfolder()"))
		subDirs = append(subDirs, child.Name)
	}

	util.Log(util.LogInfo, "GenFolder.ToNetscapeSubfolder() found", g.Name, g.Guid, len(folders), subDirs)
	return folders
}

func (g *GenChild) ToNetscape(c []GenChild) netscape.Folder {
	g.Children = c // why the fuck do i have to copy this literally what
	return g.makeFolder("GenChild.ToNetscape()")
}

func (g GenRoot) ToNetscape(genSubFolders *GenSubfolders) []netscape.Folder {
	netSubFolders := make([]netscape.Folder, 0)

	for _, folder := range *genSubFolders {
		if folder != nil {
			netSubFolders = append(netSubFolders, folder.ToNetscape(folder.Children))
		}
	}

	return netSubFolders
}
