-- FlyRank Week 3 A2 - Stage 4
-- SQL exercises performed against tasks.db

-- 1. List all tasks
SELECT * FROM tasks;

-- 2. List completed tasks
SELECT * FROM tasks WHERE done = 1;

-- 3. Count all tasks
SELECT COUNT(*) FROM tasks;

-- 4. Mark all tasks completed
UPDATE tasks SET done = 1;

-- 5. Delete all completed tasks
DELETE FROM tasks WHERE done = 1;
