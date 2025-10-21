package api

import(
	"github.com/labstack/echo/v4"	
	"net/http"
	"strconv"
	"time"
	"log"
	"miami-back-end/de-training"
	"miami-back-end/members"
	"miami-back-end/project"
	"miami-back-end/stage"
	"miami-back-end/pilot-questionnaire"
	"gopkg.in/guregu/null.v4"
)

type DeTrainingRoute struct{
	DeTrainingService *detraining.DeTrainingService
	MemberService *members.MemberService
	ProjectService *project.Service
	StageService *stage.StageService
	PQService *pilotquestionnaire.PilotQuestionnaireService

}

func NewDeTrainingRoute( DeTrainingService *detraining.DeTrainingService,MemberService *members.MemberService,ProjectService *project.Service,StageService *stage.StageService,PQService *pilotquestionnaire.PilotQuestionnaireService) *DeTrainingRoute{
	return &DeTrainingRoute{
		DeTrainingService: DeTrainingService,
		MemberService: MemberService,
		ProjectService: ProjectService,
		StageService: StageService,
		PQService: PQService,
	}
}

func (r *DeTrainingRoute)Group(g *echo.Group){
	g.GET("/:projectId",r.getDeTrainingInfo)
	g.POST("/create",r.DeTrainingCreate)
	g.PATCH("/edit-detail",r.DeEditDetail)
	g.PATCH("/edit-questionnaire",r.DeEditQuestionnaire)
	g.PATCH("/cancel",r.DeCancel)
	g.POST("/bypass",r.DeBypass)
	g.POST("/path",r.DePath)
	g.PATCH("/update-path",r.DeUpdatePath)
}

func (r *DeTrainingRoute)getDeTrainingInfo(c echo.Context) error {
	projectId := c.Param("projectId")
	deTraining,err := r.DeTrainingService.GetDeTrainingInfo(projectId)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, deTraining)
}

func (r *DeTrainingRoute)DeTrainingCreate(c echo.Context)error{
	var deinfo detraining.DeTraining
	if err := c.Bind(&deinfo); err != nil {
		return err
	}
	projectID, err := strconv.Atoi(deinfo.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	MSLink := ""
	Other := ""
	if deinfo.IsOnline == true {
		MSLink = deinfo.MSLink
	}
	if deinfo.IsOther == true {
		Other = deinfo.Other
	}
	newTraining := detraining.DeTraining{
		ProjectID: deinfo.ProjectID,
		MeetingDate: deinfo.MeetingDate,
		MeetingTime: deinfo.MeetingTime,
		IsOnline: deinfo.IsOnline,
		MSLink: MSLink,
		IsOnsite: deinfo.IsOnsite,
		IsRoom1: deinfo.IsRoom1,
		IsRoom2: deinfo.IsRoom2,
		IsRoom3: deinfo.IsRoom3,
		IsRoom4: deinfo.IsRoom4,
		IsOther: deinfo.IsOther,
		Other: Other,
		Note: deinfo.Note,
		CreatedBy: userLogin.Email,
		CreatedAt: time.Now(),
	}
	err = r.DeTrainingService.CreateDeTraining(&newTraining)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}
	member,err := r.MemberService.GetMembersOfPj(deinfo.ProjectID)
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

	DeFilePath,err := r.DeTrainingService.GetLatestDeQuestionnairePath(deinfo.ProjectID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if len(mbs) > 0 {
		err := r.DeTrainingService.AlertDeTraining(&mbs,&newTraining,PMName,DeName,pj.Name,DeFilePath)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_11"
		err = r.StageService.UpdateStage(deinfo.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated,deinfo)
}

func (r *DeTrainingRoute)DeEditDetail(c echo.Context) error {
	var p detraining.DeTraining
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
	newDe := detraining.DeTraining{
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
	err = r.DeTrainingService.DeEditdetail(&newDe)
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
		err := r.DeTrainingService.AlertEditDetailDeTraining(&mbs,&newDe,PMName,DeName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated,p)
}

func (r *DeTrainingRoute)DeEditQuestionnaire(c echo.Context) error {
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
	stage := "11"
	err = r.PQService.DetailOnChange(&p,stage)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	pthlt,err := r.DeTrainingService.GetLatestDeQuestionnairePath(p.ProjectID)
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
		err := r.DeTrainingService.AlertEditDeTrainingQuestion(&mbs,&p,PMName,pj.Name,pthlt)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		
	}
	return c.JSON(http.StatusCreated,p)
}

func (r *DeTrainingRoute)DeCancel(c echo.Context) error {
	var p detraining.DeTraining
	if err := c.Bind(&p); err != nil {
		return err
	}
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	newDe := detraining.DeTraining{
		ProjectID: p.ProjectID,
		Note: p.Note,
		ISCancel: true,
		UpdatedBy: null.StringFrom(userLogin.Email),
		UpdatedAt: null.TimeFrom(time.Now()),
	}
	err := r.DeTrainingService.CancelDeTraining(&newDe)
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
		err := r.DeTrainingService.AlertCancelDeTraining(&mbs,&newDe,PMName,DeName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated,p)
}

func (r *DeTrainingRoute)DeBypass(c echo.Context) error {
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	var q detraining.DeTraining
	if err := c.Bind(&q); err != nil {
		return err
	}
	newDe := detraining.DeTraining{
		ProjectID: q.ProjectID,
		Note: q.Note,
		CreatedBy: userLogin.Email,
		ISByPass: true,
		CreatedAt: time.Now(),
	}
	err := r.DeTrainingService.ByPassDeTraining(&newDe)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}
	stage := "stage_11"
	err = r.StageService.UpdateStage(q.ProjectID,stage)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated,q)
}

func (r *DeTrainingRoute)DePath(c echo.Context)error{
	var p pilotquestionnaire.Path
	if err := c.Bind(&p); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	err := r.DeTrainingService.InsertDeTrainingPath(&p)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated,p)
}

func (r *DeTrainingRoute)DeUpdatePath(c echo.Context)error{
	var p pilotquestionnaire.Path
	if err := c.Bind(&p); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	ap,err := r.DeTrainingService.GetAllDeTrainingPath(p.ProjectID)
	
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fileNumber := len(ap)
	err = r.DeTrainingService.UpdatePathDeTraining(&p,fileNumber)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated,p)
}