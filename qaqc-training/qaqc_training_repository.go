package qaqctraining
import(
	"miami-back-end/db"
	"log"
	"github.com/jmoiron/sqlx"
	"miami-back-end/pilot-questionnaire"
	"database/sql"
	"errors"
	"gopkg.in/guregu/null.v4"
)

type QAQCRepository struct {

}

func NewQAQCRepository() *QAQCRepository {
	return &QAQCRepository{}
}

func (qq *QAQCRepository)SendCreateQAQCTraining(q *QAQCTraining) error  {

	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO qaqc_training (
					project_id, meeting_date, meeting_time, is_online, ms_link, is_onsite,
					is_room1, is_room2, is_room3, is_room4, is_other, other, note, created_by, created_at
				) VALUES (
					:project_id, :meeting_date, :meeting_time, :is_online, :ms_link, :is_onsite,
					:is_room1, :is_room2, :is_room3, :is_room4, :is_other, :other, :note, :created_by, :created_at
				)`

	_, err = tx.NamedExec(query, q)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	detailOnChangeQuery := `INSERT INTO detail_on_change (project_id,stage) VALUES (?,'8')`
	_,err = tx.Exec(detailOnChangeQuery,q.ProjectID)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	return tx.Commit()
}

func (qq *QAQCRepository)InsertQAQCTrainingPath(p *pilotquestionnaire.Path) error  {

	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	InsertQuery := `INSERT INTO qaqc_questionnaire_path 
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

func (qq *QAQCRepository)GetQAQCTrainingInfo(PjID string) (*QAQCTrainingNullCase,error){

	query := `SELECT 
					project_id, meeting_date, meeting_time, is_online, ms_link, is_onsite,
					is_room1, is_room2, is_room3, is_room4, is_other, other, note,is_bypass,is_cancel
			FROM 
				qaqc_training
			WHERE 
				project_id = ?`
	

	var q QAQCTrainingNullCase
	err := db.DB.Get(&q,query,PjID)

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
						qaqc_questionnaire_path
					WHERE 
						project_id = ?
					ORDER BY number DESC;`

	var Pth []pilotquestionnaire.FilePath
	err = db.DB.Select(&Pth,queryPath,PjID)
	if err != nil {
		// ถ้าไม่มี path ก็ให้เป็น [] ไม่ error เพราะอาจไม่มีไฟล์ได้
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		Pth = []pilotquestionnaire.FilePath{}
	}
	q.FilePath = Pth

	queryDetailOnChange := `SELECT 
								detail_on_change
							FROM 
								detail_on_change
							WHERE 
								project_id = ? and stage = '8';`
	var detailOnChange null.String
	err = db.DB.Get(&detailOnChange,queryDetailOnChange,PjID)
	if err != nil {
		// ถ้าไม่เจอก็ไม่เป็นไร ให้เป็น null.String{}
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	q.DetailOnChange = detailOnChange

	return &q,nil
}

func (qq *QAQCRepository)GetLatestQAQCTrainingPath(PjID string) ([]pilotquestionnaire.FilePath,error)  {

	query := `SELECT 
					path,file_name,number,is_new,is_sign
				FROM 
					qaqc_questionnaire_path
				WHERE 
					project_id = ? and is_new = 1`
	
	var Pth []pilotquestionnaire.FilePath
	err := db.DB.Select(&Pth,query,PjID)
	if err != nil {
		return nil, err
	}
	return Pth,err
}

func (qq *QAQCRepository)GetAllQAQCTrainingPath(PjID string) ([]*pilotquestionnaire.FilePath,error)  {

	query := `SELECT 
					path,file_name,number,is_new,is_sign,is_training
				FROM 
					qaqc_questionnaire_path
				WHERE 
					project_id = ?`
	
	var Pth []*pilotquestionnaire.FilePath
	err := db.DB.Select(&Pth,query,PjID)
	if err != nil {
		return nil, err
	}
	return Pth,err
}

func (qq *QAQCRepository)QAQCEditdetail(q *QAQCTraining) error {
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `UPDATE qaqc_training
				SET  
					project_id = :project_id, meeting_date = :meeting_date, meeting_time = :meeting_time, is_online = :is_online, ms_link = :ms_link
					,is_onsite = :is_onsite,is_room1 = :is_room1, is_room2 = :is_room2, is_room3 = :is_room3, is_room4 = :is_room4, is_other = :is_other, other = :other, note = :note
					,updated_by = :updated_by, updated_at = :updated_at,is_cancel = 0
				WHERE 
					project_id = :project_id;`
	
	
	_, err = tx.NamedExec(query, q)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	return tx.Commit()
}

func (qq *QAQCRepository)UpdatePathQAQCTraining(p *pilotquestionnaire.Path,fileNumber int) error  {

	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	UpdateQuery := `UPDATE qaqc_questionnaire_path 
				SET 
					is_new = 0
				WHERE 
					project_id = :project_id;`

	_,err = tx.NamedExec(UpdateQuery,p)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	InsertQuery := `INSERT INTO qaqc_questionnaire_path 
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

func (qq *QAQCRepository)CancelQAQCTraining(q *QAQCTraining) error{
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE qaqc_training
				SET 
					note = :note,
					is_cancel = 1,
					updated_by = :updated_by,
					updated_at = :updated_at
				WHERE 
					project_id = :project_id;`

	_,err = tx.NamedExec(query,q)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	return tx.Commit()
}

func (qq *QAQCRepository)ByPassQAQCTraining(q *QAQCTraining) error{
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO qaqc_training 
						(project_id,note,is_bypass,created_by,created_at) 
					VALUES 
						(:project_id,:note,:is_bypass,:created_by,:created_at);`

	_,err = tx.NamedExec(query,q)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}
	return tx.Commit()
}