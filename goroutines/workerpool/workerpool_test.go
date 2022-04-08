package workerpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WorkerSuite struct {
	suite.Suite
	tasks []Task
	limit int
}

func (s *WorkerSuite) SetupSuite() {
	s.tasks = []Task{
		Site{URL: "https://www.google.com"},
		Site{URL: "https://avito.ru"},
		Site{URL: "https://www.yandex.ru"},
		Site{URL: "https://www.mail.ru"},
	}
	s.limit = 128
}

func (s *WorkerSuite) TestAtomicWorkerPool() {
	wp := NewPoolAtomic(s.limit)
	wp.Run(s.tasks)
	res := wp.Stop()

	assert.Equal(s.T(), len(s.tasks), len(res))
	assert.Equal(s.T(), int32(0), wp.GetNumInProgress())
}

func (s *WorkerSuite) TestSyncWorkerPool() {
	wp := NewPoolSync(s.limit)
	wp.Run(s.tasks)
	res := wp.Stop()

	assert.Equal(s.T(), len(s.tasks), len(res))
}

func (s *WorkerSuite) TestChannelWorkerPool() {
	wp := NewPoolChannel(s.limit)
	_, err := wp.Run(s.tasks)
	assert.NoError(s.T(), err)
	res := wp.Stop()

	assert.Equal(s.T(), len(s.tasks), len(res))
}

func TestWorkerPoolSuite(t *testing.T) {
	suite.Run(t, new(WorkerSuite))
}
