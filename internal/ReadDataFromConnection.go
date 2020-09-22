package internal

import (
	"context"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"net"
	"time"
)

func ReadDataFromConnection(
	Conn net.Conn,
	CancelFunc CancelFunc,
	CancelCtx context.Context,
	ConnectionManager rxgo.IPublishToConnectionManager,
	ConnectionId string,
	index int,
	description string,
	next func(rws goprotoextra.IReadWriterSize, cancelCtx context.Context, CancelFunc CancelFunc)) {
	messageCount := 0
	byteCount := 0
	var buffer []byte
	bufferStart := 0
	bufferEnd := 4096
	resetBuffer := func() {
		bufferStart = 0
		bufferEnd = 4096
		buffer = make([]byte, bufferEnd)
	}
	resetBuffer()
	var lastUpdate time.Time
	for {
		if CancelCtx.Err() == nil{
			now := time.Now()
			if now.Sub(lastUpdate) >= time.Second {
				lastUpdate = now
				ConnectionManager.PublishStackData(
					index,
					ConnectionId,
					description,
					rxgo.StreamDirectionInbound,
					messageCount,
					byteCount)
			}
		}
		n, err := Conn.Read(buffer[bufferStart:bufferEnd])
		if err != nil {
			switch v := err.(type) {
			case *net.OpError:
				CancelFunc(description, true, v.Err)
			default:
				CancelFunc(description, true, err)
			}

			return
		}
		if CancelCtx.Err() != nil {
			return
		}
		messageCount++
		byteCount += n
		next(gomessageblock.NewReaderWriterBlock(buffer[bufferStart:bufferStart+n]), CancelCtx, CancelFunc)
		if CancelCtx.Err() != nil {
			return
		}
		bufferStart += n
		if bufferStart >= bufferEnd-255 {
			resetBuffer()
		}
	}
}
