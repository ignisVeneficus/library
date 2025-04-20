package api

import "ignis/library/server/db/dbo"

type Tag struct {
	TagId NullNumber `json:"id"`
	Name  string     `json:"name"`
	Color NullString `json:"color"`
}

func convertDBOTagToApiTag(dbo dbo.Tag) Tag {
	return Tag{
		TagId: NullNumber{dbo.TagId},
		Name:  dbo.Name,
		Color: NullString{dbo.Color},
	}
}
func convertApiTagToDBOTag(api Tag) dbo.Tag {
	return dbo.Tag{
		TagId: api.TagId.NullInt64,
		Name:  api.Name,
		Color: api.Color.NullString,
	}
}
