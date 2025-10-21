package api

import (
	"log"
	"miami-back-end/members"
	"miami-back-end/stage"
	"net/http"
	"time"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"

)
type MemberRoute struct{
	memberService *members.MemberService
	stageService *stage.StageService
}

func NewMemberRoute(memberService *members.MemberService,stagerService *stage.StageService) *MemberRoute{
	return &MemberRoute{
		memberService: memberService,
		stageService: stagerService,
	}
}

func (r *MemberRoute)Group(g *echo.Group) {
	g.Use(Auth())
	g.GET("/employee", r.getMember)
	g.GET("/membersOfPj/:projectId", r.getMembers)
	g.POST("/addMember", r.addMember)
	g.POST("/comfirmMember/:projectId", r.confirmMember)
}

func (r *MemberRoute)getMember(c echo.Context) error{
	
	m,err := r.memberService.GetMemberByPosition();
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, m)
}

func (r *MemberRoute)getMembers(c echo.Context) error{
	projectId := c.Param("projectId")
	m,err := r.memberService.GetMembersOfPj(projectId)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(http.StatusOK, m)
}
func (r *MemberRoute)addMember(c echo.Context) error{
	
	var mop members.AddMembers
	if err := c.Bind(&mop); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	userLogin := ParseJWTCustomClaims(c.Get("user"))

	almbr,err := r.memberService.GetMemberAllByPjID(mop.ProjectID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var mops []members.MembersOfPj
	var newMember members.MembersOfPj
	for _, UID := range mop.UID {
		for _, x := range almbr {
			if x.UID == UID {
			newMember = members.MembersOfPj{
					ProjectID: mop.ProjectID,
					UID: UID,
					Position: x.Position,
					Member_status: "active",
					IsSendEmail: false,
					CreatedAt :	x.CreatedAt,
					CreatedBy: x.CreatedBy,
					UpdatedAt: null.TimeFrom(time.Now()),
					UpdatedBy: null.StringFrom(userLogin.Email),
				}
			}else{
				newMember = members.MembersOfPj{
					ProjectID: mop.ProjectID,
					UID: UID,
					Position: "M",
					Member_status: "active",
					IsSendEmail: false,
					CreatedAt : time.Now(),
					CreatedBy: userLogin.Email,
				}
			}
		}
		
		if err := newMember.ValidateCreate(); err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		mops = append(mops, newMember)
	}

	outOfMbr,err := r.memberService.GetOutOfMemberAllByPjID(mop.ProjectID,mop.UID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var Member_status string
	var IsSendEmail bool
	for _, x := range outOfMbr {
		if x.Role == "Project Manager" {
			Member_status = x.Member_status
			IsSendEmail = x.IsSendEmail
		}else{
			Member_status = "inactive"
			IsSendEmail = false
		}
		newOotOfMember := members.MembersOfPj{
			ProjectID: x.ProjectID,
			UID: x.UID,
			Position: x.Position,
			Member_status: Member_status,
			IsSendEmail: IsSendEmail,
			CreatedAt : x.CreatedAt,
			CreatedBy: x.CreatedBy,
			UpdatedAt: null.TimeFrom(time.Now()),
			UpdatedBy: null.StringFrom(userLogin.Email),
		}
		if err := newOotOfMember.ValidateCreate(); err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		mops = append(mops, newOotOfMember)
	}

	if len(mops) > 0 {
		err := r.memberService.AddMemberByPjID(&mops)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusCreated, mops)
}

func (r *MemberRoute) confirmMember(c echo.Context) error {
	
	projectId := c.Param("projectId")
	var PjName members.PjNameOfPj
	if err := c.Bind(&PjName); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	
	m, err := r.memberService.GetMemberAllByPjID(projectId)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var mActive []members.MembersOfPj
	var mInactive []members.MembersOfPj
	var PMName string
	for _, x := range m {
		if x.Position == "PO" {
			PMName = x.Firstname + " " + x.Lastname
		}
		if x.Member_status == "active" && !x.IsSendEmail {
			mActive = append(mActive, members.MembersOfPj{
				ProjectID:   x.ProjectID,
				UID:         x.UID,
				Firstname:   x.Firstname,
				Lastname:    x.Lastname,
				Email:       x.Email,
				IsSendEmail: true,
			})
		} else if x.Member_status == "inactive" && !x.IsSendEmail {
			mInactive = append(mInactive, members.MembersOfPj{
				ProjectID:   x.ProjectID,
				UID:         x.UID,
				Firstname:   x.Firstname,
				Lastname:    x.Lastname,
				Email:       x.Email,
				IsSendEmail: true,
			})
		}
	}

	if len(mActive) > 0 {
		if err := r.memberService.SendMailToMembersActive(&mActive,PjName.PjName,PMName); err != nil {
			log.Println("error sending active emails:", err)
			return c.JSON(http.StatusInternalServerError, "failed to send active member emails")
		}

		if err := r.memberService.UpdateIsSendEmail(&mActive); err != nil {
			log.Println("error updating IsSendEmail")
			return c.JSON(http.StatusInternalServerError, "failed to update IsSendEmail")
		}
		stage := "stage_1"
		err := r.stageService.UpdateStage(projectId,stage)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	if len(mInactive) > 0 {
		if err := r.memberService.SendMailToMembersInactive(&mInactive,PjName.PjName,PMName); err != nil {
			log.Println("error sending inactive emails:", err)
			return c.JSON(http.StatusInternalServerError, "failed to send inactive member emails")
		}

		if err := r.memberService.UpdateIsSendEmail(&mInactive); err != nil {
			log.Println("error updating IsSendEmail")
			return c.JSON(http.StatusInternalServerError, "failed to update IsSendEmail")
		}
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "success"})

}