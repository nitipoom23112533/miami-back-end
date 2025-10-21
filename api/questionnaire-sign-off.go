package api

import (
	"log"
	"miami-back-end/members"
	"miami-back-end/pilot-questionnaire"
	"miami-back-end/questionnaire-sign-off"
	"miami-back-end/project"
	"miami-back-end/stage"
	"net/http"
	"github.com/labstack/echo/v4"
	"strconv"

)


type QuestionnaireSORoute struct{
	SOService *questionnairesignoff.QuestionnaireSignOffService
	MemberService *members.MemberService
	ProjectService *project.Service
	StageService *stage.StageService
	PQService *pilotquestionnaire.PilotQuestionnaireService
}
func NewQuestionnaireSORoute( SOService *questionnairesignoff.QuestionnaireSignOffService,
	MemberService *members.MemberService,ProjectService *project.Service,StageService *stage.StageService,PQService *pilotquestionnaire.PilotQuestionnaireService) *QuestionnaireSORoute{
	return &QuestionnaireSORoute{
		SOService: SOService,
		MemberService: MemberService,
		ProjectService: ProjectService,
		StageService: StageService,
		PQService: PQService,
	}
}

func (r *QuestionnaireSORoute)Group(g *echo.Group) {
	g.Use(Auth())
	g.GET("/:projectId", r.getQuestionnaireSignOff)
	g.GET("/detail-on-change/:projectId", r.getDetailOnChangeQuestionnaireSignOff)
	g.PATCH("/create",r.SOcreate)
	g.PATCH("/edit",r.SOedit)
	g.POST("/create-questionnaire-sign-off",r.createQuestionnaireSignOff)
	g.POST("/edit-questionnaire-sign-off",r.editQuestionnaireSignOff)
	g.PATCH("/edit-detail-on-change",r.SOeditDetailOnChange)
}

func (r *QuestionnaireSORoute)getQuestionnaireSignOff(c echo.Context) error {
	projectId := c.Param("projectId")
	info,err := r.SOService.GetAllPilotQuestionnairePath(projectId)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK,info)
}
func (r *QuestionnaireSORoute)getDetailOnChangeQuestionnaireSignOff(c echo.Context) error {

	projectId := c.Param("projectId")
	stage := "4"
	info,err := r.SOService.GetSignOffDetailOnChange(projectId,stage)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return  c.JSON(http.StatusOK,info)
}
		
func (r *QuestionnaireSORoute)SOcreate(c echo.Context) error {
	var pth pilotquestionnaire.Path
	if err := c.Bind(&pth); err != nil {
		return err
	}
	err :=r.SOService.SelectFileToSignOff(&pth)
	if err != nil {
		return err
	}

	projectID, err := strconv.Atoi(pth.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	var IsSignOff []pilotquestionnaire.FilePath
	for _, ps := range pth.FilePath{
		if ps.IsSign == true {
			IsSignOff = append(IsSignOff,ps)
		}
	}

	member,err := r.MemberService.GetMembersOfPj(pth.ProjectID)
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

	if len(mbs) > 0 && len(IsSignOff) > 0 {
		err := r.SOService.AlertQuestionnaireSignOff(IsSignOff,&mbs,PMName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_4"
		err = r.StageService.UpdateStage(pth.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusCreated, pth)
}

func (r *QuestionnaireSORoute)createQuestionnaireSignOff(c echo.Context) error {
	var pth pilotquestionnaire.Path
	if err := c.Bind(&pth); err != nil {
		return err
	}
	ap,err := r.PQService.GetAllPilotQuestionnairePath(pth.ProjectID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fileNumber := len(ap)
	err = r.SOService.InsertQuestionnaireSignOff(&pth,fileNumber)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	projectID, err := strconv.Atoi(pth.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	var IsSignOff []pilotquestionnaire.FilePath
	for _, ps := range pth.FilePath{
			IsSignOff = append(IsSignOff,ps)
	}

	member,err := r.MemberService.GetMembersOfPj(pth.ProjectID)
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

	if len(mbs) > 0 && len(IsSignOff) > 0 {
		err := r.SOService.AlertQuestionnaireSignOff(IsSignOff,&mbs,PMName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_4"
		err = r.StageService.UpdateStage(pth.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusCreated, pth)
}

func (r *QuestionnaireSORoute)editQuestionnaireSignOff(c echo.Context) error {
	var pth pilotquestionnaire.Path
	if err := c.Bind(&pth); err != nil {
		return err
	}
	ap,err := r.PQService.GetAllPilotQuestionnairePath(pth.ProjectID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fileNumber := len(ap)
	err = r.SOService.InsertQuestionnaireSignOff(&pth,fileNumber)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	projectID, err := strconv.Atoi(pth.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	var IsSignOff []pilotquestionnaire.FilePath
	for _, ps := range pth.FilePath{
			IsSignOff = append(IsSignOff,ps)
	}

	member,err := r.MemberService.GetMembersOfPj(pth.ProjectID)
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

	if len(mbs) > 0 && len(IsSignOff) > 0 {
		err := r.SOService.AlertEditQuestionnaireSignOff(IsSignOff,&mbs,PMName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusCreated, pth)
}

func (r *QuestionnaireSORoute)SOedit(c echo.Context) error {
	var pth pilotquestionnaire.Path
	if err := c.Bind(&pth); err != nil {
		return err
	}
	err :=r.SOService.SelectFileToSignOff(&pth)
	if err != nil {
		return err
	}

	projectID, err := strconv.Atoi(pth.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	var IsSignOff []pilotquestionnaire.FilePath
	for _, ps := range pth.FilePath{
		if ps.IsSign == true {
			IsSignOff = append(IsSignOff,ps)
		}
	}

	member,err := r.MemberService.GetMembersOfPj(pth.ProjectID)
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

	if len(mbs) > 0 && len(IsSignOff) > 0 {
		err := r.SOService.AlertEditQuestionnaireSignOff(IsSignOff,&mbs,PMName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusCreated, pth)
}

func (r *QuestionnaireSORoute)SOeditDetailOnChange(c echo.Context) error {
	var d pilotquestionnaire.PilotQuestionnaire
	if err := c.Bind(&d); err != nil {
		return err
	}
	stage := "4"
	err := r.PQService.DetailOnChange(&d,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return  c.JSON(http.StatusCreated, d)
}