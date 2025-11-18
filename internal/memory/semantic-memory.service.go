package memory

import "memory/internal/persistence"

type SemanticService struct {
	vectorMemoryRepo persistence.VectorMemoryRepoInterface
}
