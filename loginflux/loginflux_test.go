package loginflux

import (
	"testing"
)

var is *InfluxServer

func init() {
	var keys = []string{"level", "commont"}
	tags := make(map[string]string, 1)
	tags["server"] = "localhostmac"
	is = &InfluxServer{
		"http://localhost:8086",
		"derek",
		"123456",
		"serverInfo",
		"logs",
		"ms",
		2,
		tags,
		keys,
	}
}

func TestLoginflux_Write(t *testing.T) {
	wr := NewLoginflux(*is)
	args := []string{
		`{"component":"influxdb","file":"logging.go","level":"info","line":63,"method":"jvole.com/influx/route.(*loggingService).Delete","took":1666,"ts":"2018-02-13T09:12:04.421056917Z"}`,
		`{"file":"transport.go","level":"debug","line":304,"method":"jvole.com/influx/route.encodeResponse","response":{"Errcode":"0","Msg":"ok","Data":null},"took":11,"ts":"2018-02-13T09:14:05.310336721Z"}`,
	}
	for _, v := range args {
		wr.Write([]byte(v))
	}
}
