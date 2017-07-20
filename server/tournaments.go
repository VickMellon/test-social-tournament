package server

import (
	"github.com/VickMellon/test-social-tournament/controller"
	"github.com/VickMellon/test-social-tournament/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"strconv"
)

type resultTournamentRequest struct {
	TournamentId string          `json:"tournamentId"`
	Winners      []*model.Winner `json:"winners"`
}

func getAnnounceTournamentHandler(s *WebServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tournamentId string
		var deposit float64
		var err error

		// parse params
		tournamentId, ok := c.GetQuery("tournamentId")
		if !ok || tournamentId == "" {
			c.String(400, "tournamentId required")
			return
		}
		depositString, ok := c.GetQuery("deposit")
		if !ok || depositString == "" {
			c.String(400, "deposit required")
			return
		}
		if deposit, err = strconv.ParseFloat(depositString, 64); err != nil || deposit == 0 {
			c.String(400, "invalid deposit value")
			return
		}

		// Do action
		if err := s.controller.CreateTournament(tournamentId, deposit); err != nil {
			if err == controller.ErrTournamentExists {
				c.String(409, err.Error())
			} else {
				c.String(500, err.Error())
			}
			return
		}

		// OK
		c.String(200, "")
	}
}

func getJoinTournamentHandler(s *WebServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tournamentId string
		var playerId string
		var backerIds []string

		// parse params
		tournamentId, ok := c.GetQuery("tournamentId")
		if !ok || tournamentId == "" {
			c.String(400, "tournamentId required")
			return
		}
		playerId, ok = c.GetQuery("playerId")
		if !ok || playerId == "" {
			c.String(400, "playerId required")
			return
		}
		// bakers
		backerIds, _ = c.GetQueryArray("backerId")

		// Do action
		if err := s.controller.JoinTournament(tournamentId, playerId, backerIds); err != nil {
			if err == controller.ErrPlayerNotFound {
				c.String(404, "Tournament or player or backer not found")
			} else if err == controller.ErrLowBalance {
				c.String(409, "Player or some backer haven't enough points to join this tournament")
			} else if err == controller.ErrAlreadyJoined {
				c.String(409, err.Error())
			} else {
				c.String(500, err.Error())
			}
			return
		}

		// OK
		c.String(200, "")
	}
}

func postResultTournamentHandler(s *WebServer) gin.HandlerFunc {
	return func(c *gin.Context) {

		// parse params
		var req resultTournamentRequest
		if err := binding.JSON.Bind(c.Request, &req); err != nil {
			c.String(400, "invalid request Json")
			return
		}
		if req.TournamentId == "" {
			c.String(400, "tournamentId required")
			return
		}
		if len(req.Winners) == 0 {
			c.String(400, "winners required")
			return
		}

		// Do action
		if err := s.controller.FinishTournament(req.TournamentId, req.Winners); err != nil {
			if err == controller.ErrPlayerNotFound {
				c.String(404, "Tournament or player or backer not found")
			} else if err == controller.ErrTournamentFinished {
				c.String(409, "Tournament already finished")
			} else {
				c.String(500, err.Error())
			}
			return
		}

		// OK
		c.String(200, "")
	}
}
