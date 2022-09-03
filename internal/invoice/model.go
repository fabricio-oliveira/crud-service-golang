package invoice

type Invoice struct {
	Id        string  `json:"id" binding:"required"`
	BillTo    string  `json:"bill_to" binding:"required"`
	Items     []Items `json:"items"`
	CreatedAt string  `json:"created_at"`
	UpdateAt  string  `json:"updated_at"`
}

type Items struct {
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unity"`
	Price       int64  `json:"price"`
}
