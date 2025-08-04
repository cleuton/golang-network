// Copyright 2023 Cleuton Sampaio.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"encoding/json"
	"log"

	"github.com/sony/sonyflake"
	"github.com/valyala/fasthttp"
)

const VERSION = "0.1.0"

var logger = log.Default()
var flake = sonyflake.NewSonyflake(sonyflake.Settings{})

// Response é a estrutura da resposta JSON
type Response struct {
	Error bool   `json:"error"`
	ID    uint64 `json:"id,omitempty"`
}

// handlerNextID gera o próximo ID e retorna como JSON
func handlerNextID(ctx *fasthttp.RequestCtx) {
	if flake == nil {
		logger.Panicf("FATAL ERROR: Couldn't generate sonyflake.NewSonyflake.\n")
	}

	id, err := flake.NextID()
	if err != nil {
		logger.Printf("ERROR: flake.NextID() failed with %s\n", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		json.NewEncoder(ctx).Encode(Response{Error: true})
		return
	}

	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(Response{Error: false, ID: id})
}

func main() {
	logger.Printf("INFO: ZAPTIDGEN REST - V %s\n - starting... (port: 8888)", VERSION)

	// Roteamento
	mux := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/nextid":
			handlerNextID(ctx)
		default:
			ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		}
	}

	// Inicia o servidor
	logger.Println("INFO: server listening at :8888")
	if err := fasthttp.ListenAndServe(":8888", mux); err != nil {
		logger.Panicf("FATAL ERROR: failed to serve: %v\n", err)
	}
}
