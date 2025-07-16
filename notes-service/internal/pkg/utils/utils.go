package utils

import (
	_"fmt"
	"strings"
)

type Field struct {
	Name  string
	Value interface{}
}

type ChangedColumns struct {
	Title       string
	Priority    string
	Description string
	Category    string
	CreatedAt   string
	CompletedAt *string
}

func (c ChangedColumns) StringifyNotEmptyFields() string {
	var parts []string
	if c.Title != "" {
		parts = append(parts, "title")
	}
	if c.Priority != "" {
		parts = append(parts, "priority")
	}
	if c.Description != "" {
		parts = append(parts, "description")
	}
	if c.Category != "" {
		parts = append(parts, "category")
	}
	if c.CreatedAt != "" {
		parts = append(parts, "createdat")
	}
	if c.CompletedAt != nil {
		parts = append(parts, "completedat")
	}
	return strings.Join(parts, ", ")
}

func FilterFields(fields []Field) ChangedColumns {
	result := ChangedColumns{}
	for _, f := range fields {
		switch f.Name {
		case "title":
			if val, ok := f.Value.(string); ok {
				result.Title = val
			}
		case "priority":
			if val, ok := f.Value.(string); ok {
				result.Priority = val
			}
		case "description":
			if val, ok := f.Value.(string); ok {
				result.Description = val
			}
		case "category":
			if val, ok := f.Value.(string); ok {
				result.Category = val
			}
		case "createdat":
			if val, ok := f.Value.(string); ok {
				result.CreatedAt = val
			}
		case "completedat":
			if val, ok := f.Value.(*string); ok {
				result.CompletedAt = val
			}
		}
	}
	return result
}