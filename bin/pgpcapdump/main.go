// Copyright (C) 2019 The go-postgresql Authors. All rights reserved.
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

func outputf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
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

	reqWriter := &bytes.Buffer{}
	resWriter := &bytes.Buffer{}

	pktSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for pkt := range pktSource.Packets() {
		fmt.Print(pkt)
		// continue
		tcpLayer := pkt.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}
		tcp, _ := tcpLayer.(*layers.TCP)
		if tcp.DstPort != postgresql.DefaultPort && tcp.SrcPort != postgresql.DefaultPort {
			continue
		}
		appLayer := pkt.ApplicationLayer()
		if appLayer == nil {
			continue
		}
		payload := appLayer.Payload()
		if tcp.DstPort == postgresql.DefaultPort {
			reqWriter.Write(payload)
		} else if tcp.SrcPort == postgresql.DefaultPort {
			resWriter.Write(payload)
		}
		fmt.Printf("Payload = [%d] %s\n", len(appLayer.Payload()), hex.EncodeToString(appLayer.Payload()))
	}

	fmt.Printf("reqWriter = %d %s\n", reqWriter.Len(), hex.EncodeToString(reqWriter.Bytes()[:128]))
	fmt.Printf("resWriter = %d %s\n", resWriter.Len(), hex.EncodeToString(resWriter.Bytes()[:128]))

	reqMsgReader := message.NewMessageReaderWith(bufio.NewReader(bytes.NewReader(reqWriter.Bytes())))

	// Handle a Start-up message.

	_, err = message.NewStartupWithReader(reqMsgReader)
	if err != nil {
		exit(err)
	}

	outputf("startup")

	for {
		msgType, err := reqMsgReader.PeekType()
		if err != nil {
			break
		}

		outputf("%s (%s)", msgType.String(), hex.EncodeToString([]byte{byte(msgType)}))

		switch msgType { // nolint:exhaustive
		case message.QueryMessage:
			if *isQueryEnabled {
				query, err := message.NewQueryWithReader(reqMsgReader)
				if err != nil {
					continue
					// exit(err)
				}
				println(query.Query)
			} else {
				err := skipMessage(reqMsgReader)
				if err != nil {
					exit(err)
				}
			}
		default:
			err := skipMessage(reqMsgReader)
			if err != nil {
				exit(err)
			}
		}
	}

	os.Exit(0)
}
