package api

import (
	"log"
	"miami-back-end/data-collection"
	"net/http"
	"time"
	"github.com/labstack/echo/v4"
	"strconv"
	"miami-back-end/project"
	"miami-back-end/members"
	"miami-back-end/stage"
	"gopkg.in/guregu/null.v4"
)

type DataCollectionRoute struct{
	DataCollectionService *datacollection.DataCollectionService
	ProjectService *project.Service
	MemberService *members.MemberService
	StageService *stage.StageService
}

func NewDataCollectionRoute(DataCollectionService *datacollection.DataCollectionService,ProjectService *project.Service,MemberService *members.MemberService,
	StageService *stage.StageService) *DataCollectionRoute{
	return &DataCollectionRoute{
		DataCollectionService: DataCollectionService,
		ProjectService: ProjectService,
		MemberService: MemberService,
		StageService: StageService,
	}
}

func (r *DataCollectionRoute)Group(g *echo.Group){
	g.Use(Auth())
	g.GET("/:projectId",r.getDataCollection)
	g.GET("/dashboard-logs/:projectId",r.getDashBoardLogs)
	g.GET("/ss-responses/:projectId",r.getSsResponses)
	g.GET("/fs-responses/:projectId",r.getFsResponses)
	g.POST("/start",r.startDataCollection)
	g.PATCH("/completed",r.completedDataCollection)

}

func (r *DataCollectionRoute)getDataCollection(c echo.Context) error {
	projectId := c.Param("projectId")
	info,err := r.DataCollectionService.GetDataCollection(projectId)
	if err != nil {
		log.Println("err",err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK,info)

}

func (r *DataCollectionRoute)getDashBoardLogs(c echo.Context) error {
	projectId := c.Param("projectId")
	logs,err := r.DataCollectionService.GetDashBoardLogs(projectId)
	if err != nil {
		log.Println("err",err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK,logs)
}

func (r *DataCollectionRoute)getSsResponses(c echo.Context) error {
	projectId := c.Param("projectId")
	logs,err := r.DataCollectionService.GetSsResponses(projectId)
	if err != nil {
		log.Println("err",err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK,logs)

}

func (r *DataCollectionRoute)getFsResponses(c echo.Context) error {
	projectId := c.Param("projectId")
	logs,err := r.DataCollectionService.GetFsResponses(projectId)
	if err != nil {
		log.Println("err",err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK,logs)
}

func (r *DataCollectionRoute)startDataCollection(c echo.Context) error {
	var dct datacollection.DataCollection
	if err := c.Bind(&dct); err != nil {
		log.Println(err)
		return err
	}
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	newDataCollection := datacollection.DataCollection{
		ProjectID: dct.ProjectID,
		IsStart: dct.IsStart,
		Quota: dct.Quota,
		Day: dct.Day,
		StartDate: dct.StartDate,
		CreatedBy: userLogin.Email,
		CreatedAt: time.Now(),

	}
	err := r.DataCollectionService.StartDataCollection(&newDataCollection)
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
	var FWName []string

	for _, mb := range member{
		if mb.Position == "PO" {
			PMName = mb.Firstname + " " + mb.Lastname
		}
		if mb.Role == "Fieldwork Manager" {
			FWName = append(FWName, mb.Firstname + " " + mb.Lastname) 
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
		err := r.DataCollectionService.AlertStartDatacollection(&mbs,&newDataCollection,PMName,FWName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_6"
		err = r.StageService.UpdateStage(dct.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

	}
	return c.JSON(http.StatusCreated,dct)
}

func (r *DataCollectionRoute)completedDataCollection(c echo.Context) error {
	var dct datacollection.DataCollection
	if err := c.Bind(&dct); err != nil {
		log.Println(err)
		return err
	}
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	newDataCollection := datacollection.DataCollection{
		ProjectID: dct.ProjectID,
		IsCompleted: dct.IsCompleted,		
		CompletedDate: dct.CompletedDate,
		UpdatedAt: null.TimeFrom(time.Now()),
		UpdatedBy: null.StringFrom(userLogin.Email),

	}
	err := r.DataCollectionService.CompletedDataCollection(&newDataCollection)
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
	var FWName []string


	for _, mb := range member{
		if mb.Position == "PO" {
			PMName = mb.Firstname + " " + mb.Lastname
		}
		if mb.Role == "Fieldwork Manager" {
			FWName = append(FWName, mb.Firstname + " " + mb.Lastname) 
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
		err := r.DataCollectionService.AlertCompletedDatacollection(&mbs,&newDataCollection,PMName,FWName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_7"
		err = r.StageService.UpdateStage(dct.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

	}
	return c.JSON(http.StatusCreated,dct)
}

