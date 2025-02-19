package driven

import metadataModel "github.com/karaMuha/go-movie/metadata/pkg"

type IMessageProducer interface {
	PublishMetadataSubmittedEvent(event metadataModel.MetadataEvent) error
}
