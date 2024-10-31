package messages

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/sirupsen/logrus"

	cloudtypes "github.com/opensourceways/xihe-server/cloud/domain"
	cloudmsg "github.com/opensourceways/xihe-server/cloud/domain/message"
	common "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	spaceappdomain "github.com/opensourceways/xihe-server/spaceapp/domain"
)

const (
	retryNum = 3

	handlerNameAddLike            = "add_like"
	handlerNameAddFork            = "add_fork"
	handlerNameAddDownload        = "add_download"
	handlerNameAddFollowing       = "add_following"
	handlerNameAddRelatedResource = "add_related_resource"
	handlerNameCreateCloud        = "create_cloud"
	handlerNameReleaseCloud       = "release_cloud"
	handlerNameCreateTraining     = "create_training"
	handlerNameCreateFinetune     = "create_finetune"
	handlerNameCreateInference    = "create_inference"
)

func Subscribe(
	ctx context.Context, handler interface{}, log *logrus.Entry,
	topic *Topics, subscriber common.Subscriber,
) (err error) {
	r := register{
		topics:     topic,
		subscriber: subscriber,
	}

	// register like
	if err = r.registerHandlerForLike(handler); err != nil {
		return
	}

	// register fork
	if err = r.registerHandlerForFork(handler); err != nil {
		return
	}

	// register download
	if err = r.registerHandlerForDownload(handler); err != nil {
		return
	}

	// register related resource
	if err = r.registerHandlerForRelatedResource(handler); err != nil {
		return
	}

	// finetune
	if err = r.registerHandlerForFinetune(handler); err != nil {
		return
	}

	// inference
	if err = r.registerHandlerForInference(handler); err != nil {
		return
	}

	// cloud
	if err = r.registerHandlerForCloud(handler); err != nil {
		return
	}
	if err = r.registerHandlerForReleasingCloud(handler); err != nil {
		return
	}

	// register end
	<-ctx.Done()

	return nil
}

type register struct {
	topics     *Topics
	subscriber common.Subscriber
}

func (r *register) registerHandlerForLike(handler interface{}) error {
	h, ok := handler.(message.LikeHandler)
	if !ok {
		return nil
	}

	return r.subscribe(r.topics.Like, handlerNameAddLike, func(b []byte, hd map[string]string) (err error) {
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

		default:
			logrus.Warn("unknown action")
			return
		}

	})
}

func (r *register) registerHandlerForFork(handler interface{}) error {
	h, ok := handler.(message.ForkHandler)
	if !ok {
		return nil
	}

	return r.subscribe(r.topics.Fork, handlerNameAddFork, func(b []byte, hd map[string]string) (err error) {
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

func (r *register) registerHandlerForDownload(handler interface{}) error {
	h, ok := handler.(message.DownloadHandler)
	if !ok {
		return nil
	}

	return r.subscribe(r.topics.Download, handlerNameAddDownload, func(b []byte, hd map[string]string) (err error) {
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

func (r *register) registerHandlerForRelatedResource(handler interface{}) error {
	h, ok := handler.(message.RelatedResourceHandler)
	if !ok {
		return nil
	}

	f := func(b []byte, hd map[string]string) (err error) {
		body := msgRelatedResources{}
		if err = json.Unmarshal(b, &body); err != nil {
			return
		}

		switch body.Action {
		case actionAdd:
			return body.handle(h.HandleEventAddRelatedResource)

		case actionRemove:
			return body.handle(h.HandleEventRemoveRelatedResource)

		default:
			logrus.Warn("unknown action")
			return nil
		}

	}

	return r.subscribe(
		r.topics.RelatedResource, handlerNameAddRelatedResource, f,
	)
}

func (r *register) registerHandlerForFinetune(handler interface{}) error {
	h, ok := handler.(message.FinetuneHandler)
	if !ok {
		return nil
	}

	f := func(b []byte, m map[string]string) (err error) {
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
	}

	return r.subscribe(r.topics.Finetune, handlerNameCreateFinetune, f)
}

func (r *register) registerHandlerForInference(handler interface{}) error {
	h, ok := handler.(message.InferenceHandler)
	if !ok {
		return nil
	}

	f := func(b []byte, m map[string]string) (err error) {
		body := msgInference{}
		if err = json.Unmarshal(b, &body); err != nil {
			return
		}

		v := spaceappdomain.InferenceIndex{}

		if v.Project.Owner, err = domain.NewAccount(body.ProjectOwner); err != nil {
			return
		}

		v.Id = body.InferenceId
		v.Project.Id = body.ProjectId
		v.LastCommit = body.LastCommit

		info := spaceappdomain.InferenceInfo{
			InferenceIndex: v,
			Requester:      body.User,
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

		default:
			logrus.Warn("unknown action")
			return nil
		}

	}

	return r.subscribe(r.topics.Inference, handlerNameCreateInference, f)
}

func (r *register) subscribe(
	topicName string, handlerName string,
	handler func(b []byte, m map[string]string) (err error),
) error {
	return r.subscriber.SubscribeWithStrategyOfRetry(
		handlerName, handler, []string{topicName}, retryNum,
	)
}

func (r *register) registerHandlerForCloud(handler interface{}) error {
	h, ok := handler.(cloudmsg.CloudMessageHandler)
	if !ok {
		return nil
	}

	f := func(b []byte, m map[string]string) (err error) {
		body := msgPodCreate{}
		if err = json.Unmarshal(b, &body); err != nil {
			return
		}

		user, err := domain.NewAccount(body.User)
		if err != nil {
			return
		}

		cloudCardsNum, err := cloudtypes.NewCloudSpecCardsNum(body.CloudCardsNum)
		if err != nil {
			return
		}

		v := cloudtypes.PodInfo{
			Pod: cloudtypes.Pod{
				Id:      body.PodId,
				CloudId: body.CloudId,
				Owner:   user,
				Image:   body.CloudImage,
			},
			CardsNum: cloudCardsNum,
		}
		if err = v.SetDefaultExpiry(); err != nil {
			return
		}

		return h.HandleEventPodSubscribe(&v)
	}

	return r.subscribe(r.topics.Cloud, handlerNameCreateCloud, f)
}

func (r *register) registerHandlerForReleasingCloud(handler interface{}) error {
	h, ok := handler.(cloudmsg.CloudMessageHandler)
	if !ok {
		return nil
	}

	f := func(b []byte, m map[string]string) (err error) {
		body := ReleasePodMsg{}
		if err = json.Unmarshal(b, &body); err != nil {
			return
		}

		return h.HandleEventPodRelease(body.PodId, body.CloudType)
	}

	return r.subscribe(r.topics.ReleaseCloud, handlerNameReleaseCloud, f)
}
