CREATE TABLE `journal`
(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(1024) NOT NULL,
    `category` TINYINT UNSIGNED NOT NULL,
    `created` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
)
