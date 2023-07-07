package muonengine

import (
	"encoding/binary"
	"fmt"
	"io"
)

type messageID uint8

const (
    MsgChoke         messageID = 0
    MsgUnchoke       messageID = 1
    MsgInterested    messageID = 2
    MsgNotInterested messageID = 3
    MsgHave          messageID = 4
    MsgBitfield      messageID = 5
    MsgRequest       messageID = 6
    MsgPiece         messageID = 7
    MsgCancel        messageID = 8
)

type Message struct {
	ID messageID
	Data []byte
}

func (m *Message) Serialize() []byte {
    if m == nil {
        return make([]byte, 4)
    }
    length := uint32(len(m.Data) + 1)
    buf := make([]byte, 4+length)
    binary.BigEndian.PutUint32(buf[0:4], length)
    buf[4] = byte(m.ID)
    copy(buf[5:], m.Data)
    return buf
}

func MessageDeserialize(r io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}

	// keep-alive
	length := binary.BigEndian.Uint32(lengthBuf)
	if length == 0 {
		return nil, nil
	}

	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, err
	}

	m := Message{
		ID:      messageID(messageBuf[0]),
		Data: messageBuf[1:],
	}

	return &m, nil
}

func FormatRequest(index, begin, length int) *Message {
	data := make([]byte, 12)
	binary.BigEndian.PutUint32(data[0:4], uint32(index))
	binary.BigEndian.PutUint32(data[4:8], uint32(begin))
	binary.BigEndian.PutUint32(data[8:12], uint32(length))
	return &Message{ID: MsgRequest, Data: data}
}

func FormatHave(index int) *Message {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(index))
	return &Message{ID: MsgHave, Data: data}
}

func ParsePiece(index int, buf []byte, msg *Message) (int, error) {
	if msg.ID != MsgPiece {
		return 0, fmt.Errorf("Expected PIECE (ID %d), got ID %d", MsgPiece, msg.ID)
	}
	if len(msg.Data) < 8 {
		return 0, fmt.Errorf("Data too short. %d < 8", len(msg.Data))
	}
	parsedIndex := int(binary.BigEndian.Uint32(msg.Data[0:4]))
	if parsedIndex != index {
		return 0, fmt.Errorf("Expected index %d, got %d", index, parsedIndex)
	}
	begin := int(binary.BigEndian.Uint32(msg.Data[4:8]))
	if begin >= len(buf) {
		return 0, fmt.Errorf("Begin offset too high. %d >= %d", begin, len(buf))
	}
	data := msg.Data[8:]
	if begin+len(data) > len(buf) {
		return 0, fmt.Errorf("Data too long [%d] for offset %d with length %d", len(data), begin, len(buf))
	}
	copy(buf[begin:], data)
	return len(data), nil
}

func ParseHave(msg *Message) (int, error) {
	if msg.ID != MsgHave {
		return 0, fmt.Errorf("Expected to have (ID %d), got ID %d", MsgHave, msg.ID)
	}
	if len(msg.Data) != 4 {
		return 0, fmt.Errorf("Expected payload length 4, got length %d", len(msg.Data))
	}
	index := int(binary.BigEndian.Uint32(msg.Data))
	return index, nil
}

type Bitfield []byte

func (bf Bitfield) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8
	return (bf[byteIndex] >> (7-offset) & 1) != 0
}

func (bf Bitfield) SetPiece(index int) {
	byteIndex := index / 8
	offset := index % 8
	bf[byteIndex] |= 1 << (7 - offset)
}


