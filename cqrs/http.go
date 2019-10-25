package cqrs

import "net/http"

type BufferedResponseWriter struct {
	headers http.Header
	statusCode int
	buf []byte
	w http.ResponseWriter
}

func NewBufferedResponseWriter(w http.ResponseWriter)*BufferedResponseWriter {
	return &BufferedResponseWriter{
		headers:    make(http.Header),
		statusCode: 200,
		w:w,
	}
}

func (b *BufferedResponseWriter) Header() http.Header {
	return b.headers
}

func (b *BufferedResponseWriter) Write(buf []byte) (int, error) {
	b.buf = append(b.buf,buf...)
	return len(buf),nil
}

func (b *BufferedResponseWriter) WriteHeader(statusCode int) {
	b.statusCode = statusCode
}

func (b *BufferedResponseWriter) Close() error {
	b.w.WriteHeader(b.statusCode)
	h := b.w.Header()
	for k, v := range b.headers {
		if len(v) == 1 {
			h.Set(k,v[0])
		}
		for _, vv := range v {
			h.Add(k,vv)
		}
	}
	_,err := b.w.Write(b.buf)
	return err
}

