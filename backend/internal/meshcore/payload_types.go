package meshcore

const (
	PayloadRequest   = 0x00
	PayloadResponse  = 0x01
	PayloadPlainText = 0x02
	PayloadAck       = 0x03
	PayloadAdvert    = 0x04
	PayloadGroupText = 0x05
	PayloadGroupData = 0x06
	PayloadAnonReq   = 0x07
	PayloadPath      = 0x08
	PayloadTrace     = 0x09
	PayloadMultipart = 0x0A
	PayloadControl   = 0x0B
	PayloadRawCustom = 0x0F
)

const (
	RouteTransportFlood  = 0x00
	RouteFlood           = 0x01
	RouteDirect          = 0x02
	RouteTransportDirect = 0x03
)

func PayloadTypeName(v int) string {
	switch v {
	case PayloadRequest:
		return "REQUEST"
	case PayloadResponse:
		return "RESPONSE"
	case PayloadPlainText:
		return "PLAIN_TEXT"
	case PayloadAck:
		return "ACK"
	case PayloadAdvert:
		return "ADVERT"
	case PayloadGroupText:
		return "GROUP_TEXT"
	case PayloadGroupData:
		return "GROUP_DATA"
	case PayloadAnonReq:
		return "ANON_REQUEST"
	case PayloadPath:
		return "RETURNED_PATH"
	case PayloadTrace:
		return "TRACE"
	case PayloadMultipart:
		return "MULTIPART"
	case PayloadControl:
		return "CONTROL"
	case PayloadRawCustom:
		return "CUSTOM"
	default:
		return "RESERVED"
	}
}

func RouteTypeName(v int) string {
	switch v {
	case RouteTransportFlood:
		return "TRANSPORT_FLOOD"
	case RouteFlood:
		return "FLOOD"
	case RouteDirect:
		return "DIRECT"
	case RouteTransportDirect:
		return "TRANSPORT_DIRECT"
	default:
		return "UNKNOWN"
	}
}

func HasTransportCodes(routeType int) bool {
	return routeType == RouteTransportFlood || routeType == RouteTransportDirect
}
