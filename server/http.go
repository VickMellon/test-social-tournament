package server

import (
	"github.com/VickMellon/test-social-tournament/controller"
	"github.com/gin-gonic/gin"
)

type WebServer struct {
	router     *gin.Engine
	controller *controller.Controller
}

func NewWebServer(controller *controller.Controller) *WebServer {
	w := &WebServer{
		router:     gin.Default(),
		controller: controller,
	}
	w.InitRoutes()
	return w
}

func (s *WebServer) InitRoutes() {

	// #1 /take?playerId=P1&points=300   takes 300 points from player P1 account
	s.router.GET("/take", getTakeHandler(s)) // should it be POST?

	// #1 /fund?playerId=P2&points=300   funds player P2 with 300 points. If no player exist should create new player
	s.router.GET("/fund", getFundHandler(s)) // should it be POST?

	// #2 /announceTournament?tournamentId=1&deposit=1000    Announce tournament specifying the entry deposit
	s.router.GET("/announceTournament", getAnnounceTournamentHandler(s)) // should it be POST?

	// #3 /joinTournament?tournamentId=1&playerId=P1&backerId=P2&backerId=P3     Join player into a tournament and is he backed by a set of backers
	s.router.GET("/joinTournament", getJoinTournamentHandler(s)) // should it be POST?

	// #4 /resultTournament Result tournament winners and prizes
	// with a POST document in format
	// {"tournamentId": "1", "winners": [{"playerId": "P1", "prize": 500}]}
	s.router.POST("/resultTournament", postResultTournamentHandler(s))

	// #5 /balance?playerId=P1  Player balance
	s.router.GET("/balance", getBalanceHandler(s))

	// #6 /reset Reset DB.
	s.router.GET("/reset", getResetHandler(s)) // should it be POST or DELETE?
}

func (s *WebServer) Run(port string) error {
	return s.router.Run(port)
}

func getResetHandler(s *WebServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Do action
		if err := s.controller.ResetDB(); err != nil {
			c.String(500, err.Error())
			return
		}

		// OK
		c.String(200, "")
	}
}
