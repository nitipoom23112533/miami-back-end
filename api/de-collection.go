package api

import(
	"github.com/labstack/echo/v4"
	"miami-back-end/de-collection"
	"miami-back-end/data-collection"
	"log"
	"net/http"
	"strconv"
	"time"
	"miami-back-end/members"
	"miami-back-end/project"
	"miami-back-end/stage"
	"gopkg.in/guregu/null.v4"

)

type DeCRoute struct{
	DeCollectionService *decollection.DeCollectionService
	ProjectService *project.Service
	MemberService *members.MemberService
	StageService *stage.StageService

}

func NewDeCRoute(DeCollectionService *decollection.DeCollectionService,ProjectService *project.Service,MemberService *members.MemberService,
	StageService *stage.StageService) *DeCRoute{
	return &DeCRoute{
		DeCollectionService: DeCollectionService,
		ProjectService: ProjectService,
		MemberService: MemberService,
		StageService: StageService,
	}
}

func (r *DeCRoute)Group(g *echo.Group){
	g.GET("/:projectId",r.getDeCollectionInfo)
	g.GET("/de-tasks/:projectId",r.getDeCollection)
	g.POST("/start",r.startDeCollection)
	g.PATCH("/completed",r.completedDeCollection)
}

func (r *DeCRoute)getDeCollectionInfo(c echo.Context) error {
	projectId := c.Param("projectId")
	DeCollection,err := r.DeCollectionService.GetDeCollectionInfo(projectId)
	if err != nil {
		log.Println("err",err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, DeCollection)
	
}

func (r *DeCRoute)getDeCollection(c echo.Context) error {
	projectId := c.Param("projectId")
	deCollection,err := r.DeCollectionService.GetDeCollection(projectId)
	if err != nil {
		log.Println("err",err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, deCollection)
	
}

func (r *DeCRoute)startDeCollection(c echo.Context) error {
	var dct datacollection.DataCollection
	if err := c.Bind(&dct); err != nil {
		log.Println(err)
		return err
	}
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	newDeCollection := datacollection.DataCollection{
		ProjectID: dct.ProjectID,
		IsStart: dct.IsStart,
		Quota: dct.Quota,
		StartDate: dct.StartDate,
		CreatedBy: userLogin.Email,
		CreatedAt: time.Now(),

	}
	err := r.DeCollectionService.StartDeCollection(&newDeCollection)
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
	var DeName []string

	for _, mb := range member{
		if mb.Position == "PO" {
			PMName = mb.Firstname + " " + mb.Lastname
		}
		if mb.Role == "Data Entry Manager" {
			DeName = append(DeName, mb.Firstname + " " + mb.Lastname) 
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
		err := r.DeCollectionService.AlertStartDecollection(&mbs,&newDeCollection,PMName,DeName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_12"
		err = r.StageService.UpdateStage(dct.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

	}
	return c.JSON(http.StatusCreated,dct)

}

func (r *DeCRoute)completedDeCollection(c echo.Context) error {
	var dct datacollection.DataCollection
	if err := c.Bind(&dct); err != nil {
		log.Println(err)
		return err
	}
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	newDeCollection := datacollection.DataCollection{
		ProjectID: dct.ProjectID,
		IsCompleted: dct.IsCompleted,		
		CompletedDate: dct.CompletedDate,
		UpdatedAt: null.TimeFrom(time.Now()),
		UpdatedBy: null.StringFrom(userLogin.Email),

	}
	err := r.DeCollectionService.CompletedDeCollection(&newDeCollection)
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
	var DeName []string

	for _, mb := range member{
		if mb.Position == "PO" {
			PMName = mb.Firstname + " " + mb.Lastname
		}
		if mb.Role == "Data Entry Manager" {
			DeName = append(DeName, mb.Firstname + " " + mb.Lastname) 
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
		err := r.DeCollectionService.AlertCompletedDecollection(&mbs,&newDeCollection,PMName,DeName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_13"
		err = r.StageService.UpdateStage(dct.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

	}
	return c.JSON(http.StatusCreated,dct)

}