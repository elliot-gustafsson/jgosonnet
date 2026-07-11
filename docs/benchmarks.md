## Benchmarks

Performance comparison of `jgosonnet` against the standard [Go implementation](https://github.com/google/go-jsonnet) (`go-jsonnet`) and the Rust implementation [jrsonnet](https://github.com/CertainLach/jrsonnet).

**Tested versions:**
* `go-jsonnet`: v0.22.0
* `jrsonnet`: 0.5.0-pre98

**Test Environment:**
* **CPU:** AMD Ryzen 7 7800X3D 8-Core Processor
* **OS:** Arch Linux (Kernel `7.0.14-arch1-1`)
* **Methodology:**  Benchmarks run via `go test -bench=. -benchtime=5x -count=3`. Values reported are the average execution times in milliseconds (ms) across the runs.

| Benchmark | `jgosonnet` [ms] | Rel | `jrsonnet` [ms] | Rel | `go-jsonnet` [ms] | Rel |
| :--- | ---: | ---: | ---: | ---: | ---: | ---: |
| [Large string join](../benchmarks/resources/large_string_join.jsonnet) | **7.48** | **1.00** | 7.66 | 1.02 | 36.36 | 4.86 |
| [Large string template](../benchmarks/resources/large_string_template.jsonnet) | 7.69 | 3.71 | **2.07** | **1.00** | *DNF* | *-* |
| [Realistic 1](../benchmarks/resources/realistic_benchmark1.jsonnet) | 10.08 | 1.34 | **7.50** | **1.00** | 3614.82 | 481.98 |
| [Realistic 2](../benchmarks/resources/realistic_benchmark2.jsonnet) | **97.08** | **1.00** | 163.53 | 1.68 | 3193.78 | 32.90 |
| [Tail call](../benchmarks/resources/tail_call.jsonnet) | 2.08 | 2.45 | **0.85** | **1.00** | 1.88 | 2.21 |
| [Inheritance recursion](../benchmarks/resources/inheritence_recursion.jsonnet) | **102.61** | **1.00** | 143.96 | 1.40 | 486.86 | 4.74 |
| [Comparisons primitives](../benchmarks/resources/comparisons_primitives.jsonnet) | **98.74** | **1.00** | 106.84 | 1.08 | 822.62 | 8.33 |

> **Note:** `go-jsonnet` failed with `exit status 2` (OS stack size exhaustion) on the `large_string_template` benchmark.
