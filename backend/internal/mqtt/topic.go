package mqtt

import (
	"fmt"
	"regexp"
	"strings"
)

var publicKeyPattern = regexp.MustCompile(`^[0-9A-Fa-f]{8,128}$`)

type TopicInfo struct {
	IATA        string `json:"iata"`
	PublisherPK string `json:"publisherPublicKey"`
	Subtopic    string `json:"subtopic"`
}

func ParseTopic(topic string) (TopicInfo, error) {
	parts := strings.Split(strings.Trim(topic, "/"), "/")
	if len(parts) != 4 {
		return TopicInfo{}, fmt.Errorf("unexpected topic shape")
	}
	if parts[0] != "meshcore" {
		return TopicInfo{}, fmt.Errorf("topic does not start with meshcore")
	}
	iata := strings.ToUpper(parts[1])
	if len(iata) != 3 {
		return TopicInfo{}, fmt.Errorf("invalid IATA %q", parts[1])
	}
	pk := strings.ToUpper(parts[2])
	if !publicKeyPattern.MatchString(pk) {
		return TopicInfo{}, fmt.Errorf("invalid public key in topic")
	}
	subtopic := strings.ToLower(parts[3])
	switch subtopic {
	case "packets", "status", "debug", "internal":
		return TopicInfo{IATA: iata, PublisherPK: pk, Subtopic: subtopic}, nil
	default:
		return TopicInfo{}, fmt.Errorf("unsupported subtopic %q", subtopic)
	}
}
