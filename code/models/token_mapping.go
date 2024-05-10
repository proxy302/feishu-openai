package models

type TokenMapping struct {
	Model
	ExternalTokenID int `json:"external_token_id"`
	InternalTokenID int `json:"internal_token_id"`
	Status          int `json:"status"`
	ExpiredOn       int `json:"expired_on"`
	CurrentCost     int `json:"current_cost"`
	LimitCost       int `json:"limit_cost"`
	Uid             int `json:"uid"`
	IsRobot         int `json:"is_robot"`
	ModelID         int `json:"model_id"`
}

func (TokenMapping) TableName() string {
	return "t_token_mapping"
}

func GetTokenMappingByID(id int) (*TokenMapping, error) {
	var result TokenMapping
	if err := db.Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
