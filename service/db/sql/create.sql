CREATE TABLE `file` (
    `id` bigint(20) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `hash` varchar(64) NOT NULL DEFAULT '' COMMENT '文件 hash ',
    `hash_algorithm` varchar(64) NOT NULL DEFAULT 'SHA256' COMMENT '文件 hash 算法',
    `size` bigint(20) DEFAULT '0' COMMENT '文件大小单位 Byte',
    `location` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
    `create_at` datetime default NOW() COMMENT '创建日期',
    `update_at` datetime default NOW() on update current_timestamp() COMMENT '更新日期',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态(可用/禁用/已删除等状态)'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user_file` (
    `id` bigint(20) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `user_id` bigint(20) NOT NULL COMMENT '用户 id',
    `file_id` varchar(64) NOT NULL DEFAULT '' COMMENT '文件id',
    `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
    `is_directory` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否是目录',
    `upload_at` datetime DEFAULT NOW() COMMENT '上传时间',
    `directory` varchar(256) NOT NULL DEFAULT '/' COMMENT '当前文件目录',
    `last_update_at` datetime DEFAULT NOW() ON UPDATE NOW() COMMENT '最后修改时间',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '文件状态(0正常1已删除2禁用)',
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user` (
    `id` bigint(20) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `role_id` bigint(20) NOT NULL DEFAULT 1 COMMENT '用户 role id',
    `name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
    `password` varchar(256) NOT NULL DEFAULT '' COMMENT '用户 encoded 密码',
    `email` varchar(64) DEFAULT '' COMMENT '邮箱',
    `phone` varchar(64) DEFAULT '' COMMENT '手机号',
    `email_validated` tinyint(1) DEFAULT 0 COMMENT '邮箱是否已验证',
    `phone_validated` tinyint(1) DEFAULT 0 COMMENT '手机号是否已验证',
    `sign_up_at` datetime DEFAULT NOW() COMMENT '注册日期',
    `last_active_at` datetime DEFAULT NOW() ON UPDATE NOW() COMMENT '最后活跃时间戳',
    `profile` text COMMENT '用户个人介绍',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '账户状态(启用/禁用/锁定/标记删除等)',
    KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `session` (
    `id` bigint(20) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `token` varchar(32) NOT NULL DEFAULT '' COMMENT '用户登录token',
    `user_id` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
    `create_at` datetime default NOW() COMMENT '创建日期',
    `expire_at` datetime default NOW() COMMENT '过期日期',
    `csrf_token` varchar(32) NOT NULL DEFAULT '' COMMENT '防止 csrf 攻击的 token',
    UNIQUE KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `group` (
    `id` bigint(20) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `create_at` datetime default NOW() COMMENT '创建日期',
    `update_at` datetime default NOW() on update current_timestamp() COMMENT '更新日期',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态(可用/禁用/已删除等状态)',
    `owner_id` bigint(20) NOT NULL,
    KEY `idx_owner_id` (`owner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user_group` (
    `id` bigint(20) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `user_id` bigint(20) NOT NULL COMMENT '用户 id',
    `group_id` bigint(20) NOT NULL COMMENT '组 id',
    UNIQUE KEY `idx_user_group` (`user_id`, `group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `role` (
    `id` bigint(20) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `name` varchar(64) NOT NULL DEFAULT '' COMMENT '角色名'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO role(id, name) VALUES (1, "超级管理员");
INSERT INTO role(id, name) VALUES (10000, "普通用户");

CREATE TABLE `group_file` (
    `id` bigint(20) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `group_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '组 id',
    `file_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '文件id',
    `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
    `upload_at` datetime DEFAULT NOW() COMMENT '上传时间',
    `last_update` datetime DEFAULT NOW() ON UPDATE NOW() COMMENT '最后修改时间',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '文件状态(0正常1已删除2禁用)',
    UNIQUE KEY `idx_group_file` (`group_id`, `file_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;