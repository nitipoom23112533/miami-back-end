package members

import (
	"miami-back-end/db"
	"github.com/jmoiron/sqlx"
)

type MemberRepository struct{

}
func NewMemberRepository() *MemberRepository{
	return &MemberRepository{}
}


func (mr *MemberRepository)GetMemberByPosition() ([]Employee,error)  {
	query := `SELECT 
					uid,firstname,lastname,position,status,email
				FROM users 
				WHERE 
					status = 'active' and employment_type = 'full-time'
				AND 
					position in ('Fieldwork Manager','Recruiter','Project Manager','Project Director','Senior Data Entry','Data Entry','Data Analyst','Data Analysis','QC Manager','Research Executive','Senior Research Executive') 
				GROUP BY 
				position,uid,firstname,lastname,status,email;`
	var employee []Employee
	err := db.DB.Select(&employee,query)
	if err != nil {
		return nil,err

	}
	return employee,err
}

func (mr *MemberRepository)GetMemberByPjID(Pjid string) ([]MembersOfPj,error) {
	query := `SELECT 
					pms.uid,pms.project_id,u.firstname,u.lastname,pms.position, u.position as role,u.email,pms.is_send_email,pms.member_status
					FROM 
						p_members AS pms 
					LEFT JOIN 
						users AS u ON u.uid = pms.uid 
					WHERE 
						pms.project_id = ? and pms.member_status = 'active'
					ORDER BY 
						pms.created_at DESC;`
	
	var members []MembersOfPj
	err := db.DB.Select(&members,query,Pjid)
	if err != nil {
		return nil,err

	}
	return members,err
}
func (mr *MemberRepository)GetMemberAllByPjID(Pjid string) ([]MembersOfPj,error) {
	query := `SELECT 
					pms.uid,pms.project_id,u.firstname,u.lastname,pms.position, u.position as role,u.email,pms.is_send_email,pms.member_status,pms.created_at,pms.created_by,pms.updated_at,pms.updated_by
					FROM 
						p_members AS pms 
					LEFT JOIN 
						users AS u ON u.uid = pms.uid 
					WHERE 
						pms.project_id = ? 
					ORDER BY 
						pms.created_at DESC;`
	var members []MembersOfPj
	err := db.DB.Select(&members,query,Pjid)
	if err != nil {
		return nil,err

	}
	return members,err
}

func (mr *MemberRepository) GetOutOfMemberAllByPjID(Pjid string,Uid []string) ([]MembersOfPj, error) {
	query := `SELECT 
					pms.uid, pms.project_id, u.firstname, u.lastname, pms.position, 
					u.position as role, u.email, pms.is_send_email, 
					pms.member_status, pms.created_at, pms.created_by, 
					pms.updated_at, pms.updated_by
				FROM 
					p_members AS pms 
				LEFT JOIN 
					users AS u ON u.uid = pms.uid 
				WHERE 
					pms.project_id = ? AND pms.uid NOT IN (?)
				ORDER BY 
					pms.created_at DESC;`

	// ใช้ sqlx.In เพื่อขยาย slice -> (?, ?, ?)
	query, args, err := sqlx.In(query,Pjid, Uid)
	if err != nil {
		return nil, err
	}

	// ทำให้ query รองรับ database driver (เช่น เปลี่ยน ? → $1 สำหรับ postgres)
	query = db.DB.Rebind(query)

	var members []MembersOfPj
	err = db.DB.Select(&members, query, args...)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func  (mr *MemberRepository)AddMemberByPjID(mops *[]MembersOfPj) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	projectID := (*mops)[0].ProjectID

	deleteQuery := `DELETE FROM p_members WHERE project_id = ?`
	_, err = db.DB.Exec(deleteQuery, projectID)
	if err != nil {
		return err
	}

	query := `INSERT INTO 
				p_members (
					project_id, uid, position, member_status, is_send_email, created_at, created_by, updated_at, updated_by
				) VALUES (
					:project_id, :uid, :position, :member_status, :is_send_email, :created_at, :created_by, :updated_at, :updated_by
				);`
									
	for _, mop := range *mops {
		_, err = tx.NamedExec(query, mop)
		if err != nil {
			return err
		}
	}
	
	return tx.Commit()

}

func (mr *MemberRepository)UpdateIsSendEmail(mops *[]MembersOfPj) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE 
					p_members
				SET 
					is_send_email = 1
				WHERE
					project_id = :project_id AND uid = :uid;`
	
	for _, mop := range *mops {
		_, err := db.DB.NamedExec(query, mop)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}