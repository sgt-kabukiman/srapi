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

func (self *PlayerLink) request() request {
	relURL := self.URI[len(BaseUrl):]

	return request{"GET", relURL, nil, nil, nil}
}

type playerCollection struct {
	Data []map[string]interface{}
}
