package pilotquestionnaire

import (
	"context"
	"fmt"
	"miami-back-end/meeting"
	"miami-back-end/members"
	"miami-back-end/mg"
	"time"
	"github.com/mailgun/mailgun-go/v4"
	"gopkg.in/guregu/null.v4"
)


type PilotQuestionnaireService struct {
	PQService *PilotQuestionnaireRepository
	MeetingService *meeting.MeetingService
	MemberService *members.MemberService

}
type PilotQuestionnaire struct {
	ProjectID 		string 			`db:"project_id" json:"project_id"`
	MeetingDate 	string 			`db:"meeting_date" json:"meeting_date"`
	MeetingTime 	string 			`db:"meeting_time" json:"meeting_time"`
	IsOnline 		bool 			`db:"is_online" json:"is_online"`
	MSLink 			string 			`db:"ms_link" json:"ms_link"`
	IsOnsite 		bool 			`db:"is_onsite" json:"is_onsite"`
	IsRoom1 		bool 			`db:"is_room1" json:"is_room1"`
	IsRoom2 		bool 			`db:"is_room2" json:"is_room2"`
	IsRoom3 		bool 			`db:"is_room3" json:"is_room3"`
	IsRoom4 		bool 			`db:"is_room4" json:"is_room4"`
	IsOther 		bool 			`db:"is_other" json:"is_other"`
	Other 			string 			`db:"other" json:"other"`
	Note 			string 			`db:"note" json:"note"`
	CreatedBy 		string 			`db:"created_by" json:"created_by"`
	CreatedAt 		time.Time 		`db:"created_at" json:"created_at"`
	UpdatedAt 		null.Time   	`db:"updated_at" json:"updatedAt"`
	UpdatedBy 		null.String 	`db:"updated_by" json:"updatedBy"`
	ISByPass 		bool 			`db:"is_bypass" json:"is_bypass"`
	ISCancel 		bool 			`db:"is_cancel" json:"is_cancel"`
	FilePath    	[]FilePath  	`json:"file_path"`
	DetailOnChange 	null.String 	`db:"detail_on_change" json:"detail_on_change"`

}

type PilotQuestionnaireNullCase struct {
	ProjectID   string      	`db:"project_id" json:"project_id"`
	MeetingDate null.String 	`db:"meeting_date" json:"meeting_date"`
	MeetingTime null.String 	`db:"meeting_time" json:"meeting_time"`
	IsOnline    bool        	`db:"is_online" json:"is_online"`
	MSLink      null.String 	`db:"ms_link" json:"ms_link"`
	IsOnsite    bool        	`db:"is_onsite" json:"is_onsite"`
	IsRoom1     bool        	`db:"is_room1" json:"is_room1"`
	IsRoom2     bool        	`db:"is_room2" json:"is_room2"`
	IsRoom3     bool        	`db:"is_room3" json:"is_room3"`
	IsRoom4     bool        	`db:"is_room4" json:"is_room4"`
	IsOther     bool        	`db:"is_other" json:"is_other"`
	Other       null.String 	`db:"other" json:"other"`
	Note        null.String 	`db:"note" json:"note"`
	CreatedBy   null.String 	`db:"created_by" json:"created_by"`
	CreatedAt   time.Time   	`db:"created_at" json:"created_at"`
	UpdatedAt   null.Time   	`db:"updated_at" json:"updatedAt"`
	UpdatedBy   null.String 	`db:"updated_by" json:"updatedBy"`
	ISByPass    bool        	`db:"is_bypass" json:"is_bypass"`
	ISCancel 	bool 			`db:"is_cancel" json:"is_cancel"`
	FilePath    []FilePath  	`json:"file_path"`
	DetailOnChange null.String 	`db:"detail_on_change" json:"detail_on_change"`

}

type Path struct {
	ProjectID 		string   	`db:"project_id" json:"project_id"`
	FilePath      	[]FilePath 	`json:"file_path"`
}
type FilePath struct {
	Path     		string `db:"path" json:"path"`
	FileName 		string `db:"file_name" json:"file_name"`
	Number   		int    `db:"number" json:"number"`
	IsNew  	 		bool   `db:"is_new" json:"is_new"`
	IsSign  		bool   `db:"is_sign" json:"is_sign"`
	ISTraining  	bool   `db:"is_training" json:"is_training"`
	IsRevised 		bool   `db:"is_revised" json:"is_revised"`


}

func NewPilotQuestionnaireService() *PilotQuestionnaireService{

	return &PilotQuestionnaireService{}
}

func (p *PilotQuestionnaireService)CreatePilotQuestionnaire(pq *PilotQuestionnaire) error{
	
	return p.PQService.CreatePilotQuestionnaire(pq)

}
func (p *PilotQuestionnaireService)ByPassPilotQuestionnaire(pq *PilotQuestionnaire) error {
	return p.PQService.ByPassPilotQuestionnaire(pq)
}
func (p *PilotQuestionnaireService)GetPilotQuestionnaireInfo(PjID string) (*PilotQuestionnaireNullCase,error){
	return p.PQService.GetPilotQuestionnaireInfo(PjID)
}

func (p *PilotQuestionnaireService)Editdetail(pq *PilotQuestionnaire) error {
	return p.PQService.Editdetail(pq)
}

func (p *PilotQuestionnaireService)DetailOnChange(pq *PilotQuestionnaire,stage string) error {
	return p.PQService.DetailOnChange(pq,stage)
}
func (p *PilotQuestionnaireService)CancelPilotQuestionnaire(pq *PilotQuestionnaire) error {
	return p.PQService.CancelPilotQuestionnaire(pq)
}
func (p *PilotQuestionnaireService)GetAllPilotQuestionnairePath(PjID string) ([]*FilePath,error) {
	return p.PQService.GetAllPilotQuestionnairePath(PjID)
}
func (p *PilotQuestionnaireService)GetLatestPilotQuestionnairePath(PjID string) ([]FilePath,error) {
	return p.PQService.GetLatestPilotQuestionnairePath(PjID)
}
func (p *PilotQuestionnaireService)InsertPathPilotQuestionnaire(pp *Path) error {
	return p.PQService.InsertPathPilotQuestionnaire(pp)
}
func (p *PilotQuestionnaireService)UpdatePathPilotQuestionnaire(pp *Path,fileNumber int) error {
	return p.PQService.UpdatePathPilotQuestionnaire(pp,fileNumber)
}

func (p *PilotQuestionnaireService)AlertInformPilotQuestionnaire(member *[]members.MembersOfPj,pq *PilotQuestionnaire, PMName string,PjName string,pthlt []FilePath) error{
	
	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง กำหนดวัน Pilot แบบสอบถามโครงการ"

	toList,err := p.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	// var senTo strings.Builder

	// for i, x := range *member {
	// 	senTo.WriteString(x.Email)
	// 	if i < len(*member)-1 {
	// 		senTo.WriteString(",")
	// 	}
	// }
	// recipient := senTo.String()
	// toList := strings.Split(recipient, ",")

	// toList := []string{"nitipoom.sa@brsth.com"}

	dateParsed, _ := time.Parse("2006-01-02", pq.MeetingDate)
	MeetingDate := dateParsed.Format("02/01/2006")

	IsRoom := meeting.Meeting{IsRoom1: pq.IsRoom1,IsRoom2: pq.IsRoom2,IsRoom3: pq.IsRoom3,IsRoom4: pq.IsRoom4,IsOther: pq.IsOther,Other: pq.Other}
	selectMeetingRoom,err := p.MeetingService.SelectMeetingRoom(&IsRoom)
	if err != nil {
		return err
	}
	selectFilePath,err := p.SelectFilePath(pthlt)
	if err != nil {
		return err
	}
	plainBody := fmt.Sprintf("%s ได้กำหนดวัน Pilot แบบสอบถามโครงการ %s ในวันที่ %s เวลา %s น. ที่ %s\n",PMName, PjName, MeetingDate, pq.MeetingTime, selectMeetingRoom)
	htmlBody := fmt.Sprintf(`<p>%s ได้กำหนดวัน  <strong>Pilot</strong> แบบสอบถามโครงการ <strong>%s</strong> ในวันที่ <strong>%s</strong> เวลา <strong>%s น.</strong></p>
	<p>ที่ <strong>%s</strong></p>`,PMName, PjName, MeetingDate, pq.MeetingTime, selectMeetingRoom)

	if pq.IsOnline {
		plainBody += "หรือ ทาง MS Team ตาม Link นี้: " + pq.MSLink + "\n"
		htmlBody += fmt.Sprintf(`<p>หรือ ทาง MS Team ตาม Link นี้: <a href="%s">%s</a></p>`, pq.MSLink, pq.MSLink)
	}

	if selectFilePath != "" {
		htmlBody += fmt.Sprintf(`<p>โดย Version ที่จะใช้ Pilot คือ:<br> %s</p>`, selectFilePath)
		plainBody += "โดย Version ที่จะใช้ Pilot คือ: "+ "\n" + selectFilePath + "\n"
	}

	
	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err

}

func (p *PilotQuestionnaireService)AlertEditDetail(member *[]members.MembersOfPj,pq *PilotQuestionnaire, PMName string,PjName string) error{

	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง เปลี่ยนแปลงวันและเวลา Pilot โครงการ"

	toList,err := p.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	dateParsed, _ := time.Parse("2006-01-02", pq.MeetingDate)
	MeetingDate := dateParsed.Format("02/01/2006")

	IsRoom := meeting.Meeting{IsRoom1: pq.IsRoom1,IsRoom2: pq.IsRoom2,IsRoom3: pq.IsRoom3,IsRoom4: pq.IsRoom4,IsOther: pq.IsOther,Other: pq.Other}
	selectMeetingRoom,err := p.MeetingService.SelectMeetingRoom(&IsRoom)
	if err != nil {
		return err
	}

	plainBody := fmt.Sprintf("%s ได้เปลี่ยนแปลงวันและเวลา Pilot โครงการ %s เป็นวันที่ %s เวลา %s น.\n",PMName, PjName, MeetingDate, pq.MeetingTime)
	plainBody += fmt.Sprintf("ที่ %s\n",selectMeetingRoom)
	htmlBody := fmt.Sprintf(`<p>%s ได้เปลี่ยนแปลงวันและเวลา <strong>Pilot</strong> โครงการ <strong>%s</strong> เป็นวันที่ <strong>%s</strong> เวลา <strong>%s น.</strong></p>`,PMName, PjName, MeetingDate, pq.MeetingTime)
	htmlBody += fmt.Sprintf(`<p>ที่ <strong>%s</strong></p>`,selectMeetingRoom)


	if pq.IsOnline == true {
		plainBody += fmt.Sprintf("หรือ ทาง MS Team ตาม Link นี้ %s\n",pq.MSLink)
		htmlBody += fmt.Sprintf(`<p>หรือ ทาง MS Team ตาม Link นี้ <a href="%s">%s</a></p>`, pq.MSLink, pq.MSLink)
	}

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
	
}

func (p *PilotQuestionnaireService)AlertEditQuestion(member *[]members.MembersOfPj,pq *PilotQuestionnaire, PMName string,PjName string,pthlt []FilePath) error{
	
	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง เปลี่ยนแปลงชุดในการ Pilot โครงการ"

	toList,err := p.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	selectFilePath,err := p.SelectFilePath(pthlt)
	if err != nil {
		return err
	}

	plainBody := fmt.Sprintf("โครงการ %s มีการเปลี่ยนแปลงชุดที่จะใช้ในการ Pilot\n",PjName)
	htmlBody := fmt.Sprintf(`<p>โครงการ <strong>%s</strong> มีการเปลี่ยนแปลงชุดที่จะใช้ในการ <strong>Pilot</strong></p>`,PjName)

	if selectFilePath != "" {
		htmlBody += fmt.Sprintf(`<p>โดย Version ที่จะใช้ Pilot ล่าสุด คือ:<br>%s</p>`, selectFilePath)
		plainBody += "โดย Version ที่จะใช้ Pilot ล่าสุด คือ: "+ "\n" + selectFilePath + "\n"
	}
	
	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err

}

func (p *PilotQuestionnaireService)AlertCancelPilotQuestionnaire(member *[]members.MembersOfPj,pq *PilotQuestionnaire, PMName string,PjName string) error{
	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง ยกเลิกการ Pilot โครงการ"

	toList,err := p.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	plainBody := fmt.Sprintf("%s ได้ยกเลิกการ Pilot โครงการ %s โดยจะแจ้งให้ทราบอีกครั้งเมื่อทำการนัดหมายใหม่\n",PMName,PjName)
	htmlBody := fmt.Sprintf(`<p>%s ได้ยกเลิกการ <strong>Pilot</strong> โครงการ <strong>%s</strong> โดยจะแจ้งให้ทราบอีกครั้งเมื่อทำการนัดหมายใหม่</p>`,PMName,PjName)

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}

func (p *PilotQuestionnaireService)SelectFilePath(pthlt []FilePath) (string,error){
	var selectFilePath string
	for _, x := range pthlt {
		if x.FileName == "" {
			continue
		}
		selectFilePath += fmt.Sprintf(`<a href="%s">%s</a><br>`, x.Path, x.FileName)
	}
	return selectFilePath,nil
}