package server

import (
	"github.com/VickMellon/test-social-tournament/controller"
	"github.com/gin-gonic/gin"
	"strconv"
)

func getTakeHandler(s *WebServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var playerId string
		var points float64
		var err error

		// parse params
		playerId, ok := c.GetQuery("playerId")
		if !ok || playerId == "" {
			c.String(400, "playerId required")
			return
		}
		pointsString, ok := c.GetQuery("points")
		if !ok || pointsString == "" {
			c.String(400, "points required")
			return
		}
		if points, err = strconv.ParseFloat(pointsString, 64); err != nil || points <= 0 {
			c.String(400, "invalid points value")
			return
		}

		// Do action
		if err := s.controller.TakePoints(playerId, points); err != nil {
			if err == controller.ErrLowBalance {
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

func getFundHandler(s *WebServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var playerId string
		var points float64
		var err error

		// parse params
		playerId, ok := c.GetQuery("playerId")
		if !ok || playerId == "" {
			c.String(400, "playerId required")
			return
		}
		pointsString, ok := c.GetQuery("points")
		if !ok || pointsString == "" {
			c.String(400, "points required")
			return
		}
		if points, err = strconv.ParseFloat(pointsString, 64); err != nil || points <= 0 {
			c.String(400, "invalid points value")
			return
		}

		// Do action
		if err := s.controller.FundPoints(playerId, points); err != nil {
			c.String(500, err.Error())
			return
		}

		// OK
		c.String(200, "")
	}
}

func getBalanceHandler(s *WebServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var playerId string
		var err error

		// parse params
		playerId, ok := c.GetQuery("playerId")
		if !ok || playerId == "" {
			c.String(400, "playerId required")
			return
		}

		// Do action
		balance, err := s.controller.GetBalance(playerId)
		if err != nil {
			if err == controller.ErrPlayerNotFound {
				c.String(404, err.Error())
			} else {
				c.String(500, err.Error())
			}
			return
		}

		// OK
		c.JSON(200, gin.H{"playerId": playerId, "balance": balance})
	}
}
