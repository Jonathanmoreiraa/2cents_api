package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
	interfaces "github.com/jonathanmoreiraa/2cents/internal/domain/repository"
	database "github.com/jonathanmoreiraa/2cents/internal/infra/database/interface"
	"github.com/shopspring/decimal"

	"gorm.io/gorm"
)

type expenseDatabase struct {
	DB *gorm.DB
}

func NewExpenseRepository(Database database.DatabaseProvider) interfaces.ExpenseRepository {
	return &expenseDatabase{DB: Database.GetDatabase()}
}

func (database *expenseDatabase) Create(ctx context.Context, expense entity.Expense, multiplePayments bool, numInstallments int, paymentDay int) (entity.Expense, error) {
	if multiplePayments {
		for i := 0; i < numInstallments; i++ {
			newExpense := expense
			newExpense.DueDate = &time.Time{}
			*newExpense.DueDate = time.Date(expense.DueDate.Year(), expense.DueDate.Month()+time.Month(i), paymentDay, 0, 0, 0, 0, expense.DueDate.Location())
			err := database.DB.Create(&newExpense).Error
			if err != nil {
				return newExpense, err
			}
		}
		return expense, nil
	}

	err := database.DB.Create(&expense).Error
	return expense, err
}

func (database *expenseDatabase) FindAll(ctx context.Context, userId int) ([]entity.Expense, error) {
	var expenses []entity.Expense

	err := database.DB.
		Where("user_id = ?", userId).
		Where("deleted_at IS NULL").
		Find(&expenses).Error
	return expenses, err
}

func (database *expenseDatabase) FindByID(ctx context.Context, id int) (entity.Expense, error) {
	var expense entity.Expense

	err := database.DB.First(&expense, id).Error
	return expense, err
}

func (database *expenseDatabase) FindByFilter(ctx context.Context, filters map[string]any) ([]entity.Expense, error) {
	var expenses []entity.Expense

	query := database.DB.
		Where("user_id = ?", filters["user_id"]).
		Where("deleted_at IS NULL")

	if description, ok := filters["description"]; ok && description != "" {
		query = query.Where("description LIKE ?", "%"+description.(string)+"%")
	}
	if dateStart, ok := filters["date_start"]; ok && dateStart != "" {
		query = query.Where("due_date >= ?", dateStart)
	}
	if dateEnd, ok := filters["date_end"]; ok && dateEnd != "" {
		query = query.Where("due_date <= ?", dateEnd)
	}
	if min, ok := filters["min"]; ok && !min.(decimal.Decimal).IsZero() {
		query = query.Where("value >= ?", min)
	}
	if max, ok := filters["max"]; ok && !max.(decimal.Decimal).IsZero() {
		query = query.Where("value <= ?", max)
	}
	if categories, ok := filters["categories"]; ok && len(categories.([]int)) > 0 {
		query = query.Where("category_id IN ?", categories)
	}
	if category, ok := filters["category_id"]; ok {
		query = query.Where("category_id = ?", category)
	}
	if status, ok := filters["status"]; ok && status != nil {
		statusStruct := status.(struct {
			Pending bool `json:"pending"`
			Paid    bool `json:"paid"`
			Overdue bool `json:"overdue"`
			DueSoon bool `json:"due_soon"`
		})
		conditions := []string{}

		if statusStruct.Pending {
			conditions = append(conditions, "(paid = 0 AND due_date > '"+time.Now().Format("2006-01-02")+"')")
		}
		if statusStruct.Paid {
			conditions = append(conditions, "paid = 1")
		}
		if statusStruct.Overdue {
			conditions = append(conditions, "(paid = 0 AND due_date < '"+time.Now().Format("2006-01-02")+"')")
		}
		if statusStruct.DueSoon {
			soonDate := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
			conditions = append(conditions, "(paid = 0 AND due_date BETWEEN '"+time.Now().Format("2006-01-02")+"' AND '"+soonDate+"')")
		}

		if len(conditions) > 0 {
			query = query.Where(strings.Join(conditions, " OR "))
		}
	}

	if filters["date_start"] != nil || filters["date_end"] != nil {
		query = query.Order("due_date ASC")
	}

	err := query.Find(&expenses).Error
	fmt.Println(query.Debug().Find(&expenses))

	return expenses, err
}

func (database *expenseDatabase) Update(ctx context.Context, expense entity.Expense) error {
	err := database.DB.Model(&expense).Updates(map[string]interface{}{
		"description": expense.Description,
		"due_date":    expense.DueDate,
		"paid":        expense.Paid,
		"value":       expense.Value,
		"category_id": expense.CategoryID,
		"updated_at":  time.Now(),
	}).Error
	return err
}

func (database *expenseDatabase) Delete(ctx context.Context, expense entity.Expense) error {
	err := database.DB.Delete(&expense).Error
	return err
}
