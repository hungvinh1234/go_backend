package model

// Pagination .
type Pagination struct {
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

// OrderBy .
type OrderBy struct {
	Order string `json:"order,omitempty"`
	By    string `json:"by,omitempty"`
}
