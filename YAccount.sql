insert into users (username, password, nickname, avatar,role) values ('admin', '$2a$10$u/T1M5H89WQ8/dnukU.fz.7z9CHgy5H8MWdiM1xPruuwVZ7redpEy', 'admin', 'https://example.com/avatar.jpg','admin');    -- 此密码为admin123的加密值

-- 插入系统基础配置示例
INSERT INTO `system_infos` (`config_key`, `config_value`, `config_type`, `description`, `status`, `created_by`) VALUES
('system_name', 'YAccount', 'string', '系统名称', 1, 1),
('system_icon', 'https://example.com/icon.png', 'string', '系统图标URL', 1, 1),
('system_logo', 'https://example.com/logo.png', 'string', '系统Logo URL', 1, 1),
('system_description', 'YAccount - 用户账户管理系统', 'string', '系统描述', 1, 1),
('system_client_id', 'client_dFhOnLNUYXXoybtqYcHL8w', 'string', '系统客户端ID', 1, 1);

-- 插入OAuth范围示例
INSERT INTO `oauth_scopes` (`name`, `description`, `is_default`) VALUES
		("read", "读取基本用户信息", true),
		("write", "修改用户信息", false),
		("profile", "访问完整用户资料", false),
		("admin", "管理员权限", false);