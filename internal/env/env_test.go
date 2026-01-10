package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	tests := []struct {
		name          string
		envKey        string
		envValue      string
		setEnv        bool
		fallback      string
		expectedValue string
	}{
		{
			name:          "env_var_exists_returns_value",
			envKey:        "TEST_STRING",
			envValue:      "actual_value",
			setEnv:        true,
			fallback:      "fallback_value",
			expectedValue: "actual_value",
		},
		{
			name:          "env_var_not_exists_returns_fallback",
			envKey:        "NONEXISTENT_VAR",
			setEnv:        false,
			fallback:      "fallback_value",
			expectedValue: "fallback_value",
		},
		{
			name:          "env_var_empty_string",
			envKey:        "EMPTY_STRING",
			envValue:      "",
			setEnv:        true,
			fallback:      "fallback",
			expectedValue: "",
		},
		{
			name:          "env_var_with_spaces",
			envKey:        "SPACES_VALUE",
			envValue:      "  value with spaces  ",
			setEnv:        true,
			fallback:      "fallback",
			expectedValue: "  value with spaces  ",
		},
		{
			name:          "empty_fallback",
			envKey:        "MISSING",
			setEnv:        false,
			fallback:      "",
			expectedValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv(tt.envKey, tt.envValue)
			}

			result := GetString(tt.envKey, tt.fallback)

			assert.Equal(t, tt.expectedValue, result)
		})
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		name          string
		envKey        string
		envValue      string
		setEnv        bool
		fallback      int
		expectedValue int
	}{
		{
			name:          "valid_positive_int",
			envKey:        "TEST_INT",
			envValue:      "42",
			setEnv:        true,
			fallback:      10,
			expectedValue: 42,
		},
		{
			name:          "valid_negative_int",
			envKey:        "NEGATIVE_INT",
			envValue:      "-100",
			setEnv:        true,
			fallback:      10,
			expectedValue: -100,
		},
		{
			name:          "valid_zero",
			envKey:        "ZERO_INT",
			envValue:      "0",
			setEnv:        true,
			fallback:      10,
			expectedValue: 0,
		},
		{
			name:          "env_var_not_exists_returns_fallback",
			envKey:        "NONEXISTENT_INT",
			setEnv:        false,
			fallback:      999,
			expectedValue: 999,
		},
		{
			name:          "invalid_int_returns_fallback",
			envKey:        "INVALID_INT",
			envValue:      "not_a_number",
			setEnv:        true,
			fallback:      50,
			expectedValue: 50,
		},
		{
			name:          "empty_string_returns_fallback",
			envKey:        "EMPTY_INT",
			envValue:      "",
			setEnv:        true,
			fallback:      100,
			expectedValue: 100,
		},
		{
			name:          "float_string_returns_fallback",
			envKey:        "FLOAT_INT",
			envValue:      "3.14",
			setEnv:        true,
			fallback:      25,
			expectedValue: 25,
		},
		{
			name:          "int_with_spaces_returns_fallback",
			envKey:        "SPACES_INT",
			envValue:      " 42 ",
			setEnv:        true,
			fallback:      30,
			expectedValue: 30,
		},
		{
			name:          "very_large_int",
			envKey:        "LARGE_INT",
			envValue:      "2147483647",
			setEnv:        true,
			fallback:      1,
			expectedValue: 2147483647,
		},
		{
			name:          "negative_fallback",
			envKey:        "MISSING_NEG",
			setEnv:        false,
			fallback:      -999,
			expectedValue: -999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv(tt.envKey, tt.envValue)
			}

			result := GetInt(tt.envKey, tt.fallback)

			assert.Equal(t, tt.expectedValue, result)
		})
	}
}

func TestGetBool(t *testing.T) {
	tests := []struct {
		name          string
		envKey        string
		envValue      string
		setEnv        bool
		fallback      bool
		expectedValue bool
	}{
		// True values
		{
			name:          "true_lowercase",
			envKey:        "BOOL_TRUE_LC",
			envValue:      "true",
			setEnv:        true,
			fallback:      false,
			expectedValue: true,
		},
		{
			name:          "true_uppercase",
			envKey:        "BOOL_TRUE_UC",
			envValue:      "TRUE",
			setEnv:        true,
			fallback:      false,
			expectedValue: true,
		},
		{
			name:          "true_mixed_case",
			envKey:        "BOOL_TRUE_MIXED",
			envValue:      "True",
			setEnv:        true,
			fallback:      false,
			expectedValue: true,
		},
		{
			name:          "true_as_t",
			envKey:        "BOOL_T",
			envValue:      "t",
			setEnv:        true,
			fallback:      false,
			expectedValue: true,
		},
		{
			name:          "true_as_T",
			envKey:        "BOOL_T_UC",
			envValue:      "T",
			setEnv:        true,
			fallback:      false,
			expectedValue: true,
		},
		{
			name:          "true_as_1",
			envKey:        "BOOL_1",
			envValue:      "1",
			setEnv:        true,
			fallback:      false,
			expectedValue: true,
		},
		// False values
		{
			name:          "false_lowercase",
			envKey:        "BOOL_FALSE_LC",
			envValue:      "false",
			setEnv:        true,
			fallback:      true,
			expectedValue: false,
		},
		{
			name:          "false_uppercase",
			envKey:        "BOOL_FALSE_UC",
			envValue:      "FALSE",
			setEnv:        true,
			fallback:      true,
			expectedValue: false,
		},
		{
			name:          "false_mixed_case",
			envKey:        "BOOL_FALSE_MIXED",
			envValue:      "False",
			setEnv:        true,
			fallback:      true,
			expectedValue: false,
		},
		{
			name:          "false_as_f",
			envKey:        "BOOL_F",
			envValue:      "f",
			setEnv:        true,
			fallback:      true,
			expectedValue: false,
		},
		{
			name:          "false_as_F",
			envKey:        "BOOL_F_UC",
			envValue:      "F",
			setEnv:        true,
			fallback:      true,
			expectedValue: false,
		},
		{
			name:          "false_as_0",
			envKey:        "BOOL_0",
			envValue:      "0",
			setEnv:        true,
			fallback:      true,
			expectedValue: false,
		},
		// Not exists
		{
			name:          "env_var_not_exists_returns_fallback_true",
			envKey:        "NONEXISTENT_BOOL_TRUE",
			setEnv:        false,
			fallback:      true,
			expectedValue: true,
		},
		{
			name:          "env_var_not_exists_returns_fallback_false",
			envKey:        "NONEXISTENT_BOOL_FALSE",
			setEnv:        false,
			fallback:      false,
			expectedValue: false,
		},
		// Invalid values
		{
			name:          "invalid_bool_returns_fallback",
			envKey:        "INVALID_BOOL",
			envValue:      "not_a_bool",
			setEnv:        true,
			fallback:      true,
			expectedValue: true,
		},
		{
			name:          "empty_string_returns_fallback",
			envKey:        "EMPTY_BOOL",
			envValue:      "",
			setEnv:        true,
			fallback:      false,
			expectedValue: false,
		},
		{
			name:          "numeric_2_returns_fallback",
			envKey:        "BOOL_2",
			envValue:      "2",
			setEnv:        true,
			fallback:      true,
			expectedValue: true,
		},
		{
			name:          "yes_returns_fallback",
			envKey:        "BOOL_YES",
			envValue:      "yes",
			setEnv:        true,
			fallback:      false,
			expectedValue: false,
		},
		{
			name:          "no_returns_fallback",
			envKey:        "BOOL_NO",
			envValue:      "no",
			setEnv:        true,
			fallback:      true,
			expectedValue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv(tt.envKey, tt.envValue)
			}

			result := GetBool(tt.envKey, tt.fallback)

			assert.Equal(t, tt.expectedValue, result)
		})
	}
}

func TestGetString_MultipleKeys(t *testing.T) {
	t.Setenv("KEY1", "value1")
	t.Setenv("KEY2", "value2")
	t.Setenv("KEY3", "value3")

	assert.Equal(t, "value1", GetString("KEY1", "fallback"))
	assert.Equal(t, "value2", GetString("KEY2", "fallback"))
	assert.Equal(t, "value3", GetString("KEY3", "fallback"))
	assert.Equal(t, "fallback", GetString("KEY4", "fallback"))
}

func TestGetInt_MultipleKeys(t *testing.T) {
	t.Setenv("NUM1", "10")
	t.Setenv("NUM2", "20")
	t.Setenv("NUM3", "30")

	assert.Equal(t, 10, GetInt("NUM1", 0))
	assert.Equal(t, 20, GetInt("NUM2", 0))
	assert.Equal(t, 30, GetInt("NUM3", 0))
	assert.Equal(t, 99, GetInt("NUM4", 99))
}

func TestGetBool_MultipleKeys(t *testing.T) {
	t.Setenv("BOOL1", "true")
	t.Setenv("BOOL2", "false")
	t.Setenv("BOOL3", "1")
	t.Setenv("BOOL4", "0")

	assert.Equal(t, true, GetBool("BOOL1", false))
	assert.Equal(t, false, GetBool("BOOL2", true))
	assert.Equal(t, true, GetBool("BOOL3", false))
	assert.Equal(t, false, GetBool("BOOL4", true))
	assert.Equal(t, true, GetBool("BOOL5", true))
}
