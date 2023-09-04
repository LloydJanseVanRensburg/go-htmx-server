UPDATE todos 
SET completed = ?, title = ?, updatedAt = CURRENT_TIMESTAMP 
WHERE id = ?;