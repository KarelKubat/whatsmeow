# Tooling on top of `go.mau.fi/whatsmeow`

This is my "handy package of stuff" to make the (excellent) Whatsapp Go library `go.mau.fi/whatsmeow` more usable in my programs.

Why may this be handy? If you run the original reference client https://godocs.io/go.mau.fi/whatsmeow#example-package then:
- Event handling is implemented using a typecast from an opaque `interface{}`. I'd prefer pre-registered event handlers, with a strict event type.
- Logging goes to `stdout`. I'd prefer it to go to a file for later inspection.

## Handlers

`github.com/KarelKubat/whatsmeow/handlers` adds typed events, handler registration and by-type-dispatching to event handlers that you provide. The usage is very similar to the `whatsmeow` client as shown in https://pkg.go.dev/go.mau.fi/whatsmeow#NewClient. But:

- `github.com/KarelKubat/whatsmeow/handlers` adds an `EventType` that enumerates the kinds of events;
- Event handlers are structured differently; they "bind" themselves to a given event type when registering;
- You can bind multiple handlers to one event. If so, they are all executed (in order of binding).
- Event dispatching is driven by the registry of available handlers. The dispatcher determines the type of the event and calls the appropriate handler(s).

### Anatomy of a handler

Here is a simple example of a simple `Message` handler:

```go
package message

import (
    "fmt"

    "github.com/KarelKubat/whatsmeow/handlers"
    "go.mau.fi/whatsmeow/types/events"
)

// The receiver can be a private, empty struct (or can hold data for processing a message).
type handler struct{}

// init() ensures the registration at start-up. Just add:
//  import _ "path/to/myhandlers/message"
// to your codebase.
func init() {
    handlers.Register(handlers.Message, &handler{})
}

// Handle is invoked when a `handlers.Message` event is seen.
func (h *handler) Handle(ev interface{}) error {
    // The argument `ev` is opaque. The event handler must convert it to a suitable type.
    m := ev.(*events.Message)

    // See the `events.Message` type for the usable fields.
    fmt.Println("Message:", m.Message.GetConversation())

    // A non-nil error return bubbles up.
    return nil
}
```

A more complete handler for the type `events.Message` can be found in https://github.com/KarelKubat/whapp/blob/main/handlers/message/message.go.

### Dispatching

Dispatching occurs through `handlers.Dispatch()`. This method matches an event against the registered handlers and, if one or more handlers are found, calls them. It returns a `nil` error or a `handlers.DispatchError`. 

The dispatch error (if any) has a field `Type` which is set to either:
- `NoHandlerFound`: there was no registered handler for this event. You might not install handlers for all possible events and want to ignore this error.
- `HandlerFailed`: A handler for the given event returned an error. When multiple handlers are bound to an event type, then the serial execution of the handlers stops when once one of them errors out.
- `UnknownEvent`: The dispatcher isn't configured to handle the event. This is a bug or it may mean that a new event type was implemented by https://github.com/tulir/whatsmeow/tree/main/types/events that the dispatcher doesn't know (yet).

The dispatcher is set as the callback as shown in in https://pkg.go.dev/go.mau.fi/whatsmeow#Client.AddEventHandler. The default `whatsmeow.EventHandler` type doesn't want an error return, so we can use an intermediate function:

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
            switch err.Type {
            case handlers.NoHandlerFound:
                // Silently ignore when there is not a handler for this event.
                return nil
            case handlers.HandlerFailed:
                // If the handler returned an error, print it.
                fmt.Fprintln(os.Stderr, err)
            case handlers.UnknownEvent:
                // If the dispatcher doesn't know of this event, pull the emergency brake.
                panic(err)
            }
        }
    })
}
```

For a real life example, see https://github.com/KarelKubat/whapp/blob/main/whapp.go.

> NOTE: `github.com/KarelKubat/whatsmeow/handlers` currently doesn't support per-client handlers. All event handlers are global; once registered, they apply to all `whatsmeow.Client`s.

## File Logging

`github.com/KarelKubat/whatsmeow/logger` implements the interface `go.mau.fi/whatsmeow/util/log` but instead of sending logging to `stdout`, it is sent to a file. The file can be "rotated-away" in the middle of a run; the logger ensures that when the logfile disappears, a new one is created.

Example:

```go
import (
    "github.com/KarelKubat/whatsmeow/logger"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/store/sqlstore"

    _ "github.com/mattn/go-sqlite3"    
)

const (
    logfile = "/tmp/my.log"
)

func main() {
    // Base logger, without a module name.
    baseLogger, err := logger.New(logger.Opts{
        Filename: logfile,  // file to send messages to
        Verbose:  true,     // when true, debug messages are sent too
        Append:   true,     // when true, previous logs are not overwritten but appended to
    })
    if err != nil { handleError(err) }
    // Sub loggers populate the messages with a module name.
    dbLogger := baseLogger.Sub("Database")
    clLogger := baseLogger.Sub("Client")

    // Store instantiation
    container, err := sqlstore.New("sqlite3", "file:store.db?_foreign_keys=on", dbLogger)
    store, err := container.GetFirstDevice()

    // Client instantiation
    client := whatsmeow.NewClient(store, clLogger)

    // ...
}
```

After this, database actions and client actions will be logged to `/tmp/my.log`.

> NOTE: This package supports neither opening loggers to output to different files (everything must go to one file), nor modifying the verbosity level. This can of course be implemented.
