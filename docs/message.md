## Package `muonengine`

The `muonengine` package provides functionality for a Go project related to a download client for a peer-to-peer network.

### Constants

-   `MsgChoke`: Represents the "Choke" message ID.
-   `MsgUnchoke`: Represents the "Unchoke" message ID.
-   `MsgInterested`: Represents the "Interested" message ID.
-   `MsgNotInterested`: Represents the "Not Interested" message ID.
-   `MsgHave`: Represents the "Have" message ID.
-   `MsgBitfield`: Represents the "Bitfield" message ID.
-   `MsgRequest`: Represents the "Request" message ID.
-   `MsgPiece`: Represents the "Piece" message ID.
-   `MsgCancel`: Represents the "Cancel" message ID.

### Type `Message`

```go
type Message struct {
    ID   messageID
    Data []byte
}
```

The `Message` struct represents a message used in the peer-to-peer protocol.

#### Fields

-   `ID`: The message ID.
-   `Data`: The data associated with the message.

### Method `Serialize` (Message)

```go
func (m *Message) Serialize() []byte
```

The `Serialize` method serializes the message into a byte slice.

#### Returns

-   The serialized message as a byte slice.

### Function `MessageDeserialize`

```go
func MessageDeserialize(r io.Reader) (*Message, error)
```

The `MessageDeserialize` function deserializes a message from an `io.Reader` and returns a `Message` object.

#### Parameters

-   `r`: The `io.Reader` to read the message from.

#### Returns

-   A pointer to the deserialized `Message` object.
-   An error if deserialization fails.

### Function `FormatRequest`

```go
func FormatRequest(index, begin, length int) *Message
```

The `FormatRequest` function formats a "Request" message with the specified index, begin offset, and length.

#### Parameters

-   `index`: The index of the piece being requested.
-   `begin`: The beginning offset of the requested block.
-   `length`: The length of the requested block.

#### Returns

-   A pointer to the formatted "Request" message.

### Function `FormatHave`

```go
func FormatHave(index int) *Message
```

The `FormatHave` function formats a "Have" message with the specified index.

#### Parameters

-   `index`: The index of the piece that the sender has.

#### Returns

-   A pointer to the formatted "Have" message.

### Function `ParsePiece`

```go
func ParsePiece(index int, buf []byte, msg *Message) (int, error)
```

The `ParsePiece` function parses a "Piece" message and copies the data to the specified buffer.

#### Parameters

-   `index`: The index of the piece being received.
-   `buf`: The buffer to copy the data into.
-   `msg`: The received "Piece" message.

#### Returns

-   The number of bytes copied to the buffer.
-   An error if parsing or copying fails.

### Function `ParseHave`

```go
func ParseHave(msg *Message) (int, error)
```

The `ParseHave` function parses a "Have" message and returns the index of the piece.

#### Parameters

-   `msg`: The received "Have" message.

#### Returns

-   The index of the piece that the sender has.
-   An error if parsing fails.

### Type `Bitfield`

```go
type Bitfield []byte
```

The `Bitfield` type represents a bitfield indicating the availability of pieces.

#### Methods

-   `HasPiece(index int) bool`: Checks if a specific piece is available in the bitfield.
-   `SetPiece(index int)`: Sets the availability of a specific piece in the bitfield.
