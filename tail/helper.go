package tail

import (
	"github.com/viant/bqtail/shared"
	"github.com/viant/bqtail/stage/activity"
	"github.com/viant/bqtail/task"
	"github.com/viant/toolbox"
	"math/rand"
	"time"
)

func getRandom(min, max int) int {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	if max-min == 0 {
		max++
	}
	return min + int(rnd.Int63())%(max-min)
}

func updateJobID(eventID, jobID string) string {
	info := activity.Parse(jobID)
	info.EventID = eventID
	return info.GetJobID()
}

func buildJobIDReplacementMap(eventID string, actions []*task.Action) map[string]interface{} {
	var result = make(map[string]interface{})
	for i, action := range actions {
		jobID, ok := action.Request[shared.JobIDKey]
		if ok {
			info := activity.Parse(toolbox.AsString(jobID))
			info.EventID = eventID
			info.Step = i + 1
			result[shared.JobIDKey] = info.GetJobID()
			break
		}
	}
	return result
}
