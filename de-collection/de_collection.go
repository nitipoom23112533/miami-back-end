package decollection

import(
	"miami-back-end/data-collection"
	"miami-back-end/members"
	"context"
	"fmt"
	"time"
	"github.com/mailgun/mailgun-go/v4"
	"miami-back-end/mg"

)

type DeCollectionService struct{
	DeCollectionRepo *DeCollectionRepo
	MemberService *members.MemberService

}
func NewDeCollectionService() *DeCollectionService{
	return &DeCollectionService{}
}

func (dc *DeCollectionService)StartDeCollection(dct *datacollection.DataCollection) error{
	return dc.DeCollectionRepo.StartDeCollection(dct)
}

func (dc *DeCollectionService)CompletedDeCollection(dct *datacollection.DataCollection) error{
	return dc.DeCollectionRepo.CompletedDeCollection(dct)
}

func (dc *DeCollectionService)GetDeCollectionInfo(PjID string) (datacollection.DataCollection,error){
	return dc.DeCollectionRepo.GetDeCollectionInfo(PjID)
}

func (dc *DeCollectionService)GetDeCollection(PjID string) ([]datacollection.SsAndFsResponses,error){
	return dc.DeCollectionRepo.GetDeCollection(PjID)
}

func (dc *DeCollectionService)AlertStartDecollection(member *[]members.MembersOfPj,dct *datacollection.DataCollection, PMName string,DeName []string,PjName string) error{
	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง เริ่มงาน Data Entry โครงการ"

	toList,err := dc.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	var dateParsed time.Time
	if dct.StartDate.Valid {
		dateParsed, _ = time.Parse("2006-01-02", dct.StartDate.String)
	} 
	StartDate := dateParsed.Format("02/01/2006")

	var dename string 
	for i, x := range DeName {
		if i != len(DeName)-1 {
			dename += x + " และ "

		}else {
			dename += x
		}
	}

	var name string
	if len(DeName) > 0 {
		name = dename
	}else{
		name = PMName
	}
	
	plainBody := fmt.Sprintf(`%s แจ้งให้ทีมงานทราบว่าได้เริ่มงาน Data Entry โครงการ %s ในวันที่ %s`,name, PjName, StartDate)
	htmlBody := fmt.Sprintf(`<P>%s แจ้งให้ทีมงานทราบว่าได้เริ่มงาน Data Entry โครงการ <strong>%s</strong> ในวันที่ <strong>%s</strong></P>`,name, PjName, StartDate)

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}

func (dc *DeCollectionService)AlertCompletedDecollection(member *[]members.MembersOfPj,dct *datacollection.DataCollection, PMName string,DeName []string,PjName string) error{
	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง งาน Data Entry โครงการ เสร็จสิ้น"

	toList,err := dc.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	var dateParsed time.Time
	if dct.CompletedDate.Valid {
		dateParsed, _ = time.Parse("2006-01-02", dct.CompletedDate.String)
	} 
	CompletedDate := dateParsed.Format("02/01/2006")

	var dename string 
	for i, x := range DeName {
		if i != len(DeName)-1 {
			dename += x + " และ "

		}else {
			dename += x
		}
	}

	var name string
	if len(DeName) > 0 {
		name = dename
	}else{
		name = PMName
	}
	
	plainBody := fmt.Sprintf(`%s แจ้งให้ทีมงานทราบว่างาน Data Entry โครงการ %s ได้เสร็จสิ้นในวันที่ %s`,name, PjName, CompletedDate)
	htmlBody := fmt.Sprintf(`<P>%s แจ้งให้ทีมงานทราบว่างาน Data Entry โครงการ <strong>%s</strong> ได้เสร็จสิ้นในวันที่ <strong>%s</strong></P>`,name, PjName, CompletedDate)

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}

