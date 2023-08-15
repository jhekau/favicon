package main

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 20 July 2023
 */
import adapter_http_v1 "github.com/jhekau/favicon/internal/adapters/http/v1"

func main(){

	switch {
	case adapter_http_v1.Run(): return
	}
}
