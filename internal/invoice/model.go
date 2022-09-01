package invoice

type Invoice struct {
	Id     int32
	Date   string
	BillTo string
	Items  []Items
}

type Items struct {
	Description string
	Quantity    int32
	Unit        string
	Price       int64
}
