# jgosonnet

A high-performance evaluator for [Jsonnet](https://jsonnet.org/).
This implementation is built to be a faster version of `go-jsonnet`, designed to efficiently handle exceptionally large files or highly complex, deeply nested configurations. See [Benchmarks](docs/benchmarks.md).

## Architecture

### Lexing and Parsing
For the frontend, `jgosonnet` utilizes the upstream `github.com/google/go-jsonnet` implementation to handle lexing and parsing. The source code is parsed by `go-jsonnet` into an Abstract Syntax Tree (AST). Reusing the upstream parser ensures strict syntactic compatibility, while allowing this project to focus purely on the execution engine and memory model for speed.

### Values and NaN-Boxing
Values uses inverted NaN-boxing to reduce memory footprint. Every value in the evaluator (numbers, nulls, booleans, strings, objects, arrays, functions, thunks) is packed into a single 64-bit unsigned integer (`uint64`).

- **Numbers** are stored natively as standard IEEE 754 float64 bit-patterns.
- **Non-numbers** leverage the vast unused bit-space of NaN (Not-a-Number) values. A specific NaN pattern is used as a tag in the upper 17 bits to denote the primitive type (e.g., String, Object, Array). The lower 47 bits store a direct virtual memory pointer (`uintptr`) to the actual data inside the Arena.

### String Interning
Object keys and string values are passed through a central String Interner during evaluation. By translating strings into 32-bit reference IDs, the evaluator avoids repetitive string allocations. String equality checks, which happen frequently during object resolution and sorting, are reduced to simple integer comparisons.

### Arena Allocation
This project uses a custom unified byte-arena for allocating evaluation state. Memory is allocated in contiguous 256KB chunks. Internal representations of Objects, Arrays, Thunks, Functions, and dynamically sized strings are aligned and written directly into this flat memory space.

This layout improves speed by replacing individual heap allocations with simple pointer arithmetic and maximizing CPU cache locality. It reduces garbage collection pressure because the GC only tracks the large byte slices rather than individual objects. Additionally, since internal references are stored as integer `uintptr` values, the GC does not traverse the evaluation graph, bypassing the overhead of GC scanning.

### Lazy Evaluation (Thunks)
Adherence to Jsonnet's lazy evaluation semantics is achieved using Thunks. A Thunk bundles an AST node with its captured lexical scope and an identifier. Computations are deferred and only evaluated when explicitly requested, such as during final JSON/YAML manifestation or when strict typing is required by standard library functions.

### Object Resolution and Layers
Jsonnet's object model supports complex inheritance (`+`, `super`, `self`) and field visibility constraints (`::`, `:::`, `:`). The architecture handles this by flattening inherited objects into `Layers` and compiling a `FieldPlan`. Instead of deeply copying or continuously merging objects during AST evaluation, the evaluator constructs a plan that points to the correct field layer and computes the final visibility mask. This flattens the inheritance tree and minimizes the computational overhead of object composition.
