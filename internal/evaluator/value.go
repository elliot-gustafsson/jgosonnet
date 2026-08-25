package evaluator

import (
	"cmp"
	"fmt"
	"math"
	"strings"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
)

type ValueType uint8

const (
	ValueTypeNone ValueType = iota
	ValueTypeNull
	ValueTypeString
	ValueTypeNumber
	ValueTypeBool
	ValueTypeObject
	ValueTypeArray
	ValueTypeFunction
	ValueTypeNativeFunction
	ValueTypeThunk
)

func (t ValueType) IsLiteral() bool {
	return t == ValueTypeString || t == ValueTypeNumber || t == ValueTypeBool || t == ValueTypeNull
}

func (t ValueType) String() string {
	switch t {
	case ValueTypeNone:
		return "none"
	case ValueTypeNull:
		return "null"
	case ValueTypeString:
		return "string"
	case ValueTypeNumber:
		return "number"
	case ValueTypeBool:
		return "boolean"
	case ValueTypeObject:
		return "object"
	case ValueTypeArray:
		return "array"
	case ValueTypeFunction, ValueTypeNativeFunction:
		return "function"
	case ValueTypeThunk:
		return "thunk"
	default:
		return fmt.Sprintf("unknown (%d)", t)
	}
}

// Value represents a NaN-boxed 64-bit value.
// Bit layout for non-float types:
//
// 63                   51 50    47 46                               0
// |----------------------|--------|---------------------------------|
// |  Float NaN Boundary  |  Type  |             Payload             |
// |       (13 bits)      |(4 bits)|            (47 bits)            |
// |----------------------|--------|---------------------------------|
//
// - Bits 63-51: Reserved for IEEE-754 Quiet NaN float64 detection.
// - Bits 50-47: ValueType enum.
// - Bits 46-0 : Payload (pointer info, interned string id, 1/0 for bools)
//
//	Note: lowest bits can be used for flags depending on the alignment
type Value uint64

const (
	nanTag           uint64 = 0xFFF8000000000000
	tagMask                 = 0xFFFF800000000000
	floatThreshold   uint64 = 1 << 51
	canonicalNaNBits uint64 = 0x7FF8000000000000
	typeShift               = 47

	ValueNone Value = 0

	// ValueFlagStringConst uint64 = 1 << 40
)

func box(t ValueType, data uint32) Value {
	// pack the type into bits 50-47, and the data into bits 0-31
	return Value((uint64(t) << typeShift) | uint64(data))
}

func boxPtr(t ValueType, p unsafe.Pointer) Value {
	return Value(uint64(t)<<typeShift | uint64(uintptr(p)))
}

func (v Value) unboxPtr() unsafe.Pointer {
	addr := uintptr(v &^ tagMask)
	return *(*unsafe.Pointer)(unsafe.Pointer(&addr))
}

type NamedValue struct {
	Key uint32

	Value
}

type CachedValue struct {
	Value

	Visible bool
}

const (
	intSize    = unsafe.Sizeof(int(0))
	intAlign   = unsafe.Alignof(int(0))
	valueSize  = unsafe.Sizeof(Value(0))
	valueAlign = unsafe.Alignof(Value(0))
)

func MakeNull() (rv Value) {
	rv = box(ValueTypeNull, 0)
	return
}

func MakeString(s string, ctx Context) (rv Value) {
	n := len(s)

	totalSize := intSize + uintptr(n)

	ptr := arena.AlignedAlloc(ctx.State.Allocator, totalSize, intAlign)

	// write len at beginning of block
	*(*int)(ptr) = n

	// write data
	bytePtr := unsafe.Add(ptr, intSize)
	copy(unsafe.Slice((*byte)(bytePtr), n), s)

	rv = boxPtr(ValueTypeString, ptr)
	return
}

func MakeStringFromBytes(in []byte, ctx Context) (rv Value) {
	n := len(in)

	totalSize := intSize + uintptr(n)

	ptr := arena.AlignedAlloc(ctx.State.Allocator, totalSize, intAlign)

	// write len at beginning of block
	*(*int)(ptr) = n

	// write data
	bytePtr := unsafe.Add(ptr, intSize)
	copy(unsafe.Slice((*byte)(bytePtr), n), in)

	rv = boxPtr(ValueTypeString, ptr)
	return
}

func MakeStringConcat(v1, v2 string, ctx Context) (rv Value) {
	n := len(v1) + len(v2)

	totalSize := intSize + uintptr(n)

	ptr := arena.AlignedAlloc(ctx.State.Allocator, totalSize, intAlign)

	// write len at beginning of block
	*(*int)(ptr) = n

	// write data
	bytePtr := unsafe.Add(ptr, intSize)

	buf := unsafe.Slice((*byte)(bytePtr), n)
	x := copy(buf, v1)
	copy(buf[x:], v2)

	rv = boxPtr(ValueTypeString, ptr)
	return
}

func AllocStringBuilder(ctx Context, capacity int) (unsafe.Pointer, []byte) {
	totalSize := intSize + uintptr(capacity)
	ptr := arena.AlignedAlloc(ctx.State.Allocator, totalSize, intAlign)

	// Create a 0-length slice pointing to the space right after the length header
	buf := unsafe.Slice((*byte)(unsafe.Add(ptr, intSize)), capacity)[:0]
	return ptr, buf
}

// FinalizeStringBuilder writes the length header and boxes the pointer into a Value.
func FinalizeStringBuilder(ptr unsafe.Pointer, length int) Value {
	*(*int)(ptr) = length
	return boxPtr(ValueTypeString, ptr)
}

func MakeStringConst(id uint32) Value {
	// set lowest bit to 1 to signal its a interned string
	return Value((uint64(ValueTypeString) << typeShift) | (uint64(id) << 1) | 1)

}

func MakeNumber(v float64) (rv Value) {
	bits := math.Float64bits(v)
	if math.IsNaN(v) {
		bits = canonicalNaNBits
	}
	rv = Value(bits ^ nanTag)
	return
}

func MakeBool(v bool) (rv Value) {
	if v {
		rv = box(ValueTypeBool, 1)
		return
	}
	rv = box(ValueTypeBool, 0)
	return
}

func MakeObjectValue(v *Object) (rv Value) {
	rv = boxPtr(ValueTypeObject, unsafe.Pointer(v))
	return
}

func MakeArraySized(n int, ctx Context) (arr []Value, rv Value) {
	if n < 0 {
		panic("MakeArraySized: called with negative size")
	}

	totalSize := intSize + (uintptr(n) * valueSize)
	ptr := arena.AlignedAlloc(ctx.State.Allocator, totalSize, valueAlign)

	*(*int)(ptr) = n

	if n > 0 {
		elemPtr := (*Value)(unsafe.Add(ptr, intSize))
		arr = unsafe.Slice(elemPtr, n)
	} else {
		arr = unsafe.Slice((*Value)(nil), 0)
	}

	rv = boxPtr(ValueTypeArray, ptr)
	return
}

func MakeArray(v []Value, ctx Context) (rv Value) {
	n := len(v)
	totalSize := intSize + (uintptr(n) * valueSize)

	ptr := arena.AlignedAlloc(ctx.State.Allocator, totalSize, valueAlign)

	*(*int)(ptr) = n

	if n > 0 {
		elemPtr := (*Value)(unsafe.Add(ptr, intSize))
		copy(unsafe.Slice(elemPtr, n), v)
	}

	return boxPtr(ValueTypeArray, ptr)
}

func MakeFunctionValue(v *Function, ctx Context) (rv Value) {
	rv = boxPtr(ValueTypeFunction, unsafe.Pointer(v))
	return
}

func MakeNativeFunctionValue(v *NativeFunction, ctx Context) (rv Value) {
	rv = boxPtr(ValueTypeNativeFunction, unsafe.Pointer(v))
	return
}

func MakeThunkValue(v *Thunk, ctx Context) (rv Value) {
	rv = boxPtr(ValueTypeThunk, unsafe.Pointer(v))
	return
}

func MakeTombstoneValue(scope int) (rv Value) {
	rv = box(ValueTypeNone, uint32(scope))
	return
}

func (v Value) Type() ValueType {
	t := uint64(v) >> typeShift
	if t >= (floatThreshold >> typeShift) {
		return ValueTypeNumber
	}
	return ValueType(t)
}

func (v Value) RefId() uint32 {
	return uint32(v)
}

func (v Value) Payload() uintptr {
	return uintptr(v &^ tagMask)
}

func (v Value) String(ctx Context) string {
	if (uint64(v) & 1) == 1 {
		return ctx.State.Interner.Get(uint32(v.Payload() >> 1))
	}

	ptr := v.unboxPtr()
	n := *(*int)(ptr)
	return unsafe.String((*byte)(unsafe.Add(ptr, intSize)), n)
}

func (v Value) Number() float64 {
	return math.Float64frombits(uint64(v) ^ nanTag)
}

func (v Value) Bool() bool {
	// Because MakeBool stored 1 or 0 in the RefId slot
	return uint32(v) == 1
}

func (v Value) Array() []Value {
	ptr := v.unboxPtr()
	n := *(*int)(ptr)
	if n == 0 {
		return nil
	}
	return unsafe.Slice((*Value)(unsafe.Add(ptr, intSize)), n)
}

func (v Value) Object() *Object {
	return (*Object)(v.unboxPtr())
}

func (v Value) Function() *Function {
	return (*Function)(v.unboxPtr())
}

func (v Value) NativeFunction() *NativeFunction {
	return (*NativeFunction)(v.unboxPtr())
}

func (v Value) Thunk() *Thunk {
	return (*Thunk)(v.unboxPtr())
}

func (v Value) Eval(ctx Context) (rv Value, err error) {
	if !v.IsThunk() {
		rv = v
		return
	}
	rv, err = v.evalThunk(ctx)
	return
}

//go:noinline
func (v Value) evalThunk(ctx Context) (Value, error) {
	thunk := v.Thunk()

	if thunk.EvalState == -1 {
		return thunk.Value, nil
	}

	return thunk.Eval(ctx)
}

type RuntimeError struct {
	err error
}

func (e RuntimeError) Error() string {
	return "RUNTIME ERROR: " + e.err.Error()
}

func MakeRuntimeError(err error) RuntimeError {
	return RuntimeError{err}
}

func TypeErrorSpecific(good, bad ValueType) error {
	return MakeRuntimeError(fmt.Errorf("Unexpected type %s, expected %s", bad, good))
}

func TypeErrorGeneral(bad ValueType) error {
	return MakeRuntimeError(fmt.Errorf("Unexpected type %s", bad))
}

func (v Value) EvalString(ctx Context) (string, error) {
	x, err := v.Eval(ctx)
	if err != nil {
		return "", err
	}
	if !x.IsString() {
		return "", TypeErrorSpecific(ValueTypeString, x.Type())
	}
	return x.String(ctx), nil
}

func (v Value) EvalNumber(ctx Context) (float64, error) {
	x, err := v.Eval(ctx)
	if err != nil {
		return 0, err
	}
	if !x.IsNumber() {
		return 0, TypeErrorSpecific(ValueTypeNumber, x.Type())
	}
	return x.Number(), nil
}

func (v Value) EvalInteger(ctx Context) (int, error) {
	x, err := v.Eval(ctx)
	if err != nil {
		return 0, err
	}
	if !x.IsNumber() {
		return 0, TypeErrorSpecific(ValueTypeNumber, x.Type())
	}
	intNum := int(x.Number())
	if float64(intNum) != x.Number() {
		return 0, MakeRuntimeError(fmt.Errorf("Expected an integer, but got %v", x.Number()))
	}
	return intNum, nil
}

func (v Value) EvalBool(ctx Context) (bool, error) {
	x, err := v.Eval(ctx)
	if err != nil {
		return false, err
	}
	if !x.IsBool() {
		return false, TypeErrorSpecific(ValueTypeBool, x.Type())
	}
	return x.Bool(), nil
}

func (v Value) EvalArray(ctx Context) ([]Value, error) {
	x, err := v.Eval(ctx)
	if err != nil {
		return nil, err
	}
	if !x.IsArray() {
		return nil, TypeErrorSpecific(ValueTypeArray, x.Type())
	}
	return x.Array(), nil
}

func (v Value) EvalObject(ctx Context) (*Object, error) {
	x, err := v.Eval(ctx)
	if err != nil {
		return nil, err
	}
	if !x.IsObject() {
		return nil, TypeErrorSpecific(ValueTypeObject, x.Type())
	}
	return x.Object(), nil
}

// func (v Value) EvalFunction(ctx Context) (*Function, error) {
// 	x, err := v.Eval(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !x.IsFunction() {
// 		return nil, TypeErrorSpecific(ValueTypeFunction, x.Type())
// 	}
// 	return x.Function(), nil
// }

func (v Value) ToString(ctx Context) (string, error) {

	switch v.Type() {
	default:
		return "", fmt.Errorf("unhandled type %s, string conversion not available", v.Type().String())
	case ValueTypeNull:
		return "null", nil
	case ValueTypeString:
		return v.String(ctx), nil
	case ValueTypeNumber:
		var p [64]byte
		res := unparseNumber(p[:0], v.Number())
		return string(res), nil
	case ValueTypeBool:
		if v.Bool() {
			return "true", nil
		}
		return "false", nil
	case ValueTypeArray:
		var b strings.Builder
		err := ManifestJson(&b, v, ctx, JsonConfigToString)
		if err != nil {
			return "", err
		}
		return b.String(), nil
	case ValueTypeObject:
		subCtx := ctx
		subCtx.Self = v
		var b strings.Builder
		err := ManifestJson(&b, v, subCtx, JsonConfigToString)
		if err != nil {
			return "", err
		}
		return b.String(), nil
	case ValueTypeThunk:
		v, err := v.Eval(ctx)
		if err != nil {
			return "", err
		}
		return v.ToString(ctx)
	}

}

func (v Value) IsNone() bool {
	return v == 0
}

func (v Value) IsLiteral() bool {
	return v.Type().IsLiteral()
}

func (v Value) IsNull() bool {
	return uint64(v)>>typeShift == uint64(ValueTypeNull)
}

func (v Value) IsString() bool {
	return uint64(v)>>typeShift == uint64(ValueTypeString)
}

func (v Value) IsNumber() bool {
	return uint64(v) >= floatThreshold
}

func (v Value) IsBool() bool {
	return uint64(v)>>typeShift == uint64(ValueTypeBool)
}

func (v Value) IsThunk() bool {
	return uint64(v)>>typeShift == uint64(ValueTypeThunk)
}

func (v Value) IsObject() bool {
	return uint64(v)>>typeShift == uint64(ValueTypeObject)
}

func (v Value) IsFunction() bool {
	t := uint64(v) >> typeShift
	return t == uint64(ValueTypeFunction) || t == uint64(ValueTypeNativeFunction)
}

func (v Value) IsArray() bool {
	return uint64(v)>>typeShift == uint64(ValueTypeArray)
}

func (v Value) IsEmpty(ctx Context) bool {
	switch v.Type() {
	default:
		return false
	case ValueTypeNull:
		return true
	case ValueTypeObject:
		return v.Object().Length(ctx) == 0
	case ValueTypeArray:
		return len(v.Array()) == 0
	}
}

func (v Value) AsStringConst() uint32 {
	if uint64(v)>>typeShift != uint64(ValueTypeString) || (uint64(v)&1) != 1 {
		return 0
	}
	return uint32(v.Payload() >> 1)
}

func (v Value) FunctionExec(args []NamedValue, ctx Context) (Value, error) {
	return v.FunctionExecEx(args, ctx, false)
}

func (v Value) FunctionExecEx(args []NamedValue, ctx Context, tailstrict bool) (Value, error) {
	v, err := v.Eval(ctx)
	if err != nil {
		return ValueNone, err
	}
	return execFunction(v, args, ctx, tailstrict)
}

func (v Value) Prune(ctx Context) (Value, error) {
	switch v.Type() {
	default:
		return ValueNone, fmt.Errorf("unhandled type (%s) in Value.Prune()", v.Type())
	case ValueTypeNull, ValueTypeString, ValueTypeNumber, ValueTypeBool:
		return v, nil
	case ValueTypeObject:
		return v.Object().Prune(ctx)
	case ValueTypeArray:
		arr := v.Array()
		res := make([]Value, 0, len(arr))
		for _, v := range arr {
			v, err := v.Eval(ctx)
			if err != nil {
				return ValueNone, err
			}
			out, err := v.Prune(ctx)
			if err != nil {
				return ValueNone, err
			}
			if out.IsEmpty(ctx) {
				continue
			}
			res = append(res, out)
		}
		return MakeArray(res, ctx), nil
	}

}

func (a Value) Equal(b Value, ctx Context) (bool, error) {
	if a.Type() != b.Type() {
		return false, nil
	}

	switch a.Type() {
	case ValueTypeNull:
		return true, nil
	case ValueTypeString:
		if a == b {
			return true, nil
		}
		return a.String(ctx) == b.String(ctx), nil
	case ValueTypeNumber:
		return a.Number() == b.Number(), nil
	case ValueTypeBool:
		return a.Bool() == b.Bool(), nil
	case ValueTypeObject:
		return a.Object().Equal(b.Object(), ctx)
	case ValueTypeArray:
		aArr := a.Array()
		bArr := b.Array()

		if len(aArr) != len(bArr) {
			return false, nil
		}

		for i, av := range aArr {
			av, err := av.Eval(ctx)
			if err != nil {
				return false, err
			}
			bv, err := bArr[i].Eval(ctx)
			if err != nil {
				return false, err
			}

			eq, err := av.Equal(bv, ctx)
			if err != nil {
				return false, err
			}

			if !eq {
				return false, nil
			}

		}

		return true, nil
	case ValueTypeThunk:
		a, err := a.Eval(ctx)
		if err != nil {
			return false, err
		}
		b, err := b.Eval(ctx)
		if err != nil {
			return false, err
		}
		return a.Equal(b, ctx)
	default:
		return false, fmt.Errorf("comparing types %s is not supported", a.Type().String())
	}

}

func (a Value) Compare(b Value, ctx Context) (int, error) {
	if a.Type() != b.Type() {
		return 0, fmt.Errorf("comparing %s with %s is not supported", a.Type().String(), b.Type().String())
	}

	switch a.Type() {
	case ValueTypeNull:
		return 0, nil
	case ValueTypeString:
		if a == b {
			return 0, nil
		}
		return cmp.Compare(a.String(ctx), b.String(ctx)), nil
	case ValueTypeNumber:
		return cmp.Compare(a.Number(), b.Number()), nil
	// case ValueTypeBool:
	// 	return 0, fmt.Errorf("comparing boolean values is not supported")
	// case ValueTypeObject:
	// 	return a.Object().Compare(b.Object(), ctx)
	case ValueTypeArray:
		aArr := a.Array()
		bArr := b.Array()

		i, j := 0, 0
		for i < len(aArr) && j < len(bArr) {

			iv, err := aArr[i].Eval(ctx)
			if err != nil {
				return 0, err
			}

			jv, err := bArr[j].Eval(ctx)
			if err != nil {
				return 0, err
			}

			x, err := iv.Compare(jv, ctx)
			if err != nil {
				return 0, err
			}

			if x != 0 {
				return x, nil
			}
			i++
			j++
		}
		return cmp.Compare(len(aArr), len(bArr)), nil
	case ValueTypeThunk:
		a, err := a.Eval(ctx)
		if err != nil {
			return 0, err
		}
		b, err := b.Eval(ctx)
		if err != nil {
			return 0, err
		}
		return a.Compare(b, ctx)
	default:
		return 0, fmt.Errorf("comparing types %s is not supported", a.Type().String())
	}

}
