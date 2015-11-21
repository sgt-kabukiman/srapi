package srapi

type Level struct {
	Id      string
	Name    string
	Weblink string
	Rules   string
	Links   []Link
}

type levelResponse struct {
	Data Level
}

func LevelById(id string) (*Level, *Error) {
	return fetchLevel(request{"GET", "/levels/" + id, nil, nil, nil})
}

func fetchLevel(request request) (*Level, *Error) {
	result := &levelResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (self *Level) Game() *Game {
	link := firstLink(self, "game")
	if link == nil {
		return nil
	}

	game, _ := fetchGame(link.request())
	return game
}

// for the 'hasLinks' interface
func (self *Level) links() []Link {
	return self.Links
}

type LevelCollection struct {
	Data       []Level
	Pagination Pagination
}

func (self *LevelCollection) levels() []*Level {
	result := make([]*Level, 0)

	for idx := range self.Data {
		result = append(result, &self.Data[idx])
	}

	return result
}

func (self *LevelCollection) NextPage() (*LevelCollection, *Error) {
	return self.fetchLink("next")
}

func (self *LevelCollection) PrevPage() (*LevelCollection, *Error) {
	return self.fetchLink("prev")
}

func (self *LevelCollection) fetchLink(name string) (*LevelCollection, *Error) {
	next := firstLink(&self.Pagination, name)
	if next == nil {
		return nil, nil
	}

	return fetchLevels(next.request())
}

// always returns a collection, even when an error is returned;
// makes other code more monadic
func fetchLevels(request request) (*LevelCollection, *Error) {
	result := &LevelCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}
