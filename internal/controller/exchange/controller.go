package exchange

import (
	"errors"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	modelexchange "github.com/andrqxa/artshare/internal/model/exchange"
	"github.com/andrqxa/artshare/internal/model/pagination"
	exchangerepository "github.com/andrqxa/artshare/internal/repository/exchange"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type Controller struct {
	repository *exchangerepository.Repository
}

func NewController() *Controller {
	return &Controller{
		repository: exchangerepository.NewRepository(),
	}
}

func (c *Controller) ListExchangeRequests(filter modelexchange.ListFilter) ([]modelexchange.Summary, pagination.Meta) {
	return c.repository.List(filter)
}

func (c *Controller) CreateExchangeRequest(input modelexchange.CreateInput) (modelexchange.Detail, error) {
	request, err := c.repository.Create(input)
	return request, controllerError(err)
}

func (c *Controller) GetExchangeRequestByID(id string) (modelexchange.Detail, error) {
	request, err := c.repository.FindByID(id)
	return request, controllerError(err)
}

func (c *Controller) UpdateExchangeRequestStatus(input modelexchange.UpdateStatusInput) (modelexchange.Detail, error) {
	current, err := c.repository.FindByID(input.ID)
	if err != nil {
		return modelexchange.Detail{}, controllerError(err)
	}
	if current.Status != modelexchange.StatusPending {
		return modelexchange.Detail{}, apperror.ErrConflict
	}

	request, err := c.repository.UpdateStatus(input)
	return request, controllerError(err)
}

func controllerError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, repoerror.ErrConflict):
		return apperror.ErrConflict
	case errors.Is(err, repoerror.ErrNotFound):
		return apperror.ErrNotFound
	default:
		return err
	}
}
