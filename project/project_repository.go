package project

import (
	"log"
	"miami-back-end/db"
	"github.com/jmoiron/sqlx"
)
type Repository struct{

}
func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) CreateMemberAndStage(x *Project, m *Member) error  {

	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	m.ProjectID = x.ID
	query := `INSERT INTO p_members (
		project_id, uid, position,member_status, created_at, created_by
	) VALUES (
		:project_id, :uid, :position,'active', :created_at, :created_by
	)`
	result, err := tx.NamedExec(query, m)
	if err != nil {
		return err
	}
	m.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	query = `INSERT INTO stage (project_id) VALUES (?)`
	_, err = tx.Exec(query,m.ProjectID)
	if err != nil {
		log.Println(err)
		return err
	}
	return tx.Commit()
}

func (r *Repository) GetProjectsByUIDAndStatus(uid string,status string,isAdmin bool) ([]Project,error)  {

	query := `SELECT
				id, name, year, status, created_at, created_by
			FROM projects_destiny
			WHERE status = ? and 
			CASE
				WHEN ? THEN TRUE
				ELSE id IN (SELECT DISTINCT project_id FROM p_members_destiny WHERE uid = ?)
			END
			ORDER BY id DESC`
	var xs []Project
	err := db.DB.Select(&xs, query, status, isAdmin, uid)
	if err != nil {
		return nil, err
	}
	if len(xs) == 0 {
		return xs ,nil
	}
	query = `SELECT p.id AS project_id, m.position, u.uid, CONCAT(u.firstname, ' ', u.lastname) AS fullname
			FROM projects_destiny AS p
			INNER JOIN p_members_destiny AS m ON m.project_id = p.id AND m.position IN ('PD', 'PO')
			INNER JOIN users AS u ON u.uid = m.uid
			WHERE p.id IN (?)
			ORDER BY p.id DESC`

	idList := make([]int64, len(xs))
	for i, x := range xs {
		idList[i] = x.ID
	}
	query, args, err := sqlx.In(query, idList)
	if err != nil {
		return nil, err
	}
	query = db.DB.Rebind(query)
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var id int64
	idx := -1
	var x struct {
		ProjectID int64
		Position  string
		UID       string
		Fullname  string
	}
	for rows.Next() {
		if err := rows.Scan(&x.ProjectID, &x.Position, &x.UID, &x.Fullname); err != nil {
			return nil, err
		}
		if x.ProjectID != id {
			id = x.ProjectID
			idx++
		}
		for id != xs[idx].ID {
			idx++
		}
		switch x.Position {
		case "PO":
			xs[idx].Owners = append(xs[idx].Owners, UserPosition{x.UID, x.Fullname})
		case "PD":
			xs[idx].Directors = append(xs[idx].Directors, UserPosition{x.UID, x.Fullname})
		}
	}
	return xs, err
}
func (r *Repository) GetProjectByID(id int64)(*Project,error){
	query := `SELECT
				id, name, code, year, status, created_at, created_by,
				updated_at, updated_by
			FROM projects_destiny
			WHERE id = ?`
	var x Project
	err := db.DB.Get(&x, query, id)
	return &x, err

}

