package main

var (
	Games = map[string]*Game{}
)

type Game struct {
	Candidates []*Candidate `json:"candidates"`
}

type Candidate struct {
	Name  string `json:"name"`
	EMail string `json:"email"`
	Photo string `json:"photo"`
}

func NewGame() *Game {
	return &Game{
		Candidates: []*Candidate{},
	}
}

func (g *Game) AddCandidate(c *Candidate) {
	// 중복 체크
	for _, v := range g.Candidates {
		if v.EMail == c.EMail {
			return
		}
	}

	g.Candidates = append(g.Candidates, c)
}

func (g *Game) RemoveCandidate(email string) {
	for i, v := range g.Candidates {
		if v.EMail == email {
			g.Candidates = append(g.Candidates[:i], g.Candidates[i+1:]...)
			return
		}
	}
}
