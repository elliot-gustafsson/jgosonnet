package tests

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_go_jsonnet_argcapture_builtin_call(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/argcapture_builtin_call.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/argcapture_builtin_call.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/argcapture_builtin_call.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_array(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/array.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/array.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/array.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_array_comp_try_iterate_over_empty_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/array_comp_try_iterate_over_empty_string.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/array_comp_try_iterate_over_empty_string.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected array
	testdata/array_comp_try_iterate_over_empty_string.jsonnet:1:1-16
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/array_comp_try_iterate_over_empty_string.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_array_comp_try_iterate_over_obj(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/array_comp_try_iterate_over_obj.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/array_comp_try_iterate_over_obj.golden
	expected := `RUNTIME ERROR: Unexpected type object, expected array
	testdata/array_comp_try_iterate_over_obj.jsonnet:1:1-16
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/array_comp_try_iterate_over_obj.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_array_comp_try_iterate_over_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/array_comp_try_iterate_over_string.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/array_comp_try_iterate_over_string.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected array
	testdata/array_comp_try_iterate_over_string.jsonnet:1:1-17
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/array_comp_try_iterate_over_string.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_array_index1(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/array_index1.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/array_index1.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/array_index1.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_array_index2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/array_index2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/array_index2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/array_index2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_array_index3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/array_index3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/array_index3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/array_index3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_array_index4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/array_index4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/array_index4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/array_index4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_array_out_of_bounds(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/array_out_of_bounds.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/array_out_of_bounds.golden
	expected := `RUNTIME ERROR: Index 0 out of bounds, not within [0, 0)
	testdata/array_out_of_bounds.jsonnet:1:1-6	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/array_out_of_bounds.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_array_out_of_bounds2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/array_out_of_bounds2.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/array_out_of_bounds2.golden
	expected := `RUNTIME ERROR: Index 3 out of bounds, not within [0, 3)
	testdata/array_out_of_bounds2.jsonnet:1:1-11	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/array_out_of_bounds2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_array_out_of_bounds3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/array_out_of_bounds3.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/array_out_of_bounds3.golden
	expected := `RUNTIME ERROR: Index -1 out of bounds, not within [0, 0)
	testdata/array_out_of_bounds3.jsonnet:1:1-7	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/array_out_of_bounds3.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_array_out_of_bounds4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/array_out_of_bounds4.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/array_out_of_bounds4.golden
	expected := `RUNTIME ERROR: Index 42 out of bounds, not within [0, 3)
	testdata/array_out_of_bounds4.jsonnet:1:1-12	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/array_out_of_bounds4.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_array_plus_bad(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/array_plus_bad.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/array_plus_bad.golden
	expected := `RUNTIME ERROR: Unexpected type number, expected array
	testdata/array_plus_bad.jsonnet:1:1-8	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/array_plus_bad.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_arrcomp(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/arrcomp.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_arrcomp2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/arrcomp2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_arrcomp3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/arrcomp3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_arrcomp4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/arrcomp4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_arrcomp5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp5.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/arrcomp5.golden
	expected := `RUNTIME ERROR: testdata/arrcomp5.jsonnet:1:14-15 Unknown variable: x
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/arrcomp5.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_arrcomp6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/arrcomp6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_arrcomp7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp7.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp7.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/arrcomp7.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_arrcomp_if(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp_if.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp_if.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/arrcomp_if.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_arrcomp_if2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp_if2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp_if2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/arrcomp_if2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_arrcomp_if3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp_if3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp_if3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/arrcomp_if3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_arrcomp_if4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp_if4.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/arrcomp_if4.golden
	expected := `RUNTIME ERROR: testdata/arrcomp_if4.jsonnet:1:33-34 Unknown variable: y
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/arrcomp_if4.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_arrcomp_if5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp_if5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp_if5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/arrcomp_if5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_arrcomp_if6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp_if6.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/arrcomp_if6.golden
	expected := `RUNTIME ERROR: x
	testdata/arrcomp_if6.jsonnet:1:20-29	$
	testdata/arrcomp_if6.jsonnet:1:1-30
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/arrcomp_if6.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_arrcomp_if7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/arrcomp_if7.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/arrcomp_if7.golden
	expected := `RUNTIME ERROR: Unexpected type number, expected boolean
	testdata/arrcomp_if7.jsonnet:1:1-29
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/arrcomp_if7.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_assert(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/assert.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/assert.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/assert.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_assert2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/assert2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/assert2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/assert2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_assert3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/assert3.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/assert3.golden
	expected := `RUNTIME ERROR: Assertion failed
	testdata/assert3.jsonnet:1:1-20
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/assert3.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_assert_equal(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/assert_equal.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/assert_equal.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/assert_equal.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_assert_equal2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/assert_equal2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/assert_equal2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/assert_equal2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_assert_equal3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/assert_equal3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/assert_equal3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/assert_equal3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_assert_equal4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/assert_equal4.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/assert_equal4.golden
	expected := `RUNTIME ERROR: Assertion failed. {"x": 1} != {"x": 2}
	testdata/assert_equal4.jsonnet:1:1-32	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/assert_equal4.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_assert_equal5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/assert_equal5.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/assert_equal5.golden
	expected := `RUNTIME ERROR: Assertion failed. "\n " != "\n"
	testdata/assert_equal5.jsonnet:1:1-29	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/assert_equal5.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_assert_equal6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/assert_equal6.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/assert_equal6.golden
	expected := `RUNTIME ERROR: Assertion failed. "\u001b[31m" != ""
	testdata/assert_equal6.jsonnet:1:1-34	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/assert_equal6.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_assert_failed(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/assert_failed.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/assert_failed.golden
	expected := `RUNTIME ERROR: Assertion failed
	testdata/assert_failed.jsonnet:1:1-19
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/assert_failed.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_assert_failed_custom(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/assert_failed_custom.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/assert_failed_custom.golden
	expected := `RUNTIME ERROR: Custom Message
	testdata/assert_failed_custom.jsonnet:1:1-38
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/assert_failed_custom.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_bad_function_call(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bad_function_call.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/bad_function_call.golden
	expected := `RUNTIME ERROR: Missing argument: x
	testdata/bad_function_call.jsonnet:1:1-18	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/bad_function_call.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_bad_function_call2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bad_function_call2.jsonnet")
	expected := `RUNTIME ERROR: function expected 1 positional argument(s), but got 2
	testdata/bad_function_call2.jsonnet:1:1-22	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/bad_function_call2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_bad_function_call_and_error(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bad_function_call_and_error.jsonnet")
	expected := `RUNTIME ERROR: function expected 1 positional argument(s), but got 2
	testdata/bad_function_call_and_error.jsonnet:1:1-38	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/bad_function_call_and_error.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_bad_index_array(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bad_index_array.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/bad_index_array.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected number
	testdata/bad_index_array.jsonnet:1:1-10	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/bad_index_array.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_bad_index_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bad_index_object.jsonnet")
	expected := `RUNTIME ERROR: unexpected index type for indexing object, expected string, got number
	testdata/bad_index_object.jsonnet:1:1-7	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/bad_index_object.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_bad_index_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bad_index_string.jsonnet")
	expected := `RUNTIME ERROR: unexpected index type for indexing string, expected number, got string
	testdata/bad_index_string.jsonnet:1:1-13	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/bad_index_string.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_binaryNot(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/binaryNot.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/binaryNot.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/binaryNot.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_binaryNot2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/binaryNot2.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type string, expected number
	testdata/binaryNot2.jsonnet:1:1-7	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/binaryNot2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_bitwise_and(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_and.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_and.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_and.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_and2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_and2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_and2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_and2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_and3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_and3.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/bitwise_and3.golden
	expected := `RUNTIME ERROR: Bitwise operator argument 1e+30 outside of range [-9007199254740991, 9007199254740991]
	testdata/bitwise_and3.jsonnet:1:1-10	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/bitwise_and3.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_bitwise_and4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_and4.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/bitwise_and4.golden
	expected := `RUNTIME ERROR: x
	testdata/bitwise_and4.jsonnet:1:5-14	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/bitwise_and4.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_bitwise_and5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_and5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_and5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_and5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_and6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_and6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_and6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_and6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_and7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_and7.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/bitwise_and7.golden
	expected := `RUNTIME ERROR: Bitwise operator argument -1e+20 outside of range [-9007199254740991, 9007199254740991]
	testdata/bitwise_and7.jsonnet:1:1-11	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/bitwise_and7.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_bitwise_or(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_or.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_or10(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or10.jsonnet")
	expected := `unhandled binary operation string | number
	testdata/bitwise_or10.jsonnet:1:1-11	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/bitwise_or10.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_bitwise_or2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_or2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_or3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_or3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_or4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_or4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_or5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_or5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_or6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_or6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_or7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or7.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or7.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_or7.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_or8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or8.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or8.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_or8.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_or9(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_or9.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/bitwise_or9.golden
	expected := `RUNTIME ERROR: Bitwise operator argument 4.611686018427388e+18 outside of range [-9007199254740991, 9007199254740991]
	testdata/bitwise_or9.jsonnet:1:1-14	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/bitwise_or9.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_bitwise_shift(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_shift.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_shift.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_shift.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_shift2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_shift2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_shift2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_shift2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_shift3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_shift3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_shift3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_shift3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_shift4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_shift4.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/bitwise_shift4.golden
	expected := `RUNTIME ERROR: Shift by negative exponent.
	testdata/bitwise_shift4.jsonnet:1:1-15	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/bitwise_shift4.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_bitwise_shift5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_shift5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_shift5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_shift5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_shift6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_shift6.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/bitwise_shift6.golden
	expected := `RUNTIME ERROR: Shift by negative exponent.
	testdata/bitwise_shift6.jsonnet:1:1-13	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/bitwise_shift6.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_bitwise_xor(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_xor.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_xor2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_xor2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_xor3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_xor3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_xor4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_xor4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_xor5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_xor5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_xor6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_xor6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_xor7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor7.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/bitwise_xor7.golden
	expected := `RUNTIME ERROR: x
	testdata/bitwise_xor7.jsonnet:1:5-14	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/bitwise_xor7.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_bitwise_xor8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor8.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor8.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_xor8.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_bitwise_xor9(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor9.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/bitwise_xor9.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/bitwise_xor9.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_block_escaping(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/block_escaping.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/block_escaping.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/block_escaping.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_block_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/block_string.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/block_string.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/block_string.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_block_string_chomped(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/block_string_chomped.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/block_string_chomped.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/block_string_chomped.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_block_string_chomped_concatted(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/block_string_chomped_concatted.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/block_string_chomped_concatted.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/block_string_chomped_concatted.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_boolean_literal(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/boolean_literal.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/boolean_literal.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/boolean_literal.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinAvg(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinAvg.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinAvg.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinAvg.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinBase64(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinBase64.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinBase64Decode(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64Decode.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64Decode.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinBase64Decode.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinBase64DecodeBytes(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64DecodeBytes.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64DecodeBytes.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinBase64DecodeBytes.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinBase64DecodeBytes_high_codepoint(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64DecodeBytes_high_codepoint.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinBase64DecodeBytes_high_codepoint.golden
	expected := `illegal base64 data at input byte 0
	testdata/builtinBase64DecodeBytes_high_codepoint.jsonnet:1:1-30	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinBase64DecodeBytes_high_codepoint.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinBase64DecodeBytes_invalid_base64_data(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64DecodeBytes_invalid_base64_data.jsonnet")
	expected := `illegal base64 data at input byte 4
	testdata/builtinBase64DecodeBytes_invalid_base64_data.jsonnet:1:1-31	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinBase64DecodeBytes_invalid_base64_data.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinBase64DecodeBytes_wrong_type(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64DecodeBytes_wrong_type.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type number, expected string
	testdata/builtinBase64DecodeBytes_wrong_type.jsonnet:1:1-25	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinBase64DecodeBytes_wrong_type.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinBase64Decode_high_codepoint(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64Decode_high_codepoint.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinBase64Decode_high_codepoint.golden
	expected := `illegal base64 data at input byte 0
	testdata/builtinBase64Decode_high_codepoint.jsonnet:1:1-25	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinBase64Decode_high_codepoint.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinBase64Decode_invalid_base64_data(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64Decode_invalid_base64_data.jsonnet")
	expected := `illegal base64 data at input byte 4
	testdata/builtinBase64Decode_invalid_base64_data.jsonnet:1:1-26	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinBase64Decode_invalid_base64_data.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinBase64Decode_wrong_type(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64Decode_wrong_type.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type number, expected string
	testdata/builtinBase64Decode_wrong_type.jsonnet:1:1-20	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinBase64Decode_wrong_type.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinBase64_byte_array(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64_byte_array.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64_byte_array.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinBase64_byte_array.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinBase64_invalid_byte_array(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64_invalid_byte_array.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinBase64_invalid_byte_array.golden
	expected := `RUNTIME ERROR: base64 encountered a non-integer value in the array, got string
	testdata/builtinBase64_invalid_byte_array.jsonnet:1:1-23	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinBase64_invalid_byte_array.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinBase64_invalid_byte_array1(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64_invalid_byte_array1.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinBase64_invalid_byte_array1.golden
	expected := `RUNTIME ERROR: base64 encountered invalid codepoint value in the array (must be 0 <= X <= 255), got -1
	testdata/builtinBase64_invalid_byte_array1.jsonnet:1:1-20	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinBase64_invalid_byte_array1.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinBase64_invalid_byte_array2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64_invalid_byte_array2.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinBase64_invalid_byte_array2.golden
	expected := `RUNTIME ERROR: base64 encountered invalid codepoint value in the array (must be 0 <= X <= 255), got 256
	testdata/builtinBase64_invalid_byte_array2.jsonnet:1:1-21	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinBase64_invalid_byte_array2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinBase64_non_string_non_array(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64_non_string_non_array.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinBase64_non_string_non_array.golden
	expected := `RUNTIME ERROR: base64 can only base64 encode strings / arrays of single bytes, got number
	testdata/builtinBase64_non_string_non_array.jsonnet:1:1-14	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinBase64_non_string_non_array.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinBase64_string_high_codepoint(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinBase64_string_high_codepoint.jsonnet")
	expected := `RUNTIME ERROR: base64 encountered invalid codepoint value in the array (must be 0 <= X <= 255), got 256
	testdata/builtinBase64_string_high_codepoint.jsonnet:1:1-17	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinBase64_string_high_codepoint.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinChar(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinChar.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinChar.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinChar.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinChar2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinChar2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinChar2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinChar2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinChar3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinChar3.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinChar3.golden
	expected := `RUNTIME ERROR: codepoints must be >= 0, got -1
	testdata/builtinChar3.jsonnet:1:1-13	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinChar3.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinChar4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinChar4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinChar4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinChar4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinChar5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinChar5.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinChar5.golden
	expected := `RUNTIME ERROR: invalid unicode codepoint, got 1.114112e+06
	testdata/builtinChar5.jsonnet:2:1-18	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinChar5.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinChar6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinChar6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinChar6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinChar6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinChar7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinChar7.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinChar7.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected number
	testdata/builtinChar7.jsonnet:1:1-16	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinChar7.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinContains(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinContains.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinContains.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinContains.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinContains2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinContains2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinContains2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinContains2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinEqualsIgnoreCase(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinEqualsIgnoreCase.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinEqualsIgnoreCase.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinEqualsIgnoreCase.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinEqualsIgnoreCase2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinEqualsIgnoreCase2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinEqualsIgnoreCase2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinEqualsIgnoreCase2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinIsDecimal(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsDecimal.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsDecimal.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinIsDecimal.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinIsDecimal2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsDecimal2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsDecimal2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinIsDecimal2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinIsEmpty(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsEmpty.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsEmpty.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinIsEmpty.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinIsEmpty1(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsEmpty1.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsEmpty1.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinIsEmpty1.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinIsEmpty2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsEmpty2.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinIsEmpty2.golden
	expected := `RUNTIME ERROR: Unexpected type number, expected string
	testdata/builtinIsEmpty2.jsonnet:1:1-16	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinIsEmpty2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinIsEven(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsEven.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsEven.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinIsEven.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinIsEven2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsEven2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsEven2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinIsEven2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinIsInteger(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsInteger.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsInteger.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinIsInteger.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinIsInteger2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsInteger2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsInteger2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinIsInteger2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinIsNull(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsNull.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsNull.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinIsNull.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinIsNull2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsNull2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsNull2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinIsNull2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinIsOdd(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsOdd.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsOdd.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinIsOdd.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinIsOdd2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsOdd2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinIsOdd2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinIsOdd2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinManifestJsonEx(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinManifestJsonEx.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinManifestJsonEx.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinManifestJsonEx.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinManifestJsonEx_cyclic(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinManifestJsonEx_cyclic.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinManifestJsonEx_cyclic.golden
	expected := `RUNTIME ERROR: max manifest depth exceeded, possible infinite recursion
	testdata/builtinManifestJsonEx_cyclic.jsonnet:1:1-32	$
	During evaluation`
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinManifestJsonEx_cyclic.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		MaxStack:    500,
	})
}

func Test_go_jsonnet_builtinMaxArray(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinMaxArray.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinMaxArray.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinMaxArray.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinMinArray(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinMinArray.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinMinArray.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinMinArray.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinObjectFieldsEx(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectFieldsEx.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectFieldsEx.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinObjectFieldsEx.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinObjectFieldsExWithHidden(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectFieldsExWithHidden.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectFieldsExWithHidden.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinObjectFieldsExWithHidden.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinObjectFieldsEx_bad(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectFieldsEx_bad.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinObjectFieldsEx_bad.golden
	expected := `RUNTIME ERROR: Unexpected type number, expected object
	testdata/builtinObjectFieldsEx_bad.jsonnet:1:1-29	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinObjectFieldsEx_bad.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinObjectFieldsEx_bad2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectFieldsEx_bad2.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinObjectFieldsEx_bad2.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected boolean
	testdata/builtinObjectFieldsEx_bad2.jsonnet:1:1-30	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinObjectFieldsEx_bad2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinObjectHasEx(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectHasEx.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectHasEx.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinObjectHasEx.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinObjectHasExBadBoolean(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectHasExBadBoolean.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinObjectHasExBadBoolean.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected boolean
	testdata/builtinObjectHasExBadBoolean.jsonnet:1:1-34	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinObjectHasExBadBoolean.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinObjectHasExBadField(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectHasExBadField.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinObjectHasExBadField.golden
	expected := `RUNTIME ERROR: Unexpected type number, expected string
	testdata/builtinObjectHasExBadField.jsonnet:1:1-31	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinObjectHasExBadField.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinObjectHasExBadObject(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectHasExBadObject.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinObjectHasExBadObject.golden
	expected := `RUNTIME ERROR: Unexpected type number, expected object
	testdata/builtinObjectHasExBadObject.jsonnet:1:1-32	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinObjectHasExBadObject.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinObjectRemoveKey(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectRemoveKey.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectRemoveKey.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinObjectRemoveKey.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinObjectRemoveKey_hidden(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectRemoveKey_hidden.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectRemoveKey_hidden.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinObjectRemoveKey_hidden.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinObjectRemoveKey_super(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectRemoveKey_super.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectRemoveKey_super.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinObjectRemoveKey_super.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinObjectRemoveKey_super_assert(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinObjectRemoveKey_super_assert.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinObjectRemoveKey_super_assert.golden
	expected := `RUNTIME ERROR: Field does not exist: x
	testdata/builtinObjectRemoveKey_super_assert.jsonnet:2:10-16	object <o1>
	testdata/builtinObjectRemoveKey_super_assert.jsonnet:2:3-22
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinObjectRemoveKey_super_assert.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinRemove(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinRemove.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinRemove.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinRemove.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinRemoveAt(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinRemoveAt.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinRemoveAt.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinRemoveAt.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinRemoveAt2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinRemoveAt2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinRemoveAt2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinRemoveAt2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinReverse(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinReverse.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinReverse.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinReverse.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinReverse_empty(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinReverse_empty.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinReverse_empty.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinReverse_empty.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinReverse_many(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinReverse_many.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinReverse_many.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinReverse_many.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinReverse_not_array(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinReverse_not_array.jsonnet")
	expected := `unexpected type passed to std.reverse (arg 0): boolean, expected array
	testdata/builtinReverse_not_array.jsonnet:1:1-19	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinReverse_not_array.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinReverse_single(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinReverse_single.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinReverse_single.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinReverse_single.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinRound(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinRound.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinRound.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinRound.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinSha1(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSha1.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSha1.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinSha1.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinSha256(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSha256.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSha256.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinSha256.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinSha3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSha3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSha3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinSha3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinSha512(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSha512.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSha512.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinSha512.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinSplitLimitR(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSplitLimitR.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSplitLimitR.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinSplitLimitR.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinSplitLimitR2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSplitLimitR2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSplitLimitR2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinSplitLimitR2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinSplitLimitR3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSplitLimitR3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSplitLimitR3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinSplitLimitR3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinSplitLimitR4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSplitLimitR4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSplitLimitR4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinSplitLimitR4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinSplitLimitR5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSplitLimitR5.jsonnet")
	expected := `RUNTIME ERROR: std.splitLimitR third parameter should be -1 or non-negative, got -2
	testdata/builtinSplitLimitR5.jsonnet:1:1-44	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinSplitLimitR5.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinSplitLimitR6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSplitLimitR6.jsonnet")
	expected := `RUNTIME ERROR: std.splitLimitR second parameter should have length 1 or greater, got 0
	testdata/builtinSplitLimitR6.jsonnet:1:1-43	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinSplitLimitR6.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinSubStr_first_param_not_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSubStr_first_param_not_string.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type number, expected string
	testdata/builtinSubStr_first_param_not_string.jsonnet:1:1-20	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinSubStr_first_param_not_string.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinSubStr_length_larger(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSubStr_length_larger.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSubStr_length_larger.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinSubStr_length_larger.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinSubStr_second_parameter_not_integer(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSubStr_second_parameter_not_integer.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinSubStr_second_parameter_not_integer.golden
	expected := `RUNTIME ERROR: substr second parameter should be an integer, got 1.200000
	testdata/builtinSubStr_second_parameter_not_integer.jsonnet:1:1-28	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinSubStr_second_parameter_not_integer.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinSubStr_second_parameter_not_number(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSubStr_second_parameter_not_number.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type string, expected number
	testdata/builtinSubStr_second_parameter_not_number.jsonnet:1:1-30	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinSubStr_second_parameter_not_number.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinSubStr_start_larger_then_size(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSubStr_start_larger_then_size.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSubStr_start_larger_then_size.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinSubStr_start_larger_then_size.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinSubStr_third_parameter_less_then_zero(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSubStr_third_parameter_less_then_zero.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinSubStr_third_parameter_less_then_zero.golden
	expected := `RUNTIME ERROR: substr third parameter should be greater than zero, got -1
	testdata/builtinSubStr_third_parameter_less_then_zero.jsonnet:1:1-27	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinSubStr_third_parameter_less_then_zero.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinSubStr_third_parameter_not_integer(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSubStr_third_parameter_not_integer.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinSubStr_third_parameter_not_integer.golden
	expected := `RUNTIME ERROR: substr third parameter should be an integer, got 1.200000
	testdata/builtinSubStr_third_parameter_not_integer.jsonnet:1:1-28	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinSubStr_third_parameter_not_integer.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinSubStr_third_parameter_not_number(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSubStr_third_parameter_not_number.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type string, expected number
	testdata/builtinSubStr_third_parameter_not_number.jsonnet:1:1-30	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinSubStr_third_parameter_not_number.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinSubstr(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSubstr.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSubstr.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinSubstr.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinSum(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSum.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinSum.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinSum.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinTrim(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinTrim.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinTrim.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinTrim.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinTrim1(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinTrim1.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinTrim1.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinTrim1.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinTrim2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinTrim2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinTrim2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinTrim2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinTrim3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinTrim3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinTrim3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinTrim3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinTrim4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinTrim4.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinTrim4.golden
	expected := `RUNTIME ERROR: Unexpected type number, expected string
	testdata/builtinTrim4.jsonnet:1:1-13	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinTrim4.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinXnor(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinXnor.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinXnor.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinXnor.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinXnor1(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinXnor1.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinXnor1.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinXnor1.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinXnor2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinXnor2.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinXnor2.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected boolean
	testdata/builtinXnor2.jsonnet:1:1-24	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinXnor2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtinXor(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinXor.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinXor.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinXor.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinXor1(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinXor1.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtinXor1.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtinXor1.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtinXor2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtinXor2.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtinXor2.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected boolean
	testdata/builtinXor2.jsonnet:1:1-23	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtinXor2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtin_acos(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_acos.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_acos.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_acos.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_asin(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_asin.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_asin.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_asin.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_atan(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_atan.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_atan.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_atan.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_ceil(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_ceil.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_ceil.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_ceil.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_cos(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_cos.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_cos.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_cos.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_escapeStringJson(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_escapeStringJson.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_escapeStringJson.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_escapeStringJson.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_exp(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_exp.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_exp.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_exp.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_exp2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_exp2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_exp2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_exp2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_exp3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_exp3.jsonnet")
	expected := `RUNTIME ERROR: Overflow
	testdata/builtin_exp3.jsonnet:1:1-14	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtin_exp3.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtin_exp4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_exp4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_exp4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_exp4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_exp5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_exp5.jsonnet")
	expected := `RUNTIME ERROR: Overflow
	testdata/builtin_exp5.jsonnet:1:1-31	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtin_exp5.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtin_exp6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_exp6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_exp6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_exp6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_exp7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_exp7.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_exp7.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_exp7.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_exp8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_exp8.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_exp8.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_exp8.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_floor(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_floor.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_floor.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_floor.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_log(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_log.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_log.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_log.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_log2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_log2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_log2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_log2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_log3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_log3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_log3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_log3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_log4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_log4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_log4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_log4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_log5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_log5.jsonnet")
	expected := `RUNTIME ERROR: Overflow
	testdata/builtin_log5.jsonnet:1:1-11	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtin_log5.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtin_log6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_log6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_log6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_log6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_log7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_log7.jsonnet")
	expected := `RUNTIME ERROR: Not a number
	testdata/builtin_log7.jsonnet:1:1-12	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtin_log7.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtin_log8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_log8.jsonnet")
	expected := `RUNTIME ERROR: Not a number
	testdata/builtin_log8.jsonnet:1:1-24	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtin_log8.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtin_lstripChars(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_lstripChars.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_lstripChars.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_lstripChars.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_manifestTomlEx(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_manifestTomlEx.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_manifestTomlEx.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_manifestTomlEx.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_manifestTomlEx_array(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_manifestTomlEx_array.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type array, expected object
	testdata/builtin_manifestTomlEx_array.jsonnet:11:10-41	object <anonymous>
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtin_manifestTomlEx_array.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtin_manifestTomlEx_cyclic(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_manifestTomlEx_cyclic.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtin_manifestTomlEx_cyclic.golden
	expected := `RUNTIME ERROR: max manifest depth exceeded, possible infinite recursion
	testdata/builtin_manifestTomlEx_cyclic.jsonnet:1:1-37	$
	During evaluation`
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtin_manifestTomlEx_cyclic.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		MaxStack:    500,
	})
}

func Test_go_jsonnet_builtin_manifestTomlEx_null(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_manifestTomlEx_null.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type null, expected object
	testdata/builtin_manifestTomlEx_null.jsonnet:2:11-42	object <anonymous>
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtin_manifestTomlEx_null.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtin_manifestYamlDoc(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_manifestYamlDoc.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_manifestYamlDoc.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_manifestYamlDoc.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_manifestYamlDoc_cyclic(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_manifestYamlDoc_cyclic.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtin_manifestYamlDoc_cyclic.golden
	expected := `RUNTIME ERROR: max manifest depth exceeded, possible infinite recursion
	testdata/builtin_manifestYamlDoc_cyclic.jsonnet:1:1-28	$
	During evaluation`
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtin_manifestYamlDoc_cyclic.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		MaxStack:    500,
	})
}

func Test_go_jsonnet_builtin_member_array(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_member_array.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_member_array.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_member_array.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_member_object_invalid(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_member_object_invalid.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type string, expected array
	testdata/builtin_member_object_invalid.jsonnet:1:1-31	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtin_member_object_invalid.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtin_member_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_member_string.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_member_string.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_member_string.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_parseInt(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_parseInt.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_parseInt.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_parseInt.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_parseInt2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_parseInt2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_parseInt2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_parseInt2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_parseInt_invalid(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_parseInt_invalid.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtin_parseInt_invalid.golden
	expected := `error: hello is not a base 10 integer
	testdata/builtin_parseInt_invalid.jsonnet:1:1-22	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtin_parseInt_invalid.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtin_parseInt_invalid_decimal(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_parseInt_invalid_decimal.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtin_parseInt_invalid_decimal.golden
	expected := `error: 123.12 is not a base 10 integer
	testdata/builtin_parseInt_invalid_decimal.jsonnet:1:1-23	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtin_parseInt_invalid_decimal.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtin_parseInt_invalid_hexadecimal(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_parseInt_invalid_hexadecimal.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtin_parseInt_invalid_hexadecimal.golden
	expected := `error: 7B316 is not a base 10 integer
	testdata/builtin_parseInt_invalid_hexadecimal.jsonnet:1:1-22	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtin_parseInt_invalid_hexadecimal.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtin_parseYaml_empty(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_parseYaml_empty.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_parseYaml_empty.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_parseYaml_empty.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_range_negative(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_range_negative.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_range_negative.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_range_negative.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_rstripChars(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_rstripChars.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_rstripChars.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_rstripChars.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_sin(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_sin.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_sin.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_sin.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_sqrt(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_sqrt.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_sqrt.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_sqrt.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_sqrt2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_sqrt2.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtin_sqrt2.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected number
	testdata/builtin_sqrt2.jsonnet:1:1-19	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtin_sqrt2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtin_stripChars(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_stripChars.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_stripChars.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_stripChars.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_stripChars_invalid(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_stripChars_invalid.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/builtin_stripChars_invalid.golden
	expected := `RUNTIME ERROR: Unexpected type object, expected string
	testdata/builtin_stripChars_invalid.jsonnet:1:1-4132	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/builtin_stripChars_invalid.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_builtin_substr_multibyte(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_substr_multibyte.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_substr_multibyte.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_substr_multibyte.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_builtin_tan(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_tan.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/builtin_tan.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/builtin_tan.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_call_number(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/call_number.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/call_number.golden
	expected := `RUNTIME ERROR: Unexpected type number, expected function
	testdata/call_number.jsonnet:1:1-5	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/call_number.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_comparisons(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/comparisons.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/comparisons.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/comparisons.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_decodeUTF8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/decodeUTF8.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/decodeUTF8.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/decodeUTF8.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_digitsep(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/digitsep.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/digitsep.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/digitsep.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_div1(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/div1.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/div1.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/div1.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_div2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/div2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/div2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/div2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_div3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/div3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/div3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/div3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_div4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/div4.jsonnet")
	expected := `RUNTIME ERROR: Overflow
	testdata/div4.jsonnet:1:1-20	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/div4.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_div_by_zero(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/div_by_zero.jsonnet")
	expected := `RUNTIME ERROR: Division by zero.
	testdata/div_by_zero.jsonnet:1:1-4	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/div_by_zero.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_dollar_bad(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/dollar_bad.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/dollar_bad.golden
	expected := `RUNTIME ERROR: testdata/dollar_bad.jsonnet:1:1-2 No top-level object found.
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/dollar_bad.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_dollar_end(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/dollar_end.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/dollar_end.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/dollar_end.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_dollar_end2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/dollar_end2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/dollar_end2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/dollar_end2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_double_thunk(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/double_thunk.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/double_thunk.golden
	expected := `RUNTIME ERROR: xxx
	testdata/double_thunk.jsonnet:1:34-35	thunk <x> from <$>
	testdata/double_thunk.jsonnet:1:37-38	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/double_thunk.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_empty_array(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/empty_array.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/empty_array.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/empty_array.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_empty_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/empty_object.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/empty_object.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/empty_object.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_empty_object_comp(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/empty_object_comp.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/empty_object_comp.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/empty_object_comp.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_encodeUTF8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/encodeUTF8.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/encodeUTF8.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/encodeUTF8.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_equals(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/equals.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/equals.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/equals.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_equals2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/equals2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/equals2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/equals2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_equals3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/equals3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/equals3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/equals3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_equals4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/equals4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/equals4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/equals4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_equals5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/equals5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/equals5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/equals5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_equals6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/equals6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/equals6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/equals6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_error(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/error.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/error.golden
	expected := `RUNTIME ERROR: 42
	testdata/error.jsonnet:1:1-11	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/error.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_error_from_array(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/error_from_array.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/error_from_array.golden
	expected := `RUNTIME ERROR: xxx
	testdata/error_from_array.jsonnet:1:1-17	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/error_from_array.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_error_from_func(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/error_from_func.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/error_from_func.golden
	expected := `RUNTIME ERROR: xxx
	testdata/error_from_func.jsonnet:1:25-32	function <foo>
	testdata/error_from_func.jsonnet:1:34-44	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/error_from_func.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_error_function_fail(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/error_function_fail.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type function, expected string
	testdata/error_function_fail.jsonnet:1:1-23	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/error_function_fail.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_error_hexnumber(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/error_hexnumber.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/error_hexnumber.golden
	expected := `RUNTIME ERROR: testdata/error_hexnumber.jsonnet:1:2-5 Did not expect: (IDENTIFIER, "x42")
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/error_hexnumber.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_error_in_method(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/error_in_method.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/error_in_method.golden
	expected := `RUNTIME ERROR: xxx
	testdata/error_in_method.jsonnet:1:23-30	function <anonymous>
	testdata/error_in_method.jsonnet:1:34-48	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/error_in_method.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_error_in_object_local(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/error_in_object_local.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/error_in_object_local.golden
	expected := `RUNTIME ERROR: xxx
	testdata/error_in_object_local.jsonnet:1:20-29	function <anonymous>
	testdata/error_in_object_local.jsonnet:1:36-46	object <anonymous>
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/error_in_object_local.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_error_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/error_object.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type object, expected string
	testdata/error_object.jsonnet:1:1-21	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/error_object.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_escaped_fields(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/escaped_fields.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/escaped_fields.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/escaped_fields.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_escaped_single_quote(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/escaped_single_quote.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/escaped_single_quote.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/escaped_single_quote.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_extvar_code(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/extvar_code.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/extvar_code.golden")
	extVars := standardExtVars
	extCodes := standardExtCode
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/extvar_code.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_extvar_error(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/extvar_error.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/extvar_error.golden
	expected := `RUNTIME ERROR: xxx
	<extvar:errorVar>:1:1-12	$
	During evaluation`
	extVars := standardExtVars
	extCodes := standardExtCode
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/extvar_error.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_extvar_hermetic(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/extvar_hermetic.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/extvar_hermetic.golden
	expected := `<extvar:UndeclaredX>:1:1-2 Unknown variable: x
	testdata/extvar_hermetic.jsonnet:1:15-40	$
	During evaluation`
	extVars := standardExtVars
	extCodes := standardExtCode
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/extvar_hermetic.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_extvar_mutually_recursive(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/extvar_mutually_recursive.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/extvar_mutually_recursive.golden")
	extVars := standardExtVars
	extCodes := standardExtCode
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/extvar_mutually_recursive.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_extvar_not_a_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/extvar_not_a_string.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/extvar_not_a_string.golden
	expected := `RUNTIME ERROR: Unexpected type number, expected string
	testdata/extvar_not_a_string.jsonnet:1:1-15	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/extvar_not_a_string.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_extvar_self_recursive(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/extvar_self_recursive.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/extvar_self_recursive.golden")
	extVars := standardExtVars
	extCodes := standardExtCode
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/extvar_self_recursive.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_extvar_static_error(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/extvar_static_error.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/extvar_static_error.golden
	expected := `<extvar:staticErrorVar>:1:1-2 Unexpected: ")" while parsing terminal
	testdata/extvar_static_error.jsonnet:1:1-29	$
	During evaluation`
	extVars := standardExtVars
	extCodes := standardExtCode
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/extvar_static_error.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_extvar_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/extvar_string.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/extvar_string.golden")
	extVars := standardExtVars
	extCodes := standardExtCode
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/extvar_string.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_extvar_unknown(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/extvar_unknown.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/extvar_unknown.golden
	expected := `undefined external variable: UNKNOWN
	testdata/extvar_unknown.jsonnet:1:1-22	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/extvar_unknown.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_false(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/false.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/false.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/false.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_fieldname_not_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/fieldname_not_string.jsonnet")
	expected := `unexpected field name type number, expected string
	testdata/fieldname_not_string.jsonnet:1:1-15	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/fieldname_not_string.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_filled_thunk(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/filled_thunk.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/filled_thunk.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/filled_thunk.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_foldl_empty(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/foldl_empty.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/foldl_empty.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/foldl_empty.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_foldl_single_element(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/foldl_single_element.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/foldl_single_element.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/foldl_single_element.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_foldl_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/foldl_string.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/foldl_string.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/foldl_string.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_foldl_various(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/foldl_various.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/foldl_various.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/foldl_various.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_foldr_empty(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/foldr_empty.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/foldr_empty.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/foldr_empty.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_foldr_single_element(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/foldr_single_element.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/foldr_single_element.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/foldr_single_element.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_foldr_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/foldr_string.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/foldr_string.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/foldr_string.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_foldr_various(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/foldr_various.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/foldr_various.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/foldr_various.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_function_call(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/function_call.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/function_call.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/function_call.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_function_capturing(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/function_capturing.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/function_capturing.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/function_capturing.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_function_in_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/function_in_object.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/function_in_object.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/function_in_object.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_function_manifested(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/function_manifested.jsonnet")
	expected := `unhandled value type: function	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/function_manifested.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_function_no_params(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/function_no_params.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/function_no_params.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/function_no_params.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_function_plus_bad(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/function_plus_bad.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type number, expected function
	testdata/function_plus_bad.jsonnet:1:1-22	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/function_plus_bad.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_function_plus_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/function_plus_string.jsonnet")
	expected := `unhandled type function, string conversion not available
	testdata/function_plus_string.jsonnet:1:1-24	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/function_plus_string.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_function_too_many_params(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/function_too_many_params.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/function_too_many_params.golden
	expected := `RUNTIME ERROR: Missing argument: x	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/function_too_many_params.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_function_with_argument(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/function_with_argument.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/function_with_argument.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/function_with_argument.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_greater(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/greater.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/greater.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/greater.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_greaterEq(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/greaterEq.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/greaterEq.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/greaterEq.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_greaterEq2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/greaterEq2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/greaterEq2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/greaterEq2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_ifthen_false(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/ifthen_false.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/ifthen_false.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/ifthen_false.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_ifthenelse_false(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/ifthenelse_false.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/ifthenelse_false.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/ifthenelse_false.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_ifthenelse_true(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/ifthenelse_true.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/ifthenelse_true.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/ifthenelse_true.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_import(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/import.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/import.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/import.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_import2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/import2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/import2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/import2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_import3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/import3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/import3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/import3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_import4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/import4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/import4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/import4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_import_block_literal(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/import_block_literal.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/import_block_literal.golden
	expected := `RUNTIME ERROR: testdata/import_block_literal.jsonnet:(1:8)-(3:4) Block string literals not allowed in imports
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/import_block_literal.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_import_computed(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/import_computed.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/import_computed.golden
	expected := `RUNTIME ERROR: testdata/import_computed.jsonnet:1:8-17 Computed imports are not allowed
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/import_computed.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_import_failure_directory(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/import_failure_directory.jsonnet")
	expected := `RUNTIME ERROR: read testdata: is a directory
	testdata/import_failure_directory.jsonnet:3:1-11	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/import_failure_directory.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_import_syntax_error(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/import_syntax_error.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/import_syntax_error.golden
	expected := `RUNTIME ERROR: testdata/syntax_error.jsonnet:2:1 Unexpected end of file
	testdata/import_syntax_error.jsonnet:1:1-30	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/import_syntax_error.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_import_twice(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/import_twice.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/import_twice.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/import_twice.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_import_various_literals(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/import_various_literals.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/import_various_literals.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/import_various_literals.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_importbin_block_literal(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/importbin_block_literal.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/importbin_block_literal.golden
	expected := `RUNTIME ERROR: testdata/importbin_block_literal.jsonnet:(1:11)-(3:4) Block string literals not allowed in imports
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/importbin_block_literal.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_importbin_computed(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/importbin_computed.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/importbin_computed.golden
	expected := `RUNTIME ERROR: testdata/importbin_computed.jsonnet:1:11-20 Computed imports are not allowed
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/importbin_computed.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_importbin_nonutf8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/importbin_nonutf8.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/importbin_nonutf8.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/importbin_nonutf8.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_importstr_block_literal(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/importstr_block_literal.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/importstr_block_literal.golden
	expected := `RUNTIME ERROR: testdata/importstr_block_literal.jsonnet:(1:11)-(3:4) Block string literals not allowed in imports
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/importstr_block_literal.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_importstr_computed(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/importstr_computed.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/importstr_computed.golden
	expected := `RUNTIME ERROR: testdata/importstr_computed.jsonnet:1:11-20 Computed imports are not allowed
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/importstr_computed.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_in(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/in.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/in.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/in.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_in2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/in2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/in2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/in2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_in3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/in3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/in3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/in3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_in4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/in4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/in4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/in4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_inf_min_number(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/inf_min_number.jsonnet")
	expected := `RUNTIME ERROR: Overflow
	testdata/inf_min_number.jsonnet:1:1-15	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/inf_min_number.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_inf_mul_number(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/inf_mul_number.jsonnet")
	expected := `RUNTIME ERROR: Overflow
	testdata/inf_mul_number.jsonnet:1:1-19	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/inf_mul_number.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_inf_sum_number(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/inf_sum_number.jsonnet")
	expected := `RUNTIME ERROR: Overflow
	testdata/inf_sum_number.jsonnet:1:1-14	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/inf_sum_number.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_insuper(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/insuper.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/insuper.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/insuper.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_insuper2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/insuper2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/insuper2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/insuper2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_insuper3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/insuper3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/insuper3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/insuper3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_insuper4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/insuper4.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/insuper4.golden
	expected := `RUNTIME ERROR: testdata/insuper4.jsonnet:1:1-13 Can't use super outside of an object.
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/insuper4.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_insuper5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/insuper5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/insuper5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/insuper5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_insuper6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/insuper6.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/insuper6.golden
	expected := `RUNTIME ERROR: testdata/insuper6.jsonnet:1:10-20 Unknown variable: undeclared
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/insuper6.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_insuper7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/insuper7.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/insuper7.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/insuper7.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_lazy(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/lazy.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/lazy.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/lazy.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_lazy_operator1(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/lazy_operator1.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/lazy_operator1.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/lazy_operator1.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_lazy_operator2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/lazy_operator2.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/lazy_operator2.golden
	expected := `RUNTIME ERROR: should happen
	testdata/lazy_operator2.jsonnet:1:9-30	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/lazy_operator2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_less(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/less.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/less.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/less.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_lessEq(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/lessEq.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/lessEq.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/lessEq.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_lessEq2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/lessEq2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/lessEq2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/lessEq2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_local_in_object_assertion(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/local_in_object_assertion.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/local_in_object_assertion.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/local_in_object_assertion.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_local_within_nested_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/local_within_nested_object.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/local_within_nested_object.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/local_within_nested_object.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_local_within_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/local_within_object.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/local_within_object.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/local_within_object.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_method_call(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/method_call.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/method_call.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/method_call.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_missing_super(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/missing_super.jsonnet")
	expected := `RUNTIME ERROR: Field does not exist: x
	testdata/missing_super.jsonnet:1:6-11	object <anonymous>
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/missing_super.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_modulo(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/modulo.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/modulo.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/modulo.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_modulo2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/modulo2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/modulo2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/modulo2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_modulo3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/modulo3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/modulo3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/modulo3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_modulo4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/modulo4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/modulo4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/modulo4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_modulo5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/modulo5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/modulo5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/modulo5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_modulo6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/modulo6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/modulo6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/modulo6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_modulo7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/modulo7.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/modulo7.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/modulo7.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_mult(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/mult.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/mult.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/mult.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_mult2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/mult2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/mult2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/mult2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_mult3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/mult3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/mult3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/mult3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_multi(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/multi.jsonnet")
	expected := readGoldenDir(t, "resources/go-jsonnet/testdata/multi.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/multi.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
		IsMulti:  true,
	})
}

func Test_go_jsonnet_multi_no_newline(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/multi_no_newline.jsonnet")
	expected := readGoldenDir(t, "resources/go-jsonnet/testdata/multi_no_newline.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:   "resources/go-jsonnet",
		File:      "testdata/multi_no_newline.jsonnet",
		Snippet:   snippet,
		Expected:  expected,
		ExtVars:   extVars,
		ExtCodes:  extCodes,
		IsMulti:   true,
		NoNewline: true,
	})
}

func Test_go_jsonnet_multi_no_newline_string_output(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/multi_no_newline_string_output.jsonnet")
	expected := readGoldenDir(t, "resources/go-jsonnet/testdata/multi_no_newline_string_output.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:      "resources/go-jsonnet",
		File:         "testdata/multi_no_newline_string_output.jsonnet",
		Snippet:      snippet,
		Expected:     expected,
		ExtVars:      extVars,
		ExtCodes:     extCodes,
		IsMulti:      true,
		NoNewline:    true,
		StringOutput: true,
	})
}

func Test_go_jsonnet_multi_string_output(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/multi_string_output.jsonnet")
	expected := readGoldenDir(t, "resources/go-jsonnet/testdata/multi_string_output.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:      "resources/go-jsonnet",
		File:         "testdata/multi_string_output.jsonnet",
		Snippet:      snippet,
		Expected:     expected,
		ExtVars:      extVars,
		ExtCodes:     extCodes,
		IsMulti:      true,
		StringOutput: true,
	})
}

// TODO: impl std.native
// func Test_go_jsonnet_native1(t *testing.T) {
// 	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/native1.jsonnet")
// 	expected := mustReadFile(t, "resources/go-jsonnet/testdata/native1.golden")
// 	extVars := map[string]string(nil)
// 	extCodes := map[string]string(nil)
// 	testPositive(t, "testdata/native1.jsonnet", snippet, expected, extVars, extCodes)
// }

// TODO: impl std.native
// func Test_go_jsonnet_native2(t *testing.T) {
// 	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/native2.jsonnet")
// 	expected := mustReadFile(t, "resources/go-jsonnet/testdata/native2.golden")
// 	extVars := map[string]string(nil)
// 	extCodes := map[string]string(nil)
// 	testPositive(t, "testdata/native2.jsonnet", snippet, expected, extVars, extCodes)
// }

// TODO: impl std.native
// func Test_go_jsonnet_native3(t *testing.T) {
// 	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/native3.jsonnet")
// 	expected := mustReadFile(t, "resources/go-jsonnet/testdata/native3.golden")
// 	extVars := map[string]string(nil)
// 	extCodes := map[string]string(nil)
// 	testPositive(t, "testdata/native3.jsonnet", snippet, expected, extVars, extCodes)
// }

// TODO: impl std.native
// func Test_go_jsonnet_native4(t *testing.T) {
// 	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/native4.jsonnet")
// 	// Expected file: resources/go-jsonnet/testdata/native4.golden
// 	expected := `RUNTIME ERROR: xxx`
// 	extVars := map[string]string(nil)
// 	extCodes := map[string]string(nil)
// 	testNegative(t, "testdata/native4.jsonnet", snippet, expected, extVars, extCodes)
// }

// TODO: impl std.native
// func Test_go_jsonnet_native5(t *testing.T) {
// 	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/native5.jsonnet")
// 	// Expected file: resources/go-jsonnet/testdata/native5.golden
// 	expected := `RUNTIME ERROR: couldn't manifest function as JSON`
// 	extVars := map[string]string(nil)
// 	extCodes := map[string]string(nil)
// 	testNegative(t, "testdata/native5.jsonnet", snippet, expected, extVars, extCodes)
// }

// TODO: impl std.native
// func Test_go_jsonnet_native6(t *testing.T) {
// 	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/native6.jsonnet")
// 	expected := mustReadFile(t, "resources/go-jsonnet/testdata/native6.golden")
// 	extVars := map[string]string(nil)
// 	extCodes := map[string]string(nil)
// 	testPositive(t, "testdata/native6.jsonnet", snippet, expected, extVars, extCodes)
// }

// TODO: impl std.native
// func Test_go_jsonnet_native7(t *testing.T) {
// 	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/native7.jsonnet")
// 	// Expected file: resources/go-jsonnet/testdata/native7.golden
// 	expected := `RUNTIME ERROR: function has no parameter y`
// 	extVars := map[string]string(nil)
// 	extCodes := map[string]string(nil)
// 	testNegative(t, "testdata/native7.jsonnet", snippet, expected, extVars, extCodes)
// }

// TODO: impl std.native
// func Test_go_jsonnet_native_error(t *testing.T) {
// 	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/native_error.jsonnet")
// 	// Expected file: resources/go-jsonnet/testdata/native_error.golden
// 	expected := `RUNTIME ERROR: native function error`
// 	extVars := map[string]string(nil)
// 	extCodes := map[string]string(nil)
// 	testNegative(t, "testdata/native_error.jsonnet", snippet, expected, extVars, extCodes)
// }

// TODO: impl std.native
// func Test_go_jsonnet_native_nonexistent(t *testing.T) {
// 	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/native_nonexistent.jsonnet")
// 	expected := mustReadFile(t, "resources/go-jsonnet/testdata/native_nonexistent.golden")
// 	extVars := map[string]string(nil)
// 	extCodes := map[string]string(nil)
// 	testPositive(t, "testdata/native_nonexistent.jsonnet", snippet, expected, extVars, extCodes)
// }

// TODO: impl std.native
// func Test_go_jsonnet_native_panic(t *testing.T) {
// 	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/native_panic.jsonnet")
// 	// Expected file: resources/go-jsonnet/testdata/native_panic.golden
// 	expected := `RUNTIME ERROR: native function "nativePanic" panicked: native function panic`
// 	extVars := map[string]string(nil)
// 	extCodes := map[string]string(nil)
// 	testNegative(t, "testdata/native_panic.jsonnet", snippet, expected, extVars, extCodes)
// }

func Test_go_jsonnet_nonexistent_import(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/nonexistent_import.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/nonexistent_import.golden
	expected := `RUNTIME ERROR: couldn't open import "no chance a file with this name exists": no match locally or in the Jsonnet library paths
	testdata/nonexistent_import.jsonnet:1:1-51	$
	During evaluation`
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/nonexistent_import.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
	})
}

func Test_go_jsonnet_nonexistent_import_crazy(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/nonexistent_import_crazy.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/nonexistent_import_crazy.golden
	expected := "RUNTIME ERROR: couldn't open import \"ąęółńśćźż \\\" ' \\n\\n\\t\\t\": no match locally or in the Jsonnet library paths\n\ttestdata/nonexistent_import_crazy.jsonnet:1:1-46\t$\n\tDuring evaluation"
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/nonexistent_import_crazy.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
	})
}

func Test_go_jsonnet_number_divided_by_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/number_divided_by_string.jsonnet")
	expected := `unhandled binary operation number / string
	testdata/number_divided_by_string.jsonnet:1:1-11	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/number_divided_by_string.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_number_leading_zero(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/number_leading_zero.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/number_leading_zero.golden
	expected := `RUNTIME ERROR: testdata/number_leading_zero.jsonnet:1:2-4 Did not expect: (NUMBER, "42")
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/number_leading_zero.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_number_times_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/number_times_string.jsonnet")
	expected := `unhandled binary operation number * string
	testdata/number_times_string.jsonnet:1:1-11	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/number_times_string.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_numeric_literal(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/numeric_literal.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/numeric_literal.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/numeric_literal.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_obj_local_right_level(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/obj_local_right_level.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/obj_local_right_level.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/obj_local_right_level.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_obj_local_right_level2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/obj_local_right_level2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/obj_local_right_level2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/obj_local_right_level2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_obj_local_right_level3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/obj_local_right_level3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/obj_local_right_level3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/obj_local_right_level3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_comp(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_comp.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_comp2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_comp2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_comp3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_comp3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_comp4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_comp4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_comp_assert(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_assert.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/object_comp_assert.golden
	expected := `RUNTIME ERROR: testdata/object_comp_assert.jsonnet:1:32-35 Object comprehension cannot have asserts
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_comp_assert.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_comp_bad_field(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_bad_field.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/object_comp_bad_field.golden
	expected := `RUNTIME ERROR: testdata/object_comp_bad_field.jsonnet:1:9-12 Object comprehensions can only have [e] fields
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_comp_bad_field.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_comp_bad_field2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_bad_field2.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/object_comp_bad_field2.golden
	expected := `RUNTIME ERROR: testdata/object_comp_bad_field2.jsonnet:1:11-14 Object comprehensions can only have [e] fields
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_comp_bad_field2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_comp_dollar(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_dollar.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_dollar.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_comp_dollar.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_comp_dollar2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_dollar2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_dollar2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_comp_dollar2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_comp_dollar3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_dollar3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_dollar3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_comp_dollar3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_comp_duplicate(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_duplicate.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/object_comp_duplicate.golden
	expected := `RUNTIME ERROR: Duplicate field name: "x"
	testdata/object_comp_duplicate.jsonnet:1:1-31
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_comp_duplicate.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_comp_err_elem(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_err_elem.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/object_comp_err_elem.golden
	expected := `RUNTIME ERROR: xxx
	testdata/object_comp_err_elem.jsonnet:1:11-22	object <anonymous>
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_comp_err_elem.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_comp_err_index(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_err_index.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/object_comp_err_index.golden
	expected := `RUNTIME ERROR: xxx
	testdata/object_comp_err_index.jsonnet:1:4-15	$
	testdata/object_comp_err_index.jsonnet:1:1-35
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_comp_err_index.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_comp_if(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_if.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_if.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_comp_if.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_comp_illegal(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_illegal.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/object_comp_illegal.golden
	expected := `RUNTIME ERROR: testdata/object_comp_illegal.jsonnet:1:15-18 Object comprehension can only have one field
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_comp_illegal.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_comp_int_index(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_int_index.jsonnet")
	expected := `unexpected field name type number, expected string
	testdata/object_comp_int_index.jsonnet:1:1-30
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_comp_int_index.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_comp_local(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_local.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_local.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_comp_local.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_comp_local2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_local2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_local2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_comp_local2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_comp_local3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_local3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_local3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_comp_local3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_comp_super(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_super.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_super.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_comp_super.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_comp_try_iterate_over_obj(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_try_iterate_over_obj.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/object_comp_try_iterate_over_obj.golden
	expected := `RUNTIME ERROR: Unexpected type object, expected array
	testdata/object_comp_try_iterate_over_obj.jsonnet:1:1-23
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_comp_try_iterate_over_obj.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_comp_try_iterate_over_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_comp_try_iterate_over_string.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/object_comp_try_iterate_over_string.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected array
	testdata/object_comp_try_iterate_over_string.jsonnet:1:1-24
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_comp_try_iterate_over_string.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_hidden(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_hidden.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_hidden.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_hidden.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_implicit_plus1(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_implicit_plus1.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_implicit_plus1.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_implicit_plus1.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_invariant(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_invariant.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_invariant10(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant10.jsonnet")
	expected := `RUNTIME ERROR: Object assertion failed.
	testdata/object_invariant10.jsonnet:1:16-28
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_invariant10.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_invariant11(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant11.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/object_invariant11.golden
	expected := `RUNTIME ERROR: Object assertion failed.
	testdata/object_invariant11.jsonnet:1:3-15
	testdata/object_invariant11.jsonnet:1:1-19	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_invariant11.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_invariant12(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant12.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant12.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_invariant12.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_invariant13(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant13.jsonnet")
	expected := `RUNTIME ERROR: x
	testdata/object_invariant13.jsonnet:1:10-19	object <anonymous>
	testdata/object_invariant13.jsonnet:1:3-19
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_invariant13.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_invariant14(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant14.jsonnet")
	expected := `RUNTIME ERROR: xxx
	testdata/object_invariant14.jsonnet:1:3-22
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_invariant14.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_invariant2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant2.jsonnet")
	expected := `RUNTIME ERROR: Object assertion failed.
	testdata/object_invariant2.jsonnet:1:3-15
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_invariant2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_invariant3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_invariant3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_invariant4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_invariant4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_invariant5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_invariant5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_invariant6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_invariant6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_invariant7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant7.jsonnet")
	expected := `RUNTIME ERROR: Field does not exist: x
	testdata/object_invariant7.jsonnet:1:16-21	object <anonymous>
	testdata/object_invariant7.jsonnet:1:9-28
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_invariant7.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_invariant8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant8.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/object_invariant8.golden
	expected := `RUNTIME ERROR: Object assertion failed.
	testdata/object_invariant8.jsonnet:1:9-27
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_invariant8.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_invariant9(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant9.jsonnet")
	expected := `RUNTIME ERROR: Object assertion failed.
	testdata/object_invariant9.jsonnet:1:16-28
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_invariant9.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_invariant_perf(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant_perf.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant_perf.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_invariant_perf.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_invariant_plus(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant_plus.jsonnet")
	expected := `RUNTIME ERROR: Object assertion failed.
	testdata/object_invariant_plus.jsonnet:1:2-14
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_invariant_plus.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_invariant_plus2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant_plus2.jsonnet")
	expected := `RUNTIME ERROR: Object assertion failed.
	testdata/object_invariant_plus2.jsonnet:1:18-30
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_invariant_plus2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_invariant_plus3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant_plus3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant_plus3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_invariant_plus3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_invariant_plus4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant_plus4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant_plus4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_invariant_plus4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_invariant_plus5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant_plus5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant_plus5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_invariant_plus5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_invariant_plus6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant_plus6.jsonnet")
	expected := `RUNTIME ERROR: yyy
	testdata/object_invariant_plus6.jsonnet:1:27-46
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_invariant_plus6.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_invariant_plus7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant_plus7.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_invariant_plus7.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_invariant_plus7.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_literal_in_array_comp(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_literal_in_array_comp.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_literal_in_array_comp.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_literal_in_array_comp.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_literal_in_object_comp(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_literal_in_object_comp.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_literal_in_object_comp.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_literal_in_object_comp.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_local(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_local.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_local.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_local.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_local_from_parent(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_local_from_parent.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_local_from_parent.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_local_from_parent.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_local_from_parent_through_local(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_local_from_parent_through_local.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_local_from_parent_through_local.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_local_from_parent_through_local.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_local_recursive(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_local_recursive.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_local_recursive.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_local_recursive.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_local_self_super(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_local_self_super.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_local_self_super.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_local_self_super.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_local_uses_local_from_outside(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_local_uses_local_from_outside.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_local_uses_local_from_outside.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_local_uses_local_from_outside.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_no_newline(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_no_newline.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_no_newline.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:   "resources/go-jsonnet",
		File:      "testdata/object_no_newline.jsonnet",
		Snippet:   snippet,
		Expected:  expected,
		ExtVars:   extVars,
		ExtCodes:  extCodes,
		NoNewline: true,
	})
}

func Test_go_jsonnet_object_plus_bad(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_plus_bad.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/object_plus_bad.golden
	expected := `RUNTIME ERROR: Unexpected type number, expected object
	testdata/object_plus_bad.jsonnet:1:1-8	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/object_plus_bad.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_object_sum(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_sum.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_sum.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_sum.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_sum2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_sum2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_sum2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_sum2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_sum3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_sum3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_sum3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_sum3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_super(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_super.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_super.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_super.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_super_deep(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_super_deep.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_super_deep.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_super_deep.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_super_within(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_super_within.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_super_within.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_super_within.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_various_field_types(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_various_field_types.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_various_field_types.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_various_field_types.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_object_within_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/object_within_object.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/object_within_object.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/object_within_object.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args10(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args10.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args10.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args10.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args11(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args11.jsonnet")
	expected := `RUNTIME ERROR: Argument x already provided
	testdata/optional_args11.jsonnet:1:1-30	$
	During evaluation`
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/optional_args11.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
	})
}

func Test_go_jsonnet_optional_args12(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args12.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args12.golden")
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args12.jsonnet",
		Snippet:  snippet,
		Expected: expected,
	})
}

func Test_go_jsonnet_optional_args13(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args13.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/optional_args13.golden
	expected := `RUNTIME ERROR: Missing argument: y
	testdata/optional_args13.jsonnet:1:1-26	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/optional_args13.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_optional_args14(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args14.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args14.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args14.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args15(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args15.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args15.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args15.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args16(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args16.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args16.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args16.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args17(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args17.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args17.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args17.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args18(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args18.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args18.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args18.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args19(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args19.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args19.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args19.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args20(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args20.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args20.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args20.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args21(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args21.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args21.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args21.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args22(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args22.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args22.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args22.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args7.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args7.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/optional_args7.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_optional_args8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args8.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/optional_args8.golden
	expected := `RUNTIME ERROR: function has no parameter y
	testdata/optional_args8.jsonnet:2:1-10	$
	During evaluation`
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/optional_args8.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
	})
}

func Test_go_jsonnet_optional_args9(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/optional_args9.jsonnet")
	expected := `RUNTIME ERROR: Argument x already provided
	testdata/optional_args9.jsonnet:1:1-26	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/optional_args9.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_or(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/or.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/or.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/or.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_or2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/or2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/or2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/or2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_or3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/or3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/or3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/or3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_or4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/or4.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/or4.golden
	expected := `RUNTIME ERROR: xxx
	testdata/or4.jsonnet:1:10-21	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/or4.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_or5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/or5.jsonnet")
	expected := `unexpected type string for || op, expected boolean
	testdata/or5.jsonnet:1:1-14	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/or5.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_or6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/or6.jsonnet")
	expected := `unexpected type string for || op, expected boolean
	testdata/or6.jsonnet:1:1-15	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/or6.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_overriding_stdlib_desugared(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/overriding_stdlib_desugared.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/overriding_stdlib_desugared.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/overriding_stdlib_desugared.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_parseJson(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/parseJson.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/parseJson.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/parseJson.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_parseYaml(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/parseYaml.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/parseYaml.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/parseYaml.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_percent_bad(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_bad.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/percent_bad.golden
	expected := `RUNTIME ERROR: Operator % cannot be used on types number and string.
	testdata/percent_bad.jsonnet:1:1-9
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/percent_bad.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_percent_bad2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_bad2.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/percent_bad2.golden
	expected := `not all arguments converted during string formatting
	testdata/percent_bad2.jsonnet:1:1-9
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/percent_bad2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_percent_bad3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_bad3.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/percent_bad3.golden
	expected := `RUNTIME ERROR: Operator % cannot be used on types function and number.
	testdata/percent_bad3.jsonnet:1:1-21
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/percent_bad3.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_percent_format_float(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_format_float.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/percent_format_float.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/percent_format_float.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_percent_format_str(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_format_str.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/percent_format_str.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/percent_format_str.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_percent_format_str2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_format_str2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/percent_format_str2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/percent_format_str2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_percent_format_str3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_format_str3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/percent_format_str3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/percent_format_str3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_percent_format_str4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_format_str4.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/percent_format_str4.golden
	expected := `not all arguments converted during string formatting
	testdata/percent_format_str4.jsonnet:1:1-20
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/percent_format_str4.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_percent_format_str5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_format_str5.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/percent_format_str5.golden
	expected := `not enough arguments for format string
	testdata/percent_format_str5.jsonnet:1:1-18
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/percent_format_str5.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_percent_format_str6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_format_str6.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/percent_format_str6.golden
	expected := `not enough arguments for format string
	testdata/percent_format_str6.jsonnet:1:1-16
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/percent_format_str6.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_percent_format_str7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_format_str7.jsonnet")
	expected := `format %d requires number
	testdata/percent_format_str7.jsonnet:1:1-15
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/percent_format_str7.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_percent_format_str8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_format_str8.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/percent_format_str8.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/percent_format_str8.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_percent_mod_int(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_mod_int.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/percent_mod_int.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/percent_mod_int.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_percent_mod_int2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_mod_int2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/percent_mod_int2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/percent_mod_int2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_percent_mod_int3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_mod_int3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/percent_mod_int3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/percent_mod_int3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_percent_mod_int4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_mod_int4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/percent_mod_int4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/percent_mod_int4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_percent_mod_int5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_mod_int5.jsonnet")
	expected := `RUNTIME ERROR: Division by zero.
	testdata/percent_mod_int5.jsonnet:1:1-7
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/percent_mod_int5.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_percent_mod_int6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/percent_mod_int6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/percent_mod_int6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/percent_mod_int6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_plus(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/plus.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/plus.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/plus.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_plus2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/plus2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/plus2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/plus2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_plus3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/plus3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/plus3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/plus3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_plus4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/plus4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/plus4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/plus4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_plus5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/plus5.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/plus5.golden
	expected := `RUNTIME ERROR: Unexpected type function, expected number
	testdata/plus5.jsonnet:1:1-19	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/plus5.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_plus6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/plus6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/plus6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/plus6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_plus7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/plus7.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/plus7.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/plus7.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_plus8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/plus8.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/plus8.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/plus8.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_plus9(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/plus9.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/plus9.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/plus9.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_positional_after_optional(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/positional_after_optional.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/positional_after_optional.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/positional_after_optional.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_pow(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/pow.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/pow.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/pow.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_pow2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/pow2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/pow2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/pow2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_pow3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/pow3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/pow3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/pow3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_pow4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/pow4.jsonnet")
	expected := `RUNTIME ERROR: Not a number
	testdata/pow4.jsonnet:1:1-17	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/pow4.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_pow5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/pow5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/pow5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/pow5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_pow6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/pow6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/pow6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/pow6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_pow7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/pow7.jsonnet")
	expected := `RUNTIME ERROR: Overflow
	testdata/pow7.jsonnet:2:1-23	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/pow7.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_pow8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/pow8.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/pow8.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected number
	testdata/pow8.jsonnet:1:1-19	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/pow8.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_pow9(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/pow9.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/pow9.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected number
	testdata/pow9.jsonnet:1:1-19	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/pow9.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_proto_object_comp(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/proto_object_comp.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/proto_object_comp.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/proto_object_comp.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_recursive_local(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/recursive_local.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/recursive_local.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/recursive_local.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_recursive_thunk(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/recursive_thunk.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/recursive_thunk.golden
	expected := `RUNTIME ERROR: xxx
	testdata/recursive_thunk.jsonnet:1:35-46	function <bar>
	testdata/recursive_thunk.jsonnet:2:16-38	function <foo>
	testdata/recursive_thunk.jsonnet:1:52-54	function <bar>
	testdata/recursive_thunk.jsonnet:2:16-38	function <foo>
	testdata/recursive_thunk.jsonnet:1:52-54	function <bar>
	testdata/recursive_thunk.jsonnet:2:16-38	function <foo>
	testdata/recursive_thunk.jsonnet:3:1-7	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/recursive_thunk.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_self(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/self.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/self.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/self.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_simple_arith1(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/simple_arith1.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/simple_arith1.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/simple_arith1.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_simple_arith2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/simple_arith2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/simple_arith2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/simple_arith2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_simple_arith3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/simple_arith3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/simple_arith3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/simple_arith3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_simple_arith_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/simple_arith_string.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/simple_arith_string.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/simple_arith_string.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_simple_arith_string2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/simple_arith_string2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/simple_arith_string2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/simple_arith_string2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_simple_arith_string3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/simple_arith_string3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/simple_arith_string3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/simple_arith_string3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_simple_arith_string_empty(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/simple_arith_string_empty.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/simple_arith_string_empty.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/simple_arith_string_empty.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_slice(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/slice.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/slice.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/slice.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_slice2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/slice2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/slice2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/slice2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_slice3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/slice3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/slice3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/slice3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_slice4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/slice4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/slice4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/slice4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_slice5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/slice5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/slice5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/slice5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_slice6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/slice6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/slice6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/slice6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_slice7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/slice7.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/slice7.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/slice7.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_stackbug_regression_test(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/stackbug-regression-test.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/stackbug-regression-test.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/stackbug-regression-test.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_stacktrace_assert(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/stacktrace_assert.jsonnet")
	expected := `RUNTIME ERROR: Object assertion failed.
	testdata/stacktrace_assert.jsonnet:1:3-15
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/stacktrace_assert.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_stacktrace_plussuper(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/stacktrace_plussuper.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type object, expected null	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/stacktrace_plussuper.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_static_error_eof(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/static_error_eof.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/static_error_eof.golden
	expected := `RUNTIME ERROR: testdata/static_error_eof.jsonnet:2:1 Expected , or ; but got end of file
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/static_error_eof.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_codepoint(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.codepoint.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.codepoint.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.codepoint.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_codepoint2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.codepoint2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.codepoint2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.codepoint2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_codepoint3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.codepoint3.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.codepoint3.golden
	expected := `RUNTIME ERROR: codepoint takes a string of length 1, got length 2
	testdata/std.codepoint3.jsonnet:1:1-20	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.codepoint3.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_codepoint4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.codepoint4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.codepoint4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.codepoint4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_codepoint5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.codepoint5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.codepoint5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.codepoint5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_codepoint6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.codepoint6.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.codepoint6.golden
	expected := `RUNTIME ERROR: codepoint takes a string of length 1, got length 0
	testdata/std.codepoint6.jsonnet:1:1-18	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.codepoint6.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_codepoint7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.codepoint7.jsonnet")
	expected := `RUNTIME ERROR: codepoint takes a string of length 1, got length 2
	testdata/std.codepoint7.jsonnet:2:1-21	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.codepoint7.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_codepoint8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.codepoint8.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.codepoint8.golden
	expected := `RUNTIME ERROR: Unexpected type number, expected string
	testdata/std.codepoint8.jsonnet:1:1-18	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.codepoint8.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_exponent(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.exponent.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.exponent.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.exponent.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_exponent2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.exponent2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.exponent2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.exponent2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_exponent3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.exponent3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.exponent3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.exponent3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_exponent4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.exponent4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.exponent4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.exponent4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_exponent5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.exponent5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.exponent5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.exponent5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_exponent6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.exponent6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.exponent6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.exponent6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_exponent7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.exponent7.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.exponent7.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.exponent7.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_filter(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.filter.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.filter.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.filter.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_filter2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.filter2.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.filter2.golden
	expected := `RUNTIME ERROR: x
	testdata/std.filter2.jsonnet:1:1-26	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.filter2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_filter3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.filter3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.filter3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.filter3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_filter4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.filter4.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.filter4.golden
	expected := `RUNTIME ERROR: Unexpected type number, expected function
	testdata/std.filter4.jsonnet:1:1-19	$
	During evaluation`
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.filter4.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
	})
}

func Test_go_jsonnet_std_filter5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.filter5.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.filter5.golden
	expected := `RUNTIME ERROR: Unexpected type number, expected array
	testdata/std.filter5.jsonnet:1:1-31	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.filter5.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_filter6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.filter6.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.filter6.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected array
	testdata/std.filter6.jsonnet:1:1-21	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.filter6.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_filter7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.filter7.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.filter7.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.filter7.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_filter8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.filter8.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.filter8.golden
	expected := `RUNTIME ERROR: Unexpected type function, expected array
	testdata/std.filter8.jsonnet:1:1-36	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.filter8.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_filter_swapped_args(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.filter_swapped_args.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.filter_swapped_args.golden
	expected := `RUNTIME ERROR: Unexpected type function, expected array
	testdata/std.filter_swapped_args.jsonnet:1:1-38	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.filter_swapped_args.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_flatmap(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.flatmap.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.flatmap.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.flatmap.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_flatmap2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.flatmap2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.flatmap2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.flatmap2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_flatmap3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.flatmap3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.flatmap3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.flatmap3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_flatmap4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.flatmap4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.flatmap4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.flatmap4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_flatmap5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.flatmap5.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.flatmap5.golden
	expected := `RUNTIME ERROR: a
	testdata/std.flatmap5.jsonnet:1:21-28	function <failWith>
	testdata/std.flatmap5.jsonnet:2:1-49	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.flatmap5.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_flatmap6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.flatmap6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.flatmap6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.flatmap6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_join(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.join.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.join.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.join.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_join2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.join2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.join2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.join2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_join3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.join3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.join3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.join3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_join4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.join4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.join4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.join4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_join5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.join5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.join5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.join5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_join6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.join6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.join6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.join6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_join7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.join7.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.join7.golden
	expected := `RUNTIME ERROR: Unexpected type array, expected string
	testdata/std.join7.jsonnet:1:1-27	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.join7.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_join8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.join8.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.join8.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected array
	testdata/std.join8.jsonnet:1:1-33	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.join8.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_length(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.length.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.length.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.length.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_length_array(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.length_array.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.length_array.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.length_array.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_length_function(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.length_function.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.length_function.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.length_function.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_length_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.length_object.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.length_object.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.length_object.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_length_object_sum(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.length_object_sum.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.length_object_sum.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.length_object_sum.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_length_object_with_hidden(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.length_object_with_hidden.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.length_object_with_hidden.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.length_object_with_hidden.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_length_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.length_string.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.length_string.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.length_string.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_lstripChars_multibyte(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.lstripChars.multibyte.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.lstripChars.multibyte.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.lstripChars.multibyte.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_makeArray(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.makeArray.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.makeArray.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.makeArray.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_makeArrayNamed(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.makeArrayNamed.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.makeArrayNamed.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.makeArrayNamed.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_makeArrayNamed2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.makeArrayNamed2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.makeArrayNamed2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.makeArrayNamed2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_makeArrayNamed3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.makeArrayNamed3.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.makeArrayNamed3.golden
	expected := `RUNTIME ERROR: function has no parameter blahblah
	testdata/std.makeArrayNamed3.jsonnet:1:1-54	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.makeArrayNamed3.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_makeArrayNamed4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.makeArrayNamed4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.makeArrayNamed4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.makeArrayNamed4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_makeArray_bad(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.makeArray_bad.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.makeArray_bad.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected number
	testdata/std.makeArray_bad.jsonnet:1:1-36	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.makeArray_bad.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_makeArray_bad2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.makeArray_bad2.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.makeArray_bad2.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected function
	testdata/std.makeArray_bad2.jsonnet:1:1-25	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.makeArray_bad2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_makeArray_noninteger(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.makeArray_noninteger.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.makeArray_noninteger.golden
	expected := `RUNTIME ERROR: Expected an integer, but got 2.5
	testdata/std.makeArray_noninteger.jsonnet:1:1-34	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.makeArray_noninteger.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_makeArray_noninteger_big(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.makeArray_noninteger_big.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.makeArray_noninteger_big.golden
	expected := `RUNTIME ERROR: Expected an integer, but got 1e+100
	testdata/std.makeArray_noninteger_big.jsonnet:1:1-36	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.makeArray_noninteger_big.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_makeArray_recursive(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.makeArray_recursive.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.makeArray_recursive.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.makeArray_recursive.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_makeArray_recursive_evalutation_order_matters(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.makeArray_recursive_evalutation_order_matters.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.makeArray_recursive_evalutation_order_matters.golden
	// expected := `RUNTIME ERROR: max stack frames exceeded.`
	expected := "500\n"
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.makeArray_recursive_evalutation_order_matters.jsonnet",
		Snippet:  snippet,
		Expected: expected,
	})
}

func Test_go_jsonnet_std_manifestYamlDoc_error(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.manifestYamlDoc_error.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.manifestYamlDoc_error.golden
	expected := `RUNTIME ERROR: foo
	testdata/std.manifestYamlDoc_error.jsonnet:1:31-42	object <anonymous>
	testdata/std.manifestYamlDoc_error.jsonnet:1:1-47	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.manifestYamlDoc_error.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_manifestYamlDoc_ok(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.manifestYamlDoc_ok.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.manifestYamlDoc_ok.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.manifestYamlDoc_ok.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_mantissa(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.mantissa.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.mantissa.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.mantissa.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_mantissa2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.mantissa2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.mantissa2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.mantissa2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_mantissa3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.mantissa3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.mantissa3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.mantissa3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_mantissa4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.mantissa4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.mantissa4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.mantissa4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_mantissa5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.mantissa5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.mantissa5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.mantissa5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_mantissa6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.mantissa6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.mantissa6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.mantissa6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_mantissa7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.mantissa7.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.mantissa7.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.mantissa7.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_maxArray(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.maxArray.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.maxArray.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.maxArray.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_maxArrayKeyF(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.maxArrayKeyF.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.maxArrayKeyF.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.maxArrayKeyF.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_maxArrayOnEmpty(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.maxArrayOnEmpty.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.maxArrayOnEmpty.golden
	expected := `Expected at least one element in array. Got none
	testdata/std.maxArrayOnEmpty.jsonnet:1:1-17	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.maxArrayOnEmpty.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_maxArrayOnEmpty2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.maxArrayOnEmpty2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.maxArrayOnEmpty2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.maxArrayOnEmpty2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_md5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.md5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.md5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.md5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_md5_2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.md5_2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.md5_2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.md5_2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_md5_3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.md5_3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.md5_3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.md5_3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_md5_4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.md5_4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.md5_4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.md5_4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_md5_5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.md5_5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.md5_5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.md5_5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_md5_6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.md5_6.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.md5_6.golden
	expected := `RUNTIME ERROR: Unexpected type number, expected string
	testdata/std.md5_6.jsonnet:1:1-12	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.md5_6.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_minArray(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.minArray.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.minArray.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.minArray.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_minArrayKeyF(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.minArrayKeyF.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.minArrayKeyF.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.minArrayKeyF.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_minArrayOnEmpty(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.minArrayOnEmpty.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.minArrayOnEmpty.golden
	expected := `Expected at least one element in array. Got none
	testdata/std.minArrayOnEmpty.jsonnet:1:1-17	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.minArrayOnEmpty.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_minArrayOnEmpty2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.minArrayOnEmpty2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.minArrayOnEmpty2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.minArrayOnEmpty2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_mod_int(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.mod_int.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.mod_int.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.mod_int.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_mod_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.mod_string.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.mod_string.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.mod_string.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_modulo(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.modulo.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.modulo.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.modulo.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_modulo2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.modulo2.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.modulo2.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected number
	testdata/std.modulo2.jsonnet:1:1-22	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.modulo2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_modulo3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.modulo3.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.modulo3.golden
	expected := `RUNTIME ERROR: Unexpected type string, expected number
	testdata/std.modulo3.jsonnet:1:1-22	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.modulo3.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_objectFields(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.objectFields.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.objectFields.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.objectFields.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_objectHasEx(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.objectHasEx.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.objectHasEx.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.objectHasEx.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_objectHasEx2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.objectHasEx2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.objectHasEx2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.objectHasEx2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_objectHasEx3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.objectHasEx3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.objectHasEx3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.objectHasEx3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_objectHasEx4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.objectHasEx4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.objectHasEx4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.objectHasEx4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_objectHasEx5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.objectHasEx5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.objectHasEx5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.objectHasEx5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.primitiveEquals.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals10(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals10.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.primitiveEquals10.golden
	expected := `RUNTIME ERROR: x
	testdata/std.primitiveEquals10.jsonnet:1:1-35	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.primitiveEquals10.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals11(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals11.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals11.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.primitiveEquals11.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals12(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals12.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals12.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.primitiveEquals12.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals13(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals13.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.primitiveEquals13.golden
	expected := `RUNTIME ERROR: primitiveEquals operates on primitive types, got array
	testdata/std.primitiveEquals13.jsonnet:1:1-28	$
	During evaluation`
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.primitiveEquals13.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
	})
}

func Test_go_jsonnet_std_primitiveEquals14(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals14.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals14.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.primitiveEquals14.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals15(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals15.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals15.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.primitiveEquals15.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals16(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals16.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals16.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.primitiveEquals16.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals17(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals17.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals17.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.primitiveEquals17.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals18(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals18.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals18.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.primitiveEquals18.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals19(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals19.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals19.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.primitiveEquals19.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.primitiveEquals2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals20(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals20.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals20.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.primitiveEquals20.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals21(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals21.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals21.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.primitiveEquals21.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.primitiveEquals3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.primitiveEquals4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.primitiveEquals5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals6.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.primitiveEquals6.golden
	expected := `RUNTIME ERROR: primitiveEquals operates on primitive types, got object
	testdata/std.primitiveEquals6.jsonnet:1:1-28	$
	During evaluation`
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.primitiveEquals6.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
	})
}

func Test_go_jsonnet_std_primitiveEquals7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals7.jsonnet")
	expected := `RUNTIME ERROR: primitiveEquals operates on primitive types, got function
	testdata/std.primitiveEquals7.jsonnet:1:1-50	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.primitiveEquals7.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals8.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals8.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.primitiveEquals8.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_primitiveEquals9(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.primitiveEquals9.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.primitiveEquals9.golden
	expected := `RUNTIME ERROR: x
	testdata/std.primitiveEquals9.jsonnet:1:1-35	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.primitiveEquals9.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_rstripChars_multibyte(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.rstripChars.multibyte.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.rstripChars.multibyte.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.rstripChars.multibyte.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_slice(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.slice.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.slice.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.slice.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_sort(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.sort.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.sort.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.sort.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_sort2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.sort2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.sort2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.sort2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_sort3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.sort3.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.sort3.golden
	expected := `RUNTIME ERROR: foo
	testdata/std.sort3.jsonnet:1:1-29	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.sort3.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_sort4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.sort4.jsonnet")
	expected := `unexpected type array, expected number
	testdata/std.sort4.jsonnet:1:1-29	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.sort4.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_thisFile(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.thisFile.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.thisFile.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)

	tmpFile := filepath.Join(os.TempDir(), "testdata", "std.thisFile")
	_ = os.MkdirAll(filepath.Dir(tmpFile), 0755)
	_ = os.WriteFile(tmpFile, []byte(snippet), 0644)
	defer os.Remove(tmpFile)

	runTest(t, TestConfig{
		WorkDir:  os.TempDir(),
		File:     "testdata/std.thisFile",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_thisFile2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.thisFile2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.thisFile2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.thisFile2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_toString(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.toString.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.toString.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.toString.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_toString2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.toString2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.toString2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.toString2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_toString3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.toString3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.toString3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.toString3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_toString4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.toString4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.toString4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.toString4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_toString5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.toString5.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/std.toString5.golden
	expected := `RUNTIME ERROR: x
	testdata/std.toString5.jsonnet:1:1-24	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/std.toString5.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_std_toString6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.toString6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.toString6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.toString6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_toString7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.toString7.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.toString7.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.toString7.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_toString8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std.toString8.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std.toString8.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std.toString8.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_in_local(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std_in_local.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std_in_local.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std_in_local.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_std_substr(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/std_substr.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/std_substr.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/std_substr.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_stdlib_smoke_test(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/stdlib_smoke_test.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/stdlib_smoke_test.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)

	tmpFile := filepath.Join(os.TempDir(), "testdata", "stdlib_smoke_test")
	_ = os.MkdirAll(filepath.Dir(tmpFile), 0755)
	_ = os.WriteFile(tmpFile, []byte(snippet), 0644)
	defer os.Remove(tmpFile)

	runTest(t, TestConfig{
		WorkDir:  os.TempDir(),
		File:     "testdata/stdlib_smoke_test",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_strReplace(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/strReplace.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/strReplace.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/strReplace.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_strReplace2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/strReplace2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/strReplace2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/strReplace2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_strReplace3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/strReplace3.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/strReplace3.golden
	expected := `RUNTIME ERROR: 'from' string must not be zero length.
	testdata/strReplace3.jsonnet:1:1-35	$
	During evaluation`
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/strReplace3.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
	})
}

func Test_go_jsonnet_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/string.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/string.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_string2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/string2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/string2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_string_comparison1(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string_comparison1.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/string_comparison1.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/string_comparison1.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_string_comparison2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string_comparison2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/string_comparison2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/string_comparison2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_string_comparison3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string_comparison3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/string_comparison3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/string_comparison3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_string_comparison4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string_comparison4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/string_comparison4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/string_comparison4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_string_comparison5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string_comparison5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/string_comparison5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/string_comparison5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_string_comparison6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string_comparison6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/string_comparison6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/string_comparison6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_string_comparison7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string_comparison7.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/string_comparison7.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/string_comparison7.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_string_divided_by_number(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string_divided_by_number.jsonnet")
	expected := `unhandled binary operation string / number
	testdata/string_divided_by_number.jsonnet:1:1-11	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/string_divided_by_number.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_string_index(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string_index.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/string_index.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/string_index.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_string_index2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string_index2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/string_index2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/string_index2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_string_index_negative(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string_index_negative.jsonnet")
	expected := `RUNTIME ERROR: Index -1 out of bounds, not within [0, 4)
	testdata/string_index_negative.jsonnet:1:1-11	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/string_index_negative.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_string_index_out_of_bounds(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string_index_out_of_bounds.jsonnet")
	expected := `RUNTIME ERROR: Index 4 out of bounds, not within [0, 4)
	testdata/string_index_out_of_bounds.jsonnet:1:1-10	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/string_index_out_of_bounds.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_string_minus_number(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string_minus_number.jsonnet")
	expected := `unhandled binary operation string - number
	testdata/string_minus_number.jsonnet:1:1-9	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/string_minus_number.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_string_plus_function(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string_plus_function.jsonnet")
	expected := `unhandled type function, string conversion not available
	testdata/string_plus_function.jsonnet:1:1-24	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/string_plus_function.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_string_times_number(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string_times_number.jsonnet")
	expected := `unhandled binary operation string * number
	testdata/string_times_number.jsonnet:1:1-9	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/string_times_number.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_string_to_bool(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/string_to_bool.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/string_to_bool.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/string_to_bool.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_super_index_desugar(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/super_index_desugar.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/super_index_desugar.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/super_index_desugar.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_supersugar(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/supersugar.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_supersugar2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/supersugar2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_supersugar3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/supersugar3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_supersugar4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/supersugar4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_supersugar5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/supersugar5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_supersugar6(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar6.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar6.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/supersugar6.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_supersugar7(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar7.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar7.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/supersugar7.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_supersugar8(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar8.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/supersugar8.golden
	expected := `RUNTIME ERROR: Object assertion failed.
	testdata/supersugar8.jsonnet:1:3-16
	During manifestation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/supersugar8.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_supersugar9(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar9.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/supersugar9.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/supersugar9.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_syntax_error(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/syntax_error.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/syntax_error.golden
	expected := `RUNTIME ERROR: testdata/syntax_error.jsonnet:2:1 Unexpected end of file
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/syntax_error.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_tailstrict(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/tailstrict.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/tailstrict.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/tailstrict.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_tailstrict2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/tailstrict2.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/tailstrict2.golden
	expected := `RUNTIME ERROR: xxx
	testdata/tailstrict2.jsonnet:1:13-20	function <e>
	testdata/tailstrict2.jsonnet:2:14-18	function <anonymous>
	testdata/tailstrict2.jsonnet:2:1-26	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/tailstrict2.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_tailstrict3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/tailstrict3.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/tailstrict3.golden
	expected := `RUNTIME ERROR: xxx
	testdata/tailstrict3.jsonnet:2:1-8	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/tailstrict3.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_tailstrict4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/tailstrict4.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/tailstrict4.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/tailstrict4.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_tailstrict5(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/tailstrict5.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/tailstrict5.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/tailstrict5.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_tailstrict_operator1(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/tailstrict_operator1.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/tailstrict_operator1.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/tailstrict_operator1.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_tailstrict_operator2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/tailstrict_operator2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/tailstrict_operator2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/tailstrict_operator2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_tailstrict_operator3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/tailstrict_operator3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/tailstrict_operator3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/tailstrict_operator3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_too_many_arguments(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/too_many_arguments.jsonnet")
	expected := `RUNTIME ERROR: function expected 3 positional argument(s), but got 4
	testdata/too_many_arguments.jsonnet:1:1-35	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/too_many_arguments.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_true(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/true.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/true.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/true.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_type_array(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/type_array.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/type_array.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/type_array.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_type_builtin_function(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/type_builtin_function.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/type_builtin_function.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/type_builtin_function.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_type_error(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/type_error.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/type_error.golden
	expected := `RUNTIME ERROR: xxx
	testdata/type_error.jsonnet:1:1-22	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/type_error.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_type_function(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/type_function.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/type_function.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/type_function.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_type_number(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/type_number.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/type_number.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/type_number.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_type_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/type_object.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/type_object.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/type_object.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_type_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/type_string.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/type_string.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/type_string.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_unary_minus(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/unary_minus.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/unary_minus.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/unary_minus.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_unary_minus2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/unary_minus2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/unary_minus2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/unary_minus2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_unary_minus3(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/unary_minus3.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/unary_minus3.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/unary_minus3.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_unary_minus4(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/unary_minus4.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type string, expected number
	testdata/unary_minus4.jsonnet:1:1-7	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/unary_minus4.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_unary_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/unary_object.jsonnet")
	expected := `RUNTIME ERROR: Unexpected type object, expected number
	testdata/unary_object.jsonnet:1:1-5	$
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/unary_object.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_unfinished_args(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/unfinished_args.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/unfinished_args.golden
	expected := `RUNTIME ERROR: testdata/unfinished_args.jsonnet:2:1 Expected a comma before next function argument, got end of file
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/unfinished_args.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_unicode(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/unicode.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/unicode.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/unicode.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_unicode2(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/unicode2.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/unicode2.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/unicode2.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_use_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/use_object.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/use_object.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/use_object.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_use_object_in_object(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/use_object_in_object.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/use_object_in_object.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/use_object_in_object.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_variable(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/variable.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/variable.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/variable.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}

func Test_go_jsonnet_variable_not_visible(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/variable_not_visible.jsonnet")
	// Expected file: resources/go-jsonnet/testdata/variable_not_visible.golden
	expected := `RUNTIME ERROR: testdata/variable_not_visible.jsonnet:1:44-50 Unknown variable: nested
	During evaluation`
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:     "resources/go-jsonnet",
		File:        "testdata/variable_not_visible.jsonnet",
		Snippet:     snippet,
		ExpectedErr: expected,
		ExtVars:     extVars,
		ExtCodes:    extCodes,
	})
}

func Test_go_jsonnet_verbatim_string(t *testing.T) {
	snippet := mustReadFile(t, "resources/go-jsonnet/testdata/verbatim_string.jsonnet")
	expected := mustReadFile(t, "resources/go-jsonnet/testdata/verbatim_string.golden")
	extVars := map[string]string(nil)
	extCodes := map[string]string(nil)
	runTest(t, TestConfig{
		WorkDir:  "resources/go-jsonnet",
		File:     "testdata/verbatim_string.jsonnet",
		Snippet:  snippet,
		Expected: expected,
		ExtVars:  extVars,
		ExtCodes: extCodes,
	})
}
