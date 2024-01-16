CREATE TABLE IF NOT EXISTS `log_details`
(
    `id`             BIGINT UNSIGNED                                           NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `primary_log_id` BIGINT UNSIGNED                                           NOT NULL,
    `created_at`     TIMESTAMP                                                 NOT NULL,
    `log_level`      ENUM ('trace', 'debug', 'info', 'warn', 'error', 'fatal') NOT NULL,
    `log`            LONGTEXT                                                  NOT NULL,
    FOREIGN KEY (`primary_log_id`) REFERENCES `primary_logs` (`id`) -- no cascade is required
) DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;
