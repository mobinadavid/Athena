package scopes

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	MaxPageSize     = 1000
	DefaultPage     = 1
	DefaultPageSize = 10
)

// BuilderModel encapsulates filtering, sorting, and pagination options for database queries
type BuilderModel struct {
	Page                  uint
	PageSize              uint
	SortBy                string
	SortOrder             string
	PrioritySorts         []PrioritySort
	Filters               map[string]interface{}
	InFilters             map[string][]interface{}
	OrFilters             map[string][]interface{}
	OrFiltersWithAndConds [][]Condition
	Likes                 map[string]interface{}
	OrLikes               map[string]interface{}
	Joins                 []string
	Relations             []string
	timeFilters           []TimeFilter
	TotalItems            int64
	M2MFilters            map[string][]interface{}
	M2MOrLikes            map[string]interface{}
	GlobalSearchCols      []string
	GlobalSearchVal       string
}

// PaginateModel encapsulates the result of a paginated query
type PaginateModel struct {
	PageSize    uint        `json:"page_size"`
	CurrentPage uint        `json:"current_page"`
	TotalPages  int64       `json:"total_pages"`
	TotalItems  int64       `json:"total_items"`
	Items       interface{} `json:"items"`
}

// TimeFilter Adding a time based filter
// TableName is the field name in model (considering relations)
// Only one of the time based fields must be non-nil
type TimeFilter struct {
	TableName string
	Type      TimeFilterType
	Value     *time.Time
}

type TimeFilterType string

const (
	CreatedBefore TimeFilterType = "created_before"
	CreatedAfter  TimeFilterType = "created_after"
	UpdatedBefore TimeFilterType = "updated_before"
	UpdatedAfter  TimeFilterType = "updated_after"
	PublishBefore TimeFilterType = "publish_before"
	PublishAfter  TimeFilterType = "publish_after"
	DueDateBefore TimeFilterType = "due_date_before"
	DueDateAfter  TimeFilterType = "due_date_after"
)

type CompareOperator uint16

const (
	EqualOperator CompareOperator = iota
	NotEqualOperator
)

type Condition struct {
	Operand  string
	Value    interface{}
	Operator CompareOperator
}

type PriorityCase struct {
	Conditions []Condition
	Priority   int
}

// PrioritySort sorts the fields by priority
type PrioritySort struct {
	Cases []PriorityCase
}

type PreloadWithCondition struct {
	FieldName  string
	Conditions *[]string
}

// validatePaginationParams checks and sets default values for pagination parameters
func (bm *BuilderModel) validatePaginationParams() {
	if bm.Page == 0 {
		bm.Page = DefaultPage
	}
	//if bm.PageSize > MaxPageSize {
	//	bm.PageSize = MaxPageSize
	//} else
	if bm.PageSize == 0 {
		bm.PageSize = DefaultPageSize
	}
}

func (bm *BuilderModel) AddTimeFilter(tf *TimeFilter) {
	if tf.Value == nil {
		return
	}

	bm.timeFilters = append(bm.timeFilters, *tf)
}

func (bm *BuilderModel) ExistsTimeFilter(tf *TimeFilter) bool {
	for _, filter := range bm.timeFilters {
		if filter.TableName == tf.TableName && filter.Type == tf.Type && filter.Value != nil {
			return true
		}
	}
	return false
}

// applySorting applies sorting to the query based on SortBy and SortOrder
func (bm *BuilderModel) applySorting(db *gorm.DB) *gorm.DB {
	// Apply each PrioritySort
	if bm.PrioritySorts != nil {
		for _, ps := range bm.PrioritySorts {
			caseParts := ""
			for _, c := range ps.Cases {
				conds := ""
				i := 0
				for _, cond := range c.Conditions {
					if i > 0 {
						conds += " AND "
					}
					switch cond.Operator {
					case EqualOperator:
						conds += fmt.Sprintf("\"%s\" = '%v'", cond.Operand, cond.Value)
					case NotEqualOperator:
						conds += fmt.Sprintf("\"%s\" != '%v'", cond.Operand, cond.Value)
					default:
						continue
					}
					i++
				}
				caseParts += fmt.Sprintf("WHEN %s THEN %d ", conds, c.Priority)
			}

			caseSQL := fmt.Sprintf("(CASE %s ELSE %d END)", caseParts, len(ps.Cases)+1)
			db = db.Order(caseSQL)
		}
	}

	// Then append SortBy
	if bm.SortBy != "" {
		// If SortBy looks like a raw SQL function/expression, just inject it
		if strings.Contains(bm.SortBy, "(") || strings.Contains(bm.SortBy, " ") {
			db = db.Order(bm.SortBy + " " + bm.SortOrder)
		} else if strings.Contains(bm.SortBy, "/") {
			db = db.Order(bm.SortBy + " " + bm.SortOrder + " NULLS LAST")
		} else if strings.Contains(bm.SortBy, ".") {
			// Handle relation.field
			parts := strings.Split(bm.SortBy, ".")
			if len(parts) == 2 {
				tableName := parts[0]
				fieldName := parts[1]
				db = db.Order(fmt.Sprintf("\"%s\".\"%s\" %s NULLS LAST", tableName, fieldName, bm.SortOrder))
			}
		} else {
			db = db.Order(fmt.Sprintf("\"%s\" %s", bm.SortBy, bm.SortOrder))
		}
	}

	return db
}

// applyJoins applies joins to the query
func (bm *BuilderModel) applyJoins(db *gorm.DB) *gorm.DB {
	for _, join := range bm.Joins {
		db = db.Joins(join)
	}
	return db
}

// applyFilters applies dynamic filtering based on the Filters map in BuilderModel
func (bm *BuilderModel) applyFilters(db *gorm.DB) *gorm.DB {
	for key, list := range bm.InFilters {
		if strings.Contains(key, ".") {
			parts := strings.Split(key, ".")
			if len(parts) == 2 {
				tableName := parts[0]
				fieldName := parts[1]

				db = db.Where(fmt.Sprintf("\"%s\".\"%s\" IN ?", tableName, fieldName), list)
			}
		} else {
			db = db.Where(fmt.Sprintf("%s IN ?", key), list)
		}
	}

	for key, value := range bm.Filters {
		if strings.Contains(key, ".") {
			parts := strings.Split(key, ".")
			if len(parts) == 2 {
				tableName := parts[0]
				fieldName := parts[1]

				db = db.Where(fmt.Sprintf(`"%s"."%s" = ?`, tableName, fieldName), value)

			}
		} else {
			db = db.Where(clause.Expr{SQL: key + " = ?", Vars: []interface{}{value}})
		}
	}

	if len(bm.timeFilters) > 0 {
		for _, tf := range bm.timeFilters {
			if tf.Value == nil {
				continue
			}

			switch tf.Type {
			case CreatedAfter:
				db = db.Where(fmt.Sprintf(`"%s"."created_at" >= ?`, tf.TableName), tf.Value)
			case CreatedBefore:
				db = db.Where(fmt.Sprintf(`"%s"."created_at" <= ?`, tf.TableName), tf.Value)
			case UpdatedAfter:
				db = db.Where(fmt.Sprintf(`"%s"."updated_at" >= ?`, tf.TableName), tf.Value)
			case UpdatedBefore:
				db = db.Where(fmt.Sprintf(`"%s"."updated_at" <= ?`, tf.TableName), tf.Value)
			case PublishAfter:
				db = db.Where(fmt.Sprintf(`"%s"."publish_at" >= ?`, tf.TableName), tf.Value)
			case PublishBefore:
				db = db.Where(fmt.Sprintf(`"%s"."publish_at" <= ?`, tf.TableName), tf.Value)
			case DueDateAfter:
				db = db.Where(fmt.Sprintf(`"%s"."due_date" >= ?`, tf.TableName), tf.Value)
			case DueDateBefore:
				db = db.Where(fmt.Sprintf(`"%s"."due_date" <= ?`, tf.TableName), tf.Value)
			default:
				continue
			}
		}
	}

	var orConds []clause.Expression
	for key, values := range bm.OrFilters {
		if len(values) == 0 {
			continue
		}

		for _, value := range values {
			if strings.Contains(key, ".") {
				parts := strings.Split(key, ".")
				if len(parts) == 2 {
					tableName := parts[0]
					fieldName := parts[1]
					orConds = append(orConds, clause.Expr{
						SQL:  fmt.Sprintf(`"%s"."%s" = ?`, tableName, fieldName),
						Vars: []interface{}{value},
					})
				}
			} else {
				orConds = append(orConds, clause.Expr{
					SQL:  fmt.Sprintf(`%s = ?`, key),
					Vars: []interface{}{value},
				})
			}
		}
	}

	if len(orConds) > 0 {
		db = db.Where(clause.Or(orConds...))
	}

	if len(bm.OrFiltersWithAndConds) > 0 {
		var orGroupExpressions []clause.Expression

		for _, andConditions := range bm.OrFiltersWithAndConds {
			if len(andConditions) == 0 {
				continue
			}

			var andGroupExpressions []clause.Expression
			for _, cond := range andConditions {
				sqlStr := ""
				var vars []interface{}

				// Handle operand with dot (e.g., "orders.status")
				if strings.Contains(cond.Operand, ".") {
					parts := strings.Split(cond.Operand, ".")
					if len(parts) == 2 {
						tableName := parts[0]
						fieldName := parts[1]
						switch cond.Operator {
						case EqualOperator:
							sqlStr = fmt.Sprintf(`"%s"."%s" = ?`, tableName, fieldName)
						case NotEqualOperator:
							sqlStr = fmt.Sprintf(`"%s"."%s" != ?`, tableName, fieldName)
						default:
							continue
						}
						vars = []interface{}{cond.Value}
					}
				} else {
					switch cond.Operator {
					case EqualOperator:
						sqlStr = fmt.Sprintf("%s = ?", cond.Operand)
					case NotEqualOperator:
						sqlStr = fmt.Sprintf("%s != ?", cond.Operand)
					default:
						continue
					}
					vars = []interface{}{cond.Value}
				}

				if sqlStr != "" {
					andGroupExpressions = append(andGroupExpressions, clause.Expr{
						SQL:  sqlStr,
						Vars: vars,
					})
				}
			}

			if len(andGroupExpressions) > 0 {
				// Group all conditions in this set with AND: (cond1 AND cond2 AND ...)
				andClause := clause.And(andGroupExpressions...)
				orGroupExpressions = append(orGroupExpressions, andClause)
			}
		}

		if len(orGroupExpressions) > 0 {
			// Combine all AND groups with OR: (group1) OR (group2) OR (...)
			db = db.Where(clause.Or(orGroupExpressions...))
		}
	}

	return db
}

// applyRelations applies relations dynamically
func (bm *BuilderModel) applyRelations(db *gorm.DB) *gorm.DB {
	for _, relation := range bm.Relations {
		db = db.Preload(relation)
	}
	return db
}

// applyFilters applies dynamic filtering based on the Filters map in BuilderModel
func (bm *BuilderModel) applyLikes(db *gorm.DB) *gorm.DB {
	// Handle normal Likes (AND conditions) we don't need it in this project
	for key, value := range bm.Likes {
		if valStr, ok := value.(string); ok {
			if strings.Contains(key, ".") {
				parts := strings.Split(key, ".")
				if len(parts) == 2 {
					tableName := parts[0]
					fieldName := parts[1]
					db = db.Where(fmt.Sprintf(`CAST("%s"."%s" AS TEXT) ILIKE ?`, tableName, fieldName), "%"+valStr+"%")
				}
			} else {
				db = db.Where(clause.Expr{
					SQL:  "CAST(" + key + " AS TEXT) ILIKE ?",
					Vars: []interface{}{"%" + valStr + "%"},
				})
			}
		}
	}

	// Handle OrLikes (OR conditions)
	if len(bm.OrLikes) > 0 {
		var orClauses []string
		var orVars []interface{}

		for key, value := range bm.OrLikes {
			if valStr, ok := value.(string); ok {
				if strings.Contains(key, ".") {
					parts := strings.Split(key, ".")
					if len(parts) == 2 {
						tableName := parts[0]
						fieldName := parts[1]
						orClauses = append(orClauses, fmt.Sprintf(`CAST("%s"."%s" AS TEXT) ILIKE ?`, tableName, fieldName))
						orVars = append(orVars, "%"+valStr+"%")
					}
				} else {
					orClauses = append(orClauses, "CAST("+key+" AS TEXT) ILIKE ?")
					orVars = append(orVars, "%"+valStr+"%")
				}
			}
		}

		if len(orClauses) > 0 {
			db = db.Where("("+strings.Join(orClauses, " OR ")+")", orVars...)
		}
	}

	return db
}

// applyGlobalSearch creates an isolated OR block for fields marked with search:"true"
func (bm *BuilderModel) applyGlobalSearch(db *gorm.DB) *gorm.DB {
	if bm.GlobalSearchVal != "" && len(bm.GlobalSearchCols) > 0 {
		var orClauses []string
		var orVars []interface{}

		for _, col := range bm.GlobalSearchCols {
			// col is already formatted as "table"."field" from the middleware
			orClauses = append(orClauses, fmt.Sprintf(`CAST(%s AS TEXT) ILIKE ?`, col))
			orVars = append(orVars, "%"+bm.GlobalSearchVal+"%")
		}

		if len(orClauses) > 0 {
			// This creates: AND (col1 ILIKE %val% OR col2 ILIKE %val%)
			db = db.Where("("+strings.Join(orClauses, " OR ")+")", orVars...)
		}
	}
	return db
}

// getTotalItems retrieves the total number of items matching the query
func (bm *BuilderModel) getTotalItems(db *gorm.DB) (int64, error) {
	var totalItems int64
	if err := db.Count(&totalItems).Error; err != nil {
		return 0, err
	}
	return totalItems, nil
}

// calculateOffset calculates the offset for pagination
func (bm *BuilderModel) calculateOffset() int {
	return int((bm.Page - 1) * bm.PageSize)
}

// QueryBuilderScope returns a GORM scope function for paginating results
func (bm *BuilderModel) QueryBuilderScope(db *gorm.DB) (*gorm.DB, error) {
	// Validate and set default pagination parameters
	bm.validatePaginationParams()

	// Apply filters and sorting
	db = bm.applyJoins(db)
	db = bm.applyFilters(db)
	db = bm.applyRelations(db)
	db = bm.applyLikes(db)
	db = bm.applyGlobalSearch(db)
	db = bm.applySorting(db)

	// Get total items count
	totalItems, err := bm.getTotalItems(db)
	if err != nil {
		return nil, err
	}
	bm.TotalItems = totalItems

	// Calculate total pages
	totalPages := (totalItems + int64(bm.PageSize) - 1) / int64(bm.PageSize)
	if bm.Page > uint(totalPages) && totalPages > 0 {
		return nil, errors.New("requested page exceeds total number of pages")
	}

	// Calculate offset and apply pagination
	offset := bm.calculateOffset()
	db = db.Offset(offset).Limit(int(bm.PageSize))

	return db, nil
}

// CreatePaginateModel constructs a PaginateModel from the query result
func (bm *BuilderModel) CreatePaginateModel(db *gorm.DB, result interface{}) (*PaginateModel, error) {
	if bm.TotalItems == 0 {
		totalItems, err := bm.getTotalItems(db)
		if err != nil {
			return nil, err
		}
		bm.TotalItems = totalItems
	}

	totalPages := (bm.TotalItems + int64(bm.PageSize) - 1) / int64(bm.PageSize)

	// Execute the query
	if err := db.Find(result).Error; err != nil {
		return nil, err
	}

	return &PaginateModel{
		PageSize:    bm.PageSize,
		CurrentPage: bm.Page,
		TotalPages:  totalPages,
		TotalItems:  bm.TotalItems,
		Items:       result,
	}, nil
}
