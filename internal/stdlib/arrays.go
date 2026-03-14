package stdlib

import (
	"fmt"
	"slices"
	"strings"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

func std_range(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.range %d, expected 2", len(args))
	}

	from := args[0]
	if !from.IsNumber() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.range (arg 0): %s, expected number", from.Type().String())
	}

	to := args[1]
	if !to.IsNumber() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.range (arg 1): %s, expected number", to.Type().String())
	}

	fromIndex := int(from.Number())
	toIndex := int(to.Number())

	if toIndex < fromIndex {
		return evaluator.MakeArray(make([]evaluator.Value, 0), ctx), nil
	}

	res := make([]evaluator.Value, 0, toIndex-fromIndex)
	for i := fromIndex; i <= toIndex; i++ {
		res = append(res, evaluator.MakeNumber(float64(i)))
	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_makeArray(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.makeArray %d, expected 2", len(args))
	}

	val := args[0]
	if !val.IsNumber() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.makeArray (arg 0): %s, expected number", val.Type().String())
	}

	f := args[1]
	if !f.IsFunction() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.makeArray (arg 1): %s, expected function", f.Type().String())
	}

	// Create the array once and mutate it to reduce object on the heap
	mapperFuncInput := []evaluator.Value{{}}

	res := make([]evaluator.Value, int(val.Number()))
	for i := range int(val.Number()) {
		v := evaluator.MakeNumber(float64(i))
		mapperFuncInput[0] = v
		out, err := f.Function(ctx)(mapperFuncInput, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		res[i] = out
	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_join(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.join %d, expected 2", len(args))
	}

	sep := args[0]
	arr := args[1]

	// TODO: Think, is pre-calculate worth it

	inputArray := arr.Array(ctx) // Avoid calling .Array() repeatedly
	inputLen := len(inputArray)

	if sep.IsString() {

		if !arr.IsArray() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.join (arg 1): %s, expected array if arg 0 is string", arr.Type().String())
		}

		totalLen := 0
		sepLen := len(sep.String(ctx))

		for i := range inputLen {
			err := evaluator.EvaluateValueStrict(&inputArray[i], ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			if !inputArray[i].IsString() {
				return evaluator.Value{}, fmt.Errorf("second parameter to std.join must be an array of strings if first argument is a string, got %s", inputArray[i].Type().String())
			}

			totalLen += len(inputArray[i].String(ctx))
		}

		if inputLen > 1 {
			totalLen += sepLen * (inputLen - 1)
		}

		var sb strings.Builder
		sb.Grow(totalLen)

		for i, v := range inputArray {

			if i > 0 {
				// Dont write separator first iteration
				_, err := sb.WriteString(sep.String(ctx))
				if err != nil {
					return evaluator.Value{}, fmt.Errorf("failed to write separator, err: %w", err)
				}
			}
			_, err := sb.WriteString(v.String(ctx))
			if err != nil {
				return evaluator.Value{}, fmt.Errorf("failed to write value, err: %w", err)
			}
		}

		return evaluator.MakeString(sb.String(), ctx), nil
	}

	if !sep.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.join (arg 0): %s, expected number or array", sep.Type().String())
	}

	sepArray := sep.Array(ctx)
	sepLen := len(sepArray)

	totalCap := 0
	for i := range inputLen {

		err := evaluator.EvaluateValueStrict(&inputArray[i], ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !inputArray[i].IsArray() {
			return evaluator.Value{}, fmt.Errorf("second parameter to std.join must be an array of arrays if first argument is an array")
		}
		totalCap += len(inputArray[i].Array(ctx))
	}

	if inputLen > 1 {
		totalCap += sepLen * (inputLen - 1)
	}

	res := make([]evaluator.Value, 0, totalCap)
	for i, v := range inputArray {

		if i > 0 {
			res = append(res, sepArray...)
		}
		res = append(res, v.Array(ctx)...)
	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_filter(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.filter %d, expected 2", len(args))
	}

	f := args[0]
	if !f.IsFunction() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.filter (arg 0): %s, expected function", f.Type().String())
	}

	arr := args[1]
	if !arr.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.filter (arg 1): %s, expected array", arr.Type().String())
	}

	inputArray := arr.Array(ctx)

	// Create the array once and mutate it to reduce object on the heap
	mapperFuncInput := []evaluator.Value{{}}

	res := []evaluator.Value{}
	for _, v := range inputArray {
		mapperFuncInput[0] = v
		out, err := f.Function(ctx)(mapperFuncInput, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if !out.IsBool() {
			return evaluator.Value{}, fmt.Errorf("unexpected responser from std.filter func: %s, expected bool", arr.Type().String())
		}

		if out.Bool() {
			res = append(res, v)
		}

	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_filterMap(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 3 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.filterMap %d, expected 2", len(args))
	}

	filterFunc := args[0]
	if !filterFunc.IsFunction() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.filterMap (arg 0): %s, expected function", filterFunc.Type().String())
	}

	mapFunc := args[1]
	if !mapFunc.IsFunction() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.filterMap (arg 1): %s, expected function", mapFunc.Type().String())
	}

	arr := args[2]
	if !arr.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.filterMap (arg 2): %s, expected array", arr.Type().String())
	}

	inputArray := arr.Array(ctx)

	// Create the array once and mutate it to reduce object on the heap
	mapperFuncInput := []evaluator.Value{{}}

	filteredArr := make([]evaluator.Value, 0, len(inputArray)/2)
	for _, v := range inputArray {
		mapperFuncInput[0] = v
		out, err := filterFunc.Function(ctx)(mapperFuncInput, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if !out.IsBool() {
			return evaluator.Value{}, fmt.Errorf("unexpected responser from std.filter func: %s, expected bool", arr.Type().String())
		}

		if out.Bool() {
			filteredArr = append(filteredArr, v)
		}

	}

	res := make([]evaluator.Value, 0, len(filteredArr))
	for _, v := range filteredArr {
		mapperFuncInput[0] = v
		out, err := mapFunc.Function(ctx)(mapperFuncInput, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		res = append(res, out)
	}

	return evaluator.MakeArray(res, ctx), nil
}

var id_func = func(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed %d, expected 1", len(args))
	}
	return args[0], nil
}

func std_uniq(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) == 0 || len(args) > 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.uniq %d, expected 1 or 2", len(args))
	}

	arr := args[0]
	if !arr.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.uniq (arg 0: %s, expected array", arr.Type().String())
	}

	var keyF evaluator.Func = id_func
	if len(args) > 1 {
		f := args[1]
		if !f.IsFunction() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.uniq (arg 1): %s, expected function", f.Type().String())
		}
		keyF = f.Function(ctx)
	}

	inputArr := arr.Array(ctx)

	// Create the array once and mutate it to reduce object on the heap
	mapperFuncInput := []evaluator.Value{{}}

	var last evaluator.Value
	res := make([]evaluator.Value, 0, len(inputArr))
	for _, v := range inputArr {

		if last.IsNone() {
			res = append(res, v)
			last = v
			continue
		}

		mapperFuncInput[0] = last
		x, err := keyF(mapperFuncInput, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		mapperFuncInput[0] = v
		y, err := keyF(mapperFuncInput, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		err = evaluator.EvaluateValueStrict(&x, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		err = evaluator.EvaluateValueStrict(&y, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if x.Type() == y.Type() {

			eq, err := evaluator.BopManifestEqual(x, y, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			if eq.Bool() {
				continue
			}

		}
		res = append(res, v)
		last = v
		continue
	}

	return evaluator.MakeArray(res, ctx), nil

}

func std_sort(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.sort %d, expected 1 or 2", len(args))
	}

	arr := args[0]
	if !arr.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.sort (arg 0: %s, expected array", arr.Type().String())
	}

	var keyF evaluator.Func
	if len(args) > 1 {
		f := args[1]
		if !f.IsFunction() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.sort (arg 1): %s, expected function", f.Type().String())
		}
		keyF = f.Function(ctx)
	}

	inputArr := arr.Array(ctx)

	res, err := sortArray(inputArr, keyF, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeArray(res, ctx), nil

}

func sortArray(arr []evaluator.Value, keyF evaluator.Func, ctx evaluator.Context) (res []evaluator.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("unexpected panic during sort: %v", r)
			}
		}
	}()

	result := slices.Clone(arr)

	// TODO: now we eval the values over and over, think abt this
	slices.SortFunc(result, func(a, b evaluator.Value) int {
		err := evaluator.EvaluateValueStrict(&a, ctx)
		if err != nil {
			panic(err)
		}

		err = evaluator.EvaluateValueStrict(&b, ctx)
		if err != nil {
			panic(err)
		}

		ar := a
		br := b

		if keyF != nil {
			ar, err = keyF([]evaluator.Value{a}, ctx)
			if err != nil {
				panic(err)
			}
			br, err = keyF([]evaluator.Value{b}, ctx)
			if err != nil {
				panic(err)
			}
		}

		err = evaluator.EvaluateValueStrict(&ar, ctx)
		if err != nil {
			panic(err)
		}

		err = evaluator.EvaluateValueStrict(&br, ctx)
		if err != nil {
			panic(err)
		}

		if ar.IsObject() || br.IsObject() {
			err := fmt.Errorf("unexpected type object")
			panic(err)
		}

		if ar.Type() != br.Type() {
			err := fmt.Errorf("unexpected type %s, expected %s", a.Type().String(), b.Type().String())
			panic(err)
		}

		greater, err := evaluator.BopGreater(ar, br, ctx)
		if err != nil {
			panic(err)
		}

		if greater.Bool() {
			return 1
		}
		return -1

	})

	return result, nil
}

// Shortcut for std.uniq(std.sort(arr)).
func std_set(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) == 0 || len(args) > 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.set %d, expected 1 or 2", len(args))
	}

	sorted, err := std_sort(args, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	uniqArgs := []evaluator.Value{sorted}
	if len(args) > 1 {
		uniqArgs = append(uniqArgs, args[1])
	}

	set, err := std_uniq(uniqArgs, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	return set, nil
}

func std_map(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.map %d, expected 2", len(args))
	}

	f := args[0]
	if !f.IsFunction() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.map (arg 0): %s, expected function", f.Type().String())
	}

	arr := args[1]
	if !arr.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.map (arg 1): %s, expected array", arr.Type().String())
	}

	inputArr := arr.Array(ctx)

	// Create the array once and mutate it to reduce object on the heap
	mapperFuncInput := []evaluator.Value{{}}

	res := make([]evaluator.Value, 0, len(inputArr))
	for _, v := range inputArr {
		mapperFuncInput[0] = v
		out, err := f.Function(ctx)(mapperFuncInput, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		res = append(res, out)
	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_member(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.member %d, expected 2", len(args))
	}

	indexable := args[0]
	arg := args[1]

	if indexable.IsString() {

		err := evaluator.EvaluateValueStrict(&arg, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if !arg.IsString() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.member (arg 1): %s, expected string", arg.Type().String())
		}

		for _, s := range indexable.String(ctx) {
			v := evaluator.MakeString(string(s), ctx)
			eq, err := evaluator.BopManifestEqual(v, arg, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			if eq.Bool() {
				return evaluator.MakeBool(true), nil
			}

		}
		return evaluator.MakeBool(false), nil
	}

	if !indexable.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.member (arg 0): %s, expected array or string", indexable.Type().String())
	}

	inputArr := indexable.Array(ctx)

	err := evaluator.EvaluateValueStrict(&arg, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	for _, v := range inputArr {

		if v.Type() != arg.Type() {
			continue
		}

		eq, err := evaluator.BopManifestEqual(v, arg, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if eq.Bool() {
			return evaluator.MakeBool(true), nil
		}
	}

	return evaluator.MakeBool(false), nil
}

func std_setMember(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.setMember %d, expected 2-3", len(args))
	}

	member := args[0]

	arr := args[1]
	if !arr.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.member (arg 1): %s, expected array", arr.Type().String())
	}

	var keyF evaluator.Func
	if len(args) > 2 {
		f := args[2]
		if !f.IsFunction() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.member (arg 2): %s, expected function", f.Type().String())
		}
		keyF = f.Function(ctx)
	}

	// Create the array once and mutate it to reduce object on the heap
	mapperFuncInput := []evaluator.Value{{}}

	for _, v := range arr.Array(ctx) {
		err := evaluator.EvaluateValueStrict(&v, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		ar := v
		br := member

		if keyF != nil {
			mapperFuncInput[0] = v
			ar, err = keyF(mapperFuncInput, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

			mapperFuncInput[0] = member
			br, err = keyF(mapperFuncInput, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
		}

		if v.Type() != member.Type() {
			continue
		}

		eq, err := evaluator.BopManifestEqual(ar, br, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if eq.Bool() {
			return evaluator.MakeBool(true), nil
		}

	}

	return evaluator.MakeBool(false), nil
}

func std_count(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.count %d, expected 2", len(args))
	}

	indexable := args[0]
	if !indexable.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.count (arg 0): %s, expected array", indexable.Type().String())
	}

	arg := args[1]
	err := evaluator.EvaluateValueStrict(&arg, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	count := 0
	for _, v := range indexable.Array(ctx) {
		err := evaluator.EvaluateValueStrict(&v, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if v.Type() != arg.Type() {
			continue
		}

		res, err := evaluator.BopManifestEqual(v, arg, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if res.Bool() {
			count++
		}
	}

	return evaluator.MakeNumber(float64(count)), nil
}

func std_slice(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {

	if len(args) != 4 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.slice %d, expected 4", len(args))
	}

	indexable := args[0]

	index := args[1]
	if !index.IsNumber() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.slice (arg 1): %s, expected number", index.Type().String())
	}
	indexInt := int(index.Number())

	end := args[2]
	if !end.IsNull() && !end.IsNumber() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.slice (arg 2): %s, expected number", end.Type().String())
	}

	step := args[3]
	if !step.IsNull() && !step.IsNumber() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.slice (arg 3): %s, expected number", step.Type().String())
	}

	if indexable.IsString() {
		x := []rune(indexable.String(ctx))

		endInt := len(x)
		if end.IsNumber() {
			endInt = int(end.Number())
		}

		stepInt := 1
		if step.IsNumber() {
			stepInt = int(step.Number())
		}

		res, err := sliceArr(x, indexInt, endInt, stepInt)
		if err != nil {
			return evaluator.Value{}, err
		}
		return evaluator.MakeString(string(res), ctx), nil
	}

	if indexable.IsArray() {
		x := indexable.Array(ctx)

		endInt := len(x)
		if end.IsNumber() {
			endInt = int(end.Number())
		}

		stepInt := 1
		if step.IsNumber() {
			stepInt = int(step.Number())
		}

		res, err := sliceArr(x, indexInt, endInt, stepInt)
		if err != nil {
			return evaluator.Value{}, err
		}
		return evaluator.MakeArray(res, ctx), nil
	}

	return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.slice (arg 0): %s, expected string or array", indexable.Type().String())
}

func std_lines(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.lines %d, expected 1", len(args))
	}

	indexable := args[0]
	if !indexable.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.lines (arg 0): %s, expected array", indexable.Type().String())
	}

	b := strings.Builder{}
	for _, v := range indexable.Array(ctx) {
		err := evaluator.EvaluateValueStrict(&v, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if !v.IsString() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.lines array: %s, expected strings", v.Type().String())
		}

		b.WriteString(v.String(ctx))
		b.WriteByte('\n')
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

func std_reverse(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.reverse %d, expected 1", len(args))
	}

	indexable := args[0]
	if !indexable.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.reverse (arg 0): %s, expected array", indexable.Type().String())
	}

	res := slices.Clone(indexable.Array(ctx))

	slices.Reverse(res)

	return evaluator.MakeArray(res, ctx), nil
}

func sliceArr[T any](arr []T, start, end, step int) ([]T, error) {

	if step <= 0 {
		return nil, fmt.Errorf("got %d but step must be greater than 0", step)
	}

	arrLen := len(arr)

	if start > arrLen {
		return []T{}, nil
	}

	end = min(end, len(arr))
	if end < 0 {
		end = max(len(arr)+end, 0)
	}

	capacity := max((end-start+step-1)/step, 0)

	res := make([]T, 0, capacity)
	for i := start; i < end; i += step {
		res = append(res, arr[i])
	}
	return res, nil
}
