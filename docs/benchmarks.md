## Benchmarks

Performance comparison of `jgosonnet` against the standard [Go implementation](https://github.com/google/go-jsonnet) (`go-jsonnet`) and the Rust implementation [jrsonnet](https://github.com/CertainLach/jrsonnet).

**Tested versions:**
* `go-jsonnet`: v0.22.0
* `jrsonnet`: 0.5.0-pre98

**Test Environment:**
* **CPU:** AMD Ryzen 7 7800X3D 8-Core Processor
* **OS:** Arch Linux (Kernel `7.1.5-arch1-2`)
* **GO:** 1.26.5
* **Methodology:**  Benchmarks run via `go test -bench=. -benchtime=5x -count=3`. Values reported are the average execution times in milliseconds (ms) across the runs.

| Benchmark | `jgosonnet` [ms] | Rel | `jrsonnet` [ms] | Rel | `go-jsonnet` [ms] | Rel |
| :--- | ---: | ---: | ---: | ---: | ---: | ---: |
| [Large string join](../benchmarks/resources/large_string_join.jsonnet) | **5.50** | **1.00** | 7.60 | 1.38 | 37.27 | 6.78 |
| [Large string template](../benchmarks/resources/large_string_template.jsonnet) | 7.14 | 3.47 | **2.06** | **1.00** | *DNF* | *-* |
| [Realistic 1](../benchmarks/resources/realistic_benchmark1.jsonnet) | **6.16** | **1.00** | 7.75 | 1.26 | 3612.12 | 586.38 |
| [Realistic 2](../benchmarks/resources/realistic_benchmark2.jsonnet) | **71.84** | **1.00** | 163.04 | 2.27 | 3176.99 | 44.22 |
| [Tail call](../benchmarks/resources/tail_call.jsonnet) | 1.76 | 2.07 | **0.85** | **1.00** | 1.86 | 2.19 |
| [Inheritance recursion](../benchmarks/resources/inheritence_recursion.jsonnet) | **81.31** | **1.00** | 143.99 | 1.77 | 485.56 | 5.97 |
| [Comparisons primitives](../benchmarks/resources/comparisons_primitives.jsonnet) | **86.55** | **1.00** | 106.56 | 1.23 | 812.20 | 9.38 |

> **Note:** `go-jsonnet` failed with `exit status 2` (OS stack size exhaustion) on the `large_string_template` benchmark.
