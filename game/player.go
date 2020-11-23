package game

type Player struct {
	Name   string
	Assist uint
	Score  uint
}

func NewPlayer(name string) *Player {
	return &Player{
		Name: name,
	}
}

func (p *Player) String() string {
	return p.Name
}
