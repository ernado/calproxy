package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var arg struct {
		listen  string
		output  string
		target  string
		deleted string
		dump    bool
		secure  bool
		cert    string
		key     string
	}
	flag.StringVar(&arg.listen, "listen", "localhost:8080", "listen address")
	flag.StringVar(&arg.output, "output", "output", "directory to store requests/responses")
	flag.StringVar(&arg.target, "target", "https://calendar.mail.ru", "target URL")
	flag.BoolVar(&arg.secure, "secure", false, "use https listener")
	flag.StringVar(&arg.cert, "cert", "server.crt", "path to certificate")
	flag.StringVar(&arg.key, "key", "server.key", "path to key")
	flag.StringVar(&arg.deleted, "deleted", "deleted.txt", "file with deleted UIDs")
	flag.BoolVar(&arg.dump, "dump", false, "dump requests/responses")

	flag.Parse()

	deleted := make(map[string]struct{})
	{
		buf, err := os.ReadFile(arg.deleted)
		if err != nil {
			panic(err)
		}
		s := bufio.NewScanner(bytes.NewReader(buf))
		for s.Scan() {
			text := s.Text()
			if text == "" || strings.HasPrefix(text, "#") {
				continue
			}
			text = strings.TrimSpace(text)
			deleted[text] = struct{}{}
		}
	}

	proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf(" -> %s %s", r.Method, r.URL)
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Printf("Error dumping request: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		u, err := url.Parse(arg.target)
		if err != nil {
			log.Printf("Error parsing target URL: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		target := u.ResolveReference(r.URL)
		r.URL = target

		dumpFilePrefix := filepath.Join(
			"output",
			fmt.Sprintf("%s_%s", time.Now().Format("20060102_150405"), r.Method),
		)
		if arg.dump {
			if err := os.WriteFile(dumpFilePrefix+".req.txt", dump, 0644); err != nil {
				log.Printf("Error writing request dump: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		start := time.Now()
		resp, err := http.DefaultTransport.RoundTrip(r)
		if err != nil {
			log.Printf("Error forwarding request: %v", err)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		fmt.Printf(" <- %s %s\n", resp.Status, time.Since(start).Round(time.Millisecond))
		dumpResp, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Printf("Error dumping response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if arg.dump {
			if err := os.WriteFile(dumpFilePrefix+".res.txt", dumpResp, 0644); err != nil {
				log.Printf("Error writing request dump: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var out bytes.Buffer

		if bytes.Contains(body, []byte("ns0:multistatus")) {
			fmt.Println(" -# Found multi-status response")
			out.WriteString(`<?xml version="1.0" encoding="utf-8" ?>`)
			status, err := DecodeMultiStatus(body)
			if err != nil {
				log.Printf("Error decoding MultiStatus: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var filtered []StatusResponse
			for _, s := range status.Responses {
				u, err := url.Parse(s.URI)
				if err != nil {
					log.Printf("Error parsing URI: %v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				base := path.Base(u.Path)
				base = strings.TrimSuffix(base, ".ics")
				if _, ok := deleted[base]; ok {
					fmt.Printf(" -#  Deleted: %s\n", base)
					filtered = append(filtered, s)
				}
			}
			if err := xml.NewEncoder(&out).Encode(status); err != nil {
				log.Printf("Error encoding MultiStatus: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			out.Write(body)
		}

		w.Write(out.Bytes())
	})

	proxyServer := &http.Server{
		Addr:    arg.listen,
		Handler: proxy,
	}

	log.Printf("Starting CalDAV proxy on %s (tls=%v), forwarding to %s", arg.listen, arg.secure, arg.target)
	if arg.secure {
		log.Fatal(proxyServer.ListenAndServeTLS(arg.cert, arg.key))
	} else {
		log.Fatal(proxyServer.ListenAndServe())
	}
}
