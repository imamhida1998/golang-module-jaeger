// Package tracing — propagation.go: fungsi generik untuk membawa trace context
// lintas proses/service, TIDAK terikat ke transport tertentu (gRPC, HTTP, queue, dll).
//
// Pemanggil: Inject(ctx) -> map[string]string, taruh isinya ke header/metadata/field
// apa pun yang dipakai transport kamu.
//
// Penerima: Extract(ctx, receivedHeaders) -> context.Context baru; pakai ctx ini untuk
// Operation() pertama di sisi penerima, supaya otomatis jadi child dari span pemanggil.
package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Inject mengambil trace context aktif dari ctx dan menuliskannya ke sebuah carrier
// map[string]string. Carrier ini generik — bisa dikonversi ke http.Header, grpc
// metadata.MD, atau field custom (misal header pesan Kafka/RabbitMQ).
//
//	headers := tracing.Inject(span.Context())
//	// http: for k, v := range headers { req.Header.Set(k, v) }
//	// grpc: md := metadata.New(headers); ctx = metadata.NewOutgoingContext(ctx, md)
func Inject(ctx context.Context) map[string]string {
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return map[string]string(carrier)
}

// Extract membaca trace context dari carrier (header/metadata yang diterima) dan
// mengembalikan context.Context baru yang membawanya. Panggil ini SEBELUM
// Operation() pertama di sisi penerima.
//
//	ctx = tracing.Extract(ctx, receivedHTTPHeadersAsMap)
//	span := tracing.GetTracing().Operation(ctx, ...) // otomatis jadi child span
func Extract(ctx context.Context, carrier map[string]string) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(carrier))
}
