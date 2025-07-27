# 🕵️‍♀️ goleakwatch

Lightweight goroutine leak checker for Go.  
Catch leaks early in dev/test environments before they creep into production.

---

## 🚀 Features

- Detects leaked goroutines using `runtime.NumGoroutine`
- Configurable thresholds, wait durations, and stack traces
- Optional goroutine dump with `pprof`
- Simple test wrapper for CI-safe assertions

---

## 📦 Installation

```bash
go get github.com/incogni23/goleakwatch@latest
