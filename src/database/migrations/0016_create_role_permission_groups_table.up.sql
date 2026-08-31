CREATE TABLE IF NOT EXISTS role_permission_groups (
    role_model_id BIGINT NOT NULL
    CONSTRAINT fk_role_permission_groups_role
    REFERENCES roles(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE,

    permission_group_model_id BIGINT NOT NULL
    CONSTRAINT fk_role_permission_groups_group
    REFERENCES permission_groups(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE,

    PRIMARY KEY (role_model_id, permission_group_model_id)
    );
