package meeting

import (
	"context"
	"fmt"
	"miami-back-end/members"
	"miami-back-end/mg"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"gopkg.in/guregu/null.v4"
	// "strings"
)



type MeetingService struct{
	MeetingRepository *MeetingRepository
	MemberService *members.MemberService

}

func NewMeetingService() *MeetingService {
	return &MeetingService{}
}


type Meeting struct {
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


}

type MeetingNullCase struct {
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

}


func (ms *MeetingService)SendMeetingInvite(mti *Meeting) error {

	return ms.MeetingRepository.SendMeetingInvite(mti)
}
func (ms *MeetingService)SendEditMeetingInvite(mti *Meeting) error {

	return ms.MeetingRepository.SendEditMeetingInvite(mti)
}
func (ms *MeetingService)GetMeetingInfo(PjID string) (*MeetingNullCase,error){

	return ms.MeetingRepository.GetMeetingInfo(PjID)
}
func (ms *MeetingService)ByPassMeeting(mti *Meeting) error {

	return ms.MeetingRepository.ByPassMeeting(mti)
}
func (ms *MeetingService)CancelMeeting(mti *Meeting) error {

	return ms.MeetingRepository.CancelMeeting(mti)
}

func (ms *MeetingService)AlertInviteMeeting(member *[]members.MembersOfPj,mif *Meeting, PMName string,PjName string) error{

	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง Kick off meeting โครงการ"

	toList,err := ms.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	dateParsed, _ := time.Parse("2006-01-02", mif.MeetingDate)
	MeetingDate := dateParsed.Format("02/01/2006")

	selectMeetingRoom,err := ms.SelectMeetingRoom(mif)
	if err != nil {
		return err
	}
	
	body := fmt.Sprintf("%s เชิญทุกท่านประชุม Kick off meeting โครงการ %s ในวันที่ %s เวลา %s น.\n",PMName,PjName,MeetingDate,mif.MeetingTime)
	body += fmt.Sprintf("ที่ %s\n",selectMeetingRoom)
	htmlBody := fmt.Sprintf(`<p>%s เชิญทุกท่านประชุม <strong>Kick off meeting</strong> โครงการ <strong>%s</strong> ในวันที่ <strong>%s</strong> เวลา <strong>%s น.</strong></p>`,PMName,PjName,MeetingDate,mif.MeetingTime)
	htmlBody += fmt.Sprintf(`<p>ที่ <strong>%s</strong></p>`,selectMeetingRoom)

	if mif.IsOnline == true {
		body += fmt.Sprintf("หรือ ทาง MS Team ตาม Link นี้ %s\n",mif.MSLink)
		htmlBody += fmt.Sprintf(`<p>หรือ ทาง MS Team ตาม Link นี้ <a href="%s">%s</a></p>`,mif.MSLink,mif.MSLink)
	}

	message := mailgun.NewMessage(sender, subject, body, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
	
}

func (ms *MeetingService)AlertEditInviteMeeting(member *[]members.MembersOfPj,mif *Meeting, PMName string,PjName string) error{

	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง เปลี่ยนแปลงการประชุม Kick off meeting โครงการ"

	toList,err := ms.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	dateParsed, _ := time.Parse("2006-01-02", mif.MeetingDate)
	MeetingDate := dateParsed.Format("02/01/2006")

	selectMeetingRoom,err := ms.SelectMeetingRoom(mif)
	if err != nil {
		return err
	}
	
	body := fmt.Sprintf("%s ได้เปลี่ยนแปลงการประชุม Kick off meeting โครงการ %s เป็นวันที่ %s เวลา %s น.\n",PMName,PjName,MeetingDate,mif.MeetingTime)
	body += fmt.Sprintf("ที่ %s\n",selectMeetingRoom)
	htmlBody := fmt.Sprintf(`<p>%s ได้เปลี่ยนแปลงการประชุม <strong>Kick off meeting</strong> โครงการ <strong>%s</strong> เป็นวันที่ <strong>%s</strong> เวลา <strong>%s น.</strong></p>`,PMName,PjName,MeetingDate,mif.MeetingTime)
	htmlBody += fmt.Sprintf(`<p>ที่ <strong>%s</strong></p>`,selectMeetingRoom)

	if mif.IsOnline == true {
		body += fmt.Sprintf("หรือ ทาง MS Team ตาม Link นี้ %s\n",mif.MSLink)
		htmlBody += fmt.Sprintf(`<p>หรือ ทาง MS Team ตาม Link นี้ <a href="%s">%s</a></p>`,mif.MSLink,mif.MSLink)
	}

	message := mailgun.NewMessage(sender, subject, body, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
	
}

func (ms *MeetingService)AlertCancelInviteMeeting(member *[]members.MembersOfPj,mif *Meeting, PMName string,PjName string) error{

	sender := "BRS Admin Official <noreply@brsth.com>"
	subject := "Miami เรื่อง ยกเลิกการประชุม Kick off meeting โครงการ"

	toList,err := ms.MemberService.EmailMemberList(member)
	if err != nil {
		return err
	}

	body := fmt.Sprintf("%s ได้ยกเลิกการประชุม Kick off meeting โครงการ %s โดยจะแจ้งให้ทราบอีกครั้งเมื่อทำการนัดหมายใหม่\n",PMName,PjName)
	htmlBody := fmt.Sprintf(`<p>%s ได้ยกเลิกการประชุม <strong>Kick off meeting</strong> โครงการ <strong>%s</strong> โดยจะแจ้งให้ทราบอีกครั้งเมื่อทำการนัดหมายใหม่</p>`,PMName,PjName)

	message := mailgun.NewMessage(sender, subject, body, toList...)
	message.SetHTML(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = mg.Client.Send(ctx, message)
	return err
	
}

func (ms *MeetingService)SelectMeetingRoom(mif *Meeting) (string,error){

	MeetingRoom := []string{"ชั้น 21 ห้องประชุมใหญ่","ชั้น 21 ห้องประชุมเล็ก","ชั้น 23 ห้อง Focus Group","ชั้น 23 ห้องพัก RD"}
	var selectMeetingRoomArray []string

	for i, x := range MeetingRoom {
		if mif.IsRoom1 == true && i == 0 {
			selectMeetingRoomArray = append(selectMeetingRoomArray, x)
		}else if mif.IsRoom2 == true && i == 1 {
			selectMeetingRoomArray = append(selectMeetingRoomArray, x)
		}else if mif.IsRoom3 == true && i == 2 {
			selectMeetingRoomArray = append(selectMeetingRoomArray, x)
		}else if mif.IsRoom4 == true && i == 3 {
			selectMeetingRoomArray = append(selectMeetingRoomArray, x)
		}
	}

	if mif.IsOther == true {
		selectMeetingRoomArray = append(selectMeetingRoomArray, mif.Other)
	}
	var selectMeetingRoom string
	
	for i, x := range selectMeetingRoomArray {
		if i != len(selectMeetingRoomArray)-1 {
			selectMeetingRoom += x + " และที่ "

		}else {
			selectMeetingRoom += x
		}
	}
	return  selectMeetingRoom,nil
}


