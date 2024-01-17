CREATE TABLE IF NOT EXISTS `primary_logs`
(
    `created_at`    TIMESTAMP                       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`    TIMESTAMP                       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at`    TIMESTAMP                       NULL,
    `id`            BIGINT UNSIGNED                 NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `uuid`          VARCHAR(255)                    NOT NULL,
    `file`          VARCHAR(255)                    NOT NULL,
    `process_state` ENUM ('init', 'done', 'failed') NOT NULL DEFAULT 'init'
) DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;