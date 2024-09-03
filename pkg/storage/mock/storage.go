// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go
//
// Generated by this command:
//
//	mockgen -source=storage.go -package=mock -destination=./mock/storage.go
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	storage "github.com/bucketeer-io/bucketeer/pkg/storage"
)

// MockIterator is a mock of Iterator interface.
type MockIterator struct {
	ctrl     *gomock.Controller
	recorder *MockIteratorMockRecorder
}

// MockIteratorMockRecorder is the mock recorder for MockIterator.
type MockIteratorMockRecorder struct {
	mock *MockIterator
}

// NewMockIterator creates a new mock instance.
func NewMockIterator(ctrl *gomock.Controller) *MockIterator {
	mock := &MockIterator{ctrl: ctrl}
	mock.recorder = &MockIteratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIterator) EXPECT() *MockIteratorMockRecorder {
	return m.recorder
}

// Cursor mocks base method.
func (m *MockIterator) Cursor() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cursor")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Cursor indicates an expected call of Cursor.
func (mr *MockIteratorMockRecorder) Cursor() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cursor", reflect.TypeOf((*MockIterator)(nil).Cursor))
}

// Next mocks base method.
func (m *MockIterator) Next(dst any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Next", dst)
	ret0, _ := ret[0].(error)
	return ret0
}

// Next indicates an expected call of Next.
func (mr *MockIteratorMockRecorder) Next(dst any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockIterator)(nil).Next), dst)
}

// MockGetter is a mock of Getter interface.
type MockGetter struct {
	ctrl     *gomock.Controller
	recorder *MockGetterMockRecorder
}

// MockGetterMockRecorder is the mock recorder for MockGetter.
type MockGetterMockRecorder struct {
	mock *MockGetter
}

// NewMockGetter creates a new mock instance.
func NewMockGetter(ctrl *gomock.Controller) *MockGetter {
	mock := &MockGetter{ctrl: ctrl}
	mock.recorder = &MockGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGetter) EXPECT() *MockGetterMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockGetter) Get(ctx context.Context, key *storage.Key, dst any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key, dst)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockGetterMockRecorder) Get(ctx, key, dst any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockGetter)(nil).Get), ctx, key, dst)
}

// GetMulti mocks base method.
func (m *MockGetter) GetMulti(ctx context.Context, keys []*storage.Key, dst any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMulti", ctx, keys, dst)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetMulti indicates an expected call of GetMulti.
func (mr *MockGetterMockRecorder) GetMulti(ctx, keys, dst any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMulti", reflect.TypeOf((*MockGetter)(nil).GetMulti), ctx, keys, dst)
}

// MockPutter is a mock of Putter interface.
type MockPutter struct {
	ctrl     *gomock.Controller
	recorder *MockPutterMockRecorder
}

// MockPutterMockRecorder is the mock recorder for MockPutter.
type MockPutterMockRecorder struct {
	mock *MockPutter
}

// NewMockPutter creates a new mock instance.
func NewMockPutter(ctrl *gomock.Controller) *MockPutter {
	mock := &MockPutter{ctrl: ctrl}
	mock.recorder = &MockPutterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPutter) EXPECT() *MockPutterMockRecorder {
	return m.recorder
}

// Put mocks base method.
func (m *MockPutter) Put(ctx context.Context, key *storage.Key, src any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, key, src)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockPutterMockRecorder) Put(ctx, key, src any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockPutter)(nil).Put), ctx, key, src)
}

// PutMulti mocks base method.
func (m *MockPutter) PutMulti(ctx context.Context, keys []*storage.Key, src any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutMulti", ctx, keys, src)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutMulti indicates an expected call of PutMulti.
func (mr *MockPutterMockRecorder) PutMulti(ctx, keys, src any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutMulti", reflect.TypeOf((*MockPutter)(nil).PutMulti), ctx, keys, src)
}

// MockGetPutter is a mock of GetPutter interface.
type MockGetPutter struct {
	ctrl     *gomock.Controller
	recorder *MockGetPutterMockRecorder
}

// MockGetPutterMockRecorder is the mock recorder for MockGetPutter.
type MockGetPutterMockRecorder struct {
	mock *MockGetPutter
}

// NewMockGetPutter creates a new mock instance.
func NewMockGetPutter(ctrl *gomock.Controller) *MockGetPutter {
	mock := &MockGetPutter{ctrl: ctrl}
	mock.recorder = &MockGetPutterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGetPutter) EXPECT() *MockGetPutterMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockGetPutter) Get(ctx context.Context, key *storage.Key, dst any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key, dst)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockGetPutterMockRecorder) Get(ctx, key, dst any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockGetPutter)(nil).Get), ctx, key, dst)
}

// GetMulti mocks base method.
func (m *MockGetPutter) GetMulti(ctx context.Context, keys []*storage.Key, dst any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMulti", ctx, keys, dst)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetMulti indicates an expected call of GetMulti.
func (mr *MockGetPutterMockRecorder) GetMulti(ctx, keys, dst any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMulti", reflect.TypeOf((*MockGetPutter)(nil).GetMulti), ctx, keys, dst)
}

// Put mocks base method.
func (m *MockGetPutter) Put(ctx context.Context, key *storage.Key, src any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, key, src)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockGetPutterMockRecorder) Put(ctx, key, src any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockGetPutter)(nil).Put), ctx, key, src)
}

// PutMulti mocks base method.
func (m *MockGetPutter) PutMulti(ctx context.Context, keys []*storage.Key, src any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutMulti", ctx, keys, src)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutMulti indicates an expected call of PutMulti.
func (mr *MockGetPutterMockRecorder) PutMulti(ctx, keys, src any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutMulti", reflect.TypeOf((*MockGetPutter)(nil).PutMulti), ctx, keys, src)
}

// MockQuerier is a mock of Querier interface.
type MockQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockQuerierMockRecorder
}

// MockQuerierMockRecorder is the mock recorder for MockQuerier.
type MockQuerierMockRecorder struct {
	mock *MockQuerier
}

// NewMockQuerier creates a new mock instance.
func NewMockQuerier(ctrl *gomock.Controller) *MockQuerier {
	mock := &MockQuerier{ctrl: ctrl}
	mock.recorder = &MockQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuerier) EXPECT() *MockQuerierMockRecorder {
	return m.recorder
}

// RunQuery mocks base method.
func (m *MockQuerier) RunQuery(ctx context.Context, query storage.Query) (storage.Iterator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunQuery", ctx, query)
	ret0, _ := ret[0].(storage.Iterator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunQuery indicates an expected call of RunQuery.
func (mr *MockQuerierMockRecorder) RunQuery(ctx, query any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunQuery", reflect.TypeOf((*MockQuerier)(nil).RunQuery), ctx, query)
}

// MockTransaction is a mock of Transaction interface.
type MockTransaction struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionMockRecorder
}

// MockTransactionMockRecorder is the mock recorder for MockTransaction.
type MockTransactionMockRecorder struct {
	mock *MockTransaction
}

// NewMockTransaction creates a new mock instance.
func NewMockTransaction(ctrl *gomock.Controller) *MockTransaction {
	mock := &MockTransaction{ctrl: ctrl}
	mock.recorder = &MockTransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransaction) EXPECT() *MockTransactionMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockTransaction) Delete(ctx context.Context, key *storage.Key) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTransactionMockRecorder) Delete(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTransaction)(nil).Delete), ctx, key)
}

// Get mocks base method.
func (m *MockTransaction) Get(ctx context.Context, key *storage.Key, dst any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key, dst)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockTransactionMockRecorder) Get(ctx, key, dst any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTransaction)(nil).Get), ctx, key, dst)
}

// GetMulti mocks base method.
func (m *MockTransaction) GetMulti(ctx context.Context, keys []*storage.Key, dst any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMulti", ctx, keys, dst)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetMulti indicates an expected call of GetMulti.
func (mr *MockTransactionMockRecorder) GetMulti(ctx, keys, dst any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMulti", reflect.TypeOf((*MockTransaction)(nil).GetMulti), ctx, keys, dst)
}

// Put mocks base method.
func (m *MockTransaction) Put(ctx context.Context, key *storage.Key, src any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, key, src)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockTransactionMockRecorder) Put(ctx, key, src any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockTransaction)(nil).Put), ctx, key, src)
}

// PutMulti mocks base method.
func (m *MockTransaction) PutMulti(ctx context.Context, keys []*storage.Key, src any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutMulti", ctx, keys, src)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutMulti indicates an expected call of PutMulti.
func (mr *MockTransactionMockRecorder) PutMulti(ctx, keys, src any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutMulti", reflect.TypeOf((*MockTransaction)(nil).PutMulti), ctx, keys, src)
}

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockClient) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClient)(nil).Close))
}

// Delete mocks base method.
func (m *MockClient) Delete(ctx context.Context, key *storage.Key) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockClientMockRecorder) Delete(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockClient)(nil).Delete), ctx, key)
}

// Get mocks base method.
func (m *MockClient) Get(ctx context.Context, key *storage.Key, dst any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key, dst)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockClientMockRecorder) Get(ctx, key, dst any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockClient)(nil).Get), ctx, key, dst)
}

// GetMulti mocks base method.
func (m *MockClient) GetMulti(ctx context.Context, keys []*storage.Key, dst any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMulti", ctx, keys, dst)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetMulti indicates an expected call of GetMulti.
func (mr *MockClientMockRecorder) GetMulti(ctx, keys, dst any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMulti", reflect.TypeOf((*MockClient)(nil).GetMulti), ctx, keys, dst)
}

// Put mocks base method.
func (m *MockClient) Put(ctx context.Context, key *storage.Key, src any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, key, src)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockClientMockRecorder) Put(ctx, key, src any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockClient)(nil).Put), ctx, key, src)
}

// PutMulti mocks base method.
func (m *MockClient) PutMulti(ctx context.Context, keys []*storage.Key, src any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutMulti", ctx, keys, src)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutMulti indicates an expected call of PutMulti.
func (mr *MockClientMockRecorder) PutMulti(ctx, keys, src any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutMulti", reflect.TypeOf((*MockClient)(nil).PutMulti), ctx, keys, src)
}

// RunInTransaction mocks base method.
func (m *MockClient) RunInTransaction(ctx context.Context, f func(storage.Transaction) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunInTransaction", ctx, f)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunInTransaction indicates an expected call of RunInTransaction.
func (mr *MockClientMockRecorder) RunInTransaction(ctx, f any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunInTransaction", reflect.TypeOf((*MockClient)(nil).RunInTransaction), ctx, f)
}

// RunQuery mocks base method.
func (m *MockClient) RunQuery(ctx context.Context, query storage.Query) (storage.Iterator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunQuery", ctx, query)
	ret0, _ := ret[0].(storage.Iterator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunQuery indicates an expected call of RunQuery.
func (mr *MockClientMockRecorder) RunQuery(ctx, query any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunQuery", reflect.TypeOf((*MockClient)(nil).RunQuery), ctx, query)
}

// MockWriter is a mock of Writer interface.
type MockWriter struct {
	ctrl     *gomock.Controller
	recorder *MockWriterMockRecorder
}

// MockWriterMockRecorder is the mock recorder for MockWriter.
type MockWriterMockRecorder struct {
	mock *MockWriter
}

// NewMockWriter creates a new mock instance.
func NewMockWriter(ctrl *gomock.Controller) *MockWriter {
	mock := &MockWriter{ctrl: ctrl}
	mock.recorder = &MockWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWriter) EXPECT() *MockWriterMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockWriter) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockWriterMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockWriter)(nil).Close))
}

// Write mocks base method.
func (m *MockWriter) Write(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockWriterMockRecorder) Write(p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockWriter)(nil).Write), p)
}

// MockReader is a mock of Reader interface.
type MockReader struct {
	ctrl     *gomock.Controller
	recorder *MockReaderMockRecorder
}

// MockReaderMockRecorder is the mock recorder for MockReader.
type MockReaderMockRecorder struct {
	mock *MockReader
}

// NewMockReader creates a new mock instance.
func NewMockReader(ctrl *gomock.Controller) *MockReader {
	mock := &MockReader{ctrl: ctrl}
	mock.recorder = &MockReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReader) EXPECT() *MockReaderMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockReader) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockReaderMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockReader)(nil).Close))
}

// Read mocks base method.
func (m *MockReader) Read(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockReaderMockRecorder) Read(p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockReader)(nil).Read), p)
}

// MockDeleter is a mock of Deleter interface.
type MockDeleter struct {
	ctrl     *gomock.Controller
	recorder *MockDeleterMockRecorder
}

// MockDeleterMockRecorder is the mock recorder for MockDeleter.
type MockDeleterMockRecorder struct {
	mock *MockDeleter
}

// NewMockDeleter creates a new mock instance.
func NewMockDeleter(ctrl *gomock.Controller) *MockDeleter {
	mock := &MockDeleter{ctrl: ctrl}
	mock.recorder = &MockDeleterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeleter) EXPECT() *MockDeleterMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockDeleter) Delete(ctx context.Context, key *storage.Key) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDeleterMockRecorder) Delete(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDeleter)(nil).Delete), ctx, key)
}

// MockObjectStorageClient is a mock of ObjectStorageClient interface.
type MockObjectStorageClient struct {
	ctrl     *gomock.Controller
	recorder *MockObjectStorageClientMockRecorder
}

// MockObjectStorageClientMockRecorder is the mock recorder for MockObjectStorageClient.
type MockObjectStorageClientMockRecorder struct {
	mock *MockObjectStorageClient
}

// NewMockObjectStorageClient creates a new mock instance.
func NewMockObjectStorageClient(ctrl *gomock.Controller) *MockObjectStorageClient {
	mock := &MockObjectStorageClient{ctrl: ctrl}
	mock.recorder = &MockObjectStorageClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockObjectStorageClient) EXPECT() *MockObjectStorageClientMockRecorder {
	return m.recorder
}

// Bucket mocks base method.
func (m *MockObjectStorageClient) Bucket(ctx context.Context, bucket string) (storage.Bucket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bucket", ctx, bucket)
	ret0, _ := ret[0].(storage.Bucket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Bucket indicates an expected call of Bucket.
func (mr *MockObjectStorageClientMockRecorder) Bucket(ctx, bucket any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bucket", reflect.TypeOf((*MockObjectStorageClient)(nil).Bucket), ctx, bucket)
}

// Close mocks base method.
func (m *MockObjectStorageClient) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockObjectStorageClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockObjectStorageClient)(nil).Close))
}

// MockBucket is a mock of Bucket interface.
type MockBucket struct {
	ctrl     *gomock.Controller
	recorder *MockBucketMockRecorder
}

// MockBucketMockRecorder is the mock recorder for MockBucket.
type MockBucketMockRecorder struct {
	mock *MockBucket
}

// NewMockBucket creates a new mock instance.
func NewMockBucket(ctrl *gomock.Controller) *MockBucket {
	mock := &MockBucket{ctrl: ctrl}
	mock.recorder = &MockBucketMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBucket) EXPECT() *MockBucketMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockBucket) Delete(ctx context.Context, key *storage.Key) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockBucketMockRecorder) Delete(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockBucket)(nil).Delete), ctx, key)
}

// Reader mocks base method.
func (m *MockBucket) Reader(ctx context.Context, environmentId, filename string) (storage.Reader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reader", ctx, environmentId, filename)
	ret0, _ := ret[0].(storage.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Reader indicates an expected call of Reader.
func (mr *MockBucketMockRecorder) Reader(ctx, environmentId, filename any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reader", reflect.TypeOf((*MockBucket)(nil).Reader), ctx, environmentId, filename)
}

// Writer mocks base method.
func (m *MockBucket) Writer(ctx context.Context, environmentId, filename string, CRC32C uint32) (storage.Writer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Writer", ctx, environmentId, filename, CRC32C)
	ret0, _ := ret[0].(storage.Writer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Writer indicates an expected call of Writer.
func (mr *MockBucketMockRecorder) Writer(ctx, environmentId, filename, CRC32C any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Writer", reflect.TypeOf((*MockBucket)(nil).Writer), ctx, environmentId, filename, CRC32C)
}

// MockObject is a mock of Object interface.
type MockObject struct {
	ctrl     *gomock.Controller
	recorder *MockObjectMockRecorder
}

// MockObjectMockRecorder is the mock recorder for MockObject.
type MockObjectMockRecorder struct {
	mock *MockObject
}

// NewMockObject creates a new mock instance.
func NewMockObject(ctrl *gomock.Controller) *MockObject {
	mock := &MockObject{ctrl: ctrl}
	mock.recorder = &MockObjectMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockObject) EXPECT() *MockObjectMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockObject) Delete(ctx context.Context, key *storage.Key) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockObjectMockRecorder) Delete(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockObject)(nil).Delete), ctx, key)
}

// Reader mocks base method.
func (m *MockObject) Reader(ctx context.Context, environmentId, filename string) (storage.Reader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reader", ctx, environmentId, filename)
	ret0, _ := ret[0].(storage.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Reader indicates an expected call of Reader.
func (mr *MockObjectMockRecorder) Reader(ctx, environmentId, filename any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reader", reflect.TypeOf((*MockObject)(nil).Reader), ctx, environmentId, filename)
}

// Writer mocks base method.
func (m *MockObject) Writer(ctx context.Context, environmentId, filename string, CRC32C uint32) (storage.Writer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Writer", ctx, environmentId, filename, CRC32C)
	ret0, _ := ret[0].(storage.Writer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Writer indicates an expected call of Writer.
func (mr *MockObjectMockRecorder) Writer(ctx, environmentId, filename, CRC32C any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Writer", reflect.TypeOf((*MockObject)(nil).Writer), ctx, environmentId, filename, CRC32C)
}
