package questionnairesignoff

import (
	"context"
	"fmt"
	"miami-back-end/members"
	"miami-back-end/mg"
	"miami-back-end/pilot-questionnaire"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"gopkg.in/guregu/null.v4"

)

type QuestionnaireSignOffService struct{
	QuestionnaireSignOffRepository *QuestionnaireSignOffRepository
	PQService *pilotquestionnaire.PilotQuestionnaireService
	MemberService *members.MemberService
}
func NewQuestionnaireSignOffService() *QuestionnaireSignOffService{
	return &QuestionnaireSignOffService{}
}

type DetailOnChangeSignOff struct {
	ProjectID    	string `db:"project_id" json:"project_id"`
	DetailOnChange 	null.String `db:"detail_on_change" json:"detail_on_change"`
	Stage 			string `db:"stage" json:"stage"`
}

func (se *QuestionnaireSignOffService)SelectFileToSignOff(pth *pilotquestionnaire.Path) error  {
	return se.QuestionnaireSignOffRepository.SelectFileToSignOff(pth)
}

func (se *QuestionnaireSignOffService)InsertQuestionnaireSignOff(pth *pilotquestionnaire.Path,fileNumber int) error  {
	return se.QuestionnaireSignOffRepository.InsertQuestionnaireSignOff(pth,fileNumber)
}


func (se *QuestionnaireSignOffService)GetAllPilotQuestionnairePath(PjID string) ([]*pilotquestionnaire.FilePath,error) {
	return se.QuestionnaireSignOffRepository.GetAllPilotQuestionnairePath(PjID)
}

func (se *QuestionnaireSignOffService)GetSignOffDetailOnChange(project_id string,stage string) (DetailOnChangeSignOff,error){
	return se.QuestionnaireSignOffRepository.GetSignOffDetailOnChange(project_id,stage)
}

func (se *QuestionnaireSignOffService)AlertQuestionnaireSignOff(pth []pilotquestionnaire.FilePath,member *[]members.MembersOfPj,PMName string,PjName string) error  {

	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง การ Sign Off แบบสอบถามโครงการ"

	toList,err := se.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	selectFilePath,err := se.PQService.SelectFilePath(pth)
	if err != nil {
		return err
	}

	plainBody := fmt.Sprintf("%s แจ้งให้ทราบว่าแบบสอบถามโครงการ %s ได้รับการ Sign Off เรียบร้อยแล้ว\n",PMName,PjName)
	htmlBody := fmt.Sprintf(`<p>%s แจ้งให้ทราบว่าแบบสอบถามโครงการ <strong>%s</strong> ได้รับการ <strong>Sign Off</strong> เรียบร้อยแล้ว</p>`,PMName, PjName)

	if selectFilePath != "" {
		htmlBody += fmt.Sprintf(`<p>โดย Version ที่ได้รับการ Sign Off คือ:<br>%s</p>`, selectFilePath)
		plainBody += "โดย Version ที่ได้รับการ Sign Off คือ: "+ "\n" + selectFilePath + "\n"
	}

	
	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}

func (se *QuestionnaireSignOffService)AlertEditQuestionnaireSignOff(pth []pilotquestionnaire.FilePath,member *[]members.MembersOfPj,PMName string,PjName string) error  {

	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง เปลี่ยนแปลงแบบสอบถามที่ Sign Off โครงการ"

	toList,err := se.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	selectFilePath,err := se.PQService.SelectFilePath(pth)
	if err != nil {
		return err
	}

	plainBody := fmt.Sprintf("%s แจ้งให้ทราบว่าแบบสอบถามโครงการ %s มีการเปลี่ยนแปลงชุดที่ Sign Off\n",PMName,PjName)
	htmlBody := fmt.Sprintf(`<p>%s แจ้งให้ทราบว่าแบบสอบถามโครงการ <strong>%s</strong> มีการเปลี่ยนแปลงชุดที่ <strong>Sign Off</strong></p>`,PMName, PjName)

	if selectFilePath != "" {
		htmlBody += fmt.Sprintf(`<p>โดย Version ที่ใช้ล่าสุด คือ:<br>%s</p><p>ทีมงานโปรดรับทราบและตรวจสอบชุด Sign Off ใหม่นี้ด้วย</p>`, selectFilePath)
		plainBody += "โดย Version ที่ใช้ล่าสุด คือ: "+ "\n" + selectFilePath + "\n" + "ทีมงานโปรดรับทราบและตรวจสอบชุด Sign Off ใหม่นี้ด้วย\n"
	}

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}

