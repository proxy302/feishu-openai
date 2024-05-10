package models

type Models struct {
	ID                int    `gorm:"primary_key" json:"id"`
	Name              string `json:"name"`
	ShowInRoot        int    `json:"show_in_root"`
	ShowInApi         int    `json:"show_in_api"`
	Ord               int    `json:"ord"`
	TokenSupplierID   int    `json:"token_supplier_id"`
	TokenSupplierList string `json:"token_supplier_list"`
}

func (Models) TableName() string {
	return "t_models"
}

func GetModelByID(id int) (*Models, error) {
	var result Models
	if err := db.Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
