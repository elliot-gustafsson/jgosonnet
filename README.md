# jgosonnet

# OBS Not feature complete yet

`jgosonnet` is a custom Jsonnet evaluator written in Go. It evaluates Jsonnet files and renders them into JSON, YAML, or Go data structures.
For parsing, `jgosonnet` uses the upstream `github.com/google/go-jsonnet` parser to generate an Abstract Syntax Tree (AST). It then evaluates this AST using a custom Go evaluation engine.

## Architecture Overview
The evaluation process has three main phases:
1. **Parsing**: Jsonnet source code is parsed into an AST using the `go-jsonnet` parser.
2. **Evaluation**: The AST is traversed by the custom runtime (`internal/evaluator`). This handles scopes, variables, lazy evaluation, and object inheritance to produce a final result.
3. **Manifestation**: The evaluated result is serialized into formats like JSON or YAML (via `EvaluateJson`, `EvaluateYaml`, `EvaluateYamlMulti`).

## Core Components
The evaluator uses specific designs to manage memory and performance:

### Arena Allocator (`Arena`)
`jgosonnet` uses an Arena-based memory management system to group related data and limit heap allocations.
- Complex types (`Objects`, `Arrays`, `Funcs`, `Thunks`, and `Scopes`) are stored in continuous slices within the Arena.
- The `Value` struct uses `uint32` IDs to reference these types in the Arena instead of Go pointers, which reduces garbage collection overhead.

### String Interner (`Interner`)
Strings are frequently used for object keys and variables in Jsonnet. `jgosonnet` maps unique strings to `uint32` IDs.
- String comparisons are performed using these integer IDs instead of comparing the string values.
- This prevents allocating duplicate strings in memory.

### Unified `Value` Representation
All data types in the runtime use a single `Value` struct that is passed by value:
- Small data types (like `float64` for numbers and booleans) are stored directly inside the struct.
- Complex types store an integer ID pointing to the `Interner` or `Arena`.
- The type is tracked using a `uint8` tag (`ValueType`).

### Lazy Evaluation (`Thunk`)
Jsonnet variables and fields are evaluated lazily. `jgosonnet` implements this using `Thunks`.
- When an expression is defined, a `Thunk` is created instead of evaluating it immediately.
- The `Thunk` stores the AST node, the scope ID, and the context (such as `self` and `super`).
- When the value is required (for example, during output generation), the `Thunk` calculates the result and caches it.


### Object Inheritance
Jsonnet objects support inheritance, mixins (`+`), and field visibility modifiers (`::`, `:::`).
- `jgosonnet` uses a layered `Object` struct to track merged keys, overridden nodes, and field visibility.
