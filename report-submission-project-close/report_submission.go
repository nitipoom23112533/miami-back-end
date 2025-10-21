package reportsubmissionprojectclose

import (
	"context"
	"fmt"
	"miami-back-end/members"
	"miami-back-end/mg"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"gopkg.in/guregu/null.v4"
)

type ReportSubmissionService struct{
	ReportSubmissionRepo *ReportSubmissionRepo
	MemberService *members.MemberService
}

type ReportSubmission struct{

	ProjectID 		string `db:"project_id" json:"project_id"`
	SubmissionDate 	null.String `db:"submission_date" json:"submission_date"`
	CloseDate 		null.String `db:"close_date" json:"close_date"`
	CreatedBy 		null.String `db:"created_by" json:"created_by"`
	CreatedAt 		null.Time `db:"created_at" json:"created_at"`
	UpdatedBy 		null.String `db:"updated_by" json:"updated_by"`
	UpdatedAt 		null.Time `db:"updated_at" json:"updated_at"`

}


func NewReportSubmissionService() *ReportSubmissionService{
	return &ReportSubmissionService{}
}

func (rs *ReportSubmissionService)GetReportSubmission(PjID string) (ReportSubmission,error){
	return rs.ReportSubmissionRepo.GetReportSubmission(PjID)
}

func (rs *ReportSubmissionService)InsertReportSubmission(submission *ReportSubmission) error{
	return rs.ReportSubmissionRepo.InsertReportSubmission(submission)
}

func (rs *ReportSubmissionService)UpdateReportSubmission(submission *ReportSubmission) error{
	return rs.ReportSubmissionRepo.UpdateReportSubmission(submission)
}

func (rs *ReportSubmissionService)AlertSubmissionClose(member *[]members.MembersOfPj,rss *ReportSubmission, PMName string,PjName string,rstype string) error{
	sender := "BRS Admin Official <noreply@brsth.com>"

	toList,err := rs.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	var submissionDateParsed time.Time
	var closeDateParsed time.Time

	if rss.SubmissionDate.Valid {
		submissionDateParsed, _ = time.Parse("2006-01-02", rss.SubmissionDate.String)
	} 
	if rss.CloseDate.Valid {
		closeDateParsed, _ = time.Parse("2006-01-02", rss.CloseDate.String)
	}
	submissionDate := submissionDateParsed.Format("02/01/2006")
	closeDate := closeDateParsed.Format("02/01/2006")

	var plainBody string
	var htmlBody string
	var subject string

	switch rstype {
		case "report":
			subject = "Miami เรื่อง Report Submission โครงการ"
			plainBody = fmt.Sprintf(`%s แจ้งให้ทีมงานทราบว่าได้ส่งมอบ Report โครงการ %s แก่ลูกค้าแล้วในวันที่ %s`,PMName, PjName, submissionDate)
			htmlBody = fmt.Sprintf(`<P>%s แจ้งให้ทีมงานทราบว่าได้ส่งมอบ Report โครงการ <strong>%s</strong> แก่ลูกค้าแล้วในวันที่ <strong>%s</strong></P>`,PMName, PjName, submissionDate)
		
		case "close":
			subject = "Miami เรื่อง โครงการ เสร็จสิ้น"
			plainBody = fmt.Sprintf(`%s แจ้งให้ทีมงานทราบว่าโครงการ %s ได้เสร็จสิ้นในวันที่ %s ขอบคุณสำหรับความร่วมมือของทีมงานทุกท่าน`,PMName, PjName, closeDate)
			htmlBody = fmt.Sprintf(`<P>%s แจ้งให้ทีมงานทราบว่าโครงการ <strong>%s</strong> ได้เสร็จสิ้นในวันที่ <strong>%s</strong> ขอบคุณสำหรับความร่วมมือของทีมงานทุกท่าน</P>`,PMName, PjName, closeDate)
	}
	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}