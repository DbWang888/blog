CREATE TABLE `blog_tag` (
  `id` INT(10) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) DEFAULT "" COMMENT '标签名称',
  `created_on` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` VARCHAR(100) DEFAULT "" COMMENT '创建人',
  `modified_on` DATETIME NOT NULL DEFAULT "1000-01-01 00:00:00" COMMENT '修改时间',
  `modified_by` VARCHAR(100) DEFAULT "" COMMENT '修改人',
  `deleted_on` DATETIME NOT NULL DEFAULT "1000-01-01 00:00:00",
  `state` TINYINT(3) DEFAULT "1" COMMENT '状态 0为禁用、1为启用'
);

CREATE TABLE `blog_article` (
  `id` INT(10) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `tag_id` INT(10) DEFAULT "0" COMMENT '标签ID',
  `title` VARCHAR(100) DEFAULT "" COMMENT '文章标题',
  `desc` VARCHAR(255) DEFAULT "" COMMENT '简述',
  `content` TEXT,
  `created_on` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` VARCHAR(100) DEFAULT "" COMMENT '创建人',
  `modified_on` DATETIME NOT NULL DEFAULT "1000-01-01 00:00:00" COMMENT '修改时间',
  `modified_by` VARCHAR(255) DEFAULT "" COMMENT '修改人',
  `deleted_on` DATETIME NOT NULL DEFAULT "1000-01-01 00:00:00",
  `state` TINYINT(3) DEFAULT "1" COMMENT '状态 0为禁用1为启用'
);

CREATE TABLE `blog_auth` (
  `id` INT(10) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(50) UNIQUE DEFAULT "" COMMENT '账号',
  `password` VARCHAR(50) DEFAULT "" COMMENT '密码',
  `created_on` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
);

ALTER TABLE `blog_tag` COMMENT = '文章标签管理';

ALTER TABLE `blog_article` COMMENT = '文章管理';



ALTER TABLE blog_article ADD FOREIGN KEY (tag_id) REFERENCES blog_tag(id);
ALTER TABLE blog_tag ADD UNIQUE KEY(`name`, created_by);
