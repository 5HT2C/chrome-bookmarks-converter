package gen_chrome

type GenChild struct {
	DateAdded    string `json:"date_added" gorm:"column:date_added"`
	DateLastUsed string `json:"date_last_used" gorm:"column:date_last_used"`
	Name         string `json:"name" gorm:"column:name"`
	Guid         string `json:"guid" gorm:"column:guid"`
	ID           string `json:"id" gorm:"column:id"`
	Type         string `json:"type" gorm:"column:type"`
	Url          string `json:"url" gorm:"column:url"`
	MetaInfo     struct {
		PowerBookmarkMeta string `json:"power_bookmark_meta" gorm:"column:power_bookmark_meta"`
	} `json:"meta_info" gorm:"column:meta_info"`
}

type GenFolder struct {
	DateAdded    string     `json:"date_added" gorm:"column:date_added"`
	DateModified string     `json:"date_modified" gorm:"column:date_modified"`
	Children     []GenChild `json:"children,omitempty" gorm:"column:children,omitempty"`
	DateLastUsed string     `json:"date_last_used" gorm:"column:date_last_used"`
	Name         string     `json:"name" gorm:"column:name"`
	Guid         string     `json:"guid" gorm:"column:guid"`
	ID           string     `json:"id" gorm:"column:id"`
	Type         string     `json:"type" gorm:"column:type"`
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
