package invoice

type Invoice struct {
	Id          string  `json:"id" binding:"required"`
	Address     string  `json:"address" binding:"required"`
	CompanyName string  `json:"company_name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Goods       []Goods `json:"goods"`
	Amount      string  `json:"amount"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type Goods struct {
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unity"`
	Price       string `json:"price"`
}
