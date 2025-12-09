table "command_parameters" {
  schema = schema.public
  column "command_id" {
    null = false
    type = bigint
  }
  column "parameter_id" {
    null = false
    type = bigint
  }
  primary_key {
    columns = [column.command_id, column.parameter_id]
  }
  foreign_key "fk_command_parameters_command" {
    columns     = [column.command_id]
    ref_columns = [table.commands.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_command_parameters_parameter" {
    columns     = [column.parameter_id]
    ref_columns = [table.parameters.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}
table "commands" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  column "deleted_at" {
    null = true
    type = timestamptz
  }
  column "name" {
    null = false
    type = text
  }
  column "description" {
    null = true
    type = text
  }
  column "executable_id" {
    null = false
    type = bigint
  }
  column "requires_root" {
    null = false
    type = boolean
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_commands_executable" {
    columns     = [column.executable_id]
    ref_columns = [table.executables.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_commands_deleted_at" {
    columns = [column.deleted_at]
  }
  unique "uni_commands_name" {
    columns = [column.name]
  }
}
table "executables" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  column "deleted_at" {
    null = true
    type = timestamptz
  }
  column "tag" {
    null = true
    type = text
  }
  column "name" {
    null = true
    type = text
  }
  column "version" {
    null = true
    type = text
  }
  column "path" {
    null = true
    type = text
  }
  column "description" {
    null = true
    type = text
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_executables_deleted_at" {
    columns = [column.deleted_at]
  }
  unique "uni_executables_tag" {
    columns = [column.tag]
  }
}
table "flag_conflicts" {
  schema = schema.public
  column "flag_id" {
    null = false
    type = bigint
  }
  column "interfer_id" {
    null = false
    type = bigint
  }
  primary_key {
    columns = [column.flag_id, column.interfer_id]
  }
  foreign_key "fk_flag_conflicts_interfer" {
    columns     = [column.interfer_id]
    ref_columns = [table.parameters.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_flag_conflicts_parameter" {
    columns     = [column.flag_id]
    ref_columns = [table.parameters.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}
table "flag_dependencies" {
  schema = schema.public
  column "flag_id" {
    null = false
    type = bigint
  }
  column "requires_id" {
    null = false
    type = bigint
  }
  primary_key {
    columns = [column.flag_id, column.requires_id]
  }
  foreign_key "fk_flag_dependencies_parameter" {
    columns     = [column.flag_id]
    ref_columns = [table.parameters.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_flag_dependencies_require" {
    columns     = [column.requires_id]
    ref_columns = [table.parameters.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}
table "flag_values" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  column "deleted_at" {
    null = true
    type = timestamptz
  }
  column "parameter_id" {
    null = false
    type = bigint
  }
  column "value" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_flag_values_parameter" {
    columns     = [column.parameter_id]
    ref_columns = [table.parameters.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_flag_values_deleted_at" {
    columns = [column.deleted_at]
  }
  index "uid_flag_value" {
    unique  = true
    columns = [column.parameter_id, column.value]
  }
}
table "job_flagvalues" {
  schema = schema.public
  column "job_id" {
    null = false
    type = bigint
  }
  column "flag_value_id" {
    null = false
    type = bigint
  }
  primary_key {
    columns = [column.job_id, column.flag_value_id]
  }
  foreign_key "fk_job_flagvalues_flag_value" {
    columns     = [column.flag_value_id]
    ref_columns = [table.flag_values.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_job_flagvalues_job" {
    columns     = [column.job_id]
    ref_columns = [table.jobs.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}
table "jobs" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  column "deleted_at" {
    null = true
    type = timestamptz
  }
  column "name" {
    null = false
    type = text
  }
  column "command_id" {
    null = false
    type = bigint
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_jobs_command" {
    columns     = [column.command_id]
    ref_columns = [table.commands.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_jobs_deleted_at" {
    columns = [column.deleted_at]
  }
  unique "uni_jobs_name" {
    columns = [column.name]
  }
}
table "parameters" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  column "deleted_at" {
    null = true
    type = timestamptz
  }
  column "flag" {
    null = false
    type = text
  }
  column "description" {
    null = true
    type = text
  }
  column "executable_tag" {
    null = false
    type = bigint
  }
  column "requires_root" {
    null = false
    type = boolean
  }
  column "requires_value" {
    null = false
    type = boolean
  }
  column "value_type" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_parameters_executable" {
    columns     = [column.executable_tag]
    ref_columns = [table.executables.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_parameters_deleted_at" {
    columns = [column.deleted_at]
  }
  index "uid_executable_parameter" {
    unique  = true
    columns = [column.flag, column.executable_tag]
  }
}
schema "public" {
  comment = "standard public schema"
}
