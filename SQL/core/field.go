package core

import "errors"

type Field struct {
	dataType     FieldType
	keys         FieldKeys
	defaultValue interface{}
	values       []interface{}
	rows         int
}

func (f *Field) NewField(dataType FieldType, keys []FieldKey, defaultValue interface{}) *Field {
	return &Field{
		dataType:     dataType,
		keys:         keys,
		defaultValue: defaultValue,
	}
}

func (f *Field) Validate() error {
	if f.keys.exist(INCREMENT) {
		if f.dataType != INT {
			return errors.New("Increment key require Data-Type is integer")
		}
		if !f.keys.exist(PRIMARY) {
			return errors.New("Increment key require primary key")
		}
	}
	if f.defaultValue != nil && f.keys.exist(UNIQUE) {
		return errors.New("Unique key not allow to set default value")
	}
	return nil
}

func (f *Field) checkType(value FieldType) bool {
	if value == f.dataType {
		return true
	}
	return false
}
func (f *Field) checkIndex(index int) bool {
	if index > f.rows {
		return false
	}
	return true
}
func (f *Field) checkKey(value interface{}) bool {
	if f.keys.exist(INCREMENT) {
		if value == nil {
			value = f.rows + 1
		}
		//if va
	}
	return false
}
