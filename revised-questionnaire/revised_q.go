package revisedquestionnaire
import(
	"miami-back-end/pilot-questionnaire"
	"miami-back-end/members"
	"time"
	"fmt"
	"context"
	"github.com/mailgun/mailgun-go/v4"
	"miami-back-end/mg"
)

type RQService struct{
	RQRepository *RQRepository
	PQService *pilotquestionnaire.PilotQuestionnaireService
	MemberService *members.MemberService

}

func NewRQService() *RQService{
	return &RQService{}
}

type DetailOnChange struct {
	ProjectID    	string `db:"project_id" json:"project_id"`
	DetailOnChange 	string `db:"detail_on_change" json:"detail_on_change"`
	Stage 			string `db:"stage" json:"stage"`
}

func (rq *RQService)GetRevisedQuestionnaire(project_id string ) ([]pilotquestionnaire.FilePath,error){
	return rq.RQRepository.GetRevisedQuestionnaire(project_id)
}
func (rq *RQService)GetRevisedQuestionnaireDetailOnChange(project_id string,stage string) (DetailOnChange,error){
	return rq.RQRepository.GetRevisedQuestionnaireDetailOnChange(project_id,stage)
}

func (rq *RQService)RQRDetailOnChange(p *pilotquestionnaire.PilotQuestionnaire,stage string) error {
	return rq.RQRepository.RQRDetailOnChange(p,stage)
}

func (rq *RQService)InsertRevisedQuestionnaire(p *pilotquestionnaire.Path,fileNumber int) error  {
	return rq.RQRepository.InsertRevisedQuestionnaire(p,fileNumber)
}
func (rq *RQService)AlertRevisedQuestionnaire(pth []pilotquestionnaire.FilePath,member *[]members.MembersOfPj,PMName string,PjName string) error  {

	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง Revised Questionnaire โครงการ"

	toList,err := rq.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	selectFilePath,err := rq.PQService.SelectFilePath(pth)
	if err != nil {
		return err
	}

	plainBody := fmt.Sprintf("%s แจ้งให้ทราบว่าโครงการ %s มีการ Revised Questionnaire หลัง Training\n",PMName,PjName)
	htmlBody := fmt.Sprintf(`<p>%s แจ้งให้ทราบว่าโครงการ <strong>%s</strong> มีการ <strong>Revised Questionnaire</strong> หลัง Training</p>`,PMName, PjName)

	if selectFilePath != "" {
		htmlBody += fmt.Sprintf(`<p>โดย Version ที่ใช้ล่าสุด คือ:<br>%s</p><p>ทีมงานโปรดรับทราบและตรวจแบบสอบถามชุดใหม่นี้ด้วย</p>`, selectFilePath)
		plainBody += "โดย Version ที่ใช้ล่าสุด คือ: "+ "\n" + selectFilePath + "\n" + "ทีมงานโปรดรับทราบและตรวจแบบสอบถามชุดใหม่นี้ด้วย\n"
	}

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}