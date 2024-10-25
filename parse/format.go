package parse

import (
	"github.com/5HT2C/chrome-bookmarks-converter/utils"
)

// PopulateOrigin will try to detect whether what browser this bookmark backup is from.
// This is kind of unreliable, and it's mostly a guess based on what kind of "default" bookmarks bar we have.
// If the user has both types of bookmark bars, we can just default to Chrome as we have no real way of knowing otherwise.
func (g Gen) PopulateOrigin() {
	gotBrowserType := func(child GenChild) bool {
		if child.IsChrome() {
			g.Origin = GenOriginChrome
			return true
		}

		if child.IsEdge() {
			g.Origin = GenOriginEdge
			return true
		}

		// We found a default folder but are unsure if it is edge or chrome, try setting other?
		if g.Origin == GenOriginUnknown && child.IsFolder() {
			g.Origin = GenOriginOther
		}

		return false
	}

	if g.Origin == GenOriginUnknown {
	l:
		for _, folder := range []*GenFolder{g.Roots.BookmarkBar, g.Roots.Synced, g.Roots.Other} {
			if gotBrowserType(folder.GenChild) {
				break l
			}

			for _, bookmark := range folder.Children {
				if gotBrowserType(bookmark) {
					break l
				}
			}
		}
	}
}
func (g GenChild) IsFolder() bool {
	return g.Type == "folder"
}

func (g GenChild) IsChrome() bool { // unreliable
	return g.IsFolder() && (g.Name == "Bookmarks Bar" || g.Name == "Bookmarks bar")
}

func (g GenChild) IsEdge() bool { // unreliable
	return g.IsFolder() && (g.Name == "Favorites Bar" || g.Name == "Favorites bar")
}

// AttrDefaultBar will return if the current folder is the default bookmark folder for either edge or chrome
func (g GenChild) AttrDefaultBar() bool {
	return g.IsFolder() && (g.IsChrome() || g.IsEdge())
}

func (g GenChild) AttrStr(m map[string]string) map[string]string {
	m = utils.MapAppend(
		m, []utils.MapCondition{
			{
				!utils.StringEmpty(g.DateAdded),
				"ADD_DATE",
				g.DateAdded,
			},
			{
				!utils.StringEmpty(g.DateModified),
				"LAST_MODIFIED",
				g.DateModified,
			},
			{
				utils.StringEmpty(g.DateModified) && !utils.StringEmpty(g.DateAdded),
				"LAST_MODIFIED",
				g.DateAdded,
			},
			{
				!utils.StringEmpty(g.DateLastUsed),
				"LAST_VISIT",
				g.DateLastUsed,
			},
			{len(g.Guid) > 0, "GUID", g.Guid},
			{g.AttrDefaultBar(), "PERSONAL_TOOLBAR_FOLDER", "true"},
		}...,
	)

	return m
}
