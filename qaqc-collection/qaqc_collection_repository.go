package qaqccollection

import(
	"miami-back-end/db"
	"miami-back-end/data-collection"
	"github.com/jmoiron/sqlx"
	"database/sql"
	"errors"
)

type QaqcCollectionRepository struct{

}
func NewQaqcCollectionRepository() *QaqcCollectionRepository{
	return &QaqcCollectionRepository{}
}

func (qr *QaqcCollectionRepository)StartQaqcCollection(dct *datacollection.DataCollection) error{
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO qaqc_collection 
					(project_id,is_start,quota,day,start_date,created_at,created_by) 
				VALUES 
					(:project_id,:is_start,:quota,:day,:start_date,:created_at,:created_by) `

	_, err = tx.NamedExec(query,dct)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (qr *QaqcCollectionRepository)CompletedQaqcCollection(dct *datacollection.DataCollection) error{
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE qaqc_collection 
					SET is_completed = :is_completed,completed_date = :completed_date,updated_by = :updated_by,updated_at = :updated_at
				WHERE 
					project_id = :project_id`

	_, err = tx.NamedExec(query,dct)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (qr *QaqcCollectionRepository) GetQaqcCollection(PjID string) (datacollection.DataCollection, error) {
	query := `SELECT 
					project_id,is_start,quota,start_date,is_completed,completed,completed_date,created_by,created_at,updated_by,updated_at
				FROM 
					qaqc_collection
				WHERE 
					project_id = ?`

	var dct datacollection.DataCollection
	err := db.DB.Get(&dct, query, PjID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ไม่พบข้อมูล -> คืน struct ว่างกับ nil error
			return datacollection.DataCollection{}, nil
		}
		// error อื่น ๆ เช่น database error
		return datacollection.DataCollection{}, err
	}

	return dct, nil
}
func (qr *QaqcCollectionRepository)GetQcCollection(Pj string,statuses []string) ([]datacollection.SsAndFsResponses,error){

	query, args, err := sqlx.In(`
		SELECT id, project_id, status
		FROM fs_responses_destiny
		WHERE project_id = ? 
		AND status IN (?);`, Pj, statuses)

	if err != nil {
		return nil, err
	}

	// sqlx.In จะขยายเป็น status IN (?,?,?,?) พร้อม args
	query = db.DB.Rebind(query)

	var fsrqc []datacollection.SsAndFsResponses
	err = db.DB.Select(&fsrqc, query, args...)
	if err != nil {
		return nil, err
	}

	return fsrqc, nil
}

func (qr *QaqcCollectionRepository)GetQaCollection(Pj string) ([]datacollection.SsAndFsResponses,error){
	
	query := `SELECT 
					id,project_id
				FROM 
					qa_destiny
				WHERE
					project_id = ?;`
	var fsrqa []datacollection.SsAndFsResponses
	err := db.DB.Select(&fsrqa, query, Pj)
	if err != nil {
		return nil, err
	}
	return fsrqa,nil
}