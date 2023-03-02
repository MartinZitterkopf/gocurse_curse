package curse

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/MartinZitterkopf/gocurse_domain/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(ctx context.Context, curse *domain.Curse) error
		GetAll(ctx context.Context, filters Fillters, limit, offset int) ([]domain.Curse, error)
		GetByID(ctx context.Context, id string) (*domain.Curse, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, name *string, startDate, endDate *time.Time) error
		Count(ctx context.Context, filters Fillters) (int, error)
	}

	repo struct {
		log *log.Logger
		db  *gorm.DB
	}
)

func NewRepo(l *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: l,
		db:  db,
	}
}

func (repo *repo) Create(ctx context.Context, curse *domain.Curse) error {

	if err := repo.db.WithContext(ctx).Create(curse).Error; err != nil {
		repo.log.Printf("error; %v", err)
		return err
	}
	repo.log.Println("curse created with id: ", curse.ID)
	return nil
}

func (repo *repo) GetAll(ctx context.Context, filters Fillters, offset, limit int) ([]domain.Curse, error) {

	var c []domain.Curse

	tx := repo.db.WithContext(ctx).Model(&c)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)

	result := tx.Order("created_at desc").Find(&c)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return nil, result.Error
	}

	return c, nil
}

func (repo *repo) GetByID(ctx context.Context, id string) (*domain.Curse, error) {

	curse := domain.Curse{ID: id}

	if err := repo.db.WithContext(ctx).First(&curse).Error; err != nil {
		repo.log.Println(err)
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound{id}
		}

		return nil, err
	}

	return &curse, nil
}

func (repo *repo) Delete(ctx context.Context, id string) error {

	curse := domain.Curse{ID: id}

	result := repo.db.WithContext(ctx).Delete(&curse)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		repo.log.Printf("curse %s doest't exists", id)
		return ErrNotFound{id}
	}

	return nil
}

func (repo *repo) Update(ctx context.Context, id string, name *string, startDate, endDate *time.Time) error {

	values := make(map[string]interface{})

	if name != nil {
		values["name"] = *name
	}

	if startDate != nil {
		values["start_date"] = *startDate
	}

	if endDate != nil {
		values["end_date"] = *endDate
	}

	result := repo.db.WithContext(ctx).Model(&domain.Curse{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		repo.log.Printf("curse %s doest't exists", id)
		return ErrNotFound{id}
	}

	return nil
}

func applyFilters(tx *gorm.DB, filters Fillters) *gorm.DB {

	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(name) like ?", filters.Name)
	}

	return tx
}

func (repo *repo) Count(ctx context.Context, filters Fillters) (int, error) {

	var count int64
	tx := repo.db.WithContext(ctx).Model(domain.Curse{})

	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
