package messages

import (
	"context"
	"encoding/json"
	"errors"

	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/sirupsen/logrus"

	bigmoddelmsg "github.com/opensourceways/xihe-server/bigmodel/domain/message"
	cloudtypes "github.com/opensourceways/xihe-server/cloud/domain"
	cloudmsg "github.com/opensourceways/xihe-server/cloud/domain/message"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	userdomain "github.com/opensourceways/xihe-server/user/domain"
)

const (
	handlerNameCreateBigModel     = "create_bigmodel"
	handlerNameAddFollowing       = "add_following"
	handlerNameAddLike            = "add_like"
	handlerNameAddFork            = "add_fork"
	handlerNameAddDownload        = "add_download"
	handlerNameAddRelatedResource = "add_related_resource"
	handlerNameCreateTraining     = "create_training"
	handlerNameCreateFinetune     = "create_finetune"
	handlerNameCreateInference    = "create_inference"
	handlerNameCreateEvaluate     = "create_evaluate"
	handlerNameCreateCloud        = "create_cloud"
)

func Subscribe(ctx context.Context, handler interface{}, log *logrus.Entry) (err error) {
	// register following
	if err = registerHandlerForFollowing(handler); err != nil {
		return
	}

	// register like
	if err = registerHandlerForLike(handler); err != nil {
		return
	}

	// register fork
	if err = registerHandlerForFork(handler); err != nil {
		return
	}

	// register download
	if err = registerHandlerForDownload(handler); err != nil {
		return
	}

	// register related resource
	if err = registerHandlerForRelatedResource(handler); err != nil {
		return
	}

	// training
	if err = registerHandlerForTraining(handler); err != nil {
		return
	}

	// finetune
	if err = registerHandlerForFinetune(handler); err != nil {
		return
	}

	// inference
	if err = registerHandlerForInference(handler); err != nil {
		return
	}

	// evaluate
	if err = registerHandlerForEvaluate(handler); err != nil {
		return
	}

	// cloud
	if err = registerHandlerForCloud(handler); err != nil {
		return
	}

	// bigmodel
	if err = registerHandlerForBigModel(handler); err != nil {
		return err
	}

	// register end
	<-ctx.Done()

	return nil
}

func registerHandlerForFollowing(handler interface{}) error {
	h, ok := handler.(message.FollowingHandler)
	if !ok {
		return nil
	}

	return subscribe(topics.Following, handlerNameAddFollowing, func(b []byte, hd map[string]string) (err error) {
		body := msgFollower{}
		if err = json.Unmarshal(b, &body); err != nil {
			return
		}

		f := &userdomain.FollowerInfo{}
		if f.User, err = userdomain.NewAccount(body.User); err != nil {
			return
		}

		if f.Follower, err = userdomain.NewAccount(body.Follower); err != nil {
			return
		}

		switch body.Action {
		case actionAdd:
			return h.HandleEventAddFollowing(f)

		case actionRemove:
			return h.HandleEventRemoveFollowing(f)
		}

		return nil
	})
}

func registerHandlerForLike(handler interface{}) error {
	h, ok := handler.(message.LikeHandler)
	if !ok {
		return nil
	}

	return subscribe(topics.Like, handlerNameAddLike, func(b []byte, hd map[string]string) (err error) {
		body := msgLike{}
		if err = json.Unmarshal(b, &body); err != nil {
			return
		}

		like := &domain.ResourceObject{}
		if err = body.Resource.toResourceObject(like); err != nil {
			return
		}

		switch body.Action {
		case actionAdd:
			return h.HandleEventAddLike(like)

		case actionRemove:
			return h.HandleEventRemoveLike(like)
		}

		return
	})
}

func registerHandlerForFork(handler interface{}) error {
	h, ok := handler.(message.ForkHandler)
	if !ok {
		return nil
	}

	return subscribe(topics.Fork, handlerNameAddFork, func(b []byte, hd map[string]string) (err error) {
		body := resourceIndex{}
		if err = json.Unmarshal(b, &body); err != nil {
			return
		}

		index := new(domain.ResourceIndex)
		if err = body.toResourceIndex(index); err != nil {
			return
		}

		return h.HandleEventFork(index)
	})
}

func registerHandlerForDownload(handler interface{}) error {
	h, ok := handler.(message.DownloadHandler)
	if !ok {
		return nil
	}

	return subscribe(topics.Download, handlerNameAddDownload, func(b []byte, hd map[string]string) (err error) {
		body := resourceObject{}
		if err = json.Unmarshal(b, &body); err != nil {
			return
		}

		obj := new(domain.ResourceObject)
		if err = body.toResourceObject(obj); err != nil {
			return
		}

		return h.HandleEventDownload(obj)
	})
}

func registerHandlerForRelatedResource(handler interface{}) error {
	h, ok := handler.(message.RelatedResourceHandler)
	if !ok {
		return nil
	}

	return subscribe(
		topics.RelatedResource, handlerNameAddRelatedResource,
		func(b []byte, hd map[string]string) (err error) {
			body := msgRelatedResources{}
			if err = json.Unmarshal(b, &body); err != nil {
				return
			}

			switch body.Action {
			case actionAdd:
				return body.handle(h.HandleEventAddRelatedResource)

			case actionRemove:
				return body.handle(h.HandleEventRemoveRelatedResource)
			}

			return nil
		})
}

func registerHandlerForTraining(handler interface{}) error {
	h, ok := handler.(message.TrainingHandler)
	if !ok {
		return nil
	}

	return subscribe(topics.Training, handlerNameCreateTraining,
		func(b []byte, hd map[string]string) (err error) {
			body := message.MsgTraining{}
			if err = json.Unmarshal(b, &body); err != nil {
				return
			}

			if body.Details["project_id"] == "" || body.Details["training_id"] == "" {
				err = errors.New("invalid message of training")

				return
			}

			v := domain.TrainingIndex{}
			if v.Project.Owner, err = domain.NewAccount(body.Details["project_owner"]); err != nil {
				return
			}

			v.Project.Id = body.Details["project_id"]
			v.TrainingId = body.Details["training_id"]

			return h.HandleEventCreateTraining(&v)
		})
}

func registerHandlerForFinetune(handler interface{}) error {
	h, ok := handler.(message.FinetuneHandler)
	if !ok {
		return nil
	}

	return subscribe(topics.Finetune, handlerNameCreateFinetune,
		func(b []byte, m map[string]string) (err error) {
			body := msgFinetune{}
			if err = json.Unmarshal(b, &body); err != nil {
				return
			}

			if body.Id == "" {
				err = errors.New("invalid message of finetune")

				return
			}

			v := domain.FinetuneIndex{Id: body.Id}
			if v.Owner, err = domain.NewAccount(body.User); err != nil {
				return
			}

			return h.HandleEventCreateFinetune(&v)
		})
}

func registerHandlerForInference(handler interface{}) error {
	h, ok := handler.(message.InferenceHandler)
	if !ok {
		return nil
	}

	return subscribe(topics.Inference, handlerNameCreateInference,
		func(b []byte, m map[string]string) (err error) {
			body := msgInference{}
			if err = json.Unmarshal(b, &body); err != nil {
				return
			}

			v := domain.InferenceIndex{}

			if v.Project.Owner, err = domain.NewAccount(body.ProjectOwner); err != nil {
				return
			}

			v.Id = body.InferenceId
			v.Project.Id = body.ProjectId
			v.LastCommit = body.LastCommit

			info := domain.InferenceInfo{
				InferenceIndex: v,
			}

			info.ProjectName, err = domain.NewResourceName(body.ProjectName)
			if err != nil {
				return
			}

			info.ResourceLevel = body.ResourceLevel

			switch body.Action {
			case actionCreate:
				return h.HandleEventCreateInference(&info)

			case actionExtend:
				return h.HandleEventExtendInferenceSurvivalTime(
					&message.InferenceExtendInfo{
						InferenceInfo: info,
						Expiry:        body.Expiry,
					},
				)
			}

			return nil
		})
}

func registerHandlerForEvaluate(handler interface{}) error {
	h, ok := handler.(message.EvaluateHandler)
	if !ok {
		return nil
	}

	return subscribe(topics.Evaluate, handlerNameCreateEvaluate,
		func(b []byte, m map[string]string) (err error) {
			body := msgEvaluate{}
			if err = json.Unmarshal(b, &body); err != nil {
				return
			}

			v := message.EvaluateInfo{}

			if v.Project.Owner, err = domain.NewAccount(body.ProjectOwner); err != nil {
				return
			}

			v.Id = body.EvaluateId
			v.Type = body.Type
			v.OBSPath = body.OBSPath
			v.Project.Id = body.ProjectId
			v.TrainingId = body.TrainingId

			return h.HandleEventCreateEvaluate(&v)
		})
}

func registerHandlerForCloud(handler interface{}) error {
	h, ok := handler.(cloudmsg.CloudMessageHandler)
	if !ok {
		return nil
	}

	return subscribe(topics.Cloud, handlerNameCreateCloud,
		func(b []byte, m map[string]string) (err error) {
			body := msgPodCreate{}
			if err = json.Unmarshal(b, &body); err != nil {
				return
			}

			user, err := domain.NewAccount(body.User)
			if err != nil {
				return
			}

			v := cloudtypes.PodInfo{
				Pod: cloudtypes.Pod{
					Id:      body.PodId,
					CloudId: body.CloudId,
					Owner:   user,
				},
			}
			v.SetDefaultExpiry()

			return h.HandleEventPodSubscribe(&v)
		})
}

func registerHandlerForBigModel(handler interface{}) error {

	return subscribe(topics.BigModel, handlerNameCreateBigModel,
		func(b []byte, m map[string]string) (err error) {
			body := bigmoddelmsg.MsgTask{}
			if err = json.Unmarshal(b, &body); err != nil {
				return
			}

			h, ok := handler.(BigModelMessageHandler)
			if !ok {
				return
			}

			switch body.Type {
			case bigmoddelmsg.MsgTypeWuKongAsyncTaskFinish:

				return h.HandleEventBigModelWuKongAsyncTaskFinish(&body)

			case bigmoddelmsg.MsgTypeWuKongAsyncTaskStart:

				return h.HandleEventBigModelWuKongAsyncTaskStart(&body)

			case bigmoddelmsg.MsgTypeWuKongInferenceStart:

				return h.HandleEventBigModelWuKongInferenceStart(&body)

			case bigmoddelmsg.MsgTypeWuKongInferenceError:

				return h.HandleEventBigModelWuKongInferenceError(&body)

			}

			return
		})
}

func subscribe(topicName string, handlerName string, handler kfklib.Handler) error {
	return kfklib.SubscribeWithStrategyOfRetry(
		handlerName, handler, []string{topicName}, 3,
	)
}
