package main

import (
    "net/http"
    "net/url"
    "path"
    "encoding/json"
	"strconv"
    "bytes"
	"html/template"
    "strings"
	"crypto/md5"
	"encoding/hex"


    "github.com/gorilla/mux"
    "github.com/Sirupsen/logrus"

)

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route


type LoxRouter struct {
    router *mux.Router
    config Configuration
    configFile string

    publicURL *url.URL


   	tmpl       *template.Template

    bundle []byte
    bundleHash string

}


func NewRouter(config Configuration, configFile string, publicURL *url.URL) (*LoxRouter, error) {

	bundle, err := Asset("ui/build/bundle.js")
	if err != nil {
		return nil, err
	}
	bundleHash := md5.Sum(bundle)

	tmpl, err := template.New("index.html").Parse(tmplIndex)
	if err != nil {
		return nil, err
	}

    loxRouter := &LoxRouter {
        router:         mux.NewRouter().StrictSlash(true),
        config:         config,
        configFile:     configFile,

        tmpl:           tmpl,
        bundle:         bundle,
        bundleHash:     hex.EncodeToString(bundleHash[:]),

        publicURL:      publicURL,
    }

    var routes = Routes{
    Route{
        "GetLoxoneConfig",
        "GET",
        "/config/loxone",
        loxRouter.getLoxoneConfig,
    },

    Route {
        "GetInfluxConfig",
        "GET",
        "/config/influx",
        loxRouter.getInfluxConfig,
    },
    Route {
        "GetMetricsConfig",
        "GET",
        "/config/metrics",
        loxRouter.getMetricsConfig,
    },
    Route {
        "Index",
        "GET",
        "/",
        loxRouter.serveIndex,
    },
    Route {
        "Bundle.js",
        "GET",
        "/bundle.js",
        loxRouter.serveBundle,
    },
    }


    loxRouter.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

    for _, route := range routes {
        loxRouter.router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(route.HandlerFunc)
    }

    return loxRouter, nil
}
func (l *LoxRouter) urlFor(p string) string {
	return path.Join(l.publicURL.Path, p)
}

func (l *LoxRouter) urlForAbs(p string) string {
	u := *l.publicURL
	u.Path = l.urlFor(p)
	return u.String()
}

func (l *LoxRouter) serveIndex(w http.ResponseWriter, r *http.Request) {

	buf := &bytes.Buffer{}

    err := l.tmpl.Execute(buf, struct {
		BaseURL    string
		BasePath   string
		BundleURL  string
		BundleHash string
	}{
		l.urlForAbs("/"),
		strings.TrimRight(l.urlFor("/"), "/"),
		l.urlFor("/bundle.js"),
		l.bundleHash,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
    w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))

	if _, err := buf.WriteTo(w); err != nil {
		logrus.Errorf("failed to write:", err)
	}

}
func (l *LoxRouter) serveBundle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("Content-Length", strconv.Itoa(len(l.bundle)))
	if l.bundleHash != "" {
		w.Header().Set("Cache-Control", "max-age=31536000")
	}

	if _, err := w.Write(l.bundle); err != nil {
		logrus.Errorf("failed to write:", err)
	}

}
func (l *LoxRouter) getLoxoneConfig(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(configuration.Loxone); err != nil {
        panic(err)
    }
}
func (l *LoxRouter) getInfluxConfig(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(configuration.InfluxDb); err != nil {
        panic(err)
    }
}
func (l *LoxRouter) getMetricsConfig(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(configuration.Metrics); err != nil {
        panic(err)
    }
}

const tmplIndex = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Loxone Influx Exporter</title>
    <style>
      html { margin: 0; }
      body { margin: 0; }
      a { color: #4078c0; text-decoration: none; }
      a:hover { text-decoration: underline; }
    </style>
  </head>
  <body style="margin: 0;">
    <div id="app"></div>
    <script type="text/javascript">
      BASE_URL = {{ .BaseURL }};
      BASE_PATH = {{ .BasePath }};
    </script>
    <script type="text/javascript" src="{{ .BundleURL }}?hash={{ .BundleHash }}"></script>
  </body>
</html>
`
