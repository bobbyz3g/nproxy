// Copyright (c) 2022 The Author Kaiser925. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("usage: nproxy <localhost:port@remote:port>")
	}

	done := make(chan struct{})

	for _, v := range args {
		localAddr, remoteAddr, ok := strings.Cut(v, "@")
		if !ok {
			log.Fatalf("error proxy argument: %s", v)
		}
		log.Printf("start listen %q proxy to %q", localAddr, remoteAddr)
		go func() {
			proxy(localAddr, remoteAddr)
		}()
	}

	<-done
}

func proxy(localAddr, remoteAddr string) {
	l, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Printf("listen %q failed: %v", localAddr, err)
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("accept connection failed: %v", err)
			return
		}

		go func() {
			defer conn.Close()

			remote, err := net.Dial("tcp", remoteAddr)
			if err != nil {
				log.Printf("dial %q failed: %v", localAddr, err)
				return
			}
			defer remote.Close()

			go func() {
				_, err = io.Copy(remote, conn)
				if err != nil {
					if err != nil {
						log.Printf("copy data from remote to local failed: %v", err)
					}
				}
			}()

			_, err = io.Copy(conn, remote)
			if err != nil {
				log.Printf("copy data from local to remote failed: %v", err)
			}
		}()

	}
}
