package api

import(
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"miami-back-end/report-submission-project-close"
	"time"
	"strconv"
	"miami-back-end/project"
	"miami-back-end/members"
	"miami-back-end/stage"
	"gopkg.in/guregu/null.v4"

)

type RsPcRoute struct{

	RsPcService *reportsubmissionprojectclose.ReportSubmissionService
	ProjectService *project.Service
	MemberService *members.MemberService
	StageService *stage.StageService


}

func NewRsPcRoute(RsPcService *reportsubmissionprojectclose.ReportSubmissionService,ProjectService *project.Service,MemberService *members.MemberService,StageService *stage.StageService) *RsPcRoute {
	return &RsPcRoute{
		RsPcService: RsPcService,
		ProjectService: ProjectService,
		MemberService: MemberService,
		StageService: StageService,
	}
}

func (r *RsPcRoute)Group(g *echo.Group){
	g.Use(Auth())
	g.GET("/:projectId",r.getReportSubmission)
	g.POST("/create",r.insertReportSubmission)
	g.PATCH("/update",r.updateReportSubmission)
}

func (r *RsPcRoute)getReportSubmission(c echo.Context) error {
	projectId := c.Param("projectId")
	reportSubmission,err := r.RsPcService.GetReportSubmission(projectId)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, reportSubmission)
}

func (r *RsPcRoute)insertReportSubmission(c echo.Context) error {
	var reportSubmission reportsubmissionprojectclose.ReportSubmission
	if err := c.Bind(&reportSubmission); err != nil {
		return err
	}
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	newRsPc := reportsubmissionprojectclose.ReportSubmission{
		ProjectID: reportSubmission.ProjectID,
		SubmissionDate:reportSubmission.SubmissionDate,
		CreatedBy: null.StringFrom(userLogin.Email),
		CreatedAt: null.TimeFrom(time.Now()),
	}

	err := r.RsPcService.InsertReportSubmission(&newRsPc)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	projectID, err := strconv.Atoi(reportSubmission.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	member,err := r.MemberService.GetMembersOfPj(reportSubmission.ProjectID)
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
		err := r.RsPcService.AlertSubmissionClose(&mbs,&newRsPc,PMName,pj.Name,"report")
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_16"
		err = r.StageService.UpdateStage(reportSubmission.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, reportSubmission)
}

func (r *RsPcRoute)updateReportSubmission(c echo.Context) error {
	var reportSubmission reportsubmissionprojectclose.ReportSubmission
	if err := c.Bind(&reportSubmission); err != nil {
		return err
	}
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	newRsPc := reportsubmissionprojectclose.ReportSubmission{
		ProjectID: reportSubmission.ProjectID,
		CloseDate: reportSubmission.CloseDate,
		UpdatedBy: null.StringFrom(userLogin.Email),
		UpdatedAt: null.TimeFrom(time.Now()),
	}
	err := r.RsPcService.UpdateReportSubmission(&newRsPc)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	projectID, err := strconv.Atoi(reportSubmission.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	member,err := r.MemberService.GetMembersOfPj(reportSubmission.ProjectID)
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
		err := r.RsPcService.AlertSubmissionClose(&mbs,&newRsPc,PMName,pj.Name,"close")
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_17"
		err = r.StageService.UpdateStage(reportSubmission.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, reportSubmission)
}