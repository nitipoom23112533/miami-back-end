package questionnairesignoff

import(
	"miami-back-end/pilot-questionnaire"
	"miami-back-end/db"
	"log"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"

)


type QuestionnaireSignOffRepository struct{


}
func NewQuestionnaireSignOffRepository() *QuestionnaireSignOffRepository{
	return &QuestionnaireSignOffRepository{}
}

func (se *QuestionnaireSignOffRepository)GetAllPilotQuestionnairePath(PjID string) ([]*pilotquestionnaire.FilePath,error)  {

	query := `SELECT 
					path,file_name,number,is_new,is_sign,is_training
				FROM 
					questionnaire_path
				WHERE 
					project_id = ? and is_new = 1
				ORDER BY number DESC`
	
	var Pth []*pilotquestionnaire.FilePath
	err := db.DB.Select(&Pth,query,PjID)
	if err != nil {
		return nil, err
	}
	return Pth,err
}

func (se *QuestionnaireSignOffRepository)GetSignOffDetailOnChange(project_id string,stage string) (DetailOnChangeSignOff,error){
	query := `SELECT 
					detail_on_change,stage
				FROM 
					detail_on_change
				WHERE 
					project_id = ? and stage = ?`
	var detailOnChange DetailOnChangeSignOff
	err := db.DB.Get(&detailOnChange,query,project_id,stage)	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ไม่เจอข้อมูล => คืน struct ว่าง กับ nil error
			return DetailOnChangeSignOff{}, nil
		}
		// กรณี error อื่น
		return DetailOnChangeSignOff{}, err
	}
	return detailOnChange,err
}

func (se *QuestionnaireSignOffRepository)SelectFileToSignOff(pth *pilotquestionnaire.Path) error  {
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE questionnaire_path 
				SET 
					is_sign = :is_sign
				WHERE 
					project_id = :project_id and number = :number;`
	for _, ps := range pth.FilePath{
		param := map[string]interface{}{
			"project_id":pth.ProjectID,
			"number":ps.Number,
			"is_sign":ps.IsSign,
		}
		_,err = tx.NamedExec(query,param)
		if err != nil {
			log.Println("SQL error:", err)
			return err
		}
	}

	detailOnChangeQuery := `INSERT INTO detail_on_change (project_id,stage) VALUES (?,'4')`
	_,err = tx.Exec(detailOnChangeQuery,pth.ProjectID)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	return tx.Commit()
}

func (se *QuestionnaireSignOffRepository)InsertQuestionnaireSignOff(p *pilotquestionnaire.Path,fileNumber int) error  {
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
						(project_id,path,number,file_name,is_new,is_sign) 
					VALUES 
						(?,?,?,?,1,1);`
	
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

