package test

type req struct {
	ID string `form:"id" binding:"required"`
}

type res struct {
	TableName string `db:"Tables_in_test" json:"table_name"`
}
