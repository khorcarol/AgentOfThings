package interests

type InterestCategory int

const (
	Sport InterestCategory = iota
	Music
	FilmTV
	Books
)

var stateName = map[InterestCategory]string{
	Sport:  "Sport",
	Music:  "Music",
	FilmTV: "FilmTV",
	Books:  "Books",
}

func String(ic InterestCategory) string {
	return stateName[ic]
}
