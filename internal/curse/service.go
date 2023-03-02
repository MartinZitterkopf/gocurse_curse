package curse

import (
	"context"
	"log"
	"time"

	"github.com/MartinZitterkopf/gocurse_domain/domain"
)

type (
	Service interface {
		Create(ctx context.Context, name, startDate, endDate string) (*domain.Curse, error)
		GetAll(ctx context.Context, filters Fillters, offset, limit int) ([]domain.Curse, error)
		GetByID(ctx context.Context, id string) (*domain.Curse, error)
		Update(ctx context.Context, id string, name, startDate, endDate *string) error
		Delete(ctx context.Context, id string) error
		Count(ctx context.Context, filters Fillters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}

	Fillters struct {
		Name string
	}
)

func NewService(l *log.Logger, r Repository) Service {
	return &service{
		log:  l,
		repo: r,
	}
}

func (s service) Create(ctx context.Context, name, startDate, endDate string) (*domain.Curse, error) {

	s.log.Println("create user service")
	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	if startDateParsed.After(endDateParsed) {
		s.log.Println(ErrEndLesserStart)
		return nil, ErrEndLesserStart
	}

	curse := &domain.Curse{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}

	if err := s.repo.Create(ctx, curse); err != nil {
		s.log.Println(err)
		return nil, err
	}

	return curse, nil
}

func (s service) GetAll(ctx context.Context, filters Fillters, offset, limit int) ([]domain.Curse, error) {

	curses, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	return curses, nil
}

func (s service) GetByID(ctx context.Context, id string) (*domain.Curse, error) {

	curse, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	return curse, nil
}

func (s service) Update(ctx context.Context, id string, name, startDate, endDate *string) error {

	var startDateParsed, endDateParsed *time.Time

	curse, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if startDate != nil {
		date, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			s.log.Println(err)
			return ErrInvalidStartDate
		}

		if date.After(curse.EndDate) {
			s.log.Println(ErrEndLesserStart)
			return ErrEndLesserStart
		}

		startDateParsed = &date
	}

	if endDate != nil {
		date, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			s.log.Println(err)
			return ErrInvalidEndDate
		}

		if curse.StartDate.After(date) {
			s.log.Println(ErrEndLesserStart)
			return ErrEndLesserStart
		}

		endDateParsed = &date
	}

	if err := s.repo.Update(ctx, id, name, startDateParsed, endDateParsed); err != nil {
		return err
	}

	return nil
}

func (s service) Delete(ctx context.Context, id string) error {

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s service) Count(ctx context.Context, filters Fillters) (int, error) {
	return s.repo.Count(ctx, filters)
}
