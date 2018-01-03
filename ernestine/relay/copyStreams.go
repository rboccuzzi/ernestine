package relay

import (
  "io"
)

const bufferSize = 1024

func CopyStreams(src io.Reader, dst io.Writer) {
  buff := make([]byte, bufferSize)

  for {
    n, err := src.Read(buff)
    if err != nil {
      if n == 0 && err == io.EOF {
        return
      }
      panic(err)
    }
    dst.Write(buff[:n])
  }
}
