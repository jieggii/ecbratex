package tests

//import (
//	"fmt"
//	"net/http"
//	"net/http/httptest"
//)
//
//type InvalidTestHTTPServer struct {
//	*httptest.Server
//
//	URLLatest               string
//	URLTimeSeries           string
//	URLTimeSeriesLast90Days string
//
//	urlPathLatest               string
//	urlPathTimeSeries           string
//	urlPathTimeSeriesLast90Days string
//}
//
//func (p *InvalidTestHTTPServer) HandlerFunc() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusInternalServerError)
//		if _, err := w.Write([]byte("")); err != nil {
//			panic(err)
//		}
//	}
//}
//
//func NewInvalidTestHTTPServer() *InvalidTestHTTPServer {
//	server := &InvalidTestHTTPServer{
//
//		urlPathLatest:               "/latest/",
//		urlPathTimeSeries:           "/time-series/",
//		urlPathTimeSeriesLast90Days: "/time-series-last-90-days/",
//	}
//	server.Server = httptest.NewServer(server.HandlerFunc())
//	server.Server = httptest.NewServer(server.HandlerFunc())
//	server.URLLatest = fmt.Sprintf("%s%s", server.URL, server.urlPathLatest)
//	server.URLTimeSeries = fmt.Sprintf("%s%s", server.URL, server.urlPathTimeSeries)
//	server.URLTimeSeriesLast90Days = fmt.Sprintf("%s%s", server.URL, server.urlPathTimeSeriesLast90Days)
//
//	return server
//}
