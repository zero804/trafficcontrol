package fakesrvr

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/traffic_monitor/tools/testcaches/fakesrvrdata"
)

// TODO config?
const readTimeout = time.Second * 10
const writeTimeout = time.Second * 10

func reqIsApplicationSystem(r *http.Request) bool {
	return r.URL.Query().Get("application") == "system"
}

func astatsHandler(fakeSrvrDataThs fakesrvrdata.Ths) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srvr := (*fakesrvrdata.FakeServerData)(fakeSrvrDataThs.Get())
		// TODO cast to System, if query string `application=system`
		b := []byte{}
		err := error(nil)
		if reqIsApplicationSystem(r) {
			system := srvr.GetSystem()
			b, err = json.MarshalIndent(&system, "", "  ") // TODO debug, change to Marshal
		} else {
			b, err = json.MarshalIndent(&srvr, "", "  ") // TODO debug, change to Marshal
		}
		if err != nil {
			w.Write([]byte(`{"error": "marshalling: ` + err.Error() + `"}`)) // TODO escape error for JSON
		}
		w.Write(b)
	}
}

func Serve(port int, fakeSrvrData fakesrvrdata.Ths) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/_astats", astatsHandler(fakeSrvrData))
	server := &http.Server{
		Addr:           ":" + strconv.Itoa(port),
		Handler:        mux,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			// TODO pass the error somewhere, somehow?
			fmt.Println("Error serving on port " + strconv.Itoa(port) + ": " + err.Error())
		}
	}()
	return server
}
