package api

import(
	"github.com/labstack/echo/v4"
	"miami-back-end/project-review"
	"log"
	"net/http"
	"strconv"
	"time"
	"miami-back-end/members"
	"miami-back-end/project"
	"gopkg.in/guregu/null.v4"
	"miami-back-end/stage"
	
)
type ProjectReviewRoute struct{
	ProjectReviewService *projectreview.ProjectReviewService
	MemberService *members.MemberService
	ProjectService *project.Service
	StageService *stage.StageService

}

func NewProjectReviewRoute(projectReviewService *projectreview.ProjectReviewService,MemberService *members.MemberService,ProjectService *project.Service,StageService *stage.StageService) *ProjectReviewRoute{
	return &ProjectReviewRoute{
		ProjectReviewService: projectReviewService,
		MemberService: MemberService,
		ProjectService: ProjectService,
		StageService: StageService,
	}
}

func (r *ProjectReviewRoute)Group(g *echo.Group) {

	g.GET("/:projectId",r.getProjectReview)
	g.POST("/create",r.createProjectReview)
	g.PATCH("/edit",r.editProjectReview)
	g.PATCH("/cancel",r.cancelProjectReview)
	g.POST("/bypass",r.byPassProjectReview)

}

func (r *ProjectReviewRoute)getProjectReview(c echo.Context) error {
	projectId := c.Param("projectId")
	prw,err := r.ProjectReviewService.GetProjectReview(projectId)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	
	return c.JSON(http.StatusOK,prw)
}

func (r *ProjectReviewRoute)createProjectReview(c echo.Context) error {
	var prw projectreview.ProjectReview
	if err := c.Bind(&prw); err != nil {
		return err
	}
	projectID, err := strconv.Atoi(prw.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	MSLink := ""
	Other := ""
	if prw.IsOnline {
		if prw.MSLink.Valid {
			MSLink = prw.MSLink.String
		}
	}

	if prw.IsOther {
		if prw.Other.Valid {
			Other = prw.Other.String
		}
	}
	newprw := projectreview.ProjectReview{
		ProjectID: prw.ProjectID,
		MeetingDate: prw.MeetingDate,
		MeetingTime: prw.MeetingTime,
		IsOnline: prw.IsOnline,
		MSLink: null.StringFrom(MSLink),
		IsOnsite: prw.IsOnsite,
		IsRoom1: prw.IsRoom1,
		IsRoom2: prw.IsRoom2,
		IsRoom3: prw.IsRoom3,
		IsRoom4: prw.IsRoom4,
		IsOther: prw.IsOther,
		Other: null.StringFrom(Other),
		Note: prw.Note,
		CreatedBy: null.StringFrom(userLogin.Email),
		CreatedAt: time.Now(),
	}
	err = r.ProjectReviewService.CreateProjectReview(&newprw)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}
	member,err := r.MemberService.GetMembersOfPj(prw.ProjectID)
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
		err := r.ProjectReviewService.AlertCreateProjectReview(&mbs,&newprw,PMName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_18"
		err = r.StageService.UpdateStage(prw.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated,prw)
}

func (r *ProjectReviewRoute)editProjectReview(c echo.Context) error {
	var prw projectreview.ProjectReview
	if err := c.Bind(&prw); err != nil {
		return err
	}
	projectID, err := strconv.Atoi(prw.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	MSLink := ""
	Other := ""
	if prw.IsOnline {
		if prw.MSLink.Valid {
			MSLink = prw.MSLink.String
		}
	}

	if prw.IsOther {
		if prw.Other.Valid {
			Other = prw.Other.String
		}
	}
	newprw := projectreview.ProjectReview{
		ProjectID: prw.ProjectID,
		MeetingDate: prw.MeetingDate,
		MeetingTime: prw.MeetingTime,
		IsOnline: prw.IsOnline,
		MSLink: null.StringFrom(MSLink),
		IsOnsite: prw.IsOnsite,
		IsRoom1: prw.IsRoom1,
		IsRoom2: prw.IsRoom2,
		IsRoom3: prw.IsRoom3,
		IsRoom4: prw.IsRoom4,
		IsOther: prw.IsOther,
		Other: null.StringFrom(Other),
		Note: prw.Note,
		UpdatedBy: null.StringFrom(userLogin.Email),
		UpdatedAt: null.TimeFrom(time.Now()),
	}
	err = r.ProjectReviewService.EditProjectReview(&newprw)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	member,err := r.MemberService.GetMembersOfPj(prw.ProjectID)
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
		err := r.ProjectReviewService.AlertEditProjectReview(&mbs,&newprw,PMName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated,prw)
}

func (r *ProjectReviewRoute)cancelProjectReview(c echo.Context) error {
	var prw projectreview.ProjectReview
	if err := c.Bind(&prw); err != nil {
		return err
	}
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	newPrw := projectreview.ProjectReview{
		ProjectID: prw.ProjectID,
		Note: prw.Note,
		ISCancel: true,
		UpdatedBy: null.StringFrom(userLogin.Email),
		UpdatedAt: null.TimeFrom(time.Now()),
	}
	err := r.ProjectReviewService.CancelProjectReview(&newPrw)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	
	projectID, err := strconv.Atoi(prw.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	member,err := r.MemberService.GetMembersOfPj(prw.ProjectID)
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
		err := r.ProjectReviewService.AlertCancelProjectReview(&mbs,&newPrw,PMName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated,prw)
}

func (r *ProjectReviewRoute)byPassProjectReview(c echo.Context) error {
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	var prw projectreview.ProjectReview
	if err := c.Bind(&prw); err != nil {
		return err
	}
	newPrw := projectreview.ProjectReview{
		ProjectID: prw.ProjectID,
		Note: prw.Note,
		CreatedBy: null.StringFrom(userLogin.Email),
		ISByPass: true,
		CreatedAt: time.Now(),
	}
	err := r.ProjectReviewService.ByPassProjectReview(&newPrw)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}
	stage := "stage_18"
	err = r.StageService.UpdateStage(prw.ProjectID,stage)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated,prw)
}