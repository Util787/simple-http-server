-- name: GetChirpsByChirpId :one
SELECT * FROM chirps
WHERE id = $1;