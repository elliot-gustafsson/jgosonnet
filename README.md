# jgosonnet

# OBS Not feature complete yet

A high-performance evaluator for [Jsonnet](https://jsonnet.org/).
This implementation is built from the ground up to be significantly faster than `go-jsonnet`. The primary motivation for this project is that `go-jsonnet` is too slow and consumes too much memory when evaluating exceptionally large files or highly complex, deeply nested configurations.

## Architecture

### Lexing and Parsing
For the frontend, `jgosonnet` utilizes the upstream `github.com/google/go-jsonnet` implementation to handle lexing and parsing. The source code is parsed by `go-jsonnet` into an Abstract Syntax Tree (AST). Reusing the upstream parser ensures strict syntactic compatibility, while allowing this project to focus purely on completely overhauling the execution engine and memory model for speed.

### String Interning
Object keys and string values are passed through a central String Interner during evaluation. By translating strings into 32-bit reference IDs (`refId`), the evaluator avoids repetitive string allocations. String equality checks—which happen millions of times during object resolution and sorting—are reduced to simple integer comparisons.

### Arena Allocation
Rather than independently allocating every JSON object, array, and function on the Go heap, the evaluator utilizes an Arena allocator. Runtime representations of Objects, Arrays, Thunks, and Functions are stored in contiguous memory slices within an evaluation `Context`. This drastically reduces the pressure on Go's Garbage Collector (GC) and improves cache locality, which is critical when parsing massive data structures.

### Lazy Evaluation (Thunks)
Adherence to Jsonnet's lazy evaluation semantics is achieved using Thunks. A Thunk bundles an AST node with its captured lexical scope and an identifier. Computations are deferred and only evaluated when explicitly requested (e.g., during final JSON/YAML manifestation or when strict typing is required by standard library functions).

### Object Resolution and Layers
Jsonnet's object model supports complex inheritance (`+`, `super`, `self`) and field visibility constraints (`::`, `:::`, `:`). The architecture handles this by flattening inherited objects into `Layers` and compiling a `FieldPlan`. Instead of deeply copying or continuously merging objects during AST evaluation, the evaluator constructs a plan that points to the correct field layer and computes the final visibility mask. This flattens the inheritance tree and minimizes the computational overhead of object composition.
