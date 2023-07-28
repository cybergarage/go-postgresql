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

// PostgreSQL: Documentation: 16: 55.2. Message Flow
// https://www.postgresql.org/docs/16/protocol-flow.html
// PostgreSQL: Documentation: 16: 55.7. Message Formats
// https://www.postgresql.org/docs/16/protocol-message-formats.html
// PostgreSQL: Documentation: 16: 55.8. Error and Notice Message Fields
// https://www.postgresql.org/docs/16/protocol-error-fields.html

import (
	"fmt"
)

// CommandComplete represents a command complete message.
type CommandComplete struct {
	*ResponseMessage
}

// NewCommandComplete returns a new command complete message instance.
func NewCommandComplete() *CommandComplete {
	return &CommandComplete{
		ResponseMessage: NewResponseMessageWith(CommandCompleteMessage),
	}
}

// NewInsertCompleteWith returns a new command complete message for insert query.
func NewCommandCompleteWith(tag string) (*CommandComplete, error) {
	msg := NewCommandComplete()
	return msg, msg.AppendString(tag)
}

// NewInsertCompleteWith returns a new command complete message for insert query.
func NewInsertCompleteWith(n int) (*CommandComplete, error) {
	msg := NewCommandComplete()
	return msg, msg.AppendString(fmt.Sprintf("INSERT 0 %d", n))
}

// NewUpdateCompleteWith returns a new command complete message for update query.
func NewUpdateCompleteWith(n int) (*CommandComplete, error) {
	msg := NewCommandComplete()
	return msg, msg.AppendString(fmt.Sprintf("UPDATE %d", n))
}

// NewSelectCompleteWith returns a new command complete message for select query.
func NewSelectCompleteWith(n int) (*CommandComplete, error) {
	msg := NewCommandComplete()
	return msg, msg.AppendString(fmt.Sprintf("SELECT %d", n))
}

// NewDeleteCompleteWith returns a new command complete message for delete query.
func NewDeleteCompleteWith(n int) (*CommandComplete, error) {
	msg := NewCommandComplete()
	return msg, msg.AppendString(fmt.Sprintf("DELETE %d", n))
}

// NewCopyCompleteWith returns a new command complete message for copy query.
func NewCopyCompleteWith(n int) (*CommandComplete, error) {
	msg := NewCommandComplete()
	return msg, msg.AppendString(fmt.Sprintf("COPY %d", n))
}

// NewCommitComplete returns a new command complete message for commit query.
func NewCommitComplete() (*CommandComplete, error) {
	msg := NewCommandComplete()
	return msg, msg.AppendString("COMMIT")
}
