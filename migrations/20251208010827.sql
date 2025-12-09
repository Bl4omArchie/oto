-- Create "executables" table
CREATE TABLE "public"."executables" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "tag" text NULL,
  "name" text NULL,
  "version" text NULL,
  "path" text NULL,
  "description" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_executables_tag" UNIQUE ("tag")
);
-- Create index "idx_executables_deleted_at" to table: "executables"
CREATE INDEX "idx_executables_deleted_at" ON "public"."executables" ("deleted_at");
-- Create "commands" table
CREATE TABLE "public"."commands" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "executable_id" bigint NOT NULL,
  "requires_root" boolean NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_commands_name" UNIQUE ("name"),
  CONSTRAINT "fk_commands_executable" FOREIGN KEY ("executable_id") REFERENCES "public"."executables" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_commands_deleted_at" to table: "commands"
CREATE INDEX "idx_commands_deleted_at" ON "public"."commands" ("deleted_at");
-- Create "parameters" table
CREATE TABLE "public"."parameters" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "flag" text NOT NULL,
  "description" text NULL,
  "executable_tag" bigint NOT NULL,
  "requires_root" boolean NOT NULL,
  "requires_value" boolean NOT NULL,
  "value_type" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_parameters_executable" FOREIGN KEY ("executable_tag") REFERENCES "public"."executables" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_parameters_deleted_at" to table: "parameters"
CREATE INDEX "idx_parameters_deleted_at" ON "public"."parameters" ("deleted_at");
-- Create index "uid_executable_parameter" to table: "parameters"
CREATE UNIQUE INDEX "uid_executable_parameter" ON "public"."parameters" ("flag", "executable_tag");
-- Create "command_parameters" table
CREATE TABLE "public"."command_parameters" (
  "command_id" bigint NOT NULL,
  "parameter_id" bigint NOT NULL,
  PRIMARY KEY ("command_id", "parameter_id"),
  CONSTRAINT "fk_command_parameters_command" FOREIGN KEY ("command_id") REFERENCES "public"."commands" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_command_parameters_parameter" FOREIGN KEY ("parameter_id") REFERENCES "public"."parameters" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "flag_conflicts" table
CREATE TABLE "public"."flag_conflicts" (
  "flag_id" bigint NOT NULL,
  "interfer_id" bigint NOT NULL,
  PRIMARY KEY ("flag_id", "interfer_id"),
  CONSTRAINT "fk_flag_conflicts_interfer" FOREIGN KEY ("interfer_id") REFERENCES "public"."parameters" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_flag_conflicts_parameter" FOREIGN KEY ("flag_id") REFERENCES "public"."parameters" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "flag_dependencies" table
CREATE TABLE "public"."flag_dependencies" (
  "flag_id" bigint NOT NULL,
  "requires_id" bigint NOT NULL,
  PRIMARY KEY ("flag_id", "requires_id"),
  CONSTRAINT "fk_flag_dependencies_parameter" FOREIGN KEY ("flag_id") REFERENCES "public"."parameters" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_flag_dependencies_require" FOREIGN KEY ("requires_id") REFERENCES "public"."parameters" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "flag_values" table
CREATE TABLE "public"."flag_values" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "parameter_id" bigint NOT NULL,
  "value" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_flag_values_parameter" FOREIGN KEY ("parameter_id") REFERENCES "public"."parameters" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_flag_values_deleted_at" to table: "flag_values"
CREATE INDEX "idx_flag_values_deleted_at" ON "public"."flag_values" ("deleted_at");
-- Create "jobs" table
CREATE TABLE "public"."jobs" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NOT NULL,
  "command_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_jobs_name" UNIQUE ("name"),
  CONSTRAINT "fk_jobs_command" FOREIGN KEY ("command_id") REFERENCES "public"."commands" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_jobs_deleted_at" to table: "jobs"
CREATE INDEX "idx_jobs_deleted_at" ON "public"."jobs" ("deleted_at");
-- Create "job_flagvalues" table
CREATE TABLE "public"."job_flagvalues" (
  "job_id" bigint NOT NULL,
  "flag_value_id" bigint NOT NULL,
  PRIMARY KEY ("job_id", "flag_value_id"),
  CONSTRAINT "fk_job_flagvalues_flag_value" FOREIGN KEY ("flag_value_id") REFERENCES "public"."flag_values" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_job_flagvalues_job" FOREIGN KEY ("job_id") REFERENCES "public"."jobs" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
