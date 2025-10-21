package pilotquestionnaire
import (
	"log"
	"miami-back-end/db"
	"errors"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"

)

type PilotQuestionnaireRepository struct{

}

func NewPilotQuestionnaireRepository() *PilotQuestionnaireRepository{

	return &PilotQuestionnaireRepository{}
}

func (pr *PilotQuestionnaireRepository)CreatePilotQuestionnaire(pq *PilotQuestionnaire) error  {

	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO pilot_questionnaire (
					project_id, meeting_date, meeting_time, is_online, ms_link, is_onsite,
					is_room1, is_room2, is_room3, is_room4, is_other, other, note, created_by, created_at
				) VALUES (
					:project_id, :meeting_date, :meeting_time, :is_online, :ms_link, :is_onsite,
					:is_room1, :is_room2, :is_room3, :is_room4, :is_other, :other, :note, :created_by, :created_at
				)`

	_, err = tx.NamedExec(query, pq)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	detailOnChangeQuery := `INSERT INTO detail_on_change (project_id,stage) VALUES (?,'3')`
	_,err = tx.Exec(detailOnChangeQuery,pq.ProjectID)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	return tx.Commit()
}

func (pr *PilotQuestionnaireRepository)ByPassPilotQuestionnaire(pq *PilotQuestionnaire) error{
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO pilot_questionnaire 
						(project_id,note,is_bypass,created_by,created_at) 
					VALUES 
						(:project_id,:note,:is_bypass,:created_by,:created_at);`

	_,err = tx.NamedExec(query,pq)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}
	return tx.Commit()
}

func (pr *PilotQuestionnaireRepository)GetPilotQuestionnaireInfo(PjID string) (*PilotQuestionnaireNullCase,error){

	query := `SELECT 
					project_id, meeting_date, meeting_time, is_online, ms_link, is_onsite,
					is_room1, is_room2, is_room3, is_room4, is_other, other, note,is_bypass,is_cancel
			FROM 
				pilot_questionnaire
			WHERE 
				project_id = ?`
	

	var PQ PilotQuestionnaireNullCase
	err := db.DB.Get(&PQ,query,PjID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ไม่พบข้อมูลในฐานข้อมูล
			return nil, nil // หรือจะ return error ใหม่ก็ได้ เช่น errors.New("meeting not found")
		}
		// error อื่น ๆ เช่น DB ล่ม หรือ query ผิด
		return nil, err
	}

	queryPath := `SELECT 
						path,file_name,number,is_new,is_sign,is_training
					FROM 
						questionnaire_path
					WHERE 
						project_id = ?
					ORDER BY number DESC;`

	var Pth []FilePath
	err = db.DB.Select(&Pth,queryPath,PjID)
	if err != nil {
		// ถ้าไม่มี path ก็ให้เป็น [] ไม่ error เพราะอาจไม่มีไฟล์ได้
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		Pth = []FilePath{}
	}
	PQ.FilePath = Pth

	queryDetailOnChange := `SELECT 
								detail_on_change
							FROM 
								detail_on_change
							WHERE 
								project_id = ? and stage = '3';`
	var detailOnChange null.String
	err = db.DB.Get(&detailOnChange,queryDetailOnChange,PjID)
	if err != nil {
		// ถ้าไม่เจอก็ไม่เป็นไร ให้เป็น null.String{}
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	PQ.DetailOnChange = detailOnChange

	return &PQ,nil
}

func (pr *PilotQuestionnaireRepository)Editdetail(pq *PilotQuestionnaire) error {
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `UPDATE pilot_questionnaire 
				SET  
					project_id = :project_id, meeting_date = :meeting_date, meeting_time = :meeting_time, is_online = :is_online, ms_link = :ms_link
					,is_onsite = :is_onsite,is_room1 = :is_room1, is_room2 = :is_room2, is_room3 = :is_room3, is_room4 = :is_room4, is_other = :is_other, other = :other, note = :note
					,updated_by = :updated_by, updated_at = :updated_at,is_cancel = 0
				WHERE 
					project_id = :project_id;`
	
	
	_, err = tx.NamedExec(query, pq)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	return tx.Commit()
}
func (pr *PilotQuestionnaireRepository)CancelPilotQuestionnaire(pq *PilotQuestionnaire) error{
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE pilot_questionnaire 
				SET 
					note = :note,
					is_cancel = 1,
					updated_by = :updated_by,
					updated_at = :updated_at
				WHERE 
					project_id = :project_id;`

	_,err = tx.NamedExec(query,pq)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	return tx.Commit()
}
func (pr *PilotQuestionnaireRepository)GetAllPilotQuestionnairePath(PjID string) ([]*FilePath,error)  {

	query := `SELECT 
					path,file_name,number,is_new,is_sign,is_training
				FROM 
					questionnaire_path
				WHERE 
					project_id = ?`
	
	var Pth []*FilePath
	err := db.DB.Select(&Pth,query,PjID)
	if err != nil {
		return nil, err
	}
	return Pth,err
}
func (pr *PilotQuestionnaireRepository)GetLatestPilotQuestionnairePath(PjID string) ([]FilePath,error)  {

	query := `SELECT 
					path,file_name,number,is_new,is_sign
				FROM 
					questionnaire_path
				WHERE 
					project_id = ? and is_new = 1`
	
	var Pth []FilePath
	err := db.DB.Select(&Pth,query,PjID)
	if err != nil {
		return nil, err
	}
	return Pth,err
}
func (pr *PilotQuestionnaireRepository)InsertPathPilotQuestionnaire(p *Path) error  {

	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	InsertQuery := `INSERT INTO questionnaire_path 
						(project_id,path,number,file_name,is_new) 
					VALUES 
						(?,?,?,?,1);`
	
	for i, ps := range p.FilePath{
		query, args, err := sqlx.In(InsertQuery, p.ProjectID, ps.Path, i+1,ps.FileName)
		if err != nil {
			return err
		}

		query = tx.Rebind(query) // ทำให้ query ใช้ placeholder ตาม DB (เช่น ?, $1, etc.)

		_, err = tx.Exec(query, args...)
		if err != nil {
			return err
		}
	}

	
	return tx.Commit()
}

func (pr *PilotQuestionnaireRepository)UpdatePathPilotQuestionnaire(p *Path,fileNumber int) error  {

	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	UpdateQuery := `UPDATE questionnaire_path 
				SET 
					is_new = 0
				WHERE 
					project_id = :project_id;`

	_,err = tx.NamedExec(UpdateQuery,p)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	InsertQuery := `INSERT INTO questionnaire_path 
						(project_id,path,number,file_name,is_new) 
					VALUES 
						(?,?,?,?,1);`
	
	n := fileNumber
	for _, ps := range p.FilePath{
		n += 1
		query, args, err := sqlx.In(InsertQuery, p.ProjectID, ps.Path, n,ps.FileName)
		if err != nil {
			return err
		}

		query = tx.Rebind(query) // ทำให้ query ใช้ placeholder ตาม DB (เช่น ?, $1, etc.)

		_, err = tx.Exec(query, args...)
		if err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

func (pr *PilotQuestionnaireRepository)DetailOnChange(pq *PilotQuestionnaire,stage string) error {
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE detail_on_change 
				SET 
					detail_on_change = ?
				WHERE 
					project_id = ? and stage = ?;`

	_,err = tx.Exec(query,pq.DetailOnChange,pq.ProjectID,stage)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	return tx.Commit()
}