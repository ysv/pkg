package httputil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Similar to net/http/httputil.DumpResponse but prettified and includes more details.
// TODO: Add ability to pass options to DumpResponse which define what to include in dump.
func DumpResponse(r *http.Response) (dump []byte, err error) {
	req := r.Request

	// Ignore panics if they happen.
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("failed to dump response")
		}
	}()

	resBodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewReader(resBodyBytes))

	var reqBodyBytes []byte
	if req.Body != nil {
		reqBody, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		reqBodyBytes, err = ioutil.ReadAll(reqBody)
		if err != nil {
			return nil, err
		}
	}

	parts := []string{
		"-- Response Status --",
		fmt.Sprintf("%v %v --", r.Proto, r.Status),
		"",
		"-- Request URL --",
		req.URL.String(),
		"",
		"-- Request Method --",
		req.Method,
		"",
		"-- Request Headers --",
		string(dumpHeaders(req.Header)),
		"",
		"-- Request Body --",
		string(dumpBody(reqBodyBytes, req.Header)),
		"",
		"-- Response Headers --",
		string(dumpHeaders(r.Header)),
		"",
		"-- Response Body --",
		string(dumpBody(resBodyBytes, r.Header)),
	}
	dump = []byte(strings.Join(parts, "\n"))

	return
}

func dumpHeaders(originalHeaders http.Header) []byte {
	headers := make(map[string]interface{}, len(originalHeaders))
	for k, origHeader := range originalHeaders {
		if len(origHeader) == 1 {
			headers[k] = origHeader[0]
		} else {
			headers[k] = origHeader
		}
	}

	headersDump, err := json.MarshalIndent(headers, "", "  ")
	if err != nil {
		return []byte("FAILED TO DUMP")
	}

	return headersDump
}

func dumpBody(originalBody []byte, headers http.Header) []byte {
	var body []byte

	if originalBody == nil {
		return body
	}

	if headers.Get("Content-Type") == "application/json" {
		var decodedBody map[string]interface{}
		if err := json.Unmarshal(originalBody, &decodedBody); err == nil {
			body, _ = json.MarshalIndent(decodedBody, "", "  ")
		}
	}

	if body == nil {
		body = originalBody
	}

	// TODO: 1. Add truncation. 2. Ignore binary body.
	return body
}

//"\n" +
//"-- 301 MOVED PERMANENTLY --\n" +
//"\n" +
//"-- Request URL --\n" +
//"https://google.com\n" +
//"\n" +
//"-- Request method --\n" +
//"GET\n" +
//"\n" +
//"-- Request headers --\n" +
//"{\"User-Agent\":\"Faraday v1.4.1\"}\n" +
//"\n" +
//"-- Request body --\n" +
//"\n" +
//"\n" +
//"-- Request sent at --\n" +
//"2021-05-07 14:08:59.75 UTC\n" +
//"\n" +
//"-- Response headers --\n" +
//"{\"location\":\"https://www.google.com/\",\"content-type\":\"text/html; charset=UTF-8\",\"date\":\"Fri, 07 May 2021 14:08:59 GMT\",\"expires\":\"Sun, 06 Jun 2021 14:08:59 GMT\",\"cache-control\":\"public, max-age=2592000\",\"server\":\"gws\",\"content-length\":\"220\",\"x-xss-protection\":\"0\",\"x-frame-options\":\"SAMEORIGIN\"}\n" +
//"\n" +
//"-- Response body --\n" +
//"<HTML><HEAD><meta http-equiv=\"content-type\" content=\"text/html;charset=utf-8\">\n" +
//"<TITLE>301 Moved</TITLE></HEAD><BODY>\n" +
//"<H1>301 Moved</H1>\n" +
//"The document has moved\n" +
//"<A HREF=\"https://www.google.com/\">here</A>.\r\n" +
//"</BODY></HTML>\r\n" +
//"\n" +
//"\n" +
//"-- Response received at --\n" +
//"2021-05-07 14:08:59.80 UTC\n" +
//"\n" +
//"-- Response received in --\n" +
//"50.18ms\n"
//lines = [
//"",
//"-- #{status} #{reason_phrase} --".upcase,
//"",
//"-- Request URL --",
//env.url.to_s,
//"",
//"-- Request method --",
//env.method.to_s.upcase,
//"",
//"-- Request headers --",
//::JSON.generate(request_headers).yield_self { |t| t.truncate(2048, omission: "... (truncated, full length: #{t.length})") },
//"",
//
//"-- Request body --",
//if request_json
//::JSON.generate(request_json)
//else
//body = env.request_body.to_s.dup
//if body.encoding.name == "ASCII-8BIT"
//"Binary (#{body.size} bytes)"
//else
//body
//end
//end.yield_self { |t| t.truncate(1024, omission: "... (truncated, full length: #{t.length})") },
//"",
//
//"-- Request sent at --",
//env.request_sent_at.strftime("%Y-%m-%d %H:%M:%S.%2N") + " UTC",
//"",
//
//"-- Response headers --",
//if response_headers
//::JSON.generate(response_headers)
//else
//env.response_headers.to_s
//end.yield_self { |t| t.truncate(2048, omission: "... (truncated, full length: #{t.length})") },
//"",
//
//"-- Response body --",
//if response_json
//::JSON.generate(response_json)
//else
//body = env.body.to_s.dup
//if body.encoding.name == "ASCII-8BIT"
//"Binary (#{body.size} bytes)"
//else
//body
//end
//end.yield_self { |t| t.truncate(2048, omission: "... (truncated, full length: #{t.length})") }
//]
//
//if env.response_received_at
//lines.concat [
//"",
//"-- Response received at --",
//env.response_received_at.strftime("%Y-%m-%d %H:%M:%S.%2N") + " UTC",
//"",
//"-- Response received in --",
//"#{((env.response_received_at.to_f - env.request_sent_at.to_f) * 1000.0).round(2)}ms"
//]
//end
