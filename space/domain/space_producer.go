package domain

type SpaceEventProducer interface {
	SendDeletedEvent(*DeleteSpaceEvent) error
}
