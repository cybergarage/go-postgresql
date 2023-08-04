// Copyright (C) 2019 Satoshi Konno. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
pgpcapdump is a dump utility for PostgreSQL packet capture file.

	NAME
	 pgpcapdump

	SYNOPSIS
	 pgpcapdump [OPTIONS] FILE

	OPTIONS
	-q      : Extract queries.

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
//nolint:forbidigo
package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/cybergarage/go-postgresql/postgresql"
	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const (
	ProgramName = "pgpcapdump"
)

func usages() {
	println("Usage:")
	println("  " + ProgramName + " [OPTIONS] FILE")
	println("")
	println("Options:")
	println("  -q      : Extract queries.")
	println("")
	println("Return Value:")
	println("  Return EXIT_SUCCESS or EXIT_FAILURE")
	os.Exit(1)
}

func exit(err error) {
	println("Error: " + err.Error())
	os.Exit(1)
}

func main() {
	isQueryEnabled := flag.Bool("q", true, "extract queries")
	flag.Parse()

	args := flag.Args()
	if len(args) <= 0 {
		usages()
	}

	pcapFilename := args[0]
	handle, err := pcap.OpenOffline(pcapFilename)
	if err != nil {
		exit(err)
	}

	skipMessage := func(reader *message.MessageReader) error {
		msg, err := message.NewMessageWithReader(reader)
		if err != nil {
			return err
		}
		_, err = msg.ReadMessageData()
		if err != nil {
			return err
		}
		return nil
	}

	pktSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for pkt := range pktSource.Packets() {
		tcpLayer := pkt.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}
		tcp, _ := tcpLayer.(*layers.TCP)
		if tcp.DstPort != postgresql.DefaultPort && tcp.SrcPort != postgresql.DefaultPort {
			continue
		}
		for _, layer := range pkt.Layers() {
			fmt.Println("PACKET LAYER:", layer.LayerType())
		}
		app := pkt.TransportLayer()
		if app == nil {
			continue
		}
		appPayload := app.LayerPayload()
		println(hex.EncodeToString(appPayload))

		msgReader := message.NewMessageReaderWith(bufio.NewReader(bytes.NewReader(appPayload)))

		for {
			msgType, err := msgReader.PeekType()
			if err != nil {
				break
			}

			switch msgType { // nolint:exhaustive
			case message.QueryMessage:
				if *isQueryEnabled {
					query, err := message.NewQueryWithReader(msgReader)
					if err != nil {
						exit(err)
					}
					println(query.Query)
				} else {
					err := skipMessage(msgReader)
					if err != nil {
						exit(err)
					}
				}
			default:
				err := skipMessage(msgReader)
				if err != nil {
					exit(err)
				}
			}
		}
	}

	os.Exit(0)
}
