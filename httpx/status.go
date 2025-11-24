package httpx

type Status struct {
	Code  int    `json:"code"`
	Label string `json:"label"`
}

var (
	StatusContinue           = Status{Code: 100, Label: "Continue"}            // RFC 9110, 15.2.1
	StatusSwitchingProtocols = Status{Code: 101, Label: "Switching Protocols"} // RFC 9110, 15.2.2
	StatusProcessing         = Status{Code: 102, Label: "Processing"}          // RFC 2518, 10.1
	StatusEarlyHints         = Status{Code: 103, Label: "Early Hints"}         // RFC 8297

	StatusOK                   = Status{Code: 200, Label: "OK"}                            // RFC 9110, 15.3.1
	StatusCreated              = Status{Code: 201, Label: "Created"}                       // RFC 9110, 15.3.2
	StatusAccepted             = Status{Code: 202, Label: "Accepted"}                      // RFC 9110, 15.3.3
	StatusNonAuthoritativeInfo = Status{Code: 203, Label: "Non-Authoritative Information"} // RFC 9110, 15.3.4
	StatusNoContent            = Status{Code: 204, Label: "No Content"}                    // RFC 9110, 15.3.5
	StatusResetContent         = Status{Code: 205, Label: "Reset Content"}                 // RFC 9110, 15.3.6
	StatusPartialContent       = Status{Code: 206, Label: "Partial Content"}               // RFC 9110, 15.3.7
	StatusMultiStatus          = Status{Code: 207, Label: "Multi-Status"}                  // RFC 4918, 11.1
	StatusAlreadyReported      = Status{Code: 208, Label: "Already Reported"}              // RFC 5842, 7.1
	StatusIMUsed               = Status{Code: 226, Label: "IM Used"}                       // RFC 3229, 10.4.1

	StatusMultipleChoices   = Status{Code: 300, Label: "Multiple Choices"}   // RFC 9110, 15.4.1
	StatusMovedPermanently  = Status{Code: 301, Label: "Moved Permanently"}  // RFC 9110, 15.4.2
	StatusFound             = Status{Code: 302, Label: "Found"}              // RFC 9110, 15.4.3
	StatusSeeOther          = Status{Code: 303, Label: "See Other"}          // RFC 9110, 15.4.4
	StatusNotModified       = Status{Code: 304, Label: "Not Modified"}       // RFC 9110, 15.4.5
	StatusUseProxy          = Status{Code: 305, Label: "Use Proxy"}          // RFC 9110, 15.4.6
	StatusTemporaryRedirect = Status{Code: 307, Label: "Temporary Redirect"} // RFC 9110, 15.4.8
	StatusPermanentRedirect = Status{Code: 308, Label: "Permanent Redirect"} // RFC 9110, 15.4.9

	StatusBadRequest                   = Status{Code: 400, Label: "Bad Request"}                     // RFC 9110, 15.5.1
	StatusUnauthorized                 = Status{Code: 401, Label: "Unauthorized"}                    // RFC 9110, 15.5.2
	StatusPaymentRequired              = Status{Code: 402, Label: "Payment Required"}                // RFC 9110, 15.5.3
	StatusForbidden                    = Status{Code: 403, Label: "Forbidden"}                       // RFC 9110, 15.5.4
	StatusNotFound                     = Status{Code: 404, Label: "Not Found"}                       // RFC 9110, 15.5.5
	StatusMethodNotAllowed             = Status{Code: 405, Label: "Method Not Allowed"}              // RFC 9110, 15.5.6
	StatusNotAcceptable                = Status{Code: 406, Label: "Not Acceptable"}                  // RFC 9110, 15.5.7
	StatusProxyAuthRequired            = Status{Code: 407, Label: "Proxy Authentication Required"}   // RFC 9110, 15.5.8
	StatusRequestTimeout               = Status{Code: 408, Label: "Request Timeout"}                 // RFC 9110, 15.5.9
	StatusConflict                     = Status{Code: 409, Label: "Conflict"}                        // RFC 9110, 15.5.10
	StatusGone                         = Status{Code: 410, Label: "Gone"}                            // RFC 9110, 15.5.11
	StatusLengthRequired               = Status{Code: 411, Label: "Length Required"}                 // RFC 9110, 15.5.12
	StatusPreconditionFailed           = Status{Code: 412, Label: "Precondition Failed"}             // RFC 9110, 15.5.13
	StatusRequestEntityTooLarge        = Status{Code: 413, Label: "Payload Too Large"}               // RFC 9110, 15.5.14
	StatusRequestURITooLong            = Status{Code: 414, Label: "URI Too Long"}                    // RFC 9110, 15.5.15
	StatusUnsupportedMediaType         = Status{Code: 415, Label: "Unsupported Media Type"}          // RFC 9110, 15.5.16
	StatusRequestedRangeNotSatisfiable = Status{Code: 416, Label: "Range Not Satisfiable"}           // RFC 9110, 15.5.17
	StatusExpectationFailed            = Status{Code: 417, Label: "Expectation Failed"}              // RFC 9110, 15.5.18
	StatusTeapot                       = Status{Code: 418, Label: "I'm a teapot"}                    // RFC 9110, 15.5.19
	StatusMisdirectedRequest           = Status{Code: 421, Label: "Misdirected Request"}             // RFC 9110, 15.5.20
	StatusUnprocessableEntity          = Status{Code: 422, Label: "Unprocessable Entity"}            // RFC 9110, 15.5.21
	StatusLocked                       = Status{Code: 423, Label: "Locked"}                          // RFC 4918, 11.3
	StatusFailedDependency             = Status{Code: 424, Label: "Failed Dependency"}               // RFC 4918, 11.4
	StatusTooEarly                     = Status{Code: 425, Label: "Too Early"}                       // RFC 8470, 5.2.
	StatusUpgradeRequired              = Status{Code: 426, Label: "Upgrade Required"}                // RFC 9110, 15.5.22
	StatusPreconditionRequired         = Status{Code: 428, Label: "Precondition Required"}           // RFC 6585, 3
	StatusTooManyRequests              = Status{Code: 429, Label: "Too Many Requests"}               // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  = Status{Code: 431, Label: "Request Header Fields Too Large"} // RFC 6585, 5
	StatusUnavailableForLegalReasons   = Status{Code: 451, Label: "Unavailable For Legal Reasons"}   // RFC 7725, 3

	StatusInternalServerError           = Status{Code: 500, Label: "Internal Server Error"}           // RFC 9110, 15.6.1
	StatusNotImplemented                = Status{Code: 501, Label: "Not Implemented"}                 // RFC 9110, 15.6.2
	StatusBadGateway                    = Status{Code: 502, Label: "Bad Gateway"}                     // RFC 9110, 15.6.3
	StatusServiceUnavailable            = Status{Code: 503, Label: "Service Unavailable"}             // RFC 9110, 15.6.4
	StatusGatewayTimeout                = Status{Code: 504, Label: "Gateway Timeout"}                 // RFC 9110, 15.6.5
	StatusHTTPVersionNotSupported       = Status{Code: 505, Label: "HTTP Version Not Supported"}      // RFC 9110, 15.6.6
	StatusVariantAlsoNegotiates         = Status{Code: 506, Label: "Variant Also Negotiates"}         // RFC 2295, 8.1
	StatusInsufficientStorage           = Status{Code: 507, Label: "Insufficient Storage"}            // RFC 4918, 11.5
	StatusLoopDetected                  = Status{Code: 508, Label: "Loop Detected"}                   // RFC 5842, 7.2
	StatusNotExtended                   = Status{Code: 510, Label: "Not Extended"}                    // RFC 2774, 7
	StatusNetworkAuthenticationRequired = Status{Code: 511, Label: "Network Authentication Required"} // RFC 6585, 6

	statusMap = map[int]*Status{
		100: &StatusContinue,
		101: &StatusSwitchingProtocols,
		102: &StatusProcessing,
		103: &StatusEarlyHints,
		200: &StatusOK,
		201: &StatusCreated,
		202: &StatusAccepted,
		203: &StatusNonAuthoritativeInfo,
		204: &StatusNoContent,
		205: &StatusResetContent,
		206: &StatusPartialContent,
		207: &StatusMultiStatus,
		208: &StatusAlreadyReported,
		226: &StatusIMUsed,
		300: &StatusMultipleChoices,
		301: &StatusMovedPermanently,
		302: &StatusFound,
		303: &StatusSeeOther,
		304: &StatusNotModified,
		305: &StatusUseProxy,
		307: &StatusTemporaryRedirect,
		308: &StatusPermanentRedirect,
		400: &StatusBadRequest,
		401: &StatusUnauthorized,
		402: &StatusPaymentRequired,
		403: &StatusForbidden,
		404: &StatusNotFound,
		405: &StatusMethodNotAllowed,
		406: &StatusNotAcceptable,
		407: &StatusProxyAuthRequired,
		408: &StatusRequestTimeout,
		409: &StatusConflict,
		410: &StatusGone,
		411: &StatusLengthRequired,
		412: &StatusPreconditionFailed,
		413: &StatusRequestEntityTooLarge,
		414: &StatusRequestURITooLong,
		415: &StatusUnsupportedMediaType,
		416: &StatusRequestedRangeNotSatisfiable,
		417: &StatusExpectationFailed,
		418: &StatusTeapot,
		421: &StatusMisdirectedRequest,
		422: &StatusUnprocessableEntity,
		423: &StatusLocked,
		424: &StatusFailedDependency,
		425: &StatusTooEarly,
		426: &StatusUpgradeRequired,
		428: &StatusPreconditionRequired,
		429: &StatusTooManyRequests,
		431: &StatusRequestHeaderFieldsTooLarge,
		451: &StatusUnavailableForLegalReasons,
		500: &StatusInternalServerError,
		501: &StatusNotImplemented,
		502: &StatusBadGateway,
		503: &StatusServiceUnavailable,
		504: &StatusGatewayTimeout,
		505: &StatusHTTPVersionNotSupported,
		506: &StatusVariantAlsoNegotiates,
		507: &StatusInsufficientStorage,
		508: &StatusLoopDetected,
		510: &StatusNotExtended,
		511: &StatusNetworkAuthenticationRequired,
	}
)

func GetStatus(code int) *Status {
	if status, ok := statusMap[code]; ok {
		return status
	}
	return nil
}

func (s *Status) Is5xx() bool {
	return s.Code >= 500 && s.Code < 600
}

func (s *Status) Is4xx() bool {
	return s.Code >= 400 && s.Code < 500
}

func (s *Status) Is3xx() bool {
	return s.Code >= 300 && s.Code < 400
}

func (s *Status) Is2xx() bool {
	return s.Code >= 200 && s.Code < 300
}

func (s *Status) Is1xx() bool {
	return s.Code >= 100 && s.Code < 200
}

func (s *Status) IsErr() bool {
	return s.Is4xx() || s.Is5xx()
}

func (s *Status) IsSuccess() bool {
	return s.Is2xx()
}
