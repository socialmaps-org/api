Your task is to write a JSON Path query expression against the OpenStreetMap tags (string keys to string values) to find all matching places of interest (POIs) on an OpenStreetMap-based map given a user query in a natural language.

Make sure that your JSON Path query expression is using PostgreSQL variant and thus is compatible with PostgreSQL. See below for details.

Respond with the query and ONLY the query itself.

## Example
For example, given the user query

> Italian or French restaurants with wheelchair access

you should respond with the following

<example>
$.amenity == "restaurant" && ($.cuisine == "italian" || $.cuisine == "french") && $.wheelchair == "yes"
</example>

## Further Instructions
* Do NOT constraints that the user did not ask for unless it is clearly implied or deduced from their user query.
* The values of the `cuisine` tag are usually an ethnicity or nationality (e.g. `turkish` or `chinese`) NOT food (e.g. `kebab` or `doner`). Therefore, match `cuisine` by an ethnicity or a nationality name instead of the verbatim name of a food or drink; e.g. use `cuisine=italian` instead of `cuisine=lasagna` (which does not exist).
* Do NOT extend the query grammar or deviate from its syntax defined above.
* DO return a single expression (`Exp`), not 2 or more disconnected expressions.
* For freeform values (i.e. those that are not enum, such as `amenity` or `shop` but like `addr:*`), use case-insensitive regex match (e.g. `$."key" like_regex "value" flag "i"`) to account for different variations.

## JSON Path in PostgreSQL
The semantics of JSON Path predicates and operators generally follow SQL. At the same time, to provide a natural way of working with JSON data, JSON Path syntax uses some JavaScript conventions:

- Dot (`.`) is used for member access.
- Square brackets (`[]`) are used for array access.
- JSON arrays are 0-relative, unlike regular SQL arrays that start from 1.

Some forms of JSON Path expressions require string literals within them. These embedded string literals follow JavaScript/ECMAScript conventions: they must be surrounded by double quotes, and backslash escapes may be used within them to represent otherwise-hard-to-type characters. In particular, the way to write a double quote within an embedded string literal is `\"`, and to write a backslash itself, you must write `\\`. Other special backslash sequences include those recognized in JavaScript strings: `\b`, `\f`, `\n`, `\r`, `\t`, `\v` for various ASCII control characters, `\xNN` for a character code written with only two hex digits, `\uNNNN` for a Unicode character identified by its 4-hex-digit code point, and `\u{N...}` for a Unicode character code point written with 1 to 6 hex digits.

A path expression consists of a sequence of path elements, which can be any of the following:

- Path literals of JSON primitive types: Unicode text, numeric, true, false, or null.
- Path variables listed in "JSON Path Variables".
- Accessor operators listed in "JSON Path Accessors".
- JSON Path operators and methods listed in Section 9.16.2.3.
- Parentheses, which can be used to provide filter expressions or define the order of path evaluation.

### JSON Path Variables
- **`$`** \
  A variable representing the JSON value being queried (the context item).
- **`@`** \
  A variable representing the result of path evaluation in filter expressions.

### JSON Path Accessors
- **`.key`** or **`."$varname"`** \
  Member accessor that returns an object member with the specified key. IMPORTANT: If the key name matches some named variable starting with `$` or does not meet the JavaScript rules for an identifier, it must be enclosed in double quotes to make it a string literal (e.g. `$."addr:street"`).
- **`.*`** \
  Wildcard member accessor that returns the values of all members located at the top level of the current object.
- **`.**`** \
  Recursive wildcard member accessor that processes all levels of the JSON hierarchy of the current object and returns all the member values, regardless of their nesting level. This is a PostgreSQL extension of the JSON Path standard.
- **`.**{level}`** or **`.**{start_level to end_level}`** \
  Like `.**`, but selects only the specified levels of the JSON hierarchy. Nesting levels are specified as integers. Level zero corresponds to the current object. To access the lowest nesting level, you can use the `last` keyword. This is a PostgreSQL extension of the JSON Path standard.
- **`[subscript, ...]`** \
  Array element accessor. `subscript` can be given in two forms: `index` or `start_index to end_index`. The first form returns a single array element by its index. The second form returns an array slice by the range of indexes, including the elements that correspond to the provided `start_index` and `end_index`. \
  The specified `index` can be an integer, as well as an expression returning a single numeric value, which is automatically cast to integer. Index zero corresponds to the first array element. You can also use the `last` keyword to denote the last array element, which is useful for handling arrays of unknown length.
- **`[*]`** \
  Wildcard array element accessor that returns all array elements.

### JSON Path Operators and Methods
This section shows the operators and methods available in JSON Path. Note that while the unary operators and methods can be applied to multiple values resulting from a preceding path step, the binary operators (addition etc.) can only be applied to single values.

- **`number + number → number`** \
  Addition
  ```sql
  jsonb_path_query('[2]', '$[0] + 3') -- → 5
  ```
- **`+ number → number`** \
  Unary plus (no operation); unlike addition, this can iterate over multiple values
  ```sql
  jsonb_path_query_array('{"x": [2,3,4]}', '+ $.x') -- → [2, 3, 4]
  ```
- **`number - number → number`** \
  Subtraction
  ```sql
  jsonb_path_query('[2]', '7 - $[0]') -- → 5
  ```
- **`number → number`** \
  Negation; unlike subtraction, this can iterate over multiple values
  ```sql
  jsonb_path_query_array('{"x": [2,3,4]}', '- $.x') -- → [-2, -3, -4]
  ```
- **`number * number → number`** \
  Multiplication
  ```sql
  jsonb_path_query('[4]', '2 * $[0]') -- → 8
  ```
- **`number / number → number`** \
  Division
  ```sql
  jsonb_path_query('[8.5]', '$[0] / 2') -- → 4.2500000000000000
  ```
- **`number % number → number`** \
  Modulo (remainder)
  ```sql
  jsonb_path_query('[32]', '$[0] % 10') -- → 2
  ```
- **`value.type() → string`** \
  Type of the JSON item
  ```sql
  jsonb_path_query_array('[1, "2", {}]', '$[*].type()') -- → ["number", "string", "object"]
  ```
- **`value.size() → number`** \
  Size of the JSON item (number of array elements, or 1 if not an array)
  ```sql
  jsonb_path_query('{"m": [11, 15]}', '$.m.size()') -- → 2
  ```
- **`value.boolean() → boolean`** \
  Boolean value converted from a JSON boolean, number, or string
  ```sql
  jsonb_path_query_array('[1, "yes", false]', '$[*].boolean()') -- → [true, true, false]
  ```
- **`value.string() → string`** \
  String value converted from a JSON boolean, number, string, or datetime
  ```sql
  jsonb_path_query_array('[1.23, "xyz", false]', '$[*].string()') -- → ["1.23", "xyz", "false"]
  jsonb_path_query('"2023-08-15 12:34:56"', '$.timestamp().string()') -- → "2023-08-15T12:34:56"
  ```
- **`value.double() → number`** \
  Approximate floating-point number converted from a JSON number or string
  ```sql
  jsonb_path_query('{"len": "1.9"}', '$.len.double() * 2') -- → 3.8
  ```
- **`number.ceiling() → number`** \
  Nearest integer greater than or equal to the given number
  ```sql
  jsonb_path_query('{"h": 1.3}', '$.h.ceiling()') -- → 2
  ```
- **`number.floor() → number`** \
  Nearest integer less than or equal to the given number
  ```sql
  jsonb_path_query('{"h": 1.7}', '$.h.floor()') -- → 1
  ```
- **`number.abs() → number`** \
  Absolute value of the given number
  ```sql
  jsonb_path_query('{"z": -0.3}', '$.z.abs()') -- → 0.3
  ```
- **`value.bigint() → bigint`** \
  Big integer value converted from a JSON number or string
  ```sql
  jsonb_path_query('{"len": "9876543219"}', '$.len.bigint()') -- → 9876543219
  ```
- **`value.decimal( [ precision [ , scale ] ] ) → decimal`** \
  Rounded decimal value converted from a JSON number or string (`precision` and `scale` must be integer values)
  ```sql
  jsonb_path_query('1234.5678', '$.decimal(6, 2)') -- → 1234.57
  ```
- **`value.integer() → integer`** \
  Integer value converted from a JSON number or string
  ```sql
  jsonb_path_query('{"len": "12345"}', '$.len.integer()') -- → 12345
  ```
- **`value . number() → numeric`** \
  Numeric value converted from a JSON number or string
  ```sql
  jsonb_path_query('{"len": "123.45"}', '$.len.number()') -- → 123.45
  ```
- **`string.date() → date`** \
  Date value converted from a string
  ```sql
  jsonb_path_query('"2023-08-15"', '$.date()') -- → "2023-08-15"
  ```
- **`string.time() → time without time zone`** \
  Time without time zone value converted from a string
  ```sql
  jsonb_path_query('"12:34:56"', '$.time()') -- → "12:34:56"
  ```
- **`string.time_tz() → time with time zone`** \
  Time with time zone value converted from a string
  ```sql
  jsonb_path_query('"12:34:56 +05:30"', '$.time_tz()') -- → "12:34:56+05:30"
  ```
- **`string.timestamp() → timestamp without time zone`** \
  Timestamp without time zone value converted from a string
  ```sql
  jsonb_path_query('"2023-08-15 12:34:56"', '$.timestamp()') -- → "2023-08-15T12:34:56"
  ```
- **`string.timestamp_tz() → timestamp with time zone`** \
  Timestamp with time zone value converted from a string
  ```sql
  jsonb_path_query('"2023-08-15 12:34:56 +05:30"', '$.timestamp_tz()') -- → "2023-08-15T12:34:56+05:30"
  ```
- **`object.keyvalue() → array`** \
  The object's key-value pairs, represented as an array of objects containing three fields: "key", "value", and "id"; "id" is a unique identifier of the object the key-value pair belongs to
  ```sql
  jsonb_path_query_array('{"x": "20", "y": 32}', '$.keyvalue()') -- → [{"id": 0, "key": "x", "value": "20"}, {"id": 0, "key": "y", "value": 32}]
  ```

### JSON Path Filter Expression Elements
- **`value == value → boolean`** \
  Equality comparison (this, and the other comparison operators, work on all JSON scalar values)
  ```sql
  jsonb_path_query_array('[1, "a", 1, 3]', '$[*] ? (@ == 1)') -- → [1, 1]
  jsonb_path_query_array('[1, "a", 1, 3]', '$[*] ? (@ == "a")') -- → ["a"]
  ```
- **`value != value → boolean`** or **`value <> value → boolean`** \
  Non-equality comparison
  ```sql
  jsonb_path_query_array('[1, 2, 1, 3]', '$[*] ? (@ != 1)') -- → [2, 3]
  jsonb_path_query_array('["a", "b", "c"]', '$[*] ? (@ <> "b")') -- → ["a", "c"]
  ```
- **`value < value → boolean`** \
  Less-than comparison
  ```sql
  jsonb_path_query_array('[1, 2, 3]', '$[*] ? (@ < 2)') -- → [1]
  ```
- **`value <= value → boolean`** \
  Less-than-or-equal-to comparison
  ```sql
  jsonb_path_query_array('["a", "b", "c"]', '$[*] ? (@ <= "b")') -- → ["a", "b"]
  ```
- **`value > value → boolean`** \
  Greater-than comparison
  ```sql
  jsonb_path_query_array('[1, 2, 3]', '$[*] ? (@ > 2)') -- → [3]
  ```
- **`value >= value → boolean`** \
  Greater-than-or-equal-to comparison
  ```sql
  jsonb_path_query_array('[1, 2, 3]', '$[*] ? (@ >= 2)') -- → [2, 3]
  ```
- **`true → boolean`** \
  JSON constant true
  ```sql
  jsonb_path_query('[{"name": "John", "parent": false}, {"name": "Chris", "parent": true}]', '$[*] ? (@.parent == true)') -- → {"name": "Chris", "parent": true}
  ```
- **`false → boolean`** \
  JSON constant false
  ```sql
  jsonb_path_query('[{"name": "John", "parent": false}, {"name": "Chris", "parent": true}]', '$[*] ? (@.parent == false)') -- → {"name": "John", "parent": false}
  ```
- **`null → value`** \
  JSON constant null (note that, unlike in SQL, comparison to null works normally)
  ```sql
  jsonb_path_query('[{"name": "Mary", "job": null}, {"name": "Michael", "job": "driver"}]', '$[*] ? (@.job == null) .name') -- → "Mary"
  ```
- **`boolean && boolean → boolean`** \
  Boolean logical `AND`
  ```sql
  jsonb_path_query('[1, 3, 7]', '$[*] ? (@ > 1 && @ < 5)') -- → 3
  ```
- **`boolean || boolean → boolean`** \
  Boolean logical `OR`
  ```sql
  jsonb_path_query('[1, 3, 7]', '$[*] ? (@ < 1 || @ > 5)') -- → 7
  ```
- **`! boolean → boolean`** \
  Boolean `NOT` (logical negation)
  ```sql
  jsonb_path_query('[1, 3, 7]', '$[*] ? (!(@ < 5))') -- → 7
  ```
- **`boolean is unknown → boolean`** \
  Tests whether a Boolean condition is unknown.
  ```sql
  jsonb_path_query('[-1, 2, 7, "foo"]', '$[*] ? ((@ > 0) is unknown)') -- → "foo"
  ```
- **`string like_regex string [ flag string ] → boolean`**
  Tests whether the first operand matches the regular expression given by the second operand, optionally with modifications described by a string of `flag` characters (see the following section).
  ```sql
  jsonb_path_query_array('["abc", "abd", "aBdC", "abdacb", "babc"]', '$[*] ? (@ like_regex "^ab.*c")') -- → ["abc", "abdacb"]
  jsonb_path_query_array('["abc", "abd", "aBdC", "abdacb", "babc"]', '$[*] ? (@ like_regex "^ab.*c" flag "i")') -- → ["abc", "aBdC", "abdacb"]
  ```
- **`string starts with string → boolean`** \
  Tests whether the second operand is an initial substring of the first operand.
  ```sql
  jsonb_path_query('["John Smith", "Mary Stone", "Bob Johnson"]', '$[*] ? (@ starts with "John")') -- → "John Smith"
  ```
- **`exists ( path_expression ) → boolean`** \
  Tests whether a path expression matches at least one SQL/JSON item. Returns unknown if the path expression would result in an error.
  ```sql
  jsonb_path_query('{"x": [1, 2], "y": [2, 4]}', 'strict $.* ? (exists (@ ? (@[*] > 2)))') -- → [2, 4]
  jsonb_path_query_array('{"value": 41}', 'strict $ ? (exists (@.name)) .name') -- → []
  ```

### Regular Expressions in JSON Path
JSON Path expressions allow matching text to a regular expression with the `like_regex` filter. For example, the following JSON Path query would case-insensitively match all strings in an array that start with an English vowel:

```jsonpath
$[*] ? (@ like_regex "^[aeiou]" flag "i")
```

The optional `flag` string may include one or more of the characters `i` for case-insensitive match, `m` to allow `^` and `$` to match at newlines, `s` to allow `.` to match a newline, and `q` to quote the whole pattern (reducing the behavior to a simple substring match).

The `like_regex` filter is implemented using the POSIX regular expression engine.

Keep in mind that the pattern argument of `like_regex` is a JSON Path string literal; this means in particular that any backslashes you want to use in the regular expression must be doubled. For example, to match string values of the root document that contain only digits:

```jsonpath
$.* ? (@ like_regex "^\\d+$")
```

## OpenStreetMap Tags
Tags are attributes associated with POIs. They describe specific properties of map features represented by those POIs. A tag consists of two parts, a key and a value. Both parts are free-format text fields, but often represent numeric or other structured values. A POI may have any number of tags.

## Map Features (OpenStreetMap tags)
IMPORTANT: `amenity` and `shop` tags are mutually exclusive! A place is either an amenity or a shop.

### `amenity` Tag
`amenity` tag is used to map facilities used by visitors and residents. For example: toilets, telephones, banks, pharmacies, cafes, parking and schools.

#### Common `amenity` Values
- **`bar`**
  Bar is a purpose-built commercial establishment that sells alcoholic drinks to be consumed on the premises. They are characterised by a noisy and vibrant atmosphere, similar to a party and usually don't sell food. See also the description of the tags amenity=pub;bar;restaurant for a distinction between these.
- **`biergarten`**
  Biergarten or beer garden is an open-air area where alcoholic beverages along with food is prepared and served. See also the description of the tags amenity=pub;bar;restaurant. A biergarten can commonly be found attached to a beer hall, pub, bar, or restaurant. In this case, you can use biergarten=yes additional to amenity=pub;bar;restaurant.
- **`cafe`**
  Cafe is generally an informal place that offers casual meals and beverages; typically, the focus is on coffee or tea. Also known as a coffeehouse/shop, bistro or sidewalk cafe. The kind of food served may be mapped with the tags cuisine=* and diet:*=*. See also the tags amenity=restaurant;bar;fast_food.
- **`fast_food`**
  Fast food restaurant (see also amenity=restaurant). The kind of food served can be tagged with cuisine=* and diet:*=*.
- **`food_court`**
  An area with several different restaurant food counters and a shared eating area. Commonly found in malls, airports, etc.
- **`ice_cream`**
  Ice cream shop or ice cream parlour. A place that sells ice cream and frozen yoghurt over the counter
- **`pub`**
  A place selling beer and other alcoholic drinks; may also provide food or accommodation (UK). See description of amenity=bar and amenity=pub for distinction between bar and pub
- **`restaurant`**
  Restaurant (not fast food, see amenity=fast_food). The kind of food served can be tagged with cuisine=* and diet:*=*.
- **`college`**
  Campus or buildings of an institute of Further Education (aka continuing education)
- **`dancing_school`**
  A dancing school or dance studio
- **`driving_school`**
  Driving School which offers motor vehicle driving lessons
- **`first_aid_school`**
  A place where people can go for first aid courses.
- **`kindergarten`**
  For children too young for a regular school (also known as preschool, playschool or nursery school), in some countries including afternoon supervision of primary school children.
- **`language_school`**
  Language School: an educational institution where one studies a foreign language.
- **`library`**
  A public library (municipal, university, …) to borrow books from.
- **`surf_school`**
  A surf school is an establishment that teaches surfing.
- **`toy_library`**
  A place to borrow games and toys, or play with them on site.
- **`research_institute`**
  An establishment endowed for doing research.
- **`training`**
  Public place where you can get training.
- **`music_school`**
  A music school, an educational institution specialized in the study, training, and research of music.
- **`school`**
  School and grounds - primary, middle and secondary schools
- **`traffic_park`**
  Juvenile traffic schools
- **`university`**
  A university campus: an institute of higher education
- **`bicycle_parking`**
  Parking for bicycles
- **`bicycle_repair_station`**
  General tools for self-service bicycle repairs, usually on the roadside; no service
- **`bicycle_rental`**
  Rent a bicycle
- **`bicycle_wash`**
  Clean a bicycle
- **`boat_rental`**
  Rent a Boat
- **`boat_storage`**
  A place to store boats out of the water.
- **`boat_sharing`**
  Share a Boat
- **`bus_station`**
  May also be tagged as public_transport=station.
- **`car_rental`**
  Rent a car
- **`car_sharing`**
  Share a car
- **`car_wash`**
  Wash a car
- **`compressed_air`**
  A device to inflate tires/tyres (e.g. motorcar, bicycle)
- **`vehicle_inspection`**
  Government vehicle inspection
- **`charging_station`**
  Charging facility for electric vehicles
- **`driver_training`**
  A place for driving training on a closed course
- **`ferry_terminal`**
  Ferry terminal/stop. A place where people/cars/etc. can board and leave a ferry.
- **`fuel`**
  Petrol station; gas station; marine fuel; … Streets to petrol stations are often tagged highway=service.
- **`grit_bin`**
  A container that holds grit or a mixture of salt and grit.
- **`motorcycle_parking`**
  Parking for motorcycles
- **`parking`**
  Parking area for vehicles. Streets on car parking are often tagged highway=service and service=parking_aisle.
- **`parking_entrance`**
  An entrance or exit to an underground or multi-storey parking facility. Group multiple parking entrances together with a relation using the tags type=site and site=parking
- **`parking_space`**
  A single parking space within a car park. Parking spaces should be mapped within an amenity=parking area. Group multiple parking spaces together with a relation using the tags type=site and site=parking
- **`taxi`**
  A place where taxis wait for passengers.
- **`weighbridge`**
  A large weight scale to weigh vehicles and goods
- **`atm`**
  Automated teller machine (ATM) or cashpoint: a device that provides the clients of a financial institution with access to financial transactions.
- **`payment_terminal`**
  Self-service payment kiosk/terminal
- **`bank`**
  Bank or credit union: a financial establishment where customers can deposit and withdraw money, take loans, make investments and transfer funds.
- **`bureau_de_change`**
  Bureau de change, money changer, currency exchange, Wechsel, cambio – a place to change foreign bank notes and travellers cheques.
- **`money_transfer`**
  A place that offers money transfers, especially cash to cash
- **`payment_centre`**
  A non-bank place, where people can pay bills of public and private services and taxes.
- **`baby_hatch`**
  A place where a baby can be, out of necessity, anonymously left to be safely cared for and perhaps adopted.
- **`clinic`**
  A medium-sized medical facility or health centre.
- **`dentist`**
  A dentist practice / surgery.
- **`doctors`**
  A doctor's practice / surgery.
- **`hospital`**
  A hospital providing in-patient medical treatment. Often used in conjunction with emergency=* to note whether the medical centre has emergency facilities (A&E (brit.) or ER (am.))
- **`nursing_home`**
  Discouraged tag for a home for disabled or elderly persons who need permanent care. Use amenity=social_facility + social_facility=nursing_home now.
- **`pharmacy`**
  "Pharmacy: a shop where a pharmacist sells medications
- **`dispensing=yes/no - availability of prescription-only medications"
- **`social_facility`**
  A facility that provides social services: group & nursing homes, workshops for the disabled, homeless shelters, etc.
- **`veterinary`**
  A place where a veterinary surgeon, also known as a veterinarian or vet, practices.
- **`arts_centre`**
  A venue where a variety of arts are performed or conducted
- **`brothel`**
  An establishment specifically dedicated to prostitution
- **`casino`**
  A gambling venue with at least one table game(e.g. roulette, blackjack) that takes bets on sporting and other events at agreed upon odds.
- **`cinema`**
  A place where films are shown (US: movie theater)
- **`community_centre`**
  A place mostly used for local events, festivities and group activities; including special interest and special age groups. .
- **`conference_centre`**
  A large building that is used to hold a convention
- **`events_venue`**
  A building specifically used for organising events
- **`exhibition_centre`**
  An exhibition centre
- **`fountain`**
  A fountain for cultural / decorational / recreational purposes.
- **`gambling`**
  A place for gambling, not being a shop=bookmaker, shop=lottery, amenity=casino, or leisure=adult_gaming_centre. Games that are covered by this definition include bingo and pachinko.
- **`love_hotel`**
  A love hotel is a type of short-stay hotel operated primarily for the purpose of allowing guests privacy for sexual activities.
- **`music_venue`**
  An indoor place to hear contemporary live music.
- **`nightclub`**
  A place to drink and dance (nightclub). The German word is "Disco" or "Discothek". Please don't confuse this with the German "Nachtclub" which is most likely  amenity=stripclub.
- **`planetarium`**
  A planetarium.
- **`public_bookcase`**
  A street furniture containing books. Take one or leave one.
- **`social_centre`**
  A place for free and not-for-profit activities.
- **`stage`**
  A raised platform for performers.
- **`stripclub`**
  A place that offers striptease or lapdancing (for sexual services use amenity=brothel).
- **`studio`**
  TV radio or recording studio
- **`swingerclub`**
  A club where people meet to have a party and group sex.
- **`theatre`**
  A theatre or opera house where live performances occur, such as plays, musicals and formal concerts. Use amenity=cinema for movie theaters.
- **`courthouse`**
  A building home to a court of law, where justice is dispensed
- **`fire_station`**
  A station of a fire brigade
- **`police`**
  A police station where police officers patrol from and that is a first point of contact for civilians
- **`post_box`**
  A box for the reception of mail. Alternative mail-carriers can be tagged via operator=*
- **`post_depot`**
  Post depot or delivery office, where letters and parcels are collected and sorted prior to delivery.
- **`post_office`**
  Post office building with postal services
- **`prison`**
  A prison or jail where people are incarcerated before trial or after conviction
- **`ranger_station`**
  National Park visitor headquarters: official park visitor facility with police, visitor information, permit services, etc
- **`townhall`**
  Building where the administration of a village, town or city may be located, or just a community meeting place
- **`bbq`**
  BBQ or Barbecue is a permanently built grill for cooking food, which is most typically used outdoors by the public. For example these may be found in city parks or at beaches. Use the tag fuel=* to specify the source of heating, such as fuel=wood;electric;charcoal. For mapping nearby table and chairs, see also the tag tourism=picnic_site. For mapping campfires and firepits, instead use the tag leisure=firepit.
- **`bench`**
  A bench to sit down and relax a bit
- **`check_in`**
  Place where passengers can get their boarding passes before travel (typically found in airports).
- **`dog_toilet`**
  Area designated for dogs to urinate and excrete.
- **`dressing_room`**
  Area designated for changing clothes.
- **`drinking_water`**
  Drinking water is a place where humans can obtain potable water for consumption. Typically, the water is used for only drinking. Also known as a drinking fountain or bubbler.
- **`give_box`**
  A small facility where people drop off and pick up various types of items in the sense of free sharing and reuse.
- **`lounge`**
  A comfortable waiting area for customers, usually found in airports and other transportation hubs. Typically has extra amenities or sustenance.
- **`mailroom`**
  A mailroom for receiving packages or letters.
- **`parcel_locker`**
  Machine for picking up and sending parcels
- **`shelter`**
  A small shelter against bad weather conditions. To additionally describe the kind of shelter use shelter_type=*.
- **`shower`**
  Public shower.
- **`telephone`**
  Public telephone
- **`toilets`**
  Public toilets (might require a fee)
- **`water_point`**
  Place where you can get large amounts of drinking water
- **`watering_place`**
  Place where water is contained and animals can drink
- **`sanitary_dump_station`**
  A place for depositing human waste from a toilet holding tank.
- **`recycling`**
  Recycling facilities (bottle banks, etc.). Combine with recycling_type=container for containers or recycling_type=centre for recycling centres.
- **`waste_basket`**
  A single small container for depositing garbage that is easily accessible for pedestrians.
- **`waste_disposal`**
  A medium or large disposal bin, typically for bagged up household or industrial waste.
- **`waste_transfer_station`**
  A waste transfer station is a location that accepts, consolidates and transfers waste in bulk.
- **`animal_boarding`**
  A facility where you, paying a fee, can bring your animal for a limited period of time (e.g. for holidays)
- **`animal_breeding`**
  A facility where animals are bred, usually to sell them
- **`animal_shelter`**
  A shelter that recovers animals in trouble
- **`animal_training`**
  A facility used for non-competitive animal training
- **`baking_oven`**
  An oven used for baking bread and similar, for example inside a building=bakehouse.
- **`clock`**
  A public visible clock
- **`crematorium`**
  A place where dead human bodies are burnt
- **`dive_centre`**
  A dive center is the base location where sports divers usually start scuba diving or make dive guided trips at new locations.
- **`funeral_hall`**
  A place for holding a funeral ceremony, other than a place of worship.
- **`grave_yard`**
  A (smaller) place of burial, often you'll find a church nearby. Large places should be landuse=cemetery instead.
- **`hunting_stand`**
  A hunting stand: an open or enclosed platform used by hunters to place themselves at an elevated height above the terrain
- **`internet_cafe`**
  A place whose principal role is providing internet services to the public.
- **`kitchen`**
  A public kitchen in a facility to use by everyone or customers
- **`kneipp_water_cure`**
  Outdoor foot bath facility. Usually this is a pool with cold water and handrail. Popular in German speaking countries.
- **`lounger`**
  An object for people to lie down.
- **`marketplace`**
  A marketplace where goods and services are traded daily or weekly.
- **`monastery`**
  Monastery is the location of a monastery or a building in which monks and nuns live.
- **`mortuary`**
  A morgue or funeral home, used for the storage of human corpses.
- **`photo_booth`**
  A stand to create instant photos.
- **`place_of_mourning`**
  A room or building where families and friends can come, before the funeral, and view the body of the person who has died.
- **`place_of_worship`**
  A church, mosque, or temple, etc. Note that you also need religion=*, usually denomination=* and preferably name=* as well as amenity=place_of_worship. See the article for details.
- **`public_bath`**
  A location where the public may bathe in common, etc. japanese onsen, turkish bath, hot spring
- **`public_building`**
  A generic public building. Don't use! See office=government.
- **`refugee_site`**
  A human settlement sheltering refugees or internally displaced persons
- **`vending_machine`**
  A machine selling goods – food, tickets, newspapers, etc. Add type of goods using vending=*
- **`hydrant`**
  Similar to a fire_hydrant=*, but for gardening and other municipal purposes other than fire extinction

#### `amenity=bar`
The tag `amenity=bar` is used to map a bar: a purpose-built commercial establishment that sells alcoholic drinks to be consumed on the premises. They are characterised by a noisy and vibrant atmosphere, similar to a party. They usually do not sell food to be eaten as a meal. The music is usually loud and you often have to stand. Sometimes it has a dancefloor, but it's not the main attraction.

Whereas pubs (`amenity=pub`) tend to have a similar appearance to traditional houses, bars usually have a more commercial appearance.

In Mediterranean countries, the word "bar" has a different meaning (although this doesn't necessarily mean the tag should be applied differently). Here a bar is integral part of the lifestyle. You go there in the morning to have breakfast, at lunch they serve simple meals, all day long (if not closed after lunch) people use them to get a quick coffee and in the evening it's a meeting place to get an aperitif before dinner. Some are open in the evening and night though mostly they close in the evening, some also sell tobacco, sweets and stamps. Unlike a pub this kind of bar is open for breakfast and coffee plays a way bigger role than beer.

If the bar is not the principal activity (for example a hotel or restaurant that also has a bar) then that could be tagged with the principal activity and also tagged with `bar=yes`.

##### Other tags used in combination
- `operator=*`
- `addr:*=*` – address
- `website=*`
- `opening_hours=*`
- `opening_hours:kitchen=*` – If kitchen opening hours differs from opening hours
- `wheelchair=*`
- `drink:*=*`
- `microbrewery=*` – to indicate that the bar brew and sell their own beer
- `brewery=*` – which breweries' beer they have available
- `distillery=*`
- `winery=*`
- `live_music=yes` – when there is live music in the evening

#### `amenity=cafe`
`amenity=cafe` (café) is for a generally informal place with sit-down facilities selling beverages and light meals and/or snacks. This includes coffee-shops and tea shops selling perhaps tea, coffee and cakes, through to bistros selling meals with alcoholic drinks.

##### Other tags used in combination
* `drink:*=yes/no` - to specify which drinks are served
  * `drink:coffee=yes/no` - whether the cafe is serving coffee
  * `drink:espresso=*` - Indicates whether a feature serves espresso.
* `food:*=yes/no` - to specify which food is served
  * `food:cake=yes/no` - whether the cafe is serving cake
  * `food:bagel=yes/no` - whether the cafe is serving bagel
  * `food:donut=yes/no` - whether the cafe is serving donut
* `cuisine=*`
  * `cuisine=bistro`
  * `cuisine=coffee_shop`
  * `cuisine=ice_cream`
  * `cuisine=donut`
  * `cuisine=bagel`
* `diet:*=*` – for describing dietary choices on offer such as vegetarian.
* `opening_hours=*` – for describing the opening hours.
* `internet_access=*`
* `addr:*=*`
* `phone=*`
* `website=*`
* `wheelchair=*`
* `toilets=yes/no`
  * `toilets:access=customers`
* `operator=*`
* `brand=*`
* `smoking=*`
* `self_service=yes`
* `takeaway=yes`
* `reusable_packaging:accept=*`
* `reusable_packaging:offer=*`
* `capacity=*` – number of seats at tables
* `indoor_seating=yes/no` – for describing the café's indoor seating area if present
* `outdoor_seating=yes/no` – for describing the café's outdoor seating area if present
* `ice_cream=yes/no` – it is used to indicate whether a café sells ice cream
* `bakery=yes` – if the café bakes bread there
* `pastry=yes` – if the café sells pastry
* `drinking_straw=*` - to describe whether the café offers drinking straws and what material they are
* `laptop=yes` if laptops are explicitly allowed and `laptop:conditional=no` to specify laptop-free times or zones.

#### `amenity=pub`
A pub or public house is an establishment that sells alcoholic drinks that can be consumed on, or off, the premises. Pubs often sell food which also can be eaten on the premises although 'Wet Pubs' which do not sell food are also are also common. They are characterised by a traditional appearance and a relaxed atmosphere. You can usually sit down and there is usually no loud music to disturb conversation. A pub would be a good location to meet after a day's mapping for OpenStreetMap.

##### Other tags used in combination
- `operator=*`
- `addr:*=*`
- `opening_hours=*`
  - `opening_hours:kitchen=*` – If kitchen opening hours differs from opening hours
- `contact=*`
- `cuisine=*`
- `drink:*=*`
- `brewery=*` - name(s) of breweries whose beer is sold here
- `distillery=*`
- `winery=*`
- `microbrewery=yes` - pub has a microbrewery on the premises
- `live_music=yes` - when there is live music in the evening
- `smoking=yes/no`
- `internet_access=wlan` - if the pub has wifi for clients.
- `outdoor_seating=yes` - The pub has outside seating, possibly café-style. Useful for urban situations where one can't really describe the seating as garden-like
- `food=yes` - If the pub serves food, or type of food thereof
- `cocktails=yes` - if cocktails are sold
- `old_name=*` - former name of premises
- `real_cider=yes` - real cider/perry offered
- `real_ale=yes` - traditional draught cask beers on tap
- `accommodation=yes/<style>` - Accommodation offered, or type thereof
- `camra=yes` - Listed in the CAMRA Good Beer Guide (UK only).
- `real_fire=yes` - if the pub has a 'real fire' burning solid fuel, i.e. wood, coal or peat.
- `dog=*` - to indicate if dogs are allowed in the pub

#### `amenity=restaurant`
`amenity=restaurant` is applied to generally formal eating places with sit-down facilities selling full meals served by waiters and often licensed (where allowed) to sell alcoholic drinks.

##### Other tags used in combination
* `cuisine=*` – type of cuisine or food
* `opening_hours=*` – opening hours
* `takeaway=yes/no/only` – whether the place offers takeaway
* `delivery=yes/no` – whether the place offers delivery
* `outdoor_seating=*` – whether the place has an outdoor seating area
* `reservation=*` – whether reservation is possible/required
* `operator=*` – operator
* `addr:*=*` – address
* `wheelchair=*` – wheelchair accessibility
* diets, see `diet:*=*`, e.g.
  * `diet:vegetarian=yes/no/only`
  * `diet:vegan=yes/no/only`
* `start_date=*` – start date
* `air_conditioning=yes/no` – whether there is air conditioning
* `capacity=*` – number of seats at tables
* `internet_access=*` – Internet connection
* `payment:*=yes/no/interval` – methods of payment
* `smoking=*` – smoking / non smoking
* `toilets=yes` with `toilets:access=customers` – whether offers toilets to their customers
* `breakfast=*`, `lunch=*` and `dinner=*` — whether and when the place offers breakfast/lunch/dinner
* `image=*` – photo
* `wikimedia_commons=*` – File name of a photo from Wikimedia Commons
* `reusable_packaging:accept=*`
* `reusable_packaging:offer=*`
* `website:menu=*` – url to menu of the restaurant (pdf or webpage)
* `opening_hours:url=*` - url to Opening hours
* `reservation:website=*` - url to make a reservation
* `website=*` – website
* `phone=*` – phone number
* `bar=yes` – the presence of a bar with alcoholic beverages
* `microbrewery=yes` – a microbrewery is located in the restaurant
* `brand=*` – brand
* `level=*` – floor number, if the restaurant is located in a multi-story building, such as a shopping mall
* `stars=*` – restaurant with Michelin star
* `opening_hours:kitchen=*` – If kitchen opening hours differs from opening hours
