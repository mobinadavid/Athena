package job

import "github.com/amirhossein2831/message-brokering/job"

const LogQueue job.Queue = "log-queue"

type LogJob struct{}

func NewLogJob() *LogJob {
	return &LogJob{}
}

func (j *LogJob) GetQueue() job.Queue {
	return LogQueue
}

func (j *LogJob) Process(payload []byte) error {
	println("Process login job" + string(payload))

	return nil
}
