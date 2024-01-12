package events

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) GetName() string {
	return e.Name
}
func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

func (e *TestEvent) GetPayLoad() interface{} {
	return e.Payload
}

type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	wg.Done()
}

type EventDispatcherTestSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	EventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.EventDispatcher = NewEventDispatcher()
	suite.handler = TestEventHandler{ID: 1}
	suite.handler2 = TestEventHandler{ID: 2}
	suite.handler3 = TestEventHandler{ID: 3}
	suite.event = TestEvent{Name: "test", Payload: "test"}
	suite.event2 = TestEvent{Name: "test2", Payload: "test2"}
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	err := suite.EventDispatcher.Register(suite.event.Name, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event.Name]))

	err = suite.EventDispatcher.Register(suite.event.Name, &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.EventDispatcher.handlers[suite.event.Name]))

	assert.Equal(suite.T(), &suite.handler, suite.EventDispatcher.handlers[suite.event.Name][0])
	assert.Equal(suite.T(), &suite.handler2, suite.EventDispatcher.handlers[suite.event.Name][1])

}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	//event 1
	err := suite.EventDispatcher.Register(suite.event.Name, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event.Name]))

	err = suite.EventDispatcher.Register(suite.event.Name, &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.EventDispatcher.handlers[suite.event.Name]))

	//event 2
	err = suite.EventDispatcher.Register(suite.event2.Name, &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event2.Name]))

	suite.EventDispatcher.Clear()
	suite.Equal(0, len(suite.EventDispatcher.handlers))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_RegisterWhenIsAlreadyRegistered() {
	err := suite.EventDispatcher.Register(suite.event.Name, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event.Name]))

	err = suite.EventDispatcher.Register(suite.event.Name, &suite.handler)
	suite.NotNil(err)
	assert.Equal(suite.T(), err, ErrEventAlreadyRegistered)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	err := suite.EventDispatcher.Register(suite.event.Name, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event.Name]))

	has := suite.EventDispatcher.Has(suite.event.Name, &suite.handler)
	suite.True(has)

	has = suite.EventDispatcher.Has(suite.event2.Name, &suite.handler)
	suite.False(has)

	has = suite.EventDispatcher.Has(suite.event.Name, &suite.handler2)
	suite.False(has)

}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Remove() {
	//event 1
	err := suite.EventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event.GetName()]))

	err = suite.EventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.EventDispatcher.handlers[suite.event.GetName()]))

	//event 2
	err = suite.EventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event2.GetName()]))

	err = suite.EventDispatcher.Remove(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event.GetName()]))
	suite.Equal(&suite.handler2, suite.EventDispatcher.handlers[suite.event.GetName()][0])

	err = suite.EventDispatcher.Remove(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(0, len(suite.EventDispatcher.handlers[suite.event.GetName()]))

	err = suite.EventDispatcher.Remove(suite.event2.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(0, len(suite.EventDispatcher.handlers[suite.event2.GetName()]))

}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	m.Called(event)
	wg.Done()
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
	mh := &MockHandler{}
	mh.On("Handle", &suite.event)
	suite.EventDispatcher.Register(suite.event.GetName(), mh)
	suite.EventDispatcher.Dispatch(&suite.event)
	mh.AssertExpectations(suite.T())
	mh.AssertNumberOfCalls(suite.T(), "Handle", 1)

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
