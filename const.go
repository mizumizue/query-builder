package query_builder

const (
	Equal       = "="
	GraterThan  = ">"
	GraterEqual = ">="
	LessThan    = "<"
	LessEqual   = "<="
	Not         = "!="
	In          = "IN"
	NotIn       = "NOT IN"
)

const (
	Question = iota
	Named
)

const (
	LeftJoin  = "LEFT JOIN"
	RightJoin = "RIGHT JOIN"
)

const (
	Asc  = "ASC"
	Desc = "DESC"
)
