// Package http — adapter TIPIS di atas tracing.Inject/Extract, khusus untuk net/http.
// Opsional: pakai ini kalau service kamu REST/HTTP. Tidak menambah dependency baru
// selain net/http standar (tidak butuh otelhttp).
package http

import (
	"context"
	"net/http"

	"github.com/funxdofficial/golang-module-jaeger/tracing"
)

// Middleware men-extract trace context dari header request masuk dan menaruhnya
// di r.Context() sebelum meneruskan ke handler berikutnya. Pasang di paling luar
// (sebelum router/handler lain).
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/approval", handler)
//	http.ListenAndServe(":8080", httptransport.Middleware(mux))
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := map[string]string{}
		for k := range r.Header {
			headers[k] = r.Header.Get(k)
		}
		ctx := tracing.Extract(r.Context(), headers)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// InjectHeaders menyalin trace context dari ctx ke header request KELUAR (client call).
// Panggil sebelum client.Do(req).
//
//	req, _ := http.NewRequestWithContext(span.Context(), "POST", url, body)
//	httptransport.InjectHeaders(span.Context(), req.Header)
//	client.Do(req)
func InjectHeaders(ctx context.Context, header http.Header) {
	for k, v := range tracing.Inject(ctx) {
		header.Set(k, v)
	}
}
