// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/service/shortener/shortener.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/service/shortener/shortener.go -destination=internal/app/mock/shortener.go -package=mock URLShortener
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	dto "github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
)

// MockURLShortener is a mock of URLShortener interface.
type MockURLShortener struct {
	ctrl     *gomock.Controller
	recorder *MockURLShortenerMockRecorder
	isgomock struct{}
}

// MockURLShortenerMockRecorder is the mock recorder for MockURLShortener.
type MockURLShortenerMockRecorder struct {
	mock *MockURLShortener
}

// NewMockURLShortener creates a new mock instance.
func NewMockURLShortener(ctrl *gomock.Controller) *MockURLShortener {
	mock := &MockURLShortener{ctrl: ctrl}
	mock.recorder = &MockURLShortenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLShortener) EXPECT() *MockURLShortenerMockRecorder {
	return m.recorder
}

// GetOriginalURL mocks base method.
func (m *MockURLShortener) GetOriginalURL(ctx context.Context, slug domain.Slug) (domain.OriginalURL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOriginalURL", ctx, slug)
	ret0, _ := ret[0].(domain.OriginalURL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOriginalURL indicates an expected call of GetOriginalURL.
func (mr *MockURLShortenerMockRecorder) GetOriginalURL(ctx, slug any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOriginalURL", reflect.TypeOf((*MockURLShortener)(nil).GetOriginalURL), ctx, slug)
}

// GetUserURLs mocks base method.
func (m *MockURLShortener) GetUserURLs(ctx context.Context) (*dto.URLPairBatch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserURLs", ctx)
	ret0, _ := ret[0].(*dto.URLPairBatch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserURLs indicates an expected call of GetUserURLs.
func (mr *MockURLShortenerMockRecorder) GetUserURLs(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserURLs", reflect.TypeOf((*MockURLShortener)(nil).GetUserURLs), ctx)
}

// ShortenURL mocks base method.
func (m *MockURLShortener) ShortenURL(ctx context.Context, original domain.OriginalURL) (domain.Slug, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShortenURL", ctx, original)
	ret0, _ := ret[0].(domain.Slug)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShortenURL indicates an expected call of ShortenURL.
func (mr *MockURLShortenerMockRecorder) ShortenURL(ctx, original any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShortenURL", reflect.TypeOf((*MockURLShortener)(nil).ShortenURL), ctx, original)
}

// ShortenURLBatch mocks base method.
func (m *MockURLShortener) ShortenURLBatch(ctx context.Context, batch *dto.OriginalURLBatch) (*dto.SlugBatch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShortenURLBatch", ctx, batch)
	ret0, _ := ret[0].(*dto.SlugBatch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShortenURLBatch indicates an expected call of ShortenURLBatch.
func (mr *MockURLShortenerMockRecorder) ShortenURLBatch(ctx, batch any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShortenURLBatch", reflect.TypeOf((*MockURLShortener)(nil).ShortenURLBatch), ctx, batch)
}
