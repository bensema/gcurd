package gcurd

type (
	WhereValue struct {
		Name  string      `json:"name"`
		Op    Op          `json:"op"`
		Value interface{} `json:"value"`
	}

	OrderBy struct {
		Direction string `json:"direction"`
		Filed     string `json:"filed"`
	}

	Pagination struct {
		Num  int `json:"num"`
		Size int `json:"size"`
	}

	Request struct {
		Where      []WhereValue
		OrderBy    OrderBy
		Pagination Pagination
	}
)
