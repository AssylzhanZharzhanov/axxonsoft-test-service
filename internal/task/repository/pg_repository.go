package repository

import (
	"context"
	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/pkg/database/postgres"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository - create a new repository
func NewRepository(db *gorm.DB) domain.TaskRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, task *domain.Task) (domain.TaskID, error) {
	var (
		db = r.db
	)

	err := db.WithContext(ctx).Table(domain.TableName).Create(&task).Error
	if err != nil {
		return 0, err
	}

	return task.ID, nil
}

func (r *repository) Get(ctx context.Context, taskID domain.TaskID) (*domain.Task, error) {
	var (
		db   = r.db
		task *domain.Task
	)

	err := db.WithContext(ctx).Table(domain.TableName).First(&task, taskID).Error
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *repository) List(ctx context.Context, criteria domain.TaskSearchCriteria) ([]*domain.Task, domain.Total, error) {
	var (
		db         = r.db
		tasks      []*domain.Task
		totalCount int64
	)

	err := db.WithContext(ctx).Scopes(postgres.Paginate(criteria.Page)).Table(domain.TableName).Find(&tasks).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, domain.Total(totalCount), nil
}

func (r *repository) Update(ctx context.Context, task *domain.Task) error {
	var (
		db = r.db
	)

	err := db.WithContext(ctx).Table(domain.TableName).Where("id = ?", task.ID).Updates(&task).Error
	if err != nil {
		return err
	}

	return nil
}
