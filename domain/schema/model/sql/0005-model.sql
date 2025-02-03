-- The model_config table is a new table that is used to store configuration 
-- data for the model.
--
-- The provider tracker relies on the model_config table. Do not modify the
-- model_config table in a patch/build release. Only make changes to this table
-- during a major/minor release.
CREATE TABLE model_config (
    "key" TEXT NOT NULL PRIMARY KEY,
    value TEXT NOT NULL
);

-- The model_constraint table is a new table that is used to store the
-- constraints that are associated with a model.
CREATE TABLE model_constraint (
    model_uuid TEXT NOT NULL PRIMARY KEY,
    constraint_uuid TEXT NOT NULL,
    CONSTRAINT fk_model_constraint_model
    FOREIGN KEY (model_uuid)
    REFERENCES model (uuid),
    CONSTRAINT fk_model_constraint_constraint
    FOREIGN KEY (constraint_uuid)
    REFERENCES "constraint" (uuid)
);

-- v_model_constraint is a view to represent the current model constraints. If
-- no constraints have been set then expect this view to be empty. There will
-- also only ever be a maximum of 1 record in this view.
CREATE VIEW v_model_constraint AS
SELECT
    c.uuid,
    c.arch,
    c.cpu_cores,
    c.cpu_power,
    c.mem,
    c.root_disk,
    c.root_disk_source,
    c.instance_role,
    c.instance_type,
    c.container_type,
    c.virt_type,
    c.allocate_public_ip,
    c.image_id
FROM model_constraint mc
INNER JOIN v_constraint c ON mc.constraint_uuid = c.uuid;

-- v_model_constraint_tag is a view of all the constraint tags set for the
-- current model. It is expected that this view can be empty.
CREATE VIEW v_model_constraint_tag AS
SELECT
    ct.constraint_uuid,
    ct.tag
FROM constraint_tag ct
JOIN "constraint" c on ct.constraint_uuid = c.uuid
JOIN model_constraint mc on mc.constraint_uuid = c.uuid;

-- v_model_constraint_space is a view of all the constraint spaces set for the
-- current model. It is expected that this view can be empty.
CREATE VIEW v_model_constraint_space AS
SELECT
    cs.constraint_uuid,
    cs.space,
    cs.exclude
FROM constraint_space cs
JOIN "constraint" c on cs.constraint_uuid = c.uuid
JOIN model_constraint mc on mc.constraint_uuid = c.uuid;

-- v_model_constraint_zone is a view of all the constraint zones set for the
-- current model. It is expected that this view can be empty.
CREATE VIEW v_model_constraint_zone AS
SELECT
    cz.constraint_uuid,
    cz.zone
FROM constraint_zone cz
JOIN "constraint" c on cz.constraint_uuid = c.uuid
JOIN model_constraint mc on mc.constraint_uuid = c.uuid;