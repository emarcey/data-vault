package server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseStringValueErrors(t *testing.T) {
	var tests = []struct {
		op        string
		vars      map[string]string
		paramName string
	}{
		{
			op:   "empty map",
			vars: map[string]string{},
		},
		{
			op: "not found",
			vars: map[string]string{
				"anything2": "something",
			},
		},
	}

	for _, given := range tests {
		t.Run(fmt.Sprintf("parseStringValue - Errors - %v", given.op), func(t *testing.T) {
			result, err := parseStringValue(given.op, given.vars, "anything")

			require.NotNil(t, err, "no error in parseStringValue: %v", err)
			require.Equal(t, result, "", "Result, %v, does not equal expected, ''", result)
		})
	}
}

func TestParseStringValueSuccess(t *testing.T) {
	var tests = []struct {
		op        string
		vars      map[string]string
		paramName string
		expected  string
	}{
		{
			op: "found",
			vars: map[string]string{
				"anything": "something",
			},
			expected: "something",
		},
	}

	for _, given := range tests {
		t.Run(fmt.Sprintf("parseStringValue - Success - %v", given.op), func(t *testing.T) {
			result, err := parseStringValue(given.op, given.vars, "anything")

			require.Nil(t, err, "error in parseStringValue: %v", err)
			require.Equal(t, result, given.expected, "Result, %v, does not equal expected, ''", result)
		})
	}
}

func TestParseIntegerUrlParamErrors(t *testing.T) {
	var tests = []struct {
		op           string
		urlParams    map[string][]string
		paramName    string
		defaultValue int
	}{
		{
			op: "too few vals",
			urlParams: map[string][]string{
				"anything": []string{},
			},
		},
		{
			op: "too many vals",
			urlParams: map[string][]string{
				"anything": []string{"abc", "def"},
			},
		},
		{
			op: "not int",
			urlParams: map[string][]string{
				"anything": []string{"25zzz"},
			},
		},
		{
			op: "not int",
			urlParams: map[string][]string{
				"anything": []string{"25.222"},
			},
		},
		{
			op: "negative",
			urlParams: map[string][]string{
				"anything": []string{"-25"},
			},
		},
	}

	for _, given := range tests {
		t.Run(fmt.Sprintf("parseIntegerUrlParam - Errors - %v", given.op), func(t *testing.T) {
			result, err := parseIntegerUrlParam(given.op, given.urlParams, "anything", 10)

			require.NotNil(t, err, "no error in parseIntegerUrlParam: %v", err)
			require.Equal(t, result, -1, "Result, %v, does not equal expected, -1", result)
		})
	}
}

func TestParseIntegerUrlParamSuccess(t *testing.T) {
	var tests = []struct {
		op           string
		urlParams    map[string][]string
		paramName    string
		defaultValue int
		expected     int
	}{
		{
			op:           "empty map",
			urlParams:    map[string][]string{},
			paramName:    "anything",
			defaultValue: 10,
			expected:     10,
		},
		{
			op: "missing val",
			urlParams: map[string][]string{
				"anything2": []string{"abc", "def"},
			},
			paramName:    "anything",
			defaultValue: 10,
			expected:     10,
		},
		{
			op: "found",
			urlParams: map[string][]string{
				"anything": []string{"25"},
			},
			paramName:    "anything",
			defaultValue: 10,
			expected:     25,
		},
		{
			op: "found - zero",
			urlParams: map[string][]string{
				"anything": []string{"0"},
			},
			paramName:    "anything",
			defaultValue: 10,
			expected:     0,
		},
	}

	for _, given := range tests {
		t.Run(fmt.Sprintf("parseIntegerUrlParam - Success - %v", given.op), func(t *testing.T) {
			result, err := parseIntegerUrlParam(given.op, given.urlParams, given.paramName, given.defaultValue)

			require.Nil(t, err, "error in parseIntegerUrlParam: %v", err)
			require.Equal(t, result, given.expected, "Result, %v, does not equal expected, %v", result, given.expected)
		})
	}
}
