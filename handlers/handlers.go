// Package handlers adds typed events, handler registration and dispatching on top of
// `go.mau.fi/whatsmeow`.
package handlers

import (
	"fmt"
	"sync"

	"go.mau.fi/whatsmeow/types/events"
)

// EventType is an enum for whatsmeow events.
type EventType int

const (
	firstEventType EventType = iota // Keep at first slot for tests

	AppState
	AppStateSyncComplete
	Archive
	BusinessName
	CallAccept
	CallOffer
	CallOfferNotice
	CallRelayLatency
	CallTerminate
	ChatPresence
	ClientOutdated
	Connected
	ConnectFailure
	Contact
	DeleteChat
	DeleteForMe
	Disconnected
	GroupInfo
	HistorySync
	IdentityChange
	JoinedGroup
	KeepAliveRestored
	KeepAliveTimeout
	LoggedOut
	MarkChatAsRead
	MediaRetry
	Message
	Mute
	OfflineSyncCompleted
	OfflineSyncPreview
	PairError
	PairSuccess
	Picture
	Pin
	Presence
	PrivacySettings
	PushName
	PushNameSetting
	QR
	QRScannedWithoutMultidevice
	Receipt
	Star
	StreamError
	StreamReplaced
	TemporaryBan
	UnarchiveChatSetting
	UndecryptableMessage
	UnknownCallEvent

	lastEventType // Keep at last slot for tests
)

// String returns the string representation of a Type.
func (t EventType) String() string {
	return []string{
		"", // unused
		"AppState",
		"AppStateSyncComplete",
		"Archive",
		"BusinessName",
		"CallAccept",
		"CallOffer",
		"CallOfferNotice",
		"CallRelayLatency",
		"CallTerminate",
		"ChatPresence",
		"ClientOutdated",
		"Connected",
		"ConnectFailure",
		"Contact",
		"DeleteChat",
		"DeleteForMe",
		"Disconnected",
		"GroupInfo",
		"HistorySync",
		"IdentityChange",
		"JoinedGroup",
		"KeepAliveRestored",
		"KeepAliveTimeout",
		"LoggedOut",
		"MarkChatAsRead",
		"MediaRetry",
		"Message",
		"Mute",
		"OfflineSyncCompleted",
		"OfflineSyncPreview",
		"PairError",
		"PairSuccess",
		"Picture",
		"Pin",
		"Presence",
		"PrivacySettings",
		"PushName",
		"PushNameSetting",
		"QR",
		"QRScannedWithoutMultidevice",
		"Receipt",
		"Star",
		"StreamError",
		"StreamReplaced",
		"TemporaryBan",
		"UnarchiveChatSetting",
		"UndecryptableMessage",
		"UnknownCallEvent",
	}[t]
}

type handler interface {
	Handle(evt interface{}) error
}

var registry = make(map[EventType][]handler)
var registryMutex sync.Mutex

// Register registers a handler for an event type. The handler must expose a method
//
//	Handle(evt interface{}) error
//
// The passed-in event to the handler method is an opaque pointer to one of the types
// of `go.mau.fi/whatsmeow/types/events`. The handler must convert it to the true event
// using a typecast. For example, a handler for the type `Message` would convert as
// follows:
//
//	import "go.mau.fi/whatsmeow/types/events"
//	func (h *myHandler) Handle(ev interface{}) error {
//	  messageEvent := ev.(*events.Message)
//	  ...
//	}
//
// More than one handlers may be registered for an event. Upon encountering the event,
// the handlers will be called in-order.
//
//	Register(Message, h1)
//	Register(Message, h2)
//	// When a `Message` is seen, first `h1.Handle(ev)` is invoked, then `h2.Handle(ev)`.
func Register(t EventType, h handler) {
	registryMutex.Lock()
	defer registryMutex.Unlock()

	if _, ok := registry[t]; !ok {
		registry[t] = []handler{}
	}
	registry[t] = append(registry[t], h)
}

type dispatchErrorType int

const (
	firstDispatchError dispatchErrorType = iota // Keep at first slot for tests

	NoHandlerFound
	HandlerFailed

	lastDispatchError // Keep at last slot for tests
)

func (d dispatchErrorType) String() string {
	return []string{
		"",
		"NoHandlerFound",
		"HandlerFailed",
	}[d]
}

// DispatchError enriches the error returned by Dispatch with an error reason, which may be
// `NoHandlerFound` or `HandlerFailed`. Example:
//
//	 if err := Dispatch(e); err != nil {
//		  if err.Type == NoHandlerFound {
//			// Log but ignore that there is no handler for this event
//			log.Println(err)
//		  } else {
//			// A handler ran, but failed
//			log.Fatalln(err)
//		  }
//	 }
type DispatchError struct {
	Type dispatchErrorType
	Err  error
}

func (d *DispatchError) Error() string {
	return d.Err.Error()
}

// Dispatch invokes registered handlers for any `EventType`. There is a `nil` error return
// IFF:
// - One or more handlers for the event type were registered,
// - They all executed without returning an error.
//
// The absence of registered handlers is returned with `err.Type == NoHandlerFound` and
// `err.Error()` stating the event type and payload.
//
// The failure of a registered handler is returned with `err.Type` == HandlerFailed`,
// `err.Err` being the underlying error, and `err.Error()` stating the handler's error.
// Invoking handlers stops when a handler returns an error; i.e., a second handler may not run
// if the first handler fails.
func Dispatch(evt interface{}) error {
	switch v := evt.(type) {
	case *events.AppState:
		return dispatch(AppState, v)
	case *events.AppStateSyncComplete:
		return dispatch(AppStateSyncComplete, v)
	case *events.Archive:
		return dispatch(Archive, v)
	case *events.BusinessName:
		return dispatch(BusinessName, v)
	case *events.CallAccept:
		return dispatch(CallAccept, v)
	case *events.CallOffer:
		return dispatch(CallOffer, v)
	case *events.CallOfferNotice:
		return dispatch(CallOfferNotice, v)
	case *events.CallRelayLatency:
		return dispatch(CallRelayLatency, v)
	case *events.CallTerminate:
		return dispatch(CallTerminate, v)
	case *events.ChatPresence:
		return dispatch(ChatPresence, v)
	case *events.ClientOutdated:
		return dispatch(ClientOutdated, v)
	case *events.Connected:
		return dispatch(Connected, v)
	case *events.ConnectFailure:
		return dispatch(ConnectFailure, v)
	case *events.Contact:
		return dispatch(Contact, v)
	case *events.DeleteChat:
		return dispatch(DeleteChat, v)
	case *events.DeleteForMe:
		return dispatch(DeleteForMe, v)
	case *events.Disconnected:
		return dispatch(Disconnected, v)
	case *events.GroupInfo:
		return dispatch(GroupInfo, v)
	case *events.HistorySync:
		return dispatch(HistorySync, v)
	case *events.JoinedGroup:
		return dispatch(JoinedGroup, v)
	case *events.IdentityChange:
		return dispatch(IdentityChange, v)
	case *events.KeepAliveRestored:
		return dispatch(KeepAliveRestored, v)
	case *events.KeepAliveTimeout:
		return dispatch(KeepAliveTimeout, v)
	case *events.LoggedOut:
		return dispatch(LoggedOut, v)
	case *events.MarkChatAsRead:
		return dispatch(MarkChatAsRead, v)
	case *events.MediaRetry:
		return dispatch(MediaRetry, v)
	case *events.Message:
		return dispatch(Message, v)
	case *events.OfflineSyncCompleted:
		return dispatch(OfflineSyncCompleted, v)
	case *events.OfflineSyncPreview:
		return dispatch(OfflineSyncPreview, v)
	case *events.PairError:
		return dispatch(PairError, v)
	case *events.PairSuccess:
		return dispatch(PairSuccess, v)
	case *events.Picture:
		return dispatch(Picture, v)
	case *events.Pin:
		return dispatch(Pin, v)
	case *events.Presence:
		return dispatch(Presence, v)
	case *events.PrivacySettings:
		return dispatch(PrivacySettings, v)
	case *events.PushName:
		return dispatch(PushName, v)
	case *events.PushNameSetting:
		return dispatch(PushNameSetting, v)
	case *events.QR:
		return dispatch(QR, v)
	case *events.QRScannedWithoutMultidevice:
		return dispatch(QRScannedWithoutMultidevice, v)
	case *events.Receipt:
		return dispatch(Receipt, v)
	case *events.Star:
		return dispatch(Star, v)
	case *events.StreamError:
		return dispatch(StreamError, v)
	case *events.StreamReplaced:
		return dispatch(StreamReplaced, v)
	case *events.TemporaryBan:
		return dispatch(TemporaryBan, v)
	case *events.UnarchiveChatsSetting:
		return dispatch(UnarchiveChatSetting, v)
	case *events.UndecryptableMessage:
		return dispatch(UndecryptableMessage, v)
	case *events.UnknownCallEvent:
		return dispatch(UnknownCallEvent, v)
	default:
		return fmt.Errorf("unknown event %+v, can't dispatch", v)
	}
}

func dispatch(t EventType, ev interface{}) error {
	if handlers, ok := registry[t]; ok {
		for _, h := range handlers {
			if err := h.Handle(ev); err != nil {
				return &DispatchError{
					Type: HandlerFailed,
					Err:  err,
				}
			}
		}
		return nil
	}
	return &DispatchError{
		Type: NoHandlerFound,
		Err:  fmt.Errorf("no handler for event %v (payload: %+v)", t, ev),
	}
}
