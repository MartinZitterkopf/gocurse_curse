package curse

import (
	"context"
	"errors"

	"github.com/MartinZitterkopf/gocurse_library_response/response"
	"github.com/MartinZitterkopf/gocurse_meta/meta"
)

type (
	Controller func(ctx context.Context, request interface{}) (response interface{}, err error)

	Endpoints struct {
		Create  Controller
		GetAll  Controller
		GetByID Controller
		Update  Controller
		Delete  Controller
	}

	CreateReq struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	GetAllReq struct {
		Name  string
		Limit int
		Page  int
	}

	GetByIDReq struct {
		ID string
	}

	UpdateReq struct {
		ID        string
		Name      *string `json:"name"`
		StartDate *string `json:"start_date"`
		EndDate   *string `json:"end_date"`
	}

	DeleteReq struct {
		ID string
	}

	Config struct {
		PageLimDefault string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Create:  makeCreateEndpoint(s),
		GetAll:  makeGetAllEnpoint(s, config),
		GetByID: makeGetByIDEnpoint(s),
		Update:  makeUpdateEnpoint(s),
		Delete:  makeDeleteEnpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateReq)

		if req.Name == "" {
			return nil, response.BadRequest(ErrNameRequired.Error())
		}

		if req.StartDate == "" {
			return nil, response.BadRequest(ErrStartRequired.Error())
		}

		if req.EndDate == "" {
			return nil, response.BadRequest(ErrEndRequired.Error())
		}

		curse, err := s.Create(ctx, req.Name, req.StartDate, req.EndDate)
		if err != nil {

			if err == ErrInvalidStartDate || err == ErrInvalidEndDate || err == ErrEndLesserStart {
				return nil, response.BadRequest(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", curse, nil), nil
	}
}

func makeGetAllEnpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllReq)

		filters := Fillters{
			Name: req.Name,
		}

		count, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, config.PageLimDefault)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		curses, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Ok("success", curses, meta), nil
	}
}

func makeGetByIDEnpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDReq)

		curse, err := s.GetByID(ctx, req.ID)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.Ok("success", curse, nil), nil
	}
}

func makeUpdateEnpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateReq)

		if req.Name != nil && *req.Name == "" {
			return nil, response.BadRequest(ErrNameRequired.Error())
		}

		if req.StartDate != nil && *req.StartDate == "" {
			return nil, response.BadRequest(ErrStartRequired.Error())
		}

		if req.EndDate != nil && *req.EndDate == "" {
			return nil, response.BadRequest(ErrEndRequired.Error())
		}

		err := s.Update(ctx, req.ID, req.Name, req.StartDate, req.EndDate)
		if err != nil {
			if err == ErrInvalidStartDate || err == ErrInvalidEndDate || err == ErrEndLesserStart {
				return nil, response.BadRequest(err.Error())
			}

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.Ok("success", nil, nil), nil
	}
}

func makeDeleteEnpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteReq)

		err := s.Delete(ctx, req.ID)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.Ok("success", nil, nil), nil
	}
}
