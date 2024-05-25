package pgx

// query for Projects Storage
const (
	queryMigrateP = `CREATE TABLE IF NOT EXISTS projects
(
    "id" UUID NOT NULL UNIQUE,
    "name" VARCHAR NOT NULL,
    "created_by" UUID NOT NULL ,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("created_by") REFERENCES users(id) ON DELETE CASCADE
);`

	queryGetMyProjects = `SELECT p.id, p.name, p.created_by 
FROM projects AS p
WHERE created_by = $1;`

	queryGetProjectsByID = `SELECT p.id, p.name, p.created_by
FROM projects AS p
WHERE p.id = $1;`

	queryGetMyProjectsByName = `SELECT p.id, p.name, p.created_by 
FROM projects AS p
WHERE name = $1 and created_by = $2;`

	queryCreateProjects = `INSERT INTO projects (id, name, created_by) VALUES ($1, $2, $3);`

	queryUpdateProjectName = `UPDATE projects SET name = $1 WHERE id = $2;`

	queryDeleteProject = `DELETE FROM projects WHERE id = $1;`
)

// query for Todos Storage
const (
	queryMigrateT = `CREATE TABLE IF NOT EXISTS todos (
    "id" UUID PRIMARY KEY NOT NULL UNIQUE,
    "name" VARCHAR NOT NULL UNIQUE,
    "description" VARCHAR NOT NULL,
    "is_completed" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_by" UUID NOT NULL,
    "project_id" UUID NOT NULL,
    "column" VARCHAR NOT NULL,
    FOREIGN KEY (project_id, "column") REFERENCES project_columns(project_id, name),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS todos_created_by_index ON todos(created_by);
`
	queryCreateTodo  = `INSERT INTO todos (id, name, description, created_by, project_id, "column" )VALUES ($1, $2, $3, $4, $5, $6)`
	queryTodoGetByID = `SELECT created_by, name, id, description FROM todos WHERE id = $1`
	queryGetAllTodos = `SELECT t.id, t.name, t.description, t.is_completed, t.project_id
FROM todos AS t ;`
	queryUpdateTodo = `UPDATE todos
		SET name = $1, description = $2, is_completed = $3
		WHERE id = $4 AND created_by = $5 and project_id = $6`
	queryDeleteTodo = `DELETE FROM todos WHERE id = $1 and created_by = $2 and "column" = $3`
)

// query for Users Storage
const (
	queryInsertInto = `INSERT INTO users (id, login, encrypted_password) VALUES ($1, $2, $3);`

	queryGetByID = `SELECT u.id, u.login, u.encrypted_password
FROM users AS u
WHERE u.id = $1;`

	queryUpdatePassword = `UPDATE users SET encrypted_password = $1 WHERE id = $2;`

	queryGetByLogin = `SELECT u.id, u.login, u.encrypted_password
FROM users AS u
WHERE u.login = $1;`

	queryMigrateU = `CREATE TABLE IF NOT EXISTS users
(
    id UUID PRIMARY KEY NOT NULL UNIQUE ,
    login VARCHAR NOT NULL UNIQUE ,
    encrypted_password VARCHAR NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS users_login_idx ON users (login);`

	queryGetAllUsers = `SELECT u.id, u.login
FROM users AS u;`
)

// query for Columns Storage
const (
	queryMigrateC = `CREATE TABLE IF NOT EXISTS project_columns
(
    "project_id" UUID NOT NULL,
    "name" VARCHAR,
    "order" INT,
    PRIMARY KEY (project_id, name),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);`
	queryInsertC         = `INSERT INTO project_columns (project_id, name, "order") VALUES ($1, $2, $3);`
	queryDeleteC         = `DELETE FROM project_columns WHERE name = $1 and project_id = $2;`
	queryGetColumnByName = `SELECT * FROM project_columns WHERE name = $1 and project_id = $2;`
	queryUpdateColumns   = `UPDATE project_columns SET name = $1, "order" = $2 
                       WHERE name = $3 and project_id = $4;`
	queryGetAllColumns = `SELECT * FROM project_columns WHERE project_id = $1;`
)
