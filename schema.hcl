schema "main" {}

table "schedule" {
	schema = schema.main

	column "id" {
		null           = false
		type           = integer
		auto_increment = true
	}

	column "start" {
		null = false
		type = integer
	}

	column "end" {
		null = true
		type = integer
	}

	column "name" {
		null = false
		type = text
	}

	primary_key {
		columns = [column.id]
	}

	index "idx_schedule_start" {
		columns = [column.start]
		unique  = false
	}
}

table "tribes" {
	schema = schema.main

	column "id" {
		null           = false
		type           = integer
		auto_increment = true
	}

	column "created_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "updated_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "name" {
		null = false
		type = text
	}

	column "short" {
		null = true
		type = text
	}

	column "dpsg" {
		null = true
		type = text
	}

	column "image_id" {
		null = true
		type = text
	}

	column "email_domain" {
		null = true
		type = text
	}

	column "stavo_email" {
		null = true
		type = text
	}

	primary_key {
		columns = [column.id]
	}

	foreign_key "image_id" {
		columns     = [column.image_id]
		ref_columns = [table.images.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}

	index "idx_tribes_name" {
		columns = [column.name]
		unique  = true
	}
}

table "station_categories" {
	schema = schema.main

	column "id" {
		null           = false
		type           = integer
		auto_increment = true
	}

	column "name" {
		null = false
		type = text
	}

	column "max" {
		null    = false
		type    = integer
		default = 0
	}

	primary_key {
		columns = [column.id]
	}

	index "idx_station_categories_name" {
		columns = [column.name]
		unique  = true
	}
}

table "stations" {
	schema = schema.main

	column "id" {
		null           = false
		type           = integer
		auto_increment = true
	}

	column "created_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "created_by" {
		null = true
		type = integer
	}

	column "updated_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "updated_by" {
		null = true
		type = integer
	}

	column "name" {
		null = false
		type = text
	}

	column "abbr" {
		null = true
		type = text
	}

	column "size" {
		null = false
		type = integer
		default = 0
	}

	column "tribe_id" {
		null = false
		type = integer
	}

	column "category_id" {
		null = true
		type = integer
	}

	column "lati" {
		null = true
		type = real
	}

	column "long" {
		null = true
		type = real
	}

	column "image_id" {
		null = true
		type = text
	}

	column "description" {
		null = true
		type = text
	}

	column "requirements" {
		null = true
		type = text
	}

	primary_key {
		columns = [column.id]
	}

	foreign_key "tribe_id" {
		columns     = [column.tribe_id]
		ref_columns = [table.tribes.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}

	foreign_key "image_id" {
		columns     = [column.image_id]
		ref_columns = [table.images.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}

	foreign_key "created_by" {
		columns     = [column.created_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}

	foreign_key "updated_by" {
		columns     = [column.updated_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}

	foreign_key "category_id" {
		columns     = [column.category_id]
		ref_columns = [table.station_categories.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}

	index "idx_stations_abbr" {
		columns = [column.abbr]
		unique  = true
		where   = "abbr is not null"
	}
}

table "groups" {
	schema = schema.main

	column "id" {
		null           = false
		type           = integer
		auto_increment = true
	}

	column "created_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "created_by" {
		null = true
		type = integer
	}

	column "updated_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "updated_by" {
		null = true
		type = integer
	}

	column "name" {
		null = false
		type = text
	}

	column "abbr" {
		null = true
		type = text
	}

	column "size" {
		null = false
		type = integer
		default = 0
	}

	column "comment" {
		null = true
		type = text
	}

	column "grouping" {
		null = false
		type = integer
	}

	column "tribe_id" {
		null = false
		type = integer
	}

	column "image_id" {
		null = true
		type = text
	}

	primary_key {
		columns = [column.id]
	}

	foreign_key "tribe_id" {
		columns     = [column.tribe_id]
		ref_columns = [table.tribes.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}

	foreign_key "image_id" {
		columns     = [column.image_id]
		ref_columns = [table.images.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}

	foreign_key "created_by" {
		columns     = [column.created_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}

	foreign_key "updated_by" {
		columns     = [column.updated_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}

	index "idx_groups_abbr" {
		columns = [column.abbr]
		unique  = true
		where   = "abbr is not null"
	}
}

table "images" {
	schema = schema.main

	column "id" {
		null           = false
		type           = integer
		auto_increment = true
	}

	column "created_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "created_by" {
		null = true
		type = integer
	}

	column "filepath" {
		null = false
		type = text
	}

	column "tribe_id" {
		null = true
		type = integer
	}

	column "station_id" {
		null = true
		type = integer
	}

	primary_key {
		columns = [column.id]
	}

	foreign_key "created_by" {
		columns     = [column.created_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}

	foreign_key "tribe_id" {
		columns     = [column.tribe_id]
		ref_columns = [table.tribes.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}

	foreign_key "station_id" {
		columns     = [column.station_id]
		ref_columns = [table.stations.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}

	index "idx_image_filepath" {
		columns = [column.filepath]
		unique  = true
	}
}

table "points_to_stations" {
	schema = schema.main

	column "id" {
		null           = false
		type           = integer
		auto_increment = true
	}

	column "created_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "created_by" {
		null = true
		type = integer
	}

	column "updated_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "updated_by" {
		null = true
		type = integer
	}

	column "group_id" {
		null = false
		type = integer
	}

	column "station_id" {
		null = false
		type = integer
	}

	primary_key {
		columns = [column.id]
	}

	foreign_key "group_id" {
		columns     = [column.group_id]
		ref_columns = [table.groups.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}

	foreign_key "station_id" {
		columns     = [column.station_id]
		ref_columns = [table.stations.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}

	foreign_key "created_by" {
		columns     = [column.created_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}

	foreign_key "updated_by" {
		columns     = [column.updated_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}

	index "idx_pts" {
		columns = [column.group_id, column.station_id]
		unique  = true
	}
}

table "points_to_groups" {
	schema = schema.main

	column "id" {
		null           = false
		type           = integer
		auto_increment = true
	}

	column "created_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "created_by" {
		null = true
		type = integer
	}

	column "updated_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "updated_by" {
		null = true
		type = integer
	}

	column "station_id" {
		null = false
		type = integer
	}

	column "group_id" {
		null = false
		type = integer
	}

	primary_key {
		columns = [column.id]
	}

	foreign_key "station_id" {
		columns     = [column.station_id]
		ref_columns = [table.stations.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}

	foreign_key "group_id" {
		columns     = [column.group_id]
		ref_columns = [table.groups.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}

	foreign_key "created_by" {
		columns     = [column.created_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}

	foreign_key "updated_by" {
		columns     = [column.updated_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}

	index "idx_ptg" {
		columns = [column.station_id, column.group_id]
		unique  = true
	}
}

table "group_roles" {
	schema = schema.main

	column "id" {
		null           = false
		type           = integer
		auto_increment = true
	}

	column "created_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "created_by" {
		null = true
		type = integer
	}

	column "updated_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "updated_by" {
		null = true
		type = integer
	}

	column "user_id" {
		null = false
		type = integer
	}

	column "group_id" {
		null = false
		type = integer
	}

	column "group_role" {
		null = false
		type = integer
	}

	primary_key {
		columns = [column.id]
	}

	index "idx_group_roles_user" {
		columns = [column.user_id]
		unique  = false
	}

	index "idx_group_roles_user_group" {
		columns = [column.user_id, column.group_id]
		unique  = true
	}

	foreign_key "user_id" {
		columns     = [column.user_id]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}

	foreign_key "group_id" {
		columns     = [column.group_id]
		ref_columns = [table.groups.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}

	foreign_key "created_by" {
		columns     = [column.created_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}

	foreign_key "updated_by" {
		columns     = [column.updated_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}
}

table "station_roles" {
	schema = schema.main

	column "id" {
		null           = false
		type           = integer
		auto_increment = true
	}

	column "created_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "created_by" {
		null = true
		type = integer
	}

	column "updated_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "updated_by" {
		null = true
		type = integer
	}

	column "user_id" {
		null = false
		type = integer
	}

	column "station_id" {
		null = false
		type = integer
	}

	column "station_role" {
		null = false
		type = integer
	}

	primary_key {
		columns = [column.id]
	}

	index "idx_station_roles_user" {
		columns = [column.user_id]
		unique  = false
	}

	index "idx_station_roles_user_station" {
		columns = [column.user_id, column.station_id]
		unique  = true
	}

	foreign_key "user_id" {
		columns     = [column.user_id]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}

	foreign_key "station_id" {
		columns     = [column.station_id]
		ref_columns = [table.stations.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}

	foreign_key "created_by" {
		columns     = [column.created_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}

	foreign_key "updated_by" {
		columns     = [column.updated_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}
}

table "tribe_roles" {
	schema = schema.main

	column "id" {
		null           = false
		type           = integer
		auto_increment = true
	}

	column "created_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "created_by" {
		null = true
		type = integer
	}

	column "updated_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "updated_by" {
		null = true
		type = integer
	}

	column "user_id" {
		null = false
		type = integer
	}

	column "tribe_id" {
		null = false
		type = integer
	}

	column "tribe_role" {
		null = false
		type = integer
	}

	column "accepted_at" {
		null = true
		type = integer
	}

	primary_key {
		columns = [column.id]
	}

	index "idx_tribe_roles_user" {
		columns = [column.user_id]
		unique  = false
	}

	index "idx_tribe_roles_user_tribe" {
		columns = [column.user_id, column.tribe_id]
		unique  = true
	}

	foreign_key "user_id" {
		columns     = [column.user_id]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}

	foreign_key "tribe_id" {
		columns     = [column.tribe_id]
		ref_columns = [table.tribes.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}

	foreign_key "created_by" {
		columns     = [column.created_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}

	foreign_key "updated_by" {
		columns     = [column.updated_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}
}

table "users" {
	schema = schema.main

	column "id" {
		null           = false
		type           = integer
		auto_increment = true
	}

	column "ext_id" {
		null = true
		type = text
	}

	column "username" {
		null = false
		type = text
	}

	column "email" {
		null = false
		type = text
	}

	column "last_login" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "created_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "role" {
		null    = false
		type    = integer
		default = 0
	}

	column "firstname" {
		null = false
		type = text
	}

	column "lastname" {
		null = false
		type = text
	}

	primary_key {
		columns = [column.id]
	}

	index "idx_users_email" {
		columns = [column.email]
		unique  = true
	}

	index "idx_users_ext_id" {
		columns = [column.ext_id]
		unique  = true
		where   = "ext_id is not null"
	}
}

table tribe_icons {
	schema = schema.main

	column "id" {
		null           = false
		type           = integer
		auto_increment = true
	}

	column "created_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "created_by" {
		null = true
		type = integer
	}

	column "image" {
		null = false
		type = blob
	}

	primary_key {
		columns = [column.id]
	}

	foreign_key "id" {
		columns     = [column.id]
		ref_columns = [table.tribes.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}

	foreign_key "created_by" {
		columns     = [column.created_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}
}

table user_icons {
	schema = schema.main

	column "id" {
		null           = false
		type           = integer
		auto_increment = true
	}

	column "created_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "image" {
		null = false
		type = blob
	} 

	primary_key {
		columns = [column.id]
	}

	foreign_key "id" {
		columns     = [column.id]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}
}

table settings {
	schema = schema.main

	column "id" {
		null           = false
		type           = integer
		auto_increment = true
	}

	column "updated_at" {
		null    = false
		type    = integer
		default = "(unixepoch())"
	}

	column "updated_by" {
		null = true
		type = integer
	}

	column "data" {
		null = false
		type = blob
	}

	primary_key {
		columns = [column.id]
	}

	foreign_key "updated_by" {
		columns     = [column.updated_by]
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = SET_NULL
	}
}
