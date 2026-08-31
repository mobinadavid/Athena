CREATE TABLE IF NOT EXISTS role_user
(
    user_model_id bigint NOT NULL,
    role_model_id bigint NOT NULL,
    PRIMARY KEY (user_model_id, role_model_id),
    CONSTRAINT fk_role_user_user_model
    FOREIGN KEY (user_model_id)
    REFERENCES users (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
    CONSTRAINT fk_role_user_role_model
    FOREIGN KEY (role_model_id)
    REFERENCES roles (id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
    );