# Copyright (C) 2019 The go-postgresql Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: certgen clean

all: certgen

certgen: clean
	openssl req -new -x509 -days 365 -nodes -text -out server.crt -keyout server.key -subj "/CN=dbhost.yourdomain.com"
	chmod og-rwx server.key
	openssl req -new -nodes -text -out root.csr -keyout root.key -subj "/CN=root.yourdomain.com"
	chmod og-rwx root.key
	openssl x509 -req -in root.csr -text -days 3650 -extfile /etc/ssl/openssl.cnf -signkey root.key -out root.crt
	openssl req -new -nodes -text -out server.csr -keyout server.key -subj "/CN=dbhost.yourdomain.com"
	chmod og-rwx server.key
	openssl x509 -req -in server.csr -text -days 365 -CA root.crt -CAkey root.key -CAcreateserial -out server.crt

clean:
	rm -f server.crt server.csr server.key root.crt root.csr root.key root.srl