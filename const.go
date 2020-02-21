package query_builder

const (
	Equal           = "="           // tag.operator.eq
	GraterThan      = ">"           // tag.operator.gt
	GraterThanEqual = ">="          // tag.operator.gte
	LessThan        = "<"           // tag.operator.lt
	LessThanEqual   = "<="          // tag.operator.lte
	NotEqual        = "!="          // tag.operator.ne
	Like            = "LIKE"        // tag.operator.like
	NotLike         = "NOT LIKE"    // tag.operator.not-like
	IsNull          = "IS NULL"     // tag.operator.is-null
	IsNotNull       = "IS NOT NULL" // tag.operator.not-null
	In              = "IN"
	NotIn           = "NOT IN"
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

const (
	DBTag       = "db"
	TableTag    = "table"
	SearchTag   = "search"
	OperatorTag = "operator"
)
