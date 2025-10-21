package api

import (
	"log"
	"miami-back-end/project"
	"miami-back-end/members"
	"net/http"
	"strconv"
	"time"
	"github.com/labstack/echo/v4"
)

type ProjectRoute struct {
	projectService *project.Service
	memberService *members.MemberService
}
func NewProjectRoute(service *project.Service,memberService *members.MemberService) *ProjectRoute {
	return &ProjectRoute{
		projectService: service,
		memberService: memberService,
	}
}
func (r *ProjectRoute) Group(g *echo.Group)  {
	g.Use(Auth())
	g.POST("/create-member-and-stage", r.createMemberAndStageRoute)
	g.GET("/getproject", r.getProject)
}

func (r *ProjectRoute) createMemberAndStageRoute(c echo.Context) error  {
	var body project.Project
	if err := c.Bind(&body); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	var POinfo project.UserPosition
	for i, x := range body.Owners {
		if i == 0 {
			POinfo.UID = x.UID
		}
	}
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	newProject := project.Project{
		ID: body.ID,
		Name: body.Name,
		Year: body.Year,
		Status: "active",
		CreatedAt: time.Now(),
		CreatedBy: userLogin.Email,
	}
	if err := newProject.ValidateCreate(); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
		
	}
	owner := project.Member{
		UID: POinfo.UID,
		Position: "PO",
		CreatedAt: time.Now(),
		CreatedBy: userLogin.Email,
	}
	checkNullPJ, err := r.memberService.GetMembersOfPj(strconv.FormatInt(newProject.ID, 10))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	if checkNullPJ == nil {
		err = r.projectService.CreateMemberAndStage(&newProject, &owner)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusCreated, body)

}
func (r *ProjectRoute) getProject(c echo.Context) error {
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	status := c.QueryParam("status")
	if status == "" {
		status = "active"
	}
	xs,err := r.projectService.GetProjectsByUIDAndStatus(userLogin.UID,status,userLogin.IsAdmin);
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, xs);
}
