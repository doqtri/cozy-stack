package jobs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cozy/cozy-stack/pkg/consts"
	"github.com/cozy/cozy-stack/pkg/jobs"
	_ "github.com/cozy/cozy-stack/pkg/jobs/workers" // import all workers
	"github.com/cozy/cozy-stack/web/jsonapi"
	"github.com/cozy/cozy-stack/web/middlewares"
	"github.com/labstack/echo"
)

const typeTextEventStream = "text/event-stream"

type (
	apiJob struct {
		j *jobs.JobInfos
	}
	apiJobRequest struct {
		Arguments json.RawMessage  `json:"arguments"`
		Options   *jobs.JobOptions `json:"options"`
	}

	apiQueue struct {
		Count      int `json:"count"`
		workerType string
	}
)

func (j *apiJob) ID() string                             { return j.j.ID }
func (j *apiJob) Rev() string                            { return "" }
func (j *apiJob) DocType() string                        { return consts.Jobs }
func (j *apiJob) SetID(_ string)                         {}
func (j *apiJob) SetRev(_ string)                        {}
func (j *apiJob) Relationships() jsonapi.RelationshipMap { return nil }
func (j *apiJob) Included() []jsonapi.Object             { return nil }
func (j *apiJob) SelfLink() string {
	return "/jobs/" + j.j.WorkerType + "/" + j.j.ID
}
func (j *apiJob) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.j)
}

func (q *apiQueue) ID() string                             { return q.workerType }
func (q *apiQueue) Rev() string                            { return "" }
func (q *apiQueue) DocType() string                        { return consts.Queues }
func (q *apiQueue) SetID(_ string)                         {}
func (q *apiQueue) SetRev(_ string)                        {}
func (q *apiQueue) Relationships() jsonapi.RelationshipMap { return nil }
func (q *apiQueue) Included() []jsonapi.Object             { return nil }
func (q *apiQueue) SelfLink() string {
	return "/jobs/queue/" + q.workerType
}

func getQueue(c echo.Context) error {
	instance := middlewares.GetInstance(c)
	workerType := c.Param("worker-type")
	count, err := instance.JobsBroker().QueueLen(workerType)
	if err != nil {
		return err
	}
	o := &apiQueue{
		workerType: workerType,
		Count:      count,
	}
	return jsonapi.Data(c, http.StatusOK, o, nil)
}

func pushJob(c echo.Context) error {
	instance := middlewares.GetInstance(c)

	req := &apiJobRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	job, ch, err := instance.JobsBroker().PushJob(&jobs.JobRequest{
		WorkerType: c.Param("worker-type"),
		Options:    req.Options,
		Message: &jobs.Message{
			Type: jobs.JSONEncoding,
			Data: req.Arguments,
		},
	})
	if err != nil {
		return wrapJobsError(err)
	}

	accept := c.Request().Header.Get("Accept")
	if accept != typeTextEventStream {
		return jsonapi.Data(c, http.StatusAccepted, &apiJob{job}, nil)
	}

	w := c.Response().Writer
	w.Header().Set("Content-Type", typeTextEventStream)
	w.WriteHeader(200)
	if err := streamJob(job, w); err != nil {
		return nil
	}
	for job = range ch {
		if err := streamJob(job, w); err != nil {
			return nil
		}
	}
	return nil
}

// Routes sets the routing for the jobs service
func Routes(router *echo.Group) {
	router.POST("/queue/:worker-type", pushJob)
	router.GET("/queue/:worker-type", getQueue)
}

func streamJob(job *jobs.JobInfos, w http.ResponseWriter) error {
	b, err := json.Marshal(job)
	if err != nil {
		return err
	}
	s := fmt.Sprintf("event: %s\r\ndata: %s\r\n\r\n", job.State, b)
	_, err = w.Write([]byte(s))
	if err != nil {
		return err
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

func wrapJobsError(err error) error {
	switch err {
	case jobs.ErrUnknownWorker:
		return jsonapi.NotFound(err)
	}
	return err
}