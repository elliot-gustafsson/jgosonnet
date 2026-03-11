package stdlib

import "fmt"

func StdSlice[T any](arr []T, start, end, step int) ([]T, error) {

	if step <= 0 {
		return nil, fmt.Errorf("got %d but step must be greater than 0", step)
	}

	end = min(end, len(arr))
	capacity := (end - start + step - 1) / step

	res := make([]T, 0, capacity)
	for i := start; i < end; i += step {
		res = append(res, arr[i])
	}
	return res, nil
}
