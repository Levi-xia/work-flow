CREATE TABLE IF NOT EXISTS `process_define` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL DEFAULT '',
    `code` VARCHAR(255) NOT NULL DEFAULT '',
    `user_id` INT(11) NOT NULL DEFAULT 0,
    `version` INT NOT NULL DEFAULT 0,
    `content` TEXT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) AUTO_INCREMENT = 100000;

ALTER TABLE `process_define` add column `user_id` INT(11) NOT NULL DEFAULT 0 AFTER `code`;

CREATE TABLE IF NOT EXISTS `process_instance` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `process_define_id` INT NOT NULL DEFAULT 0,
    `user_id` INT(11) NOT NULL DEFAULT 0,
    `status` VARCHAR(32) NOT NULL DEFAULT '',
    `variables` JSON,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) AUTO_INCREMENT = 100000;

ALTER TABLE `process_instance` add column `user_id` INT(11) NOT NULL DEFAULT 0 AFTER `process_define_id`;

CREATE TABLE IF NOT EXISTS `process_task` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `process_instance_id` INT NOT NULL DEFAULT 0,
    `form_instance_id` INT NOT NULL DEFAULT 0 AFTER `process_instance_id`
    `name` VARCHAR(255) NOT NULL DEFAULT '',
    `code` VARCHAR(255) NOT NULL DEFAULT '',
    `user_id` INT(11) NOT NULL DEFAULT 0,
    `status` VARCHAR(32) NOT NULL DEFAULT '',
    `variables` JSON,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) AUTO_INCREMENT = 100000;

ALTER TABLE `process_task` add column `user_id` INT(11) NOT NULL DEFAULT 0 AFTER `code`;

CREATE TABLE IF NOT EXISTS `action_record` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `action_define_id` INT NOT NULL DEFAULT 0,
    `process_task_id` INT NOT NULL DEFAULT 0,
    `input` TEXT NOT NULL,
    `output` TEXT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY `idx_action_define_id` (`action_define_id`),
    KEY `idx_process_task_id` (`process_task_id`)
) AUTO_INCREMENT = 100000;

CREATE TABLE IF NOT EXISTS `action_define` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL DEFAULT '',
    `code` VARCHAR(255) NOT NULL DEFAULT '',
    `user_id` INT(11) NOT NULL DEFAULT 0,
    `version` INT NOT NULL DEFAULT 0,
    `protocol` VARCHAR(32) NOT NULL DEFAULT '',
    `content` TEXT NOT NULL,
    `input_structs` TEXT NOT NULL,
    `output_checks` TEXT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY `uk_code_version` (`code`, `version`)
) AUTO_INCREMENT = 100000;

ALTER TABLE `action_define` add column `user_id` INT(11) NOT NULL DEFAULT 0 AFTER `code`;

CREATE TABLE IF NOT EXISTS `form_define` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL DEFAULT '',
    `code` VARCHAR(255) NOT NULL DEFAULT '',
    `user_id` INT(11) NOT NULL DEFAULT 0,
    `version` INT NOT NULL DEFAULT 0,
    `form_structure` TEXT NOT NULL,
    `component_structure` TEXT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY `uk_code_version` (`code`, `version`)
) AUTO_INCREMENT = 100000;

ALTER TABLE `form_define` add column `user_id` INT(11) NOT NULL DEFAULT 0 AFTER `code`;

CREATE TABLE IF NOT EXISTS `form_instance` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `form_define_id` INT NOT NULL DEFAULT 0,
    `form_data` TEXT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) AUTO_INCREMENT = 100000;
