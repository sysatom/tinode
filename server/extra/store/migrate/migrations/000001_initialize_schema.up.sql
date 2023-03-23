CREATE TABLE `chatbot_action`
(
    `id`         int unsigned                                              NOT NULL AUTO_INCREMENT,
    `uid`        char(25) COLLATE utf8mb4_unicode_ci                       NOT NULL,
    `topic`      char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `seqid`      int                                                       NOT NULL,
    `value`      varchar(256) COLLATE utf8mb4_unicode_ci                   NOT NULL,
    `state`      tinyint                                                   NOT NULL,
    `created_at` datetime                                                  NOT NULL,
    `updated_at` datetime                                                  NOT NULL,
    PRIMARY KEY (`id`),
    KEY `uid` (`uid`, `topic`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE `chatbot_behavior`
(
    `id`         int unsigned                            NOT NULL AUTO_INCREMENT,
    `uid`        char(25) COLLATE utf8mb4_unicode_ci     NOT NULL,
    `flag`       varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `count`      int                                     NOT NULL,
    `extra`      json DEFAULT NULL,
    `created_at` datetime                                NOT NULL,
    `updated_at` datetime                                NOT NULL,
    PRIMARY KEY (`id`),
    KEY `uid` (`uid`),
    KEY `flag` (`flag`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE `chatbot_configs`
(
    `id`         int unsigned                            NOT NULL AUTO_INCREMENT,
    `uid`        char(25) COLLATE utf8mb4_unicode_ci     NOT NULL,
    `topic`      char(25) COLLATE utf8mb4_unicode_ci     NOT NULL,
    `key`        varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `value`      json                                    NOT NULL,
    `created_at` datetime                                NOT NULL,
    `updated_at` datetime                                NOT NULL,
    PRIMARY KEY (`id`),
    KEY `uid` (`uid`, `topic`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE `chatbot_counters`
(
    `id`         int unsigned                                              NOT NULL AUTO_INCREMENT,
    `uid`        char(25) COLLATE utf8mb4_unicode_ci                       NOT NULL,
    `topic`      char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `flag`       varchar(100) COLLATE utf8mb4_unicode_ci                   NOT NULL,
    `digit`      bigint                                                    NOT NULL,
    `status`     int                                                       NOT NULL,
    `created_at` datetime                                                  NOT NULL,
    `updated_at` datetime                                                  NOT NULL,
    PRIMARY KEY (`id`),
    KEY `uid` (`uid`, `topic`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE `chatbot_data`
(
    `id`         int unsigned                            NOT NULL AUTO_INCREMENT,
    `uid`        char(25) COLLATE utf8mb4_unicode_ci     NOT NULL,
    `topic`      char(25) COLLATE utf8mb4_unicode_ci     NOT NULL,
    `key`        varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `value`      json                                    NOT NULL,
    `created_at` datetime                                NOT NULL,
    `updated_at` datetime                                NOT NULL,
    PRIMARY KEY (`id`),
    KEY `uid` (`uid`, `topic`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE `chatbot_form`
(
    `id`         int unsigned                                              NOT NULL AUTO_INCREMENT,
    `form_id`    varchar(100) COLLATE utf8mb4_unicode_ci                   NOT NULL,
    `uid`        char(25) COLLATE utf8mb4_unicode_ci                       NOT NULL,
    `topic`      char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `schema`     json                                                      NOT NULL,
    `values`     json                                                      NOT NULL,
    `state`      tinyint                                                   NOT NULL,
    `created_at` datetime                                                  NOT NULL,
    `updated_at` datetime                                                  NOT NULL,
    PRIMARY KEY (`id`),
    KEY `form_id` (`form_id`),
    KEY `uid` (`uid`, `topic`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE `chatbot_key_result_values`
(
    `id`            int unsigned NOT NULL AUTO_INCREMENT,
    `key_result_id` int DEFAULT NULL,
    `value`         int          NOT NULL,
    `created_at`    datetime     NOT NULL,
    `updated_at`    datetime     NOT NULL,
    PRIMARY KEY (`id`),
    KEY `key_result_id` (`key_result_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE `chatbot_key_results`
(
    `id`            int unsigned                                                   NOT NULL AUTO_INCREMENT,
    `uid`           char(25) COLLATE utf8mb4_unicode_ci                            NOT NULL,
    `topic`         char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci      NOT NULL,
    `objective_id`  bigint                                                         NOT NULL,
    `sequence`      bigint                                                         NOT NULL,
    `title`         varchar(100) COLLATE utf8mb4_unicode_ci                        NOT NULL,
    `memo`          varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `initial_value` int                                                            NOT NULL,
    `target_value`  int                                                            NOT NULL,
    `current_value` int                                                            NOT NULL,
    `value_mode`    tinyint                                                        NOT NULL,
    `tag`           varchar(100) COLLATE utf8mb4_unicode_ci                        NOT NULL,
    `created_at`    datetime                                                       NOT NULL,
    `updated_at`    datetime                                                       NOT NULL,
    PRIMARY KEY (`id`),
    KEY `uid` (`uid`, `topic`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE `chatbot_oauth`
(
    `id`         int unsigned                                              NOT NULL AUTO_INCREMENT,
    `uid`        char(25) COLLATE utf8mb4_unicode_ci                       NOT NULL,
    `topic`      char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `name`       varchar(100) COLLATE utf8mb4_unicode_ci                   NOT NULL,
    `type`       varchar(50) COLLATE utf8mb4_unicode_ci                    NOT NULL,
    `token`      varchar(256) COLLATE utf8mb4_unicode_ci                   NOT NULL,
    `extra`      json                                                      NOT NULL,
    `created_at` datetime                                                  NOT NULL,
    `updated_at` datetime                                                  NOT NULL,
    PRIMARY KEY (`id`),
    KEY `uid` (`uid`, `topic`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE `chatbot_objectives`
(
    `id`            int unsigned                                                  NOT NULL AUTO_INCREMENT,
    `uid`           char(25) COLLATE utf8mb4_unicode_ci                           NOT NULL,
    `topic`         char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci     NOT NULL,
    `sequence`      bigint                                                        NOT NULL,
    `title`         varchar(100) COLLATE utf8mb4_unicode_ci                       NOT NULL,
    `memo`          varchar(1000) COLLATE utf8mb4_unicode_ci                      NOT NULL,
    `motive`        varchar(1000) COLLATE utf8mb4_unicode_ci                      NOT NULL,
    `feasibility`   varchar(1000) COLLATE utf8mb4_unicode_ci                      NOT NULL,
    `is_plan`       tinyint                                                       NOT NULL,
    `plan_start`    bigint                                                        NOT NULL,
    `plan_end`      bigint                                                        NOT NULL,
    `total_value`   int                                                           NOT NULL,
    `current_value` int                                                           NOT NULL,
    `tag`           varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_data`  datetime                                                      NOT NULL,
    `updated_date`  datetime                                                      NOT NULL,
    PRIMARY KEY (`id`),
    KEY `uid` (`uid`, `topic`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE `chatbot_page`
(
    `id`         int unsigned                                              NOT NULL AUTO_INCREMENT,
    `page_id`    varchar(100) COLLATE utf8mb4_unicode_ci                   NOT NULL,
    `uid`        char(25) COLLATE utf8mb4_unicode_ci                       NOT NULL,
    `topic`      char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `type`       varchar(100) COLLATE utf8mb4_unicode_ci                   NOT NULL,
    `schema`     json                                                      NOT NULL,
    `state`      tinyint                                                   NOT NULL,
    `created_at` datetime                                                  NOT NULL,
    `updated_at` datetime                                                  NOT NULL,
    PRIMARY KEY (`id`),
    KEY `page_id` (`page_id`),
    KEY `uid` (`uid`, `topic`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE `chatbot_session`
(
    `id`         int unsigned                                                  NOT NULL AUTO_INCREMENT,
    `uid`        char(25) COLLATE utf8mb4_unicode_ci                           NOT NULL,
    `topic`      char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci     NOT NULL,
    `rule_id`    varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `init`       json                                                          NOT NULL,
    `values`     json                                                          NOT NULL,
    `state`      tinyint                                                       NOT NULL,
    `created_at` datetime                                                      NOT NULL,
    `updated_at` datetime                                                      NOT NULL,
    PRIMARY KEY (`id`),
    KEY `uid` (`uid`, `topic`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE `chatbot_todos`
(
    `id`                int unsigned                                                  NOT NULL AUTO_INCREMENT,
    `uid`               char(25) COLLATE utf8mb4_unicode_ci                           NOT NULL,
    `topic`             char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci     NOT NULL,
    `sequence`          bigint                                                        NOT NULL,
    `content`           varchar(1000) COLLATE utf8mb4_unicode_ci                      NOT NULL,
    `category`          varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `remark`            varchar(100) COLLATE utf8mb4_unicode_ci                       NOT NULL,
    `priority`          bigint                                                        NOT NULL,
    `is_remind_at_time` tinyint                                                       NOT NULL,
    `remind_at`         bigint                                                        NOT NULL,
    `repeat_method`     varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `repeat_rule`       varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `repeat_end_at`     bigint                                                        NOT NULL,
    `complete`          tinyint                                                       NOT NULL,
    `created_at`        datetime                                                      NOT NULL,
    `updated_at`        datetime                                                      NOT NULL,
    PRIMARY KEY (`id`),
    KEY `uid` (`uid`, `topic`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE `chatbot_url`
(
    `id`         int unsigned                            NOT NULL AUTO_INCREMENT,
    `flag`       varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `url`        varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
    `state`      tinyint                                 NOT NULL,
    `view_count` int                                     NOT NULL DEFAULT '0',
    `created_at` datetime                                NOT NULL,
    `updated_at` datetime                                NOT NULL,
    PRIMARY KEY (`id`),
    KEY `flag` (`flag`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE `chatbot_counter_records`
(
    `counter_id` int unsigned NOT NULL AUTO_INCREMENT,
    `digit`      int          NOT NULL,
    `created_at` datetime     NOT NULL,
    PRIMARY KEY (`counter_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE `chatbot_instruct`
(
    `id`         int unsigned                                                 NOT NULL AUTO_INCREMENT,
    `no`         char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci    NOT NULL,
    `uid`        char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci    NOT NULL,
    `object`     varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `bot`        varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `flag`       varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `content`    json                                                         NOT NULL,
    `priority`   int                                                          NOT NULL,
    `state`      tinyint                                                      NOT NULL,
    `expire_at`  datetime                                                     NOT NULL,
    `created_at` datetime                                                     NOT NULL,
    `updated_at` datetime                                                     NOT NULL,
    PRIMARY KEY (`id`),
    KEY `uid` (`uid`),
    KEY `no` (`no`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
