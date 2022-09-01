package invoice

type Invoice struct {
	Id        int32   `json:"id" binding:"required"`
	BillTo    string  `json:"bill_to" binding:"required"`
	Items     []Items `json:"items"`
	CreatedAt string
	UpdateAt  string
}

type Items struct {
	Description string `json:"description"`
	Quantity    int32  `json:"quantity"`
	Unit        string `json:"unity"`
	Price       int64  `json:"price"`
}
