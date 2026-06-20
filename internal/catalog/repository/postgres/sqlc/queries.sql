-- name: CreateServiceCategory :one
INSERT INTO service_categories (name)
VALUES ($1)
RETURNING
    id,
    name,
    created_at,
    updated_at;


-- name: UpdateServiceCategory :one
UPDATE service_categories
SET name       = $2,
    updated_at = now()
WHERE id = $1
RETURNING
    id,
    name,
    created_at,
    updated_at;


-- name: ListServiceCategories :many
SELECT id,
       name,
       created_at,
       updated_at
FROM service_categories
WHERE sqlc.arg('search')::text = ''
   OR strpos(lower(name), lower(sqlc.arg('search')::text)) > 0
ORDER BY
    CASE WHEN sqlc.arg('sort_order')::text = 'asc' THEN name END ASC,
    CASE WHEN sqlc.arg('sort_order')::text = 'desc' THEN name END DESC,
    id ASC
LIMIT sqlc.arg('count')::int
OFFSET sqlc.arg('offset')::int;


-- name: DeleteServiceCategory :one
DELETE FROM service_categories
WHERE id = $1
RETURNING id;


-- name: CreateInputCharacteristicsTemplate :one
INSERT INTO input_characteristics_templates (name, description)
VALUES ($1, $2)
RETURNING
    id,
    name,
    description,
    created_at,
    updated_at;


-- name: PatchInputCharacteristicsTemplate :one
UPDATE input_characteristics_templates
SET name = CASE
               WHEN sqlc.arg('set_name')::bool THEN sqlc.arg('name')::varchar
               ELSE name
           END,
    description = CASE
                      WHEN sqlc.arg('set_description')::bool THEN sqlc.narg('description')::text
                      ELSE description
                  END,
    updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING
    id,
    name,
    description,
    created_at,
    updated_at;


-- name: GetInputCharacteristicsTemplateByID :one
SELECT id,
       name,
       description,
       created_at,
       updated_at
FROM input_characteristics_templates
WHERE id = $1;

-- name: DeleteInputCharacteristicsTemplate :one
DELETE FROM input_characteristics_templates
WHERE id = $1
RETURNING id;

-- name: ListInputCharacteristicsByTemplateID :many
SELECT ic.id,
       ic.name,
       ic.type,
       ic.created_at,
       ic.updated_at
FROM input_characteristics ic
JOIN input_characteristic_template_items icti
  ON icti.input_characteristic_id = ic.id
WHERE icti.template_id = $1
ORDER BY ic.id;


-- name: ListInputCharacteristicsTemplates :many
SELECT ict.id,
       ict.name,
       ict.description,
       ict.created_at,
       ict.updated_at
FROM input_characteristics_templates ict
WHERE (
        sqlc.arg('search')::text = ''
        OR strpos(lower(ict.name), lower(sqlc.arg('search')::text)) > 0
    )
  AND (
        cardinality(sqlc.arg('input_characteristic_ids')::int[]) = 0
        OR EXISTS (
            SELECT 1
            FROM input_characteristic_template_items item
            WHERE item.template_id = ict.id
              AND item.input_characteristic_id = ANY(sqlc.arg('input_characteristic_ids')::int[])
        )
    )
ORDER BY
    CASE WHEN sqlc.arg('sort_order')::text = 'asc' THEN ict.name END ASC,
    CASE WHEN sqlc.arg('sort_order')::text = 'desc' THEN ict.name END DESC,
    ict.id ASC
LIMIT sqlc.arg('count')::int
OFFSET sqlc.arg('offset')::int;


-- name: CreateInputCharacteristic :one
INSERT INTO input_characteristics (name, type)
VALUES ($1, $2)
RETURNING
    id,
    name,
    type,
    created_at,
    updated_at;


-- name: PatchInputCharacteristic :one
UPDATE input_characteristics
SET name = CASE
               WHEN sqlc.arg('set_name')::bool THEN sqlc.arg('name')::varchar
               ELSE name
           END,
    type = CASE
               WHEN sqlc.arg('set_type')::bool THEN sqlc.arg('type')::input_characteristics_type
               ELSE type
           END,
    updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING
    id,
    name,
    type,
    created_at,
    updated_at;


-- name: GetInputCharacteristicByID :one
SELECT ic.id,
       ic.name,
       ic.type,
       COALESCE(
           array_agg(icti.template_id ORDER BY icti.template_id)
               FILTER (WHERE icti.template_id IS NOT NULL),
           '{}'
       )::int[] AS template_ids,
       ic.created_at,
       ic.updated_at
FROM input_characteristics ic
LEFT JOIN input_characteristic_template_items icti ON icti.input_characteristic_id = ic.id
WHERE ic.id = $1
GROUP BY ic.id, ic.name, ic.type, ic.created_at, ic.updated_at;


-- name: ListInputCharacteristics :many
SELECT ic.id,
       ic.name,
       ic.type,
       ic.created_at,
       ic.updated_at
FROM input_characteristics ic
WHERE sqlc.arg('search')::text = ''
   OR strpos(lower(ic.name), lower(sqlc.arg('search')::text)) > 0
ORDER BY ic.id
LIMIT sqlc.arg('count')::int
OFFSET sqlc.arg('offset')::int;


-- name: CreateInputCharacteristicTemplateItem :exec
INSERT INTO input_characteristic_template_items (template_id, input_characteristic_id)
VALUES ($1, $2)
ON CONFLICT (template_id, input_characteristic_id) DO NOTHING;


-- name: DeleteInputCharacteristicTemplateItems :exec
DELETE FROM input_characteristic_template_items
WHERE input_characteristic_id = $1;

-- name: DeleteInputCharacteristicTemplateItemsByTemplateID :exec
DELETE FROM input_characteristic_template_items
WHERE template_id = $1;


-- name: CreateService :one
INSERT INTO services (name, base_price, description, type, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING
    id,
    name,
    base_price,
    description,
    type,
    status,
    created_at,
    updated_at;


-- name: PatchService :one
UPDATE services
SET name = CASE
               WHEN sqlc.arg('set_name')::bool THEN sqlc.arg('name')::varchar
               ELSE name
           END,
    base_price = CASE
                     WHEN sqlc.arg('set_base_price')::bool THEN sqlc.arg('base_price')::numeric
                     ELSE base_price
                 END,
    description = CASE
                      WHEN sqlc.arg('set_description')::bool THEN sqlc.narg('description')::text
                      ELSE description
                  END,
    updated_at = now()
WHERE id = sqlc.arg('id')
  AND type = sqlc.arg('type')::service_type
RETURNING
    id,
    name,
    base_price,
    description,
    type,
    status,
    created_at,
    updated_at;


-- name: GetServiceByIDAndType :one
SELECT id,
       name,
       base_price,
       description,
       type,
       status,
       created_at,
       updated_at
FROM services
WHERE id = sqlc.arg('id')
  AND type = sqlc.arg('type')::service_type;


-- name: ListServices :many
SELECT id,
       name,
       base_price,
       description,
       type
FROM services
WHERE type = sqlc.arg('type')::service_type
  AND (
        sqlc.arg('search')::text = ''
        OR strpos(lower(name), lower(sqlc.arg('search')::text)) > 0
    )
ORDER BY id
LIMIT sqlc.arg('count')::int
OFFSET sqlc.arg('offset')::int;


-- name: CreateServiceInputCharacteristic :exec
INSERT INTO service_input_characteristics (
    input_characteristics_id,
    service_id,
    is_required,
    sort_order
)
VALUES ($1, $2, $3, $4);


-- name: DeleteServiceInputCharacteristics :exec
DELETE FROM service_input_characteristics
WHERE service_id = $1;


-- name: ListServiceInputCharacteristics :many
SELECT ic.id,
       ic.name,
       ic.type,
       sic.is_required,
       sic.sort_order
FROM service_input_characteristics sic
JOIN input_characteristics ic ON ic.id = sic.input_characteristics_id
WHERE sic.service_id = $1
ORDER BY sic.sort_order NULLS LAST, ic.id;


-- name: ListServiceCategoriesByServiceID :many
SELECT sc.id,
       sc.name,
       sc.created_at,
       sc.updated_at
FROM service_category_services scs
JOIN service_categories sc ON sc.id = scs.service_category_id
WHERE scs.service_id = $1
ORDER BY sc.name, sc.id;


-- name: ListServiceModifiersByServiceID :many
SELECT id,
       service_id,
       name,
       selection_type,
       sort_order,
       is_required
FROM service_modifiers
WHERE service_id = $1
ORDER BY sort_order NULLS LAST, id;


-- name: ListServiceModifierValues :many
SELECT id,
       name,
       service_modifier_id,
       additional_price,
       is_active,
       sort_order,
       created_at,
       updated_at
FROM service_modifier_values
WHERE service_modifier_id = $1
ORDER BY sort_order NULLS LAST, id;


-- name: DeleteService :one
DELETE FROM services
WHERE id = $1
RETURNING id;
