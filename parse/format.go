package parse

import (
	"github.com/5HT2C/chrome-bookmarks-converter/util"
)

func (g *Gen) PopulateRoots() *GenSubfolders {
	return &GenSubfolders{g.Roots.BookmarkBar, g.Roots.Synced, g.Roots.Other}
}

// PopulateOrigin will try to detect whether what browser this bookmark backup is from.
// This is kind of unreliable, and it's mostly a guess based on what kind of "default" bookmarks bar we have.
// If the user has both types of bookmark bars, we can just default to Chrome as we have no real way of knowing otherwise.
func (g *Gen) PopulateOrigin() GenOrigin {
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
		for _, folder := range *g.PopulateRoots() {
			if gotBrowserType(folder.GenChild) {
				util.Log(util.LogInfo, "PopulateOrigin() gotBrowserType", g.Origin, folder.GenChild)
				break l
			}

			for _, bookmark := range folder.Children {
				if gotBrowserType(bookmark) {
					util.Log(util.LogInfo, "PopulateOrigin() gotBrowserType", g.Origin, folder.GenChild)
					break l
				}
			}
		}
	}

	return g.Origin
}

func (g *GenChild) IsFolder() bool {
	return g.Type == "folder"
}

func (g *GenChild) IsChrome() bool { // unreliable
	return g.IsFolder() && (g.Name == "Bookmarks Bar" || g.Name == "Bookmarks bar")
}

func (g *GenChild) IsEdge() bool { // unreliable
	return g.IsFolder() && (g.Name == "Favorites Bar" || g.Name == "Favorites bar")
}

// AttrDefaultBar will return if the current folder is the default bookmark folder for either edge or chrome
func (g *GenChild) AttrDefaultBar() bool {
	return g.IsFolder() && (g.IsChrome() || g.IsEdge())
}

// AttrStr will generate the attributes for the Netscape bookmark file format.
// Docs: https://learn.microsoft.com/en-us/previous-versions/windows/internet-explorer/ie-developer/platform-apis/aa753582(v=vs.85)
func (g *GenChild) AttrStr(m map[string]string) map[string]string {
	m = util.MapAppend(
		m, []util.MapCondition{
			{
				!util.StringEmpty(g.DateAdded),
				"ADD_DATE",
				g.DateAdded,
			},
			{
				!util.StringEmpty(g.DateLastUsed),
				"LAST_VISIT", // LAST_VISIT is supposed to be before LAST_MODIFIED, not sure why
				g.DateLastUsed,
			},
			{
				!util.StringEmpty(g.DateModified) || !util.StringEmpty(g.DateAdded),
				"LAST_MODIFIED",
				util.StringOrDefault(g.DateModified, g.DateAdded),
			},
			{len(g.Guid) > 0, "GUID", g.Guid},
			{g.AttrDefaultBar(), "PERSONAL_TOOLBAR_FOLDER", "true"},
		}...,
	)

	return m
}
