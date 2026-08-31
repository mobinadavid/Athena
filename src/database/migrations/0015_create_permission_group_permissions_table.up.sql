CREATE TABLE IF NOT EXISTS permission_group_permissions (
    permission_group_model_id BIGINT NOT NULL
    CONSTRAINT fk_permission_group_permissions_group
    REFERENCES permission_groups(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE,

    permission_model_id BIGINT NOT NULL
    CONSTRAINT fk_permission_group_permissions_permission
    REFERENCES permissions(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE,

    PRIMARY KEY (permission_group_model_id, permission_model_id)
    );
