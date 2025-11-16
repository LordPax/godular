package core

import (
	"fmt"
	"io"
)

type MockFunc func(...any) any

type MockedMethod struct {
	Method string
	Params []any
	Called bool
	Func   MockFunc
}

type IMockable interface {
	MockMethod(method string, mockFunc MockFunc)
	MethodCalled(method string, params ...any)
	IsMethodCalled(method string) bool
	GetMethodParams(method string) []any
	IsParamsEqual(method string, params ...any) bool
	CallFunc(method string) any
}

type Mockable struct {
	mockedMethod map[string]*MockedMethod
}

func NewMockable() *Mockable {
	return &Mockable{
		mockedMethod: make(map[string]*MockedMethod),
	}
}

func (m *Mockable) MockMethod(method string, mockFunc MockFunc) {
	m.mockedMethod[method] = &MockedMethod{method, nil, false, mockFunc}
}

func (m *Mockable) MethodCalled(method string, params ...any) {
	mock, ok := m.mockedMethod[method]
	if !ok {
		return
	}

	mock.Params = params
	mock.Called = true

	return
}

func (m *Mockable) IsMethodCalled(method string) bool {
	called, ok := m.mockedMethod[method]
	if !ok {
		return false
	}
	return called.Called
}

func (m *Mockable) GetMethodParams(method string) []any {
	called, ok := m.mockedMethod[method]
	if !ok {
		return nil
	}
	return called.Params
}

func (m *Mockable) IsParamsEqual(method string, params ...any) bool {
	called, ok := m.mockedMethod[method]
	if !ok {
		return false
	}
	if len(called.Params) != len(params) {
		return false
	}
	for i, param := range params {
		if !m.compareParams(called.Params[i], param) {
			return false
		}
	}
	return true
}

func (m *Mockable) compareParams(a, b any) bool {
	switch aTyped := a.(type) {
	case int, string:
		if a != b {
			fmt.Println("Param mismatch:", a, b)
			return false
		}
	case map[string]any:
		bTyped, ok := b.(map[string]any)
		for k, v := range aTyped {
			if !ok || !m.compareParams(v, bTyped[k]) {
				return false
			}
		}
	case []any:
		bTyped, ok := b.([]any)
		if !ok || len(aTyped) != len(bTyped) {
			return false
		}
		for i, v := range aTyped {
			if !m.compareParams(v, bTyped[i]) {
				return false
			}
		}
	case io.Reader:
		// Skip comparison for io.Reader types
		return true
	default:
		return a == b
	}

	return true
}

func (m *Mockable) CallFunc(method string) any {
	called, ok := m.mockedMethod[method]
	if !ok {
		return nil
	}
	if called.Func != nil {
		return called.Func(called.Params...)
	}
	return nil
}
