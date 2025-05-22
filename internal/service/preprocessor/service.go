package preprocessor

import (
	"strings"

	model "agregator/preprocessor/internal/model/kafka"

	"agregator/preprocessor/internal/interfaces"
)

type Preprocessor struct {
	input  chan model.Item
	output chan model.Item
	logger interfaces.Logger
}

func New(bufferSize int, logger interfaces.Logger) *Preprocessor {
	return &Preprocessor{
		input:  make(chan model.Item, bufferSize),
		output: make(chan model.Item, bufferSize),
		logger: logger,
	}
}

func (p *Preprocessor) Start() {
	defer close(p.output)
	for model := range p.input {
		p.logger.Debug("Processing:", model.Link)
		processedModel := p.Process(model)
		p.output <- processedModel
	}
}

func (p *Preprocessor) Process(text model.Item) model.Item {
	l := text.Description
	text.Title = p.clearHTML(text.Title, false)
	text.Description = p.clearHTML(text.Description, false)
	text.FullText = p.clearHTML(text.FullText, true)
	if !strings.Contains(text.FullText, "Доступно только описание") {
		if p.cosineSimilarity(text.FullText, text.Description) > 0.9 && len(text.Description) > 256 {
			text.Description = text.Description[:256]
		}
		text.Description = strings.TrimSpace(text.Description) + "..."
	} else if len(l) > len(text.Description) {
		text.FullText = p.clearHTML(l, true)
		if len(text.Description) > 256 {
			text.Description = text.Description[:256]
		}
	}
	return text
}

func (p *Preprocessor) Input() chan<- model.Item {
	return p.input
}

func (p *Preprocessor) Output() <-chan model.Item {
	return p.output
}
