package meeting

import (
	"log"
	"miami-back-end/db"
	"errors"
	"database/sql"

)
type MeetingRepository struct{

}

func NewMeetingSRepository () *MeetingRepository{
	return &MeetingRepository{}
}

func (mr *MeetingRepository)SendMeetingInvite(mti *Meeting) error  {

	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO meeting (
					project_id, meeting_date, meeting_time, is_online, ms_link, is_onsite,
					is_room1, is_room2, is_room3, is_room4, is_other, other, note, created_by, created_at
				) VALUES (
					:project_id, :meeting_date, :meeting_time, :is_online, :ms_link, :is_onsite,
					:is_room1, :is_room2, :is_room3, :is_room4, :is_other, :other, :note, :created_by, :created_at
				)`

	_, err = tx.NamedExec(query, mti)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	return tx.Commit()
}
func (mr *MeetingRepository)SendEditMeetingInvite(mti *Meeting) error  {

	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `UPDATE meeting 
				SET  
					project_id = :project_id, meeting_date = :meeting_date, meeting_time = :meeting_time, is_online = :is_online, ms_link = :ms_link
					,is_onsite = :is_onsite,is_room1 = :is_room1, is_room2 = :is_room2, is_room3 = :is_room3, is_room4 = :is_room4, is_other = :is_other, other = :other, note = :note
					,updated_by = :updated_by, updated_at = :updated_at,is_cancel = 0
				WHERE 
					project_id = :project_id;`
	
	
	_, err = tx.NamedExec(query, mti)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	return tx.Commit()
}

func (mr *MeetingRepository)GetMeetingInfo(PjID string) (*MeetingNullCase,error){

	query := `SELECT 
					project_id, meeting_date, meeting_time, is_online, ms_link, is_onsite,
					is_room1, is_room2, is_room3, is_room4, is_other, other, note,is_bypass,is_cancel
			FROM 
				meeting
			WHERE 
				project_id = ?`

	var meeting MeetingNullCase
	err := db.DB.Get(&meeting,query,PjID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ไม่พบข้อมูลในฐานข้อมูล
			return nil, nil // หรือจะ return error ใหม่ก็ได้ เช่น errors.New("meeting not found")
		}
		// error อื่น ๆ เช่น DB ล่ม หรือ query ผิด
		return nil, err
	}
	return &meeting,err
}

func (mr *MeetingRepository)ByPassMeeting(mti *Meeting) error{
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO meeting 
						(project_id,note,is_bypass,created_by,created_at) 
					VALUES 
						(:project_id,:note,:is_bypass,:created_by,:created_at);`

	_,err = tx.NamedExec(query,mti)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}
	return tx.Commit()
}

func (mr *MeetingRepository)CancelMeeting(mti *Meeting) error{
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE meeting 
				SET 
					note = :note,
					is_cancel = 1,
					updated_by = :updated_by,
					updated_at = :updated_at
				WHERE 
					project_id = :project_id;`

	_,err = tx.NamedExec(query,mti)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	return tx.Commit()
}
