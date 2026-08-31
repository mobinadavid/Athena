CREATE TABLE IF NOT EXISTS admin_role
(
    admin_model_id bigint NOT NULL,
    role_model_id  bigint NOT NULL,
    PRIMARY KEY (admin_model_id, role_model_id),
    CONSTRAINT fk_admin_role_admin_model
        FOREIGN KEY (admin_model_id)
        REFERENCES admins (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_admin_role_role_model
        FOREIGN KEY (role_model_id)
        REFERENCES roles (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    );