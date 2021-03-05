package question

type Question interface {
	CheckAnswer(answer string) float32
	GetMaxScore() float32
}
