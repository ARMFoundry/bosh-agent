// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/cloudfoundry/bosh-agent/platform/windows/disk"
)

type FakeWindowsDiskFormatter struct {
	FormatStub        func(diskNumber, partitionNumber string) error
	formatMutex       sync.RWMutex
	formatArgsForCall []struct {
		diskNumber      string
		partitionNumber string
	}
	formatReturns struct {
		result1 error
	}
	formatReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeWindowsDiskFormatter) Format(diskNumber string, partitionNumber string) error {
	fake.formatMutex.Lock()
	ret, specificReturn := fake.formatReturnsOnCall[len(fake.formatArgsForCall)]
	fake.formatArgsForCall = append(fake.formatArgsForCall, struct {
		diskNumber      string
		partitionNumber string
	}{diskNumber, partitionNumber})
	fake.recordInvocation("Format", []interface{}{diskNumber, partitionNumber})
	fake.formatMutex.Unlock()
	if fake.FormatStub != nil {
		return fake.FormatStub(diskNumber, partitionNumber)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.formatReturns.result1
}

func (fake *FakeWindowsDiskFormatter) FormatCallCount() int {
	fake.formatMutex.RLock()
	defer fake.formatMutex.RUnlock()
	return len(fake.formatArgsForCall)
}

func (fake *FakeWindowsDiskFormatter) FormatArgsForCall(i int) (string, string) {
	fake.formatMutex.RLock()
	defer fake.formatMutex.RUnlock()
	return fake.formatArgsForCall[i].diskNumber, fake.formatArgsForCall[i].partitionNumber
}

func (fake *FakeWindowsDiskFormatter) FormatReturns(result1 error) {
	fake.FormatStub = nil
	fake.formatReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeWindowsDiskFormatter) FormatReturnsOnCall(i int, result1 error) {
	fake.FormatStub = nil
	if fake.formatReturnsOnCall == nil {
		fake.formatReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.formatReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeWindowsDiskFormatter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.formatMutex.RLock()
	defer fake.formatMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeWindowsDiskFormatter) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ disk.WindowsDiskFormatter = new(FakeWindowsDiskFormatter)
