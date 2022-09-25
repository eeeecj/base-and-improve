package core

type Serialize interface {
	Deserialized()
	Serialized()
}

type FieldType string

type FieldTypes []FieldType

func (f FieldTypes) exist(target FieldType) bool {
	for _, ff := range f {
		if ff == target {
			return true
		}
	}
	return false
}

const (
	INT     FieldType = "int"
	VARCHAR FieldType = "string"
	FLOAT   FieldType = "float"
)

var TYPE_MAP = map[string]FieldType{
	"int":     INT,
	"float":   FLOAT,
	"string":  VARCHAR,
	"INT":     INT,
	"FLOAT":   FLOAT,
	"VARCHAR": VARCHAR,
}

type FieldKey string
type FieldKeys []FieldKey

func (f FieldKeys) exist(target FieldKey) bool {
	for _, ff := range f {
		if ff == target {
			return true
		}
	}
	return false
}

const (
	PRIMARY   = "PRIMARY KEY"
	INCREMENT = "AUTO_INCREMENT"
	UNIQUE    = "UNIQUE"
	NOT_NULL  = "NOT NULL"
	NULL      = "NULL"
)
