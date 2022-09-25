# Tooliing on top of `go.mau.fi/whatsmeow`

This is my "handy package of stuff" to make the (excellent) Whatsapp Go library `go.mau.fi/whatsmeow` more usable in my programs.

## Handlers

`github.com/KarelKubat/whatsmeow/handlers` adds typed events, handler registration and by-type-dispatching to event handlers that you provide. To use it, instantiate a `whatsmeow` client as shown in https://pkg.go.dev/go.mau.fi/whatsmeow#NewClient. But event handlers are structured differently. Here is a canonical example of a simple `Message` handler:

```go
package message

import (
    "fmt"

    "github.com/KarelKubat/whatsmeow/handlers"
    "go.mau.fi/whatsmeow/types/events"
)

// The receiver can be a private, empty struct (or can hold data for processing a message).
type handler struct{}

// init() ensures the registration at start-up.
func init() {
    handlers.Register(handlers.Message, &handler{})
}

// Handle is invoked when a `handlers.Message` event is seen.
func (h *handler) Handle(ev interface{}) error {
    // The argument `ev` is an opaque event. The event handler must convert it to a
    // suitable type.
    m := ev.(*events.Message)

    // See the `events.Message` type for the usable fields.
    fmt.Println("Message:", m.Message.GetConversation())

    // A non-nil error return bubbles up.
    return nil
}
```

For the dispatching part, `handlers.Dispatch()` is set as the callback in https://pkg.go.dev/go.mau.fi/whatsmeow#Client.AddEventHandler. The default `EventHandler` type doesn't want an error return, so we can use an intermediate function:

```go
import (
    "github.com/KarelKubat/whatsmeow/handlers"
    "go.mau.fi/whatsmeow"
)
func main() {
    // Set up the container and device store as shown in
    // https://pkg.go.dev/go.mau.fi/whatsmeow#NewClient. Then:
    client := whatsmeow.NewClient(deviceStore, nil)
    client.AddEventHandler(func(e interface{}) {
        if err := handlers.Dispatch(e); err != nil {
            // Silently ignore when there is not a handler for this event. Else,
            // print the error.
            if err.Type == handlers.HandlerFailed {
                fmt.Fprintln(os.Stderr, err)
            }
        }
    })
}
```
NOTE: `github.com/KarelKubat/whatsmeow/handlers` currently doesn't support per-client handlers. All event handlers are global; once registered, they apply to all `whatsmeow.Client`s.
