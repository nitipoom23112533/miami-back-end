package stage

import (
	"log"
	"miami-back-end/db"
	"fmt"
)

type StageRepository struct{

}

func NewStageRepository() *StageRepository{
	return &StageRepository{}
}

func (sr *StageRepository)GetStage(projectId string) (*Stage,error){
	query := `SELECT 
				project_id,stage_1,stage_2,stage_3,stage_4,stage_5,stage_5_1,stage_6,stage_7,stage_8,stage_9,
				stage_10,stage_11,stage_12,stage_13,stage_14,stage_15,stage_16,stage_17,stage_18
			FROM 
			 	stage 
			WHERE project_id = ?;`
	var s Stage
	err := db.DB.Get(&s,query,projectId)
	return &s,err
}

func (sr *StageRepository)UpdateStage(project_id string,stage string) error{
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := fmt.Sprintf(`
		UPDATE stage 
		SET 
			%s = 1 
		WHERE 
			project_id = ?;
	`, stage)
	_,err = tx.Exec(query,project_id)
	if err != nil{
		log.Println(err)
		return err
	}
	return tx.Commit()
}
