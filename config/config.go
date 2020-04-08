package config

var FormatString = "string"
var FormatJson = "json"

var FieldTypeString = "string"
var FieldTypeInt = "int"
var FieldTypeFloat = "float"
var FieldTypeBool = "bool"
var FieldTypeTime = "datetime"

// CheckType 检查是否是有效的类型
func CheckType(name string) bool {
	if name == FieldTypeString {
		return true
	}
	if name == FieldTypeInt {
		return true
	}
	if name == FieldTypeFloat {
		return true
	}
	if name == FieldTypeBool {
		return true
	}
	if name == FieldTypeTime {
		return true
	}
	return false
}
