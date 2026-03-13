# go-forms

## Examples

```go
package main

import (
    "log/slog"
    "net/http"
    "os"

    "github.com/schlauerlauer/go-forms"
)

type formData struct {
    Name  string `schema:"name" validate:"required" mod:"trim,sanitize"`
    Value int64  `schema:"value" validate:"gte=0,lte=100"`
}

func main() {
    // setup go-forms
    formProcessor, err := forms.NewFormProcessor()
    if err != nil {
        slog.Error("Error setting up form processor", "err", err)
        os.Exit(1)
    }

    // setup routing
    router := http.NewServeMux()
    router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {

        // process form data
        var data formData
        if err := formProcessor.ProcessForm(&data, r); err != nil {
            slog.Error("ProcessForm", "err", err.Error())
            return
        }

        slog.Info("form processed successfully", "data", data)
    })

    // start server
    if err := http.ListenAndServe("127.0.0.1:8080", router); err != nil {
        slog.Error("Server error", "err", err)
        os.Exit(1)
    }
}
```
