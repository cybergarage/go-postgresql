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

package message

// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html

// Type represents a message type.
type Type byte

const (
	NoneMessage = 0x00
)

// Frontend (F).
const (
	BindMessage                Type = 'B'
	CancelRequestMessage       Type = ' ' // Int32(16)
	CloseMessage               Type = 'C'
	CopyFailMessage            Type = 'f'
	DescribeMessage            Type = 'D'
	ExecuteMessage             Type = 'E'
	FlushMessage               Type = 'H'
	FunctionCallMessage        Type = 'F'
	GSSENCRequestMessage       Type = ' ' // Int32(8)
	GSSResponseMessage         Type = 'p'
	ParseMessage               Type = 'P'
	PasswordMessage            Type = 'p'
	QueryMessage               Type = 'F'
	SASLInitialResponseMessage Type = 'p'
	SASLResponseMessage        Type = 'p'
	SSLRequestMessage          Type = ' ' // Int32(8)
	StartupMessage             Type = ' ' // Int32
	SyncMessage                Type = 'S'
	TerminateMessage           Type = 'F'
)

// Backend (B).
const (
	AuthenticationOkMessage                Type = 'R'
	AuthenticationKerberosV5Message        Type = 'R'
	AuthenticationCleartextPasswordMessage Type = 'R'
	AuthenticationMD5PasswordMessage       Type = 'R'
	AuthenticationGSSMessage               Type = 'R'
	AuthenticationSSPIMessage              Type = 'R'
	AuthenticationSASLMessage              Type = 'R'
	AuthenticationSASLContinueMessage      Type = 'R'
	AuthenticationSASLFinalMessage         Type = 'R'
	BackendKeyDataMessage                  Type = 'K'
	BindCompleteMessage                    Type = '2'
	CloseCompleteMessage                   Type = '3'
	CommandCompleteMessage                 Type = 'C'
	CopyInResponseMessage                  Type = 'G'
	CopyOutResponseMessage                 Type = 'H'
	CopyBothResponseMessage                Type = 'W'
	DataRowMessage                         Type = 'D'
	EmptyQueryResponseMessage              Type = 'I'
	ErrorResponseMessage                   Type = 'E'
	FunctionCallResponseMessage            Type = 'V'
	NegotiateProtocolVersionMessage        Type = 'v'
	NoDataMessage                          Type = 'n'
	NoticeResponseMessage                  Type = 'N'
	NotificationResponseMessage            Type = 'A'
	ParameterDescriptionMessage            Type = 't'
	ParameterStatusMessage                 Type = 'S'
	ParseCompleteMessage                   Type = '1'
	PortalSuspendedMessage                 Type = 's'
	ReadyForQueryMessage                   Type = 'Z'
	RowDescriptionMessage                  Type = 'T'
)

// Both (F & B).
const (
	CopyDataMessage Type = 'd'
	CopyDoneMessage Type = 'c'
)
