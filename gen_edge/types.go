package gen_edge

type GenChild struct {
	DateAdded    string `json:"date_added" gorm:"column:date_added"`
	DateLastUsed string `json:"date_last_used" gorm:"column:date_last_used"`
	ShowIcon     bool   `json:"show_icon" gorm:"column:show_icon"`
	Name         string `json:"name" gorm:"column:name"`
	Guid         string `json:"guid" gorm:"column:guid"`
	ID           string `json:"id" gorm:"column:id"`
	Type         string `json:"type" gorm:"column:type"`
	Url          string `json:"url" gorm:"column:url"`
	Source       string `json:"source" gorm:"column:source"`
}

type Gen struct {
	Checksum string `json:"checksum" gorm:"column:checksum"`
	Roots    struct {
		Other struct {
			DateAdded    string     `json:"date_added" gorm:"column:date_added"`
			DateModified string     `json:"date_modified" gorm:"column:date_modified"`
			Children     []GenChild `json:"children" gorm:"column:children"`
			DateLastUsed string     `json:"date_last_used" gorm:"column:date_last_used"`
			Name         string     `json:"name" gorm:"column:name"`
			Guid         string     `json:"guid" gorm:"column:guid"`
			ID           string     `json:"id" gorm:"column:id"`
			Source       string     `json:"source" gorm:"column:source"`
			Type         string     `json:"type" gorm:"column:type"`
		} `json:"other" gorm:"column:other"`
		Synced struct {
			DateAdded    string     `json:"date_added" gorm:"column:date_added"`
			DateModified string     `json:"date_modified" gorm:"column:date_modified"`
			Children     []GenChild `json:"children" gorm:"column:children"`
			DateLastUsed string     `json:"date_last_used" gorm:"column:date_last_used"`
			Name         string     `json:"name" gorm:"column:name"`
			Guid         string     `json:"guid" gorm:"column:guid"`
			ID           string     `json:"id" gorm:"column:id"`
			Source       string     `json:"source" gorm:"column:source"`
			Type         string     `json:"type" gorm:"column:type"`
		} `json:"synced" gorm:"column:synced"`
		BookmarkBar struct {
			DateAdded    string     `json:"date_added" gorm:"column:date_added"`
			DateModified string     `json:"date_modified" gorm:"column:date_modified"`
			Children     []GenChild `json:"children" gorm:"column:children"`
			DateLastUsed string     `json:"date_last_used" gorm:"column:date_last_used"`
			Name         string     `json:"name" gorm:"column:name"`
			Guid         string     `json:"guid" gorm:"column:guid"`
			ID           string     `json:"id" gorm:"column:id"`
			Source       string     `json:"source" gorm:"column:source"`
			Type         string     `json:"type" gorm:"column:type"`
		} `json:"bookmark_bar" gorm:"column:bookmark_bar"`
	} `json:"roots" gorm:"column:roots"`
	Version int `json:"version" gorm:"column:version"`
}
