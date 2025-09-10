# 更新评论数量
DELIMITER $$

-- 创建在插入评论后更新评论总数的触发器
CREATE TRIGGER after_comment_insert
    AFTER INSERT ON news_comments
    FOR EACH ROW
BEGIN
    -- 更新资讯的评论总数（+1）
    UPDATE news_info
    SET comment_count = comment_count + 1
    WHERE id = NEW.news_id;
END$$

-- 创建在删除评论后更新评论总数的触发器
CREATE TRIGGER after_comment_delete
    AFTER DELETE ON news_comments
    FOR EACH ROW
BEGIN
    -- 更新资讯的评论总数（-1）
    UPDATE news_info
    SET comment_count = GREATEST(comment_count - 1, 0)
    WHERE id = OLD.news_id;
END$$

DELIMITER ;