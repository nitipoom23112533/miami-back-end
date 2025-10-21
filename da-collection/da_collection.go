package dacollection

import(
	"miami-back-end/data-collection"
	"miami-back-end/members"
	"context"
	"fmt"
	"time"
	"github.com/mailgun/mailgun-go/v4"
	"miami-back-end/mg"

)

type DaCollectionService struct{
	DaCollectionRepo *DaCollectionRepo
	MemberService *members.MemberService


}
func NewDaCollectionService() *DaCollectionService{
	return &DaCollectionService{}
}

func (dc *DaCollectionService)StartDaCollection(dct *datacollection.DataCollection) error{
	return dc.DaCollectionRepo.StartDaCollection(dct)
}

func (dc *DaCollectionService)CompletedDaCollection(dct *datacollection.DataCollection) error{
	return dc.DaCollectionRepo.CompletedDaCollection(dct)

}

func (dc *DaCollectionService)GetDaCollection(PjID string) ([]datacollection.SsAndFsResponses,error){
	return dc.DaCollectionRepo.GetDaCollection(PjID)
}

func (dc *DaCollectionService)GetDaCollectionInfo(PjID string) (datacollection.DataCollection,error){
	return dc.DaCollectionRepo.GetDaCollectionInfo(PjID)
}

func (dc *DaCollectionService)AlertStartDacollection(member *[]members.MembersOfPj,dct *datacollection.DataCollection, PMName string,DaName []string,PjName string) error{
	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง เริ่มงาน Data Analysis โครงการ"

	toList,err := dc.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	var dateParsed time.Time
	if dct.StartDate.Valid {
		dateParsed, _ = time.Parse("2006-01-02", dct.StartDate.String)
	} 
	StartDate := dateParsed.Format("02/01/2006")

	var daname string 
	for i, x := range DaName {
		if i != len(DaName)-1 {
			daname += x + " และ "

		}else {
			daname += x
		}
	}

	var name string
	if len(DaName) > 0 {
		name = daname
	}else{
		name = PMName
	}
	
	plainBody := fmt.Sprintf(`%s แจ้งให้ทีมงานทราบว่าได้เริ่มงาน Data Analysis โครงการ %s ในวันที่ %s`,name, PjName, StartDate)
	htmlBody := fmt.Sprintf(`<P>%s แจ้งให้ทีมงานทราบว่าได้เริ่มงาน Data Analysis โครงการ <strong>%s</strong> ในวันที่ <strong>%s</strong></P>`,name, PjName, StartDate)

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}

func (dc *DaCollectionService)AlertCompletedDacollection(member *[]members.MembersOfPj,dct *datacollection.DataCollection, PMName string,DaName []string,PjName string) error{
	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง งาน Data Analysis โครงการ เสร็จสิ้น"

	toList,err := dc.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}
	var dateParsed time.Time
	if dct.CompletedDate.Valid {
		dateParsed, _ = time.Parse("2006-01-02", dct.CompletedDate.String)
	} 
	CompletedDate := dateParsed.Format("02/01/2006")

	var daname string 
	for i, x := range DaName {
		if i != len(DaName)-1 {
			daname += x + " และ "

		}else {
			daname += x
		}
	}

	var name string
	if len(DaName) > 0 {
		name = daname
	}else{
		name = PMName
	}
	
	plainBody := fmt.Sprintf(`%s แจ้งให้ทีมงานทราบว่างาน Data Analysis โครงการ %s ได้เสร็จสิ้นในวันที่ %s`,name, PjName, CompletedDate)
	htmlBody := fmt.Sprintf(`<P>%s แจ้งให้ทีมงานทราบว่างาน Data Analysis โครงการ <strong>%s</strong> ได้เสร็จสิ้นในวันที่ <strong>%s</strong></P>`,name, PjName, CompletedDate)

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}