package web

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

const (
	ctxSyncDataType = iota
	ctxTimestamp
	ctxNodeIp
	ctxChartType
)

func syncDataType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxSyncDataType,
			chi.URLParam(r, "dataType"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getSyncTypeCtx retrieves the syncDataType data from the request context.
// If not set, the return value is an empty string.
func getSyncDataTypeCtx(r *http.Request) string {
	syncType, ok := r.Context().Value(ctxSyncDataType).(string)
	if !ok {
		log.Trace("sync type not set")
		return ""
	}
	return syncType
}

func addTimestampToCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxTimestamp,
			chi.URLParam(r, "timestamp"))
			next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// chartTypeCtx returns a http.HandlerFunc that embeds the value at the url
// part {charttype} into the request context.
func chartTypeCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxChartType,
			chi.URLParam(r, "charttype"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTitmestampCtx(r *http.Request) int64 {
	timestampStr, ok := r.Context().Value(ctxTimestamp).(string)
	if !ok {
		return 0
	}
	timestamp, _ := strconv.ParseInt(timestampStr, 10, 64)
	return timestamp
}

func addNodeIPToCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxNodeIp,
			chi.URLParam(r, "address"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getNodeIPFromCtx(r *http.Request) string {
	address, ok := r.Context().Value(ctxNodeIp).(string)
	if !ok {
		return ""
	}
	return address
}
// getChartTypeCtx retrieves the ctxChart data from the request context.
// If not set, the return value is an empty string.
func getChartTypeCtx(r *http.Request) string {
	chartType, ok := r.Context().Value(ctxChartType).(string)
	if !ok {
		log.Trace("chart type not set")
		return ""
	}
	return chartType
}
