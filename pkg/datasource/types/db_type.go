package types

//go:generate stringer -type=DBType
type DBType int16

const (
	_ DBType = iota
	DbTypeUnknown
	DbTypeMySQL
	DbTypePostgreSQL
	DbTypeSQLServer
	DbTypeOracle
)
