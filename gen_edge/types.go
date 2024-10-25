package gen_edge

// GenChild represents one bookmark.
// DateModified is only used for GenFolder
type GenChild struct {
	DateAdded    string `json:"date_added" gorm:"column:date_added"`
	DateModified string `json:"date_modified" gorm:"column:date_modified"`
	DateLastUsed string `json:"date_last_used" gorm:"column:date_last_used"`
	Name         string `json:"name" gorm:"column:name"`
	Guid         string `json:"guid" gorm:"column:guid"`
	ID           string `json:"id" gorm:"column:id"`
	Type         string `json:"type" gorm:"column:type"`
	Url          string `json:"url" gorm:"column:url"`
	Source       string `json:"source" gorm:"column:source"`
	ShowIcon     bool   `json:"show_icon" gorm:"column:show_icon"`
}

type GenFolder struct {
	GenChild
	Children []GenChild `json:"children,omitempty" gorm:"column:children,omitempty"`
}

type Gen struct {
	Checksum string `json:"checksum" gorm:"column:checksum"`
	Roots    struct {
		Other       *GenFolder `json:"other,omitempty" gorm:"column:other,omitempty"`
		Synced      *GenFolder `json:"synced,omitempty" gorm:"column:synced,omitempty"`
		BookmarkBar *GenFolder `json:"bookmark_bar,omitempty" gorm:"column:bookmark_bar,omitempty"`
	} `json:"roots" gorm:"column:roots"`
	Version int `json:"version" gorm:"column:version"`
}
