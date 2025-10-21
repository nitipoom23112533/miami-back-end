package api

import (
	"log"
	"miami-back-end/stage"
	"github.com/labstack/echo/v4"
	"net/http"
)

type StageRoute struct{
	StageService *stage.StageService
}

func NewStageRoute(stageService *stage.StageService) *StageRoute {
	return &StageRoute{
		StageService: stageService,
	}
}

func (r *StageRoute)Group(g *echo.Group){
	g.Use(Auth())
	g.GET("/:projectId", r.getStage)
}

func (r *StageRoute)getStage(c echo.Context) error {
	projectId := c.Param("projectId")
	stage, err := r.StageService.GetStage(projectId)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, stage)
}