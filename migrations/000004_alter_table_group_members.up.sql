ALTER TABLE group_members DROP CONSTRAINT group_members_role_check;

ALTER TABLE group_members ADD CONSTRAINT group_members_role_check 
CHECK (role IN ('owner', 'participant'));