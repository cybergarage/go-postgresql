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

package protocol

// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html

// MessageType represents a message type.
type MessageType byte

// Frontend (F).
const (
	BindMessage                MessageType = 'B'
	CancelRequestMessage       MessageType = ' ' // Int32(16)
	CloseMessage               MessageType = 'C'
	CopyFailMessage            MessageType = 'f'
	DescribeMessage            MessageType = 'D'
	ExecuteMessage             MessageType = 'E'
	FlushMessage               MessageType = 'H'
	FunctionCallMessage        MessageType = 'F'
	GSSENCRequestMessage       MessageType = ' ' // Int32(8)
	GSSResponseMessage         MessageType = 'p'
	ParseMessage               MessageType = 'P'
	PasswordMessage            MessageType = 'p'
	QueryMessage               MessageType = 'F'
	SASLInitialResponseMessage MessageType = 'p'
	SASLResponseMessage        MessageType = 'p'
	SSLRequestMessage          MessageType = ' ' // Int32(8)
	StartupMessage             MessageType = ' ' // Int32
	SyncMessage                MessageType = 'S'
	TerminateMessage           MessageType = 'F'
)

// Backend (B).
const (
	AuthenticationOkMessage                MessageType = 'R'
	AuthenticationKerberosV5Message        MessageType = 'R'
	AuthenticationCleartextPasswordMessage MessageType = 'R'
	AuthenticationMD5PasswordMessage       MessageType = 'R'
	AuthenticationGSSMessage               MessageType = 'R'
	AuthenticationSSPIMessage              MessageType = 'R'
	AuthenticationSASLMessage              MessageType = 'R'
	AuthenticationSASLContinueMessage      MessageType = 'R'
	AuthenticationSASLFinalMessage         MessageType = 'R'
	BackendKeyDataMessage                  MessageType = 'K'
	BindCompleteMessage                    MessageType = '2'
	CloseCompleteMessage                   MessageType = '3'
	CommandCompleteMessage                 MessageType = 'C'
	CopyInResponseMessage                  MessageType = 'G'
	CopyOutResponseMessage                 MessageType = 'H'
	CopyBothResponseMessage                MessageType = 'W'
	DataRowMessage                         MessageType = 'D'
	EmptyQueryResponseMessage              MessageType = 'I'
	ErrorResponseMessage                   MessageType = 'E'
	FunctionCallResponseMessage            MessageType = 'V'
	NegotiateProtocolVersionMessage        MessageType = 'v'
	NoDataMessage                          MessageType = 'n'
	NoticeResponseMessage                  MessageType = 'N'
	NotificationResponseMessage            MessageType = 'A'
	ParameterDescriptionMessage            MessageType = 't'
	ParameterStatusMessage                 MessageType = 'S'
	ParseCompleteMessage                   MessageType = 'B'
	PortalSuspendedMessage                 MessageType = 's'
	ReadyForQueryMessage                   MessageType = 'Z'
	RowDescriptionMessage                  MessageType = 'T'
)

// Both (F & B).
const (
	CopyDataMessage MessageType = 'd'
	CopyDoneMessage MessageType = 'c'
)
