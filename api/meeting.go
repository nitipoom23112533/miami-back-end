package api

import (
	"log"
	"miami-back-end/meeting"
	"miami-back-end/members"
	"miami-back-end/project"
	"miami-back-end/stage"
	"net/http"
	"strconv"
	"time"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
)

type MeetingRoute struct {

	MeetingService *meeting.MeetingService
	MemberService *members.MemberService
	ProjectService *project.Service
	StageService *stage.StageService
}

func NewMeetingRoute(meetingService *meeting.MeetingService,memberService *members.MemberService,projectService *project.Service,stageService *stage.StageService) *MeetingRoute {
	return &MeetingRoute{
		MeetingService: meetingService,
		MemberService: memberService,
		ProjectService: projectService,
		StageService: stageService,
	}
}

func (r *MeetingRoute)Group(g *echo.Group){
	g.Use(Auth())
	g.POST("/invite",r.inviteMeeting)
	g.GET("/:projectId",r.getMeetingInfo)
	g.PATCH("/edit",r.editInviteMetting)
	g.POST("/bypass",r.bypassMeeting)
	g.PATCH("/cancel",r.cancelMeeting)

}

func (r *MeetingRoute)inviteMeeting(c echo.Context) error {

	var m meeting.Meeting
	if err := c.Bind(&m); err != nil {
		return err
	}
	projectID, err := strconv.Atoi(m.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	MSLink := ""
	Other := ""
	if m.IsOnline == true {
		MSLink = m.MSLink
	}
	if m.IsOther == true {
		Other = m.Other
	}
	newMeeting := meeting.Meeting{
		ProjectID: m.ProjectID,
		MeetingDate: m.MeetingDate,
		MeetingTime: m.MeetingTime,
		IsOnline: m.IsOnline,
		MSLink: MSLink,
		IsOnsite: m.IsOnsite,
		IsRoom1: m.IsRoom1,
		IsRoom2: m.IsRoom2,
		IsRoom3: m.IsRoom3,
		IsRoom4: m.IsRoom4,
		IsOther: m.IsOther,
		Other: Other,
		Note: m.Note,
		CreatedBy: userLogin.Email,
		CreatedAt: time.Now(),
	}
	err = r.MeetingService.SendMeetingInvite(&newMeeting)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	member,err := r.MemberService.GetMembersOfPj(m.ProjectID)
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
		err := r.MeetingService.AlertInviteMeeting(&mbs,&newMeeting,PMName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		stage := "stage_2"
		err = r.StageService.UpdateStage(m.ProjectID,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

	}

	return c.JSON(http.StatusCreated,m)
}

func (r *MeetingRoute)getMeetingInfo(c echo.Context) error {
	projectId := c.Param("projectId")
	info,err := r.MeetingService.GetMeetingInfo(projectId)
	if err != nil {
		log.Println("err",err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK,info)
}

func (r *MeetingRoute)editInviteMetting(c echo.Context) error {
	var m meeting.Meeting
	if err := c.Bind(&m); err != nil {
		return err
	}
	projectID, err := strconv.Atoi(m.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	MSLink := ""
	Other := ""
	if m.IsOnline == true {
		MSLink = m.MSLink
	}
	if m.IsOther == true {
		Other = m.Other
	}
	newMeeting := meeting.Meeting{
		ProjectID: m.ProjectID,
		MeetingDate: m.MeetingDate,
		MeetingTime: m.MeetingTime,
		IsOnline: m.IsOnline,
		MSLink: MSLink,
		IsOnsite: m.IsOnsite,
		IsRoom1: m.IsRoom1,
		IsRoom2: m.IsRoom2,
		IsRoom3: m.IsRoom3,
		IsRoom4: m.IsRoom4,
		IsOther: m.IsOther,
		Other: Other,
		Note: m.Note,
		UpdatedBy: null.StringFrom(userLogin.Email),
		UpdatedAt: null.TimeFrom(time.Now()),
	}
	err = r.MeetingService.SendEditMeetingInvite(&newMeeting)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	member,err := r.MemberService.GetMembersOfPj(m.ProjectID)
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
		err := r.MeetingService.AlertEditInviteMeeting(&mbs,&newMeeting,PMName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated,m)
}

func (r *MeetingRoute)bypassMeeting(c echo.Context) error {

	userLogin := ParseJWTCustomClaims(c.Get("user"))
	var m meeting.Meeting
	if err := c.Bind(&m); err != nil {
		return err
	}
	newMeeting := meeting.Meeting{
		ProjectID: m.ProjectID,
		Note: m.Note,
		CreatedBy: userLogin.Email,
		ISByPass: true,
		CreatedAt: time.Now(),
	}
	err := r.MeetingService.ByPassMeeting(&newMeeting)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}
	stage := "stage_2"
	err = r.StageService.UpdateStage(m.ProjectID,stage)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated,m)
}

func (r *MeetingRoute)cancelMeeting(c echo.Context) error {
	var m meeting.Meeting
	if err := c.Bind(&m); err != nil {
		return err
	}
	userLogin := ParseJWTCustomClaims(c.Get("user"))
	newMeeting := meeting.Meeting{
		ProjectID: m.ProjectID,
		Note: m.Note,
		ISCancel: true,
		UpdatedBy: null.StringFrom(userLogin.Email),
		UpdatedAt: null.TimeFrom(time.Now()),
	}
	err := r.MeetingService.CancelMeeting(&newMeeting)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	
	projectID, err := strconv.Atoi(m.ProjectID)
	pj,err := r.ProjectService.GetProjectByID(int64(projectID))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	member,err := r.MemberService.GetMembersOfPj(m.ProjectID)
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
		err := r.MeetingService.AlertCancelInviteMeeting(&mbs,&newMeeting,PMName,pj.Name)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated,m)
}
