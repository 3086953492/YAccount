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
