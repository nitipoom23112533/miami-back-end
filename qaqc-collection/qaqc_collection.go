package qaqccollection

import(
	"miami-back-end/data-collection"
	"miami-back-end/members"
	"context"
	"fmt"
	"time"
	"github.com/mailgun/mailgun-go/v4"
	"miami-back-end/mg"
)

type QaqcCollectionService struct{
	QaqcCollectionRepository *QaqcCollectionRepository
	MemberService *members.MemberService
}

func NewQaqcCollectionService() *QaqcCollectionService{
	return &QaqcCollectionService{}
}

func (qs *QaqcCollectionService)StartQaqcCollection(dct *datacollection.DataCollection) error{
	return qs.QaqcCollectionRepository.StartQaqcCollection(dct)
}

func (qs *QaqcCollectionService)CompletedQaqcCollection(dct *datacollection.DataCollection) error{
	return qs.QaqcCollectionRepository.CompletedQaqcCollection(dct)
}

func (qs *QaqcCollectionService)GetQaqcCollection(PjID string) (datacollection.DataCollection,error){
	return qs.QaqcCollectionRepository.GetQaqcCollection(PjID)
}

func (qs *QaqcCollectionService)GetQcCollection(Pj string,) ([]datacollection.SsAndFsResponses,error){

	statuses := []string{"ยังไม่ได้ตรวจ", "รอ QC ตรวจซ้ำ", "คืนซ่อม", "รอ QCM"}
	return qs.QaqcCollectionRepository.GetQcCollection(Pj,statuses)
}

func (qs *QaqcCollectionService)GetQaCollection(Pj string,) ([]datacollection.SsAndFsResponses,error){

	return qs.QaqcCollectionRepository.GetQaCollection(Pj)
}
func (qs *QaqcCollectionService)AlertStartQaqccollection(member *[]members.MembersOfPj,dct *datacollection.DataCollection, PMName string,QaqcName []string,PjName string) error{
	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง เริ่มงาน QA/QC ข้อมูลโครงการ"

	toList,err := qs.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	var dateParsed time.Time
	if dct.StartDate.Valid {
		dateParsed, _ = time.Parse("2006-01-02", dct.StartDate.String)
	} 
	StartDate := dateParsed.Format("02/01/2006")


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
	
	plainBody := fmt.Sprintf(`%s แจ้งให้ทีมงานทราบว่าได้เริ่ม QA/QC ข้อมูลโครงการ %s ในวันที่ %s`,name, PjName, StartDate)
	htmlBody := fmt.Sprintf(`<P>%s แจ้งให้ทีมงานทราบว่าได้เริ่ม QA/QC ข้อมูลโครงการ <strong>%s</strong> ในวันที่ <strong>%s</strong></P>`,name, PjName, StartDate)

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}

func (qs *QaqcCollectionService)AlertCompletedQaqccollection(member *[]members.MembersOfPj,dct *datacollection.DataCollection, PMName string,QaqcName []string,PjName string) error{
	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง งาน QA/QC ข้อมูลโครงการ เสร็จสิ้น"

	toList,err := qs.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	var dateParsed time.Time
	if dct.CompletedDate.Valid {
		dateParsed, _ = time.Parse("2006-01-02", dct.CompletedDate.String)
	} 
	CompletedDate := dateParsed.Format("02/01/2006")

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
	
	plainBody := fmt.Sprintf(`%s แจ้งให้ทีมงานทราบว่างาน QA/QC ข้อมูลโครงการ  %s ได้เสร็จสิ้นในวันที่ %s`,name, PjName, CompletedDate)
	htmlBody := fmt.Sprintf(`<P>%s แจ้งให้ทีมงานทราบว่างาน QA/QC ข้อมูลโครงการ <strong>%s</strong> ได้เสร็จสิ้นในวันที่ <strong>%s</strong></P>`,name, PjName, CompletedDate)

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}