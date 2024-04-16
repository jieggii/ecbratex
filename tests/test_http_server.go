package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
)

const (
	filenameLatest               = "eurofxref-daily.xml"
	filenameTimeSeries           = "eurofxref-hist.xml"
	filenameTimeSeriesLast90Days = "eurofxref-hist-90d.xml"
)

type TestHTTPServer struct {
	*httptest.Server

	URLLatest               string
	URLTimeSeries           string
	URLTimeSeriesLast90Days string

	urlPathLatest               string
	urlPathTimeSeries           string
	urlPathTimeSeriesLast90Days string

	fsPathLatest               string
	fsPathTimeSeries           string
	fsPathTimeSeriesLast90Days string

	isBroken bool
}

func (p *TestHTTPServer) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if p.isBroken {
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte("dude, I am broken!")); err != nil {
				panic(err)
			}
			return
		}

		var (
			data []byte
			err  error
		)
		switch r.URL.Path {
		case p.urlPathLatest:
			data, err = os.ReadFile(p.fsPathLatest)
		case p.urlPathTimeSeries:
			data, err = os.ReadFile(p.fsPathTimeSeries)
		case p.urlPathTimeSeriesLast90Days:
			data, err = os.ReadFile(p.fsPathTimeSeriesLast90Days)
		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				panic(err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(data); err != nil {
			panic(err)
		}
	}
}

func NewTestHTTPServer(testDataPath string, isBroken bool) *TestHTTPServer {
	server := &TestHTTPServer{
		Server: nil,

		URLLatest:               "",
		URLTimeSeries:           "",
		URLTimeSeriesLast90Days: "",

		urlPathLatest:               "/latest/",
		urlPathTimeSeries:           "/time-series/",
		urlPathTimeSeriesLast90Days: "/time-series-last-90-days/",

		fsPathLatest:               path.Join(testDataPath, filenameLatest),
		fsPathTimeSeries:           path.Join(testDataPath, filenameTimeSeries),
		fsPathTimeSeriesLast90Days: path.Join(testDataPath, filenameTimeSeriesLast90Days),

		isBroken: isBroken,
	}

	server.Server = httptest.NewServer(server.HandlerFunc())
	server.URLLatest = fmt.Sprintf("%s%s", server.URL, server.urlPathLatest)
	server.URLTimeSeries = fmt.Sprintf("%s%s", server.URL, server.urlPathTimeSeries)
	server.URLTimeSeriesLast90Days = fmt.Sprintf("%s%s", server.URL, server.urlPathTimeSeriesLast90Days)
	return server
}
