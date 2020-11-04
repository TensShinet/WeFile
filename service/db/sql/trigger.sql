use wefile;
delimiter //
CREATE TRIGGER after_user_files_delete
AFTER DELETE
ON wefile.user_files FOR EACH ROW
BEGIN
    SET
    @file_id = OLD.file_id,
    @is_directory = OLD.is_directory;

    IF @is_directory = 0 THEN
        SELECT count
        INTO @file_count
        FROM files
        WHERE id = @file_id
        FOR UPDATE;

        SET @file_count = @file_count - 1;

        UPDATE files
        SET count = @file_count
        WHERE id = @file_id;
    end if;
end//


CREATE TRIGGER after_group_files_delete
AFTER DELETE
ON wefile.group_files FOR EACH ROW
BEGIN
    SET
    @file_id = OLD.file_id,
    @is_directory = OLD.is_directory;

    IF @is_directory = 0 THEN
        SELECT count
        INTO @file_count
        FROM files
        WHERE id = @file_id
            FOR UPDATE;

        SET @file_count = @file_count - 1;

        UPDATE files
        SET count = @file_count
        WHERE id = @file_id;
    end if;
end//

delimiter ;

