package project
import (
	"time"

	"gopkg.in/guregu/null.v4"
	validation "github.com/go-ozzo/ozzo-validation"
)

type Service struct {
	Projectrepository *Repository
}
func NewService() *Service {
	return &Service{}
}
// Project struct
type Project struct {
	ID        int64       `db:"id" json:"id"`
	Name      string      `db:"name" json:"name"`
	Code      null.Int    `db:"code" json:"code"`
	Year      int64       `db:"year" json:"year"`
	Status    string      `db:"status" json:"status"`
	CreatedAt time.Time   `db:"created_at" json:"createdAt"`
	CreatedBy string      `db:"created_by" json:"createdBy"`
	UpdatedAt null.Time   `db:"updated_at" json:"updatedAt"`
	UpdatedBy null.String `db:"updated_by" json:"updatedBy"`

	Directors []UserPosition `json:"directors"`
	Owners    []UserPosition `json:"owners"`
}

// UserPosition struct
type UserPosition struct {
	UID      string `json:"uid"`
	Fullname string `json:"fullname"`
}

// ValidateCreate func
func (p *Project) ValidateCreate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Year, validation.Required, validation.Min(0)),
		validation.Field(&p.Status, validation.Required),
		validation.Field(&p.CreatedAt, validation.Required),
		validation.Field(&p.CreatedBy, validation.Required),
	)
}
func (s *Service) CreateMemberAndStage(x *Project ,m *Member) error {
	return s.Projectrepository.CreateMemberAndStage(x,m)
}
func (s *Service) GetProjectsByUIDAndStatus(uid string,status string,isAdmin bool) ([]Project,error)  {
	return s.Projectrepository.GetProjectsByUIDAndStatus(uid,status,isAdmin)
}
func (s *Service) GetProjectByID(id int64)(*Project,error){
	return s.Projectrepository.GetProjectByID(id)
}