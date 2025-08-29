CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(255) NOT NULL,
  `status` int DEFAULT '1' COMMENT '1:正常 0:禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  `nickname` varchar(50) DEFAULT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `role` varchar(50) DEFAULT 'user',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  KEY `idx_status_deleted` (`status`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';

insert into users (username, password, nickname, avatar,role) values ('admin', '$2a$10$u/T1M5H89WQ8/dnukU.fz.7z9CHgy5H8MWdiM1xPruuwVZ7redpEy', 'admin', 'https://example.com/avatar.jpg','admin');    -- 此密码为admin123的加密值

CREATE TABLE `system_infos` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `config_key` varchar(100) NOT NULL COMMENT '配置键名',
  `config_value` text COMMENT '配置值',
  `config_type` varchar(50) DEFAULT 'string' COMMENT '配置类型：string, number, boolean, json',
  `description` varchar(255) DEFAULT NULL COMMENT '配置描述',
  `status` int DEFAULT '1' COMMENT '状态：1启用 0禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_by` bigint unsigned DEFAULT NULL COMMENT '创建人ID',
  `updated_by` bigint unsigned DEFAULT NULL COMMENT '更新人ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_config_key` (`config_key`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_system_infos_created_by` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_system_infos_updated_by` FOREIGN KEY (`updated_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统配置表';

-- 插入系统基础配置示例
INSERT INTO `system_infos` (`config_key`, `config_value`, `config_type`, `description`, `status`, `created_by`) VALUES
('system_name', 'YAccount', 'string', '系统名称', 1, 1),
('system_icon', 'https://example.com/icon.png', 'string', '系统图标URL', 1, 1),
('system_logo', 'https://example.com/logo.png', 'string', '系统Logo URL', 1, 1),
('system_description', 'YAccount - 用户账户管理系统', 'string', '系统描述', 1, 1);