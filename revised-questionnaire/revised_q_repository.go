package revisedquestionnaire

import(
	"miami-back-end/db"
	"github.com/jmoiron/sqlx"
	"miami-back-end/pilot-questionnaire"
	"log"
	"database/sql"
	"errors"
)

type RQRepository struct{

}

func NewRQRepository() *RQRepository{
	return &RQRepository{}
}

func (rqr *RQRepository)GetRevisedQuestionnaire(project_id string ) ([]pilotquestionnaire.FilePath,error){
	query := `SELECT 
					path,file_name,number,is_new,is_sign,is_training,is_revised
				FROM 
					questionnaire_path
				WHERE 
					project_id = ? and is_new = 1
				ORDER BY number DESC`
	
	var Pth []pilotquestionnaire.FilePath
	err := db.DB.Select(&Pth,query,project_id)
	if err != nil {
		return nil, err
	}
	return Pth,err
}

func (rqr *RQRepository)GetRevisedQuestionnaireDetailOnChange(project_id string,stage string) (DetailOnChange,error){
	query := `SELECT 
					detail_on_change,stage
				FROM 
					detail_on_change
				WHERE 
					project_id = ? and stage = ?`
	var detailOnChange DetailOnChange
	err := db.DB.Get(&detailOnChange,query,project_id,stage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ไม่เจอข้อมูล => คืน struct ว่าง กับ nil error
			return DetailOnChange{}, nil
		}
		// กรณี error อื่น
		return DetailOnChange{}, err
	}

return detailOnChange,nil
}

func (rqr *RQRepository)InsertRevisedQuestionnaire(p *pilotquestionnaire.Path,fileNumber int) error  {
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
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
						(project_id,path,number,file_name,is_new,is_sign,is_training,is_revised) 
					VALUES 
						(?,?,?,?,1,1,1,1);`
	
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

func (rqr *RQRepository)RQRDetailOnChange(p *pilotquestionnaire.PilotQuestionnaire,stage string) error {
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `INSERT INTO 
			detail_on_change (project_id, detail_on_change, stage) 
			VALUES 
				(?, ?, ?)
			ON DUPLICATE KEY UPDATE 
				detail_on_change = VALUES(detail_on_change);`

	_, err = tx.Exec(query, p.ProjectID, p.DetailOnChange, stage)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	return tx.Commit()
}