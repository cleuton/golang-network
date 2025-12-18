<img src="./go_vs_node.png" height=400>

---

<img src="../../golangnetwork-logo.png" height=300>

---

[**Cleuton Sampaio**](https://linkedin.com/in/cleutonsampaio)

# Go vs NodeJS

Many people are discovering **Golang** as a high-performance alternative with low memory consumption and a small **footprint** for web applications and **REST** APIs. And it really is! I see many posts and reports of people converting **NodeJS** and **Python** apps to **Go**. But… is there really any gain in migrating from NodeJS to Go? Or is it just “gut feeling”?

> In software, **footprint** refers to the consumption of system resources such as RAM, CPU, and disk space that an application uses during execution. In cloud environments, this is crucial because costs are typically based on resource consumption. A smaller footprint means less memory, processing, and storage usage, which translates directly into lower operational costs. In addition, applications with a smaller footprint tend to be faster, more scalable, and more efficient, allowing more services to run on the same infrastructure and respond better under load. This is especially important in the cloud, where you pay for consumed resources and efficiency directly impacts the solution’s budget and performance.

**Let’s see what the benchmark says**

The website [**Benchmarks game**](https://benchmarksgame-team.pages.debian.net/benchmarksgame/fastest/go-node.html) has a comparative statistic, which I will summarize here:

1. **Performance:** Go demonstrates significantly superior performance in most computational benchmarks (e.g. `fannkuch-redux`, `n-body`, `mandelbrot`, `k-nucleotide`, `regex-redux`). Go programs generally have much lower execution times (`secs`) compared to Node.js programs. Node.js consistently appears among the worst performers in terms of execution speed for these types of tasks.
2. **Memory Usage:** Go uses much less RAM than Node.js for the same benchmarks. Memory values (`mem`) for Node.js are often in the hundreds of megabytes or even gigabytes (e.g. 111,055 KB, 448,524 KB, 1,897,566 KB), while Go programs operate with tens or a few hundreds of megabytes (e.g. 3,846 KB, 168,186 KB, 416,195 KB).
3. **CPU Efficiency:** Computationally intensive benchmarks show that optimized Go programs can use multiple CPU cores more effectively, often with `cpu secs` close to or greater than real time (`secs`), indicating good parallelism. Node.js, being single-threaded with an event loop, does not scale in the same way in these scenarios.
4. **Deployment (Artifact Size):**

   * **Go:** Produces a single **static executable binary** after compilation. The size of this binary is relatively small.
   * **Node.js:** Requires distribution of the **application source code** and the **`node_modules`** folder containing all third-party dependencies. The size of the `node_modules` folder can vary significantly.

This performance test is a comparison I used to provide evidence to a client who is undecided about migrating from NodeJS to Golang. In fact, this is an ID generation test using **Sony Snowflake**, implemented in **Golang** (using `FastHTTP`) and in **JavaScript** (using `NodeJS` and `Express`).

The source codes are in the [**`servers` folder**](./performance-test/servers/).

I used [**K6**](https://k6.io/) to run the tests and created `shell` and `python` scripts to collect results and metrics.

The test parameters I used were:

* Ramp Up: 30 seconds up to 50 VUs.
* 1-minute test.
* Ramp Down: 30 seconds down to 0.
* Thinking time of 0.05s.

**Workers / processes**

1. **Go (`go_fasthttp`):**

   * The Go code using `fasthttp` does not explicitly create multiple processes or threads to handle concurrent requests in the provided code. The `fasthttp.ListenAndServe` function internally handles concurrency, but the exact number of goroutines or worker threads is not fixed in the code and is managed by the `fasthttp` library itself in an optimized way. Generally, Go libraries such as `fasthttp` use a lightweight concurrency model (goroutines) and an operating system thread pool, scaling according to demand. The number of effective “workers” is dynamic and not specified as a fixed number in the code.

2. **Node.js (`nodejs_express`):**

   * The Node.js code was modified to use the `cluster` module. The line `const numCPUs = require('os').cpus().length;` determines the number of CPUs available on the system.
   * The loop `for (let i = 0; i < numCPUs; i++) { cluster.fork(); }` explicitly creates **one child worker process per available CPU core**.
   * Therefore, if the system where the test was executed has, for example, 4 CPU cores, the Node.js version started **4 workers** (in addition to 1 master process).

* **Go (`go_fasthttp`):** The number of workers is not explicitly defined in the code; it is managed internally by the `fasthttp` library.
* **Node.js (`nodejs_express`):** The number of workers is equal to the number of CPU cores of the machine where the test was executed (for example, 4 workers on a 4-core machine).

Runs on Linux only, for now!

## To run the tests:

1. Install **Python** dependencies (I recommend creating a virtual environment with `venv`):

```
pip install -r requirements.txt
```

2. Install k6: [https://k6.io/docs/getting-started/installation/](https://k6.io/docs/getting-started/installation/)

3. Grant execution permission to the scripts:

```
chmod +x run_tests.sh
chmod +x monitoring/resource_monitor.sh
```

## Run tests

```
./run_tests.sh
```

## Analyze results

```
python3 analysis/analyze_results.py
```

## Final results

After running it several times, I arrived at these results:

Based on the analysis of the provided benchmark and the specific results of your performance test between Go (`go_fasthttp`) and Node.js (`nodejs_express`), here is a summary of the most important comparisons:

**Comparison Summary: Go vs. Node.js**

This specific test compared `go_fasthttp` with `nodejs_express` (using clustering). The results consistently showed that, for this use case (sequential ID generation via a REST API):

1. **Performance (Response Time):** Go (`go_fasthttp`) was significantly faster, with average response times almost 3 times lower (484 ms vs. 1299 ms).
2. **CPU Efficiency:** Go used less CPU on average (3.63% vs. 7.47%) and had much lower peak usage (7.90% vs. 25.70%).
3. **Memory Efficiency:** Go consumed less RAM on average (6682 MB vs. 7327 MB) and showed much more stable memory usage (standard deviation of 22 MB vs. 244 MB).
4. **Deployment:** Go produces a single compact executable binary (7.28 MB). Node.js requires the runtime environment (`node`) and the dependencies (`node_modules`, which totaled 2.14 MB in this case, but still require the Node.js runtime to function).

So, **Go** proved to be an option with lower consumption, smaller footprint, and better performance than JavaScript. Note that we deliberately chose an application that is not even **CPU-intensive** for this test, whose memory consumption should be small. 
