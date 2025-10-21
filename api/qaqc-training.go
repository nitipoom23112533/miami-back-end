package api

import(
	"miami-back-end/qaqc-training"
	"miami-back-end/members"
	"miami-back-end/project"
	"miami-back-end/stage"
	"miami-back-end/pilot-questionnaire"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"time"
	"gopkg.in/guregu/null.v4"
)

type QAQCTrainingRoute struct{
	QAQCService *qaqctraining.QAQCService
	MemberService *members.MemberService
	ProjectService *project.Service
	StageService *stage.StageService
	PQService *pilotquestionnaire.PilotQuestionnaireService
}

func NewQAQCTrainingRoute(QAQCService *qaqctraining.QAQCService,MemberService *members.MemberService,ProjectService *project.Service,StageService *stage.StageService,
	PQService *pilotquestionnaire.PilotQuestionnaireService ) *QAQCTrainingRoute{
	return &QAQCTrainingRoute{
		QAQCService: QAQCService,
		MemberService: MemberService,
		ProjectService: ProjectService,
		StageService: StageService,
		PQService: PQService,
	}
}

func (r *QAQCTrainingRoute)Group(g *echo.Group){
	g.Use(Auth())
	g.GET("/:projectId",r.getQAQCTrainingInfo)
	g.POST("/create", r.QAQCCreate)
	g.PATCH("/edit-detail",r.QAQCEditDetail)
	g.PATCH("/edit-questionnaire",r.QAQCEditQuestionnaire)
	g.PATCH("/cancel",r.QAQCCancel)
	g.POST("/bypass",r.QAQCBypass)
	g.POST("/path", r.QAQCPath)
	g.POST("/update-path", r.QAQCUpdatePath)
}

func (r *QAQCTrainingRoute)getQAQCTrainingInfo(c echo.Context) error {
	projectId := c.Param("projectId")
	info,err := r.QAQCService.GetQAQCTrainingInfo(projectId)
	if err != nil {
		log.Println("err",err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK,info)
}

func (r *QAQCTrainingRoute)QAQCCreate(c echo.Context)error{
	var qaqcinfo qaqctraining.QAQCTraining
	if err := c.Bind(&qaqcinfo); err != nil {
		return err
	}
	
	projectID, err := strconv.Atoi(qaqcinfo.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	MSLink := ""
	Other := ""
	if qaqcinfo.IsOnline == true {
		MSLink = qaqcinfo.MSLink
	}
	if qaqcinfo.IsOther == true {
		Other = qaqcinfo.Other
	}
	newTraining := qaqctraining.QAQCTraining{
		ProjectID: qaqcinfo.ProjectID,
		MeetingDate: qaqcinfo.MeetingDate,
		MeetingTime: qaqcinfo.MeetingTime,
		IsOnline: qaqcinfo.IsOnline,
		MSLink: MSLink,
		IsOnsite: qaqcinfo.IsOnsite,
		IsRoom1: qaqcinfo.IsRoom1,
		IsRoom2: qaqcinfo.IsRoom2,
		IsRoom3: qaqcinfo.IsRoom3,
		IsRoom4: qaqcinfo.IsRoom4,
		IsOther: qaqcinfo.IsOther,
		Other: Other,
		Note: qaqcinfo.Note,
		CreatedBy: userLogin.Email,
		CreatedAt: time.Now(),
	}
	err = r.QAQCService.SendCreateQAQCTraining(&newTraining)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}
	pthlt,err := r.QAQCService.GetLatestQAQCTrainingPath(qaqcinfo.ProjectID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	
	member,err := r.MemberService.GetMembersOfPj(qaqcinfo.ProjectID)
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
		err := r.QAQCService.AlertQAQCTraining(&mbs,&newTraining,PMName,QaqcName,pj.Name,pthlt)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_8"
		err = r.StageService.UpdateStage(qaqcinfo.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated,qaqcinfo)
}

func (r *QAQCTrainingRoute)QAQCEditDetail(c echo.Context) error {
	var p qaqctraining.QAQCTraining
	if err := c.Bind(&p); err != nil {
		return err
	}
	projectID, err := strconv.Atoi(p.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	MSLink := ""
	Other := ""
	if p.IsOnline == true {
		MSLink = p.MSLink
	}
	if p.IsOther == true {
		Other = p.Other
	}
	newpilotquestionnaire := qaqctraining.QAQCTraining{
		ProjectID: p.ProjectID,
		MeetingDate: p.MeetingDate,
		MeetingTime: p.MeetingTime,
		IsOnline: p.IsOnline,
		MSLink: MSLink,
		IsOnsite: p.IsOnsite,
		IsRoom1: p.IsRoom1,
		IsRoom2: p.IsRoom2,
		IsRoom3: p.IsRoom3,
		IsRoom4: p.IsRoom4,
		IsOther: p.IsOther,
		Other: Other,
		Note: p.Note,
		UpdatedBy: null.StringFrom(userLogin.Email),
		UpdatedAt: null.TimeFrom(time.Now()),
	}
	err = r.QAQCService.QAQCEditdetail(&newpilotquestionnaire)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	member,err := r.MemberService.GetMembersOfPj(p.ProjectID)
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
		err := r.QAQCService.AlertEditDetailQAQCTraining(&mbs,&newpilotquestionnaire,PMName,QaqcName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated,p)
}

func (r *QAQCTrainingRoute)QAQCUpdatePath(c echo.Context)error{
	var p pilotquestionnaire.Path
	if err := c.Bind(&p); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	ap,err := r.QAQCService.GetAllQAQCTrainingPath(p.ProjectID)
	
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fileNumber := len(ap)
	err = r.QAQCService.UpdatePathQAQCTraining(&p,fileNumber)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated,p)
}

func (r *QAQCTrainingRoute)QAQCEditQuestionnaire(c echo.Context) error {
	var p pilotquestionnaire.PilotQuestionnaire
	if err := c.Bind(&p); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	projectID, err := strconv.Atoi(p.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	stage := "8"
	err = r.PQService.DetailOnChange(&p,stage)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	pthlt,err := r.QAQCService.GetLatestQAQCTrainingPath(p.ProjectID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	member,err := r.MemberService.GetMembersOfPj(p.ProjectID)
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
		err := r.QAQCService.AlertEditQAQCTrainingQuestion(&mbs,&p,PMName,pj.Name,pthlt)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		
	}
	return c.JSON(http.StatusCreated,p)
}

func (r *QAQCTrainingRoute)QAQCCancel(c echo.Context) error {
	var p qaqctraining.QAQCTraining
	if err := c.Bind(&p); err != nil {
		return err
	}
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	newPilotQuestionnaire := qaqctraining.QAQCTraining{
		ProjectID: p.ProjectID,
		Note: p.Note,
		ISCancel: true,
		UpdatedBy: null.StringFrom(userLogin.Email),
		UpdatedAt: null.TimeFrom(time.Now()),
	}
	err := r.QAQCService.CancelQAQCTraining(&newPilotQuestionnaire)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	
	projectID, err := strconv.Atoi(p.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	member,err := r.MemberService.GetMembersOfPj(p.ProjectID)
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
		err := r.QAQCService.AlertCancelQAQCTraining(&mbs,&newPilotQuestionnaire,PMName,QaqcName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated,p)
}

func (r *QAQCTrainingRoute)QAQCBypass(c echo.Context) error {
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	var q qaqctraining.QAQCTraining
	if err := c.Bind(&q); err != nil {
		return err
	}
	newPilotQuestionnaire := qaqctraining.QAQCTraining{
		ProjectID: q.ProjectID,
		Note: q.Note,
		CreatedBy: userLogin.Email,
		ISByPass: true,
		CreatedAt: time.Now(),
	}
	err := r.QAQCService.ByPassQAQCTraining(&newPilotQuestionnaire)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}
	stage := "stage_8"
	err = r.StageService.UpdateStage(q.ProjectID,stage)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated,q)
}

func (r *QAQCTrainingRoute)QAQCPath(c echo.Context)error{
	var p pilotquestionnaire.Path
	if err := c.Bind(&p); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	err := r.QAQCService.InsertQAQCTrainingPath(&p)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated,p)
}