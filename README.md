# jgosonnet

A high-performance evaluator for [Jsonnet](https://jsonnet.org/).
This implementation is built to be a faster version of `go-jsonnet`, designed to efficiently handle exceptionally large files or highly complex, deeply nested configurations.

## Architecture

### Lexing and Parsing
For the frontend, `jgosonnet` utilizes the upstream `github.com/google/go-jsonnet` implementation to handle lexing and parsing. The source code is parsed by `go-jsonnet` into an Abstract Syntax Tree (AST). Reusing the upstream parser ensures strict syntactic compatibility, while allowing this project to focus purely on the execution engine and memory model for speed.

### Values and NaN-Boxing
Values uses inverted NaN-boxing to reduce memory footprint. Every value in the evaluator (numbers, booleans, strings, objects, arrays, thunks) is packed into a single 64-bit unsigned integer (`uint64`).

* **Numbers** are stored natively as standard IEEE 754 float64 bit-patterns.
* **Non-numbers** leverage the vast unused bit-space of NaN (Not-a-Number) values. A specific NaN pattern is used as a tag in the upper 16 bits to denote the primitive type (e.g., String, Object, Array). The lower 32 bits store an index (`refId`) pointing to the actual data inside the Arena.

### String Interning
Object keys and string values are passed through a central String Interner during evaluation. By translating strings into 32-bit reference IDs (`refId`), the evaluator avoids repetitive string allocations. String equality checks, which happen frequently during object resolution and sorting, are reduced to simple integer comparisons.

### Arena Allocation
Rather than independently allocating every JSON object, array, and function on the Go heap, the evaluator utilizes an Arena allocator. Runtime representations of Objects, Arrays, Thunks, and Functions are stored in contiguous memory slices within an evaluation `Context`. This reduces the pressure on Go's Garbage Collector (GC) and improves cache locality, which is critical when parsing massive data structures.

### Lazy Evaluation (Thunks)
Adherence to Jsonnet's lazy evaluation semantics is achieved using Thunks. A Thunk bundles an AST node with its captured lexical scope and an identifier. Computations are deferred and only evaluated when explicitly requested, such as during final JSON/YAML manifestation or when strict typing is required by standard library functions.

### Object Resolution and Layers
Jsonnet's object model supports complex inheritance (`+`, `super`, `self`) and field visibility constraints (`::`, `:::`, `:`). The architecture handles this by flattening inherited objects into `Layers` and compiling a `FieldPlan`. Instead of deeply copying or continuously merging objects during AST evaluation, the evaluator constructs a plan that points to the correct field layer and computes the final visibility mask. This flattens the inheritance tree and minimizes the computational overhead of object composition.
