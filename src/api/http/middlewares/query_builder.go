package middlewares

import (
	"athena/src/database/scopes"
	"athena/src/pkg/logger"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm/utils"

	"reflect"
	"strconv"
	"strings"
	"time"
)

func QueryParametersBuilderMiddleware(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		builder := &scopes.BuilderModel{}
		v := reflect.ValueOf(model)
		var modelDbName string
		// converting to pointer, as method TableName is defined on a pointer
		if v.Kind() != reflect.Ptr {
			v = reflect.New(reflect.TypeOf(model))
		}

		// getting table name in database for db operations
		tblNameMethod := v.MethodByName("TableName")
		if tblNameMethod.IsValid() {
			modelDbName = tblNameMethod.Call(nil)[0].String()
		} else {
			logger.LogErrorDirect("failed to get model table name", nil,
				zap.String("middleware", "QueryParametersBuilderMiddleware"),
			)
		}

		page, err := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(scopes.DefaultPage)))
		if err != nil {
			page = scopes.DefaultPage
		}
		builder.Page = uint(page)

		pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(scopes.DefaultPageSize)))
		if err != nil {
			pageSize = scopes.DefaultPageSize
		}
		builder.PageSize = uint(pageSize)

		sortOrder := c.Query("sort_order")
		if sortOrder != "desc" && sortOrder != "asc" {
			sortOrder = "desc"
		}
		builder.SortOrder = sortOrder

		if createdAfterStr := c.Query("created_after"); createdAfterStr != "" {
			createdAfter, err := time.Parse(time.RFC3339, createdAfterStr)
			if err == nil {
				builder.AddTimeFilter(&scopes.TimeFilter{
					TableName: modelDbName,
					Type:      scopes.CreatedAfter,
					Value:     &createdAfter,
				})
			}
		}

		if createdBeforeStr := c.Query("created_before"); createdBeforeStr != "" {
			createdBefore, err := time.Parse(time.RFC3339, createdBeforeStr)
			if err == nil {
				builder.AddTimeFilter(&scopes.TimeFilter{
					TableName: modelDbName,
					Type:      scopes.CreatedBefore,
					Value:     &createdBefore,
				})
			}
		}

		if updatedAfterStr := c.Query("updated_after"); updatedAfterStr != "" {
			updatedAfter, err := time.Parse(time.RFC3339, updatedAfterStr)
			if err == nil {
				builder.AddTimeFilter(&scopes.TimeFilter{
					TableName: modelDbName,
					Type:      scopes.UpdatedAfter,
					Value:     &updatedAfter,
				})
			}
		}

		if updatedBeforeStr := c.Query("updated_before"); updatedBeforeStr != "" {
			updatedBefore, err := time.Parse(time.RFC3339, updatedBeforeStr)
			if err == nil {
				builder.AddTimeFilter(&scopes.TimeFilter{
					TableName: modelDbName,
					Type:      scopes.UpdatedBefore,
					Value:     &updatedBefore,
				})
			}
		}

		if dueDateAfterStr := c.Query("due_date_after"); dueDateAfterStr != "" {
			dueDateAfter, err := time.Parse(time.RFC3339, dueDateAfterStr)
			if err == nil {
				builder.AddTimeFilter(&scopes.TimeFilter{
					TableName: modelDbName,
					Type:      scopes.DueDateAfter,
					Value:     &dueDateAfter,
				})
			}
		}

		if dueDateBeforeStr := c.Query("due_date_before"); dueDateBeforeStr != "" {
			dueDateBefore, err := time.Parse(time.RFC3339, dueDateBeforeStr)
			if err == nil {
				builder.AddTimeFilter(&scopes.TimeFilter{
					TableName: modelDbName,
					Type:      scopes.DueDateBefore,
					Value:     &dueDateBefore,
				})
			}
		}
		//test := reflect.TypeOf(model).Name()
		filters := make(map[string]interface{})
		InFilters := make(map[string][]interface{})
		likes := make(map[string]interface{})
		M2MFilters := make(map[string][]interface{})
		M2MOrLikes := make(map[string]interface{})
		modelType := reflect.TypeOf(model)

		sortBy := c.DefaultQuery("sort_by", "created_at")
		filtersFromQuery := parseFilters(c)
		likesFromQuery := c.QueryMap("likes")
		globalSearchVal := strings.TrimSpace(c.Query("global_search"))
		var globalSearchCols []string

		// Handle relation.updated_after and relation.updated_before
		for key, value := range c.Request.URL.Query() {
			if strings.Contains(key, ".") {
				parts := strings.SplitN(key, ".", 2)
				if len(parts) != 2 {
					continue
				}
				relation, subKey := parts[0], parts[1]
				if subKey == "updated_before" || subKey == "updated_after" || subKey == "created_before" || subKey == "created_after" {
					t, err := time.Parse(time.RFC3339, value[0])
					if err != nil {
						continue
					}
					// Build a special key, e.g., "ProjectFinancialDetail.updated_before"
					var tf scopes.TimeFilter
					tf.TableName = relation
					switch subKey {
					case "created_before":
						tf.Type = scopes.CreatedBefore
						tf.Value = &t
					case "created_after":
						tf.Type = scopes.CreatedAfter
						tf.Value = &t
					case "updated_before":
						tf.Type = scopes.UpdatedBefore
						tf.Value = &t
					case "updated_after":
						tf.Type = scopes.UpdatedAfter
						tf.Value = &t
					default:
						continue
					}

					builder.AddTimeFilter(&tf)
					if !utils.Contains(builder.Joins, relation) {
						builder.Joins = append(builder.Joins, relation)
					}
				}
			}
		}

		// --- Main Field Loop ---
		for i := 0; i < modelType.NumField(); i++ {
			field := modelType.Field(i)
			if jsonTag, ok := field.Tag.Lookup("json"); ok && jsonTag != "" {
				jsonTag = strings.Split(jsonTag, ",")[0]

				// 1. Handle Filters
				if filterable, ok := field.Tag.Lookup("filter"); ok && filterable == "true" {
					if valueList, exists := filtersFromQuery[jsonTag]; exists {
						filterName := modelDbName + "." + jsonTag
						InFilters[filterName] = valueList

					}
				}

				// 2. Handle Joins
				if connectable, ok := field.Tag.Lookup("join"); ok && connectable == "true" {
					for key, valueList := range filtersFromQuery {
						if strings.Contains(key, ".") {
							prefix := strings.SplitN(key, ".", 2)[0]
							parts := strings.SplitN(key, ".", 2)
							if len(parts) != 2 {
								continue
							}
							prefix, subFieldJson := parts[0], parts[1]
							if field.Name == prefix {
								// Get the nested type
								nestedType := field.Type
								if nestedType.Kind() == reflect.Ptr {
									nestedType = nestedType.Elem()
								}

								// handling slices
								if nestedType.Kind() == reflect.Slice || nestedType.Kind() == reflect.Array {
									nestedType = nestedType.Elem()
									if nestedType.Kind() == reflect.Ptr {
										nestedType = nestedType.Elem()
									}
								}

								// Check if the subField is filterable
								for i := 0; i < nestedType.NumField(); i++ {
									nestedField := nestedType.Field(i)
									jsonTag := strings.Split(nestedField.Tag.Get("json"), ",")[0]
									if jsonTag == subFieldJson && nestedField.Tag.Get("filter") == "true" {
										if strings.Contains(field.Tag.Get("gorm"), "many2many") {
											for _, tag := range strings.Split(field.Tag.Get("gorm"), ";") {
												if strings.HasPrefix(tag, "many2many:") {
													M2MFilters[key] = valueList
													break
												}
											}
										} else {
											InFilters[key] = valueList
											if !utils.Contains(builder.Joins, field.Name) {
												builder.Joins = append(builder.Joins, field.Name)
											}
										}
										break
									}
								}
							}
						}
					}

					for key, value := range likesFromQuery {
						if strings.Contains(key, ".") {
							parts := strings.SplitN(key, ".", 2)
							if len(parts) != 2 {
								continue
							}
							prefix, subFieldJson := parts[0], parts[1]
							if field.Name != prefix {
								continue
							}

							nestedType := field.Type
							if nestedType.Kind() == reflect.Ptr {
								nestedType = nestedType.Elem()
							}

							// handling slices
							if nestedType.Kind() == reflect.Slice || nestedType.Kind() == reflect.Array {
								nestedType = nestedType.Elem()
								if nestedType.Kind() == reflect.Ptr {
									nestedType = nestedType.Elem()
								}
							}

							for i := 0; i < nestedType.NumField(); i++ {
								nestedField := nestedType.Field(i)
								jsonTag := strings.Split(nestedField.Tag.Get("json"), ",")[0]
								if jsonTag == subFieldJson && nestedField.Tag.Get("like") == "true" {
									if strings.Contains(field.Tag.Get("gorm"), "many2many") {
										for _, tag := range strings.Split(field.Tag.Get("gorm"), ";") {
											if strings.HasPrefix(tag, "many2many:") {
												M2MOrLikes[key] = value
												break
											}
										}
									} else {
										likes[key] = value
										if !utils.Contains(builder.Joins, field.Name) {
											builder.Joins = append(builder.Joins, field.Name)
										}
									}
									break
								}
							}
						}
					}
					// Check if sort_by targets a field in this relation
					if strings.Contains(sortBy, ".") {
						parts := strings.SplitN(sortBy, ".", 2)
						if len(parts) == 2 {
							prefix, subFieldJson := parts[0], parts[1]
							if field.Name == prefix {
								nestedType := field.Type
								if nestedType.Kind() == reflect.Ptr {
									nestedType = nestedType.Elem()
								}

								// handling slices
								if nestedType.Kind() == reflect.Slice || nestedType.Kind() == reflect.Array {
									nestedType = nestedType.Elem()
									if nestedType.Kind() == reflect.Ptr {
										nestedType = nestedType.Elem()
									}
								}

								for i := 0; i < nestedType.NumField(); i++ {
									nestedField := nestedType.Field(i)
									jsonTag := strings.Split(nestedField.Tag.Get("json"), ",")[0]
									if jsonTag == subFieldJson && nestedField.Tag.Get("sort") == "true" {
										builder.SortBy = sortBy
										if !utils.Contains(builder.Joins, field.Name) {
											builder.Joins = append(builder.Joins, field.Name)
										}
										break
									}
								}
							}
						}
					}
				}

				// 3. Handle Likes
				if likeable, ok := field.Tag.Lookup("like"); ok && likeable == "true" {
					if value, exists := likesFromQuery[jsonTag]; exists {
						filterName := modelDbName + "." + jsonTag
						likes[filterName] = value
					}
				}

				// 4. Handle Sort
				if sortable, ok := field.Tag.Lookup("sort"); ok && sortable == "true" {
					if sortBy == jsonTag {
						builder.SortBy = modelDbName + "." + sortBy
					}
				}
				//5. Handle Global Search
				if searchable, ok := field.Tag.Lookup("search"); ok && searchable == "true" {
					// Safely format as "table_name"."column_name" to prevent ambiguity
					filterName := fmt.Sprintf(`"%s"."%s"`, modelDbName, jsonTag)
					globalSearchCols = append(globalSearchCols, filterName)
				}
			}
		}

		builder.Filters = filters
		builder.OrLikes = likes
		builder.InFilters = InFilters
		builder.M2MFilters = M2MFilters
		builder.M2MOrLikes = M2MOrLikes
		builder.GlobalSearchVal = globalSearchVal
		builder.GlobalSearchCols = globalSearchCols

		c.Set("query_parameters_builder", builder)
		c.Next()
	}
}

func parseFilters(c *gin.Context) map[string][]interface{} {
	query := c.Request.URL.Query()
	filters := make(map[string][]interface{})

	// Iterate over all query params
	for key, values := range query {
		if strings.HasPrefix(key, "filters[") {
			// Extract "id" from "filters[id]"
			innerKey := strings.TrimSuffix(strings.TrimPrefix(key, "filters["), "]")
			filters[innerKey] = []interface{}{}

			for _, value := range values {
				// Split by comma to handle "val1,val2"
				if strings.Contains(value, ",") {
					splitValues := strings.Split(value, ",")
					for _, sv := range splitValues {
						// TrimSpace cleans up accidental spaces like "val1, val2"
						filters[innerKey] = append(filters[innerKey], strings.TrimSpace(sv))
					}
				} else {
					filters[innerKey] = append(filters[innerKey], strings.TrimSpace(value))
				}
			}
		}
	}

	return filters
}
