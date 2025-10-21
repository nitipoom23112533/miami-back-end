package api

import(
	"github.com/labstack/echo/v4"
	"miami-back-end/qaqc-collection"
	"miami-back-end/data-collection"
	"time"
	"log"
	"net/http"
	"strconv"
	"miami-back-end/members"
	"miami-back-end/project"
	"miami-back-end/stage"
	"gopkg.in/guregu/null.v4"
)

type QAqcRoute struct{
	QaqcCollectionService *qaqccollection.QaqcCollectionService
	ProjectService *project.Service
	MemberService *members.MemberService
	StageService *stage.StageService

}

func NewQAqcRoute(QaqcCollectionService *qaqccollection.QaqcCollectionService,ProjectService *project.Service,MemberService *members.MemberService,StageService *stage.StageService) *QAqcRoute{
	return &QAqcRoute{
		QaqcCollectionService: QaqcCollectionService,
		ProjectService: ProjectService,
		MemberService: MemberService,
		StageService: StageService,
	}
}

func (r *QAqcRoute)Group(g *echo.Group){
	// g.Use(Auth())
	g.GET("/:projectId",r.getDataCollection)
	g.GET("/fs-responses-qaqc/:projectId",r.getFsResponses)
	g.POST("/start",r.startQaqcCollection)
	g.PATCH("/completed",r.completedQaqcCollection)

}

func (r *QAqcRoute)getFsResponses(c echo.Context) error {
	projectId := c.Param("projectId")
	fsrqc,err := r.QaqcCollectionService.GetQcCollection(projectId)
	if err != nil {
		log.Println("err",err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fsrqa,err := r.QaqcCollectionService.GetQaCollection(projectId)
	if err != nil {
		log.Println("err",err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"qc_collection": fsrqc,
		"qa_collection": fsrqa,
	})
}

func (r *QAqcRoute)getDataCollection(c echo.Context) error {
	projectId := c.Param("projectId")
	dataCollection,err := r.QaqcCollectionService.GetQaqcCollection(projectId)
	if err != nil {
		log.Println("err",err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataCollection)
	
}

func (r *QAqcRoute)startQaqcCollection(c echo.Context) error {
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
	err := r.QaqcCollectionService.StartQaqcCollection(&newDataCollection)
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
	var QaqcName []string


	for _, mb := range member{
		if mb.Position == "PO" {
			PMName = mb.Firstname + " " + mb.Lastname
		}
		if mb.Role == "QC Manager" || mb.Role == "QA Manager" {
			QaqcName = append(QaqcName, mb.Firstname + " " + mb.Lastname) 
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
		err := r.QaqcCollectionService.AlertStartQaqccollection(&mbs,&newDataCollection,PMName,QaqcName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_9"
		err = r.StageService.UpdateStage(dct.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

	}
	return c.JSON(http.StatusCreated,dct)
}

func (r *QAqcRoute)completedQaqcCollection(c echo.Context) error {
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
	err := r.QaqcCollectionService.CompletedQaqcCollection(&newDataCollection)
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
	var QaqcName []string

	for _, mb := range member{
		if mb.Position == "PO" {
			PMName = mb.Firstname + " " + mb.Lastname
		}
		if mb.Role == "QC Manager" || mb.Role == "QA Manager" {
			QaqcName = append(QaqcName, mb.Firstname + " " + mb.Lastname) 
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
		err := r.QaqcCollectionService.AlertCompletedQaqccollection(&mbs,&newDataCollection,PMName,QaqcName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_10"
		err = r.StageService.UpdateStage(dct.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

	}
	return c.JSON(http.StatusCreated,dct)
}