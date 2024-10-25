package gen_chrome

import (
	"github.com/5HT2C/chrome-bookmarks-converter/utils"
)

func (g GenChild) IsFolder() bool {
	return g.Type == "folder"
}

func (g GenChild) AttrDefaultBar() bool {
	return g.IsFolder() &&
		(g.Name == "Favorites Bar" || g.Name == "Bookmarks Bar")
}

func (g GenChild) AttrStr(m map[string]string) map[string]string {
	m = utils.MapAppend(
		m, []utils.MapCondition{
			{
				g.DateAdded != "" && g.DateAdded != "0",
				"ADD_DATE",
				g.DateAdded,
			},
			{
				g.DateModified != "" && g.DateModified != "0",
				"LAST_MODIFIED",
				g.DateModified,
			},
			{len(g.Guid) > 0, "GUID", g.Guid},
			{g.AttrDefaultBar(), "PERSONAL_TOOLBAR_FOLDER", "true"},
		}...,
	)

	return m
}
