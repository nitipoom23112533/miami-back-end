package project
import (
	"time"
	validation "github.com/go-ozzo/ozzo-validation"
)


// Member struct
type Member struct {
	ID        int64     `db:"id" json:"id"`
	ProjectID int64     `db:"project_id" json:"projectID"`
	UID       string    `db:"uid" json:"uid"`
	Position  string    `db:"position" json:"position"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	CreatedBy string    `db:"created_by" json:"createdBy"`

	Firstname string `db:"firstname" json:"firstname"`
	Lastname  string `db:"lastname" json:"lastname"`
	Email     string `db:"email" json:"email"`
 	  
}

// ValidateCreate func
func (x *Member) ValidateCreate() error {
	return validation.ValidateStruct(x,
		validation.Field(&x.ProjectID, validation.Required),
		validation.Field(&x.UID, validation.Required),
		validation.Field(&x.Position, validation.Required),
		validation.Field(&x.CreatedAt, validation.Required),
		validation.Field(&x.CreatedBy, validation.Required),
	)
}