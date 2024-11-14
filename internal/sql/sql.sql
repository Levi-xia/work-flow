CREATE TABLE IF NOT EXISTS `process_define` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL DEFAULT '',
    `code` VARCHAR(255) NOT NULL DEFAULT '',
    `version` INT NOT NULL DEFAULT 0,
    `content` TEXT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) AUTO_INCREMENT = 100000;

CREATE TABLE IF NOT EXISTS `process_instance` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `process_define_id` INT NOT NULL DEFAULT 0,
    `status` VARCHAR(32) NOT NULL DEFAULT '',
    `variables` JSON,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) AUTO_INCREMENT = 100000;

CREATE TABLE IF NOT EXISTS `process_task` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `process_instance_id` INT NOT NULL DEFAULT 0,
    `name` VARCHAR(255) NOT NULL DEFAULT '',
    `code` VARCHAR(255) NOT NULL DEFAULT '',
    `status` VARCHAR(32) NOT NULL DEFAULT '',
    `variables` JSON,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) AUTO_INCREMENT = 100000;