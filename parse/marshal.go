package parse

import (
	"fmt"

	"github.com/5HT2C/chrome-bookmarks-converter/utils"
	"github.com/virtualtam/netscape-go/v2"
)

func (g Gen) ToNetscape() *netscape.Document {
	g.PopulateOrigin()

	return &netscape.Document{
		Title: fmt.Sprintf("Bookmarks-%s-%s", g.Origin, g.Checksum),
		Root: netscape.Folder{
			Subfolders: g.Roots.ToNetscape(),
		},
	}
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

func (g GenChild) Description() string {
	sep := " - "
	if len(g.Type) == 0 && len(g.Guid) == 0 {
		sep = ""
	}

	return fmt.Sprintf("%s%s%s", g.Type, sep, g.Guid)
}

func (g GenChild) ToNetscape() netscape.Bookmark {
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
		bookmarks = append(bookmarks, child.ToNetscape())
	}

	return bookmarks
}

func (g *GenFolder) ToNetscape() netscape.Folder {
	return netscape.Folder{
		CreatedAt:   utils.StringToTime(g.DateAdded),
		UpdatedAt:   utils.StringToTime(g.DateModified),
		Description: g.Description(),
		Name:        g.Name,
		Attributes:  g.AttrStr(nil),
		Bookmarks:   g.Bookmarks(),
	}
}

func (g GenRoot) ToNetscape() []netscape.Folder {
	netSubFolders := make([]netscape.Folder, 0)

	for _, folder := range []*GenFolder{g.BookmarkBar, g.Synced, g.Other} {
		if folder != nil {
			netSubFolders = append(netSubFolders, folder.ToNetscape())
		}
	}

	return netSubFolders
}
