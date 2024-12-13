package infrastructure

import (
	"errors"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/filescan/domain"
	"github.com/opensourceways/xihe-server/filescan/domain/primitive"
)

type AuditTopics struct {
	CreateDocModerationTask     string `json:"create_doc_moderation_task"     required:"true"`
	CreateReadmeModerationTask  string `json:"create_readme_moderation_task"  required:"true"`
	CreatePictureModerationTask string `json:"create_picture_moderation_task" required:"true"`
	CreateVideoModerationTask   string `json:"create_video_moderation_task"   required:"true"`
	CreateAudioModerationTask   string `json:"create_audio_moderation_task"   required:"true"`
}

type AuditSupportedFileType struct {
	MarkDown []string `json:"markdown"`
	Document []string `json:"document"`
	Audio    []string `json:"audio"`
	Video    []string `json:"video"`
	Image    []string `json:"image"`
	Code     []string `json:"code"`
}

type moderationEventPublisher struct {
	topics      AuditTopics
	publisher   message.Publisher
	markdownExt sets.Set[string]
	documentExt sets.Set[string]
	audioExt    sets.Set[string]
	videoExt    sets.Set[string]
	imageExt    sets.Set[string]
	codeExt     sets.Set[string]
}

type AuditPublisherConfig struct {
	Topics   AuditTopics            `json:"topics"`
	FileType AuditSupportedFileType `json:"file_type"`
}

func NewModerationEventPublisher(config *AuditPublisherConfig, p message.Publisher) domain.ModerationEventPublisher {
	return moderationEventPublisher{
		topics:      config.Topics,
		publisher:   p,
		markdownExt: sets.New[string](config.FileType.MarkDown...),
		documentExt: sets.New[string](config.FileType.Document...),
		audioExt:    sets.New[string](config.FileType.Audio...),
		videoExt:    sets.New[string](config.FileType.Video...),
		imageExt:    sets.New[string](config.FileType.Image...),
		codeExt:     sets.New[string](config.FileType.Code...),
	}
}

func (m moderationEventPublisher) Publish(event domain.ModerationEvent) error {
	fileType := m.getFileType(event.File)

	switch fileType.Value() {
	case primitive.MarkdownFileType.Value(), primitive.CodeFileType.Value():
		return m.publisher.Publish(m.topics.CreateReadmeModerationTask, event, nil)
	case primitive.DocumentFileType.Value():
		return m.publisher.Publish(m.topics.CreateDocModerationTask, event, nil)
	case primitive.AudioFileType.Value():
		return m.publisher.Publish(m.topics.CreateAudioModerationTask, event, nil)
	case primitive.VideoFileType.Value():
		return m.publisher.Publish(m.topics.CreateVideoModerationTask, event, nil)
	case primitive.ImageFileType.Value():
		return m.publisher.Publish(m.topics.CreatePictureModerationTask, event, nil)
	}

	return errors.New("unsupported file type")
}

func (m moderationEventPublisher) getFileType(file string) primitive.FileType {
	ext := strings.ToLower(filepath.Ext(file))

	if m.markdownExt.Has(ext) {
		return primitive.MarkdownFileType
	} else if m.documentExt.Has(ext) {
		return primitive.DocumentFileType
	} else if m.audioExt.Has(ext) {
		return primitive.AudioFileType
	} else if m.videoExt.Has(ext) {
		return primitive.VideoFileType
	} else if m.imageExt.Has(ext) {
		return primitive.ImageFileType
	} else if m.codeExt.Has(ext) {
		return primitive.CodeFileType
	}

	return primitive.UnKnownFileType
}
