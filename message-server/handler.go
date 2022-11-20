package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
)

var _ message.EventHandler = (*handler)(nil)

const sleepTime = 100 * time.Millisecond

type likeHanler interface {
	AddLike(*domain.ResourceIndex) error
	RemoveLike(*domain.ResourceIndex) error
}

type relatedResourceHanler struct {
	Add    func(*domain.ReverselyRelatedResourceInfo) error
	Remove func(*domain.ReverselyRelatedResourceInfo) error
}

type handler struct {
	log *logrus.Entry

	maxRetry         int
	trainingEndpoint string

	user      app.UserService
	model     app.ModelMessageService
	dataset   app.DatasetMessageService
	project   app.ProjectMessageService
	training  app.TrainingService
	evaluate  app.EvaluateMessageService
	inference app.InferenceMessageService
}

func (h *handler) HandleEventAddRelatedResource(info *message.RelatedResource) error {
	rh := h.getHandlerForEventRelatedResource(info)
	if rh.Add == nil {
		return errors.New("unknown Reversely Related Resource")
	}

	data := h.getParameterForEventRelatedResource(info)

	return h.do(func(bool) (err error) {
		return rh.Add(&data)
	})
}

func (h *handler) HandleEventRemoveRelatedResource(info *message.RelatedResource) error {
	rh := h.getHandlerForEventRelatedResource(info)
	if rh.Remove == nil {
		return errors.New("unknown Reversely Related Resource")
	}

	data := h.getParameterForEventRelatedResource(info)

	return h.do(func(bool) (err error) {
		return rh.Remove(&data)
	})
}

func (h *handler) getParameterForEventRelatedResource(
	info *message.RelatedResource,
) domain.ReverselyRelatedResourceInfo {
	return domain.ReverselyRelatedResourceInfo{
		Promoter: &info.Promoter.ResourceIndex,
		Resource: &info.Resource.ResourceIndex,
	}
}

func (h *handler) getHandlerForEventRelatedResource(
	info *message.RelatedResource,
) (v relatedResourceHanler) {
	pt := info.Promoter.Type.ResourceType()

	model := domain.ResourceTypeModel.ResourceType()
	project := domain.ResourceTypeProject.ResourceType()
	dataset := domain.ResourceTypeDataset.ResourceType()

	switch info.Resource.Type.ResourceType() {
	case dataset:
		switch pt {
		case model:
			v.Add = h.dataset.AddRelatedModel
			v.Remove = h.dataset.RemoveRelatedModel

		case project:
			v.Add = h.dataset.AddRelatedProject
			v.Remove = h.dataset.RemoveRelatedProject
		}

	case model:
		switch pt {
		case project:
			v.Add = h.model.AddRelatedProject
			v.Remove = h.model.RemoveRelatedProject

		case dataset:
			v.Add = h.model.AddRelatedDataset
			v.Remove = h.model.RemoveRelatedDataset
		}

	case project:
		switch pt {
		case model:
			v.Add = h.project.AddRelatedModel
			v.Remove = h.project.RemoveRelatedModel

		case dataset:
			v.Add = h.project.AddRelatedDataset
			v.Remove = h.project.RemoveRelatedDataset
		}
	}

	return
}

func (h *handler) HandleEventAddFollowing(f *domain.FollowerInfo) error {
	return h.do(func(bool) (err error) {
		if err = h.user.AddFollower(f); err == nil {
			return
		}

		if _, ok := err.(repository.ErrorDuplicateCreating); ok {
			err = nil
		}

		return
	})
}

func (h *handler) HandleEventRemoveFollowing(f *domain.FollowerInfo) (err error) {
	return h.do(func(bool) error {
		return h.user.RemoveFollower(f)
	})
}

func (h *handler) HandleEventAddLike(obj *domain.ResourceObject) error {
	lh := h.getHandlerForEventLike(obj.Type)

	return h.handleEventLike(obj, "adding", lh.AddLike)
}

func (h *handler) HandleEventRemoveLike(obj *domain.ResourceObject) (err error) {
	lh := h.getHandlerForEventLike(obj.Type)

	return h.handleEventLike(obj, "removing", lh.RemoveLike)
}

func (h *handler) handleEventLike(
	obj *domain.ResourceObject, op string,
	f func(*domain.ResourceIndex) error,
) (err error) {
	return h.do(func(bool) (err error) {
		if err = f(&obj.ResourceIndex); err != nil {
			if isResourceNotExists(err) {
				h.log.Errorf(
					"handle event of %s like for owner:%s, rid:%s, err:%v",
					op, obj.Owner.Account(), obj.Id, err,
				)

				err = nil
			}
		}

		return
	})
}

func (h *handler) getHandlerForEventLike(t domain.ResourceType) likeHanler {
	switch t.ResourceType() {
	case domain.ResourceTypeProject.ResourceType():
		return h.project

	case domain.ResourceTypeDataset.ResourceType():
		return h.dataset

	case domain.ResourceTypeModel.ResourceType():
		return h.model
	}

	return nil
}

func (h *handler) HandleEventFork(index *domain.ResourceIndex) error {
	return h.do(func(bool) (err error) {
		if err = h.project.IncreaseFork(index); err != nil {
			if isResourceNotExists(err) {
				h.log.Errorf(
					"handle event of fork for owner:%s, rid:%s, err:%v",
					index.Owner.Account(), index.Id, err,
				)

				err = nil
			}
		}

		return
	})
}

func (h *handler) HandleEventCreateTraining(info *domain.TrainingIndex) error {
	// wait for the sync of model and dataset
	time.Sleep(10 * time.Second)

	return h.retry(
		func(lastChance bool) error {
			retry, err := h.training.CreateTrainingJob(
				info, h.trainingEndpoint, lastChance,
			)
			if err != nil {
				h.log.Errorf(
					"handle training(%s/%s/%s) failed, err:%s",
					info.Project.Owner.Account(), info.Project.Id,
					info.TrainingId, err.Error(),
				)

				if !retry {
					return nil
				}
			}

			return err
		},
		10*time.Second,
	)
}

func (h *handler) HandleEventCreateInference(info *domain.InferenceInfo) error {
	return h.do(func(bool) error {
		err := h.inference.CreateInferenceInstance(info)
		if err != nil {
			h.log.Error(err)
		}

		return err
	})
}

func (h *handler) HandleEventExtendInferenceSurvivalTime(info *message.InferenceExtendInfo) error {
	return h.do(func(bool) error {
		err := h.inference.ExtendSurvivalTime(info)
		if err != nil {
			h.log.Error(err)
		}

		return err
	})
}

func (h *handler) HandleEventCreateEvaluate(info *message.EvaluateInfo) error {
	return h.do(func(bool) error {
		err := h.evaluate.CreateEvaluateInstance(info)
		if err != nil {
			h.log.Error(err)
		}

		return err
	})
}

func (h *handler) do(f func(bool) error) (err error) {
	return h.retry(f, sleepTime)
}

func (h *handler) retry(f func(bool) error, interval time.Duration) (err error) {
	n := h.maxRetry - 1

	if err = f(n <= 0); err == nil || n <= 0 {
		return
	}

	for i := 1; i < n; i++ {
		time.Sleep(interval)

		if err = f(false); err == nil {
			return
		}
	}

	time.Sleep(interval)

	return f(true)
}

func (h *handler) errMaxRetry(err error) error {
	return fmt.Errorf("exceed max retry num, last err:%v", err)
}

func isResourceNotExists(err error) bool {
	_, ok := err.(repository.ErrorResourceNotExists)

	return ok
}
