package services

import (
	"context"
	"fmt"

	"github.com/lantoniomiranda/shitreader/internal/store"
)

type AssociationService struct {
	associationStore store.AssociationStore
}

func NewAssociationService(associtationStore store.AssociationStore) *AssociationService {
	return &AssociationService{
		associationStore: associtationStore,
	}
}

func (s *AssociationService) Associate() error {
	ctx := context.Background()
	err := s.associationStore.AssociateRecordsFields(ctx)
	if err != nil {
		return fmt.Errorf("error doing associations: %w", err)
	}
	return nil
}

func (s *AssociationService) AssociateRecordTypes(filePath string, sheetName string) error {
	ctx := context.Background()
	err := s.associationStore.AssociateRecordsRecordTypes(ctx, filePath, sheetName)
	if err != nil {
		return fmt.Errorf("error associating record types: %w", err)
	}
	return nil
}

func (s *AssociationService) AssociateSteps(filePath string, sheetName string) error {
	ctx := context.Background()
	err := s.associationStore.AssociateStepsHeaderTypesAndRecords(ctx, filePath, sheetName)
	if err != nil {
		return fmt.Errorf("error associating steps: %w", err)
	}
	return nil
}
