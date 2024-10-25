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
	Source       string `json:"source,omitempty" gorm:"column:source,omitempty"`
	ShowIcon     bool   `json:"show_icon,omitempty" gorm:"column:show_icon,omitempty"`
	MetaInfo     struct {
		PowerBookmarkMeta string `json:"power_bookmark_meta,omitempty" gorm:"column:power_bookmark_meta,omitempty"`
	} `json:"meta_info,omitempty" gorm:"column:meta_info,omitempty"`
}

type GenFolder struct {
	GenChild
	Children []GenChild `json:"children,omitempty" gorm:"column:children,omitempty"`
}

type GenRoot struct {
	Other       *GenFolder `json:"other,omitempty" gorm:"column:other,omitempty"`
	Synced      *GenFolder `json:"synced,omitempty" gorm:"column:synced,omitempty"`
	BookmarkBar *GenFolder `json:"bookmark_bar,omitempty" gorm:"column:bookmark_bar,omitempty"`
}

type Gen struct {
	Checksum string  `json:"checksum" gorm:"column:checksum"`
	Roots    GenRoot `json:"roots" gorm:"column:roots"`
	Version  int     `json:"version" gorm:"column:version"`
}
