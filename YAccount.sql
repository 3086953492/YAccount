insert into users (username, password, nickname, avatar,role) values ('admin', '$2a$10$u/T1M5H89WQ8/dnukU.fz.7z9CHgy5H8MWdiM1xPruuwVZ7redpEy', 'admin', 'https://example.com/avatar.jpg','admin');    -- 此密码为admin123的加密值

-- 插入系统基础配置示例
INSERT INTO `system_infos` (`config_key`, `config_value`, `config_type`, `description`, `status`, `created_by`) VALUES
('system_name', 'YAccount', 'string', '系统名称', 1, 1),
('system_icon', 'https://example.com/icon.png', 'string', '系统图标URL', 1, 1),
('system_logo', 'https://example.com/logo.png', 'string', '系统Logo URL', 1, 1),
('system_description', 'YAccount - 用户账户管理系统', 'string', '系统描述', 1, 1);