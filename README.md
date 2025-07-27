# 🕵️‍♀️ goleakwatch

[![CI](https://github.com/incogni23/goleakwatch/workflows/CI/badge.svg)](https://github.com/incogni23/goleakwatch/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/incogni23/goleakwatch)](https://goreportcard.com/report/github.com/incogni23/goleakwatch)
[![Go Version](https://img.shields.io/github/go-mod/go-version/incogni23/goleakwatch)](https://go.dev/)
[![License](https://img.shields.io/github/license/incogni23/goleakwatch)](https://github.com/incogni23/goleakwatch/blob/main/LICENSE)

Lightweight goroutine leak checker for Go.  
Catch leaks early in dev/test environments before they creep into production.

---

## 🚀 Features

- **Detects leaked goroutines** using `runtime.NumGoroutine`
- **Configurable thresholds**, wait durations, and stack traces
- **Optional goroutine dump** with `pprof`
- **Context support** for cancellation and timeout control
- **Benchmarking utilities** for performance measurement
- **Simple test wrapper** for CI-safe assertions
- **Memory usage tracking** in benchmarks

---

## 📦 Installation

```bash
go get github.com/incogni23/goleakwatch@latest
```

---

## 🎯 Quick Start

### Basic Usage
```go
import "github.com/incogni23/goleakwatch"

err := goleakwatch.DefaultCheck(func() {
    // Your code here
    go someFunction()
})
if err != nil {
    log.Printf("Leak detected: %v", err)
}
```

### With Custom Configuration
```go
err := goleakwatch.Check(func() {
    // Your code here
}, goleakwatch.Config{
    Threshold:   2,                    // Allow 2 extra goroutines
    Wait:        500 * time.Millisecond,
    EnableTrace: true,                 // Dump stack traces
    Timeout:     10 * time.Second,     // 10 second timeout
})
```

### In Tests
```go
func TestMyFunction(t *testing.T) {
    goleakwatch.WithTest(t, func() {
        // Your test code here
    })
}
```

### Benchmarking
```go
result := goleakwatch.Benchmark("myFunction", func() {
    // Function to benchmark
})
fmt.Println(result) // Prints detailed metrics
```

---

## 📚 Documentation

- [API Documentation](docs/API.md) - Complete API reference
- [Examples](examples/) - Working code examples
- [Best Practices](docs/API.md#best-practices) - Usage guidelines

---

## 🔧 Configuration

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `Threshold` | `int` | `1` | Max allowed goroutine difference |
| `Wait` | `time.Duration` | `200ms` | Wait time after function runs |
| `EnableTrace` | `bool` | `true` | Dump goroutine trace if leak suspected |
| `Out` | `io.Writer` | `os.Stderr` | Where to write pprof dump |
| `Timeout` | `time.Duration` | `5s` | Timeout for the entire check operation |

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
