package operationtraining
import(
	"miami-back-end/db"
	"log"
	"errors"
	"database/sql"
	"miami-back-end/pilot-questionnaire"
	"gopkg.in/guregu/null.v4"
)
type  OperationTrainingRepository struct{

}

func NewOperationTrainingRepository() *OperationTrainingRepository{
	return &OperationTrainingRepository{}
}

func (or *OperationTrainingRepository)SendTrainingInvite(t *Training) error  {

	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO operation_training (
					project_id, meeting_date, meeting_time, is_online, ms_link, is_onsite,
					is_room1, is_room2, is_room3, is_room4, is_other, other, note, created_by, created_at
				) VALUES (
					:project_id, :meeting_date, :meeting_time, :is_online, :ms_link, :is_onsite,
					:is_room1, :is_room2, :is_room3, :is_room4, :is_other, :other, :note, :created_by, :created_at
				)`

	_, err = tx.NamedExec(query, t)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	detailOnChangeQuery := `INSERT INTO detail_on_change (project_id,stage) VALUES (?,'5')`
	_,err = tx.Exec(detailOnChangeQuery,t.ProjectID)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	return tx.Commit()
}

func (or *OperationTrainingRepository)GetTrainingInfo(PjID string) (*TrainingNullCase,error){

	query := `SELECT 
					project_id, meeting_date, meeting_time, is_online, ms_link, is_onsite,
					is_room1, is_room2, is_room3, is_room4, is_other, other, note,is_bypass,is_cancel
			FROM 
				operation_training
			WHERE 
				project_id = ?`

	var PQ TrainingNullCase
	err := db.DB.Get(&PQ,query,PjID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ไม่พบข้อมูลในฐานข้อมูล
			return nil, nil // หรือจะ return error ใหม่ก็ได้ เช่น errors.New("meeting not found")
		}
		// error อื่น ๆ เช่น DB ล่ม หรือ query ผิด
		return nil, err
	}
	return &PQ,err
}

func (or *OperationTrainingRepository)GetAllPilotQuestionnairePathIsSign(PjID string) ([]*pilotquestionnaire.FilePath,error)  {

	query := `SELECT 
					path,file_name,number,is_new,is_sign,is_training
				FROM 
					questionnaire_path
				WHERE 
					project_id = ? and is_sign = 1
				ORDER BY number DESC;`
	
	var Pth []*pilotquestionnaire.FilePath
	err := db.DB.Select(&Pth,query,PjID)
	if err != nil {
		return nil, err
	}
	return Pth,err
}

func (or *OperationTrainingRepository)GetTrainingInfoISSign(PjID string) (*TrainingNullCase,error){

	query := `SELECT 
					project_id, meeting_date, meeting_time, is_online, ms_link, is_onsite,
					is_room1, is_room2, is_room3, is_room4, is_other, other, note,is_bypass,is_cancel
			FROM 
				operation_training
			WHERE 
				project_id = ?`
	

	var PQ TrainingNullCase
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
						project_id = ? and is_sign = 1
					ORDER BY number DESC;`

	var Pth []pilotquestionnaire.FilePath
	// err = db.DB.Select(&Pth,queryPath,PjID)
	// if err != nil {
	// 	return nil, err
	// }
	err = db.DB.Select(&Pth,queryPath,PjID)
	if err != nil {
		// ถ้าไม่มี path ก็ให้เป็น [] ไม่ error เพราะอาจไม่มีไฟล์ได้
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		Pth = []pilotquestionnaire.FilePath{}
	}
	PQ.FilePath = Pth

	queryDetailOnChange := `SELECT 
								detail_on_change
							FROM 
								detail_on_change
							WHERE 
								project_id = ? and stage = '5';`
	var detailOnChange null.String
	// err = db.DB.Get(&detailOnChange,queryDetailOnChange,PjID)
	// if err != nil {
	// 	return nil, err
	// }
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

func (or *OperationTrainingRepository)Editdetail(t *Training) error {
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `UPDATE operation_training 
				SET  
					project_id = :project_id, meeting_date = :meeting_date, meeting_time = :meeting_time, is_online = :is_online, ms_link = :ms_link
					,is_onsite = :is_onsite,is_room1 = :is_room1, is_room2 = :is_room2, is_room3 = :is_room3, is_room4 = :is_room4, is_other = :is_other, other = :other, note = :note
					,updated_by = :updated_by, updated_at = :updated_at,is_cancel = 0
				WHERE 
					project_id = :project_id;`
	
	
	_, err = tx.NamedExec(query, t)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	return tx.Commit()
}

func (or *OperationTrainingRepository)ByPassTrraining(t *Training) error{
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO operation_training 
						(project_id,note,is_bypass,created_by,created_at) 
					VALUES 
						(:project_id,:note,:is_bypass,:created_by,:created_at);`

	_,err = tx.NamedExec(query,t)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}
	return tx.Commit()
}

func (or *OperationTrainingRepository)CancelTraining(t *Training) error{
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE operation_training 
				SET 
					note = :note,
					is_cancel = 1,
					updated_by = :updated_by,
					updated_at = :updated_at
				WHERE 
					project_id = :project_id;`

	_,err = tx.NamedExec(query,t)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	return tx.Commit()
}

func (or *OperationTrainingRepository)SelectFileToTraining(t *Training) error  {
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE questionnaire_path 
				SET 
					is_training = :is_training
				WHERE 
					project_id = :project_id and number = :number;`
	for _, ps := range t.FilePath{
		param := map[string]interface{}{
			"project_id":t.ProjectID,
			"number":ps.Number,
			"is_training":ps.ISTraining,
		}
		_,err = tx.NamedExec(query,param)
		if err != nil {
			log.Println("SQL error:", err)
			return err
		}
	}

	return tx.Commit()

}

func (or *OperationTrainingRepository)UpdateTrainingPath(pth *pilotquestionnaire.Path) error{
	tx,err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	isnewUpdateQuery := `UPDATE questionnaire_path 
							SET 
								is_new = 0
							WHERE 
								project_id = ? and is_new = 1;`
	_,err = tx.Exec(isnewUpdateQuery,pth.ProjectID)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	query := `INSERT INTO 
					questionnaire_path (project_id,path,file_name,number,is_new,is_sign,is_training)
				VALUES 
					(:project_id,:path,:file_name,:number,:is_new,:is_sign,:is_training);`

	for _, ps := range pth.FilePath{
		param := map[string]interface{}{
			"project_id":pth.ProjectID,
			"path":ps.Path,
			"file_name":ps.FileName,
			"number":ps.Number,
			"is_new":ps.IsNew,
			"is_sign":ps.IsSign,
			"is_training":ps.ISTraining,
		}
		_,err = tx.NamedExec(query,param)
		if err != nil {
			log.Println("SQL error:", err)
			return err
		}
	}
	return tx.Commit()
}

func (or *OperationTrainingRepository)GetLatestTrainingPath(PjID string) ([]pilotquestionnaire.FilePath,error)  {

	query := `SELECT 
					path,file_name,number,is_new,is_sign,is_training
				FROM 
					questionnaire_path
				WHERE 
					project_id = ? and is_training = 1`
	
	var Pth []pilotquestionnaire.FilePath
	err := db.DB.Select(&Pth,query,PjID)
	if err != nil {
		return nil, err
	}
	return Pth,err
}

func (or *OperationTrainingRepository)DetailOnChangeTraining(dtl *Training) error {
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE operation_training 
				SET 
					detail_on_change = :detail_on_change,is_cancel = 0
				WHERE 
					project_id = :project_id;`

	_,err = tx.NamedExec(query,dtl)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}
	return tx.Commit()
}