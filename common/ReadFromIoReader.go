package common

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gomessageblock"
	"io"
)

const BufferSize = 4096

func ReadFromIoReader(
	context string,
	reader io.Reader,
	CancelCtx context.Context,
	cancelFunc context.CancelFunc,
	rxNextHandler goCommsDefinitions.IRxNextHandler,
) {

	var buffer []byte
	bufferStart := 0
	bufferEnd := BufferSize
	resetBuffer := func() {
		bufferStart = 0
		bufferEnd = BufferSize
		buffer = make([]byte, bufferEnd)
	}
	resetBuffer()
	for CancelCtx.Err() == nil {
		if CancelCtx.Err() != nil {
			return
		}
		n, err := reader.Read(buffer[bufferStart:bufferEnd])
		if err != nil {
			rxNextHandler.OnError(err)
			rxNextHandler.OnComplete()
			cancelFunc()
			//ConnectionCancelFunc(fmt.Sprintf("ReadFromIoReader(%v)", context), true, err)
			return
		}
		if CancelCtx.Err() != nil {
			return
		}
		rxNextHandler.OnSendData(gomessageblock.NewReaderWriterBlock(buffer[bufferStart : bufferStart+n]))
		bufferStart += n
		if bufferStart >= bufferEnd-255 {
			resetBuffer()
		}
	}
}
