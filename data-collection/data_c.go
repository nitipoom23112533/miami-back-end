package datacollection

import (
	"time"

	"gopkg.in/guregu/null.v4"
	"miami-back-end/members"
	"miami-back-end/mg"
	"fmt"
	"context"
	"github.com/mailgun/mailgun-go/v4"
)

type DataCollectionService struct{
	DataCollectionRepository *DataCollectionRepository
	MemberService *members.MemberService


}

func NewDataCollectionService() *DataCollectionService{
	return &DataCollectionService{}
}

type DataCollection struct{
	ProjectID 			string `db:"project_id" json:"project_id"`
	IsStart 			bool `db:"is_start" json:"is_start"`
	Quota 				null.Int `db:"quota" json:"quota"`
	Day 				null.Int `db:"day" json:"day"`
	StartDate 			null.String `db:"start_date" json:"start_date"`
	IsCompleted 		bool `db:"is_completed" json:"is_completed"`
	Completed 			null.Int `db:"completed" json:"completed"`
	CompletedDate 		null.String `db:"completed_date" json:"completed_date"`
	CreatedBy 			string `db:"created_by" json:"created_by"`
	CreatedAt 			time.Time `db:"created_at" json:"created_at"`
	UpdatedBy 			null.String `db:"updated_by" json:"updated_by"`
	UpdatedAt 			null.Time `db:"updated_at" json:"updated_at"`

}

type DashboardLogs struct{
	ProjectID 			string `db:"project_id" json:"project_id"`
	Mn 					null.Int `db:"mn" json:"mn"`
	Qc 					null.Int `db:"qc" json:"qc"`
	Qa 					null.Int `db:"qa" json:"qa"`
	Fw 					null.Int `db:"fw" json:"fw"`
	De 					null.Int `db:"de" json:"de"`
	Da 					null.Int `db:"da" json:"da"`
	Doc 				null.Int `db:"doc" json:"doc"`
	DocReject 			null.Int `db:"doc_reject" json:"doc_reject"`
}
type SsAndFsResponses struct {
	Id 					int64 `db:"id" json:"id"`
	ProjectID 			string `db:"project_id" json:"project_id"`
	Status 				string `db:"status" json:"status"`
}
	

func (dc *DataCollectionService)StartDataCollection(dct *DataCollection) error{
	return dc.DataCollectionRepository.StartDataCollection(dct)

}

func (dc *DataCollectionService)CompletedDataCollection(dct *DataCollection) error{
	return dc.DataCollectionRepository.CompletedDataCollection(dct)
}


func (dc *DataCollectionService)GetDataCollection(PjID string) (*DataCollection,error){
	return dc.DataCollectionRepository.GetDataCollection(PjID)
}

func (dc *DataCollectionService)GetDashBoardLogs(PjID string) (*DashboardLogs,error){
	return dc.DataCollectionRepository.GetDashBoardLogs(PjID)

}

func (dc *DataCollectionService)GetSsResponses(Pj string) (*[]SsAndFsResponses,error){
	return dc.DataCollectionRepository.GetSsResponses(Pj)
}

func (dc *DataCollectionService)GetFsResponses(Pj string) (*[]SsAndFsResponses,error){
	return  dc.DataCollectionRepository.GetFsResponses(Pj)
}

func (dc *DataCollectionService)AlertStartDatacollection(member *[]members.MembersOfPj,dct *DataCollection, PMName string,FWName []string,PjName string) error{
	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง Data Collection เริ่มเก็บข้อมูลโครงการ"

	toList,err := dc.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	var dateParsed time.Time
	if dct.StartDate.Valid {
		dateParsed, _ = time.Parse("2006-01-02", dct.StartDate.String)
	} 
	StartDate := dateParsed.Format("02/01/2006")

	var fwname string 

	for i, x := range FWName {
		if i != len(FWName)-1 {
			fwname += x + " และ "

		}else {
			fwname += x
		}
	}

	var name string
	if len(FWName) > 0 {
		name = fwname
	}else{
		name = PMName
	}
	
	plainBody := fmt.Sprintf(`%s แจ้งให้ทีมงานทราบว่าการเริ่มเก็บข้อมูลโครงการ %s เริ่มวันที่ %s`,name, PjName, StartDate)
	htmlBody := fmt.Sprintf(`<P>%s แจ้งให้ทีมงานทราบว่าการเริ่มเก็บข้อมูลโครงการ <strong>%s</strong> เริ่มวันที่ <strong>%s</strong></P>`,name, PjName, StartDate)

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}

func (dc *DataCollectionService)AlertCompletedDatacollection(member *[]members.MembersOfPj,dct *DataCollection, PMName string,FWName []string,PjName string) error{
	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง Data Collection เสร็จสิ้นเก็บข้อมูลโครงการ"

	toList,err := dc.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	var dateParsed time.Time
	if dct.CompletedDate.Valid {
		dateParsed, _ = time.Parse("2006-01-02", dct.CompletedDate.String)
	} 
	CompletedDate := dateParsed.Format("02/01/2006")

	var fwname string 

	for i, x := range FWName {
		if i != len(FWName)-1 {
			fwname += x + " และ "

		}else {
			fwname += x
		}
	}

	var name string
	if len(FWName) > 0 {
		name = fwname
	}else{
		name = PMName
	}
	
	plainBody := fmt.Sprintf(`%s แจ้งให้ทีมงานทราบว่าการเก็บข้อมูลโครงการ %s ได้เสร็จสิ้นในวันที่ %s`,name, PjName, CompletedDate)
	htmlBody := fmt.Sprintf(`<P>%s แจ้งให้ทีมงานทราบว่าการเก็บข้อมูลโครงการ <strong>%s</strong> ได้เสร็จสิ้นในวันที่ <strong>%s</strong></P>`,name, PjName, CompletedDate)

	message := mailgun.NewMessage(sender, subject, plainBody, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
}