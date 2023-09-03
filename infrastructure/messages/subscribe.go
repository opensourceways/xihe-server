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
	handlerNameBigModel        = "bigmodel"
	handlerNameFollowing       = "following"
	handlerNameLike            = "like"
	handlerNameFork            = "fork"
	handlerNameDownload        = "download"
	handlerNameRelatedResource = "related_resource"
	handlerNameTraining        = "training"
	handlerNameFinetune        = "finetune"
	handlerNameInference       = "inference"
	handlerNameEvaluate        = "evaluate"
	handlerNameCloud           = "cloud"
)

func Subscribe(ctx context.Context, handler interface{}, log *logrus.Entry) error {
	// following
	if err := registerHandler(
		topics.Following, handlerNameFollowing, handlerForFollowing(handler),
	); err != nil {
		return err
	}

	// like
	if err := registerHandler(
		topics.Like, handlerNameLike, handlerForLike(handler),
	); err != nil {
		return err
	}

	// fork
	if err := registerHandler(
		topics.Fork, handlerNameFork, handlerForFork(handler),
	); err != nil {
		return err
	}

	// download
	if err := registerHandler(
		topics.Download, handlerNameDownload, handlerForDownload(handler),
	); err != nil {
		return err
	}

	// related resource
	if err := registerHandler(
		topics.RelatedResource, handlerNameRelatedResource, handlerForRelatedResource(handler),
	); err != nil {
		return err
	}

	// training
	if err := registerHandler(
		topics.Training, handlerNameTraining, handlerForTraining(handler),
	); err != nil {
		return err
	}

	// finetune
	if err := registerHandler(
		topics.Finetune, handlerNameFinetune, handlerForFinetune(handler),
	); err != nil {
		return err
	}

	// inference
	if err := registerHandler(
		topics.Inference, handlerNameInference, handlerForInference(handler),
	); err != nil {
		return err
	}

	// evaluate
	if err := registerHandler(
		topics.Evaluate, handlerNameEvaluate, handlerForEvaluate(handler),
	); err != nil {
		return err
	}

	// cloud
	if err := registerHandler(
		topics.Cloud, handlerNameCloud, handlerForCloud(handler),
	); err != nil {
		return err
	}

	// bigmodel
	if err := registerHandler(
		topics.BigModel, handlerNameBigModel, handlerForBigModel(handler),
	); err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func handlerForFollowing(handler interface{}) kfklib.Handler {
	return func(b []byte, h map[string]string) (err error) {
		body := msgFollower{}
		if err = json.Unmarshal(b, &body); err != nil {
			return
		}

		hd, ok := handler.(message.FollowingHandler)
		if !ok {
			return errors.New("internal error, FollowingHandler assert error")
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
			return hd.HandleEventAddFollowing(f)

		case actionRemove:
			return hd.HandleEventRemoveFollowing(f)
		}

		return nil
	}
}

func handlerForLike(handler interface{}) kfklib.Handler {
	return func(b []byte, h map[string]string) (err error) {
		body := msgLike{}
		if err = json.Unmarshal(b, &body); err != nil {
			return
		}

		hd, ok := handler.(message.LikeHandler)
		if !ok {
			return errors.New("internal error, LikeHandler assert error")
		}

		like := &domain.ResourceObject{}
		if err = body.Resource.toResourceObject(like); err != nil {
			return
		}

		switch body.Action {
		case actionAdd:
			return hd.HandleEventAddLike(like)

		case actionRemove:
			return hd.HandleEventRemoveLike(like)
		}

		return
	}
}

func handlerForFork(handler interface{}) kfklib.Handler {
	return func(b []byte, h map[string]string) (err error) {
		body := resourceIndex{}
		if err = json.Unmarshal(b, &body); err != nil {
			return
		}

		hd, ok := handler.(message.ForkHandler)
		if !ok {
			return errors.New("internal error, ForkHandler assert error")
		}

		index := new(domain.ResourceIndex)
		if err = body.toResourceIndex(index); err != nil {
			return
		}

		return hd.HandleEventFork(index)
	}
}

func handlerForDownload(handler interface{}) kfklib.Handler {
	return func(b []byte, h map[string]string) (err error) {
		body := resourceObject{}
		if err = json.Unmarshal(b, &body); err != nil {
			return
		}

		hd, ok := handler.(message.DownloadHandler)
		if !ok {
			return errors.New("internal error, DownloadHandler assert error")
		}

		obj := new(domain.ResourceObject)
		if err = body.toResourceObject(obj); err != nil {
			return
		}

		return hd.HandleEventDownload(obj)
	}
}

func handlerForRelatedResource(handler interface{}) kfklib.Handler {
	return func(b []byte, h map[string]string) (err error) {
		body := msgRelatedResources{}
		if err = json.Unmarshal(b, &body); err != nil {
			return
		}

		hd, ok := handler.(message.RelatedResourceHandler)
		if !ok {
			return errors.New("internal error, RelatedResourceHandler assert error")
		}

		switch body.Action {
		case actionAdd:
			return body.handle(hd.HandleEventAddRelatedResource)

		case actionRemove:
			return body.handle(hd.HandleEventRemoveRelatedResource)
		}

		return nil
	}
}

func handlerForTraining(handler interface{}) kfklib.Handler {
	return func(b []byte, h map[string]string) (err error) {
		body := message.MsgTraining{}
		if err = json.Unmarshal(b, &body); err != nil {
			return
		}

		hd, ok := handler.(message.TrainingHandler)
		if !ok {
			return errors.New("internal error, TrainingHandler assert error")
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

		return hd.HandleEventCreateTraining(&v)
	}
}

func handlerForFinetune(handler interface{}) kfklib.Handler {
	return func(b []byte, h map[string]string) (err error) {
		body := msgFinetune{}
		if err = json.Unmarshal(b, &body); err != nil {
			return
		}

		hd, ok := handler.(message.FinetuneHandler)
		if !ok {
			return errors.New("internal error, FinetuneHandler assert error")
		}

		if body.Id == "" {
			err = errors.New("invalid message of finetune")

			return
		}

		v := domain.FinetuneIndex{Id: body.Id}
		if v.Owner, err = domain.NewAccount(body.User); err != nil {
			return
		}

		return hd.HandleEventCreateFinetune(&v)
	}
}

func handlerForInference(handler interface{}) kfklib.Handler {
	return func(b []byte, h map[string]string) (err error) {
		body := msgInference{}
		if err = json.Unmarshal(b, &body); err != nil {
			return
		}

		hd, ok := handler.(message.InferenceHandler)
		if !ok {
			return errors.New("internal error, InferenceHandler assert error")
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
			return hd.HandleEventCreateInference(&info)

		case actionExtend:
			return hd.HandleEventExtendInferenceSurvivalTime(
				&message.InferenceExtendInfo{
					InferenceInfo: info,
					Expiry:        body.Expiry,
				},
			)
		}

		return nil
	}
}

func handlerForEvaluate(handler interface{}) kfklib.Handler {
	return func(b []byte, h map[string]string) (err error) {
		body := msgEvaluate{}
		if err := json.Unmarshal(b, &body); err != nil {
			return err
		}

		hd, ok := handler.(message.EvaluateHandler)
		if !ok {
			return errors.New("internal error, CloudMessageHandler assert error")
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

		return hd.HandleEventCreateEvaluate(&v)
	}
}

func handlerForCloud(handler interface{}) kfklib.Handler {
	return func(b []byte, h map[string]string) error {
		body := msgPodCreate{}
		if err := json.Unmarshal(b, &body); err != nil {
			return err
		}

		hd, ok := handler.(cloudmsg.CloudMessageHandler)
		if !ok {
			return errors.New("internal error, CloudMessageHandler assert error")
		}

		user, err := domain.NewAccount(body.User)
		if err != nil {
			return err
		}

		v := cloudtypes.PodInfo{
			Pod: cloudtypes.Pod{
				Id:      body.PodId,
				CloudId: body.CloudId,
				Owner:   user,
			},
		}
		v.SetDefaultExpiry()

		return hd.HandleEventPodSubscribe(&v)
	}
}

func handlerForBigModel(handler interface{}) kfklib.Handler {
	return func(b []byte, h map[string]string) error {
		body := bigmoddelmsg.MsgTask{}
		if err := json.Unmarshal(b, &body); err != nil {
			return err
		}

		hd, ok := handler.(BigModelMessageHandler)
		if !ok {
			return errors.New("internal error, BigModelMessageHandler assert error")
		}

		switch body.Type {
		case bigmoddelmsg.MsgTypeWuKongAsyncTaskFinish:

			return hd.HandleEventBigModelWuKongAsyncTaskFinish(&body)

		case bigmoddelmsg.MsgTypeWuKongAsyncTaskStart:

			return hd.HandleEventBigModelWuKongAsyncTaskStart(&body)

		case bigmoddelmsg.MsgTypeWuKongInferenceStart:

			return hd.HandleEventBigModelWuKongInferenceStart(&body)

		case bigmoddelmsg.MsgTypeWuKongInferenceError:

			return hd.HandleEventBigModelWuKongInferenceError(&body)

		}

		return nil
	}

}

func registerHandler(topic, handlerName string, h kfklib.Handler) error {
	return kfklib.SubscribeWithStrategyOfRetry(
		handlerName+"-"+topic, h, []string{topic}, 3,
	)
}
