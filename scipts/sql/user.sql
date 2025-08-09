CREATE TABLE `user`
(
    `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID，主键',
    `username`      VARCHAR(50)  NOT NULL DEFAULT '' COMMENT '用户名',
    `email`         VARCHAR(100) NOT NULL DEFAULT '' COMMENT '邮箱，唯一',
    `password_hash` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '密码哈希值',
    `created_at`    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `status`        TINYINT      NOT NULL DEFAULT 1 COMMENT '状态, 1正常-1删除',
    PRIMARY KEY (`id`)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表';