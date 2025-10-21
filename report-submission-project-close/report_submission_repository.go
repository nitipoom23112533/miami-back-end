package reportsubmissionprojectclose

import (
	"miami-back-end/db"
	"database/sql"
	"errors"
)

type ReportSubmissionRepo struct{
	
}

func NewReportSubmissionRepo() *ReportSubmissionRepo{
	return &ReportSubmissionRepo{}
}

func (rs *ReportSubmissionRepo)GetReportSubmission(PjID string) (ReportSubmission,error){
	query := `SELECT 
					project_id, submission_date, close_date
				FROM 
					report_submission
				WHERE 
					project_id = ?`
	var submission ReportSubmission
	err := db.DB.Get(&submission,query,PjID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ไม่พบข้อมูล -> คืน struct ว่างกับ nil error
			return ReportSubmission{}, nil
		}
		// error อื่น ๆ เช่น database error
		return ReportSubmission{}, err
	}

	return submission,nil
}

func (rs *ReportSubmissionRepo)InsertReportSubmission(submission *ReportSubmission) error{
	tx, err := db.DB.Beginx() // ใช้ Beginx() ถ้าใช้ sqlx
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO report_submission 
					(project_id,submission_date,created_at,created_by) 
				VALUES 
					(:project_id,:submission_date,:created_at,:created_by) `

	_, err = tx.NamedExec(query,submission)
	if err != nil {
		return err
	}

	return tx.Commit()

}

func (rs *ReportSubmissionRepo)UpdateReportSubmission(submission *ReportSubmission) error{
	tx, err := db.DB.Beginx()

	if err != nil {
		return err
	}
	defer tx.Rollback()

query := `UPDATE report_submission 
				SET 
					close_date = :close_date,
					updated_at = :updated_at,
					updated_by = :updated_by
				WHERE 
					project_id = :project_id`

	_, err = tx.NamedExec(query,submission)
	if err != nil {
		return err
	}

	return tx.Commit()
}