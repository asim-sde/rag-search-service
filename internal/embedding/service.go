package embedding

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// GenerateEmbedding mocks an external LLM call (e.g. OpenAI text-embedding-ada-002)
func (s *Service) GenerateEmbedding(text string) []float32 {
	length := float32(len(text)%10) / 10.0
	hash := float32(0)
	for _, c := range text {
		hash += float32(c)
	}
	hashVal := float32(int(hash)%100) / 100.0

	return []float32{length, hashVal, 0.5}
}
