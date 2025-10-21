package stage


type StageService struct {
	StageRepository *StageRepository
}

type Stage struct {
	ProjectID string `db:"project_id" json:"project_id"`
	Stage1 bool `db:"stage_1" json:"stage_1"`
	Stage2 bool `db:"stage_2" json:"stage_2"`
	Stage3 bool `db:"stage_3" json:"stage_3"`
	Stage4 bool `db:"stage_4" json:"stage_4"`
	Stage5 bool `db:"stage_5" json:"stage_5"`
	Stage5_1 bool `db:"stage_5_1" json:"stage_5_1"`
	Stage6 bool `db:"stage_6" json:"stage_6"`
	Stage7 bool `db:"stage_7" json:"stage_7"`
	Stage8 bool `db:"stage_8" json:"stage_8"`
	Stage9 bool `db:"stage_9" json:"stage_9"`
	Stage10 bool `db:"stage_10" json:"stage_10"`
	Stage11 bool `db:"stage_11" json:"stage_11"`
	Stage12 bool `db:"stage_12" json:"stage_12"`
	Stage13 bool `db:"stage_13" json:"stage_13"`
	Stage14 bool `db:"stage_14" json:"stage_14"`
	Stage15 bool `db:"stage_15" json:"stage_15"`
	Stage16 bool `db:"stage_16" json:"stage_16"`
	Stage17 bool `db:"stage_17" json:"stage_17"`
	Stage18 bool `db:"stage_18" json:"stage_18"`
}

func NewStageService() *StageService {
	return &StageService{}
}

func (s *StageService)UpdateStage(projectId string,stage string) error {
	return s.StageRepository.UpdateStage(projectId,stage)
}

func (s *StageService) GetStage(projectId string) (*Stage,error) {
	return s.StageRepository.GetStage(projectId)

}
