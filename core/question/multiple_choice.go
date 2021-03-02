package question

type SingleChoice struct {
	Id      string  `json:"id"`
	Content string  `json:"content"`
	Mark    float32 `json:"mark"`
}

type MultipleChoice struct {
	Prompt  string         `json:"prompt"`
	Choices []SingleChoice `json:"choices"`
}

func (q *MultipleChoice) CheckAnswer(answer string) float32 {
	for _, choice := range q.Choices {
		if choice.Id == answer {
			return choice.Mark
		}
	}
	return 0
}
