# sekret

A simple Go package for handling sensitive data safely. Automatically masks secrets when printing, logging, or serializing to prevent accidental exposure.

Inspired by Rust's [secrecy](https://crates.io/crates/secrecy) crate.

## Features

- **Safe printing**: Secrets are masked when used with `fmt.Printf`, `fmt.Sprintf`, and similar operations
- **JSON/YAML safety**: Automatic masking during marshaling, with support for unmarshaling
- **Simple API**: Minimal interface with deliberate secret exposure

## Installation

```bash
go get github.com/lex-unix/sekret
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/yourusername/sekret"
)

func main() {
    password := sekret.New("supersecret123")

    // Safe - prints: Password: ******
    fmt.Printf("Password: %s\n", password)

    // Deliberate access when you need the actual value
    actualPassword := password.ExposeSecret()
    fmt.Printf("Connecting with: %s\n", actualPassword)
}
```

### Configuration Structs

```go
type Config struct {
    Username string              `json:"username" yaml:"username"`
    Password sekret.Sekret[string] `json:"password" yaml:"password"`
    APIKey   sekret.Sekret[string] `json:"api_key" yaml:"api_key"`
}

func main() {
    config := Config{
        Username: "admin",
        Password: sekret.New("secret123"),
        APIKey:   sekret.New("sk-1234567890"),
    }

    // Safe logging - secrets are masked
    fmt.Printf("Config: %+v\n", config)
    // Output: Config: {Username:admin Password:****** APIKey:******}
}
```

### JSON/YAML Support

```go
// JSON marshaling masks secrets
data, _ := json.Marshal(config)
fmt.Printf("JSON: %s\n", string(data))
// Output: {"username":"admin","password":"******","api_key":"******"}

// JSON unmarshaling works normally
var newConfig Config
json.Unmarshal([]byte(`{"username":"user","password":"pass123"}`), &newConfig)
fmt.Printf("Password: %s\n", newConfig.Password.ExposeSecret()) // pass123
```
