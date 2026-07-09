package stdlib

import (
	"fmt"
	"slices"
	"strings"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

func std_range(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	from, err := args[0].EvalInteger(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	to, err := args[1].EvalInteger(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if to < from {
		_, arrVal := evaluator.MakeArraySized(0, ctx)
		return arrVal, nil
	}

	length := int(to-from) + 1
	res, arrVal := evaluator.MakeArraySized(length, ctx)
	for i := range length {
		res[i] = evaluator.MakeNumber(float64(from + i))
	}

	return arrVal, nil
}

func std_makeArray(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	size, err := args[0].EvalInteger(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	mapperFunc, err := args[1].EvalFunction(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	allArgs := make([]evaluator.NamedValue, size)

	res, arrVal := evaluator.MakeArraySized(size, ctx)
	for i := range size {

		allArgs[i] = evaluator.NamedValue{Value: evaluator.MakeNumber(float64(i))}

		n, _ := ctx.Registry.GoCallbackNodes.New()
		n.Func = mapperFunc
		n.Args = allArgs[i : i+1]

		res[i] = evaluator.MakeThunk(evaluator.NewThunk(n, 0, ctx), ctx)

	}

	return arrVal, nil
}

func std_join(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	sep, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	inputArray, err := args[1].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	inputLen := len(inputArray)

	if sep.IsString() {

		totalLen := 0
		validCount := 0
		sepLen := len(sep.String(ctx))

		for i := range inputLen {
			v, err := inputArray[i].Eval(ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			if v.IsNull() {
				continue
			}
			if !v.IsString() {
				return evaluator.Value{}, evaluator.TypeErrorSpecific(evaluator.ValueTypeString, v.Type())
			}

			totalLen += len(v.String(ctx))
			validCount++
		}

		if validCount > 1 {
			totalLen += sepLen * (validCount - 1)
		}

		var sb strings.Builder
		sb.Grow(totalLen)

		hasWritten := false
		for _, v := range inputArray {

			v, err := v.Eval(ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

			if v.IsNull() {
				continue
			}

			if hasWritten {
				sb.WriteString(sep.String(ctx))
			}
			sb.WriteString(v.String(ctx))
			hasWritten = true
		}

		return evaluator.MakeString(sb.String(), ctx), nil
	}

	if !sep.IsArray() {
		return evaluator.Value{}, evaluator.TypeErrorSpecific(evaluator.ValueTypeArray, sep.Type())
	}

	sepArray := sep.Array(ctx)
	sepLen := len(sepArray)

	totalCap := 0
	validCount := 0
	for i := range inputLen {

		v, err := inputArray[i].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if v.IsNull() {
			continue
		}

		if !v.IsArray() {
			return evaluator.Value{}, evaluator.TypeErrorSpecific(evaluator.ValueTypeArray, v.Type())
		}
		totalCap += len(v.Array(ctx))
		validCount++
	}

	if validCount > 1 {
		totalCap += sepLen * (validCount - 1)
	}

	hasWritten := false
	res, arrVal := evaluator.MakeArraySized(totalCap, ctx)
	idx := 0
	for _, v := range inputArray {

		v, err := v.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if v.IsNull() {
			continue
		}

		if hasWritten {
			idx += copy(res[idx:], sepArray)
		}
		idx += copy(res[idx:], v.Array(ctx))
		hasWritten = true
	}

	return arrVal, nil
}

func std_deepJoin(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	arrVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if arrVal.IsString() {
		return evaluator.MakeString(arrVal.String(ctx), ctx), nil
	}

	if !arrVal.IsArray() {
		return evaluator.Value{}, evaluator.TypeErrorSpecific(evaluator.ValueTypeArray, arrVal.Type())
	}
	arr := arrVal.Array(ctx)

	var b strings.Builder

	err = flattenArray(&b, arr, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeString(b.String(), ctx), nil
}

func flattenArray(b *strings.Builder, arr []evaluator.Value, ctx evaluator.Context) error {

	for _, v := range arr {
		v, err := v.Eval(ctx)
		if err != nil {
			return err
		}
		if v.IsArray() {
			err := flattenArray(b, v.Array(ctx), ctx)
			if err != nil {
				return err
			}
			continue
		}

		if v.IsString() {
			b.WriteString(v.String(ctx))
			continue
		}

		return fmt.Errorf("Expected string or array, got %s", v.Type().String())
	}

	return nil
}

func std_filter(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	mapperFunc, err := args[0].EvalFunction(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	inputArray, err := args[1].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	mapperFuncInput := []evaluator.NamedValue{{}}

	res := []evaluator.Value{}
	for _, v := range inputArray {
		mapperFuncInput[0] = evaluator.NamedValue{Value: v}
		out, err := mapperFunc.Exec(mapperFuncInput, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		b, err := out.EvalBool(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if b {
			res = append(res, v)
		}

	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_flatMap(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	mapFunc, err := args[0].EvalFunction(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	arr, err := args[1].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	mapFuncArgs := []evaluator.NamedValue{{}}

	res := make([]evaluator.Value, 0, len(arr))
	for _, v := range arr {

		mapFuncArgs[0] = evaluator.NamedValue{Value: v}

		out, err := mapFunc.Exec(mapFuncArgs, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		outArr, err := out.EvalArray(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		res = append(res, outArr...)
	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_filterMap(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	filterFunc, err := args[0].EvalFunction(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	mapFunc, err := args[1].EvalFunction(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	inputArray, err := args[2].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	funcArgs := []evaluator.NamedValue{{}}

	filteredArr := make([]evaluator.Value, 0, len(inputArray)/2)
	for _, v := range inputArray {
		funcArgs[0] = evaluator.NamedValue{Value: v}
		out, err := filterFunc.Exec(funcArgs, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		b, err := out.EvalBool(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if b {
			filteredArr = append(filteredArr, v)
		}

	}

	res := make([]evaluator.Value, 0, len(filteredArr))
	for _, v := range filteredArr {
		funcArgs[0] = evaluator.NamedValue{Value: v}
		out, err := mapFunc.Exec(funcArgs, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		res = append(res, out)
	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_uniq(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	arr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	var keyF evaluator.Function
	if !args[1].IsNone() {
		f, err := args[1].EvalFunction(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		keyF = f
	}

	// Create the array once and mutate it to reduce objects on the heap
	mapperFuncInput := []evaluator.NamedValue{{}}

	var last evaluator.Value
	res := make([]evaluator.Value, 0, len(arr))
	for _, v := range arr {

		if last.IsNone() {
			res = append(res, v)
			last = v
			continue
		}

		var x evaluator.Value
		var y evaluator.Value
		var err error

		if !keyF.Noop() {

			mapperFuncInput[0] = evaluator.NamedValue{Value: last}
			x, err = keyF.Exec(mapperFuncInput, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

			mapperFuncInput[0] = evaluator.NamedValue{Value: v}
			y, err = keyF.Exec(mapperFuncInput, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
		} else {
			x = last
			y = v
		}

		x, err = x.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		y, err = y.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if x.Type() == y.Type() {

			eq, err := x.Equal(y, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			if eq {
				continue
			}

		}
		res = append(res, v)
		last = v
		continue
	}

	return evaluator.MakeArray(res, ctx), nil

}

func std_sort(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	arr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	var keyF evaluator.Function
	if !args[1].IsNone() {
		f, err := args[1].EvalFunction(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		keyF = f
	}

	res, err := sortArray(arr, keyF, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeArray(res, ctx), nil

}

func sortArray(arr []evaluator.Value, keyF evaluator.Function, ctx evaluator.Context) (res []evaluator.Value, err error) {
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

	mapperFuncInput := []evaluator.NamedValue{{}}

	// TODO: now we eval the values over and over, think abt this
	slices.SortStableFunc(result, func(a, b evaluator.Value) int {

		if !keyF.Noop() {
			mapperFuncInput[0] = evaluator.NamedValue{Value: a}
			a, err = keyF.Exec(mapperFuncInput, ctx)
			if err != nil {
				panic(err)
			}
			mapperFuncInput[0] = evaluator.NamedValue{Value: b}
			b, err = keyF.Exec(mapperFuncInput, ctx)
			if err != nil {
				panic(err)
			}
		}

		a, err := a.Eval(ctx)
		if err != nil {
			panic(err)
		}

		b, err = b.Eval(ctx)
		if err != nil {
			panic(err)
		}

		if a.IsObject() || b.IsObject() {
			err := fmt.Errorf("unexpected type object")
			panic(err)
		}

		if a.Type() != b.Type() {
			err := fmt.Errorf("unexpected type %s, expected %s", a.Type().String(), b.Type().String())
			panic(err)
		}

		x, err := a.Compare(b, ctx)
		if err != nil {
			panic(err)
		}
		return x

	})

	return result, nil
}

// Shortcut for std.uniq(std.sort(arr)).
func std_set(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	sorted, err := std_sort(args, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	uniqArgs := []evaluator.NamedValue{{Value: sorted}, args[1]}

	set, err := std_uniq(uniqArgs, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	return set, nil
}

func std_map(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	mapFunc, err := args[0].EvalFunction(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	arr, err := args[1].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	allArgs := make([]evaluator.NamedValue, len(arr))

	res, arrVal := evaluator.MakeArraySized(len(arr), ctx)
	for i, v := range arr {

		allArgs[i] = evaluator.NamedValue{Value: v}

		n, _ := ctx.Registry.GoCallbackNodes.New()
		n.Func = mapFunc
		n.Args = allArgs[i : i+1]

		res[i] = evaluator.MakeThunk(evaluator.NewThunk(n, 0, ctx), ctx)

	}

	return arrVal, nil
}

func std_mapWithIndex(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	mapFunc, err := args[0].EvalFunction(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	arr, err := args[1].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	// Create the array once and mutate it to reduce objects on the heap
	// mapperFuncInput := []evaluator.NamedValue{{}, {}}

	// res := make([]evaluator.Value, 0, len(arr))
	// for i, v := range arr {
	// 	mapperFuncInput[0] = evaluator.NamedValue{Value: evaluator.MakeNumber(float64(i))}
	// 	mapperFuncInput[1] = evaluator.NamedValue{Value: v}
	// 	out, err := mapFunc.Exec(mapperFuncInput, ctx)
	// 	if err != nil {
	// 		return evaluator.Value{}, err
	// 	}
	// 	res = append(res, out)
	// }

	allArgs := make([]evaluator.NamedValue, len(arr)*2)

	res, arrVal := evaluator.MakeArraySized(len(arr), ctx)
	for i, v := range arr {
		idx := i * 2

		allArgs[idx] = evaluator.NamedValue{Value: evaluator.MakeNumber(float64(i))}
		allArgs[idx+1] = evaluator.NamedValue{Value: v}

		n, _ := ctx.Registry.GoCallbackNodes.New()
		n.Func = mapFunc
		n.Args = allArgs[idx : idx+2]

		res[i] = evaluator.MakeThunk(evaluator.NewThunk(n, 0, ctx), ctx)

	}

	return arrVal, nil
}

func std_member(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	indexable, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	arg, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if indexable.IsString() {

		if !arg.IsString() {
			return evaluator.Value{}, evaluator.TypeErrorSpecific(evaluator.ValueTypeString, arg.Type())
		}
		idx := strings.Index(indexable.String(ctx), arg.String(ctx))
		return evaluator.MakeBool(idx >= 0), nil
	}

	if !indexable.IsArray() {
		return evaluator.Value{}, evaluator.TypeErrorSpecific(evaluator.ValueTypeArray, arg.Type())
	}

	inputArr := indexable.Array(ctx)

	for _, v := range inputArr {

		v, err := v.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if v.Type() != arg.Type() {
			continue
		}

		eq, err := v.Equal(arg, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if eq {
			return evaluator.MakeBool(true), nil
		}
	}

	return evaluator.MakeBool(false), nil
}

func std_setMember(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 3 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.setMember %d, expected 3", len(args))
	}

	member, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	arr, err := args[1].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	var keyF evaluator.Function
	if !args[2].IsNone() {
		f, err := args[2].EvalFunction(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		keyF = f
	}

	mapperFuncInput := []evaluator.NamedValue{{}}

	for _, v := range arr {
		v, err := v.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		ar := v
		br := member

		if !keyF.Noop() {
			mapperFuncInput[0] = evaluator.NamedValue{Value: v}
			ar, err = keyF.Exec(mapperFuncInput, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

			mapperFuncInput[0] = evaluator.NamedValue{Value: member}
			br, err = keyF.Exec(mapperFuncInput, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
		}

		if v.Type() != member.Type() {
			continue
		}

		eq, err := ar.Equal(br, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if eq {
			return evaluator.MakeBool(true), nil
		}

	}

	return evaluator.MakeBool(false), nil
}

func std_count(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	indexable, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	arg, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	count := 0
	for _, v := range indexable {
		v, err := v.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if v.Type() != arg.Type() {
			continue
		}

		res, err := v.Equal(arg, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if res {
			count++
		}
	}

	return evaluator.MakeNumber(float64(count)), nil
}

func std_slice(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	indexable, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	index, err := args[1].EvalNumber(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	indexInt := int(index)

	end, err := args[2].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !end.IsNull() && !end.IsNumber() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.slice (arg 2): %s, expected number", end.Type().String())
	}

	step, err := args[3].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
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

func std_lines(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.lines %d, expected 1", len(args))
	}

	indexable, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !indexable.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.lines (arg 0): %s, expected array", indexable.Type().String())
	}

	b := strings.Builder{}
	for _, v := range indexable.Array(ctx) {
		v, err := v.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if v.IsNull() {
			continue
		}

		if !v.IsString() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.lines array: %s, expected strings", v.Type().String())
		}

		b.WriteString(v.String(ctx))
		b.WriteByte('\n')
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

func std_reverse(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.reverse %d, expected 1", len(args))
	}

	indexable, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !indexable.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.reverse (arg 0): %s, expected array", indexable.Type().String())
	}

	arr := indexable.Array(ctx)
	res, arrVal := evaluator.MakeArraySized(len(arr), ctx)
	for i, v := range arr {
		res[len(arr)-1-i] = v
	}

	return arrVal, nil
}

func std_foldl(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	// std.foldl(func, arr, init)
	if len(args) != 3 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.foldl %d, expected 3", len(args))
	}

	fVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !fVal.IsFunction() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.foldl (arg 0): %s, expected function", fVal.Type().String())
	}
	foldFunc := fVal.Function(ctx)
	foldFuncArgs := []evaluator.NamedValue{{}, {}}

	state, err := args[2].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	arrVal, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if arrVal.IsString() {
		for _, v := range arrVal.String(ctx) {
			str := evaluator.MakeString(string(v), ctx)

			foldFuncArgs[0] = evaluator.NamedValue{Value: state}
			foldFuncArgs[1] = evaluator.NamedValue{Value: str}

			val, err := foldFunc.Exec(foldFuncArgs, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

			state = val
		}

		return state, nil
	}

	if !arrVal.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.foldl (arg 1): %s, expected string or array", arrVal.Type().String())
	}

	for _, v := range arrVal.Array(ctx) {

		foldFuncArgs[0] = evaluator.NamedValue{Value: state}
		foldFuncArgs[1] = evaluator.NamedValue{Value: v}

		val, err := foldFunc.Exec(foldFuncArgs, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		state = val
	}

	return state, nil
}

func std_foldr(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	// std.foldr(func, arr, init)
	if len(args) != 3 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.foldr %d, expected 3", len(args))
	}

	fVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !fVal.IsFunction() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.foldr (arg 0): %s, expected function", fVal.Type().String())
	}
	foldFunc := fVal.Function(ctx)
	foldFuncArgs := []evaluator.NamedValue{{}, {}}

	state, err := args[2].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	arrVal, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if arrVal.IsString() {
		runes := []rune(arrVal.String(ctx))

		for i := len(runes) - 1; i >= 0; i-- {
			v := runes[i]

			str := evaluator.MakeString(string(v), ctx)

			foldFuncArgs[0] = evaluator.NamedValue{Value: str}
			foldFuncArgs[1] = evaluator.NamedValue{Value: state}

			val, err := foldFunc.Exec(foldFuncArgs, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

			state = val
		}

		return state, nil
	}

	if !arrVal.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.foldr (arg 1): %s, expected string or array", arrVal.Type().String())
	}

	arr := arrVal.Array(ctx)
	for i := len(arr) - 1; i >= 0; i-- {
		v := arr[i]

		foldFuncArgs[0] = evaluator.NamedValue{Value: v}
		foldFuncArgs[1] = evaluator.NamedValue{Value: state}

		val, err := foldFunc.Exec(foldFuncArgs, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		state = val
	}

	return state, nil
}

func std_sum(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.sum %d, expected 1", len(args))
	}

	arr, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if !arr.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.sum (arg 0): %s, expected array", arr.Type().String())
	}

	var sum float64
	for _, v := range arr.Array(ctx) {
		v, err := v.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !v.IsNumber() {
			return evaluator.Value{}, fmt.Errorf("unexpected type in std.sum arr: %s, expected number", arr.Type().String())
		}
		sum += v.Number()
	}
	return evaluator.MakeNumber(sum), nil
}

func std_flattenArrays(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.flattenArrays %d, expected 1", len(args))
	}

	arr, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if !arr.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.flattenArrays (arg 0): %s, expected array", arr.Type().String())
	}

	res := []evaluator.Value{}
	for _, v := range arr.Array(ctx) {
		v, err := v.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !v.IsArray() {
			return evaluator.Value{}, fmt.Errorf("unexpected type in std.flattenArrays arr: %s, expected array", v.Type().String())
		}
		res = append(res, v.Array(ctx)...)
	}
	return evaluator.MakeArray(res, ctx), nil
}

func std_flattenDeepArray(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.flattenDeepArray %d, expected 1", len(args))
	}

	arrVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if !arrVal.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.flattenDeepArray (arg 0): %s, expected array", arrVal.Type().String())
	}
	arr := arrVal.Array(ctx)

	res := make([]evaluator.Value, 0, len(arr))
	for _, v := range arr {
		v, err := v.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if v.IsArray() {
			res, err = flattedDeepArray(res, v.Array(ctx), ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			continue
		}
		res = append(res, v)
	}
	return evaluator.MakeArray(res, ctx), nil
}

func flattedDeepArray(state, v []evaluator.Value, ctx evaluator.Context) ([]evaluator.Value, error) {
	for _, x := range v {
		x, err := x.Eval(ctx)
		if err != nil {
			return nil, err
		}
		if x.IsArray() {
			state, err = flattedDeepArray(state, x.Array(ctx), ctx)
			if err != nil {
				return nil, err
			}
			continue
		}
		state = append(state, x)
	}
	return state, nil
}

func std_repeat(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.repeat %d, expected 2", len(args))
	}

	whatVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	countVal, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !countVal.IsNumber() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.repeat (arg 1): %s, expected number", countVal.Type().String())
	}
	count := int(countVal.Number())

	if whatVal.IsString() {
		what := whatVal.String(ctx)
		sb := strings.Builder{}
		sb.Grow(count * len(what))
		for range count {
			sb.WriteString(what)
		}
		return evaluator.MakeString(sb.String(), ctx), nil
	}

	if whatVal.IsArray() {
		what := whatVal.Array(ctx)
		res, arrVal := evaluator.MakeArraySized(count*len(what), ctx)
		for i := 0; i < count; i++ {
			copy(res[i*len(what):], what)
		}
		return arrVal, nil
	}

	return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.repeat (arg 0): %s, expected array or string", whatVal.Type().String())
}

func std_setUnion(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 3 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.setUnion %d, expected 3", len(args))
	}

	aVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !aVal.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.setUnion (arg 0): %s, expected array", aVal.Type().String())
	}

	bVal, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !bVal.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.setUnion (arg 1): %s, expected array", bVal.Type().String())
	}

	var keyF evaluator.Function
	var funcArgs []evaluator.NamedValue
	if !args[2].IsNone() {
		f, err := args[2].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !f.IsFunction() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.setUnion (arg 1): %s, expected function", f.Type().String())
		}
		keyF = f.Function(ctx)
		funcArgs = []evaluator.NamedValue{{}}
	}

	aArr := aVal.Array(ctx)
	bArr := bVal.Array(ctx)

	i, j := 0, 0

	res := make([]evaluator.Value, 0, (len(aArr)+len(bArr))/2)
	for i < len(aArr) && j < len(bArr) {

		aV, err := aArr[i].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		bV, err := bArr[j].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		var x int

		if !keyF.Noop() {

			funcArgs[0] = evaluator.NamedValue{Value: aV}
			aVC, err := keyF.Exec(funcArgs, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			aVC, err = aVC.Eval(ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

			funcArgs[0] = evaluator.NamedValue{Value: bV}
			bVC, err := keyF.Exec(funcArgs, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			bVC, err = bVC.Eval(ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

			x, err = aVC.Compare(bVC, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

		} else {
			x, err = aV.Compare(bV, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
		}

		if x < 0 {
			res = append(res, aV)
			i++
		} else if x > 0 {
			res = append(res, bV)
			j++
		} else { // x == 0
			res = append(res, aV)
			i++
			j++
		}
	}

	if i < len(aArr) {
		for i < len(aArr) {
			aV, err := aArr[i].Eval(ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			// TODO: dont append duplicates
			res = append(res, aV)
			i++
		}
	}

	if j < len(bArr) {
		for j < len(bArr) {
			bV, err := bArr[j].Eval(ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			// TODO: dont append duplicates
			res = append(res, bV)
			j++
		}
	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_setInter(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 3 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.setInter %d, expected 3", len(args))
	}

	aVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !aVal.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.setInter (arg 0): %s, expected array", aVal.Type().String())
	}

	bVal, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !bVal.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.setInter (arg 1): %s, expected array", bVal.Type().String())
	}

	var keyF evaluator.Function
	var funcArgs []evaluator.NamedValue
	if !args[2].IsNone() {
		f, err := args[2].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !f.IsFunction() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.setInter (arg 1): %s, expected function", f.Type().String())
		}
		keyF = f.Function(ctx)
		funcArgs = []evaluator.NamedValue{{}}
	}

	aArr := aVal.Array(ctx)
	bArr := bVal.Array(ctx)

	i, j := 0, 0

	res := make([]evaluator.Value, 0, (len(aArr)+len(bArr))/2)
	for i < len(aArr) && j < len(bArr) {

		aV, err := aArr[i].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		bV, err := bArr[j].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		var x int

		if !keyF.Noop() {

			funcArgs[0] = evaluator.NamedValue{Value: aV}
			aVC, err := keyF.Exec(funcArgs, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			aVC, err = aVC.Eval(ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

			funcArgs[0] = evaluator.NamedValue{Value: bV}
			bVC, err := keyF.Exec(funcArgs, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			bVC, err = bVC.Eval(ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

			x, err = aVC.Compare(bVC, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

		} else {
			x, err = aV.Compare(bV, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
		}

		if x < 0 {
			i++
		} else if x > 0 {
			j++
		} else { // x == 0
			res = append(res, aV)
			i++
			j++
		}
	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_setDiff(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 3 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.setDiff %d, expected 3", len(args))
	}

	aArr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	bArr, err := args[1].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	var keyF evaluator.Function
	var funcArgs []evaluator.NamedValue
	if !args[2].IsNone() {
		f, err := args[2].EvalFunction(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		keyF = f
		funcArgs = []evaluator.NamedValue{{}}
	}

	i, j := 0, 0
	res := make([]evaluator.Value, 0, len(aArr))
	for i < len(aArr) && j < len(bArr) {
		aV, err := aArr[i].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		bV, err := bArr[j].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		var x int
		if !keyF.Noop() {
			funcArgs[0] = evaluator.NamedValue{Value: aV}
			aVC, err := keyF.Exec(funcArgs, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

			aVC, err = aVC.Eval(ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

			funcArgs[0] = evaluator.NamedValue{Value: bV}
			bVC, err := keyF.Exec(funcArgs, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			bVC, err = bVC.Eval(ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

			x, err = aVC.Compare(bVC, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
		} else {
			x, err = aV.Compare(bV, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
		}

		if x < 0 {
			res = append(res, aV)
			i++
		} else if x > 0 {
			j++
		} else { // x == 0
			i++
			j++
		}
	}

	if i < len(aArr) {
		for i < len(aArr) {
			aV, err := aArr[i].Eval(ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			res = append(res, aV)
			i++
		}
	}
	return evaluator.MakeArray(res, ctx), nil
}

func std_find(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	valueVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	arr, err := args[1].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	res := []evaluator.Value{}
	for i, v := range arr {
		v, err := v.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if valueVal.Type() != v.Type() {
			continue
		}

		eq, err := valueVal.Equal(v, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if eq {
			res = append(res, evaluator.MakeNumber(float64(i)))
		}
	}

	return evaluator.MakeArray(res, ctx), nil
}

func std_any(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	arr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	for _, v := range arr {
		b, err := v.EvalBool(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if b {
			return evaluator.MakeBool(true), nil
		}
	}

	return evaluator.MakeBool(false), nil
}

func std_all(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.all %d, expected 1", len(args))
	}

	arr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	for _, v := range arr {
		b, err := v.EvalBool(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if !b {
			return evaluator.MakeBool(false), nil
		}
	}

	return evaluator.MakeBool(true), nil
}

func std_avg(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	arr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	var total float64
	for _, v := range arr {
		n, err := v.EvalNumber(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		total += n
	}

	avg := total / float64(len(arr))
	return evaluator.MakeNumber(avg), nil
}

func std_minArray(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	return minMaxArray(args, ctx, false, "std.minArray")
}

func std_maxArray(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	return minMaxArray(args, ctx, true, "std.maxArray")
}

func minMaxArray(args []evaluator.NamedValue, ctx evaluator.Context, max bool, name string) (evaluator.Value, error) {

	arr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	var keyF evaluator.Function
	var funcArgs []evaluator.NamedValue
	if !args[1].IsNone() {
		f, err := args[1].EvalFunction(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		keyF = f
		funcArgs = []evaluator.NamedValue{{}}
	}

	var defaultValue evaluator.Value
	if !args[2].IsNone() {
		d, err := args[2].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		defaultValue = d
	}

	if len(arr) == 0 {
		if defaultValue.IsNone() {
			return evaluator.Value{}, fmt.Errorf("Expected at least one element in array. Got none")
		}
		return defaultValue, nil
	}

	var last evaluator.Value
	for _, v := range arr {
		v, err := v.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if last.IsNone() {
			last = v
			continue
		}

		x, y := last, v
		if !keyF.Noop() {
			funcArgs[0] = evaluator.NamedValue{Value: last}
			x, err = keyF.Exec(funcArgs, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

			funcArgs[0] = evaluator.NamedValue{Value: v}
			y, err = keyF.Exec(funcArgs, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
		}

		if x.Type() != y.Type() {
			return evaluator.Value{}, fmt.Errorf("Unexpected type %s, expected %s", y.Type().String(), x.Type().String())
		}

		cmp, err := x.Compare(y, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if max {
			if cmp < 0 {
				last = v
				continue
			}
		} else {
			if cmp > 0 {
				last = v
				continue
			}
		}

	}

	return last, nil
}

func std_contains(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.contains %d, expected 2", len(args))
	}

	arr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	elem, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	for _, v := range arr {
		v, err := v.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if elem.Type() != v.Type() {
			continue
		}

		eq, err := elem.Equal(v, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if eq {
			return evaluator.MakeBool(true), nil
		}
	}

	return evaluator.MakeBool(false), nil
}

func std_remove(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.remove %d, expected 2", len(args))
	}

	arr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	elem, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	i := 0
	res := make([]evaluator.Value, 0, len(arr))
	for _, v := range arr {
		i++

		v, err := v.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if elem.Type() != v.Type() {
			continue
		}

		eq, err := elem.Equal(v, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if eq {
			// Break on first match
			break
		}
		res = append(res, v)
	}

	// Add rest
	res = append(res, arr[i:]...)

	return evaluator.MakeArray(res, ctx), nil
}

func std_removeAt(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.removeAt %d, expected 2", len(args))
	}

	arr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	idx, err := args[1].EvalInteger(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if idx > len(arr)-1 {
		return evaluator.Value{}, fmt.Errorf("idx %d is out of bound for array with length %d", idx, len(arr))
	}

	res, arrVal := evaluator.MakeArraySized(len(arr)-1, ctx)
	copy(res, arr[:idx])
	copy(res[idx:], arr[idx+1:])

	return arrVal, nil
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
