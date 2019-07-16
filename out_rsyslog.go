package main

import (
	"C"
	"encoding/json"
	"fmt"
	"log/syslog"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
)

const (
	pluginName = "rsyslog"
)

type Logger struct {
	network string
	addr    string
	tag     string
	writer  *syslog.Writer
}

func (l *Logger) connect() error {
	if l.writer != nil {
		return nil
	}

	writer, err := syslog.Dial(l.network, l.addr, syslog.LOG_WARNING|syslog.LOG_DAEMON, l.tag)
	if err != nil {
		return err
	}
	l.writer = writer

	return nil
}

func (l *Logger) Info(m string) error {
	if err := l.connect(); err != nil {
		return err
	}
	return l.writer.Info(m)
}

func NewLogger(network, addr, tag string) (*Logger, error) {
	logger := Logger{
		network: network,
		addr:    addr,
		tag:     tag,
	}

	return &logger, nil
}

//export FLBPluginRegister
func FLBPluginRegister(ctx unsafe.Pointer) int {
	return output.FLBPluginRegister(ctx, pluginName, pluginName)
}

//export FLBPluginInit
func FLBPluginInit(plugin unsafe.Pointer) int {
	network := output.FLBPluginConfigKey(plugin, "network")
	addr := output.FLBPluginConfigKey(plugin, "addr")
	tag := output.FLBPluginConfigKey(plugin, "tag")

	fmt.Printf("plugin=%s network=%s addr=%s tag=%s\n", pluginName, network, addr, tag)

	logger, err := NewLogger(network, addr, tag)
	if err != nil {
		return output.FLB_ERROR
	}

	output.FLBPluginSetContext(plugin, unsafe.Pointer(logger))

	return output.FLB_OK
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	return output.FLB_OK
}

//export FLBPluginFlushCtx
func FLBPluginFlushCtx(ctx, data unsafe.Pointer, length C.int, tag *C.char) int {
	logger := (*Logger)(output.FLBPluginGetContext(ctx).(unsafe.Pointer))
	dec := output.NewDecoder(data, int(length))

	logger.Out(ctx, data, length, tag)

	for {
		ret, _, record := output.GetRecord(dec)
		if ret != 0 {
			break
		}

		m := map[string]interface{}{}
		for k, v := range record {
			m[fmt.Sprintf("%s", k)] = fmt.Sprintf("%s", v)
		}

		jsonBody, err := json.Marshal(m)
		if err != nil {
			fmt.Printf("plugin=%s err=%+v\n", pluginName, err)
			return output.FLB_ERROR
		}

		logger.Info(string(jsonBody))
	}

	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	return output.FLB_OK
}

func main() {
}
