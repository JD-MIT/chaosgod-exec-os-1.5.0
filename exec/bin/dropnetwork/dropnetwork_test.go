/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"testing"

	"github.com/chaosblade-io/chaosblade-spec-go/channel"
	"github.com/chaosblade-io/chaosblade-spec-go/spec"

	"github.com/chaosblade-io/chaosblade-exec-os/exec/bin"
)

func Test_startDropNet_failed(t *testing.T) {
	var exitCode int
	bin.ExitFunc = func(code int) {
		exitCode = code
	}
	tests := []struct {
		sourceIp        string
		destinationIp   string
		sourcePort      string
		destinationPort string
		stringPattern   string
		networkTraffic  string
	}{
		{"", "", "", "", "", ""},
	}

	for _, tt := range tests {
		startDropNet(tt.sourceIp, tt.destinationIp, tt.sourcePort, tt.destinationPort, tt.stringPattern, tt.networkTraffic)
		if exitCode != 1 {
			t.Errorf("unexpected result: %d, expected result: %d", exitCode, 1)
		}
	}
}

func Test_handleDropSpecifyPort(t *testing.T) {
	type input struct {
		sourceIp        string
		destinationIp   string
		sourcePort      string
		destinationPort string
		stringPattern   string
		networkTraffic  string
		response        *spec.Response
	}
	type expect struct {
		exitCode   int
		invokeTime int
	}

	tests := []struct {
		input  input
		expect expect
	}{
		{input{"", "", "80", "", "", "", spec.ResponseFailWithFlags(spec.CommandIptablesNotFound)},
			expect{1, 1}},
		{input{"", "", "", "80", "", "", spec.ResponseFailWithFlags(spec.CommandIptablesNotFound)},
			expect{1, 1}},
		{input{"", "", "80", "", "", "", spec.ReturnSuccess("success")},
			expect{0, 0}},
	}

	var exitCode int
	bin.ExitFunc = func(code int) {
		exitCode = code
	}
	var invokeTime int
	stopDropNetFunc = func(sourceIp, destinationIp, sourcePort, destinationPort, stringPattern, networkTraffic string) {
		invokeTime++
	}
	for _, tt := range tests {
		cl = channel.NewMockLocalChannel()
		mockChannel := cl.(*channel.MockLocalChannel)
		mockChannel.RunFunc = func(ctx context.Context, script, args string) *spec.Response {
			return tt.input.response
		}
		handleDropSpecifyPort(tt.input.sourceIp, tt.input.destinationIp, tt.input.sourcePort, tt.input.destinationPort, tt.input.stringPattern, tt.input.networkTraffic, context.Background())
		if exitCode != tt.expect.exitCode {
			t.Errorf("unexpected result: %d, expected result: %d", exitCode, tt.expect.exitCode)
		}
	}
}
