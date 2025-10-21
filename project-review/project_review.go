package projectreview

import(
	"time"
	"gopkg.in/guregu/null.v4"
	"miami-back-end/members"
	"miami-back-end/meeting"
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"miami-back-end/mg"
)

type ProjectReviewService struct{
	ProjectReviewRepo *ProjectReviewRepo
	MeetingService *meeting.MeetingService
	MemberService *members.MemberService
}


func NewProjectReviewService() *ProjectReviewService{
	return &ProjectReviewService{}

}

type ProjectReview struct{

	ProjectID   string      `db:"project_id" json:"project_id"`
	MeetingDate null.String `db:"meeting_date" json:"meeting_date"`
	MeetingTime null.String `db:"meeting_time" json:"meeting_time"`
	IsOnline    bool        `db:"is_online" json:"is_online"`
	MSLink      null.String `db:"ms_link" json:"ms_link"`
	IsOnsite    bool        `db:"is_onsite" json:"is_onsite"`
	IsRoom1     bool        `db:"is_room1" json:"is_room1"`
	IsRoom2     bool        `db:"is_room2" json:"is_room2"`
	IsRoom3     bool        `db:"is_room3" json:"is_room3"`
	IsRoom4     bool        `db:"is_room4" json:"is_room4"`
	IsOther     bool        `db:"is_other" json:"is_other"`
	Other       null.String `db:"other" json:"other"`
	Note        null.String `db:"note" json:"note"`
	CreatedBy   null.String `db:"created_by" json:"created_by"`
	CreatedAt   time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt   null.Time   `db:"updated_at" json:"updatedAt"`
	UpdatedBy   null.String `db:"updated_by" json:"updatedBy"`
	ISByPass    bool        `db:"is_bypass" json:"is_bypass"`
	ISCancel 	bool 		`db:"is_cancel" json:"is_cancel"`
}

func (pr *ProjectReviewService)CreateProjectReview(prw *ProjectReview) error{
	return pr.ProjectReviewRepo.CreateProjectReview(prw)
}

func (pr *ProjectReviewService)EditProjectReview(prw *ProjectReview) error{
	return pr.ProjectReviewRepo.EditProjectReview(prw)
}

func (pr *ProjectReviewService)GetProjectReview(PjID string) (*ProjectReview,error){
	return pr.ProjectReviewRepo.GetProjectReview(PjID)
}

func (pr *ProjectReviewService)CancelProjectReview(prw *ProjectReview) error{
	return pr.ProjectReviewRepo.CancelProjectReview(prw)
}

func (pr *ProjectReviewService)ByPassProjectReview(prw *ProjectReview) error{
	return pr.ProjectReviewRepo.ByPassProjectReview(prw)
}

func (pr *ProjectReviewService)AlertCreateProjectReview(member *[]members.MembersOfPj,prw *ProjectReview, PMName string,PjName string) error{

	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง กำหนดวัน Project Review โครงการ"

	toList,err := pr.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	dateParsed, _ := time.Parse("2006-01-02", prw.MeetingDate.String)
	MeetingDate := dateParsed.Format("02/01/2006")

	ISRoom := meeting.Meeting{IsRoom1: prw.IsRoom1,IsRoom2: prw.IsRoom2,IsRoom3: prw.IsRoom3,IsRoom4: prw.IsRoom4,IsOther: prw.IsOther,Other: prw.Other.String}
	selectMeetingRoom,err := pr.MeetingService.SelectMeetingRoom(&ISRoom)
	if err != nil {
		return err
	}
	plainBody := fmt.Sprintf("%s ได้กำหนดวัน Project Review โครงการ %s ในวันที่ %s เวลา %s น.\nที่ %s\n",
	PMName, PjName, MeetingDate, prw.MeetingTime.String, selectMeetingRoom)

	htmlBody := fmt.Sprintf(`<p>%s ได้กำหนดวัน <strong>Project Review</strong> โครงการ <strong>%s</strong> ในวันที่ <strong>%s</strong> เวลา <strong>%s น.</strong></p>
	<p>ที่ <strong>%s</strong></p>`,PMName, PjName, MeetingDate, prw.MeetingTime.String, selectMeetingRoom)

	if prw.IsOnline {
		plainBody += "หรือ ทาง MS Team ตาม Link นี้: " + prw.MSLink.String + "\n"
		htmlBody += fmt.Sprintf(`<p>หรือ ทาง MS Team ตาม Link นี้: <a href="%s">%s</a></p>`, prw.MSLink.String, prw.MSLink.String)
	}

	
	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}

func (pr *ProjectReviewService)AlertEditProjectReview(member *[]members.MembersOfPj,prw *ProjectReview, PMName string,PjName string) error{

	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง เปลี่ยนแปลงวัน Project Review โครงการ"

	toList,err := pr.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	dateParsed, _ := time.Parse("2006-01-02", prw.MeetingDate.String)
	MeetingDate := dateParsed.Format("02/01/2006")

	ISRoom := meeting.Meeting{IsRoom1: prw.IsRoom1,IsRoom2: prw.IsRoom2,IsRoom3: prw.IsRoom3,IsRoom4: prw.IsRoom4,IsOther: prw.IsOther,Other: prw.Other.String}
	selectMeetingRoom,err := pr.MeetingService.SelectMeetingRoom(&ISRoom)
	if err != nil {
		return err
	}
	plainBody := fmt.Sprintf("%s ได้เปลี่ยนแปลงวัน Project Review โครงการ %s เป็นวันที่ %s เวลา %s น.\nที่ %s\n",
	PMName, PjName, MeetingDate, prw.MeetingTime.String, selectMeetingRoom)

	htmlBody := fmt.Sprintf(`<p>%s ได้เปลี่ยนแปลงวัน <strong>Project Review</strong> โครงการ <strong>%s</strong> เป็นวันที่ <strong>%s</strong> เวลา <strong>%s น.</strong></p>
	<p>ที่ <strong>%s</strong></p>`,PMName, PjName, MeetingDate, prw.MeetingTime.String, selectMeetingRoom)

	if prw.IsOnline {
		plainBody += "หรือ ทาง MS Team ตาม Link นี้: " + prw.MSLink.String + "\n"
		htmlBody += fmt.Sprintf(`<p>หรือ ทาง MS Team ตาม Link นี้: <a href="%s">%s</a></p>`, prw.MSLink.String, prw.MSLink.String)
	}

	
	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}

func (pr *ProjectReviewService)AlertCancelProjectReview(member *[]members.MembersOfPj,prw *ProjectReview, PMName string,PjName string) error{
	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง ยกเลิก Project Review โครงการ"

	toList,err := pr.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	plainBody := fmt.Sprintf("%s ได้ยกเลิก Project Review โครงการ %s โดยจะแจ้งให้ทราบอีกครั้งเมื่อทำการนัดหมายใหม่\n",PMName,PjName)
	htmlBody := fmt.Sprintf(`<p>%s ได้ยกเลิก <strong>Project Review</strong> โครงการ <strong>%s</strong> โดยจะแจ้งให้ทราบอีกครั้งเมื่อทำการนัดหมายใหม่</p>`,PMName,PjName)

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}



