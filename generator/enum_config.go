package generator

import (
	"fmt"
	"strconv"
	"strings"
)

// EnumConfigValue holds a configuration value with its validity flag.
type EnumConfigValue struct {
	Value interface{}
	Valid bool
}

// GetBool returns the boolean value if valid, otherwise returns the default value.
func (v *EnumConfigValue) GetBool(defaultValue bool) bool {
	if v.Valid {
		if b, ok := v.Value.(bool); ok {
			return b
		}
	}
	return defaultValue
}

// GetString returns the string value if valid, otherwise returns the default value.
func (v *EnumConfigValue) GetString(defaultValue string) string {
	if v.Valid {
		if s, ok := v.Value.(string); ok {
			return s
		}
	}
	return defaultValue
}


// EnumConfig holds configuration options specific to a single enum.
// These options can be specified inline via annotations and override global GeneratorConfig.
type EnumConfig struct {
	// Bool options
	NoPrefix        EnumConfigValue `json:"no_prefix"`
	NoIota          EnumConfigValue `json:"no_iota"`
	LowercaseLookup EnumConfigValue `json:"lowercase_lookup"`
	CaseInsensitive EnumConfigValue `json:"case_insensitive"`
	Marshal         EnumConfigValue `json:"marshal"`
	SQL             EnumConfigValue `json:"sql"`
	SQLInt          EnumConfigValue `json:"sql_int"`
	Flag            EnumConfigValue `json:"flag"`
	Names           EnumConfigValue `json:"names"`
	Values          EnumConfigValue `json:"values"`
	LeaveSnakeCase  EnumConfigValue `json:"leave_snake_case"`
	Ptr             EnumConfigValue `json:"ptr"`
	SQLNullInt      EnumConfigValue `json:"sql_null_int"`
	SQLNullStr      EnumConfigValue `json:"sql_null_str"`
	MustParse       EnumConfigValue `json:"must_parse"`
	ForceLower      EnumConfigValue `json:"force_lower"`
	ForceUpper      EnumConfigValue `json:"force_upper"`
	NoComments      EnumConfigValue `json:"no_comments"`
	NoParse         EnumConfigValue `json:"no_parse"`

	// String options
	Prefix   EnumConfigValue `json:"prefix"`
	
	// Slice/map options (not supported inline for simplicity)
	// BuildTags         []string
	// ReplacementNames  map[string]string
	// TemplateFileNames []string
}

// NewEnumConfig creates a new EnumConfig with default values.
func NewEnumConfig() *EnumConfig {
	return &EnumConfig{}
}

// ParseAnnotation parses a single annotation string (e.g., "@marshal", "@marshal:true", "@prefix=\"My\"")
// and updates the EnumConfig accordingly.
func (ec *EnumConfig) ParseAnnotation(annotation string) error {
	annotation = strings.TrimSpace(annotation)
	if annotation == "" {
		return nil
	}
	
	// Remove @ prefix
	if !strings.HasPrefix(annotation, "@") {
		return fmt.Errorf("annotation must start with @: %s", annotation)
	}
	annotation = annotation[1:]
	
	// Check for key:value format (e.g., @marshal:true, @marshal:false)
	if strings.Contains(annotation, ":") {
		parts := strings.SplitN(annotation, ":", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		// Parse boolean value
		if value == "true" || value == "false" {
			boolValue, _ := strconv.ParseBool(value)
			return ec.setBoolOption(key, boolValue)
		}
		
		// String value (could be quoted)
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || 
			(value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}
		
		return ec.setStringOption(key, value)
	}
	
	// Check for key=value format (legacy style, e.g., @prefix="My")
	if strings.Contains(annotation, "=") {
		parts := strings.SplitN(annotation, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		// Remove quotes if present
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || 
			(value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}
		
		return ec.setStringOption(key, value)
	}
	
	// Boolean flag without explicit value (defaults to true)
	return ec.setBoolOption(annotation, true)
}

// setBoolOption sets a boolean option in the EnumConfig.
func (ec *EnumConfig) setBoolOption(key string, value bool) error {
	switch key {
	case "noprefix":
		ec.NoPrefix = EnumConfigValue{Value: value, Valid: true}
	case "noiota":
		ec.NoIota = EnumConfigValue{Value: value, Valid: true}
	case "lower":
		ec.LowercaseLookup = EnumConfigValue{Value: value, Valid: true}
	case "nocase":
		ec.CaseInsensitive = EnumConfigValue{Value: value, Valid: true}
		if value {
			ec.LowercaseLookup = EnumConfigValue{Value: true, Valid: true} // nocase forces lower
		}
	case "marshal":
		ec.Marshal = EnumConfigValue{Value: value, Valid: true}
	case "sql":
		ec.SQL = EnumConfigValue{Value: value, Valid: true}
	case "sqlint":
		ec.SQLInt = EnumConfigValue{Value: value, Valid: true}
	case "flag":
		ec.Flag = EnumConfigValue{Value: value, Valid: true}
	case "names":
		ec.Names = EnumConfigValue{Value: value, Valid: true}
	case "values":
		ec.Values = EnumConfigValue{Value: value, Valid: true}
	case "nocamel":
		ec.LeaveSnakeCase = EnumConfigValue{Value: value, Valid: true}
	case "ptr":
		ec.Ptr = EnumConfigValue{Value: value, Valid: true}
	case "sqlnullint":
		ec.SQLNullInt = EnumConfigValue{Value: value, Valid: true}
	case "sqlnullstr":
		ec.SQLNullStr = EnumConfigValue{Value: value, Valid: true}
	case "mustparse":
		ec.MustParse = EnumConfigValue{Value: value, Valid: true}
	case "forcelower":
		ec.ForceLower = EnumConfigValue{Value: value, Valid: true}
	case "forceupper":
		ec.ForceUpper = EnumConfigValue{Value: value, Valid: true}
	case "nocomments":
		ec.NoComments = EnumConfigValue{Value: value, Valid: true}
	case "noparse":
		ec.NoParse = EnumConfigValue{Value: value, Valid: true}
	default:
		return fmt.Errorf("unknown annotation: @%s", key)
	}
	
	return nil
}

// setStringOption sets a string option in the EnumConfig.
func (ec *EnumConfig) setStringOption(key, value string) error {
	switch key {
	case "prefix":
		ec.Prefix = EnumConfigValue{Value: value, Valid: true}
	default:
		return fmt.Errorf("unknown annotation with value: @%s=%s", key, value)
	}
	return nil
}
