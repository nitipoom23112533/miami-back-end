package api
import (
	"miami-back-end/pilot-questionnaire"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"time"
	"miami-back-end/members"
	"miami-back-end/project"
	"miami-back-end/stage"
	"gopkg.in/guregu/null.v4"
	
)


type pilotQuesttionnairRoute struct{
	PQService *pilotquestionnaire.PilotQuestionnaireService
	ProjectService *project.Service
	MemberService *members.MemberService
	StageService *stage.StageService

}

func NewPilotQuesttionaireRoute(PQRoutService *pilotquestionnaire.PilotQuestionnaireService,ProjectService *project.Service,MemberService *members.MemberService,StageService *stage.StageService ) *pilotQuesttionnairRoute{
	return &pilotQuesttionnairRoute{
		PQService: PQRoutService,
		ProjectService: ProjectService,
		MemberService: MemberService,
		StageService: StageService,


	}
}

func (r *pilotQuesttionnairRoute)Group(g *echo.Group) {
	g.Use(Auth())
	g.GET("/:projectId",r.getPQInfo)
	g.POST("/inform",r.PQInform)
	g.POST("/bypass",r.PQBypass)
	g.PATCH("/edit-detail",r.PQEditDetail)
	g.PATCH("/edit-questionnaire",r.PQEditQuestionnaire)
	g.PATCH("/cancel",r.PQCancel)
	g.POST("/path",r.PQPath)
	g.POST("/update-path",r.PQUpdatePath)
}

func (r *pilotQuesttionnairRoute)PQInform(c echo.Context) error {

	var q pilotquestionnaire.PilotQuestionnaire
	if err := c.Bind(&q); err != nil {
		return err
	}

	projectID, err := strconv.Atoi(q.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	userLogin := ParseJWTCustomClaims(c.Get("user"))
	MSLink := ""
	Other := ""
	if q.IsOnline == true {
		MSLink = q.MSLink
	}
	if q.IsOther == true {
		Other = q.Other
	}
	newPilotQuestionnaire := pilotquestionnaire.PilotQuestionnaire{
		ProjectID: q.ProjectID,
		MeetingDate: q.MeetingDate,
		MeetingTime: q.MeetingTime,
		IsOnline: q.IsOnline,
		MSLink: MSLink,
		IsOnsite: q.IsOnsite,
		IsRoom1: q.IsRoom1,
		IsRoom2: q.IsRoom2,
		IsRoom3: q.IsRoom3,
		IsRoom4: q.IsRoom4,
		IsOther: q.IsOther,
		Other: Other,
		Note: q.Note,
		CreatedBy: userLogin.Email,
		CreatedAt: time.Now(),
	}
	err = r.PQService.CreatePilotQuestionnaire(&newPilotQuestionnaire)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	member,err := r.MemberService.GetMembersOfPj(q.ProjectID)
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

	pthlt,err := r.PQService.GetLatestPilotQuestionnairePath(q.ProjectID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if len(mbs) > 0 {
		err := r.PQService.AlertInformPilotQuestionnaire(&mbs,&newPilotQuestionnaire,PMName,pj.Name,pthlt)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_3"
		err = r.StageService.UpdateStage(q.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated,q)
}

func (r *pilotQuesttionnairRoute)PQBypass(c echo.Context) error {

	userLogin := ParseJWTCustomClaims(c.Get("user"))
	var p pilotquestionnaire.PilotQuestionnaire
	if err := c.Bind(&p); err != nil {
		return err
	}
	newPilotQuestionnaire := pilotquestionnaire.PilotQuestionnaire{
		ProjectID: p.ProjectID,
		Note: p.Note,
		CreatedBy: userLogin.Email,
		ISByPass: true,
		CreatedAt: time.Now(),
	}
	err := r.PQService.ByPassPilotQuestionnaire(&newPilotQuestionnaire)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}
	stage := "stage_3"
	err = r.StageService.UpdateStage(p.ProjectID,stage)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated,p)

}

func (r *pilotQuesttionnairRoute)getPQInfo(c echo.Context) error {
	projectId := c.Param("projectId")
	info,err := r.PQService.GetPilotQuestionnaireInfo(projectId)
	if err != nil {
		log.Println("err",err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK,info)
}
func (r *pilotQuesttionnairRoute)PQEditDetail(c echo.Context) error {
	var p pilotquestionnaire.PilotQuestionnaire
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
	newpilotquestionnaire := pilotquestionnaire.PilotQuestionnaire{
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
	err = r.PQService.Editdetail(&newpilotquestionnaire)
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
		err := r.PQService.AlertEditDetail(&mbs,&newpilotquestionnaire,PMName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated,p)
}

func (r *pilotQuesttionnairRoute)PQCancel(c echo.Context) error {

	var p pilotquestionnaire.PilotQuestionnaire
	if err := c.Bind(&p); err != nil {
		return err
	}
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	newPilotQuestionnaire := pilotquestionnaire.PilotQuestionnaire{
		ProjectID: p.ProjectID,
		Note: p.Note,
		ISCancel: true,
		UpdatedBy: null.StringFrom(userLogin.Email),
		UpdatedAt: null.TimeFrom(time.Now()),
	}
	err := r.PQService.CancelPilotQuestionnaire(&newPilotQuestionnaire)
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
		err := r.PQService.AlertCancelPilotQuestionnaire(&mbs,&newPilotQuestionnaire,PMName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated,p)
}

func (r *pilotQuesttionnairRoute)PQEditQuestionnaire(c echo.Context) error {
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
	stage := "3"
	err = r.PQService.DetailOnChange(&p,stage)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	pthlt,err := r.PQService.GetLatestPilotQuestionnairePath(p.ProjectID)
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
		err := r.PQService.AlertEditQuestion(&mbs,&p,PMName,pj.Name,pthlt)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		
	}
	return c.JSON(http.StatusCreated,p)
}

func (r *pilotQuesttionnairRoute)PQPath(c echo.Context) error{
	var p pilotquestionnaire.Path
	if err := c.Bind(&p); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	err := r.PQService.InsertPathPilotQuestionnaire(&p)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated,p)
}

func (r *pilotQuesttionnairRoute)PQUpdatePath(c echo.Context) error{
	var p pilotquestionnaire.Path
	if err := c.Bind(&p); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	ap,err := r.PQService.GetAllPilotQuestionnairePath(p.ProjectID)
	
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fileNumber := len(ap)
	err = r.PQService.UpdatePathPilotQuestionnaire(&p,fileNumber)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated,p)
}