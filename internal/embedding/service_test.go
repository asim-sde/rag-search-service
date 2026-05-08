package embedding

import "testing"

func TestGenerateEmbedding(t *testing.T) {
	s := NewService()
	emb1 := s.GenerateEmbedding("distributed caching")
	emb2 := s.GenerateEmbedding("distributed caching")
	
	if len(emb1) != 3 {
		t.Errorf("Expected embedding of length 3, got %d", len(emb1))
	}
	
	for i := range emb1 {
		if emb1[i] != emb2[i] {
			t.Errorf("Expected deterministic embeddings, got %v and %v", emb1, emb2)
		}
	}
}
