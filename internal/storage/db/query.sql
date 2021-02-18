-- query.sql

-- name: GetSpecies :many
Select * from species;

-- name: GetSpeciesByID :one
Select * from species where id = $1;

-- name: InsertSpecies :exec
Insert into species(id, name, common_name, created_at) values($1, $2, $3, now());

-- name: UpdateSpecies :exec
Update species set name = $1, common_name=$2 where id = $3;

-- name: GetVolunteers :many
Select * from volunteers;

-- name: GetVolunteerByID :one
Select * from volunteers where id = $1;

-- name: InsertVolunteer :exec
Insert into volunteers(fname, lname, role_id, created_at, updated_at) values($1, $2, $3, NOW(), NOW());

-- name: UpdateVolunteer :exec
Update volunteers set fname  = $1, lname=$2, role_id=$3 where id = $4;

-- GetStates :many
Select * from states;

-- name: GetStateByID :one
Select * from states where id = $1;

-- name: InsertState :exec
Insert into states(category_id, name, created_at, updated_at) values($1, $2, now(), now());

-- name: UpdateState :exec
Update states set category_id = $1, name= $2 where id = $1;

-- name: GetObservations :many
Select * from observations;

-- name: GetObservationByID :one
Select * from observations where id = $1;

-- name: InsertObservation :exec
Insert into observations(link, created_at, updated_at) values($1,  now(), now());

-- name: UpdateObservation :exec
Update observations set link = $1 where id = $2;


-- name: GetObservationState :many
Select * from observation_state;

-- name: GetObservationStateByID :one
Select * from observation_state where id = $1;

-- name: InsertObservationState :exec
Insert into observation_state(observation_id, volunteer_id, state_id, comment, created_at, updated_at) values($1, $2, $3, $4, now(), now());

-- name: UpdateObservationState :exec
Update observation_state set observation_id = $1, volunteer_id= $2, state_id = $3, comment=$4 where id = $1;


-- name: GetNotifications :many
Select * from notifications;

-- name: GetNotification :one
Select * from notifications where id = $1;

-- name: InsertNotification :exec
Insert into notifications(observation_id, comment, link, created_at, updated_at) values($1, $2, $3, now(), now());

-- name: UpdateNotification :exec
Update notifications set observation_id = $1,comment = $2, link = $3 where id = $4;


-- name: GetCategories :many
Select * from categories;

-- name: GetCategory :one
Select * from categories where id = $1;

-- name: InsertCategory :exec
Insert into categories(name, created_at, updated_at) values ($1, now(), now());

-- name: UpdateCategory :exec
Update categories set name=$1 where id = $2;
