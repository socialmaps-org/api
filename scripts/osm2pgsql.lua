-- luacheck: globals osm2pgsql

local element = osm2pgsql.define_table({
	name = "element",
	ids = { type = "any", type_column = "osm_type", id_column = "osm_id" },
	columns = {
		{ column = "location", type = "point", not_null = true, projection = 4326 },
		{ column = "tags", type = "jsonb", not_null = true },

		-- Generated columns (Postgres computes; osm2pgsql MUST NOT fill them):
		{
			column = "lon",
			sql_type = 'DOUBLE PRECISION NOT NULL GENERATED ALWAYS AS (ST_X("location")) STORED',
			create_only = true,
		},
		{
			column = "lat",
			sql_type = 'DOUBLE PRECISION NOT NULL GENERATED ALWAYS AS (ST_Y("location")) STORED',
			create_only = true,
		},
		{
			column = "name",
			sql_type = "TEXT NOT NULL GENERATED ALWAYS AS (tags ->> 'name') STORED",
			create_only = true,
		},
	},
})

local function process_element(object, location)
	local a = {
		location = location,
		tags = object.tags,
	}

	if object.tags.name and (object.tags.amenity or object.tags.shop) then
		element:insert(a)
	end
end

function osm2pgsql.process_node(object)
	process_element(object, object:as_point())
end

function osm2pgsql.process_way(object)
	if object.is_closed then
		process_element(object, object:as_polygon():centroid())
	end
end

function osm2pgsql.process_relation(object)
	local multipolygon = object:as_multipolygon()
	if not multipolygon:is_null() then
		process_element(object, multipolygon:centroid())
	end
end
