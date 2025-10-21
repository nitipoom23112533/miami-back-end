package projectreview

import(
	"miami-back-end/db"
	"database/sql"
	"errors"
	"log"
)

type ProjectReviewRepo struct{
	
}

func NewProjectReviewRepo() *ProjectReviewRepo{
	return &ProjectReviewRepo{}

}

func (pr *ProjectReviewRepo)CreateProjectReview(prw *ProjectReview) error{
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO project_review 
					(project_id,meeting_date,meeting_time,is_online,ms_link,is_onsite,is_room1,is_room2,is_room3,is_room4,is_other,other,note,created_by,created_at) 
				VALUES 
					(:project_id,:meeting_date,:meeting_time,:is_online,:ms_link,:is_onsite,:is_room1,:is_room2,:is_room3,:is_room4,:is_other,:other,:note,:created_by,:created_at);`

	_, err = tx.NamedExec(query,prw)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (pr *ProjectReviewRepo)EditProjectReview(prw *ProjectReview) error{
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE project_review 
				SET 
					meeting_date = :meeting_date,
					meeting_time = :meeting_time,
					is_online = :is_online,
					ms_link = :ms_link,
					is_onsite = :is_onsite,
					is_room1 = :is_room1,
					is_room2 = :is_room2,
					is_room3 = :is_room3,
					is_room4 = :is_room4,
					is_other = :is_other,
					other = :other,
					note = :note,
					updated_by = :updated_by,
					updated_at = :updated_at,
					is_cancel = 0
				WHERE 
					project_id = :project_id;`
	
	_, err = tx.NamedExec(query,prw)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (pr *ProjectReviewRepo)GetProjectReview(PjID string) (*ProjectReview,error){
	query := `SELECT 
					project_id,meeting_date,meeting_time,is_online,ms_link,is_onsite,is_room1,is_room2,is_room3,is_room4,is_other
					,other,note,created_by,created_at,updated_by,updated_at,is_bypass,is_cancel
				FROM 
					project_review
				WHERE 
					project_id = ?`
	var prw ProjectReview
	err := db.DB.Get(&prw,query,PjID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ไม่พบข้อมูลในฐานข้อมูล
			return nil, nil 
		}
		
		return nil, err
	}
	return &prw,nil
}

func (pr *ProjectReviewRepo)CancelProjectReview(prw *ProjectReview) error{
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE project_review
				SET 
					note = :note,
					is_cancel = 1,
					updated_by = :updated_by,
					updated_at = :updated_at
				WHERE 
					project_id = :project_id;`

	_,err = tx.NamedExec(query,prw)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}

	return tx.Commit()
}

func (pr *ProjectReviewRepo)ByPassProjectReview(prw *ProjectReview) error{
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO project_review 
						(project_id,note,is_bypass,created_by,created_at) 
					VALUES 
						(:project_id,:note,:is_bypass,:created_by,:created_at);`

	_,err = tx.NamedExec(query,prw)
	if err != nil {
		log.Println("SQL error:", err)
		return err
	}
	return tx.Commit()
}

