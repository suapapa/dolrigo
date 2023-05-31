package main

var (
	Games = map[string]*Game{}
)

type Game struct {
	Candidates []*Candidate `json:"candidates"`
}

type Candidate struct {
	Name    string `json:"name"`
	EMail   string `json:"email"`
	Picture string `json:"picture"`
}

func NewGame() *Game {
	return &Game{
		Candidates: []*Candidate{},
	}
}

func (g *Game) AddCandidate(c *Candidate) {
	g.Candidates = append(g.Candidates, c)
}
