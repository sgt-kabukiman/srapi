// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

type Player struct {
	User  *User
	Guest *Guest
}

type PlayerLink struct {
	Relation string `json:"rel"`
	Id       string
	Name     string
	URI      string `json:"uri"`
}

func (self *PlayerLink) request(filter filter, sort *Sorting) request {
	relURL := self.URI[len(BaseUrl):]

	return request{"GET", relURL, filter, sort, nil}
}

type playerCollection struct {
	Data []map[string]interface{}
}
