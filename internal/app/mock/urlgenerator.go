// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/service/urlgenerator/urlgenerator.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/service/urlgenerator/urlgenerator.go -destination=internal/app/mock/urlgenerator.go -package=mock URLGenerator
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
)

// MockURLGenerator is a mock of URLGenerator interface.
type MockURLGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockURLGeneratorMockRecorder
	isgomock struct{}
}

// MockURLGeneratorMockRecorder is the mock recorder for MockURLGenerator.
type MockURLGeneratorMockRecorder struct {
	mock *MockURLGenerator
}

// NewMockURLGenerator creates a new mock instance.
func NewMockURLGenerator(ctrl *gomock.Controller) *MockURLGenerator {
	mock := &MockURLGenerator{ctrl: ctrl}
	mock.recorder = &MockURLGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLGenerator) EXPECT() *MockURLGeneratorMockRecorder {
	return m.recorder
}

// GenerateSlug mocks base method.
func (m *MockURLGenerator) GenerateSlug(ctx context.Context, original domain.OriginalURL) domain.Slug {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateSlug", ctx, original)
	ret0, _ := ret[0].(domain.Slug)
	return ret0
}

// GenerateSlug indicates an expected call of GenerateSlug.
func (mr *MockURLGeneratorMockRecorder) GenerateSlug(ctx, original any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateSlug", reflect.TypeOf((*MockURLGenerator)(nil).GenerateSlug), ctx, original)
}

// GenerateSlugs mocks base method.
func (m *MockURLGenerator) GenerateSlugs(ctx context.Context, originals []domain.OriginalURL) ([]domain.Slug, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateSlugs", ctx, originals)
	ret0, _ := ret[0].([]domain.Slug)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateSlugs indicates an expected call of GenerateSlugs.
func (mr *MockURLGeneratorMockRecorder) GenerateSlugs(ctx, originals any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateSlugs", reflect.TypeOf((*MockURLGenerator)(nil).GenerateSlugs), ctx, originals)
}

// IsValidSlug mocks base method.
func (m *MockURLGenerator) IsValidSlug(slug domain.Slug) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsValidSlug", slug)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsValidSlug indicates an expected call of IsValidSlug.
func (mr *MockURLGeneratorMockRecorder) IsValidSlug(slug any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsValidSlug", reflect.TypeOf((*MockURLGenerator)(nil).IsValidSlug), slug)
}
