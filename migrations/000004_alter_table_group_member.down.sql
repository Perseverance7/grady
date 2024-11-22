-- Удаляем новое ограничение
ALTER TABLE group_members DROP CONSTRAINT group_members_role_check;

-- Восстанавливаем старое ограничение
ALTER TABLE group_members ADD CONSTRAINT group_members_role_check 
CHECK (role IN ('teacher', 'student'));
