package repositories

import (
	"database/sql"
	"errors"
	"fluxton/errs"
	"fluxton/models"
	"fluxton/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/samber/do"
)

type FormFieldRepository struct {
	db *sqlx.DB
}

func NewFormFieldRepository(injector *do.Injector) (*FormFieldRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &FormFieldRepository{db: db}, nil
}

func (r *FormFieldRepository) ListForForm(formUUID uuid.UUID) ([]models.FormField, error) {
	query := "SELECT * FROM fluxton.form_fields WHERE form_uuid = $1;"

	rows, err := r.db.Queryx(query, formUUID)
	if err != nil {
		return nil, utils.FormatError(err, "select", utils.GetMethodName())
	}
	defer rows.Close()

	var forms []models.FormField
	for rows.Next() {
		var form models.FormField
		if err := rows.StructScan(&form); err != nil {
			return nil, utils.FormatError(err, "scan", utils.GetMethodName())
		}
		forms = append(forms, form)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.FormatError(err, "iterate", utils.GetMethodName())
	}

	return forms, nil
}

func (r *FormFieldRepository) GetByUUID(formUUID uuid.UUID) (models.FormField, error) {
	query := "SELECT %s FROM fluxton.form_fields WHERE uuid = $1"
	query = fmt.Sprintf(query, utils.GetColumns[models.FormField]())

	var form models.FormField
	err := r.db.Get(&form, query, formUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.FormField{}, errs.NewNotFoundError("form.error.notFound")
		}

		return models.FormField{}, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return form, nil
}

func (r *FormFieldRepository) ExistsByUUID(formFieldUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.form_fields WHERE uuid = $1)"

	var exists bool
	err := r.db.Get(&exists, query, formFieldUUID)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *FormFieldRepository) ExistsByAnyLabelForForm(labels []string, formUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.form_fields WHERE label = ANY($1) AND form_uuid = $2)"

	var exists bool
	err := r.db.Get(&exists, query, pq.Array(labels), formUUID)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *FormFieldRepository) ExistsByLabelForForm(label string, formUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.form_fields WHERE label = $1 AND form_uuid = $2)"

	var exists bool
	err := r.db.Get(&exists, query, label, formUUID)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *FormFieldRepository) Create(formField *models.FormField) (*models.FormField, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, utils.FormatError(err, "transactionBegin", utils.GetMethodName())
	}

	query := `
    INSERT INTO fluxton.form_fields (
        form_uuid,
        label,
        type,
        description,
        is_required,
        options,
        min_length,
        max_length,
        min_value,
        max_value,
        pattern,
        default_value,
        start_date,
        end_date,
        date_format
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
    )
    RETURNING uuid
`

	queryErr := tx.QueryRowx(
		query,
		formField.FormUuid,
		formField.Label,
		formField.Type,
		formField.Description,
		formField.IsRequired,
		formField.Options,
		formField.MinLength,
		formField.MaxLength,
		formField.MinValue,
		formField.MaxValue,
		formField.Pattern,
		formField.DefaultValue,
		formField.StartDate,
		formField.EndDate,
		formField.DateFormat,
	).Scan(&formField.Uuid)

	if queryErr != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, utils.FormatError(queryErr, "insert", utils.GetMethodName())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, utils.FormatError(err, "transactionCommit", utils.GetMethodName())
	}

	return formField, nil
}

func (r *FormFieldRepository) CreateMany(formFields []models.FormField, formUUID uuid.UUID) ([]models.FormField, error) {
	createdFields := make([]models.FormField, 0, len(formFields))
	for i, formField := range formFields {
		formField.FormUuid = formUUID

		createdField, err := r.Create(&formField)
		if err != nil {
			return nil, fmt.Errorf("could not create form field at index %d: %v", i, err)
		}

		createdFields = append(createdFields, *createdField)
	}

	return createdFields, nil
}

func (r *FormFieldRepository) Update(formField *models.FormField) (*models.FormField, error) {
	query := `
		UPDATE fluxton.form_fields 
		SET 
		    label = :label, 
		    description = :description, 
		    type = :type, 
		    is_required = :is_required, 
		    options = :options, 
		    updated_at = :updated_at
		WHERE uuid = :uuid`

	res, err := r.db.NamedExec(query, formField)
	if err != nil {
		return &models.FormField{}, utils.FormatError(err, "update", utils.GetMethodName())
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &models.FormField{}, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	return formField, nil
}

func (r *FormFieldRepository) Delete(formFieldUUID uuid.UUID) (bool, error) {
	query := "DELETE FROM fluxton.form_fields WHERE uuid = $1"
	res, err := r.db.Exec(query, formFieldUUID)
	if err != nil {
		return false, utils.FormatError(err, "delete", utils.GetMethodName())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	return rowsAffected == 1, nil
}
