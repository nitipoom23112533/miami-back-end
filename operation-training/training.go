package operationtraining
import(
	"time"
	"gopkg.in/guregu/null.v4"
	"miami-back-end/members"
	"miami-back-end/meeting"
	"context"
	"github.com/mailgun/mailgun-go/v4"
	"miami-back-end/mg"
	"miami-back-end/pilot-questionnaire"
	"fmt"
)

type  OperationTrainingServive struct{
	OperationTrainingRepository *OperationTrainingRepository
	MeetingService *meeting.MeetingService
	PQService *pilotquestionnaire.PilotQuestionnaireService
	MemberService *members.MemberService
}

func NewOperationTrainingService() *OperationTrainingServive{
	return &OperationTrainingServive{}
}

type Training struct {
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

type TrainingNullCase struct {
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

func (ot *OperationTrainingServive)SendTrainingInvite(t *Training) error {

	return ot.OperationTrainingRepository.SendTrainingInvite(t)
}

func (ot *OperationTrainingServive)GetTrainingInfo(PjID string) (*TrainingNullCase,error){
	return ot.OperationTrainingRepository.GetTrainingInfo(PjID)

}
func (ot *OperationTrainingServive)GetAllPilotQuestionnairePathIsSign(PjID string) ([]*pilotquestionnaire.FilePath,error) {
	return ot.OperationTrainingRepository.GetAllPilotQuestionnairePathIsSign(PjID)
}
func (ot *OperationTrainingServive)GetTrainingInfoISSign(PjID string) (*TrainingNullCase,error) {
	return ot.OperationTrainingRepository.GetTrainingInfoISSign(PjID)
}


func (ot *OperationTrainingServive)EditDetailTraining(t *Training) error {
	return ot.OperationTrainingRepository.Editdetail(t)
}

func (ot *OperationTrainingServive)ByPassTrraining(t *Training) error {
	return ot.OperationTrainingRepository.ByPassTrraining(t)
}

func (ot *OperationTrainingServive)CancelTraining(t *Training) error {
	return ot.OperationTrainingRepository.CancelTraining(t)
}

func (ot *OperationTrainingServive)SelectFileToTraining(t *Training) error {
	return ot.OperationTrainingRepository.SelectFileToTraining(t)
}

func (ot *OperationTrainingServive)UpdateTrainingPath(pth *pilotquestionnaire.Path) error{
	return ot.OperationTrainingRepository.UpdateTrainingPath(pth)
}

func (ot *OperationTrainingServive)DetailOnChangeTraining(dtl *Training) error {
	return ot.OperationTrainingRepository.DetailOnChangeTraining(dtl)
}
func (ot *OperationTrainingServive)GetLatestTrainingPath(PjID string) ([]pilotquestionnaire.FilePath,error) {
	return ot.OperationTrainingRepository.GetLatestTrainingPath(PjID)
}


func (ot *OperationTrainingServive)AlertInviteTraining(member *[]members.MembersOfPj,t *Training, PMName string,PjName string,pthlt []pilotquestionnaire.FilePath) error{

	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง กำหนดวัน training แบบสอบถามโครงการ"

	toList,err := ot.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	dateParsed, _ := time.Parse("2006-01-02", t.MeetingDate)
	MeetingDate := dateParsed.Format("02/01/2006")

	ISRoom := meeting.Meeting{IsRoom1: t.IsRoom1,IsRoom2: t.IsRoom2,IsRoom3: t.IsRoom3,IsRoom4: t.IsRoom4,IsOther: t.IsOther,Other: t.Other}
	selectMeetingRoom,err := ot.MeetingService.SelectMeetingRoom(&ISRoom)
	if err != nil {
		return err
	}

	selectFilePath,err := ot.PQService.SelectFilePath(pthlt)
	if err != nil {
		return err
	}

	plainBody := fmt.Sprintf("%s ได้กำหนดวัน training แบบสอบถามโครงการ %s ในวันที่ %s เวลา %s น.\nที่ %s\n",
	PMName, PjName, MeetingDate, t.MeetingTime, selectMeetingRoom)

	htmlBody := fmt.Sprintf(`<p>%s ได้กำหนดวัน <strong>training</strong> แบบสอบถามโครงการ <strong>%s</strong> ในวันที่ <strong>%s</strong> เวลา <strong>%s น.</strong></p>
	<p>ที่ <strong>%s</strong></p>`,PMName, PjName, MeetingDate, t.MeetingTime, selectMeetingRoom)

	if t.IsOnline {
		plainBody += "หรือ ทาง MS Team ตาม Link นี้: " + t.MSLink + "\n"
		htmlBody += fmt.Sprintf(`<p>หรือ ทาง MS Team ตาม Link นี้: <a href="%s">%s</a></p>`, t.MSLink, t.MSLink)
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

func (ot *OperationTrainingServive)AlertEditDetailTraining(member *[]members.MembersOfPj,t *Training, PMName string,PjName string) error{

	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง เปลี่ยนแปลงวันและเวลา Training โครงการ"

	toList,err := ot.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	dateParsed, _ := time.Parse("2006-01-02", t.MeetingDate)
	MeetingDate := dateParsed.Format("02/01/2006")

	ISRoom := meeting.Meeting{IsRoom1: t.IsRoom1,IsRoom2: t.IsRoom2,IsRoom3: t.IsRoom3,IsRoom4: t.IsRoom4,IsOther: t.IsOther,Other: t.Other}
	selectMeetingRoom,err := ot.MeetingService.SelectMeetingRoom(&ISRoom)
	if err != nil {
		return err
	}

	plainBody := fmt.Sprintf("%s ได้เปลี่ยนแปลงวันและเวลา Training โครงการ %s เป็นวันที่ %s เวลา %s น.\n",PMName, PjName, MeetingDate, t.MeetingTime)
	plainBody += fmt.Sprintf("ที่ %s\n",selectMeetingRoom)
	htmlBody := fmt.Sprintf(`<p>%s ได้เปลี่ยนแปลงวันและเวลา <strong>Training</strong> โครงการ <strong>%s</strong> เป็นวันที่ <strong>%s</strong> เวลา <strong>%s น.</strong></p>`,PMName, PjName, MeetingDate, t.MeetingTime)
	htmlBody += fmt.Sprintf(`<p>ที่ <strong>%s</strong></p>`,selectMeetingRoom)


	
	if t.IsOnline == true {
		plainBody += fmt.Sprintf("หรือ ทาง MS Team ตาม Link นี้ %s\n",t.MSLink)
		htmlBody += fmt.Sprintf(`<p>หรือ ทาง MS Team ตาม Link นี้ <a href="%s">%s</a></p>`, t.MSLink, t.MSLink)
	}

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
	
}
func (ot *OperationTrainingServive)AlertEditQuestionTraining(pthlt []pilotquestionnaire.FilePath,member *[]members.MembersOfPj,PMName string,PjName string) error  {

	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง เปลี่ยนแปลงชุดในการ Training โครงการ"

	toList,err := ot.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	selectFilePath,err := ot.PQService.SelectFilePath(pthlt)
	if err != nil {
		return err
	}

	plainBody := fmt.Sprintf("%s แจ้งให้ทราบว่าโครงการ %s มีการเปลี่ยนแปลงชุดที่จะใช้ในการ Training\n",PMName,PjName)
	htmlBody := fmt.Sprintf(`<p>%s แจ้งให้ทราบว่าโครงการ <strong>%s</strong> มีการเปลี่ยนแปลงชุดที่จะใช้ในการ <strong>Training</strong></p>`,PMName, PjName)

	if selectFilePath != "" {
		htmlBody += fmt.Sprintf(`<p>โดย Version ที่ใช้ล่าสุด คือ:<br>%s</p><p>ทีมงานโปรดรับทราบและตรวจสอบชุด Training ที่จะใช้ใหม่นี้ด้วย</p>`, selectFilePath)
		plainBody += "โดย Version ที่ใช้ล่าสุด คือ: "+ "\n" + selectFilePath + "\n" + "ทีมงานโปรดรับทราบและตรวจสอบชุด Training ที่จะใช้ใหม่นี้ด้วย\n"
	}

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}

func (ot *OperationTrainingServive)AlertCancelTraining(member *[]members.MembersOfPj,t *Training, PMName string,PjName string) error{
	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง ยกเลิกการ Training โครงการ"

	toList,err := ot.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	plainBody := fmt.Sprintf(`%s ได้ยกเลิกการ Training โครงการ %s โดยจะแจ้งให้ทราบอีกครั้งเมื่อทำการนัดหมายใหม่`,PMName, PjName)
	htmlBody := fmt.Sprintf(`<P>%s ได้ยกเลิกการ <strong>Training</strong> โครงการ <strong>%s</strong> โดยจะแจ้งให้ทราบอีกครั้งเมื่อทำการนัดหมายใหม่</P>`,PMName, PjName)

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}

