
type CommonInterest struct {
	Category    string
	Description string
}

type Peer struct {
	CommonInterests []CommonInterest
	Photo           string
	Name            string
}

// Function to get a list of peers with common interests
func getPeers() []CommonInterest {
	return []CommonInterest{
		{"Films/TV", "Sci-fi, Action movies"},
		{"Music", "Rock, Pop, Classical"},
	}
}

// Function to get a list of connected peers, which includes photo and name
func getConnectedPeers() []Peer {
	return []Peer{
		{
			CommonInterests: []CommonInterest{
				{"Films/TV", "Action movies, Thrillers"},
				{"Music", "Jazz, Blues"},
			},
			Photo: "http://example.com/photo1.jpg",
			Name:  "Alice",
		},
		{
			CommonInterests: []CommonInterest{
				{"Sports", "Soccer, Tennis"},
				{"Books", "Science fiction, Historical"},
			},
			Photo: "http://example.com/photo2.jpg",
			Name:  "Bob",
		},
	}
}