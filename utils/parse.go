package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// CastInput formats user input values safely
func CastInput(input string) interface{} {
	// If it's a number, cast to float64 (standard for JSON dynamics)
	if val, err := strconv.ParseFloat(input, 64); err == nil {
		return val
	}
	return input
}

// InterpolateVariables replaces tokens like {user_name} with context values
func InterpolateVariables(text string, ctx map[string]interface{}) string {
	re := regexp.MustCompile(`{(\w+)}`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		key := match[1 : len(match)-1] // strip brackets
		if val, ok := ctx[key]; ok {
			return fmt.Sprintf("%v", val)
		}
		return match
	})
}

// EvaluateExpression evaluates simple conditions like "context.user_age >= 18" safely
func EvaluateExpression(expression string, ctx map[string]interface{}) bool {
	expr := strings.TrimSpace(expression)
	expr = strings.ReplaceAll(expr, "context.", "")

	// Recognized operators
	operators := []string{">=", "<=", ">", "<", "==", "!="}
	var op string
	for _, o := range operators {
		if strings.Contains(expr, o) {
			op = o
			break
		}
	}
	if op == "" {
		return false
	}

	parts := strings.Split(expr, op)
	if len(parts) != 2 {
		return false
	}

	key := strings.TrimSpace(parts[0])
	compareStr := strings.TrimSpace(parts[1])

	val, exists := ctx[key]
	if !exists {
		return false
	}

	switch v := val.(type) {
	case float64:
		compareVal, err := strconv.ParseFloat(compareStr, 64)
		if err != nil {
			return false
		}
		switch op {
		case ">=":
			return v >= compareVal
		case "<=":
			return v <= compareVal
		case ">":
			return v > compareVal
		case "<":
			return v < compareVal
		case "==":
			return v == compareVal
		case "!=":
			return v != compareVal
		}
	case string:
		// Strip possible quotes in comparisons
		compareVal := strings.Trim(compareStr, "\"'")
		switch op {
		case "==":
			return v == compareVal
		case "!=":
			return v != compareVal
		}
	}

	return false
}
