CREATE TABLE IF NOT EXISTS telemetry (
  packet_uid                   UBIGINT,
  packet_4cc                   TEXT,
  game_delta_time              FLOAT,
  game_frame_count             UBIGINT,
  game_total_time              FLOAT,
  game_mode                    USMALLINT,
  shiftlights_fraction         FLOAT,
  shiftlights_rpm_end          FLOAT,
  shiftlights_rpm_start        FLOAT,
  shiftlights_rpm_valid        BOOLEAN,
  stage_current_distance       DOUBLE,
  stage_current_time           FLOAT,
  stage_length                 DOUBLE,
  stage_previous_split_time    FLOAT,
  stage_result_time            FLOAT,
  stage_result_time_penalty    FLOAT,
  stage_result_status          USMALLINT,
  stage_progress               FLOAT,
  stage_shakedown              BOOLEAN,
  vehicle_acceleration_x       FLOAT,
  vehicle_acceleration_y       FLOAT,
  vehicle_acceleration_z       FLOAT,
  vehicle_brake                FLOAT,
  vehicle_brake_temperature_bl FLOAT,
  vehicle_brake_temperature_br FLOAT,
  vehicle_brake_temperature_fl FLOAT,
  vehicle_brake_temperature_fr FLOAT,
  vehicle_clutch               FLOAT,
  vehicle_cluster_abs          BOOLEAN,
  vehicle_cp_forward_speed_bl  FLOAT,
  vehicle_cp_forward_speed_br  FLOAT,
  vehicle_cp_forward_speed_fl  FLOAT,
  vehicle_cp_forward_speed_fr  FLOAT,
  vehicle_engine_rpm_current   FLOAT,
  vehicle_engine_rpm_idle      FLOAT,
  vehicle_engine_rpm_max       FLOAT,
  vehicle_forward_direction_x  FLOAT,
  vehicle_forward_direction_y  FLOAT,
  vehicle_forward_direction_z  FLOAT,
  vehicle_gear_index           USMALLINT,
  vehicle_gear_index_neutral   USMALLINT,
  vehicle_gear_index_reverse   USMALLINT,
  vehicle_gear_maximum         USMALLINT,
  vehicle_handbrake            FLOAT,
  vehicle_hub_position_bl      FLOAT,
  vehicle_hub_position_br      FLOAT,
  vehicle_hub_position_fl      FLOAT,
  vehicle_hub_position_fr      FLOAT,
  vehicle_hub_velocity_bl      FLOAT,
  vehicle_hub_velocity_br      FLOAT,
  vehicle_hub_velocity_fl      FLOAT,
  vehicle_hub_velocity_fr      FLOAT,
  vehicle_id                   USMALLINT,
  vehicle_class_id             USMALLINT,
  vehicle_manufacturer_id      USMALLINT,
  vehicle_left_direction_x     FLOAT,
  vehicle_left_direction_y     FLOAT,
  vehicle_left_direction_z     FLOAT,
  vehicle_position_x           FLOAT,
  vehicle_position_y           FLOAT,
  vehicle_position_z           FLOAT,
  vehicle_speed                FLOAT,
  vehicle_steering             FLOAT,
  vehicle_throttle             FLOAT,
  vehicle_transmission_speed   FLOAT,
  vehicle_tyre_state_bl        USMALLINT,
  vehicle_tyre_state_br        USMALLINT,
  vehicle_tyre_state_fl        USMALLINT,
  vehicle_tyre_state_fr        USMALLINT,
  vehicle_up_direction_x       FLOAT,
  vehicle_up_direction_y       FLOAT,
  vehicle_up_direction_z       FLOAT,
  vehicle_velocity_x           FLOAT,
  vehicle_velocity_y           FLOAT,
  vehicle_velocity_z           FLOAT,
  location_id                  USMALLINT,
  route_id                     USMALLINT
);

COMMENT ON COLUMN telemetry.packet_uid IS 'A rolling unique identifier for the current packet. Can be used to order and drop received packets. [count]';
COMMENT ON COLUMN telemetry.packet_4cc IS 'An automatic channel that exports the packet''s 4CC code. This allows the use of packet headers to receive UDP on a single socket with reliable network connections.';
COMMENT ON COLUMN telemetry.game_delta_time IS 'Time spent since last frame. [second]';
COMMENT ON COLUMN telemetry.game_frame_count IS 'Frame count in game since boot. [count]';
COMMENT ON COLUMN telemetry.game_total_time IS 'Time spent in game since boot. [second]';
COMMENT ON COLUMN telemetry.game_mode IS 'Game mode unique identifier. See "game_mode" table.';
COMMENT ON COLUMN telemetry.shiftlights_fraction IS 'For shift lights, from 0 ("vehicle_engine_rpm_current"="shiftlights_rpm_start") to 1 ("vehicle_engine_rpm_current"="shiftlights_rpm_end").';
COMMENT ON COLUMN telemetry.shiftlights_rpm_end IS 'Shift lights end (i.e. optimal shift) at "vehicle_engine_rpm_current" value. [revolution per minute]';
COMMENT ON COLUMN telemetry.shiftlights_rpm_start IS 'Shift lights start at "vehicle_engine_rpm_current" value. [revolution per minute]';
COMMENT ON COLUMN telemetry.shiftlights_rpm_valid IS 'Are shift lights RPM data valid: "vehicle_engine_rpm_current", "shiftlights_rpm_start", "shiftlights_rpm_end"';
COMMENT ON COLUMN telemetry.stage_current_distance IS 'Distance reached on current stage. [metre]';
COMMENT ON COLUMN telemetry.stage_current_time IS 'Time spent on current stage. [second]';
COMMENT ON COLUMN telemetry.stage_length IS 'Total length of current stage. [metre]';
COMMENT ON COLUMN telemetry.stage_previous_split_time IS 'Split time of previous sector. Value unspecified if no sector completed. [second]';
COMMENT ON COLUMN telemetry.stage_result_time IS 'Result time. Has a small delay to update after crossing finish line, but guaranteed correct by the time the race telemetry session has ended. Does not take into account retirement or disqualification. [second]';
COMMENT ON COLUMN telemetry.stage_result_time_penalty IS 'Total time penalty gained [second]';
COMMENT ON COLUMN telemetry.stage_result_status IS 'Unique identifier for stage result status (not finished, finished, disqualified etc.). See "stage_result_status table"';
COMMENT ON COLUMN telemetry.stage_progress IS 'Percentage of stage progress 0 to 1 (during race). Value is unspecified before start line and after finish line.';
COMMENT ON COLUMN telemetry.stage_shakedown IS 'Is in a shakedown event.';
COMMENT ON COLUMN telemetry.vehicle_acceleration_x IS 'Car acceleration X component, positive left. [metre per second squared]';
COMMENT ON COLUMN telemetry.vehicle_acceleration_y IS 'Car acceleration Y component, positive up. [metre per second squared]';
COMMENT ON COLUMN telemetry.vehicle_acceleration_z IS 'Car acceleration Z component, positive forward. [metre per second squared]';
COMMENT ON COLUMN telemetry.vehicle_brake IS 'Brake pedal after assists and overrides, 0 (off) to 1 (full).';
COMMENT ON COLUMN telemetry.vehicle_brake_temperature_bl IS 'Brake temperature, back left. [degree Celsius]';
COMMENT ON COLUMN telemetry.vehicle_brake_temperature_br IS 'Brake temperature, back right. [degree Celsius]';
COMMENT ON COLUMN telemetry.vehicle_brake_temperature_fl IS 'Brake temperature, front left. [degree Celsius]';
COMMENT ON COLUMN telemetry.vehicle_brake_temperature_fr IS 'Brake temperature, front right. [degree Celsius]';
COMMENT ON COLUMN telemetry.vehicle_clutch IS 'Clutch pedal after assists and overrides, 0 (off) to 1 (full).';
COMMENT ON COLUMN telemetry.vehicle_cluster_abs IS 'Anti-lock Braking System light active on vehicle cluster.';
COMMENT ON COLUMN telemetry.vehicle_cp_forward_speed_bl IS 'Contact patch forward speed, back left. [metre per second]';
COMMENT ON COLUMN telemetry.vehicle_cp_forward_speed_br IS 'Contact patch forward speed, back right. [metre per second]';
COMMENT ON COLUMN telemetry.vehicle_cp_forward_speed_fl IS 'Contact patch forward speed, front left. [metre per second]';
COMMENT ON COLUMN telemetry.vehicle_cp_forward_speed_fr IS 'Contact patch forward speed, front right. [metre per second]';
COMMENT ON COLUMN telemetry.vehicle_engine_rpm_current IS 'Engine rotation rate, current. [revolution per minute]';
COMMENT ON COLUMN telemetry.vehicle_engine_rpm_idle IS 'Engine rotation rate, at idle. [revolution per minute]';
COMMENT ON COLUMN telemetry.vehicle_engine_rpm_max IS 'Engine rotation rate, maximum. [revolution per minute]';
COMMENT ON COLUMN telemetry.vehicle_forward_direction_x IS 'Car forward unit vector X component, positive left.';
COMMENT ON COLUMN telemetry.vehicle_forward_direction_y IS 'Car forward unit vector Y component, positive up.';
COMMENT ON COLUMN telemetry.vehicle_forward_direction_z IS 'Car forward unit vector Z component, positive forward.';
COMMENT ON COLUMN telemetry.vehicle_gear_index IS 'Gear index or value of "vehicle_gear_index_neutral" or "vehicle_gear_index_reverse"';
COMMENT ON COLUMN telemetry.vehicle_gear_index_neutral IS '"vehicle_gear_index" if gearbox in Neutral.';
COMMENT ON COLUMN telemetry.vehicle_gear_index_reverse IS '"vehicle_gear_index" if gearbox in Reverse.';
COMMENT ON COLUMN telemetry.vehicle_gear_maximum IS 'Number of forward gears.';
COMMENT ON COLUMN telemetry.vehicle_handbrake IS 'Handbrake after assists and overrides, 0 (off) to 1 (full).';
COMMENT ON COLUMN telemetry.vehicle_hub_position_bl IS 'Wheel hub height displacement, back left, positive up. [metre]';
COMMENT ON COLUMN telemetry.vehicle_hub_position_br IS 'Wheel hub height displacement, back right, positive up. [metre]';
COMMENT ON COLUMN telemetry.vehicle_hub_position_fl IS 'Wheel hub height displacement, front left, positive up. [metre]';
COMMENT ON COLUMN telemetry.vehicle_hub_position_fr IS 'Wheel hub height displacement, front right, positive up. [metre]';
COMMENT ON COLUMN telemetry.vehicle_hub_velocity_bl IS 'Wheel hub vertical velocity, back left, positive up. [metre per second]';
COMMENT ON COLUMN telemetry.vehicle_hub_velocity_br IS 'Wheel hub vertical velocity, back right, positive up. [metre per second]';
COMMENT ON COLUMN telemetry.vehicle_hub_velocity_fl IS 'Wheel hub vertical velocity, front left, positive up. [metre per second]';
COMMENT ON COLUMN telemetry.vehicle_hub_velocity_fr IS 'Wheel hub vertical velocity, front right, positive up. [metre per second]';
COMMENT ON COLUMN telemetry.vehicle_id IS 'Vehicle unique identifier. See "vehicles" table.';
COMMENT ON COLUMN telemetry.vehicle_class_id IS 'Vehicle class unique identifier. See "vehicle_classes" table.';
COMMENT ON COLUMN telemetry.vehicle_manufacturer_id IS 'Vehicle manufacturer unique identifier. See "vehicle_manufacturers" table.';
COMMENT ON COLUMN telemetry.vehicle_left_direction_x IS 'Car left unit vector X component, positive left.';
COMMENT ON COLUMN telemetry.vehicle_left_direction_y IS 'Car left unit vector Y component, positive up.';
COMMENT ON COLUMN telemetry.vehicle_left_direction_z IS 'Car left unit vector Z component, positive forward.';
COMMENT ON COLUMN telemetry.vehicle_position_x IS 'Car position X component, positive left. [metre]';
COMMENT ON COLUMN telemetry.vehicle_position_y IS 'Car position Y component, positive up. [metre]';
COMMENT ON COLUMN telemetry.vehicle_position_z IS 'Car position Z component, positive forward. [metre]';
COMMENT ON COLUMN telemetry.vehicle_speed IS 'Car body speed. [metre per second]';
COMMENT ON COLUMN telemetry.vehicle_steering IS 'Steering after assists and overrides, -1 (full left) to 1 (full right).';
COMMENT ON COLUMN telemetry.vehicle_throttle IS 'Throttle pedal after assists and overrides, 0 (off) to 1 (full).';
COMMENT ON COLUMN telemetry.vehicle_transmission_speed IS 'Car speed at wheel/road due to transmission (for speedo use). NB. May differ from "vehicle_speed". [metre per second]';
COMMENT ON COLUMN telemetry.vehicle_tyre_state_bl IS 'Unique identifier for tyre state (e.g. undamaged, punctured), back left. See "vehicle_tyre_state" table.';
COMMENT ON COLUMN telemetry.vehicle_tyre_state_br IS 'Unique identifier for tyre state (e.g. undamaged, punctured), back right. See "vehicle_tyre_state" table.';
COMMENT ON COLUMN telemetry.vehicle_tyre_state_fl IS 'Unique identifier for tyre state (e.g. undamaged, punctured), front left. See "vehicle_tyre_state" table.';
COMMENT ON COLUMN telemetry.vehicle_tyre_state_fr IS 'Unique identifier for tyre state (e.g. undamaged, punctured), front right. See "vehicle_tyre_state" table.';
COMMENT ON COLUMN telemetry.vehicle_up_direction_x IS 'Car up unit vector X component, positive left.';
COMMENT ON COLUMN telemetry.vehicle_up_direction_y IS 'Car up unit vector Y component, positive up.';
COMMENT ON COLUMN telemetry.vehicle_up_direction_z IS 'Car up unit vector Z component, positive forward.';
COMMENT ON COLUMN telemetry.vehicle_velocity_x IS 'Car velocity X component, positive left. [metre per second]';
COMMENT ON COLUMN telemetry.vehicle_velocity_y IS 'Car velocity Y component, positive up. [metre per second]';
COMMENT ON COLUMN telemetry.vehicle_velocity_z IS 'Car velocity Z component, positive forward. [metre per second]';
COMMENT ON COLUMN telemetry.location_id IS 'Location unique identifier. See "locations" table.';
COMMENT ON COLUMN telemetry.route_id IS 'Route unique identifier. See "routes" table.';

CREATE TABLE IF NOT EXISTS vehicle_classes (
  id USMALLINT PRIMARY KEY,
  name TEXT
);

INSERT OR IGNORE INTO vehicle_classes(id, name) VALUES
  (1, 'H1 (FWD)'),
  (2, 'H2 (FWD)'),
  (3, 'H2 (RWD)'),
  (4, 'H3 (RWD)'),
  (5, 'Group B (RWD)'),
  (6, 'Group B (4WD)'),
  (7, 'NR4/R4'),
  (8, 'Rally4'),
  (9, 'S1600'),
  (10, 'S2000'),
  (11, 'F2 Kit Car'),
  (12, 'Group A'),
  (13, 'World Rally Car 1997 - 2011'),
  (14, 'Rally2'),
  (17, 'World Rally Car 2017 - 2021'),
  (19, 'WRC'),
  (20, 'JWRC'),
  (21, 'WRC2'),
  (23, 'World Rally Car 2012 - 2016');

CREATE TABLE IF NOT EXISTS vehicle_manufacturers (
  id USMALLINT PRIMARY KEY,
  name TEXT
);

INSERT OR IGNORE INTO vehicle_manufacturers(id, name) VALUES
  (1, 'lng_manufacturer_name_test_manufacturer_1'),
  (2, 'lng_manufacturer_name_test_manufacturer_2'),
  (3, 'lng_manufacturer_name_test_manufacturer_3'),
  (4, 'Citroën'),
  (5, 'Ford'),
  (6, 'Hyundai'),
  (7, 'Peugeot'),
  (8, 'ŠKODA'),
  (9, 'Volkswagen'),
  (10, 'Proton'),
  (11, 'Mitsubishi'),
  (12, 'SUBARU'),
  (13, 'SEAT'),
  (14, 'Lancia'),
  (15, 'Toyota'),
  (16, 'Vauxhall'),
  (17, 'Renault'),
  (18, 'MG'),
  (19, 'Suzuki'),
  (20, 'Opel'),
  (21, 'Audi'),
  (22, 'BMW'),
  (23, 'Nissan'),
  (24, 'Porsche'),
  (25, 'DJM Motorsport'),
  (26, 'Chrysler'),
  (27, 'Fiat'),
  (28, 'Lada'),
  (29, 'Talbot'),
  (30, 'Alpine'),
  (31, 'DS Automobiles'),
  (32, 'MINI'),
  (33, 'Datsun'),
  (34, 'Builder'),
  (35, 'Hillman');

CREATE TABLE IF NOT EXISTS vehicles (
  id USMALLINT PRIMARY KEY,
  class USMALLINT REFERENCES vehicle_classes(id),
  manufacturer USMALLINT REFERENCES vehicle_manufacturers(id),
  name TEXT,
  builder BOOLEAN
);

INSERT OR IGNORE INTO vehicles(id, class, manufacturer, name, builder) VALUES
  (4, 21, 4, 'Citroën C3 Rally2', false),
  (5, 21, 5, 'Ford Fiesta Rally2', false),
  (7, 8, 7, 'Peugeot 208 Rally4', false),
  (8, 21, 8, 'ŠKODA Fabia Rally2 Evo', false),
  (9, 21, 9, 'Volkswagen Polo GTI R5', false),
  (12, 13, 4, 'Citroën C4 WRC', false),
  (13, 13, 4, 'Citroën Xsara WRC', false),
  (14, 13, 5, 'Ford Focus RS Rally 2001', false),
  (15, 13, 5, 'Ford Focus RS Rally 2008', false),
  (16, 13, 7, 'Peugeot 206 Rally', false),
  (17, 13, 8, 'ŠKODA Fabia WRC', false),
  (18, 13, 12, 'SUBARU Impreza 2001', false),
  (19, 13, 12, 'SUBARU Impreza 1998', false),
  (20, 13, 12, 'SUBARU Impreza 2008', false),
  (21, 13, 13, 'SEAT Córdoba WRC', false),
  (22, 12, 5, 'Ford Escort RS Cosworth', false),
  (23, 12, 14, 'Lancia Delta HF Integrale', false),
  (24, 13, 11, 'Mitsubishi Lancer Evolution VI', false),
  (26, 12, 12, 'SUBARU Legacy RS', false),
  (29, 11, 7, 'Peugeot 306 Maxi', false),
  (30, 11, 9, 'Volkswagen Golf IV Kit Car', false),
  (31, 11, 13, 'SEAT Ibiza Kit Car', false),
  (32, 11, 16, 'Vauxhall Astra Rally Car', false),
  (33, 11, 5, 'Ford Escort Mk 6 Maxi', false),
  (34, 11, 17, 'Renault Maxi Mégane', false),
  (35, 10, 27, 'Fiat Grande Punto Abarth S2000', false),
  (36, 10, 7, 'Peugeot 207 S2000', false),
  (41, 9, 4, 'Citroën Saxo Super 1600', false),
  (42, 9, 17, 'Renault Clio S1600', false),
  (43, 9, 5, 'Ford Puma S1600', false),
  (44, 9, 4, 'Citroën C2 Super 1600', false),
  (45, 8, 5, 'Ford Fiesta MK8 Rally4', false),
  (46, 8, 20, 'Opel Adam R2', false),
  (48, 8, 17, 'Renault Twingo II', false),
  (50, 6, 21, 'Audi Sport quattro S1 (E2)', false),
  (51, 6, 5, 'Ford RS200', false),
  (52, 6, 14, 'Lancia Delta S4', false),
  (53, 6, 18, 'MG Metro 6R4', false),
  (54, 6, 7, 'Peugeot 205 T16 Evo 2', false),
  (55, 5, 22, 'BMW M1 Procar Rally', false),
  (57, 5, 14, 'Lancia 037 Evo 2', false),
  (58, 5, 20, 'Opel Manta 400', false),
  (59, 5, 24, 'Porsche 911 SC RS', false),
  (60, 7, 11, 'Mitsubishi Lancer Evolution X', false),
  (61, 7, 12, 'SUBARU WRX STI NR4', false),
  (62, 7, 25, 'McRae R4', false),
  (63, 4, 22, 'BMW M3 Evo Rally', false),
  (64, 4, 14, 'Lancia Stratos', false),
  (65, 4, 5, 'Ford Sierra Cosworth RS500', false),
  (67, 3, 35, 'Hillman Avenger', false),
  (68, 4, 20, 'Opel Ascona 400', false),
  (69, 4, 17, 'Renault 5 Turbo', false),
  (70, 3, 5, 'Ford Escort MK2', false),
  (71, 3, 27, 'Fiat 131 Abarth Rally', false),
  (73, 3, 20, 'Opel Kadett C GT/E', false),
  (74, 3, 29, 'Talbot Sunbeam Lotus', false),
  (75, 3, 30, 'Alpine Renault A110 1600 S', false),
  (76, 2, 7, 'Peugeot 205 GTI', false),
  (77, 2, 7, 'Peugeot 309 GTI', false),
  (78, 1, 16, 'Vauxhall Nova Sport', false),
  (79, 2, 9, 'Volkswagen Golf GTI', false),
  (81, 1, 14, 'Lancia Fulvia HF', false),
  (82, 1, 32, 'MINI Cooper S', false),
  (83, 12, 11, 'Mitsubishi Galant VR4', false),
  (85, 14, 7, 'Peugeot 208 T16 R5', false),
  (86, 13, 32, 'MINI Countryman Rally Edition', false),
  (87, 10, 20, 'Opel Corsa S2000', false),
  (89, 12, 12, 'SUBARU Impreza 1995', false),
  (94, 4, 5, 'Ford Escort MK2 McRae Motorsport', false),
  (95, 17, 9, 'Volkswagen Polo 2017', false),
  (99, 21, 6, 'Hyundai i20 N Rally2', false),
  (100, 21, 34, 'WRC2 Builder Vehicle', true),
  (101, 19, 34, 'WRC Builder Vehicle', true),
  (102, 20, 34, 'Junior WRC Builder Vehicle', true),
  (103, 19, 5, 'Ford Puma Rally1 HYBRID', false),
  (104, 19, 6, 'Hyundai i20 N Rally1 HYBRID', false),
  (105, 19, 15, 'Toyota GR Yaris Rally1 HYBRID', false),
  (106, 20, 5, 'Ford Fiesta Rally3', false),
  (107, 21, 8, 'ŠKODA Fabia RS Rally2', false),
  (112, 14, 5, 'Ford Fiesta R5 MK7 Evo 2', false),
  (116, 14, 9, 'Volkswagen Polo GTI R5', false),
  (117, 14, 8, 'ŠKODA Fabia Rally2 Evo', false),
  (118, 17, 5, 'Ford Fiesta WRC', false),
  (119, 14, 8, 'ŠKODA Fabia RS Rally2', false),
  (120, 21, 15, 'Toyota GR Yaris Rally2', false),
  (121, 20, 5, 'Ford Fiesta Rally3 Evo', false),
  (124, 19, 5, 'Ford Puma Rally1 HYBRID', false),
  (125, 19, 6, 'Hyundai i20 N Rally1 HYBRID', false),
  (126, 19, 15, 'Toyota GR Yaris Rally1 HYBRID', false),
  (128, 8, 4, 'Citroën C2 R2 Max', false),
  (129, 11, 4, 'Citroën Xsara Kit Car', false),
  (130, 17, 4, 'Citroën C3 WRC', false),
  (132, 23, 4, 'Citroën DS3 WRC ''12', false),
  (133, 23, 9, 'Volkswagen Polo R WRC 2013', false),
  (134, 23, 32, 'MINI John Cooper Works WRC', false),
  (136, 9, 7, 'Peugeot 206 S1600', false);

CREATE TABLE IF NOT EXISTS vehicle_tyre_states (
  id USMALLINT PRIMARY KEY,
  name TEXT
);

INSERT OR IGNORE INTO vehicle_tyre_states(id, name) VALUES
  (0, 'undamaged'),
  (1, 'punctured'),
  (2, 'burst');

CREATE TABLE IF NOT EXISTS locations (
  id USMALLINT PRIMARY KEY,
  name TEXT
);

INSERT OR IGNORE INTO locations(id, name) VALUES
  (5, 'RALLY MEDITERRANEO'),
  (6, 'VODAFONE RALLY DE PORTUGAL'),
  (7, 'RALLY ITALIA SARDEGNA'),
  (8, 'RALLY ESTONIA'),
  (9, 'RALLY SCANDIA'),
  (12, 'GUANAJUATO RALLY MÉXICO'),
  (13, 'RALLY CHILE BIO BÍO'),
  (14, 'AGON BY AOC RALLY PACIFICO'),
  (15, 'SECTO RALLY FINLAND'),
  (16, 'CROATIA RALLY'),
  (17, 'RALLYE MONTE-CARLO'),
  (18, 'RALLY SWEDEN'),
  (24, 'EKO ACROPOLIS RALLY GREECE'),
  (25, 'FORUM8 RALLY JAPAN'),
  (26, 'SAFARI RALLY KENYA'),
  (27, 'FANATEC RALLY OCEANIA'),
  (28, 'RALLY IBERIA'),
  (29, 'CENTRAL EUROPE RALLY'),
  (30, 'TET RALLY LATVIA'),
  (31, 'ORLEN 80TH RALLY POLAND');

CREATE TABLE IF NOT EXISTS routes (
  id USMALLINT PRIMARY KEY,
  name TEXT
);

INSERT OR IGNORE INTO routes(id, name) VALUES
  (24, 'Asco'),
  (25, 'Ponte'),
  (26, 'Monte Cinto'),
  (28, 'Albarello'),
  (29, 'Capannace'),
  (30, 'Serra Di Cuzzioli'),
  (32, 'Maririe'),
  (33, 'Poggiola'),
  (34, 'Monte Alloradu'),
  (35, 'Ravin de Finelio'),
  (36, 'Cabanella'),
  (37, 'Moltifao'),
  (57, 'Talao'),
  (58, 'Talanghilirair'),
  (59, 'Sungai Kunit'),
  (60, 'Sangir Balai Janggo'),
  (61, 'South Solok'),
  (63, 'Kebun Raya Solok'),
  (64, 'Batukangkung'),
  (65, 'Abai'),
  (66, 'Moearaikoer'),
  (67, 'Bidaralam'),
  (68, 'Loeboekmalaka'),
  (69, 'Gunung Tujuh'),
  (85, 'Holtjønn'),
  (86, 'Hengeltjønn'),
  (87, 'Fyresvatn'),
  (88, 'Russvatn'),
  (89, 'Tovsli'),
  (91, 'Kottjønn'),
  (93, 'Fordol'),
  (94, 'Fyresdal'),
  (95, 'Ljosdalstjønn'),
  (96, 'Dagtrolltjønn'),
  (97, 'Tovslioytjorn'),
  (98, 'Bergsøytjønn'),
  (99, 'Baião'),
  (100, 'Caminha'),
  (101, 'Fridão'),
  (103, 'Marão'),
  (104, 'Ponte de Lima'),
  (105, 'Viana do Castelo'),
  (106, 'Ervideiro'),
  (108, 'Celeiro'),
  (109, 'Touca'),
  (110, 'Vila Boa'),
  (111, 'Carrazedo'),
  (112, 'Anjos'),
  (113, 'Rena Majore'),
  (114, 'Monte Olia'),
  (115, 'Littichedda'),
  (117, 'Alà del Sardi'),
  (118, 'Mamone'),
  (119, 'Li Pinnenti'),
  (120, 'Malti'),
  (122, 'Bassacutena'),
  (123, 'Bortigiadas'),
  (124, 'Sa Mela'),
  (125, 'Monte Muvri'),
  (126, 'Monte Acuto'),
  (127, 'El Chocolate'),
  (128, 'Otates'),
  (129, 'Ortega'),
  (130, 'Las Minas'),
  (131, 'Ibarrilla'),
  (133, 'Derramadero'),
  (134, 'El Brinco'),
  (136, 'Guanajuatito'),
  (137, 'Alfaro'),
  (138, 'Mesa Cuata'),
  (139, 'San Diego'),
  (140, 'El Mosquito'),
  (155, 'Bio Bío'),
  (156, 'Pulpería'),
  (157, 'Río Lía'),
  (159, 'María Las Cruces'),
  (160, 'Las Pataguas'),
  (161, 'Rere'),
  (162, 'El Poñen'),
  (163, 'Laja'),
  (164, 'Yumbel'),
  (165, 'Río Claro'),
  (166, 'Hualqui'),
  (167, 'Chivilingo'),
  (204, 'Leustu'),
  (205, 'Lahdenkylä'),
  (206, 'Saakoski'),
  (207, 'Maahi'),
  (208, 'Painaa'),
  (209, 'Peltola'),
  (216, 'Päijälä'),
  (217, 'Ruokolahti'),
  (218, 'Honkanen'),
  (219, 'Venkajärvi'),
  (220, 'Vehmas'),
  (221, 'Hatanpää'),
  (229, 'Hof-Finnskog'),
  (230, 'Åsnes'),
  (231, 'Spikbrenna'),
  (232, 'Lauksjøen'),
  (233, 'Åslia'),
  (234, 'Knapptjernet'),
  (241, 'Vargasen'),
  (242, 'Lövstaholm'),
  (243, 'Älgsjön'),
  (244, 'Ekshärad'),
  (245, 'Stora Jangen'),
  (246, 'Sunne'),
  (252, 'La Bollène-Vésubie - Peïra Cava'),
  (253, 'Peïra Cava - La Bollène-Vésubie'),
  (254, 'La Bollène-Vésubie - Col de Turini'),
  (255, 'Pra d''Alart'),
  (256, 'La Maïris'),
  (257, 'Baisse de Patronel'),
  (263, 'Saint-Léger-les-Mélèzes - La Bâtie-Neuve'),
  (264, 'La Bâtie-Neuve - Saint-Léger-les-Mélèzes'),
  (265, 'Moissière'),
  (266, 'Ancelle'),
  (267, 'Ravin de Coste Belle'),
  (268, 'Les Borels'),
  (274, 'Santes Creus'),
  (275, 'Valldossera'),
  (276, 'Campdasens'),
  (277, 'Pontils'),
  (278, 'Montagut'),
  (279, 'Aiguamúrcia'),
  (285, 'Alforja'),
  (286, 'Les Irles'),
  (287, 'L''Argentera'),
  (288, 'Les Voltes'),
  (289, 'Montclar'),
  (290, 'Botareli'),
  (296, 'Gravia'),
  (297, 'Prosilio'),
  (298, 'Mariolata'),
  (299, 'Karoutes'),
  (300, 'Viniani'),
  (301, 'Delphi'),
  (302, 'Eptalofos'),
  (303, 'Lilea'),
  (304, 'Parnassós'),
  (305, 'Bauxites'),
  (306, 'Drosochori'),
  (307, 'Amfissa'),
  (308, 'Lake Mikawa'),
  (309, 'Kudarisawa'),
  (310, 'Oninotaira'),
  (311, 'Okuwacho'),
  (312, 'Habu Dam'),
  (313, 'Habucho'),
  (314, 'Nenoue Plateau'),
  (315, 'Tegano'),
  (316, 'Higashino'),
  (317, 'Hokono Lake'),
  (318, 'Nenoue Highlands'),
  (319, 'Nakatsugawa'),
  (320, 'Otepää'),
  (321, 'Rebaste'),
  (322, 'Nüpli'),
  (323, 'Truuta'),
  (324, 'Koigu'),
  (325, 'Kooraste'),
  (326, 'Elva'),
  (327, 'Metsalaane'),
  (328, 'Vahessaare'),
  (329, 'Külaaseme'),
  (330, 'Vissi'),
  (331, 'Vellavere'),
  (332, 'Malewa'),
  (333, 'Tarambete'),
  (334, 'Moi North'),
  (335, 'Marula'),
  (336, 'Wileli'),
  (337, 'Kingono'),
  (338, 'Soysambu'),
  (339, 'Mbaruk'),
  (340, 'Sugunoi'),
  (341, 'Nakuru'),
  (342, 'Kanyawa'),
  (343, 'Kanyawa - Nakura'),
  (344, 'Bliznec'),
  (345, 'Trakošćan'),
  (346, 'Vrbno'),
  (347, 'Zagorska Sela'),
  (348, 'Kumrovec'),
  (349, 'Grdanjci'),
  (350, 'Stojdraga'),
  (351, 'Mali Lipovec'),
  (352, 'Hartje'),
  (353, 'Kostanjevac'),
  (354, 'Krašić'),
  (355, 'Petruš Vrh'),
  (356, 'Oakleigh'),
  (357, 'Doctors Hill'),
  (358, 'Mangapai'),
  (359, 'Brynderwyn'),
  (360, 'Taipuha'),
  (361, 'Mareretu'),
  (362, 'Waiwera'),
  (363, 'Tahekeroa'),
  (364, 'Noakes Hill'),
  (365, 'Orewa'),
  (366, 'Tahekeroa - Orewa'),
  (367, 'Makarau'),
  (368, 'Rouské'),
  (369, 'Lukoveček'),
  (370, 'Raztoka'),
  (371, 'Žabárna'),
  (372, 'Provodovice'),
  (373, 'Chvalčov'),
  (374, 'Vítová'),
  (375, 'Brusné'),
  (376, 'Libosváry'),
  (377, 'Rusava'),
  (378, 'Osíčko'),
  (379, 'Příkazy'),
  (380, 'Vecpils'),
  (381, 'Kaģene'),
  (382, 'Mazilmāja'),
  (383, 'Ķirsīts'),
  (384, 'Baznīca'),
  (385, 'Stroķacs'),
  (386, 'Podnieki'),
  (387, 'Dinsdurbe'),
  (388, 'Kalvene'),
  (389, 'Cērpi'),
  (390, 'Krote'),
  (391, 'Kalēti'),
  (392, 'Swietajno'),
  (393, 'Jelonek'),
  (394, 'Gajrowskie'),
  (395, 'Pietrasze'),
  (396, 'Dybowo'),
  (397, 'Chełchy'),
  (398, 'Mikolajki'),
  (399, 'Zawada'),
  (400, 'Gmina Mragowo'),
  (401, 'Czerwonki'),
  (402, 'Kosewo'),
  (403, 'Probark'),
  (440, 'Briançonnet - Entrevaux'),
  (441, 'Entrevaux - Briançonnet'),
  (442, 'Les Vénières'),
  (443, 'Parbiou'),
  (444, 'Le champ'),
  (445, 'Pertus'),
  (446, 'Fafe'),
  (447, 'Vila Pouca'),
  (448, 'Barbosa'),
  (449, 'Passos'),
  (450, 'Moreira do Rei'),
  (451, 'Ruivães');

CREATE TABLE IF NOT EXISTS game_modes (
  id USMALLINT PRIMARY KEY,
  name TEXT
);

INSERT OR IGNORE INTO game_modes(id, name) VALUES
  (0, 'time_trial'),
  (1, 'quick_play_multiplayer'),
  (2, 'quick_play_solo'),
  (3, 'esports'),
  (4, 'career'),
  (5, 'test_drive'),
  (6, 'moments'),
  (7, 'clubs'),
  (8, 'championship'),
  (9, 'rally_school');

CREATE TABLE IF NOT EXISTS stage_result_states (
  id USMALLINT PRIMARY KEY,
  name TEXT
);

INSERT OR IGNORE INTO stage_result_states(id, name) VALUES
  (0, 'not_finished'),
  (1, 'finished'),
  (2, 'timed_out_stage'),
  (3, 'terminally_damaged'),
  (4, 'retired'),
  (5, 'disqualified'),
  (6, 'unknown');