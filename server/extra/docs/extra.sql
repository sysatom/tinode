/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE='NO_AUTO_VALUE_ON_ZERO', SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# auth
# ------------------------------------------------------------

CREATE TABLE `auth` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uname` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `userid` bigint NOT NULL,
  `scheme` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
  `authlvl` int NOT NULL,
  `secret` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `expires` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `auth_userid_scheme` (`userid`,`scheme`),
  UNIQUE KEY `auth_uname` (`uname`),
  CONSTRAINT `auth_ibfk_1` FOREIGN KEY (`userid`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# chatbot_action
# ------------------------------------------------------------

CREATE TABLE `chatbot_action` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `topic` char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `seqid` int NOT NULL,
  `value` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `state` tinyint NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`,`topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# chatbot_behavior
# ------------------------------------------------------------

CREATE TABLE `chatbot_behavior` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `flag` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `count` int NOT NULL,
  `extra` json DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  KEY `flag` (`flag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# chatbot_configs
# ------------------------------------------------------------

CREATE TABLE `chatbot_configs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `topic` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` json NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`,`topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# chatbot_counters
# ------------------------------------------------------------

CREATE TABLE `chatbot_counters` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `topic` char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `flag` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `digit` bigint NOT NULL,
  `status` int NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`,`topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# chatbot_data
# ------------------------------------------------------------

CREATE TABLE `chatbot_data` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `topic` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` json NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`,`topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# chatbot_form
# ------------------------------------------------------------

CREATE TABLE `chatbot_form` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `form_id` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `uid` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `topic` char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `schema` json NOT NULL,
  `values` json NOT NULL,
  `state` tinyint NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `form_id` (`form_id`),
  KEY `uid` (`uid`,`topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# chatbot_key_result_values
# ------------------------------------------------------------

CREATE TABLE `chatbot_key_result_values` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `key_result_id` int DEFAULT NULL,
  `value` int NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `key_result_id` (`key_result_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# chatbot_key_results
# ------------------------------------------------------------

CREATE TABLE `chatbot_key_results` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `topic` char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `objective_id` bigint NOT NULL,
  `sequence` bigint NOT NULL,
  `title` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `memo` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `initial_value` int NOT NULL,
  `target_value` int NOT NULL,
  `current_value` int NOT NULL,
  `value_mode` tinyint NOT NULL,
  `tag` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`,`topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# chatbot_oauth
# ------------------------------------------------------------

CREATE TABLE `chatbot_oauth` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `topic` char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `token` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `extra` json NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`,`topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# chatbot_objectives
# ------------------------------------------------------------

CREATE TABLE `chatbot_objectives` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `topic` char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `sequence` bigint NOT NULL,
  `title` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `memo` varchar(1000) COLLATE utf8mb4_unicode_ci NOT NULL,
  `motive` varchar(1000) COLLATE utf8mb4_unicode_ci NOT NULL,
  `feasibility` varchar(1000) COLLATE utf8mb4_unicode_ci NOT NULL,
  `is_plan` tinyint NOT NULL,
  `plan_start` bigint NOT NULL,
  `plan_end` bigint NOT NULL,
  `total_value` int NOT NULL,
  `current_value` int NOT NULL,
  `tag` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_data` datetime NOT NULL,
  `updated_date` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`,`topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# chatbot_page
# ------------------------------------------------------------

CREATE TABLE `chatbot_page` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `page_id` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `uid` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `topic` char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `schema` json NOT NULL,
  `state` tinyint NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `page_id` (`page_id`),
  KEY `uid` (`uid`,`topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# chatbot_session
# ------------------------------------------------------------

CREATE TABLE `chatbot_session` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `topic` char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `rule_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `init` json NOT NULL,
  `values` json NOT NULL,
  `state` tinyint NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`,`topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# chatbot_todos
# ------------------------------------------------------------

CREATE TABLE `chatbot_todos` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `topic` char(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `sequence` bigint NOT NULL,
  `content` varchar(1000) COLLATE utf8mb4_unicode_ci NOT NULL,
  `category` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `remark` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `priority` bigint NOT NULL,
  `is_remind_at_time` tinyint NOT NULL,
  `remind_at` bigint NOT NULL,
  `repeat_method` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `repeat_rule` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `repeat_end_at` bigint NOT NULL,
  `complete` tinyint NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`,`topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# chatbot_url
# ------------------------------------------------------------

CREATE TABLE `chatbot_url` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `flag` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `url` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `state` tinyint NOT NULL,
  `view_count` int NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `flag` (`flag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# credentials
# ------------------------------------------------------------

CREATE TABLE `credentials` (
  `id` int NOT NULL AUTO_INCREMENT,
  `createdat` datetime(3) NOT NULL,
  `updatedat` datetime(3) NOT NULL,
  `deletedat` datetime(3) DEFAULT NULL,
  `method` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `synthetic` varchar(192) COLLATE utf8mb4_unicode_ci NOT NULL,
  `userid` bigint NOT NULL,
  `resp` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `done` tinyint NOT NULL DEFAULT '0',
  `retries` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `credentials_uniqueness` (`synthetic`),
  KEY `userid` (`userid`),
  CONSTRAINT `credentials_ibfk_1` FOREIGN KEY (`userid`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# dellog
# ------------------------------------------------------------

CREATE TABLE `dellog` (
  `id` int NOT NULL AUTO_INCREMENT,
  `topic` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `deletedfor` bigint NOT NULL DEFAULT '0',
  `delid` int NOT NULL,
  `low` int NOT NULL,
  `hi` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `dellog_topic_delid_deletedfor` (`topic`,`delid`,`deletedfor`),
  KEY `dellog_topic_deletedfor_low_hi` (`topic`,`deletedfor`,`low`,`hi`),
  KEY `dellog_deletedfor` (`deletedfor`),
  CONSTRAINT `dellog_ibfk_1` FOREIGN KEY (`topic`) REFERENCES `topics` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# devices
# ------------------------------------------------------------

CREATE TABLE `devices` (
  `id` int NOT NULL AUTO_INCREMENT,
  `userid` bigint NOT NULL,
  `hash` char(16) COLLATE utf8mb4_unicode_ci NOT NULL,
  `deviceid` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `platform` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `lastseen` datetime NOT NULL,
  `lang` varchar(8) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `devices_hash` (`hash`),
  KEY `userid` (`userid`),
  CONSTRAINT `devices_ibfk_1` FOREIGN KEY (`userid`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# filemsglinks
# ------------------------------------------------------------

CREATE TABLE `filemsglinks` (
  `id` int NOT NULL AUTO_INCREMENT,
  `createdat` datetime(3) NOT NULL,
  `fileid` bigint NOT NULL,
  `msgid` int DEFAULT NULL,
  `topic` char(25) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `userid` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fileid` (`fileid`),
  KEY `msgid` (`msgid`),
  KEY `topic` (`topic`),
  KEY `userid` (`userid`),
  CONSTRAINT `filemsglinks_ibfk_1` FOREIGN KEY (`fileid`) REFERENCES `fileuploads` (`id`) ON DELETE CASCADE,
  CONSTRAINT `filemsglinks_ibfk_2` FOREIGN KEY (`msgid`) REFERENCES `messages` (`id`) ON DELETE CASCADE,
  CONSTRAINT `filemsglinks_ibfk_3` FOREIGN KEY (`topic`) REFERENCES `topics` (`name`) ON DELETE CASCADE,
  CONSTRAINT `filemsglinks_ibfk_4` FOREIGN KEY (`userid`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# fileuploads
# ------------------------------------------------------------

CREATE TABLE `fileuploads` (
  `id` bigint NOT NULL,
  `createdat` datetime(3) NOT NULL,
  `updatedat` datetime(3) NOT NULL,
  `userid` bigint DEFAULT NULL,
  `status` int NOT NULL,
  `mimetype` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `size` bigint NOT NULL,
  `location` varchar(2048) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fileuploads_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# kvmeta
# ------------------------------------------------------------

CREATE TABLE `kvmeta` (
  `key` char(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` text COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# messages
# ------------------------------------------------------------

CREATE TABLE `messages` (
  `id` int NOT NULL AUTO_INCREMENT,
  `createdat` datetime(3) NOT NULL,
  `updatedat` datetime(3) NOT NULL,
  `deletedat` datetime(3) DEFAULT NULL,
  `delid` int DEFAULT '0',
  `seqid` int NOT NULL,
  `topic` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `from` bigint NOT NULL,
  `head` json DEFAULT NULL,
  `content` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `messages_topic_seqid` (`topic`,`seqid`),
  CONSTRAINT `messages_ibfk_1` FOREIGN KEY (`topic`) REFERENCES `topics` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# subscriptions
# ------------------------------------------------------------

CREATE TABLE `subscriptions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `createdat` datetime(3) NOT NULL,
  `updatedat` datetime(3) NOT NULL,
  `deletedat` datetime(3) DEFAULT NULL,
  `userid` bigint NOT NULL,
  `topic` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `delid` int DEFAULT '0',
  `recvseqid` int DEFAULT '0',
  `readseqid` int DEFAULT '0',
  `modewant` char(8) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `modegiven` char(8) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `private` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `subscriptions_topic_userid` (`topic`,`userid`),
  KEY `userid` (`userid`),
  KEY `subscriptions_topic` (`topic`),
  KEY `subscriptions_deletedat` (`deletedat`),
  CONSTRAINT `subscriptions_ibfk_1` FOREIGN KEY (`userid`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# topics
# ------------------------------------------------------------

CREATE TABLE `topics` (
  `id` int NOT NULL AUTO_INCREMENT,
  `createdat` datetime(3) NOT NULL,
  `updatedat` datetime(3) NOT NULL,
  `state` smallint NOT NULL DEFAULT '0',
  `stateat` datetime(3) DEFAULT NULL,
  `touchedat` datetime(3) DEFAULT NULL,
  `name` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `usebt` tinyint DEFAULT '0',
  `owner` bigint NOT NULL DEFAULT '0',
  `access` json DEFAULT NULL,
  `seqid` int NOT NULL DEFAULT '0',
  `delid` int DEFAULT '0',
  `public` json DEFAULT NULL,
  `trusted` json DEFAULT NULL,
  `tags` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `topics_name` (`name`),
  KEY `topics_owner` (`owner`),
  KEY `topics_state_stateat` (`state`,`stateat`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# topictags
# ------------------------------------------------------------

CREATE TABLE `topictags` (
  `id` int NOT NULL AUTO_INCREMENT,
  `topic` char(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `tag` varchar(96) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `topictags_userid_tag` (`topic`,`tag`),
  KEY `topictags_tag` (`tag`),
  CONSTRAINT `topictags_ibfk_1` FOREIGN KEY (`topic`) REFERENCES `topics` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# users
# ------------------------------------------------------------

CREATE TABLE `users` (
  `id` bigint NOT NULL,
  `createdat` datetime(3) NOT NULL,
  `updatedat` datetime(3) NOT NULL,
  `state` smallint NOT NULL DEFAULT '0',
  `stateat` datetime(3) DEFAULT NULL,
  `access` json DEFAULT NULL,
  `lastseen` datetime DEFAULT NULL,
  `useragent` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `public` json DEFAULT NULL,
  `trusted` json DEFAULT NULL,
  `tags` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `users_state_stateat` (`state`,`stateat`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# usertags
# ------------------------------------------------------------

CREATE TABLE `usertags` (
  `id` int NOT NULL AUTO_INCREMENT,
  `userid` bigint NOT NULL,
  `tag` varchar(96) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `usertags_userid_tag` (`userid`,`tag`),
  KEY `usertags_tag` (`tag`),
  CONSTRAINT `usertags_ibfk_1` FOREIGN KEY (`userid`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
