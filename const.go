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
	Like        = "LIKE"
	NotLike     = "NOT LIKE"
)

const (
	Question = iota
	DollarNumber
	Named
)

const (
	LeftJoin  = "LEFT JOIN"
	RightJoin = "RIGHT JOIN"
	InnerJoin = "INNER JOIN"
)

const (
	Asc  = "ASC"
	Desc = "DESC"
)
