local element = osm2pgsql.define_table({
    name = 'element',
    ids = { type = 'any', type_column = 'osm_type', id_column = 'osm_id' },
    columns = {
        { column = 'name' },
        { column = 'class', not_null = true },
        { column = 'subclass' },
        { column = 'tags', type = 'jsonb', not_null = true },
        { column = 'location', type = 'point', not_null = true, projection = 4326 },

        -- Generated columns (Postgres computes; osm2pgsql must not fill them)
        { column = 'lon',
          sql_type = 'DOUBLE PRECISION NOT NULL GENERATED ALWAYS AS (ST_X("location")) STORED',
          create_only = true },
        { column = 'lat',
          sql_type = 'DOUBLE PRECISION NOT NULL GENERATED ALWAYS AS (ST_Y("location")) STORED',
          create_only = true },
}})

function process_element(object, location)
    local a = {
        name = object.tags.name,
        location = location,
        tags = object.tags
    }

    if not object.tags.name then
            return
        end

    if object.tags.amenity then
        a.class = 'amenity'
        a.subclass = object.tags.amenity
    elseif object.tags.shop then
        a.class = 'shop'
        a.subclass = object.tags.shop
    else
        return
    end

    element:insert(a)
end

function osm2pgsql.process_node(object)
    process_element(object, object:as_point())
end

function osm2pgsql.process_way(object)
    if object.is_closed and object.tags.building then
        process_element(object, object:as_polygon():centroid())
    end
end
