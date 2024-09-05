/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package messageadapter

import (
	"fmt"

	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/organization/domain/message"
)

// MessageAdapter creates a new messageAdapter instance with the given Topics.
func MessageAdapter(topic *Topics) *messageAdapter {
	return &messageAdapter{topics: *topic}
}

type messageAdapter struct {
	topics Topics
}

// SendComputilityUserJoinEvent sends a User Join event message to the corresponding topic.
func (p *messageAdapter) SendComputilityUserJoinEvent(e message.EventMessage) error {
	return send(p.topics.ComputilityUserJoined, e)
}

// SendComputilityUserRemoveEvent sends a User Remove event message to the corresponding topic.
func (p *messageAdapter) SendComputilityUserRemoveEvent(e message.EventMessage) error {
	return send(p.topics.ComputilityUserRemoved, e)
}

// SendComputilityOrgDeleteEvent sends a Org Delete event message to the corresponding topic.
func (p *messageAdapter) SendComputilityOrgDeleteEvent(e message.EventMessage) error {
	return send(p.topics.ComputilityOrgDeleted, e)
}

func (p *messageAdapter) SendOrgInviteEvent(e message.EventMessage) error {
	return send(p.topics.OrgInvite, e)
}

func (p *messageAdapter) SendOrgRequestEvent(e message.EventMessage) error {
	return send(p.topics.OrgRequest, e)
}

func (p *messageAdapter) SendOrgRequestRejectEvent(e message.EventMessage) error {
	return send(p.topics.OrgRequestReject, e)
}

func send(topic string, v message.EventMessage) error {
	body, err := v.Message()
	if err != nil {
		return fmt.Errorf("send msg topic:%s err: %w", topic, err)
	}

	err = kfklib.Publish(topic, nil, body)
	if err != nil {
		logrus.Errorf("send msg topic:%s err: %s", topic, err)
		err = fmt.Errorf("send publish topic:%s err: %w", topic, err)
	}

	logrus.Infof("send msg topic:%s succeed", topic)

	return err
}
