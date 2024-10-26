package parse

import (
	"fmt"
	"strings"
	"time"

	"github.com/5HT2C/chrome-bookmarks-converter/util"
	"github.com/virtualtam/netscape-go/v2"
)

var (
	rootTime   = time.Now()
	rootName   = fmt.Sprintf("root-folder-%v", rootTime.UnixNano())
	rootFolder = GenChild{
		Name:         rootName,
		DateAdded:    fmt.Sprintf("%v", rootTime.Unix()),
		DateModified: fmt.Sprintf("%v", rootTime.Unix()),
		Type:         "folder",
		Source:       GenOriginImported.String(),
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
			Description: g.Description(genTitle),
			Subfolders:  g.Roots.ToNetscape(g.PopulateRoots()),
		},
	}
}

func (g *GenChild) CollectChildren() GenChildren {
	children := make([]*GenChild, 0)

	if !g.IsFolder() {
		children = append(children, g)
	}

	for _, child := range g.Children {
		if child.IsFolder() {
			children = append(children, child.CollectChildren()...) // wheeeee recursion :3
			continue
		}

		children = append(children, &child)
	}

	return children
}

func (g *Gen) CollectChildren() GenChildren {
	children := make(GenChildren, 0)

	for _, folder := range *g.PopulateRoots() {
		folder.Children = append(make([]GenChild, 0), folder.Children...) // unnecessary imo
		children = append(children, folder.CollectChildren()...)
	}

	return children
}

func (g *GenChildren) ToNetscapeUnique() *netscape.Document {
	rootChildren := make([]GenChild, 0)
	r := &Gen{Roots: GenRoot{}}
	r.Origin = GenOriginImported

	for _, child := range *g {
		rootChildren = append(rootChildren, *child)
	}

	r.Roots.BookmarkBar = &GenFolder{rootFolder, rootChildren}

	genTitle := fmt.Sprintf(
		"Bookmarks-%s%s",
		r.Origin,
		util.StringConditional("-"+r.Checksum, "", !util.StringEmpty(r.Checksum)),
	)
	util.Log(util.LogInfo, "Gen.ToNetscape() title", genTitle)

	return &netscape.Document{ // TODO: Cleanup and use Gen.ToNetscape()
		Title: genTitle,
		Root: netscape.Folder{
			Name:        genTitle,
			Description: r.Description(genTitle),
			Subfolders:  r.Roots.ToNetscape(r.PopulateRoots()),
		},
	}
}

func (g *GenChildren) CollectUnique() (GenChildren, GenChildren, GenChildren) {
	dupeUrls := make(map[string]map[string]GenChildren) // [url][name]GenChildren

	for _, child := range *g {
		// Check for duplicate URL + Name
		if bkDupe, ok := dupeUrls[child.Url]; ok {
			if bkDupeName, ok := bkDupe[child.Name]; ok { // [url] exists
				bkDupe[child.Name] = append(bkDupeName, child) // [url][name] exists
			} else {
				bkDupe[child.Name] = append(make(GenChildren, 0), child) // [url][name] doesn't exist
			}
		} else { // [url] doesn't exist
			bkDupe = make(map[string]GenChildren)
			bkDupeName := make(GenChildren, 0)
			bkDupeName = append(bkDupeName, child)

			bkDupe[child.Name] = bkDupeName
			dupeUrls[child.Url] = bkDupe
		}
	}

	bkUniqueCol := make(GenChildren, 0)
	bkExistsCol := make(GenChildren, 0)
	bkPreferred := make(GenChildren, 0)

	for childUrl, children := range dupeUrls {
		for childName, childrenDedupe := range children {
			if len(childrenDedupe) > 1 {
				bkExistsCol = append(bkExistsCol, childrenDedupe...)

				var highestB *GenChild
				var highestF *GenChild
				scoreB := int64(0)
				scoreF := int64(0)

				for _, child := range childrenDedupe {
					score := child.DoWeWantToKeepThisScore()
					if child.IsFolder() {
						if highestF == nil || score > scoreF {
							highestF = child
							scoreF = score
							continue
						}
					} else {
						if highestB == nil || score > scoreB {
							highestB = child
							scoreB = score
						}
					}
				}

				for scorePref, childPref := range map[int64]*GenChild{scoreB: highestB, scoreF: highestF} {
					if childPref == nil {
						continue
					}

					bkPreferred = append(bkPreferred, highestB)
					util.Log(util.LogInfo, "GenChildren.CollectUnique() preferred", scorePref, childPref.Name, childPref.Guid)
				}

				util.Log(util.LogWarn, "GenChildren.CollectUnique() duplicate", len(childrenDedupe), childUrl, childName)
			} else {
				bkUniqueCol = append(bkUniqueCol, childrenDedupe...)

				util.Log(util.LogInfo, "GenChildren.CollectUnique() unique", len(childrenDedupe), childUrl, childName)
			}
		}
	}

	return bkUniqueCol, bkExistsCol, bkPreferred
}

func (g GenOrigin) String() string {
	switch g {
	case GenOriginUnknown:
		fallthrough
	case GenOriginChrome:
		return "chrome"
	case GenOriginEdge:
		return "edge"
	case GenOriginImported:
		return "imported"
	case GenOriginOther:
		return "other"
	default:
		return "unknown"
	}
}

func (g *Gen) Description(title string) string {
	return strings.Join([]string{
		util.StringOrDefault(
			title, util.StringOrDefault(
				g.Checksum, g.Origin.String(),
			),
		),
	}, " - ")
}

func (g *GenChild) Description(d string) string {
	return strings.Join([]string{g.Type, util.StringOrDefault(g.Guid, d)}, " - ")
}

func (g *GenChild) makeBookmark() netscape.Bookmark {
	b := netscape.Bookmark{
		Title:       g.Name,
		URL:         g.Url,
		Description: g.Description(g.Source),
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
		Description: g.Description(g.Source),
		Name:        g.Name,
		Attributes:  g.AttrStr(nil),
		Bookmarks:   g.ToNetscapeBookmarks(),  // we fucking love recurse around here
		Subfolders:  g.ToNetscapeSubfolder(g), // whoooo yeah!!!! recursive loops!!!!!
	}

	util.Log(util.LogInfo, caller+" made folder", g.Name, g.Guid)
	return f
}

func (g *GenChild) DoWeWantToKeepThisScore() (n int64) {
	n += int64(len(g.Children))
	n += util.StringEmptyScore(string(util.StringOrDefault(g.Guid, "0")[0]))
	n += util.StringEmptyScore(string(util.StringOrDefault(g.Type, "0")[0]))
	n += util.StringEmptyScore(g.DateAdded)
	n += util.StringEmptyScore(g.DateModified)
	n += util.StringEmptyScore(g.DateLastUsed)
	n += util.StringEmptyScore(g.Name)
	n += util.StringEmptyScore(g.Url)
	n += int64(-len(g.Source)) // lower score because it's probably imported from somewhere
	return n
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
