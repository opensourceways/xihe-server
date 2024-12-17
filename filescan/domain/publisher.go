package domain

type ModerationEventPublisher interface {
	Publish(ModerationEvent) error
}
