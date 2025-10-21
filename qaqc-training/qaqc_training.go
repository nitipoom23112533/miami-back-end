package qaqctraining
import(
	"time"
	"gopkg.in/guregu/null.v4"
	"miami-back-end/pilot-questionnaire"
	"miami-back-end/members"
	"miami-back-end/meeting"
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"miami-back-end/mg"
)

type QAQCService struct {
	QAQCRepo *QAQCRepository
	MeetingService *meeting.MeetingService
	PQService *pilotquestionnaire.PilotQuestionnaireService
	MemberService *members.MemberService
}

func NewQAQCService() *QAQCService {
	return &QAQCService{}
}

type QAQCTraining struct {
	ProjectID string 		`db:"project_id" json:"project_id"`
	MeetingDate string 		`db:"meeting_date" json:"meeting_date"`
	MeetingTime string 		`db:"meeting_time" json:"meeting_time"`
	IsOnline bool 			`db:"is_online" json:"is_online"`
	MSLink string 			`db:"ms_link" json:"ms_link"`
	IsOnsite bool 			`db:"is_onsite" json:"is_onsite"`
	IsRoom1 bool 			`db:"is_room1" json:"is_room1"`
	IsRoom2 bool 			`db:"is_room2" json:"is_room2"`
	IsRoom3 bool 			`db:"is_room3" json:"is_room3"`
	IsRoom4 bool 			`db:"is_room4" json:"is_room4"`
	IsOther bool 			`db:"is_other" json:"is_other"`
	Other string 			`db:"other" json:"other"`
	Note string 			`db:"note" json:"note"`
	CreatedBy string 		`db:"created_by" json:"created_by"`
	CreatedAt time.Time 	`db:"created_at" json:"created_at"`
	UpdatedAt null.Time   	`db:"updated_at" json:"updatedAt"`
	UpdatedBy null.String 	`db:"updated_by" json:"updatedBy"`
	ISByPass bool 			`db:"is_bypass" json:"is_bypass"`
	ISCancel bool 			`db:"is_cancel" json:"is_cancel"`
	FilePath    []pilotquestionnaire.FilePath  	`json:"file_path"`
	DetailOnChange null.String 	`db:"detail_on_change" json:"detail_on_change"`
}

type QAQCTrainingNullCase struct {
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
	FilePath    []pilotquestionnaire.FilePath  	`json:"file_path"`
	DetailOnChange null.String 	`db:"detail_on_change" json:"detail_on_change"`
}

func (qq *QAQCService)SendCreateQAQCTraining(q *QAQCTraining) error {
	return qq.QAQCRepo.SendCreateQAQCTraining(q)

}

func (qq *QAQCService)InsertQAQCTrainingPath(p *pilotquestionnaire.Path) error {
	return qq.QAQCRepo.InsertQAQCTrainingPath(p)
}

func (qq *QAQCService)GetQAQCTrainingInfo(PjId string) (*QAQCTrainingNullCase,error){
	return qq.QAQCRepo.GetQAQCTrainingInfo(PjId)
}

func (qq *QAQCService)GetAllQAQCTrainingPath(PjID string) ([]*pilotquestionnaire.FilePath,error){
	return qq.QAQCRepo.GetAllQAQCTrainingPath(PjID)
}

func (qq *QAQCService)QAQCEditdetail(q *QAQCTraining) error {
	return qq.QAQCRepo.QAQCEditdetail(q)
}

func (qq *QAQCService)UpdatePathQAQCTraining(p *pilotquestionnaire.Path,fileNumber int) error {
	return qq.QAQCRepo.UpdatePathQAQCTraining(p,fileNumber)
}

func (qq *QAQCService)CancelQAQCTraining(q *QAQCTraining) error {
	return qq.QAQCRepo.CancelQAQCTraining(q)
}

func (qq *QAQCService)ByPassQAQCTraining(q *QAQCTraining) error {
	return qq.QAQCRepo.ByPassQAQCTraining(q)
}


func (qq *QAQCService)GetLatestQAQCTrainingPath(PjID string) ([]pilotquestionnaire.FilePath,error) {
	return qq.QAQCRepo.GetLatestQAQCTrainingPath(PjID)
}


func (qq *QAQCService)AlertQAQCTraining(member *[]members.MembersOfPj,q *QAQCTraining, PMName string,QaqcName []string,PjName string,pthlt []pilotquestionnaire.FilePath) error{

	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง กำหนดวัน Logic Check QA/QC โครงการ"

	toList,err := qq.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	dateParsed, _ := time.Parse("2006-01-02", q.MeetingDate)
	MeetingDate := dateParsed.Format("02/01/2006")

	var qaqcname string 
	for i, x := range QaqcName {
		if i != len(QaqcName)-1 {
			qaqcname += x + " และ "

		}else {
			qaqcname += x
		}
	}

	var name string
	if len(QaqcName) > 0 {
		name = qaqcname
	}else{
		name = PMName
	}

	ISRoom := meeting.Meeting{IsRoom1: q.IsRoom1,IsRoom2: q.IsRoom2,IsRoom3: q.IsRoom3,IsRoom4: q.IsRoom4,IsOther: q.IsOther,Other: q.Other}
	selectMeetingRoom,err := qq.MeetingService.SelectMeetingRoom(&ISRoom)
	if err != nil {
		return err
	}

	selectFilePath,err := qq.PQService.SelectFilePath(pthlt)
	if err != nil {
		return err
	}

	plainBody := fmt.Sprintf("%s ได้กำหนดวัน training แบบสอบถาม Logic Check ของ QA/QC โครงการ %s ในวันที่ %s เวลา %s น.\nที่ %s\n",
	name, PjName, MeetingDate, q.MeetingTime, selectMeetingRoom)

	htmlBody := fmt.Sprintf(`<p>%s ได้กำหนดวัน <strong>training</strong> แบบสอบถาม Logic Check ของ QA/QC โครงการ <strong>%s</strong> ในวันที่ <strong>%s</strong> เวลา <strong>%s น.</strong></p>
	<p>ที่ <strong>%s</strong></p>`,name, PjName, MeetingDate, q.MeetingTime, selectMeetingRoom)

	if q.IsOnline {
		plainBody += "หรือ ทาง MS Team ตาม Link นี้: " + q.MSLink + "\n"
		htmlBody += fmt.Sprintf(`<p>หรือ ทาง MS Team ตาม Link นี้: <a href="%s">%s</a></p>`, q.MSLink, q.MSLink)
	}

	if selectFilePath != "" {
		htmlBody += fmt.Sprintf(`<p>โดย Version ที่จะใช้ training คือ:<br>%s</p>`, selectFilePath)
		plainBody += "โดย Version ที่จะใช้ training คือ: " + "\n" + selectFilePath + "\n"
	}

	
	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}

func (qq *QAQCService)AlertEditDetailQAQCTraining(member *[]members.MembersOfPj,q *QAQCTraining, PMName string,QaqcName []string,PjName string) error{

	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง เปลี่ยนแปลงวันและเวลา Logic Check QA/QC โครงการ"

	toList,err := qq.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	dateParsed, _ := time.Parse("2006-01-02", q.MeetingDate)
	MeetingDate := dateParsed.Format("02/01/2006")

	var qaqcname string 
	for i, x := range QaqcName {
		if i != len(QaqcName)-1 {
			qaqcname += x + " และ "

		}else {
			qaqcname += x
		}
	}

	var name string
	if len(QaqcName) > 0 {
		name = qaqcname
	}else{
		name = PMName
	}


	IsRoom := meeting.Meeting{IsRoom1: q.IsRoom1,IsRoom2: q.IsRoom2,IsRoom3: q.IsRoom3,IsRoom4: q.IsRoom4,IsOther: q.IsOther,Other: q.Other}
	selectMeetingRoom,err := qq.MeetingService.SelectMeetingRoom(&IsRoom)
	if err != nil {
		return err
	}

	plainBody := fmt.Sprintf("%s ได้เปลี่ยนแปลงวันและเวลา Logic Check QA/QC training โครงการ %s เป็นวันที่ %s เวลา %s น.\n",name, PjName, MeetingDate, q.MeetingTime)
	plainBody += fmt.Sprintf("ที่ %s\n",selectMeetingRoom)
	htmlBody := fmt.Sprintf(`<p>%s ได้เปลี่ยนแปลงวันและเวลา <strong>Logic Check QA/QC training</strong> โครงการ <strong>%s</strong> เป็นวันที่ <strong>%s</strong> เวลา <strong>%s น.</strong></p>`,name, PjName, MeetingDate, q.MeetingTime)
	htmlBody += fmt.Sprintf(`<p>ที่ <strong>%s</strong></p>`,selectMeetingRoom)


	if q.IsOnline == true {
		plainBody += fmt.Sprintf("หรือ ทาง MS Team ตาม Link นี้ %s\n",q.MSLink)
		htmlBody += fmt.Sprintf(`<p>หรือ ทาง MS Team ตาม Link นี้ <a href="%s">%s</a></p>`, q.MSLink, q.MSLink)
	}

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
	
}

func (qq *QAQCService)AlertEditQAQCTrainingQuestion(member *[]members.MembersOfPj,pq *pilotquestionnaire.PilotQuestionnaire, PMName string,PjName string,pthlt []pilotquestionnaire.FilePath) error{
	
	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง เปลี่ยนแปลงชุดในการ Logic Check QA/QC โครงการ"

	toList,err := qq.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	selectFilePath,err := qq.PQService.SelectFilePath(pthlt)
	if err != nil {
		return err
	}

	plainBody := fmt.Sprintf("โครงการ %s มีการเปลี่ยนแปลงชุดที่จะใช้ Logic Check Training\n",PjName)
	htmlBody := fmt.Sprintf(`<p>โครงการ <strong>%s</strong> มีการเปลี่ยนแปลงชุดที่จะใช้ <strong>Logic Check Training</strong></p>`,PjName)

	if selectFilePath != "" {
		htmlBody += fmt.Sprintf(`<p>โดย Version ที่จะใช้ training ล่าสุด คือ:<br>%s</p>`, selectFilePath)
		plainBody += "โดย Version ที่จะใช้ training ล่าสุด คือ: "+ "\n" + selectFilePath + "\n"
	}
	
	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err

}

func (qq *QAQCService)AlertCancelQAQCTraining(member *[]members.MembersOfPj,q *QAQCTraining, PMName string,QaqcName []string,PjName string) error{
	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง ยกเลิกการ Logic Check QA/QC โครงการ"

	toList,err := qq.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	var qaqcname string 
	for i, x := range QaqcName {
		if i != len(QaqcName)-1 {
			qaqcname += x + " และ "

		}else {
			qaqcname += x
		}
	}

	var name string
	if len(QaqcName) > 0 {
		name = qaqcname
	}else{
		name = PMName
	}

	plainBody := fmt.Sprintf("%s ได้ยกเลิกการ Logic Check Training โครงการ %s โดยจะแจ้งให้ทราบอีกครั้งเมื่อทำการนัดหมายใหม่\n",name,PjName)
	htmlBody := fmt.Sprintf(`<p>%s ได้ยกเลิกการ <strong>Logic Check Training</strong> โครงการ <strong>%s</strong> โดยจะแจ้งให้ทราบอีกครั้งเมื่อทำการนัดหมายใหม่</p>`,name,PjName)

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}