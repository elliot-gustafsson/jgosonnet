package tests

import (
	"testing"
)

func Test_jsonnet_arith_bool(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/arith_bool.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "arith_bool.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_arith_float(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/arith_float.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "arith_float.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_arith_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/arith_string.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "arith_string.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_array(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/array.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "array.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_array_comparison(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/array_comparison.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/array_comparison.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "array_comparison.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_array_comparison2(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/array_comparison2.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/array_comparison2.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "array_comparison2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_assert(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/assert.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "assert.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_binary(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/binary.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "binary.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_comments(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/comments.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "comments.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_condition(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/condition.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "condition.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_digitsep(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/digitsep.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/digitsep.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "digitsep.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_dos_line_endings(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/dos_line_endings.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/dos_line_endings.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "dos_line_endings.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_error_01(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.01.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.01.jsonnet.golden
	expected := `RUNTIME ERROR: foo
	error.01.jsonnet:17:29-40	function <bananas>
	error.01.jsonnet:18:29-39	function <oranges>
	error.01.jsonnet:19:28-38	function <apples>
	error.01.jsonnet:20:1-10	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.01.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_02(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.02.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.02.jsonnet.golden
	expected := `RUNTIME ERROR: Foo.
	error.02.jsonnet:17:1-13	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.02.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_03(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.03.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.03.jsonnet.golden
	expected := `RUNTIME ERROR: foo
	error.03.jsonnet:17:21-32	object <ErrObj>
	error.03.jsonnet:18:1-9	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.03.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_04(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.04.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.04.jsonnet.golden
	expected := `RUNTIME ERROR: foo
	error.04.jsonnet:17:21-32	object <ErrObj>
	During manifestation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.04.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_05(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.05.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.05.jsonnet.golden
	expected := `RUNTIME ERROR: foo
	error.05.jsonnet:17:21-32	object <ErrObj>
	During manifestation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.05.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_06(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.06.jsonnet")
	expected := `RUNTIME ERROR: Division by zero.
	error.06.jsonnet:18:22-25	function <f>
	error.06.jsonnet:19:1-4	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.06.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_07(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.07.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.07.jsonnet.golden
	expected := `RUNTIME ERROR: sarcasm
	error.07.jsonnet:17:29-35	function <third>
	error.07.jsonnet:19:1-6	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.07.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_08(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.08.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type object, expected string
	error.08.jsonnet:18:1-8	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.08.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_args_commafodder(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.args_commafodder.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.args_commafodder.jsonnet.golden
	expected := `RUNTIME ERROR: error.args_commafodder.jsonnet:1:1-4 Unknown variable: foo
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.args_commafodder.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

// Note: go-jsonnet allows this, it just trucates the fraction
// func Test_jsonnet_error_array_fractional_index(t *testing.T) {
// 	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.array_fractional_index.jsonnet")
// 	expected := `RUNTIME ERROR: array index was not integer: 1.5
// 	error.array_fractional_index.jsonnet:17:1-15	$
// 	During evaluation`
// 	extVars := map[string]string{"var1": "test"}
// 	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
// 	runTest(t, TestConfig{
// 		WorkDir:     "resources/jsonnet-cpp/test_suite",
// 		File:        "error.array_fractional_index.jsonnet",
// 		Snippet:     snippet,
// 		ExpectedErr: expected,
// 		ExtVars:     extVars,
// 		ExtCodes:    extCodes,
// 	})
// }

func Test_jsonnet_error_array_index_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.array_index_string.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type string, expected number
	error.array_index_string.jsonnet:17:1-14	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.array_index_string.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_array_large_index(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.array_large_index.jsonnet")
	expected := `RUNTIME ERROR: Index -9223372036854775808 out of bounds, not within [0, 3)
	error.array_large_index.jsonnet:17:1-26	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.array_large_index.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_array_recursive_manifest(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.array_recursive_manifest.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.array_recursive_manifest.jsonnet.golden
	expected := `RUNTIME ERROR: max manifest depth exceeded, possible infinite recursion	During manifestation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.array_recursive_manifest.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
		MaxStack:    500,
	})
}

func Test_jsonnet_error_assert_fail1(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.assert.fail1.jsonnet")
	expected := `RUNTIME ERROR: Assertion failed
	error.assert.fail1.jsonnet:(20:1)-(22:5)
	error.assert.fail1.jsonnet:(18:1)-(22:5)	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.assert.fail1.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_assert_fail2(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.assert.fail2.jsonnet")
	expected := `RUNTIME ERROR: foo was not equal to bar
	error.assert.fail2.jsonnet:(20:1)-(22:5)
	error.assert.fail2.jsonnet:(18:1)-(22:5)	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.assert.fail2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_assert_equal_obj(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.assert_equal_obj.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.assert_equal_obj.jsonnet.golden
	expected := `RUNTIME ERROR: Assertion failed. {"a": 1} != {"b": 1}
	error.assert_equal_obj.jsonnet:17:1-36	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.assert_equal_obj.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_assert_equal_str(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.assert_equal_str.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.assert_equal_str.jsonnet.golden
	expected := `RUNTIME ERROR: Assertion failed. "one\ntwo" != "three\nfour\n"
	error.assert_equal_str.jsonnet:17:1-45	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.assert_equal_str.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_comprehension_spec_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.comprehension_spec_object.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type object, expected array
	error.comprehension_spec_object.jsonnet:17:1-22
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.comprehension_spec_object.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_comprehension_spec_object2(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.comprehension_spec_object2.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type object, expected array
	error.comprehension_spec_object2.jsonnet:17:1-32
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.comprehension_spec_object2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_computed_field_scope(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.computed_field_scope.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.computed_field_scope.jsonnet.golden
	expected := `RUNTIME ERROR: error.computed_field_scope.jsonnet:17:21-22 Unknown variable: x
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.computed_field_scope.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_decodeUTF8_float(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.decodeUTF8_float.jsonnet")
	expected := `RUNTIME ERROR: Expected an integer, but got 17.5
	error.decodeUTF8_float.jsonnet:1:1-23	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.decodeUTF8_float.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_decodeUTF8_nan(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.decodeUTF8_nan.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type string, expected number
	error.decodeUTF8_nan.jsonnet:1:1-24	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.decodeUTF8_nan.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_divide_zero(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.divide_zero.jsonnet")
	expected := `RUNTIME ERROR: Division by zero.
	error.divide_zero.jsonnet:17:1-8	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.divide_zero.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_equality_function(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.equality_function.jsonnet")
	expected := `comparing types function is not supported
	error.equality_function.jsonnet:17:1-33	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.equality_function.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_field_not_exist(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.field_not_exist.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.field_not_exist.jsonnet.golden
	expected := `RUNTIME ERROR: Field does not exist: y
	error.field_not_exist.jsonnet:17:1-11	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.field_not_exist.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_flatMap_array_typecheck(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.flatMap_array_typecheck.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type string, expected array
	error.flatMap_array_typecheck.jsonnet:1:1-44	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.flatMap_array_typecheck.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_flatMap_seq_typecheck(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.flatMap_seq_typecheck.jsonnet")
	expected := `RUNTIME ERROR: std.flatMap second param must be array / string, got object
	error.flatMap_seq_typecheck.jsonnet:1:1-51	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.flatMap_seq_typecheck.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_flatMap_string_typecheck(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.flatMap_string_typecheck.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type number, expected string
	error.flatMap_string_typecheck.jsonnet:1:1-49	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.flatMap_string_typecheck.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_format_too_few_values(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.format.too_few_values.jsonnet")
	expected := `not enough arguments for format string
	error.format.too_few_values.jsonnet:1:1-18
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.format.too_few_values.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_function_duplicate_arg(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.function_duplicate_arg.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.function_duplicate_arg.jsonnet.golden
	expected := `RUNTIME ERROR: Argument x already provided
	error.function_duplicate_arg.jsonnet:17:1-29	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.function_duplicate_arg.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

// Note: This is valid in go-jsonnet, so this behaviour is skipped
// func Test_jsonnet_error_function_duplicate_param(t *testing.T) {
// 	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.function_duplicate_param.jsonnet")
// 	expected := `arg (0) with no default arg had no value passed	During evaluation`
// 	extVars := map[string]string{"var1": "test"}
// 	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
// 	runTest(t, TestConfig{
// 		WorkDir:     "resources/jsonnet-cpp/test_suite",
// 		File:        "error.function_duplicate_param.jsonnet",
// 		Snippet:     snippet,
// 		ExpectedErr: expected,
// 		ExtVars:     extVars,
// 		ExtCodes:    extCodes,
// 	})
// }

func Test_jsonnet_error_function_infinite_default(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.function_infinite_default.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.function_infinite_default.jsonnet.golden
	expected := `RUNTIME ERROR: infinite loop detected
	error.function_infinite_default.jsonnet:17:20-21	function <anonymous>
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}

	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.function_infinite_default.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_function_no_default_arg(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.function_no_default_arg.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.function_no_default_arg.jsonnet.golden
	expected := `RUNTIME ERROR: Missing argument: b	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.function_no_default_arg.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_function_too_many_args(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.function_too_many_args.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.function_too_many_args.jsonnet.golden
	expected := `RUNTIME ERROR: function expected 2 positional argument(s), but got 3
	error.function_too_many_args.jsonnet:19:1-13	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.function_too_many_args.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_import_empty(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.import_empty.jsonnet")
	expected := `RUNTIME ERROR: couldn't open import "": the empty string is not a valid filename
	error.import_empty.jsonnet:17:1-10	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.import_empty.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_import_static_check_failure(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.import_static-check-failure.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.import_static-check-failure.jsonnet.golden
	expected := `RUNTIME ERROR: lib/static_check_failure.jsonnet:2:1-2 Unknown variable: x
	error.import_static-check-failure.jsonnet:1:1-42	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.import_static-check-failure.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_import_syntax_error(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.import_syntax-error.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.import_syntax-error.jsonnet.golden
	expected := `RUNTIME ERROR: lib/syntax_error.jsonnet:1:1 Unterminated String
	error.import_syntax-error.jsonnet:1:1-34	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.import_syntax-error.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_inside_equals_array(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.inside_equals_array.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.inside_equals_array.jsonnet.golden
	expected := `RUNTIME ERROR: foobar
	error.inside_equals_array.jsonnet:19:1-7	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.inside_equals_array.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_inside_equals_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.inside_equals_object.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.inside_equals_object.jsonnet.golden
	expected := `RUNTIME ERROR: foobar
	error.inside_equals_object.jsonnet:18:22-36	object <B>
	error.inside_equals_object.jsonnet:19:1-7	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.inside_equals_object.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_inside_tostring_array(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.inside_tostring_array.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.inside_tostring_array.jsonnet.golden
	expected := `RUNTIME ERROR: foobar
	error.inside_tostring_array.jsonnet:17:1-28	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.inside_tostring_array.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_inside_tostring_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.inside_tostring_object.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.inside_tostring_object.jsonnet.golden
	expected := `RUNTIME ERROR: foobar
	error.inside_tostring_object.jsonnet:17:12-26	object <anonymous>
	error.inside_tostring_object.jsonnet:17:1-33	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.inside_tostring_object.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_integer_conversion(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.integer_conversion.jsonnet")
	expected := `RUNTIME ERROR: Bitwise operator argument 9.007199254740992e+15 outside of range [-9007199254740991, 9007199254740991]
	error.integer_conversion.jsonnet:3:1-16	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.integer_conversion.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_integer_left_shift(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.integer_left_shift.jsonnet")
	expected := `RUNTIME ERROR: Bitwise operator argument 4.611686018427388e+18 outside of range [-9007199254740991, 9007199254740991]
	error.integer_left_shift.jsonnet:3:1-17	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.integer_left_shift.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

// Note: This works in go-jsonnet, so that behaviour also works here
func Test_jsonnet_error_integer_left_shift_runtime(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.integer_left_shift_runtime.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.integer_left_shift_runtime.jsonnet.golden
	// expected := `RUNTIME ERROR: numeric value outside safe integer range for bitwise operation.`
	expected := "-9223372036854775808\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "error.integer_left_shift_runtime.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_error_invariant_avoid_output_change(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.invariant.avoid_output_change.jsonnet")
	expected := `RUNTIME ERROR: Object assertion failed.
	error.invariant.avoid_output_change.jsonnet:18:3-24
	During manifestation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.invariant.avoid_output_change.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_invariant_equality(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.invariant.equality.jsonnet")
	expected := `RUNTIME ERROR: Object assertion failed.
	error.invariant.equality.jsonnet:17:3-15
	error.invariant.equality.jsonnet:17:1-35	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.invariant.equality.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_invariant_option(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.invariant.option.jsonnet")
	expected := `RUNTIME ERROR: Option "d" not in ["a", "b", "c"].
	error.invariant.option.jsonnet:(19:3)-(20:59)
	During manifestation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.invariant.option.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_invariant_simple(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.invariant.simple.jsonnet")
	expected := `RUNTIME ERROR: Object assertion failed.
	error.invariant.simple.jsonnet:18:3-15
	During manifestation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.invariant.simple.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_invariant_simple2(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.invariant.simple2.jsonnet")
	expected := `RUNTIME ERROR: my error message
	error.invariant.simple2.jsonnet:18:3-37
	During manifestation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.invariant.simple2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_invariant_simple3(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.invariant.simple3.jsonnet")
	expected := `RUNTIME ERROR: my error message
	error.invariant.simple3.jsonnet:18:10-34	object <anonymous>
	error.invariant.simple3.jsonnet:18:3-34
	During manifestation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.invariant.simple3.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_manifest_toml_null_value(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.manifest_toml_null_value.jsonnet")
	expected := `unsupported value type null for toml manifestation
	error.manifest_toml_null_value.jsonnet:17:1-54	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.manifest_toml_null_value.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_manifest_toml_wrong_type(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.manifest_toml_wrong_type.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type array, expected object
	error.manifest_toml_wrong_type.jsonnet:17:1-29	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.manifest_toml_wrong_type.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_negative_shfit(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.negative_shfit.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.negative_shfit.jsonnet.golden
	expected := `RUNTIME ERROR: Shift by negative exponent.
	error.negative_shfit.jsonnet:1:1-13	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.negative_shfit.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_obj_assert_fail1(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.obj_assert.fail1.jsonnet")
	expected := `RUNTIME ERROR: Object assertion failed.
	error.obj_assert.fail1.jsonnet:20:16-29
	error.obj_assert.fail1.jsonnet:20:1-49	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.obj_assert.fail1.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_obj_assert_fail2(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.obj_assert.fail2.jsonnet")
	expected := `RUNTIME ERROR: foo was not equal to bar
	error.obj_assert.fail2.jsonnet:20:16-65
	error.obj_assert.fail2.jsonnet:20:1-85	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.obj_assert.fail2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_obj_recursive(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.obj_recursive.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.obj_recursive.jsonnet.golden
	expected := `RUNTIME ERROR: max manifest depth exceeded, possible infinite recursion	During manifestation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.obj_recursive.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
		MaxStack:    500,
	})
}

func Test_jsonnet_error_obj_recursive_manifest(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.obj_recursive_manifest.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.obj_recursive_manifest.jsonnet.golden
	expected := `RUNTIME ERROR: max manifest depth exceeded, possible infinite recursion	During manifestation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.obj_recursive_manifest.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
		MaxStack:    500,
	})
}

func Test_jsonnet_error_overflow(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.overflow.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.overflow.jsonnet.golden
	expected := `(*ast.LiteralNumber) failed to parse float val (1e309), err: strconv.ParseFloat: parsing "1e309": value out of range
	error.overflow.jsonnet:17:1-6	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.overflow.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_overflow2(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.overflow2.jsonnet")
	expected := `RUNTIME ERROR: Overflow
	error.overflow2.jsonnet:17:1-19	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.overflow2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_overflow3(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.overflow3.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.overflow3.jsonnet.golden
	expected := `(*ast.LiteralNumber) failed to parse float val (1e309), err: strconv.ParseFloat: parsing "1e309": value out of range
	error.overflow3.jsonnet:17:1-6	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.overflow3.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_array_comma(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.array_comma.jsonnet")
	expected := `RUNTIME ERROR: error.parse.array_comma.jsonnet:17:7-8 Expected a comma before next array element
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.array_comma.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_deep_array_nesting(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.deep_array_nesting.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.parse.deep_array_nesting.jsonnet.golden
	expected := `RUNTIME ERROR: max manifest depth exceeded, possible infinite recursion	During manifestation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.deep_array_nesting.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
		MaxStack:    500,
	})
}

func Test_jsonnet_error_parse_function_arg_positional_after_named(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.function_arg_positional_after_named.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.parse.function_arg_positional_after_named.jsonnet.golden
	expected := `RUNTIME ERROR: error.parse.function_arg_positional_after_named.jsonnet:19:10-11 Positional argument after a named argument is not allowed
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.function_arg_positional_after_named.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_import_not_literal(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.import_not_literal.jsonnet")
	expected := `RUNTIME ERROR: error.parse.import_not_literal.jsonnet:17:8-28 Computed imports are not allowed
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.import_not_literal.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_import_text_block(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.import_text_block.jsonnet")
	expected := `RUNTIME ERROR: error.parse.import_text_block.jsonnet:(17:8)-(20:4) Block string literals not allowed in imports
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.import_text_block.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_index_unterminated(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.index_unterminated.jsonnet")
	expected := `RUNTIME ERROR: error.parse.index_unterminated.jsonnet:18:1 Unexpected end of file
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.index_unterminated.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_method_plus(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.method_plus.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.parse.method_plus.jsonnet.golden
	expected := `RUNTIME ERROR: error.parse.method_plus.jsonnet:17:15-16 Cannot use +: syntax sugar in a method: a
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.method_plus.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_object_comma(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.object_comma.jsonnet")
	expected := `RUNTIME ERROR: error.parse.object_comma.jsonnet:17:11-12 Expected a comma before next field
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.object_comma.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_object_comprehension_local_clash(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.object_comprehension_local_clash.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.parse.object_comprehension_local_clash.jsonnet.golden
	expected := `RUNTIME ERROR: error.parse.object_comprehension_local_clash.jsonnet:17:21-22 Duplicate local var: x
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.object_comprehension_local_clash.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_object_local_clash(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.object_local_clash.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.parse.object_local_clash.jsonnet.golden
	expected := `RUNTIME ERROR: error.parse.object_local_clash.jsonnet:17:21-22 Duplicate local var: x
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.object_local_clash.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_self_in_computed_field(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.self_in_computed_field.jsonnet")
	expected := `RUNTIME ERROR: error.parse.self_in_computed_field.jsonnet:17:15-19 Unexpected: "self" while parsing field definition
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.self_in_computed_field.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_static_error_bad_number(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.static_error_bad_number.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.parse.static_error_bad_number.jsonnet.golden
	expected := `RUNTIME ERROR: error.parse.static_error_bad_number.jsonnet:17:1-2 Unexpected: "." while parsing terminal
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.static_error_bad_number.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_string_invalid_escape(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.string.invalid_escape.jsonnet")
	expected := `RUNTIME ERROR: error.parse.string.invalid_escape.jsonnet:17:1-5 error.parse.string.invalid_escape.jsonnet:17:1-5 Unknown escape sequence in string literal: \o
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.string.invalid_escape.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_string_invalid_escape_unicode_non_hex(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.string.invalid_escape_unicode_non_hex.jsonnet")
	expected := `RUNTIME ERROR: error.parse.string.invalid_escape_unicode_non_hex.jsonnet:17:1-9 error.parse.string.invalid_escape_unicode_non_hex.jsonnet:17:1-9 Unicode escape sequence was malformed: \u00
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.string.invalid_escape_unicode_non_hex.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_string_invalid_escape_unicode_short(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.string.invalid_escape_unicode_short.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.parse.string.invalid_escape_unicode_short.jsonnet.golden
	expected := `RUNTIME ERROR: error.parse.string.invalid_escape_unicode_short.jsonnet:17:1 Unterminated String
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.string.invalid_escape_unicode_short.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_string_invalid_escape_unicode_short2(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.string.invalid_escape_unicode_short2.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.parse.string.invalid_escape_unicode_short2.jsonnet.golden
	expected := `RUNTIME ERROR: error.parse.string.invalid_escape_unicode_short2.jsonnet:17:1-8 error.parse.string.invalid_escape_unicode_short2.jsonnet:17:1-8 Truncated unicode escape sequence in string literal.
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.string.invalid_escape_unicode_short2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_string_invalid_escape_unicode_short3(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.string.invalid_escape_unicode_short3.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.parse.string.invalid_escape_unicode_short3.jsonnet.golden
	expected := `RUNTIME ERROR: error.parse.string.invalid_escape_unicode_short3.jsonnet:17:1 Unterminated String
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.string.invalid_escape_unicode_short3.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_string_unfinished(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.string.unfinished.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.parse.string.unfinished.jsonnet.golden
	expected := `RUNTIME ERROR: error.parse.string.unfinished.jsonnet:17:1 Unterminated String
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.string.unfinished.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_string_unfinished2(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.string.unfinished2.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.parse.string.unfinished2.jsonnet.golden
	expected := `RUNTIME ERROR: error.parse.string.unfinished2.jsonnet:17:1 Unterminated String
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.string.unfinished2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_string_multi_no_newline(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.string_multi_no_newline.jsonnet")
	expected := `RUNTIME ERROR: error.parse.string_multi_no_newline.jsonnet:17:1 Text block requires new line after |||.
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.string_multi_no_newline.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_text_block_bad_whitespace(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.text_block_bad_whitespace.jsonnet")
	expected := `RUNTIME ERROR: error.parse.text_block_bad_whitespace.jsonnet:17:1 Text block not terminated with |||
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.text_block_bad_whitespace.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_text_block_eof(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.text_block_eof.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.parse.text_block_eof.jsonnet.golden
	expected := `RUNTIME ERROR: error.parse.text_block_eof.jsonnet:17:1 Unexpected EOF
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.text_block_eof.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_text_block_indent_spaces(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.text_block_indent_spaces.jsonnet")
	expected := `RUNTIME ERROR: error.parse.text_block_indent_spaces.jsonnet:17:1 Text block not terminated with |||
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.text_block_indent_spaces.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_text_block_not_terminated(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse.text_block_not_terminated.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.parse.text_block_not_terminated.jsonnet.golden
	expected := `RUNTIME ERROR: error.parse.text_block_not_terminated.jsonnet:17:1 Text block not terminated with |||
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse.text_block_not_terminated.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_parse_json(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.parse_json.jsonnet")
	expected := `failed to parse json, err: invalid character 'b' looking for beginning of value
	error.parse_json.jsonnet:1:1-29	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.parse_json.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_recursive_function_nonterm(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.recursive_function_nonterm.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.recursive_function_nonterm.jsonnet.golden
	expected := `RUNTIME ERROR: max stack frames exceeded.
	error.recursive_function_nonterm.jsonnet:18:3-7	function <f>
	error.recursive_function_nonterm.jsonnet:20:1-6	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}

	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.recursive_function_nonterm.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_recursive_import(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.recursive_import.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.recursive_import.jsonnet.golden
	expected := `RUNTIME ERROR: infinite loop detected
	error.recursive_import.jsonnet:17:15-54	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.recursive_import.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_recursive_object_non_term(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.recursive_object_non_term.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.recursive_object_non_term.jsonnet.golden
	expected := `RUNTIME ERROR: max stack frames exceeded.
	error.recursive_object_non_term.jsonnet:20:9-15	object <Fib>
	error.recursive_object_non_term.jsonnet:23:1-18	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.recursive_object_non_term.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_sanity(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.sanity.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.sanity.jsonnet.golden
	expected := `RUNTIME ERROR: Assertion failed. 1 != 2
	error.sanity.jsonnet:17:1-22	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.sanity.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_static_error_self(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.static_error_self.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.static_error_self.jsonnet.golden
	expected := `RUNTIME ERROR: error.static_error_self.jsonnet:17:2-6 Can't use self outside of an object.
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.static_error_self.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_static_error_super(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.static_error_super.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.static_error_super.jsonnet.golden
	expected := `RUNTIME ERROR: error.static_error_super.jsonnet:17:2-7 Can't use super outside of an object.
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.static_error_super.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_static_error_var_not_exist(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.static_error_var_not_exist.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.static_error_var_not_exist.jsonnet.golden
	expected := `RUNTIME ERROR: error.static_error_var_not_exist.jsonnet:17:16-20 Unknown variable: tmp2
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.static_error_var_not_exist.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_std_join_types1(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.std_join_types1.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type array, expected string
	error.std_join_types1.jsonnet:17:1-26	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.std_join_types1.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_std_join_types2(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.std_join_types2.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type string, expected array
	error.std_join_types2.jsonnet:17:1-31	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.std_join_types2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

// Note: go-jsonnet allows this, just returns an empty array
// func Test_jsonnet_error_std_makeArray_negative(t *testing.T) {
// 	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.std_makeArray_negative.jsonnet")
// 	expected := `RUNTIME ERROR: makeArray requires size >= 0, got -10
// 	error.std_makeArray_negative.jsonnet:17:1-37	$
// 	During evaluation`
// 	extVars := map[string]string{"var1": "test"}
// 	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
// 	runTest(t, TestConfig{
// 		WorkDir:     "resources/jsonnet-cpp/test_suite",
// 		File:        "error.std_makeArray_negative.jsonnet",
// 		Snippet:     snippet,
// 		ExpectedErr: expected,
// 		ExtVars:     extVars,
// 		ExtCodes:    extCodes,
// 	})
// }

func Test_jsonnet_error_std_maxArray(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.std_maxArray.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.std_maxArray.jsonnet.golden
	expected := `Expected at least one element in array. Got none
	error.std_maxArray.jsonnet:1:1-17	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.std_maxArray.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_std_minArray(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.std_minArray.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.std_minArray.jsonnet.golden
	expected := `Expected at least one element in array. Got none
	error.std_minArray.jsonnet:1:1-17	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.std_minArray.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_std_parseJson_nodigitsep(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.std_parseJson.nodigitsep.jsonnet")
	expected := `failed to parse json, err: invalid character '_' after top-level value
	error.std_parseJson.nodigitsep.jsonnet:1:1-25	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.std_parseJson.nodigitsep.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_std_parseYaml1(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.std_parseYaml1.jsonnet")
	expected := `failed to parse yaml, err: yaml: mapping values are not allowed in this context
	error.std_parseYaml1.jsonnet:1:1-23	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.std_parseYaml1.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_top_level_func(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.top_level_func.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.top_level_func.jsonnet.golden
	expected := "RUNTIME ERROR: Missing argument: name	During evaluation"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.top_level_func.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_trace_one_param(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.trace_one_param.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.trace_one_param.jsonnet.golden
	expected := `RUNTIME ERROR: Missing argument: rest
	error.trace_one_param.jsonnet:19:6-7	object <anonymous>
	During manifestation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.trace_one_param.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_trace_three_param(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.trace_three_param.jsonnet")
	expected := `RUNTIME ERROR: function expected 2 positional argument(s), but got 3
	error.trace_three_param.jsonnet:19:6-7	object <anonymous>
	During manifestation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.trace_three_param.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_trace_two_param(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.trace_two_param.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type number, expected string
	error.trace_two_param.jsonnet:19:6-7	object <anonymous>
	During manifestation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.trace_two_param.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_trace_zero_param(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.trace_zero_param.jsonnet")
	// Expected file: resources/jsonnet-cpp/test_suite/error.trace_zero_param.jsonnet.golden
	expected := `RUNTIME ERROR: Missing argument: str
	error.trace_zero_param.jsonnet:19:6-7	object <anonymous>
	During manifestation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.trace_zero_param.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_verbatim_import(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.verbatim_import.jsonnet")
	expected := `RUNTIME ERROR: couldn't open import "C:\\can't possibly exist~": no match locally or in the Jsonnet library paths
	error.verbatim_import.jsonnet:22:1-36	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.verbatim_import.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_error_wrong_type(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/error.wrong_type.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type number, expected string
	error.wrong_type.jsonnet:1:1-18	$
	During evaluation`
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:     "resources/jsonnet-cpp/test_suite",
		File:        "error.wrong_type.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_jsonnet_fmt_idempotence_issue(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/fmt_idempotence_issue.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/fmt_idempotence_issue.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "fmt_idempotence_issue.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_fmt_no_trailing_newline(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/fmt_no_trailing_newline.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/fmt_no_trailing_newline.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "fmt_no_trailing_newline.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_fmt_trailing_c_style(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/fmt_trailing_c_style.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/fmt_trailing_c_style.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "fmt_trailing_c_style.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_fmt_trailing_multiple_comments(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/fmt_trailing_multiple_comments.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/fmt_trailing_multiple_comments.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "fmt_trailing_multiple_comments.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_fmt_trailing_newlines(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/fmt_trailing_newlines.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/fmt_trailing_newlines.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "fmt_trailing_newlines.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_fmt_trailing_newlines2(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/fmt_trailing_newlines2.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/fmt_trailing_newlines2.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "fmt_trailing_newlines2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_fmt_trailing_same_line_comment(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/fmt_trailing_same_line_comment.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/fmt_trailing_same_line_comment.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "fmt_trailing_same_line_comment.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_format(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/format.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "format.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_formatter(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/formatter.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/formatter.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "formatter.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_formatting_braces(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/formatting_braces.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/formatting_braces.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "formatting_braces.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_formatting_braces2(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/formatting_braces2.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/formatting_braces2.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "formatting_braces2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_formatting_braces3(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/formatting_braces3.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "formatting_braces3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_functions(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/functions.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "functions.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_import(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/import.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "import.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_import_sorting(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/import_sorting.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "import_sorting.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_import_sorting_by_filename(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/import_sorting_by_filename.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "import_sorting_by_filename.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_import_sorting_crazy(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/import_sorting_crazy.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "import_sorting_crazy.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_import_sorting_function_sugar(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/import_sorting_function_sugar.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "import_sorting_function_sugar.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_import_sorting_group_ends(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/import_sorting_group_ends.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "import_sorting_group_ends.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_import_sorting_groups(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/import_sorting_groups.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "import_sorting_groups.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_import_sorting_multiple_binds_and_comments(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/import_sorting_multiple_binds_and_comments.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "import_sorting_multiple_binds_and_comments.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_import_sorting_multiple_in_local(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/import_sorting_multiple_in_local.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "import_sorting_multiple_in_local.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_import_sorting_unicode(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/import_sorting_unicode.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "import_sorting_unicode.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_import_sorting_with_license(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/import_sorting_with_license.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "import_sorting_with_license.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_invariant(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/invariant.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "invariant.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_invariant_manifest(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/invariant_manifest.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/invariant_manifest.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "invariant_manifest.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_local(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/local.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "local.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_merge(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/merge.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "merge.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

// TODO: impl std.native
// func Test_jsonnet_native_not_found(t *testing.T) {
// 	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/native_not_found.jsonnet")
// 	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/native_not_found.jsonnet.golden")
// 	extVars := map[string]string{"var1": "test"}
// 	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
// 	testPositive(t, "native_not_found.jsonnet", snippet, expected, extVars, extCodes)
// }

func Test_jsonnet_null(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/null.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "null.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/object.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "object.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_oop(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/oop.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "oop.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_oop_extra(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/oop_extra.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "oop_extra.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_parseJson_long_array_gc_test(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/parseJson_long_array_gc_test.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/parseJson_long_array_gc_test.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "parseJson_long_array_gc_test.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_parsing_edge_cases(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/parsing_edge_cases.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "parsing_edge_cases.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_parsing_error(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/parsing_error.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "parsing_error.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_precedence(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/precedence.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "precedence.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_recursive_function(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/recursive_function.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "recursive_function.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_recursive_import_ok(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/recursive_import_ok.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "recursive_import_ok.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_recursive_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/recursive_object.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "recursive_object.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_safe_integer_conversion(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/safe_integer_conversion.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/safe_integer_conversion.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "safe_integer_conversion.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_sanity(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/sanity.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/sanity.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "sanity.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_sanity2(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/sanity2.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/sanity2.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "sanity2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_shebang(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/shebang.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "shebang.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_slice_sugar(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/slice.sugar.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "slice.sugar.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_std_all_hidden(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/std_all_hidden.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "std_all_hidden.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_stdlib(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/stdlib.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/stdlib.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "stdlib.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_text_block(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/text_block.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "text_block.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_tla_simple(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/tla.simple.jsonnet")
	expected := "true\n"
	tlaVars := map[string]string{"var1": "test"}
	tlaCodes := map[string]string{"var2": "{x:1,y:2}"}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "tla.simple.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		TLAVars:  tlaVars,
		TLACodes: tlaCodes,
	})
}

func Test_jsonnet_trace(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/trace.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/trace.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "trace.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_unicode(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/unicode.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "unicode.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_unicode_bmp(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/unicode_bmp.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/unicode_bmp.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "unicode_bmp.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_unix_line_endings(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/unix_line_endings.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/unix_line_endings.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "unix_line_endings.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_unparse(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/unparse.jsonnet")
	expected := mustReadFile(t, "resources/jsonnet-cpp/test_suite/unparse.jsonnet.golden")
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "unparse.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_jsonnet_verbatim_strings(t *testing.T) {
	snippet := mustReadFile(t, "resources/jsonnet-cpp/test_suite/verbatim_strings.jsonnet")
	expected := "true\n"
	extVars := map[string]string{"var1": "test"}
	extCodes := map[string]string{"var2": `{"x": 1, "y": 2}`}
	runTest(t, TestConfig{
		WorkDir:  "resources/jsonnet-cpp/test_suite",
		File:     "verbatim_strings.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}
