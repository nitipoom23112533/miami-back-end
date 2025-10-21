package api
import(
	"miami-back-end/operation-training"
	"miami-back-end/pilot-questionnaire"
	// "miami-back-end/questionnaire-sign-off"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"time"
	"miami-back-end/members"
	"miami-back-end/stage"
	"miami-back-end/project"
	"gopkg.in/guregu/null.v4"
)

type TrainingRoute struct{
	TrainingService *operationtraining.OperationTrainingServive
	StageService *stage.StageService
	MemberService *members.MemberService
	ProjectService *project.Service
	PQService *pilotquestionnaire.PilotQuestionnaireService
	// SOService *questionnairesignoff.QuestionnaireSignOffService
}

func NewTrainingRoute(trainingService *operationtraining.OperationTrainingServive,stageService *stage.StageService,
	memberService *members.MemberService,projectService *project.Service,pqService *pilotquestionnaire.PilotQuestionnaireService) *TrainingRoute {
	return &TrainingRoute{
		TrainingService: trainingService,
		StageService: stageService,
		MemberService: memberService,
		ProjectService: projectService,
		PQService: pqService,
	}
}

func (r *TrainingRoute)Group(g *echo.Group){
	g.Use(Auth())
	g.GET("/:projectId", r.getTrainingInfo)
	g.GET("/is-sign/:projectId", r.getAllPilotQuestionnairePathIsSign)
	g.POST("/invite",r.sendInviteTraining)
	g.PATCH("/edit-detail",r.editDetailTraining)
	g.POST("/edit-question",r.editQuestionTraining)
	g.PATCH("/edit-detail-on-change",r.editDetailTrainingOnChange)
	g.POST("/bypass",r.bypassTraining)
	g.PATCH("/cancel",r.cancelTraining)
}
func (r *TrainingRoute)sendInviteTraining(c echo.Context) error {
	var t operationtraining.Training
	if err := c.Bind(&t); err != nil {
		return err
	}
	projectID, err := strconv.Atoi(t.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	MSLink := ""
	Other := ""
	if t.IsOnline == true {
		MSLink = t.MSLink
	}
	if t.IsOther == true {
		Other = t.Other
	}
	newTraining := operationtraining.Training{
		ProjectID: t.ProjectID,
		MeetingDate: t.MeetingDate,
		MeetingTime: t.MeetingTime,
		IsOnline: t.IsOnline,
		MSLink: MSLink,
		IsOnsite: t.IsOnsite,
		IsRoom1: t.IsRoom1,
		IsRoom2: t.IsRoom2,
		IsRoom3: t.IsRoom3,
		IsRoom4: t.IsRoom4,
		IsOther: t.IsOther,
		Other: Other,
		Note: t.Note,
		CreatedBy: userLogin.Email,
		CreatedAt: time.Now(),
	}
	err = r.TrainingService.SendTrainingInvite(&newTraining)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	NewIsTraining := operationtraining.Training{
		ProjectID: t.ProjectID,
		FilePath: t.FilePath,
	}		

	err = r.TrainingService.SelectFileToTraining(&NewIsTraining)
	if err != nil {
		return err
	}

	member,err := r.MemberService.GetMembersOfPj(t.ProjectID)
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

	pthlt,err := r.TrainingService.GetLatestTrainingPath(t.ProjectID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(mbs) > 0 {
		err := r.TrainingService.AlertInviteTraining(&mbs,&newTraining,PMName,pj.Name,pthlt)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_5"
		err = r.StageService.UpdateStage(t.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

	}

	return c.JSON(http.StatusCreated,t)
}
func (r *TrainingRoute)getTrainingInfo(c echo.Context) error {
	projectId := c.Param("projectId")
	info,err := r.TrainingService.GetTrainingInfoISSign(projectId)
	if err != nil {
		log.Println("err",err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK,info)
}

func (r *TrainingRoute)getAllPilotQuestionnairePathIsSign(c echo.Context) error {
	projectId := c.Param("projectId")
	info,err := r.TrainingService.GetAllPilotQuestionnairePathIsSign(projectId)
	if err != nil {
		log.Println("err",err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK,info)
}

func (r *TrainingRoute)editDetailTraining(c echo.Context) error {

	var t operationtraining.Training
	if err := c.Bind(&t); err != nil {
		return err
	}
	projectID, err := strconv.Atoi(t.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	MSLink := ""
	Other := ""
	if t.IsOnline == true {
		MSLink = t.MSLink
	}
	if t.IsOther == true {
		Other = t.Other
	}
	newTraining := operationtraining.Training{
		ProjectID: t.ProjectID,
		MeetingDate: t.MeetingDate,
		MeetingTime: t.MeetingTime,
		IsOnline: t.IsOnline,
		MSLink: MSLink,
		IsOnsite: t.IsOnsite,
		IsRoom1: t.IsRoom1,
		IsRoom2: t.IsRoom2,
		IsRoom3: t.IsRoom3,
		IsRoom4: t.IsRoom4,
		IsOther: t.IsOther,
		Other: Other,
		Note: t.Note,
		UpdatedBy: null.StringFrom(userLogin.Email),
		UpdatedAt: null.TimeFrom(time.Now()),
	}
	err = r.TrainingService.EditDetailTraining(&newTraining)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	member,err := r.MemberService.GetMembersOfPj(t.ProjectID)
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

	if len(mbs) > 0 {
		err := r.TrainingService.AlertEditDetailTraining(&mbs,&newTraining,PMName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated,t)
}

func (r *TrainingRoute)editQuestionTraining(c echo.Context) error {
	var pth pilotquestionnaire.Path
	if err := c.Bind(&pth); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
		// return err
	}

	ap,err := r.PQService.GetAllPilotQuestionnairePath(pth.ProjectID)
	
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fileNumber := len(ap) + 1

	var path []pilotquestionnaire.FilePath
		
	for _, ps := range pth.FilePath{

		path = append(path, pilotquestionnaire.FilePath{
			Path: ps.Path,
			FileName: ps.FileName,
			Number: fileNumber,
			IsNew: true,
			IsSign: true,
			ISTraining: true,
		})
		fileNumber++
	}
	pth.FilePath = path

	err = r.TrainingService.UpdateTrainingPath(&pth)
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

	var ISTraining []pilotquestionnaire.FilePath
	for _, ps := range pth.FilePath{
		if ps.ISTraining == true {
			ISTraining = append(ISTraining,ps)
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

	if len(mbs) > 0 && len(ISTraining) > 0 {
		err := r.TrainingService.AlertEditQuestionTraining(ISTraining,&mbs,PMName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusCreated,pth)
	// return  nil
}

func (r *TrainingRoute)editDetailTrainingOnChange(c echo.Context) error {
	var dtl operationtraining.Training
	if err := c.Bind(&dtl); err != nil {
		return err
	}
	newDetailOnchange := pilotquestionnaire.PilotQuestionnaire{
		ProjectID: dtl.ProjectID,
		DetailOnChange: dtl.DetailOnChange,
	}
	stage := "5"
	err := r.PQService.DetailOnChange(&newDetailOnchange,stage)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated,dtl)
}

func (r *TrainingRoute)bypassTraining(c echo.Context) error {
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	var t operationtraining.Training
	if err := c.Bind(&t); err != nil {
		return err
	}
	newPilotQuestionnaire := operationtraining.Training{
		ProjectID: t.ProjectID,
		Note: t.Note,
		CreatedBy: userLogin.Email,
		ISByPass: true,
		CreatedAt: time.Now(),
	}
	err := r.TrainingService.ByPassTrraining(&newPilotQuestionnaire)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}
	stage := "stage_5"
	err = r.StageService.UpdateStage(t.ProjectID,stage)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated,t)
}

func (r *TrainingRoute)cancelTraining(c echo.Context) error {
	var t operationtraining.Training
	if err := c.Bind(&t); err != nil {
		return err
	}
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	newTraining := operationtraining.Training{
		ProjectID: t.ProjectID,
		Note: t.Note,
		ISCancel: true,
		UpdatedBy: null.StringFrom(userLogin.Email),
		UpdatedAt: null.TimeFrom(time.Now()),
	}
	err := r.TrainingService.CancelTraining(&newTraining)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	
	projectID, err := strconv.Atoi(t.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	member,err := r.MemberService.GetMembersOfPj(t.ProjectID)
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

	if len(mbs) > 0 {
		err := r.TrainingService.AlertCancelTraining(&mbs,&newTraining,PMName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated,t)
}