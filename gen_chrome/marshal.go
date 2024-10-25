package gen_chrome

import (
	"fmt"

	"github.com/5HT2C/chrome-bookmarks-converter/utils"
	"github.com/virtualtam/netscape-go/v2"
)

func (g *Gen) ToNetscape() netscape.Document {
	return netscape.Document{
		Title: fmt.Sprintf("Bookmarks-%v-%s", g.Version, g.Checksum),
		Root: netscape.Folder{
			Subfolders: g.Roots.ToNetScape(),
		},
	}
}

func (g GenChild) Description() string {
	sep := " - "
	if len(g.Type) == 0 && len(g.Guid) == 0 {
		sep = ""
	}

	return fmt.Sprintf("%s%s%s", g.Type, sep, g.Guid)
}

func (g GenChild) ToNetScape() netscape.Bookmark {
	return netscape.Bookmark{
		CreatedAt:   utils.StringToTime(g.DateAdded),
		UpdatedAt:   utils.StringToTime(g.DateModified),
		Title:       g.Name,
		URL:         g.Url,
		Description: g.Description(),
		Private:     false,
		Tags:        nil,
		Attributes:  g.AttrStr(nil),
	}
}

func (g *GenFolder) Bookmarks() []netscape.Bookmark {
	bookmarks := make([]netscape.Bookmark, 0)
	for _, child := range g.Children {
		bookmarks = append(bookmarks, child.ToNetScape())
	}

	return bookmarks
}

func (g *GenFolder) ToNetScape() netscape.Folder {
	return netscape.Folder{
		CreatedAt:   utils.StringToTime(g.DateAdded),
		UpdatedAt:   utils.StringToTime(g.DateModified),
		Description: g.Description(),
		Name:        g.Name,
		Attributes:  g.AttrStr(nil),
		Bookmarks:   g.Bookmarks(),
	}
}

func (g *GenRoot) ToNetScape() []netscape.Folder {
	netSubFolders := make([]netscape.Folder, 0)

	for _, folder := range []*GenFolder{g.BookmarkBar, g.Synced, g.Other} {
		if folder != nil {
			netSubFolders = append(netSubFolders, folder.ToNetScape())
		}
	}

	return netSubFolders
}
