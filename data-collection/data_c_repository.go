package datacollection

import(
	"miami-back-end/db"
	"errors"
	"database/sql"
)

type DataCollectionRepository struct{

}

func NewDataCollectionRepository() *DataCollectionRepository{
	return &DataCollectionRepository{}
}

func (dr *DataCollectionRepository)StartDataCollection(dct *DataCollection) error{
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO data_collection 
					(project_id,is_start,quota,day,start_date,created_at,created_by) 
				VALUES 
					(:project_id,:is_start,:quota,:day,:start_date,:created_at,:created_by) `

	_, err = tx.NamedExec(query,dct)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (dr *DataCollectionRepository)CompletedDataCollection(dct *DataCollection) error{
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE data_collection 
					SET is_completed = :is_completed,completed_date = :completed_date,updated_by = :updated_by,updated_at = :updated_at
				WHERE 
					project_id = :project_id`

	_, err = tx.NamedExec(query,dct)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (dr *DataCollectionRepository)GetDataCollection(PjID string) (*DataCollection,error){
	query := `SELECT 
					project_id,is_start,quota,start_date,is_completed,completed,completed_date,created_by,created_at,updated_by,updated_at
				FROM 
					data_collection
				WHERE 
					project_id = ?`
	var dct DataCollection
	err := db.DB.Get(&dct, query, PjID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ไม่พบข้อมูล -> คืน nil, nil
			return nil, nil
		}
		// error อื่น ๆ เช่น database error
		return nil, err
	}
	return &dct,err

}

func  (dr *DataCollectionRepository)GetDashBoardLogs(PjID string) (*DashboardLogs,error)  {
	query := `SELECT 
					project_id,mn,qc,qa,fw,de,da,doc,doc_reject
				FROM 
					dashboard_logs_destiny
				WHERE
					project_id = ?;`
	var dbl DashboardLogs
	err := db.DB.Get(&dbl, query, PjID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ไม่พบข้อมูล -> คืน nil, nil
			return nil, nil
		}
		// error อื่น ๆ เช่น database error
		return nil, err
	}
	return  &dbl,err
	
}

func (dr *DataCollectionRepository)GetSsResponses(Pj string) (*[]SsAndFsResponses,error){
	query := `SELECT 
					id,project_id
				FROM 
					ss_responses_destiny
				WHERE
					project_id = ?;`
	var ssr []SsAndFsResponses
	err := db.DB.Select(&ssr, query, Pj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {	
			return nil, nil
		}
		return nil, err
	}
	return &ssr,err
}

func (dr *DataCollectionRepository)GetFsResponses(Pj string) (*[]SsAndFsResponses,error){
	query := `SELECT 
					id,project_id,status
				FROM 
					fs_responses_destiny
				WHERE
					project_id = ?;`
	var fsr []SsAndFsResponses
	err := db.DB.Select(&fsr, query, Pj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {	
			return nil, nil
		}
		return nil, err
	}
	return &fsr,err
}