package subscriber

import (
	"context"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/component/pubsub"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/common/asyncjob"
)

type subJob struct {
	Name string
	Hdl  func(ctx context.Context, msg *pubsub.Message) error
}

type GroupJob interface {
	Run(ctx context.Context) error
}

type pbEngine struct {
	name   string
	ac     appctx.AppContext
	logger appctx.Logger
}

func NewPBEngine(ac appctx.AppContext) *pbEngine {
	return &pbEngine{
		name: "pb-engine",
		ac:   ac,
	}
}

func (engine *pbEngine) Start() error {
	engine.logger = engine.ac.Logger(engine.name)
	engine.startSubTopic(common.TopicUserLikedItem, true,
		IncreaseLikedCountAfterUserLikeItem(engine.ac),
	)

	engine.startSubTopic(common.TopicUserUnlikedItem, true,
		DecreaseLikedCountAfterUserUnlikeItem(engine.ac),
	)

	return nil
}

func (engine *pbEngine) startSubTopic(topic pubsub.Topic, isConcurrent bool, jobs ...subJob) error {
	ps := engine.ac.MustGet(common.PluginPubSub).(core.PubSubComponent)

	c, _ := ps.Subscribe(context.Background(), topic)
	for _, item := range jobs {
		engine.logger.Info("Setup subscriber :", item.Name)
	}

	getJobHandler := func(job *subJob, msg *pubsub.Message) asyncjob.JobHandler {
		return func(ctx context.Context) error {
			engine.logger.Infof("Run job [%s] - Value: %v", job.Name, msg.Data())
			return job.Hdl(ctx, msg)
		}
	}

	go func() {
		for {
			msg := <-c

			jobHdls := make([]asyncjob.Job, len(jobs))
			for index := range jobs {
				jobHdlIdnex := getJobHandler(&jobs[index], msg)
				jobHdls[index] = asyncjob.NewJob(jobHdlIdnex, asyncjob.WithName(jobs[index].Name))
			}

			group := asyncjob.NewGroup(isConcurrent, jobHdls...)
			if err := group.Run(context.Background()); err != nil {
				engine.logger.Error(err)
			}
		}
	}()

	return nil
}
