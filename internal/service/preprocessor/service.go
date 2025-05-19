package preprocessor

import (
	"log"
	"strings"

	model "agregator/preprocessor/internal/model/kafka"
)

type Preprocessor struct {
	input  chan model.Item
	output chan model.Item
}

func New(bufferSize int) *Preprocessor {
	return &Preprocessor{
		input:  make(chan model.Item, bufferSize),
		output: make(chan model.Item, bufferSize),
	}
}

func (p *Preprocessor) Start() {
	defer close(p.output)
	for model := range p.input {
		log.Default().Println("Processing:", model.Link)
		processedModel := p.Process(model)
		p.output <- processedModel
	}
}

func (p *Preprocessor) Process(text model.Item) model.Item {
	text.Title = p.clearHTML(text.Title, false)
	text.Description = p.clearHTML(text.Description, false)
	text.FullText = p.clearHTML(text.FullText, true)
	if text.Description == p.clearHTML(text.FullText, false) && !strings.Contains(text.FullText, "Доступно только описание") {
		if len(text.Description) > 256 {
			text.Description = text.Description[:256]
		}
		text.Description = strings.TrimSpace(text.Description) + "..."
	}
	return text
}

func (p *Preprocessor) Input() chan<- model.Item {
	return p.input
}

func (p *Preprocessor) Output() <-chan model.Item {
	return p.output
}
