package mocks

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Annotator" can be found in github.com/hypnoglow/chronologist/internal/grafana
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	grafana "github.com/hypnoglow/chronologist/internal/grafana"

	testify_assert "github.com/stretchr/testify/assert"
)

//AnnotatorMock implements github.com/hypnoglow/chronologist/internal/grafana.Annotator
type AnnotatorMock struct {
	t minimock.Tester

	DeleteAnnotationFunc       func(p context.Context, p1 int) (r error)
	DeleteAnnotationCounter    uint64
	DeleteAnnotationPreCounter uint64
	DeleteAnnotationMock       mAnnotatorMockDeleteAnnotation

	GetAnnotationsFunc       func(p context.Context, p1 grafana.GetAnnotationsParams) (r grafana.Annotations, r1 error)
	GetAnnotationsCounter    uint64
	GetAnnotationsPreCounter uint64
	GetAnnotationsMock       mAnnotatorMockGetAnnotations

	SaveAnnotationFunc       func(p context.Context, p1 grafana.Annotation) (r error)
	SaveAnnotationCounter    uint64
	SaveAnnotationPreCounter uint64
	SaveAnnotationMock       mAnnotatorMockSaveAnnotation
}

//NewAnnotatorMock returns a mock for github.com/hypnoglow/chronologist/internal/grafana.Annotator
func NewAnnotatorMock(t minimock.Tester) *AnnotatorMock {
	m := &AnnotatorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DeleteAnnotationMock = mAnnotatorMockDeleteAnnotation{mock: m}
	m.GetAnnotationsMock = mAnnotatorMockGetAnnotations{mock: m}
	m.SaveAnnotationMock = mAnnotatorMockSaveAnnotation{mock: m}

	return m
}

type mAnnotatorMockDeleteAnnotation struct {
	mock             *AnnotatorMock
	mockExpectations *AnnotatorMockDeleteAnnotationParams
}

//AnnotatorMockDeleteAnnotationParams represents input parameters of the Annotator.DeleteAnnotation
type AnnotatorMockDeleteAnnotationParams struct {
	p  context.Context
	p1 int
}

//Expect sets up expected params for the Annotator.DeleteAnnotation
func (m *mAnnotatorMockDeleteAnnotation) Expect(p context.Context, p1 int) *mAnnotatorMockDeleteAnnotation {
	m.mockExpectations = &AnnotatorMockDeleteAnnotationParams{p, p1}
	return m
}

//Return sets up a mock for Annotator.DeleteAnnotation to return Return's arguments
func (m *mAnnotatorMockDeleteAnnotation) Return(r error) *AnnotatorMock {
	m.mock.DeleteAnnotationFunc = func(p context.Context, p1 int) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Annotator.DeleteAnnotation method
func (m *mAnnotatorMockDeleteAnnotation) Set(f func(p context.Context, p1 int) (r error)) *AnnotatorMock {
	m.mock.DeleteAnnotationFunc = f
	m.mockExpectations = nil
	return m.mock
}

//DeleteAnnotation implements github.com/hypnoglow/chronologist/internal/grafana.Annotator interface
func (m *AnnotatorMock) DeleteAnnotation(p context.Context, p1 int) (r error) {
	atomic.AddUint64(&m.DeleteAnnotationPreCounter, 1)
	defer atomic.AddUint64(&m.DeleteAnnotationCounter, 1)

	if m.DeleteAnnotationMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.DeleteAnnotationMock.mockExpectations, AnnotatorMockDeleteAnnotationParams{p, p1},
			"Annotator.DeleteAnnotation got unexpected parameters")

		if m.DeleteAnnotationFunc == nil {

			m.t.Fatal("No results are set for the AnnotatorMock.DeleteAnnotation")

			return
		}
	}

	if m.DeleteAnnotationFunc == nil {
		m.t.Fatal("Unexpected call to AnnotatorMock.DeleteAnnotation")
		return
	}

	return m.DeleteAnnotationFunc(p, p1)
}

//DeleteAnnotationMinimockCounter returns a count of AnnotatorMock.DeleteAnnotationFunc invocations
func (m *AnnotatorMock) DeleteAnnotationMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeleteAnnotationCounter)
}

//DeleteAnnotationMinimockPreCounter returns the value of AnnotatorMock.DeleteAnnotation invocations
func (m *AnnotatorMock) DeleteAnnotationMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeleteAnnotationPreCounter)
}

type mAnnotatorMockGetAnnotations struct {
	mock             *AnnotatorMock
	mockExpectations *AnnotatorMockGetAnnotationsParams
}

//AnnotatorMockGetAnnotationsParams represents input parameters of the Annotator.GetAnnotations
type AnnotatorMockGetAnnotationsParams struct {
	p  context.Context
	p1 grafana.GetAnnotationsParams
}

//Expect sets up expected params for the Annotator.GetAnnotations
func (m *mAnnotatorMockGetAnnotations) Expect(p context.Context, p1 grafana.GetAnnotationsParams) *mAnnotatorMockGetAnnotations {
	m.mockExpectations = &AnnotatorMockGetAnnotationsParams{p, p1}
	return m
}

//Return sets up a mock for Annotator.GetAnnotations to return Return's arguments
func (m *mAnnotatorMockGetAnnotations) Return(r grafana.Annotations, r1 error) *AnnotatorMock {
	m.mock.GetAnnotationsFunc = func(p context.Context, p1 grafana.GetAnnotationsParams) (grafana.Annotations, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Annotator.GetAnnotations method
func (m *mAnnotatorMockGetAnnotations) Set(f func(p context.Context, p1 grafana.GetAnnotationsParams) (r grafana.Annotations, r1 error)) *AnnotatorMock {
	m.mock.GetAnnotationsFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetAnnotations implements github.com/hypnoglow/chronologist/internal/grafana.Annotator interface
func (m *AnnotatorMock) GetAnnotations(p context.Context, p1 grafana.GetAnnotationsParams) (r grafana.Annotations, r1 error) {
	atomic.AddUint64(&m.GetAnnotationsPreCounter, 1)
	defer atomic.AddUint64(&m.GetAnnotationsCounter, 1)

	if m.GetAnnotationsMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetAnnotationsMock.mockExpectations, AnnotatorMockGetAnnotationsParams{p, p1},
			"Annotator.GetAnnotations got unexpected parameters")

		if m.GetAnnotationsFunc == nil {

			m.t.Fatal("No results are set for the AnnotatorMock.GetAnnotations")

			return
		}
	}

	if m.GetAnnotationsFunc == nil {
		m.t.Fatal("Unexpected call to AnnotatorMock.GetAnnotations")
		return
	}

	return m.GetAnnotationsFunc(p, p1)
}

//GetAnnotationsMinimockCounter returns a count of AnnotatorMock.GetAnnotationsFunc invocations
func (m *AnnotatorMock) GetAnnotationsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetAnnotationsCounter)
}

//GetAnnotationsMinimockPreCounter returns the value of AnnotatorMock.GetAnnotations invocations
func (m *AnnotatorMock) GetAnnotationsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetAnnotationsPreCounter)
}

type mAnnotatorMockSaveAnnotation struct {
	mock             *AnnotatorMock
	mockExpectations *AnnotatorMockSaveAnnotationParams
}

//AnnotatorMockSaveAnnotationParams represents input parameters of the Annotator.SaveAnnotation
type AnnotatorMockSaveAnnotationParams struct {
	p  context.Context
	p1 grafana.Annotation
}

//Expect sets up expected params for the Annotator.SaveAnnotation
func (m *mAnnotatorMockSaveAnnotation) Expect(p context.Context, p1 grafana.Annotation) *mAnnotatorMockSaveAnnotation {
	m.mockExpectations = &AnnotatorMockSaveAnnotationParams{p, p1}
	return m
}

//Return sets up a mock for Annotator.SaveAnnotation to return Return's arguments
func (m *mAnnotatorMockSaveAnnotation) Return(r error) *AnnotatorMock {
	m.mock.SaveAnnotationFunc = func(p context.Context, p1 grafana.Annotation) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Annotator.SaveAnnotation method
func (m *mAnnotatorMockSaveAnnotation) Set(f func(p context.Context, p1 grafana.Annotation) (r error)) *AnnotatorMock {
	m.mock.SaveAnnotationFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SaveAnnotation implements github.com/hypnoglow/chronologist/internal/grafana.Annotator interface
func (m *AnnotatorMock) SaveAnnotation(p context.Context, p1 grafana.Annotation) (r error) {
	atomic.AddUint64(&m.SaveAnnotationPreCounter, 1)
	defer atomic.AddUint64(&m.SaveAnnotationCounter, 1)

	if m.SaveAnnotationMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SaveAnnotationMock.mockExpectations, AnnotatorMockSaveAnnotationParams{p, p1},
			"Annotator.SaveAnnotation got unexpected parameters")

		if m.SaveAnnotationFunc == nil {

			m.t.Fatal("No results are set for the AnnotatorMock.SaveAnnotation")

			return
		}
	}

	if m.SaveAnnotationFunc == nil {
		m.t.Fatal("Unexpected call to AnnotatorMock.SaveAnnotation")
		return
	}

	return m.SaveAnnotationFunc(p, p1)
}

//SaveAnnotationMinimockCounter returns a count of AnnotatorMock.SaveAnnotationFunc invocations
func (m *AnnotatorMock) SaveAnnotationMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SaveAnnotationCounter)
}

//SaveAnnotationMinimockPreCounter returns the value of AnnotatorMock.SaveAnnotation invocations
func (m *AnnotatorMock) SaveAnnotationMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SaveAnnotationPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AnnotatorMock) ValidateCallCounters() {

	if m.DeleteAnnotationFunc != nil && atomic.LoadUint64(&m.DeleteAnnotationCounter) == 0 {
		m.t.Fatal("Expected call to AnnotatorMock.DeleteAnnotation")
	}

	if m.GetAnnotationsFunc != nil && atomic.LoadUint64(&m.GetAnnotationsCounter) == 0 {
		m.t.Fatal("Expected call to AnnotatorMock.GetAnnotations")
	}

	if m.SaveAnnotationFunc != nil && atomic.LoadUint64(&m.SaveAnnotationCounter) == 0 {
		m.t.Fatal("Expected call to AnnotatorMock.SaveAnnotation")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AnnotatorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *AnnotatorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *AnnotatorMock) MinimockFinish() {

	if m.DeleteAnnotationFunc != nil && atomic.LoadUint64(&m.DeleteAnnotationCounter) == 0 {
		m.t.Fatal("Expected call to AnnotatorMock.DeleteAnnotation")
	}

	if m.GetAnnotationsFunc != nil && atomic.LoadUint64(&m.GetAnnotationsCounter) == 0 {
		m.t.Fatal("Expected call to AnnotatorMock.GetAnnotations")
	}

	if m.SaveAnnotationFunc != nil && atomic.LoadUint64(&m.SaveAnnotationCounter) == 0 {
		m.t.Fatal("Expected call to AnnotatorMock.SaveAnnotation")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *AnnotatorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *AnnotatorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.DeleteAnnotationFunc == nil || atomic.LoadUint64(&m.DeleteAnnotationCounter) > 0)
		ok = ok && (m.GetAnnotationsFunc == nil || atomic.LoadUint64(&m.GetAnnotationsCounter) > 0)
		ok = ok && (m.SaveAnnotationFunc == nil || atomic.LoadUint64(&m.SaveAnnotationCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.DeleteAnnotationFunc != nil && atomic.LoadUint64(&m.DeleteAnnotationCounter) == 0 {
				m.t.Error("Expected call to AnnotatorMock.DeleteAnnotation")
			}

			if m.GetAnnotationsFunc != nil && atomic.LoadUint64(&m.GetAnnotationsCounter) == 0 {
				m.t.Error("Expected call to AnnotatorMock.GetAnnotations")
			}

			if m.SaveAnnotationFunc != nil && atomic.LoadUint64(&m.SaveAnnotationCounter) == 0 {
				m.t.Error("Expected call to AnnotatorMock.SaveAnnotation")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

//AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
//it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *AnnotatorMock) AllMocksCalled() bool {

	if m.DeleteAnnotationFunc != nil && atomic.LoadUint64(&m.DeleteAnnotationCounter) == 0 {
		return false
	}

	if m.GetAnnotationsFunc != nil && atomic.LoadUint64(&m.GetAnnotationsCounter) == 0 {
		return false
	}

	if m.SaveAnnotationFunc != nil && atomic.LoadUint64(&m.SaveAnnotationCounter) == 0 {
		return false
	}

	return true
}
