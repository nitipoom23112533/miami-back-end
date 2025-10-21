package api

import (
	"github.com/labstack/echo/v4"
	"miami-back-end/da-collection"
	"miami-back-end/project"
	"miami-back-end/members"
	"miami-back-end/stage"
	"log"
	"net/http"
	"strconv"
	"time"
	"miami-back-end/data-collection"
	"gopkg.in/guregu/null.v4"
)

type DaCRoute struct{
	DaCollectionService *dacollection.DaCollectionService
	ProjectService *project.Service
	MemberService *members.MemberService
	StageService *stage.StageService

}

func NewDaCRoute( DaCollectionService *dacollection.DaCollectionService,ProjectService *project.Service,MemberService *members.MemberService,
	StageService *stage.StageService) *DaCRoute{
	return &DaCRoute{
		DaCollectionService: DaCollectionService,
		ProjectService: ProjectService,
		MemberService: MemberService,
		StageService: StageService,
	}
}

func (r *DaCRoute)Group(g *echo.Group){
	g.GET("/:projectId",r.getDaCollectionInfo)
	g.GET("/da-tasks/:projectId",r.getDaCollection)
	g.POST("/start",r.startDaCollection)
	g.PATCH("/completed",r.completedDaCollection)

}

func (r *DaCRoute)getDaCollectionInfo(c echo.Context) error {
	projectId := c.Param("projectId")
	DeCollection,err := r.DaCollectionService.GetDaCollectionInfo(projectId)
	if err != nil {
		log.Println("err",err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, DeCollection)
}

func (r *DaCRoute)getDaCollection(c echo.Context) error {
	projectId := c.Param("projectId")
	deCollection,err := r.DaCollectionService.GetDaCollection(projectId)
	if err != nil {
		log.Println("err",err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, deCollection)
}

func (r *DaCRoute)startDaCollection(c echo.Context) error {
	var dct datacollection.DataCollection
	if err := c.Bind(&dct); err != nil {
		log.Println(err)
		return err
	}
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	newDaCollection := datacollection.DataCollection{
		ProjectID: dct.ProjectID,
		IsStart: dct.IsStart,
		Quota: dct.Quota,
		StartDate: dct.StartDate,
		CreatedBy: userLogin.Email,
		CreatedAt: time.Now(),

	}
	err := r.DaCollectionService.StartDaCollection(&newDaCollection)
	if err != nil {
		log.Println(err)
		return err
	}

	projectID, err := strconv.Atoi(dct.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	member,err := r.MemberService.GetMembersOfPj(dct.ProjectID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var mbs []members.MembersOfPj
	var PMName string
	var DaName []string


	for _, mb := range member{
		if mb.Position == "PO" {
			PMName = mb.Firstname + " " + mb.Lastname
		}
		if mb.Role == "Data Analyst" {
			DaName = append(DaName, mb.Firstname + " " + mb.Lastname) 
		}
		mbs = append(mbs, members.MembersOfPj{
			ProjectID:   mb.ProjectID,
			UID:         mb.UID,
			Firstname:   mb.Firstname,
			Lastname:    mb.Lastname,
			Email:       mb.Email,
			IsSendEmail: true,
		})
	}

	if len(mbs) > 0 {
		err := r.DaCollectionService.AlertStartDacollection(&mbs,&newDaCollection,PMName,DaName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_14"
		err = r.StageService.UpdateStage(dct.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

	}
	return c.JSON(http.StatusCreated,dct)
}

func (r *DaCRoute)completedDaCollection(c echo.Context) error {
	var dct datacollection.DataCollection
	if err := c.Bind(&dct); err != nil {
		log.Println(err)
		return err
	}
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	newDaCollection := datacollection.DataCollection{
		ProjectID: dct.ProjectID,
		IsCompleted: dct.IsCompleted,		
		CompletedDate: dct.CompletedDate,
		UpdatedAt: null.TimeFrom(time.Now()),
		UpdatedBy: null.StringFrom(userLogin.Email),

	}
	err := r.DaCollectionService.CompletedDaCollection(&newDaCollection)
	if err != nil {
		log.Println(err)
		return err
	}

	projectID, err := strconv.Atoi(dct.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	member,err := r.MemberService.GetMembersOfPj(dct.ProjectID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var mbs []members.MembersOfPj
	var PMName string
	var DaName []string

	for _, mb := range member{
		if mb.Position == "PO" {
			PMName = mb.Firstname + " " + mb.Lastname
		}
		if mb.Role == "Data Analyst" {
			DaName = append(DaName, mb.Firstname + " " + mb.Lastname) 
		}
		mbs = append(mbs, members.MembersOfPj{
			ProjectID:   mb.ProjectID,
			UID:         mb.UID,
			Firstname:   mb.Firstname,
			Lastname:    mb.Lastname,
			Email:       mb.Email,
			IsSendEmail: true,
		})
	}

	if len(mbs) > 0 {
		err := r.DaCollectionService.AlertCompletedDacollection(&mbs,&newDaCollection,PMName,DaName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_15"
		err = r.StageService.UpdateStage(dct.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

	}
	return c.JSON(http.StatusCreated,dct)
}