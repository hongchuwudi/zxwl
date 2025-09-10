-- MySQL dump 10.13  Distrib 8.0.41, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: zxwl
-- ------------------------------------------------------
-- Server version	8.0.41

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `admission_batches`
--

DROP TABLE IF EXISTS `admission_batches`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admission_batches` (
  `id` int NOT NULL COMMENT '主键id',
  `name` varchar(100) NOT NULL COMMENT '批次名称',
  `batch_info` varchar(100) DEFAULT NULL COMMENT '批次信息',
  `batch_type` varchar(50) DEFAULT NULL COMMENT '批次大类名称',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='录取类型';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admission_category_mapping`
--

DROP TABLE IF EXISTS `admission_category_mapping`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admission_category_mapping` (
  `id` int NOT NULL AUTO_INCREMENT,
  `old_category_id` int NOT NULL COMMENT '源类型ID',
  `new_category_id` int NOT NULL COMMENT '目标类型ID',
  `province_id` int unsigned DEFAULT NULL COMMENT '省份ID，如果为NULL则表示适用于所有省份',
  `start_year` int DEFAULT NULL COMMENT '开始年份',
  `end_year` int DEFAULT NULL COMMENT '结束年份',
  `mapping_comment` varchar(25) DEFAULT NULL COMMENT '映射关系注释',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_mapping` (`old_category_id`,`new_category_id`,`province_id`,`start_year`,`end_year`)
) ENGINE=InnoDB AUTO_INCREMENT=53 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='高考类型映射表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admission_special`
--

DROP TABLE IF EXISTS `admission_special`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admission_special` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `year` smallint NOT NULL COMMENT '年份(1900-2155)',
  `school_id` int NOT NULL COMMENT '学校ID',
  `province_id` int NOT NULL COMMENT '省份ID',
  `special_id` int DEFAULT NULL COMMENT '专业表ID',
  `type` int NOT NULL COMMENT '高考类型ID',
  `batch` int NOT NULL COMMENT '批次类型ID',
  `zslx` int DEFAULT NULL COMMENT '招生类型ID',
  `max_score` smallint unsigned DEFAULT NULL COMMENT '最高分(0-65535)',
  `min_score` smallint unsigned DEFAULT NULL COMMENT '最低分(0-65535)',
  `min_rank` int DEFAULT NULL COMMENT '最低位次',
  `max_rank` int DEFAULT NULL COMMENT '最大位次',
  `min_section` varchar(20) DEFAULT NULL COMMENT '最低位次',
  `admission_count` int unsigned DEFAULT NULL COMMENT '录取人数',
  `average_score` smallint unsigned DEFAULT NULL COMMENT '平均分(0-65535)',
  `batch_diff` smallint DEFAULT NULL COMMENT '批次线差(-32768~32767)',
  `special_info` varchar(500) DEFAULT NULL COMMENT '专业补充信息',
  `subject_requirement` text COMMENT '选科要求文本',
  `subject_codes` varchar(50) DEFAULT NULL COMMENT '选科科目代码',
  `remark` text COMMENT '备注',
  `score_range_flag` int DEFAULT '0' COMMENT '是否分数段(1是2否)',
  `sp_id` int NOT NULL COMMENT '专业ID',
  PRIMARY KEY (`id`,`year`),
  KEY `idx_main_query` (`year`,`province_id`,`type`,`batch`) COMMENT '主要查询组合索引',
  KEY `idx_score_range` (`min_score`,`max_score`) COMMENT '分数范围索引',
  KEY `idx_zslx_batch` (`zslx`,`batch`) COMMENT '招生类型-批次索引',
  KEY `idx_admission_count` (`admission_count`) COMMENT '录取人数索引',
  KEY `idx_id` (`id`),
  KEY `idx_school_special` (`school_id`,`special_id`) COMMENT '学校-专业联合索引',
  KEY `idx_province_type` (`province_id`,`type`)
) ENGINE=InnoDB AUTO_INCREMENT=736273 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='高校招生专业信息表(优化版)'
/*!50100 PARTITION BY RANGE (`year`)
(PARTITION p2020 VALUES LESS THAN (2021) ENGINE = InnoDB,
 PARTITION p2021 VALUES LESS THAN (2022) ENGINE = InnoDB,
 PARTITION p2022 VALUES LESS THAN (2023) ENGINE = InnoDB,
 PARTITION p2023 VALUES LESS THAN (2024) ENGINE = InnoDB,
 PARTITION p2024 VALUES LESS THAN (2025) ENGINE = InnoDB,
 PARTITION p2025 VALUES LESS THAN (2026) ENGINE = InnoDB,
 PARTITION p2026 VALUES LESS THAN (2027) ENGINE = InnoDB,
 PARTITION p2027 VALUES LESS THAN (2028) ENGINE = InnoDB,
 PARTITION p_max VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admission_types`
--

DROP TABLE IF EXISTS `admission_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admission_types` (
  `id` int NOT NULL COMMENT '主键id',
  `name` varchar(100) NOT NULL COMMENT '类型名称',
  `type_info` varchar(100) DEFAULT NULL COMMENT '类型信息',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='高考类型';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admission_universities`
--

DROP TABLE IF EXISTS `admission_universities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admission_universities` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `school_id` int DEFAULT NULL COMMENT '关联的学校ID',
  `province_id` int DEFAULT NULL COMMENT '省份ID',
  `year` int DEFAULT NULL COMMENT '年份（如2024）',
  `score_type` varchar(10) DEFAULT NULL COMMENT '高考类型(1文科,2理科...逻辑外键外表）',
  `min_score` int DEFAULT NULL COMMENT '最低录取分数',
  PRIMARY KEY (`id`),
  KEY `idx_school_id` (`school_id`),
  KEY `idx_province_year` (`province_id`,`year`),
  KEY `idx_score_type` (`score_type`),
  KEY `idx_score_school` (`school_id`),
  KEY `idx_score_province` (`province_id`),
  KEY `idx_score_year` (`year`)
) ENGINE=InnoDB AUTO_INCREMENT=231457 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='省份录取分数线表，记录各省份历年录取最低分';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admission_zhaosheng_type`
--

DROP TABLE IF EXISTS `admission_zhaosheng_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admission_zhaosheng_type` (
  `id` int NOT NULL COMMENT '主键id',
  `name` varchar(100) NOT NULL COMMENT '招生类型名称',
  `zslx_info` varchar(100) DEFAULT NULL COMMENT '招生类型信息',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='招生类型';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `chat_messages`
--

DROP TABLE IF EXISTS `chat_messages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chat_messages` (
  `id` int NOT NULL AUTO_INCREMENT,
  `school_id` int NOT NULL,
  `email` varchar(255) NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_school_time` (`school_id`,`created_at`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `college`
--

DROP TABLE IF EXISTS `college`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `college` (
  `name` varchar(255) DEFAULT NULL,
  `label` varchar(255) DEFAULT NULL,
  `location` varchar(255) DEFAULT NULL,
  `id` int NOT NULL AUTO_INCREMENT,
  `badge` varchar(255) DEFAULT NULL,
  `motto` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `colleges`
--

DROP TABLE IF EXISTS `colleges`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `colleges` (
  `belong` varchar(255) DEFAULT NULL COMMENT '隶属于什么部门',
  `city_name` varchar(255) DEFAULT NULL COMMENT '学校所在位置',
  `dual_class_name` varchar(255) DEFAULT NULL COMMENT '双一流',
  `f985` int DEFAULT NULL COMMENT '是否是985',
  `f211` int DEFAULT NULL COMMENT '是否是211',
  `hightitle` varchar(255) DEFAULT NULL COMMENT '大学名',
  `level_name` varchar(255) DEFAULT NULL COMMENT '是否是本科',
  `nature_name` varchar(255) DEFAULT NULL COMMENT '公办?',
  `paiming` varchar(255) NOT NULL COMMENT '排名',
  `type_name` varchar(255) DEFAULT NULL COMMENT '综合类',
  `isbiaoshi` int(10) unsigned zerofill NOT NULL DEFAULT '0000000000',
  `school_id` int DEFAULT NULL COMMENT '照片',
  `id` int NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_colleges_school_id` (`school_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2969 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `commom_cities`
--

DROP TABLE IF EXISTS `commom_cities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `commom_cities` (
  `id` int NOT NULL COMMENT '市级行政区划代码(前4位)',
  `city_name` varchar(50) NOT NULL COMMENT '城市名称',
  `province_id` int NOT NULL COMMENT '所属省级代码(前2位)',
  `city_level` tinyint DEFAULT NULL COMMENT '城市级别(1:直辖市,2:副省级城市,3:地级市,4:县级市)',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='市级行政区划表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `common_districts`
--

DROP TABLE IF EXISTS `common_districts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `common_districts` (
  `id` int NOT NULL COMMENT '区县级行政区划代码(全6位)',
  `district_name` varchar(50) NOT NULL COMMENT '区县名称',
  `city_id` int NOT NULL COMMENT '所属市级代码(前4位)',
  `district_type` tinyint DEFAULT NULL COMMENT '区县类型(1:市辖区,2:县,3:县级市,4:自治县,5:旗)',
  `is_urban` tinyint(1) DEFAULT '0' COMMENT '是否城区',
  PRIMARY KEY (`id`),
  KEY `idx_city_id` (`city_id`) COMMENT '城市ID索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='区县级行政区划表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `common_provinces`
--

DROP TABLE IF EXISTS `common_provinces`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `common_provinces` (
  `id` int unsigned NOT NULL COMMENT '省份ID',
  `name` varchar(50) NOT NULL COMMENT '省份全称',
  `short_name` varchar(10) NOT NULL COMMENT '省份简称',
  `code` varchar(6) NOT NULL COMMENT '行政区划代码',
  `pinyin` varchar(100) NOT NULL COMMENT '拼音',
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='中国省份信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `familyshare`
--

DROP TABLE IF EXISTS `familyshare`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `familyshare` (
  `familyemail` varchar(255) NOT NULL,
  `myemail` varchar(255) NOT NULL,
  `id` int NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `logtextuser`
--

DROP TABLE IF EXISTS `logtextuser`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `logtextuser` (
  `email` varchar(255) DEFAULT NULL,
  `date` datetime DEFAULT NULL,
  `operation` varchar(255) DEFAULT NULL,
  `id` int NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=488 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `moni`
--

DROP TABLE IF EXISTS `moni`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `moni` (
  `college` varchar(255) DEFAULT NULL,
  `benke` int DEFAULT NULL,
  `major1` varchar(255) DEFAULT NULL,
  `major2` varchar(255) DEFAULT NULL,
  `major3` varchar(255) DEFAULT NULL,
  `major4` varchar(255) DEFAULT NULL,
  `major5` varchar(255) DEFAULT NULL,
  `major6` varchar(255) DEFAULT NULL,
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(255) DEFAULT NULL,
  `yitianzhuanye` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `news_comments`
--

DROP TABLE IF EXISTS `news_comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `news_comments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '评论ID，主键自增',
  `news_id` bigint unsigned NOT NULL COMMENT '关联的资讯ID',
  `parent_id` bigint unsigned DEFAULT '0' COMMENT '父评论ID（0表示顶级评论）',
  `commenter_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '评论者名称',
  `commenter_email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '评论者邮箱',
  `comment_content` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '评论内容',
  `comment_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '评论时间',
  `is_approved` tinyint(1) DEFAULT '0' COMMENT '审核状态：0-待审核，1-已审核',
  `like_count` int DEFAULT '0' COMMENT '点赞数',
  `reply_count` int DEFAULT '0' COMMENT '回复数',
  `ip_address` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '评论者IP地址',
  `user_agent` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '用户浏览器信息',
  PRIMARY KEY (`id`),
  KEY `idx_news_id` (`news_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_comment_time` (`comment_time`),
  KEY `idx_commenter_email` (`commenter_email`),
  KEY `idx_is_approved` (`is_approved`),
  CONSTRAINT `fk_news_comments_news` FOREIGN KEY (`news_id`) REFERENCES `news_info` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=95 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='资讯评论表';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
/*!50032 DROP TRIGGER IF EXISTS after_comment_insert */;
DELIMITER ;;
/*!50003 CREATE*/ /*!50017 DEFINER=`root`@`localhost`*/ /*!50003 TRIGGER `after_comment_insert` AFTER INSERT ON `news_comments` FOR EACH ROW BEGIN
    -- 更新资讯的评论总数（+1）
    UPDATE news_info
    SET comment_count = comment_count + 1
    WHERE id = NEW.news_id;
END */;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
/*!50032 DROP TRIGGER IF EXISTS after_comment_delete */;
DELIMITER ;;
/*!50003 CREATE*/ /*!50017 DEFINER=`root`@`localhost`*/ /*!50003 TRIGGER `after_comment_delete` AFTER DELETE ON `news_comments` FOR EACH ROW BEGIN
    -- 更新资讯的评论总数（-1）
    UPDATE news_info
    SET comment_count = GREATEST(comment_count - 1, 0)
    WHERE id = OLD.news_id;
END */;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;

--
-- Table structure for table `news_info`
--

DROP TABLE IF EXISTS `news_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `news_info` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID，自增长',
  `province_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '省份ID，来自API的province_id字段',
  `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资讯标题',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '资讯描述/摘要',
  `keywords` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '关键词',
  `content` longtext COLLATE utf8mb4_unicode_ci COMMENT '资讯完整内容（HTML格式）',
  `video_detail` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '视频详情',
  `video_type` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '视频类型',
  `video_img` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '视频图片',
  `from_source` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '资讯来源（简略）',
  `news_num` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '资讯编号或阅读量，来自API的num字段',
  `is_push` tinyint(1) DEFAULT '2' COMMENT '是否推送：1-是，2-否',
  `is_top` tinyint(1) DEFAULT '1' COMMENT '置顶状态：2-置顶，1-普通',
  `style_type` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '样式类型，来自API的style.type字段',
  `style_url` longtext COLLATE utf8mb4_unicode_ci COMMENT '样式图片URL，来自API的style.url字段',
  `publish_time` timestamp NULL DEFAULT NULL COMMENT '发布时间，来自API的times字段（Unix时间戳转换）',
  `add_time` timestamp NULL DEFAULT NULL COMMENT '添加时间',
  `card_school_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '关联学校ID',
  `card_live_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '关联直播ID',
  `class_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '分类名称',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录最后更新时间',
  `publisher_email` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '创建者名称',
  `like_count` int DEFAULT '0' COMMENT '点赞数',
  `share_count` int DEFAULT '0' COMMENT '分享数',
  `favorite_count` int DEFAULT '0' COMMENT '收藏数',
  `comment_count` int DEFAULT '0' COMMENT '评论总数',
  PRIMARY KEY (`id`),
  KEY `idx_province_id` (`province_id`),
  KEY `idx_publish_time` (`publish_time`),
  KEY `idx_is_top` (`is_top`),
  KEY `idx_is_push` (`is_push`)
) ENGINE=InnoDB AUTO_INCREMENT=282 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='资讯数据表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `policy`
--

DROP TABLE IF EXISTS `policy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `policy` (
  `title` varchar(255) DEFAULT NULL,
  `contenet` text,
  `id` int NOT NULL AUTO_INCREMENT,
  `foreword` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `professional`
--

DROP TABLE IF EXISTS `professional`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `professional` (
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `level1_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `salaryavg` int DEFAULT NULL,
  `limit_year` varchar(255) DEFAULT NULL,
  `fivesalaryavg` int DEFAULT NULL,
  `level2_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `boy_rate` varchar(255) DEFAULT NULL,
  `girl_rate` varchar(255) DEFAULT NULL,
  `level3_name` varchar(255) DEFAULT NULL,
  `id` int NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1886 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `profile`
--

DROP TABLE IF EXISTS `profile`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `profile` (
  `address` varchar(255) DEFAULT NULL,
  `graduate` year DEFAULT NULL,
  `picture` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `sex` int DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `id` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `name` (`name`),
  KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sch_profile`
--

DROP TABLE IF EXISTS `sch_profile`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sch_profile` (
  `school_id` int NOT NULL,
  `introduce` text,
  `hightitle` varchar(255) DEFAULT NULL,
  `gong_major` varchar(255) DEFAULT NULL,
  `li_major` varchar(255) DEFAULT NULL,
  `jing_major` varchar(255) DEFAULT NULL,
  `guan_major` varchar(255) DEFAULT NULL,
  `yi_major` varchar(255) DEFAULT NULL,
  `fa_major` varchar(255) DEFAULT NULL,
  `jiao_major` varchar(255) DEFAULT NULL,
  `yishu_major` varchar(255) DEFAULT NULL,
  `picture` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`school_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `score_exam_info`
--

DROP TABLE IF EXISTS `score_exam_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `score_exam_info` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `exam_year` smallint unsigned NOT NULL COMMENT '考试年份',
  `exam_name` varchar(88) NOT NULL COMMENT '考试名称(如"陕西-2025-文科")',
  `province_id` int unsigned NOT NULL COMMENT '省份ID,参考省份表',
  `type_id` smallint NOT NULL COMMENT '考试类型(1:理科,2:文科,3:综合...),参考录取类型表',
  `batch_id` int NOT NULL COMMENT '录取批次,参考录取批次表',
  `full_score` decimal(5,1) unsigned NOT NULL COMMENT '满分分值',
  `total_num` int unsigned NOT NULL COMMENT '高考总人数',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_exam_name_type_year` (`exam_name`,`exam_year`,`type_id`,`batch_id`),
  KEY `idx_year_type` (`exam_year`,`type_id`),
  KEY `idx_year_batch` (`exam_year`,`batch_id`)
) ENGINE=InnoDB AUTO_INCREMENT=332 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='考试基本信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `score_section`
--

DROP TABLE IF EXISTS `score_section`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `score_section` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `exam_year` smallint unsigned NOT NULL COMMENT '考试年份',
  `batch_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '批次名称(如"本科一批")',
  `exam_id` int NOT NULL COMMENT '关联考试ID',
  `score_range` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分数范围(冗余字段，如"600-750")',
  `min_score` decimal(5,1) unsigned NOT NULL COMMENT '最低分(包含)',
  `max_score` decimal(5,1) unsigned NOT NULL COMMENT '最高分(包含)',
  `rank_range` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '排名范围(冗余字段，如"100-150")',
  `rank_start` int unsigned NOT NULL COMMENT '排名起始(包含)',
  `rank_end` int unsigned NOT NULL COMMENT '排名结束(包含)',
  `total_students` mediumint unsigned NOT NULL COMMENT '该分段总人数',
  PRIMARY KEY (`id`,`exam_year`),
  KEY `idx_exam_year` (`exam_year`),
  KEY `idx_score_range` (`min_score`,`max_score`),
  KEY `idx_batch_score` (`exam_id`,`min_score`),
  KEY `idx_year_batch` (`exam_year`,`batch_name`),
  KEY `idx_cover_query` (`exam_year`,`batch_name`,`min_score`,`max_score`,`rank_start`,`rank_end`)
) ENGINE=InnoDB AUTO_INCREMENT=170916 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='一分一段主表'
/*!50100 PARTITION BY RANGE (`exam_year`)
(PARTITION p2019 VALUES LESS THAN (2020) ENGINE = InnoDB,
 PARTITION p2020 VALUES LESS THAN (2021) ENGINE = InnoDB,
 PARTITION p2021 VALUES LESS THAN (2022) ENGINE = InnoDB,
 PARTITION p2022 VALUES LESS THAN (2023) ENGINE = InnoDB,
 PARTITION p2023 VALUES LESS THAN (2024) ENGINE = InnoDB,
 PARTITION p2024 VALUES LESS THAN (2025) ENGINE = InnoDB,
 PARTITION p2025 VALUES LESS THAN (2026) ENGINE = InnoDB,
 PARTITION p2026 VALUES LESS THAN (2027) ENGINE = InnoDB,
 PARTITION p_future VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `special_content`
--

DROP TABLE IF EXISTS `special_content`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `special_content` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `special_id` int NOT NULL COMMENT '专业ID',
  `content_type` tinyint NOT NULL COMMENT '1-专业介绍 2-就业方向 3-专业描述 4-主要课程',
  `content` longtext COMMENT 'HTML内容',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_special_id` (`special_id`),
  KEY `idx_content_type` (`content_type`)
) ENGINE=InnoDB AUTO_INCREMENT=7381 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='专业大文本内容表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `special_detail`
--

DROP TABLE IF EXISTS `special_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `special_detail` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(20) NOT NULL COMMENT '专业代码，如080901',
  `name` varchar(100) NOT NULL COMMENT '专业名称',
  `degree` varchar(100) DEFAULT NULL COMMENT '授予学位',
  `direction` varchar(500) DEFAULT NULL COMMENT '专业方向',
  `type` varchar(50) DEFAULT NULL COMMENT '学科门类，如工学',
  `type_detail` varchar(50) DEFAULT NULL COMMENT '专业类，如计算机类',
  `limit_year` varchar(12) DEFAULT NULL COMMENT '修业年限',
  `level1` int DEFAULT NULL COMMENT '一级分类代码',
  `level1_name` varchar(50) DEFAULT NULL COMMENT '分类1名称',
  `level2` int DEFAULT NULL COMMENT '二级分类代码',
  `level2_name` varchar(50) DEFAULT NULL COMMENT '分类2名称',
  `level3` int DEFAULT NULL COMMENT '三级分类代码',
  `level3_name` varchar(50) DEFAULT NULL COMMENT '分类3名称',
  `employment_rate` decimal(5,2) DEFAULT NULL COMMENT '就业比例',
  `avg_salary` int DEFAULT NULL COMMENT '平均薪资',
  `top_industry` varchar(100) DEFAULT NULL COMMENT '最多就业行业',
  `top_position` varchar(100) DEFAULT NULL COMMENT '最多就业岗位',
  `top_area` varchar(100) DEFAULT NULL COMMENT '最多就业地区',
  `monthly_views` varchar(20) DEFAULT NULL COMMENT '月浏览量',
  `total_views` varchar(20) DEFAULT NULL COMMENT '总浏览量',
  `subject_requirements` varchar(100) DEFAULT NULL COMMENT '选科建议',
  `gender_ratio` varchar(50) DEFAULT NULL COMMENT '男女比例',
  `celebrities` text COMMENT '知名校友',
  `courses` longtext COMMENT '主要课程内容',
  `content` longtext COMMENT '专业介绍内容',
  `career_prospects` longtext COMMENT '就业方向内容',
  `description` longtext COMMENT '专业描述内容',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_code` (`code`),
  KEY `idx_name` (`name`),
  KEY `idx_type` (`type`),
  KEY `idx_type_detail` (`type_detail`)
) ENGINE=InnoDB AUTO_INCREMENT=10049 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='专业基本信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `special_employment_rate`
--

DROP TABLE IF EXISTS `special_employment_rate`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `special_employment_rate` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `special_id` varchar(32) NOT NULL COMMENT '专业ID',
  `year` varchar(10) NOT NULL COMMENT '年份',
  `rate` varchar(20) NOT NULL COMMENT '就业率范围',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_special_id` (`special_id`),
  KEY `idx_year` (`year`)
) ENGINE=InnoDB AUTO_INCREMENT=2513 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='专业就业率表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `special_famous_school`
--

DROP TABLE IF EXISTS `special_famous_school`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `special_famous_school` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `special_id` varchar(32) NOT NULL COMMENT '专业ID',
  `school_name` varchar(100) NOT NULL COMMENT '学校名称',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_special_id` (`special_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5862 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='专业名校示例表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `special_impression_tag`
--

DROP TABLE IF EXISTS `special_impression_tag`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `special_impression_tag` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `special_id` varchar(32) NOT NULL COMMENT '专业ID',
  `image_url` varchar(255) NOT NULL COMMENT '图片URL(开头为:https://static-data.gaokao.cn/)',
  `keyword` varchar(50) NOT NULL COMMENT '关键词',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_special_id` (`special_id`),
  KEY `idx_keyword` (`keyword`)
) ENGINE=InnoDB AUTO_INCREMENT=4271 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='专业印象标签表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `special_job_distribution`
--

DROP TABLE IF EXISTS `special_job_distribution`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `special_job_distribution` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `special_id` varchar(32) NOT NULL COMMENT '专业ID',
  `distribution_type` tinyint NOT NULL COMMENT '1-行业 2-地区 3-岗位',
  `name` varchar(100) DEFAULT NULL COMMENT '名称(行业/地区)',
  `position` varchar(200) DEFAULT NULL COMMENT '岗位名称(仅type=3时使用)',
  `job_description` text COMMENT '工作描述(仅type=3时使用)',
  `rate` varchar(20) NOT NULL COMMENT '占比',
  `sort` int DEFAULT NULL COMMENT '排序',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_special_id` (`special_id`),
  KEY `idx_distribution_type` (`distribution_type`)
) ENGINE=InnoDB AUTO_INCREMENT=16424 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='专业就业分布表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `special_salary_data`
--

DROP TABLE IF EXISTS `special_salary_data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `special_salary_data` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `special_id` varchar(32) NOT NULL COMMENT '专业ID',
  `salary_type` tinyint NOT NULL COMMENT '1-专业薪资 2-所有专业薪资',
  `salary_year` tinyint NOT NULL COMMENT '1-毕业1年 2-毕业2年 3-毕业3年 4-毕业5年',
  `salary_value` int NOT NULL COMMENT '薪资数值',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_special_type_year` (`special_id`,`salary_type`,`salary_year`),
  KEY `idx_special_id` (`special_id`)
) ENGINE=InnoDB AUTO_INCREMENT=22625 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='专业薪资数据表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `special_video`
--

DROP TABLE IF EXISTS `special_video`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `special_video` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `special_id` int NOT NULL COMMENT '专业ID',
  `title` varchar(200) NOT NULL COMMENT '视频标题',
  `cover_image` varchar(255) DEFAULT NULL COMMENT '封面图URL',
  `video_url` varchar(255) NOT NULL COMMENT '视频URL',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `school_id` varchar(32) DEFAULT NULL COMMENT '学校ID',
  `school_special_id` varchar(32) DEFAULT NULL COMMENT '学校专业ID',
  `ranks` int DEFAULT NULL COMMENT '排名',
  `url_type` tinyint DEFAULT NULL COMMENT 'URL类型',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_special_id` (`special_id`),
  KEY `idx_school_id` (`school_id`),
  KEY `idx_rank` (`ranks`)
) ENGINE=InnoDB AUTO_INCREMENT=5732 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='专业视频表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `universities`
--

DROP TABLE IF EXISTS `universities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `universities` (
  `school_id` varchar(10) NOT NULL COMMENT '学校唯一标识符',
  `name` varchar(100) NOT NULL COMMENT '学校全称',
  `short` varchar(50) DEFAULT NULL COMMENT '学校简称',
  `address` text COMMENT '校区地址（多校区用逗号分隔）',
  `area` int DEFAULT NULL COMMENT '占地面积（单位：亩）',
  `belong` varchar(100) DEFAULT NULL COMMENT '隶属单位',
  `central` enum('1','2') DEFAULT NULL COMMENT '是否中央部属（1是 2否）',
  `city_id` char(6) DEFAULT NULL COMMENT '所在城市行政区划代码',
  `city_name` varchar(50) DEFAULT NULL COMMENT '所在城市名称',
  `code_enroll` varchar(20) DEFAULT NULL COMMENT '招生代码',
  `content` text COMMENT '学校简介',
  `dual_class_name` varchar(50) DEFAULT NULL COMMENT '双一流类别',
  `email` varchar(100) DEFAULT NULL COMMENT '联系邮箱',
  `f211` enum('1','2') DEFAULT NULL COMMENT '是否211（1是 2否）',
  `f985` enum('1','2') DEFAULT NULL COMMENT '是否985（1是 2否）',
  `gbh_url` varchar(255) DEFAULT NULL COMMENT '教育部高校信息链接',
  `last_updated` datetime DEFAULT NULL COMMENT '最后更新时间',
  `level_name` varchar(50) DEFAULT NULL COMMENT '办学层次（本科/专科）',
  `logo_src` varchar(255) DEFAULT NULL COMMENT '校徽图片URL',
  `motto` text COMMENT '校训',
  `nature_name` varchar(50) DEFAULT NULL COMMENT '办学性质（公办/民办）',
  `num_academician` int DEFAULT '0' COMMENT '院士数量',
  `num_doctor` int DEFAULT '0' COMMENT '博士点数量',
  `num_lab` int DEFAULT '0' COMMENT '实验室数量',
  `num_library` varchar(50) DEFAULT NULL COMMENT '藏书量（含单位）',
  `num_master` int DEFAULT '0' COMMENT '硕士点数量',
  `num_subject` int DEFAULT '0' COMMENT '重点学科数量',
  `phone` varchar(100) DEFAULT NULL COMMENT '联系电话（多个用逗号分隔）',
  `postcode` varchar(20) DEFAULT NULL COMMENT '邮政编码',
  `province_id` char(2) DEFAULT NULL COMMENT '省份代码',
  `province_name` varchar(50) DEFAULT NULL COMMENT '省份名称',
  `qs_rank` varchar(10) DEFAULT NULL COMMENT 'QS国内排名',
  `qs_world` varchar(10) DEFAULT NULL COMMENT 'QS世界排名',
  `recommend_master_rate` decimal(5,2) DEFAULT NULL COMMENT '推荐免试研究生比例',
  `ruanke_rank` varchar(10) DEFAULT NULL COMMENT '软科排名',
  `school_site` varchar(255) DEFAULT NULL COMMENT '学校官网',
  `school_special_num` int DEFAULT NULL COMMENT '开设专业数量',
  `school_type_name` varchar(50) DEFAULT NULL COMMENT '学校类型',
  `site` varchar(255) DEFAULT NULL COMMENT '招生网地址',
  `town_name` varchar(50) DEFAULT NULL COMMENT '所在区县名称',
  `type_name` varchar(50) DEFAULT NULL COMMENT '学校类别（综合/农林类等）',
  `us_rank` varchar(10) DEFAULT NULL COMMENT 'USNews排名',
  `vocational` enum('1','2') DEFAULT NULL COMMENT '是否职业院校（1是 2否）',
  `xiaoyuan` varchar(255) DEFAULT NULL COMMENT '校园全景URL',
  `xyh_rank` varchar(10) DEFAULT NULL COMMENT '校友会排名',
  `yjszs` varchar(255) DEFAULT NULL COMMENT '研究生招生信息链接',
  `zs_code` varchar(20) DEFAULT NULL COMMENT '院校代码',
  PRIMARY KEY (`school_id`),
  KEY `idx_province` (`province_name`),
  KEY `idx_city` (`city_name`),
  KEY `idx_type` (`type_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='高校基本信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `universities_colleges_departments`
--

DROP TABLE IF EXISTS `universities_colleges_departments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `universities_colleges_departments` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `school_id` int DEFAULT NULL COMMENT '关联的学校ID',
  `campus_name` varchar(100) DEFAULT NULL COMMENT '校区名称（如"校本部"）',
  `department_id` int DEFAULT NULL COMMENT '院系ID',
  `department_name` varchar(100) DEFAULT NULL COMMENT '院系名称（如"建筑学院"）',
  PRIMARY KEY (`id`),
  KEY `idx_school_id` (`school_id`),
  KEY `idx_campus` (`campus_name`),
  KEY `idx_department` (`department_name`)
) ENGINE=InnoDB AUTO_INCREMENT=417 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='学院与系表，记录学校的组织架构';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `universities_detail`
--

DROP TABLE IF EXISTS `universities_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `universities_detail` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '学校ID，唯一标识符',
  `name` varchar(100) NOT NULL COMMENT '学校全称（如"清华大学"）',
  `motto` varchar(100) DEFAULT NULL COMMENT '校训（如"自强不息，厚德载物"）',
  `address` varchar(200) DEFAULT NULL COMMENT '详细地址（如"北京市海淀区清华大学"）',
  `province_id` int DEFAULT NULL COMMENT '省份ID（如"11"代表北京）',
  `city_id` int DEFAULT NULL COMMENT '城市ID（如"1101"代表北京市）',
  `city_name` varchar(50) DEFAULT NULL COMMENT '城市名称（如"北京市"）',
  `county_id` int DEFAULT NULL COMMENT '区县ID（如"110108"代表海淀区）',
  `town_name` varchar(50) DEFAULT NULL COMMENT '乡镇/街道名称',
  `postcode` varchar(20) DEFAULT NULL COMMENT '邮政编码',
  `phone` varchar(100) DEFAULT NULL COMMENT '联系电话（可能包含多个，用逗号分隔）',
  `email` varchar(100) DEFAULT NULL COMMENT '联系邮箱',
  `site` varchar(200) DEFAULT NULL COMMENT '招生网站URL',
  `school_site` varchar(200) DEFAULT NULL COMMENT '学校官网URL',
  `belong` varchar(100) DEFAULT NULL COMMENT '隶属单位（如"教育部"）',
  `nature_name` varchar(50) DEFAULT NULL COMMENT '办学性质（如"公办"）',
  `type_name` varchar(50) DEFAULT NULL COMMENT '学校类型（如"综合类"）',
  `level_name` varchar(50) DEFAULT NULL COMMENT '办学层次（如"本科"）',
  `create_date` varchar(20) DEFAULT NULL COMMENT '创建年份（如"1911"）',
  `area` decimal(10,2) DEFAULT NULL COMMENT '占地面积（单位：亩）',
  `num_academician` int DEFAULT NULL COMMENT '院士人数',
  `num_doctor` int DEFAULT NULL COMMENT '博士点数量',
  `num_master` int DEFAULT NULL COMMENT '硕士点数量',
  `num_lab` int DEFAULT NULL COMMENT '实验室数量',
  `num_library` varchar(50) DEFAULT NULL COMMENT '图书馆藏书量（可能包含单位，如"548万"）',
  `num_subject` int DEFAULT NULL COMMENT '学科数量',
  `content` text COMMENT '学校简介/历史沿革',
  `f211` tinyint DEFAULT NULL COMMENT '是否211高校（1是，2否）',
  `f985` tinyint DEFAULT NULL COMMENT '是否985高校（1是，2否）',
  `logo_url` varchar(88) DEFAULT NULL COMMENT 'logo的静态url',
  `dual_class_name` varchar(50) DEFAULT NULL COMMENT '双一流类别（如"双一流"）',
  `qs_world` varchar(20) DEFAULT NULL COMMENT 'QS世界排名',
  `us_rank` varchar(20) DEFAULT NULL COMMENT 'US News排名',
  `ruanke_rank` varchar(20) DEFAULT NULL COMMENT '软科排名',
  `xyh_rank` varchar(20) DEFAULT NULL COMMENT '校友会排名',
  `status` varchar(2) DEFAULT NULL COMMENT '状态标识（1表示正常）',
  `add_time` datetime DEFAULT NULL COMMENT '数据添加时间',
  `last_updated` datetime DEFAULT NULL COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_universities_name` (`name`),
  KEY `idx_universities_province` (`province_id`),
  KEY `idx_universities_city` (`city_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3740 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='大学基本信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `universities_discipline_rankings`
--

DROP TABLE IF EXISTS `universities_discipline_rankings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `universities_discipline_rankings` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `school_id` int DEFAULT NULL COMMENT '关联的学校ID',
  `rank_level` varchar(10) DEFAULT NULL COMMENT '评估等级（如"A+","A","A-"等）',
  `count` int DEFAULT NULL COMMENT '该等级学科数量',
  PRIMARY KEY (`id`),
  KEY `idx_school_id` (`school_id`),
  KEY `idx_rank_level` (`rank_level`)
) ENGINE=InnoDB AUTO_INCREMENT=1889 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='学科评估表，记录学校各等级学科数量';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `universities_dual_class_subjects`
--

DROP TABLE IF EXISTS `universities_dual_class_subjects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `universities_dual_class_subjects` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `school_id` int DEFAULT NULL COMMENT '关联的学校ID',
  `subject_name` varchar(100) DEFAULT NULL COMMENT '学科名称（如"计算机科学与技术"）',
  `subject_id` int DEFAULT NULL COMMENT '学科ID',
  PRIMARY KEY (`id`),
  KEY `idx_school_id` (`school_id`),
  KEY `idx_subject_name` (`subject_name`),
  KEY `idx_dualclass_school` (`school_id`),
  KEY `idx_dualclass_subject` (`subject_name`)
) ENGINE=InnoDB AUTO_INCREMENT=509 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='双一流学科表，记录学校入选双一流的学科';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `universities_special_programs`
--

DROP TABLE IF EXISTS `universities_special_programs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `universities_special_programs` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `school_id` int DEFAULT NULL COMMENT '关联的学校ID',
  `special_id` int DEFAULT NULL COMMENT '专业ID',
  `special_name` varchar(100) DEFAULT NULL COMMENT '专业名称（如"建筑学"）',
  `level_name` varchar(50) DEFAULT NULL COMMENT '专业层次（如"本科(普通)"）',
  `limit_year` varchar(20) DEFAULT NULL COMMENT '学制（如"四年"或"五年"）',
  `nation_feature` varchar(2) DEFAULT NULL COMMENT '是否国家级特色专业（1是，2否）',
  `nation_first_class` varchar(2) DEFAULT NULL COMMENT '是否国家一流专业（1是，2否）',
  `province_feature` varchar(2) DEFAULT NULL COMMENT '是否省级特色专业（1是，2否）',
  `ruanke_level` varchar(10) DEFAULT NULL COMMENT '软科评级（如"A+"）',
  `ruanke_rank` int DEFAULT NULL COMMENT '软科排名',
  `xueke_rank_score` varchar(10) DEFAULT NULL COMMENT '学科评估结果（如"A+"）',
  `year` varchar(10) DEFAULT NULL COMMENT '数据年份（如"2022"）',
  `is_important` varchar(2) DEFAULT NULL COMMENT '是否重点专业（1是，2否）',
  PRIMARY KEY (`id`),
  KEY `idx_school_id` (`school_id`),
  KEY `idx_special_name` (`special_name`),
  KEY `idx_ruanke_rank` (`ruanke_rank`),
  KEY `idx_special_school` (`school_id`)
) ENGINE=InnoDB AUTO_INCREMENT=28014 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='专业表，记录学校开设的各类专业及其属性';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `universities_video`
--

DROP TABLE IF EXISTS `universities_video`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `universities_video` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `school_id` int NOT NULL COMMENT '学校ID',
  `video_type` tinyint DEFAULT '1' COMMENT '视频类型(1-PC,2-移动...)',
  `title` varchar(100) DEFAULT NULL COMMENT '视频标题',
  `url` varchar(200) NOT NULL COMMENT '视频URL',
  `img_url` varchar(200) DEFAULT NULL COMMENT '封面图片URL',
  `url_type` tinyint DEFAULT '1' COMMENT 'URL类型',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_school_id` (`school_id`),
  KEY `idx_video_type` (`video_type`)
) ENGINE=InnoDB AUTO_INCREMENT=3199 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='院校视频信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` int DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `pwd` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `token` varchar(255) DEFAULT NULL COMMENT '用户认证令牌',
  PRIMARY KEY (`id`),
  KEY `name` (`name`),
  KEY `email` (`email`),
  KEY `idx_user_id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
/*!50032 DROP TRIGGER IF EXISTS insert_user_trigger */;
DELIMITER ;;
/*!50003 CREATE*/ /*!50017 DEFINER=`root`@`localhost`*/ /*!50003 TRIGGER `insert_user_trigger` AFTER INSERT ON `user` FOR EACH ROW BEGIN
INSERT INTO profile (id, name, email)
VALUES (NEW.id, NEW.name, NEW.email);
END */;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
/*!50032 DROP TRIGGER IF EXISTS update_user_trigger */;
DELIMITER ;;
/*!50003 CREATE*/ /*!50017 DEFINER=`root`@`localhost`*/ /*!50003 TRIGGER `update_user_trigger` AFTER UPDATE ON `user` FOR EACH ROW BEGIN
UPDATE profile
SET name = NEW.name, email = NEW.email
WHERE id = NEW.id;
END */;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;

--
-- Table structure for table `user_id`
--

DROP TABLE IF EXISTS `user_id`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_id` (
  `id` int NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_visits`
--

DROP TABLE IF EXISTS `user_visits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_visits` (
  `email` varchar(255) NOT NULL,
  `school_id` int NOT NULL,
  `last_visit` timestamp NOT NULL,
  PRIMARY KEY (`email`,`school_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary view structure for view `v_special_admission_summary`
--

DROP TABLE IF EXISTS `v_special_admission_summary`;
/*!50001 DROP VIEW IF EXISTS `v_special_admission_summary`*/;
SET @saved_cs_client     = @@character_set_client;
/*!50503 SET character_set_client = utf8mb4 */;
/*!50001 CREATE VIEW `v_special_admission_summary` AS SELECT 
 1 AS `special_id`,
 1 AS `special_code`,
 1 AS `special_name`,
 1 AS `subject_category`,
 1 AS `subject_class`,
 1 AS `degree`,
 1 AS `school_count`,
 1 AS `avg_admission_score`,
 1 AS `min_admission_score`,
 1 AS `max_admission_score`,
 1 AS `total_admission_count`,
 1 AS `avg_employment_rate`,
 1 AS `avg_salary`*/;
SET character_set_client = @saved_cs_client;

--
-- Temporary view structure for view `v_special_employment_analysis`
--

DROP TABLE IF EXISTS `v_special_employment_analysis`;
/*!50001 DROP VIEW IF EXISTS `v_special_employment_analysis`*/;
SET @saved_cs_client     = @@character_set_client;
/*!50503 SET character_set_client = utf8mb4 */;
/*!50001 CREATE VIEW `v_special_employment_analysis` AS SELECT 
 1 AS `special_id`,
 1 AS `special_name`,
 1 AS `avg_employment_rate`,
 1 AS `avg_salary`,
 1 AS `famous_schools`,
 1 AS `top_industries`,
 1 AS `top_positions`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `verify_codes`
--

DROP TABLE IF EXISTS `verify_codes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `verify_codes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(88) DEFAULT NULL,
  `code` varchar(20) DEFAULT NULL,
  `expires_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping routines for database 'zxwl'
--
/*!50003 DROP PROCEDURE IF EXISTS `Campus_filter` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `Campus_filter`(IN p_goal_year int, IN p_type_id int, IN p_province_id int,
                                                     IN p_equivalent_score_start decimal(8, 2),
                                                     IN p_equivalent_score_end decimal(8, 2), IN p_curr_rank_start int,
                                                     IN p_curr_rank_end int, IN p_score_gap decimal(8, 2),
                                                     IN p_is_benke int, IN p_goal_province_id int, IN p_level2 text,
                                                     IN p_level3 text, IN p_salary decimal(10, 2))
BEGIN
    SELECT SQL_BUFFER_RESULT
        ud.id as university_id,
        ud.name as school_name,
        ud.address as school_address,
        ud.ruanke_rank as ruanke_rank,
        ud.xyh_rank as xyh_rank,
        cp.name as province_name,
        at.name as type_name,
        au.min_score as min_score,
        
        CASE WHEN ud.ruanke_rank <> '' THEN 1 ELSE 0 END as has_rk_rank,
        CASE WHEN ud.xyh_rank <> '' THEN 1 ELSE 0 END as has_xyh_rank
    FROM admission_universities as au
             LEFT JOIN universities_detail as ud ON ud.id = au.school_id
             LEFT JOIN admission_types as at ON at.id = au.score_type
             LEFT JOIN common_provinces as cp ON cp.id = ud.province_id
    WHERE
        au.year = p_goal_year
      AND au.province_id = p_province_id
      AND au.score_type = p_type_id
      AND au.min_score <= (p_equivalent_score_start + p_equivalent_score_end) / 2 + p_score_gap
    ORDER BY
        au.min_score DESC,      
        has_rk_rank DESC,       
        has_xyh_rank DESC,      
        ud.ruanke_rank ASC,     
        ud.xyh_rank ASC         
    LIMIT 250;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `get_equivalent_score` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `get_equivalent_score`(IN current_year int, IN province_id int, IN category_id int,
                                                            IN current_score decimal(6, 2), IN target_year int,
                                                            IN batch_id int, OUT curr_rank_start int,
                                                            OUT curr_rank_end int,
                                                            OUT equivalent_score_start decimal(6, 2),
                                                            OUT equivalent_score_end decimal(6, 2), OUT debug_info text)
BEGIN
    
    DECLARE current_exam_id, target_exam_id INT;            
    DECLARE rank_start_current, rank_end_current BIGINT;    
    DECLARE total_current, total_target BIGINT;             
    DECLARE ratio_start, ratio_end DECIMAL(10,8);           
    DECLARE target_rank_start, target_rank_end BIGINT;      
    DECLARE target_category_id INT;                         
    DECLARE is_special_city BOOLEAN;                        
    DECLARE batch_name VARCHAR(20);                         
    SET debug_info = '';

    
    SET is_special_city = (province_id IN (11, 12, 31)); 

    SET debug_info = CONCAT(debug_info, '步骤0: 是否为特殊城市: ', IF(is_special_city, '是', '否'), '\n');

    
    
    SELECT old_category_id INTO target_category_id
    FROM admission_category_mapping
    WHERE new_category_id = category_id
      AND target_year < start_year
      AND admission_category_mapping.province_id = province_id
    LIMIT 1;

    
    IF target_category_id IS NULL THEN
        SET target_category_id = category_id; 
    END IF;
    IF target_category_id IS NULL THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '转化考试类型失败';
    END IF;

    SET debug_info = CONCAT(debug_info, '步骤1: 考试类型转化: ', category_id, ' --> ', COALESCE(target_category_id, 'NULL'), '\n');

    
    IF is_special_city THEN
        
        IF batch_id IS NULL THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '特殊城市需要指定批次ID';
        END IF;

        SELECT id INTO current_exam_id
        FROM score_exam_info
        WHERE score_exam_info.province_id = province_id
          AND exam_year = current_year
          AND score_exam_info.type_id = category_id
          AND score_exam_info.batch_id = batch_id
        LIMIT 1;
    ELSE
        
        SELECT id INTO current_exam_id
        FROM score_exam_info
        WHERE score_exam_info.province_id = province_id
          AND exam_year = current_year
          AND score_exam_info.type_id = category_id
          AND score_exam_info.batch_id = 3
        LIMIT 1;
    END IF;

    SET debug_info = CONCAT(debug_info, '步骤2: 考试id: ', COALESCE(current_exam_id, 'NULL'),
                             ', 考试名称: ', COALESCE((SELECT exam_name FROM score_exam_info WHERE id = current_exam_id), '未找到'), '\n');

    
    SELECT rank_start, rank_end
    INTO rank_start_current, rank_end_current
    FROM score_section
    WHERE exam_id = current_exam_id AND min_score = current_score;
    set curr_rank_start = rank_start_current;
    set curr_rank_end = rank_end_current;

    SET debug_info = CONCAT(debug_info, '步骤3: 获取当前分数对应的排名区间: ',
                             COALESCE(rank_start_current, 'NULL'), ' - ', COALESCE(rank_end_current, 'NULL'), '\n');

    
    SELECT MAX(rank_end) INTO total_current
    FROM score_section
    WHERE exam_id = current_exam_id;

    SET debug_info = CONCAT(debug_info, '步骤4: 获取当前年份总人数: ', COALESCE(total_current, 'NULL'), '\n');
   
    
    SELECT score_section.batch_name INTO batch_name
    from score_section
    WHERE exam_id = current_exam_id AND min_score <= current_score and max_score >= current_score;
    SET debug_info = CONCAT(debug_info, '步骤5: 获取当前年份考试分数是否过线: ', COALESCE(batch_name, 'NULL'), '\n');
    
    
    IF is_special_city THEN
        
        SELECT id INTO target_exam_id
        FROM score_exam_info
        WHERE score_exam_info.province_id = province_id
          AND exam_year = target_year
          AND score_exam_info.type_id = target_category_id
          AND score_exam_info.batch_id = batch_id
        LIMIT 1;
    ELSE
        
        SELECT id INTO target_exam_id
        FROM score_exam_info
        WHERE score_exam_info.province_id = province_id
          AND exam_year = target_year
          AND score_exam_info.type_id = target_category_id
          AND score_exam_info.batch_id = 3
        LIMIT 1;
    END IF;

    SET debug_info = CONCAT(debug_info, '步骤6: 目标考试id: ', COALESCE(target_exam_id, 'NULL'),
                             ', 考试名称: ', COALESCE((SELECT exam_name FROM score_exam_info WHERE id = target_exam_id), '未找到'), '\n');

    
    SELECT MAX(rank_end) INTO total_target
    FROM score_section
    WHERE exam_id = target_exam_id;

    SET debug_info = CONCAT(debug_info, '步骤7: 获取目标年份总人数: ', COALESCE(total_target, 'NULL'), '\n');

    
    SET ratio_start = rank_start_current / total_current;
    SET ratio_end = rank_end_current / total_current;

    SET target_rank_start = CEIL(ratio_start * total_target);
    SET target_rank_end = CEIL(ratio_end * total_target);

    SET debug_info = CONCAT(debug_info, '步骤8: 当前排名在目标年份中的等效排名: ',
                             COALESCE(target_rank_start, 'NULL'), ' - ', COALESCE(target_rank_end, 'NULL'), '\n');

    
    SELECT min_score
    INTO equivalent_score_start
    FROM score_section
    WHERE exam_id = target_exam_id
      AND rank_start <= target_rank_start
      AND rank_end >= target_rank_start
    LIMIT 1;

    
    IF equivalent_score_start IS NULL THEN
        SELECT min_score
        INTO equivalent_score_start
        FROM score_section
        WHERE exam_id = target_exam_id
        ORDER BY ABS(((rank_start + rank_end) / 2) - target_rank_start)
        LIMIT 1;
    END IF;

    
    SELECT min_score
    INTO equivalent_score_end
    FROM score_section
    WHERE exam_id = target_exam_id
      AND rank_start <= target_rank_end
      AND rank_end >= target_rank_end
    LIMIT 1;

    
    IF equivalent_score_end IS NULL THEN
        SELECT min_score
        INTO equivalent_score_end
        FROM score_section
        WHERE exam_id = target_exam_id
        ORDER BY ABS(((rank_start + rank_end) / 2) - target_rank_end)
        LIMIT 1;
    END IF;

    SET debug_info = CONCAT(debug_info, '步骤9和10: 当前分数在目标年份中的等效分数: ',
                             COALESCE(equivalent_score_start, 'NULL'), ' - ', COALESCE(equivalent_score_end, 'NULL'), '\n');

    
    IF equivalent_score_start > equivalent_score_end THEN
        SET @temp_score = equivalent_score_start;
        SET equivalent_score_start = equivalent_score_end;
        SET equivalent_score_end = @temp_score;
    END IF;
    
    set debug_info = CONCAT(debug_info, '步骤11: 确保分数顺序正确: ',
                             COALESCE(equivalent_score_start, 'NULL'), ' - ', COALESCE(equivalent_score_end, 'NULL'), '\n');

END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `sp_filter` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `sp_filter`(IN p_goal_year int, IN p_type_id int, IN p_province_id int,
                                                 IN p_equivalent_score_start decimal(8, 2),
                                                 IN p_equivalent_score_end decimal(8, 2), IN p_curr_rank_start int,
                                                 IN p_curr_rank_end int, IN score_gap decimal(8, 2), IN p_is_benke int,
                                                 IN p_goal_province_id int, IN p_level2 text, IN p_level3 text,
                                                 IN p_salary decimal(10, 2))
BEGIN
    
    DECLARE v_level2_list TEXT;
    DECLARE v_level3_list TEXT;

    
    SET v_level2_list = p_level2;
    SET v_level3_list = p_level3;

    
    SET v_level2_list = REPLACE(REPLACE(REPLACE(v_level2_list, '(', ''), ')', ''), '''', '');
    SET v_level3_list = REPLACE(REPLACE(REPLACE(v_level3_list, '(', ''), ')', ''), '''', '');

    SELECT
        
        ud.id                   as school_id,                    
        ud.name                 as school_name,                  
        ud.address              as school_address,               
        
        sd.id                   as special_id,                   
        sd.name                 as special_name,                 
        sd.level1_name          as special_level1_name,          
        sd.level2_name          as special_level2_name,          
        sd.level3_name          as special_level3_name,          
        sd.avg_salary           as special_avg_salary,           
        ads.special_info        as special_info,                 
        ads.remark              as admission_special_remark,     
        sd.code                 as special_code,                 
        sd.subject_requirements as special_subject_requirements, 
        GROUP_CONCAT(DISTINCT sit.keyword) as special_keywords,  

        
        ads.id                  as admission_id,                 
        ads.min_rank            as score_min_rank,               
        ads.max_rank            as score_max_rank,               
        ads.min_score           as score_min_score,              
        ads.max_score           as score_max_score,              
        ads.average_score       as score_avg_score,              
        ads.admission_count     as admission_count,              

        
        cp.name                 as common_province_name,         
        adzt.name               as zslx_name,                    
        adb.name                as batch_name                    

    FROM admission_special ads
        LEFT JOIN common_provinces cp ON ads.province_id = cp.id
        LEFT JOIN admission_zhaosheng_type adzt ON ads.zslx = adzt.id
        LEFT JOIN admission_batches adb ON ads.batch = adb.id
        LEFT JOIN admission_types adt ON ads.type = adt.id
        LEFT JOIN special_detail sd ON ads.special_id = sd.id
        LEFT JOIN universities_detail ud ON ads.school_id = ud.id
             JOIN special_impression_tag sit ON sd.id = sit.special_id

    WHERE ads.year = p_goal_year
      AND ads.type = p_type_id
      AND ads.province_id = p_province_id
      AND ((ads.min_score < ((p_equivalent_score_end + p_equivalent_score_start) / 2 + score_gap))
        OR (ads.min_rank > p_curr_rank_start))
      AND (p_goal_province_id = 0 OR ud.province_id = p_goal_province_id)  
      AND (FIND_IN_SET(sd.level2_name, v_level2_list) > 0                  
        OR FIND_IN_SET(sd.level3_name, v_level3_list) > 0
        OR (v_level2_list = '' AND v_level3_list = ''))
      AND (p_salary = 0 OR sd.avg_salary >= p_salary)                      

    GROUP BY
        ud.id, ud.name, ud.address,
        sd.id, sd.name, sd.level1_name, sd.level2_name, sd.level3_name, sd.avg_salary,
        ads.special_info, ads.remark, sd.code, sd.subject_requirements,
        ads.id, ads.min_rank, ads.max_rank, ads.min_score, ads.max_score,
        ads.average_score, ads.admission_count,
        cp.name, adzt.name, adb.name
    ORDER BY ads.min_score DESC, ads.min_rank ASC, ud.ruanke_rank, ud.us_rank
    LIMIT 1000;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;

--
-- Final view structure for view `v_special_admission_summary`
--

/*!50001 DROP VIEW IF EXISTS `v_special_admission_summary`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `v_special_admission_summary` AS select `sd`.`id` AS `special_id`,`sd`.`code` AS `special_code`,`sd`.`name` AS `special_name`,`sd`.`type` AS `subject_category`,`sd`.`type_detail` AS `subject_class`,`sd`.`degree` AS `degree`,count(distinct `aspec`.`school_id`) AS `school_count`,avg(`aspec`.`average_score`) AS `avg_admission_score`,min(`aspec`.`min_score`) AS `min_admission_score`,max(`aspec`.`max_score`) AS `max_admission_score`,sum(`aspec`.`admission_count`) AS `total_admission_count`,avg(`sd`.`employment_rate`) AS `avg_employment_rate`,avg(`sd`.`avg_salary`) AS `avg_salary` from (`special_detail` `sd` left join `admission_special` `aspec` on((`sd`.`id` = `aspec`.`special_id`))) group by `sd`.`id`,`sd`.`code`,`sd`.`name`,`sd`.`type`,`sd`.`type_detail`,`sd`.`degree` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `v_special_employment_analysis`
--

/*!50001 DROP VIEW IF EXISTS `v_special_employment_analysis`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `v_special_employment_analysis` AS select `sd`.`id` AS `special_id`,`sd`.`name` AS `special_name`,avg(cast(`sd`.`employment_rate` as decimal(5,2))) AS `avg_employment_rate`,avg(`sd`.`avg_salary`) AS `avg_salary`,group_concat(distinct `sfs`.`school_name` order by `sfs`.`school_name` ASC separator ', ') AS `famous_schools`,(select group_concat(distinct `sjd`.`name` order by `sjd`.`rate` DESC separator ', ') from `special_job_distribution` `sjd` where ((`sjd`.`special_id` = `sd`.`id`) and (`sjd`.`distribution_type` = 1)) limit 3) AS `top_industries`,(select group_concat(distinct `sjd`.`position` order by `sjd`.`rate` DESC separator ', ') from `special_job_distribution` `sjd` where ((`sjd`.`special_id` = `sd`.`id`) and (`sjd`.`distribution_type` = 3)) limit 3) AS `top_positions` from (`special_detail` `sd` left join `special_famous_school` `sfs` on((`sd`.`id` = `sfs`.`special_id`))) group by `sd`.`id`,`sd`.`name` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-09-10 18:07:41
