DROP TABLE IF EXISTS `account_info`;
CREATE TABLE `account_info`
(
    `id`               bigint(20) NOT NULL AUTO_INCREMENT,
    `account_name`     varchar(100) CHARACTER SET utf8 COLLATE utf8_bin NULL DEFAULT NULL COMMENT '户
主姓名',
    `account_no`       varchar(100) CHARACTER SET utf8 COLLATE utf8_bin NULL DEFAULT NULL COMMENT '银行
卡号',
    `account_password` varchar(100) CHARACTER SET utf8 COLLATE utf8_bin NULL DEFAULT NULL COMMENT
        '帐户密码',
    `account_balance`  double NULL DEFAULT NULL COMMENT '帐户余额',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8 COLLATE = utf8_bin ROW_FORMAT =
Dynamic;
INSERT INTO `account_info`
VALUES (2, '张三的账户', '1', '', 10000);

-- 分别在bank1、bank2库中创建undo_log表，此表为seata框架使用：
DROP TABLE IF EXISTS `undo_log`;
CREATE TABLE `undo_log`
(
    `id`            bigint(20) NOT NULL AUTO_INCREMENT,
    `branch_id`     bigint(20) NOT NULL,
    `xid`           varchar(128)                                            NOT NULL,
    `context`       varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
    `rollback_info` longblob                                                NOT NULL,
    `log_status`    int(11) NOT NULL,
    `log_created`   datetime                                                NOT NULL,
    `log_modified`  datetime                                                NOT NULL,
    `ext`           varchar(100) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY             `idx_unionkey` (`xid`,`branch_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

SET
FOREIGN_KEY_CHECKS = 1;