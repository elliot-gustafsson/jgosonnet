package stdlib

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

const codepointMax = 0x10FFFF

func liftString(f func(string) string, name string) evaluator.Func {
	return func(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
		if len(args) != 1 {
			return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to %s: %d, expected 1", name, len(args))
		}
		if !args[0].IsString() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to %s (arg 0): %s, expected string", name, args[0].Type().String())
		}
		res := f(args[0].String(ctx))
		return evaluator.MakeString(res, ctx), nil
	}
}

func liftString2(f func(string, string) string, name string) evaluator.Func {
	return func(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
		if len(args) != 2 {
			return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to %s: %d, expected 2", name, len(args))
		}
		if !args[0].IsString() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to %s (arg 0): %s, expected string", name, args[0].Type().String())
		}
		if !args[1].IsString() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to %s (arg 1): %s, expected string", name, args[1].Type().String())
		}
		res := f(args[0].String(ctx), args[1].String(ctx))
		return evaluator.MakeString(res, ctx), nil
	}
}

func liftStringToValueErr(f func(string) (evaluator.Value, error), name string) evaluator.Func {
	return func(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
		if len(args) != 1 {
			return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to %s: %d, expected 1", name, len(args))
		}
		if !args[0].IsString() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to %s (arg 0): %s, expected string", name, args[0].Type().String())
		}
		return f(args[0].String(ctx))
	}
}

var std_trim = liftString(strings.TrimSpace, "std.trim")
var std_rstripChars = liftString2(strings.TrimRight, "std.rstripChars")
var std_lstripChars = liftString2(strings.TrimLeft, "std.lstripChars")

var std_md5 = liftString(func(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}, "std.md5")
var std_sha1 = liftString(func(s string) string {
	hash := sha1.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}, "std.sha1")
var std_sha256 = liftString(func(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}, "std.sha256")
var std_sha512 = liftString(func(s string) string {
	hash := sha512.Sum512([]byte(s))
	return hex.EncodeToString(hash[:])
}, "std.sha512")
var std_sha3 = liftString(func(s string) string {
	hash := sha3.Sum512([]byte(s))
	return hex.EncodeToString(hash[:])
}, "std.sha3")
var std_base64 = liftString(func(s string) string {
	hash := base64.StdEncoding.EncodeToString([]byte(s))
	return hash
}, "std.sha3")

var std_parseInt = liftStringToValueErr(func(s string) (evaluator.Value, error) {
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return evaluator.Value{}, fmt.Errorf("failed to parse float val (%s), err: %w", s, err)
	}
	return evaluator.MakeNumber(num), nil
}, "std.parseInt")

func std_isEmpty(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {

	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.isEmpty: %d, expected 1", len(args))
	}

	arg := args[0]
	if !arg.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type %s, expected string (std.isEmpty arg 0)", arg.Type().String())
	}

	res := arg.String(ctx) == ""

	return evaluator.MakeBool(res), nil
}

func std_codepoint(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {

	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.codepoint: %d, expected 1", len(args))
	}

	arg := args[0]
	if !arg.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type %s, expected string (std.codepoint arg 0)", arg.Type().String())
	}

	str := arg.String(ctx)
	if len(str) != 1 {
		return evaluator.Value{}, fmt.Errorf("codepoint takes a string of length 1, got length %d", len(str))
	}

	return evaluator.MakeNumber(float64(str[0])), nil
}

func std_char(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {

	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.char: %d, expected 1", len(args))
	}

	arg := args[0]
	if !arg.IsNumber() {
		return evaluator.Value{}, fmt.Errorf("unexpected type %s, expected number (std.char arg 0)", arg.Type().String())
	}
	num := arg.Number()

	if num > codepointMax {
		return evaluator.Value{}, fmt.Errorf("invalid unicode codepoint, got %v", num)
	} else if num < 0 {
		return evaluator.Value{}, fmt.Errorf("codepoints must be >= 0, got %v", num)
	}

	return evaluator.MakeString(string(rune(num)), ctx), nil
}

func std_stringChars(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {

	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.stringChars: %d, expected 1", len(args))
	}

	arg := args[0]
	if !arg.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type %s, expected string (std.stringChars arg 0)", arg.Type().String())
	}

	res := make([]evaluator.Value, 0, len(arg.String(ctx)))
	for v := range strings.SplitSeq(arg.String(ctx), "") {
		res = append(res, evaluator.MakeString(v, ctx))
	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_startsWith(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {

	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.startsWith: %d, expected 2", len(args))
	}

	full := args[0]
	if !full.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.startsWith (arg 0): %s, expected string", full.Type().String())
	}

	prefix := args[1]
	if !prefix.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.startsWith (arg 1): %s, expected string", prefix.Type().String())
	}

	res := strings.HasPrefix(full.String(ctx), prefix.String(ctx))

	return evaluator.MakeBool(res), nil
}

func std_endsWith(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {

	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.endsWith: %d, expected 2", len(args))
	}

	full := args[0]
	if !full.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.endsWith (arg 0): %s, expected string", full.Type().String())
	}

	prefix := args[1]
	if !prefix.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.endsWith (arg 1): %s, expected string", prefix.Type().String())
	}

	res := strings.HasSuffix(full.String(ctx), prefix.String(ctx))

	return evaluator.MakeBool(res), nil
}

func std_substr(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {

	if len(args) != 3 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.substr: %d, expected 2", len(args))
	}

	fullVal := args[0]
	if !fullVal.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.substr (arg 0): %s, expected string", fullVal.Type().String())
	}
	full := fullVal.String(ctx)

	fromVal := args[1]
	if !fromVal.IsNumber() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.substr (arg 1): %s, expected number", fromVal.Type().String())
	}
	from := int(fromVal.Number())
	if from < 0 {
		return evaluator.Value{}, fmt.Errorf("std.substr (arg 1) must be greater than zero, got %d", from)
	}

	toVal := args[2]
	if !toVal.IsNumber() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.substr (arg 2): %s, expected number", toVal.Type().String())
	}
	to := int(toVal.Number())
	if to < 0 {
		return evaluator.Value{}, fmt.Errorf("std.substr (arg 2) must be greater than zero, got %d", to)
	}

	if from > len(full) {
		return evaluator.MakeString("", ctx), nil
	}

	if to > len(full)-1 {
		to = len(full)
	}

	res := full[from:to]

	return evaluator.MakeString(res, ctx), nil
}

func std_findSubstr(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {

	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.findSubstr: %d, expected 2", len(args))
	}

	substrVal := args[0]
	if !substrVal.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.findSubstr (arg 0): %s, expected string", substrVal.Type().String())
	}
	substr := substrVal.String(ctx)

	fullVal := args[1]
	if !fullVal.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.findSubstr (arg 1): %s, expected string", fullVal.Type().String())
	}
	full := fullVal.String(ctx)

	res := []evaluator.Value{}

	offset := 0
	for {
		i := strings.Index(full[offset:], substr)
		if i == -1 {
			break
		}
		res = append(res, evaluator.MakeNumber(float64(offset+i)))
		offset += i + len(substr)
	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_split(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {

	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.split: %d, expected 2", len(args))
	}

	fullVal := args[0]
	if !fullVal.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.split (arg 0): %s, expected string", fullVal.Type().String())
	}
	full := fullVal.String(ctx)

	splitVal := args[1]
	if !splitVal.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.split (arg 1): %s, expected string", splitVal.Type().String())
	}
	split := splitVal.String(ctx)

	res := []evaluator.Value{}
	for _, v := range strings.Split(full, split) {
		res = append(res, evaluator.MakeString(v, ctx))
	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_splitLimit(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {

	if len(args) != 3 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.splitLimit: %d, expected 3", len(args))
	}

	fullVal := args[0]
	if !fullVal.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.splitLimit (arg 0): %s, expected string", fullVal.Type().String())
	}
	full := fullVal.String(ctx)

	splitVal := args[1]
	if !splitVal.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.splitLimit (arg 1): %s, expected string", splitVal.Type().String())
	}
	split := splitVal.String(ctx)

	maxCountVal := args[2]
	if !maxCountVal.IsNumber() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.splitLimit (arg 2): %s, expected number", maxCountVal.Type().String())
	}
	maxCount := maxCountVal.Number()

	res := []evaluator.Value{}
	for _, v := range strings.SplitN(full, split, int(maxCount)) {
		res = append(res, evaluator.MakeString(v, ctx))
	}

	return evaluator.MakeArray(res, ctx), nil
}
