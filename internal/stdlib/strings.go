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
	"slices"
	"strings"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

const codepointMax = 0x10FFFF

func liftString(f func(string) string) evaluator.Func {
	return func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
		a, err := args[0].EvalString(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		res := f(a)
		return evaluator.MakeString(res, ctx), nil
	}
}

func liftStringErr(f func(string) (string, error)) evaluator.Func {
	return func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
		a, err := args[0].EvalString(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		res, err := f(a)
		if err != nil {
			return evaluator.Value{}, err
		}
		return evaluator.MakeString(res, ctx), nil
	}
}

func liftString2(f func(string, string) string) evaluator.Func {
	return func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
		a, err := args[0].EvalString(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		b, err := args[1].EvalString(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		res := f(a, b)
		return evaluator.MakeString(res, ctx), nil
	}
}

var std_trim = liftString(strings.TrimSpace)
var std_stripChars = liftString2(strings.Trim)
var std_rstripChars = liftString2(strings.TrimRight)
var std_lstripChars = liftString2(strings.TrimLeft)

var std_md5 = liftString(func(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
})
var std_sha1 = liftString(func(s string) string {
	hash := sha1.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
})
var std_sha256 = liftString(func(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
})
var std_sha512 = liftString(func(s string) string {
	hash := sha512.Sum512([]byte(s))
	return hex.EncodeToString(hash[:])
})
var std_sha3 = liftString(func(s string) string {
	hash := sha3.Sum512([]byte(s))
	return hex.EncodeToString(hash[:])
})

func std_base64(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	inputVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	var toEncode []byte
	if inputVal.IsString() {
		toEncode = []byte(inputVal.String(ctx))
	} else if inputVal.IsArray() {
		arr := inputVal.Array(ctx)
		toEncode = make([]byte, 0, len(arr))
		for _, v := range arr {
			v, err := v.Eval(ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

			numInt := int(v.Number())
			if !v.IsNumber() || float64(numInt) != v.Number() {
				err := fmt.Errorf("base64 encountered a non-integer value in the array, got %s", v.Type())
				return evaluator.Value{}, evaluator.MakeRuntimeError(err)
			}

			if numInt < 0 || 255 < numInt {
				err := fmt.Errorf("base64 encountered invalid codepoint value in the array (must be 0 <= X <= 255), got %d", numInt)
				return evaluator.Value{}, evaluator.MakeRuntimeError(err)
			}
			toEncode = append(toEncode, byte(numInt))
		}
	} else {
		err := fmt.Errorf("base64 can only base64 encode strings / arrays of single bytes, got %s", inputVal.Type())
		return evaluator.Value{}, evaluator.MakeRuntimeError(err)
	}

	res := base64.StdEncoding.EncodeToString(toEncode)
	return evaluator.MakeString(res, ctx), nil
}

var std_base64Decode = liftStringErr(func(s string) (string, error) {
	out, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(out), err
})

func std_base64DecodeBytes(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	arg, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	out, err := base64.StdEncoding.DecodeString(arg)
	if err != nil {
		return evaluator.Value{}, err
	}
	res := make([]evaluator.Value, 0, len(out))
	for _, v := range out {
		res = append(res, evaluator.MakeNumber(float64(v)))
	}
	return evaluator.MakeArray(res, ctx), nil
}

var std_asciiLower = liftString(strings.ToLower)
var std_asciiUpper = liftString(strings.ToUpper)

func std_isEmpty(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	arg, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeBool(len(arg) == 0), nil
}

func std_codepoint(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	arg, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if len(arg) != 1 {
		return evaluator.Value{}, evaluator.MakeRuntimeError(fmt.Errorf("codepoint takes a string of length 1, got length %d", len(arg)))
	}

	return evaluator.MakeNumber(float64(arg[0])), nil
}

func std_char(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	num, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if num > codepointMax {
		return evaluator.Value{}, evaluator.MakeRuntimeError(fmt.Errorf("invalid unicode codepoint, got %v", num))
	} else if num < 0 {
		return evaluator.Value{}, evaluator.MakeRuntimeError(fmt.Errorf("codepoints must be >= 0, got %v", num))
	}

	return evaluator.MakeString(string(rune(num)), ctx), nil
}

func std_stringChars(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	arg, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	res := make([]evaluator.Value, 0, len(arg))
	for v := range strings.SplitSeq(arg, "") {
		res = append(res, evaluator.MakeString(v, ctx))
	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_startsWith(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	full, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	prefix, err := args[1].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	res := strings.HasPrefix(full, prefix)

	return evaluator.MakeBool(res), nil
}

func std_endsWith(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	full, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	prefix, err := args[1].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	res := strings.HasSuffix(full, prefix)

	return evaluator.MakeBool(res), nil
}

func std_substr(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	fullVal, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	full := []rune(fullVal)

	fromFloat, err := args[1].EvalNumber(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	from := int(fromFloat)
	if float64(from) != fromFloat {
		return evaluator.Value{}, evaluator.MakeRuntimeError(fmt.Errorf("substr second parameter should be an integer, got %f", fromFloat))
	}

	if from < 0 {
		return evaluator.Value{}, evaluator.MakeRuntimeError(fmt.Errorf("substr second parameter should be greater than zero, got %d", from))
	}

	lenFloat, err := args[2].EvalNumber(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	length := int(lenFloat)
	if float64(length) != lenFloat {
		return evaluator.Value{}, evaluator.MakeRuntimeError(fmt.Errorf("substr third parameter should be an integer, got %f", lenFloat))
	}

	if length < 0 {
		return evaluator.Value{}, evaluator.MakeRuntimeError(fmt.Errorf("substr third parameter should be greater than zero, got %d", length))
	}

	if from > len(full) {
		return evaluator.MakeString("", ctx), nil
	}

	to := from + length

	if to > len(full)-1 {
		res := full[from:]
		return evaluator.MakeString(string(res), ctx), nil
	}

	res := full[from:to]
	return evaluator.MakeString(string(res), ctx), nil
}

func std_findSubstr(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	substr, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	full, err := args[1].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	res := []evaluator.Value{}

	if substr == "" {
		return evaluator.MakeArray(res, ctx), nil
	}

	offset := 0
	for {
		i := strings.Index(full[offset:], substr)
		if i == -1 {
			break
		}
		res = append(res, evaluator.MakeNumber(float64(offset+i)))
		offset += i + 1
	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_strReplace(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	full, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	from, err := args[1].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	to, err := args[2].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	res := strings.ReplaceAll(full, from, to)

	return evaluator.MakeString(res, ctx), nil
}

func std_split(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	full, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	split, err := args[1].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	res := []evaluator.Value{}
	for _, v := range strings.Split(full, split) {
		res = append(res, evaluator.MakeString(v, ctx))
	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_splitLimit(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	full, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	split, err := args[1].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	maxSplits, err := args[2].EvalInteger(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	var arr []string
	if maxSplits < 0 {
		arr = strings.Split(full, split)
	} else {
		arr = strings.SplitN(full, split, maxSplits+1)
	}

	res := []evaluator.Value{}
	for _, v := range arr {
		res = append(res, evaluator.MakeString(v, ctx))
	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_splitLimitR(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	full, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	split, err := args[1].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	maxSplits, err := args[2].EvalInteger(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if maxSplits < 0 {
		arr := strings.Split(full, split)
		res := make([]evaluator.Value, 0, len(arr))
		for _, v := range arr {
			res = append(res, evaluator.MakeString(v, ctx))
		}
		return evaluator.MakeArray(res, ctx), nil
	}

	s := full

	res := make([]evaluator.Value, 0, maxSplits)
	for range maxSplits {
		idx := strings.LastIndex(s, split)
		if idx < 0 {
			break // No more separators found
		}
		x := s[idx+len(split):]
		res = append(res, evaluator.MakeString(x, ctx))
		// Truncate the string
		s = s[:idx]
	}

	res = append(res, evaluator.MakeString(s, ctx))

	slices.Reverse(res)

	return evaluator.MakeArray(res, ctx), nil
}

var std_escapeStringBash = liftString(func(s string) string {
	var b strings.Builder
	b.WriteByte('\'')

	last := 0
	for i := 0; i < len(s); i++ {
		if s[i] != '\'' {
			continue
		}
		b.WriteString(s[last:i])
		b.WriteString(`'"'"'`)
		last = i + 1
	}
	b.WriteString(s[last:])

	b.WriteByte('\'')
	return b.String()

})

var std_escapeStringDollars = liftString(func(s string) string {
	return strings.ReplaceAll(s, "$", "$$")
})

var std_escapeStringXML = liftString(func(s string) string {
	var b strings.Builder
	last := 0
	for i := 0; i < len(s); i++ {
		var r string
		switch s[i] {
		default:
			continue
		case '<':
			r = "&lt;"
		case '>':
			r = "&gt;"
		case '&':
			r = "&amp;"
		case '"':
			r = "&quot;"
		case '\'':
			r = "&apos;"
		}
		b.WriteString(s[last:i])
		b.WriteString(r)
		last = i + 1
	}
	b.WriteString(s[last:])
	return b.String()
})

func std_escapeStringJson(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	strVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !strVal.IsString() {
		return evaluator.Value{}, evaluator.TypeErrorSpecific(evaluator.ValueTypeString, strVal.Type())
	}

	var b strings.Builder

	err = evaluator.ManifestJson(&b, strVal, ctx, evaluator.JsonConfigMinified)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

func std_equalsIgnoreCase(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	b, err := args[1].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeBool(strings.EqualFold(a, b)), nil
}
