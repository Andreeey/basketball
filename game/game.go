package game

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	MaxPlayers = 5
	MaxBench   = 7
	//DurationSec = 20
	DurationSec   = 240
	AttackMaxMsec = 4800

	AttackSuccessPercentage      = 70
	Attack3ScorePercentage       = 20
	PlayerSubstitutionPercentage = 50
)

type Game struct {
	sync.Mutex

	Name       string
	Teams      []*Team
	TopScorer  *Player
	TopAssist  *Player
	attackTeam int

	topScoreChan chan *Game
	cancelFunc   context.CancelFunc
	log          *log.Logger
	rand         *rand.Rand
}

func (g *Game) Start() {
	g.log.Printf("Starting %s\n", g.Name)
	var ctx context.Context
	ctx, g.cancelFunc = context.WithTimeout(context.Background(), time.Second*DurationSec)
	go g.simulateAttacks(ctx)
}

func (g *Game) End() {
	g.log.Printf("Canceling %s\n", g.Name)
	g.cancelFunc()
}

func (g *Game) PrintResults() {
	g.log.Printf("Game has finished.\n")
	for _, t := range g.Teams {
		g.log.Printf("\t[%s] Attacks: %d Score: %d\n", t, t.Attacks, t.Score)
		for _, p := range append(t.Players, t.Bench...) {
			g.log.Printf("\t\t[%s] Score: %d Assist: %d\n", p, p.Score, p.Assist)
		}
	}
}

func (g *Game) simulateAttacks(ctx context.Context) {
	timer := time.NewTimer(time.Millisecond * time.Duration(g.rand.Intn(AttackMaxMsec)))

	for {
		select {
		case <-timer.C:
			g.simulateAttack()
			timer = time.NewTimer(time.Millisecond * time.Duration(g.rand.Intn(AttackMaxMsec)))
		case <-ctx.Done():
			timer.Stop()
			g.PrintResults()
			return
		}
	}
}

func (g *Game) simulateAttack() {
	g.Lock()
	defer g.Unlock()
	team := g.Teams[g.attackTeam]
	if g.rand.Intn(100) < AttackSuccessPercentage {
		// Attack succeeded
		var score uint = 2
		if g.rand.Intn(100) < Attack3ScorePercentage {
			score = 3
		}

		player := team.Players[g.rand.Intn(len(team.Players))]
		assist := team.Players[g.rand.Intn(len(team.Players))]
		team.Scores(score, player, assist)
		g.updateTopScorerAndAssist()

		g.log.Printf("%s attack succeeded with score: %d Player: %s Assist player: %s\n", team, score, player, assist)
	} else {
		g.log.Printf("%s attack failed\n", team)
	}
	team.Attacks++

	if g.rand.Intn(100) < PlayerSubstitutionPercentage {
		// Substitute player
		g.log.Println("Coach decided to substitute player")
		g.simulateSubstitution(team)
	}

	g.attackTeam++
	if g.attackTeam >= len(g.Teams) {
		g.attackTeam = 0
	}
}

func (g *Game) simulateSubstitution(team *Team) {
	playerKey := g.rand.Intn(len(team.Players))
	benchKey := g.rand.Intn(len(team.Bench))
	player := team.Players[playerKey]
	bench := team.Bench[benchKey]
	team.Players[playerKey] = bench
	team.Bench[benchKey] = player
	g.log.Printf("Player %s was substituted by %s\n", player, bench)
}

func (g *Game) getTopScorer() (player *Player) {
	for _, t := range g.Teams {
		p := t.GetTopScorer()
		if player == nil || p.Score > player.Score {
			player = p
		}
	}
	return
}

func (g *Game) getTopAssist() (player *Player) {
	for _, t := range g.Teams {
		p := t.GetTopAssist()
		if player == nil || p.Assist > player.Assist {
			player = p
		}
	}
	return
}

func (g *Game) updateTopScorerAndAssist() {
	var haveUpdate bool
	p := g.getTopScorer()
	if p != nil {
		if g.TopScorer == nil {
			g.TopScorer = p
			haveUpdate = true
		} else if p.Score > g.TopScorer.Score {
			g.TopScorer = p
			haveUpdate = true
		}
	}

	p = g.getTopAssist()
	if p != nil {
		if g.TopAssist == nil {
			g.TopAssist = p
			haveUpdate = true
		} else if p.Assist > g.TopAssist.Assist {
			g.TopAssist = p
			haveUpdate = true
		}
	}

	if haveUpdate {
		g.topScoreChan <- g
	}
}

func New(name string, topScoreChan chan *Game) *Game {
	game := Game{
		Name:         name,
		topScoreChan: topScoreChan,
		log:          log.New(os.Stdout, fmt.Sprintf("[Game %s]: ", name), log.Ltime),
		rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	game.Teams = make([]*Team, 2)
	game.Teams[0] = NewTeam("Team A")
	game.Teams[1] = NewTeam("Team B")

	return &game
}
