package meshcore

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var hexCleaner = regexp.MustCompile(`[^0-9a-fA-F]`)

type ParsedPacket struct {
	PacketHash      string   `json:"packetHash"`
	RawHex          string   `json:"rawHex"`
	RouteType       int      `json:"routeType"`
	RouteTypeName   string   `json:"routeTypeName"`
	PayloadType     int      `json:"payloadType"`
	PayloadTypeName string   `json:"payloadTypeName"`
	PayloadVersion  int      `json:"payloadVersion"`
	TransportCodes  []byte   `json:"transportCodes,omitempty"`
	HashSize        int      `json:"hashSize"`
	HopCount        int      `json:"hopCount"`
	PathBytes       []byte   `json:"pathBytes"`
	PathChunks      []string `json:"pathChunks"`
	Payload         []byte   `json:"payload"`
	InvalidForMap   bool     `json:"invalidForMap"`
	InvalidReason   string   `json:"invalidReason,omitempty"`
}

type Advert struct {
	PublicKey      string   `json:"publicKey"`
	Timestamp      uint32   `json:"timestamp"`
	Flags          byte     `json:"flags"`
	NodeType       int      `json:"nodeType"`
	Role           string   `json:"role"`
	Latitude       *float64 `json:"latitude,omitempty"`
	Longitude      *float64 `json:"longitude,omitempty"`
	Name           string   `json:"name,omitempty"`
	LocationSource string   `json:"locationSource,omitempty"`
}

func ParseHexPacket(rawHex string) (ParsedPacket, error) {
	clean := strings.ToUpper(hexCleaner.ReplaceAllString(strings.TrimSpace(rawHex), ""))
	if clean == "" {
		return ParsedPacket{}, errors.New("empty raw packet hex")
	}
	if len(clean)%2 != 0 {
		return ParsedPacket{}, fmt.Errorf("raw packet hex has odd length")
	}
	raw, err := hex.DecodeString(clean)
	if err != nil {
		return ParsedPacket{}, fmt.Errorf("decode raw packet hex: %w", err)
	}
	return ParsePacket(raw)
}

func ParsePacket(raw []byte) (ParsedPacket, error) {
	if len(raw) < 2 {
		return ParsedPacket{}, fmt.Errorf("packet too short: %d bytes", len(raw))
	}

	header := raw[0]
	routeType := int(header & 0x03)
	payloadType := int((header >> 2) & 0x0F)
	payloadVersion := int((header >> 6) & 0x03)

	offset := 1
	var transport []byte
	if HasTransportCodes(routeType) {
		if len(raw) < offset+4+1 {
			return ParsedPacket{}, fmt.Errorf("packet too short for transport codes")
		}
		transport = append([]byte(nil), raw[offset:offset+4]...)
		offset += 4
	}

	if len(raw) < offset+1 {
		return ParsedPacket{}, fmt.Errorf("packet missing path length")
	}
	pathLength := raw[offset]
	offset++

	hashSize := int(pathLength>>6) + 1
	hopCount := int(pathLength & 0x3F)
	pathBytesLen := hashSize * hopCount
	if len(raw) < offset+pathBytesLen {
		return ParsedPacket{}, fmt.Errorf("packet path too short: need %d bytes, have %d", pathBytesLen, len(raw)-offset)
	}

	pathBytes := append([]byte(nil), raw[offset:offset+pathBytesLen]...)
	offset += pathBytesLen

	parsed := ParsedPacket{
		PacketHash:      PacketHash(raw),
		RawHex:          strings.ToUpper(hex.EncodeToString(raw)),
		RouteType:       routeType,
		RouteTypeName:   RouteTypeName(routeType),
		PayloadType:     payloadType,
		PayloadTypeName: PayloadTypeName(payloadType),
		PayloadVersion:  payloadVersion,
		TransportCodes:  transport,
		HashSize:        hashSize,
		HopCount:        hopCount,
		PathBytes:       pathBytes,
		PathChunks:      chunkPath(pathBytes, hashSize),
		Payload:         append([]byte(nil), raw[offset:]...),
	}

	if hashSize == 4 && payloadType != PayloadTrace {
		parsed.InvalidForMap = true
		parsed.InvalidReason = "reserved_4_byte_path_non_trace"
	}

	return parsed, nil
}

func chunkPath(path []byte, hashSize int) []string {
	if hashSize <= 0 || len(path) == 0 {
		return nil
	}
	out := make([]string, 0, len(path)/hashSize)
	for i := 0; i+hashSize <= len(path); i += hashSize {
		out = append(out, strings.ToUpper(hex.EncodeToString(path[i:i+hashSize])))
	}
	return out
}

func ParseAdvertPayload(payload []byte) (Advert, bool, error) {
	if len(payload) < 32+4+64 {
		return Advert{}, false, nil
	}

	publicKey := strings.ToUpper(hex.EncodeToString(payload[:32]))
	timestamp := binary.LittleEndian.Uint32(payload[32:36])
	appdata := payload[100:]

	advert := Advert{
		PublicKey: publicKey,
		Timestamp: timestamp,
		Role:      "unknown",
	}

	if len(appdata) == 0 {
		return advert, true, nil
	}

	flags := appdata[0]
	advert.Flags = flags
	advert.NodeType = int(flags & 0x0F)
	advert.Role = NodeRoleFromType(advert.NodeType)

	offset := 1
	if flags&0x10 != 0 {
		if len(appdata) < offset+8 {
			return advert, true, fmt.Errorf("advert location flag set but payload too short")
		}
		latRaw := int32(binary.LittleEndian.Uint32(appdata[offset : offset+4]))
		offset += 4
		lngRaw := int32(binary.LittleEndian.Uint32(appdata[offset : offset+4]))
		offset += 4
		lat := float64(latRaw) / 1_000_000
		lng := float64(lngRaw) / 1_000_000
		if lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180 {
			advert.Latitude = &lat
			advert.Longitude = &lng
			advert.LocationSource = "advert"
		}
	}

	if flags&0x20 != 0 {
		offset += 2
	}
	if flags&0x40 != 0 {
		offset += 2
	}
	if offset > len(appdata) {
		return advert, true, fmt.Errorf("advert appdata flags exceed payload")
	}
	if flags&0x80 != 0 && offset < len(appdata) {
		nameBytes := trimNUL(appdata[offset:])
		if utf8.Valid(nameBytes) {
			advert.Name = strings.TrimSpace(string(nameBytes))
		}
	}

	return advert, true, nil
}

func NodeRoleFromType(nodeType int) string {
	switch nodeType {
	case 1:
		return "companion"
	case 2:
		return "repeater"
	case 3:
		return "room_server"
	case 4:
		return "sensor"
	default:
		return "unknown"
	}
}

func trimNUL(in []byte) []byte {
	end := len(in)
	for end > 0 && in[end-1] == 0 {
		end--
	}
	return in[:end]
}
