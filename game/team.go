package game

import (
	"fmt"
)

type Team struct {
	Name    string
	Score   uint
	Attacks uint
	Players []*Player
	Bench   []*Player
}

func (t *Team) Scores(score uint, player, assist *Player) {
	t.Score = t.Score + score
	player.Score = player.Score + score
	if player != assist {
		assist.Assist += score
	}
}

func (t *Team) GetTopScorer() (player *Player) {
	for _, p := range append(t.Players, t.Bench...) {
		if player == nil || p.Score > player.Score {
			player = p
		}
	}
	return
}

func (t *Team) GetTopAssist() (player *Player) {
	for _, p := range append(t.Players, t.Bench...) {
		if player == nil || p.Assist > player.Assist {
			player = p
		}
	}
	return
}

func (t *Team) String() string {
	return t.Name
}

func NewTeam(name string) *Team {
	t := &Team{
		Name: name,
	}
	t.Players, t.Bench = make([]*Player, MaxPlayers), make([]*Player, MaxBench)
	for i := 0; i < MaxPlayers; i++ {
		t.Players[i] = NewPlayer(fmt.Sprintf("Player %d (%s)", i+1, name))
	}
	for i := 0; i < MaxBench; i++ {
		t.Bench[i] = NewPlayer(fmt.Sprintf("Player %d (%s)", MaxPlayers+i+1, name))
	}
	return t
}
