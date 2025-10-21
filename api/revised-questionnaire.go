package api

import (
	"log"

	"github.com/labstack/echo/v4"
	// "miami-back-end/revised-questionnaire"
	"miami-back-end/members"
	"miami-back-end/pilot-questionnaire"
	"miami-back-end/project"
	"miami-back-end/revised-questionnaire"
	"miami-back-end/stage"
	"net/http"
	"strconv"
)
	


 type RQRoute struct{
	RQService *revisedquestionnaire.RQService
	PQService *pilotquestionnaire.PilotQuestionnaireService
	StageService *stage.StageService
	ProjectService *project.Service
	MemberService *members.MemberService
	// SOService *questionnairesignoff.QuestionnaireSignOffService



 }

 func NewRQRoute(RQService *revisedquestionnaire.RQService,PQService *pilotquestionnaire.PilotQuestionnaireService,
	StageService *stage.StageService,ProjectService *project.Service,MemberService *members.MemberService) *RQRoute {
	return &RQRoute{
		RQService: RQService,
		PQService: PQService,
		StageService: StageService,
		ProjectService: ProjectService,
		MemberService: MemberService,
		// SOService: SOService,
	}
 }

 func (r *RQRoute)Group(g *echo.Group){
	g.Use(Auth())
	g.GET("/:projectId", r.getRQInfo)
	g.GET("/detail-on-change/:projectId", r.getRQDetailOnChange)
	g.POST("/create-detail", r.createRQDetail)
	g.POST("/create-question", r.createRQQuestion)
}

func (r *RQRoute)getRQInfo(c echo.Context) error {
	projectId := c.Param("projectId")
	info,err := r.RQService.GetRevisedQuestionnaire(projectId)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK,info)
 }

 func (r *RQRoute)getRQDetailOnChange(c echo.Context) error {
	projectId := c.Param("projectId")
	stage := "5_1"
	info,err := r.RQService.GetRevisedQuestionnaireDetailOnChange(projectId,stage)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return  c.JSON(http.StatusOK,info)
 }

 func (r *RQRoute)createRQDetail(c echo.Context) error {

	var rqd pilotquestionnaire.PilotQuestionnaire
	if err := c.Bind(&rqd); err != nil {
		return err
	}
	stage := "5_1"
	check,err := r.RQService.GetRevisedQuestionnaireDetailOnChange(rqd.ProjectID,stage)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if (check == revisedquestionnaire.DetailOnChange{}) {
		err := r.RQService.RQRDetailOnChange(&rqd,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		
	}else{
		err := r.PQService.DetailOnChange(&rqd,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return  c.JSON(http.StatusCreated,rqd)
 }

 func (r *RQRoute)createRQQuestion(c echo.Context) error {
	var rq pilotquestionnaire.Path
	if err := c.Bind(&rq); err != nil {
		return err
	}
	ap,err := r.PQService.GetAllPilotQuestionnairePath(rq.ProjectID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fileNumber := len(ap)
	err = r.RQService.InsertRevisedQuestionnaire(&rq,fileNumber)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	projectID, err := strconv.Atoi(rq.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	member,err := r.MemberService.GetMembersOfPj(rq.ProjectID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var mbs []members.MembersOfPj
	var PMName string

	for _, mb := range member{
		if mb.Position == "PO" {
			PMName = mb.Firstname + " " + mb.Lastname
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

	if len(mbs) > 0 && len(rq.FilePath) > 0 {
		err := r.RQService.AlertRevisedQuestionnaire(rq.FilePath,&mbs,PMName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	stage := "stage_5_1"
	err = r.StageService.UpdateStage(rq.ProjectID,stage)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return  c.JSON(http.StatusCreated,rq)
 }