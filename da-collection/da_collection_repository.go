package dacollection
import(
	"miami-back-end/db"
	"miami-back-end/data-collection"
	"database/sql"
	"errors"
)

type DaCollectionRepo struct{

}
func NewDaCollectionRepository() *DaCollectionRepo{
	return &DaCollectionRepo{}
}

func (dr *DaCollectionRepo)StartDaCollection(dct *datacollection.DataCollection) error{
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO da_collection 
					(project_id,is_start,quota,start_date,created_at,created_by) 
				VALUES 
					(:project_id,:is_start,:quota,:start_date,:created_at,:created_by) `

	_, err = tx.NamedExec(query,dct)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (dr *DaCollectionRepo)CompletedDaCollection(dct *datacollection.DataCollection) error{
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE da_collection 
					SET is_completed = :is_completed,completed_date = :completed_date,updated_by = :updated_by,updated_at = :updated_at
				WHERE 
					project_id = :project_id`

	_, err = tx.NamedExec(query,dct)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (dr *DaCollectionRepo)GetDaCollection(PjID string) ([]datacollection.SsAndFsResponses,error){
	query := `SELECT 
					id,project_id
				FROM 
					da_destiny
				WHERE
					project_id = ?;`
	var dac []datacollection.SsAndFsResponses
	err := db.DB.Select(&dac, query, PjID)
	if err != nil {
		return nil, err
	}
	return dac, nil
}

func (dr *DaCollectionRepo)GetDaCollectionInfo(PjID string) (datacollection.DataCollection,error){
	query := `SELECT 
					project_id,is_start,quota,start_date,is_completed,completed,completed_date,created_by,created_at,updated_by,updated_at
				FROM 
					da_collection
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