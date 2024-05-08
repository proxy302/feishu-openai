package models

type Token struct {
	Model
	Value            string `json:"value"`
	TokenSuppliserID int    `json:"token_supplier_id"`
}

func (Token) TableName() string {
	return "t_token"
}

func GetTokenByID(id int) (*Token, error) {
	var result Token
	if err := db.Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
