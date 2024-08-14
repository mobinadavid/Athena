package services

import (
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/repositories"
	"github.com/google/uuid"
)

type IIGPService interface {
	GetList(params *scopes.QueryBuilderModel) ([]*models.IGPModel, int64, error)
	GetByUuid(uuid *uuid.UUID) (*models.IGPModel, error)
	GetForVerification(uuid *uuid.UUID, issuerReferenceNumber string) (*models.IGPModel, error)
	Create(IGP *models.IGPModel) (*models.IGPModel, error)
	Update(IGP *models.IGPModel) (*models.IGPModel, error)
}

type IGPService struct {
	IIGPRepository repositories.IIGPRepository
}

func (service *IGPService) GetList(params *scopes.QueryBuilderModel) ([]*models.IGPModel, int64, error) {
	return service.IIGPRepository.GetList(params)
}

func (service *IGPService) GetByUuid(uuid *uuid.UUID) (*models.IGPModel, error) {
	return service.IIGPRepository.GetByUuid(uuid)
}

func (service *IGPService) Create(IGP *models.IGPModel) (*models.IGPModel, error) {
	return service.IIGPRepository.Create(IGP)
}

func (service *IGPService) Update(IGP *models.IGPModel) (*models.IGPModel, error) {
	return service.IIGPRepository.Update(IGP)
}

func (service *IGPService) GetForVerification(uuid *uuid.UUID, issuerReferenceNumber string) (*models.IGPModel, error) {
	return service.IIGPRepository.GetForVerification(uuid, issuerReferenceNumber)
}
