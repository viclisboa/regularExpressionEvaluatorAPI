package repository

import (
	"errors"
	"fmt"
	_gorm "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/viclisboa/regularExpressionEvaluatorAPI/model"
)

type ExpressionInterface interface {
	GetAllExpressions() ([]model.Expression, error)
	GetExpressionById(expressionId int) (model.Expression, error)
	CreateExpression(definition string) error
	SaveExpression(expressionId int, definition string) error
}

var _ ExpressionInterface = (*Repository)(nil)

type Repository struct {
	db *_gorm.DB
}

func NewRepository(database, connectionString string) (Repository, error) {
	db, err := _gorm.Open(database, connectionString)
	if err != nil {
		panic(err)
	}

	db.DropTableIfExists("expression")

	return Repository{
		db: db,
	}, nil
}

func (r *Repository) GetAllExpressions() ([]model.Expression, error) {
	var expressions []model.Expression
	result := r.db.Model(model.Expression{}).
		Find(&expressions)

	if result.Error != nil {
		fmt.Println("Failed to execute query", "err", result.Error)
		return nil, errors.New("failed to execute query")
	}

	return expressions, nil
}

func (r *Repository) GetExpressionById(expressionId int) (model.Expression, error) {
	var expression model.Expression
	result := r.db.First(&expression, []int{expressionId})

	if result.Error != nil {
		fmt.Println("Failed to execute query", "err", result.Error)
		return model.Expression{}, errors.New("failed to execute query")
	}

	return expression, nil
}

func (r *Repository) CreateExpression(definition string) error {
	expression := model.Expression{
		Definition: definition,
	}

	createResponse := r.db.Create(expression)

	if createResponse.Error != nil {
		fmt.Println(fmt.Sprintf("error while trying to create user, err: %s", createResponse.Error.Error()))
		return createResponse.Error
	}

	fmt.Println("user created successfully")
	return nil
}

func (r *Repository) SaveExpression(expressionId int, definition string) error {
	var expression model.Expression
	result := r.db.Model(&expression).Where("id = ?", expressionId).Update("definition", definition)

	if result.Error != nil {
		fmt.Println("Failed to execute query", "err", result.Error)
		return errors.New("failed to execute query")
	}
	return nil
}
