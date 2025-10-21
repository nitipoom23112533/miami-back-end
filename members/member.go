package members

import (
	"context"
	"fmt"
	"miami-back-end/mg"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/mailgun/mailgun-go/v4"
	"gopkg.in/guregu/null.v4"
)

type MemberService struct {
	MemberRepository *MemberRepository

}
func NewMemberService() *MemberService {
	return &MemberService{}

}

type Employee struct {

	UID       string `db:"uid" json:"uid"`
	Firstname string `db:"firstname" json:"firstname"`
	Lastname  string `db:"lastname" json:"lastname"`
	Position  string `db:"position" json:"position"`
	Email     string `db:"email" json:"email"`
	Status    string `db:"status" json:"status"`
	
}

type MembersOfPj struct {
	ProjectID string `db:"project_id" json:"project_id"`
	UID       string `db:"uid" json:"uid"`
	Firstname string `db:"firstname" json:"firstname"`
	Lastname  string `db:"lastname" json:"lastname"`
	Position  string `db:"position" json:"position"`
	Role 	  string `db:"role" json:"role"`
	Member_status string `db:"member_status" json:"member_status"`
	CreatedAt time.Time   `db:"created_at" json:"createdAt"`
	CreatedBy string `db:"created_by" json:"createdBy"`
	Email     string `db:"email" json:"email"`
	IsSendEmail bool `db:"is_send_email" json:"is_send_email"`
	UpdatedAt null.Time   `db:"updated_at" json:"updatedAt"`
	UpdatedBy null.String `db:"updated_by" json:"updatedBy"`

	
}
type AddMembers struct{

	ProjectID string `db:"project_id" json:"project_id"`
	UID []string `db:"uid" json:"uid"`
}
type PjNameOfPj struct {
    PjName string `json:"PjName"`
}

// ValidateCreate func
func (mop *MembersOfPj) ValidateCreate() error {
	return validation.ValidateStruct(mop,
		validation.Field(&mop.ProjectID, validation.Required),
		validation.Field(&mop.UID, validation.Required),
		validation.Field(&mop.Position, validation.Required),
		validation.Field(&mop.CreatedAt, validation.Required),
		validation.Field(&mop.CreatedBy, validation.Required),
	)
}
func (ms *MemberService)GetMemberByPosition() ([]Employee,error)  {

	return ms.MemberRepository.GetMemberByPosition()

}
func (ms *MemberService)GetMembersOfPj(Pjid string) ([]MembersOfPj,error)  {

	return ms.MemberRepository.GetMemberByPjID(Pjid)

}
func (ms *MemberService)GetMemberAllByPjID(Pjid string) ([]MembersOfPj,error)  {

	return ms.MemberRepository.GetMemberAllByPjID(Pjid)

}
func (ms *MemberService)GetOutOfMemberAllByPjID(Pjid string,Uid []string) ([]MembersOfPj,error)  {

	return ms.MemberRepository.GetOutOfMemberAllByPjID(Pjid,Uid)

}
func (ms *MemberService)AddMemberByPjID(mops *[]MembersOfPj) error{

	return ms.MemberRepository.AddMemberByPjID(mops)
}
func (ms *MemberService)SendMailToMembersActive(mops *[]MembersOfPj, PjName string ,PMName string) error{

	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง กำหนดทีมงานโครงการ"

	toList,err := ms.EmailMemberList(mops)
	if err != nil {
		return err
	}

	plainBody := fmt.Sprintf("%s ได้ทำการตั้งโครงการ %s และได้กำหนดให้คุณเป็นผู้ร่วมโครงการ\n",PMName,PjName)
	htmlBody := fmt.Sprintf(`<p>%s ได้ทำการตั้งโครงการ <strong>%s</strong> และได้กำหนดให้คุณเป็นผู้ร่วมโครงการ</p>`,PMName, PjName)

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}
func (ms *MemberService)SendMailToMembersInactive(mops *[]MembersOfPj, PjName string,PMName string) error{

	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง นำทีมงานออกจากโครงการ"

	toList,err := ms.EmailMemberList(mops)
	if err != nil {
		return err
	}

	plainBody := fmt.Sprintf("%s นำคุณออกจากโครงการ %s แล้ว",PMName,PjName)
	htmlBody := fmt.Sprintf(`<p>%s นำคุณออกจากโครงการ <strong>%s</strong> แล้ว</p>`,PMName, PjName)

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}

func (ms *MemberService)UpdateIsSendEmail(mops *[]MembersOfPj) error  {

	return ms.MemberRepository.UpdateIsSendEmail(mops)
}

func (ms *MemberService)EmailMemberList(mbr *[]MembersOfPj) ([]string,error) {
	var emailList []string

	for _, x := range *mbr {
		emailList = append(emailList, x.Email)
	}

	return emailList, nil
}
