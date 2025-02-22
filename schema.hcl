schema "main" {}

// TODO
// table "events" {
// 	schema = schema.main

// 	column "id" {
// 		null           = false
// 		type           = integer
// 		auto_increment = true
// 	}

// 	primary_key {
// 		columns = [column.id]
// 	}
// }

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

	column "updated_at" {
		null = false
		type = integer
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

table "stations" {
	schema = schema.main

	column "id" {
		null           = false
		type           = integer
		auto_increment = true
	}

	column "created_at" {
		null = false
		type = integer
	}

	column "updated_at" {
		null = false
		type = integer
	}

	column "name" {
		null = false
		type = text
	}

	column "short" {
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

	primary_key {
		columns = [column.id]
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
		null = false
		type = integer
	}

	column "updated_at" {
		null = false
		type = integer
	}

	column "name" {
		null = false
		type = text
	}

	column "short" {
		null = true
		type = text
	}

	column "size" {
		null = true
		type = integer
		default = 0
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

	primary_key {
		columns = [column.id]
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
		null = false
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
		null = false
		type = integer
	}

	column "updated_at" {
		null = false
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
		null = false
		type = integer
	}

	column "updated_at" {
		null = false
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

	index "idx_ptg" {
		columns = [column.station_id, column.group_id]
		unique  = true
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
		null = false
		type = integer
	}

	column "created_at" {
		null = false
		type = integer
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

table user_icons {
	schema = schema.main

	column "id" {
		null = false
		type = integer
	}

	column "created_at" {
		null = false
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
		ref_columns = [table.users.column.id]
		on_update   = NO_ACTION
		on_delete   = CASCADE
	}
}
