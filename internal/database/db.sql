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
  (5, 'Rally Mediterraneo'),
  (6, 'Vodafone Rally De Portugal'),
  (7, 'Rally Italia Sardegna'),
  (8, 'Rally Estonia'),
  (9, 'Rally Scandia'),
  (12, 'Guanajuato Rally México'),
  (13, 'Rally Chile Bio Bío'),
  (14, 'Agon By Aoc Rally Pacifico'),
  (15, 'Secto Rally Finland'),
  (16, 'Croatia Rally'),
  (17, 'Rallye Monte-Carlo'),
  (18, 'Rally Sweden'),
  (24, 'Eko Acropolis Rally Greece'),
  (25, 'Forum8 Rally Japan'),
  (26, 'Safari Rally Kenya'),
  (27, 'Fanatec Rally Oceania'),
  (28, 'Rally Iberia'),
  (29, 'Central Europe Rally'),
  (30, 'Tet Rally Latvia'),
  (31, 'Orlen 80th Rally Poland');

CREATE TABLE IF NOT EXISTS routes (
  id USMALLINT PRIMARY KEY,
  location_id USMALLINT REFERENCES locations(id),
  name TEXT
);

INSERT OR IGNORE INTO routes(id, location_id, name) VALUES
  (24, 5, 'Asco'),
  (25, 5, 'Ponte'),
  (26, 5, 'Monte Cinto'),
  (28, 5, 'Albarello'),
  (29, 5, 'Capannace'),
  (30, 5, 'Serra Di Cuzzioli'),
  (32, 5, 'Maririe'),
  (33, 5, 'Poggiola'),
  (34, 5, 'Monte Alloradu'),
  (35, 5, 'Ravin de Finelio'),
  (36, 5, 'Cabanella'),
  (37, 5, 'Moltifao'),
  (57, 14, 'Talao'),
  (58, 14, 'Talanghilirair'),
  (59, 14, 'Sungai Kunit'),
  (60, 14, 'Sangir Balai Janggo'),
  (61, 14, 'South Solok'),
  (63, 14, 'Kebun Raya Solok'),
  (64, 14, 'Batukangkung'),
  (65, 14, 'Abai'),
  (66, 14, 'Moearaikoer'),
  (67, 14, 'Bidaralam'),
  (68, 14, 'Loeboekmalaka'),
  (69, 14, 'Gunung Tujuh'),
  (85, 9, 'Holtjønn'),
  (86, 9, 'Hengeltjønn'),
  (87, 9, 'Fyresvatn'),
  (88, 9, 'Russvatn'),
  (89, 9, 'Tovsli'),
  (91, 9, 'Kottjønn'),
  (93, 9, 'Fordol'),
  (94, 9, 'Fyresdal'),
  (95, 9, 'Ljosdalstjønn'),
  (96, 9, 'Dagtrolltjønn'),
  (97, 9, 'Tovslioytjorn'),
  (98, 9, 'Bergsøytjønn'),
  (99, 6, 'Baião'),
  (100, 6, 'Caminha'),
  (101, 6, 'Fridão'),
  (103, 6, 'Marão'),
  (104, 6, 'Ponte de Lima'),
  (105, 6, 'Viana do Castelo'),
  (106, 6, 'Ervideiro'),
  (108, 6, 'Celeiro'),
  (109, 6, 'Touca'),
  (110, 6, 'Vila Boa'),
  (111, 6, 'Carrazedo'),
  (112, 6, 'Anjos'),
  (113, 7, 'Rena Majore'),
  (114, 7, 'Monte Olia'),
  (115, 7, 'Littichedda'),
  (117, 7, 'Alà del Sardi'),
  (118, 7, 'Mamone'),
  (119, 7, 'Li Pinnenti'),
  (120, 7, 'Malti'),
  (122, 7, 'Bassacutena'),
  (123, 7, 'Bortigiadas'),
  (124, 7, 'Sa Mela'),
  (125, 7, 'Monte Muvri'),
  (126, 7, 'Monte Acuto'),
  (127, 12, 'El Chocolate'),
  (128, 12, 'Otates'),
  (129, 12, 'Ortega'),
  (130, 12, 'Las Minas'),
  (131, 12, 'Ibarrilla'),
  (133, 12, 'Derramadero'),
  (134, 12, 'El Brinco'),
  (136, 12, 'Guanajuatito'),
  (137, 12, 'Alfaro'),
  (138, 12, 'Mesa Cuata'),
  (139, 12, 'San Diego'),
  (140, 12, 'El Mosquito'),
  (155, 13, 'Bio Bío'),
  (156, 13, 'Pulpería'),
  (157, 13, 'Río Lía'),
  (159, 13, 'María Las Cruces'),
  (160, 13, 'Las Pataguas'),
  (161, 13, 'Rere'),
  (162, 13, 'El Poñen'),
  (163, 13, 'Laja'),
  (164, 13, 'Yumbel'),
  (165, 13, 'Río Claro'),
  (166, 13, 'Hualqui'),
  (167, 13, 'Chivilingo'),
  (204, 15, 'Leustu'),
  (205, 15, 'Lahdenkylä'),
  (206, 15, 'Saakoski'),
  (207, 15, 'Maahi'),
  (208, 15, 'Painaa'),
  (209, 15, 'Peltola'),
  (216, 15, 'Päijälä'),
  (217, 15, 'Ruokolahti'),
  (218, 15, 'Honkanen'),
  (219, 15, 'Venkajärvi'),
  (220, 15, 'Vehmas'),
  (221, 15, 'Hatanpää'),
  (229, 9, 'Hof-Finnskog'),
  (230, 9, 'Åsnes'),
  (231, 9, 'Spikbrenna'),
  (232, 9, 'Lauksjøen'),
  (233, 9, 'Åslia'),
  (234, 9, 'Knapptjernet'),
  (241, 18, 'Vargasen'),
  (242, 18, 'Lövstaholm'),
  (243, 18, 'Älgsjön'),
  (244, 18, 'Ekshärad'),
  (245, 18, 'Stora Jangen'),
  (246, 18, 'Sunne'),
  (252, 17, 'La Bollène-Vésubie - Peïra Cava'),
  (253, 17, 'Peïra Cava - La Bollène-Vésubie'),
  (254, 17, 'La Bollène-Vésubie - Col de Turini'),
  (255, 17, 'Pra d''Alart'),
  (256, 17, 'La Maïris'),
  (257, 17, 'Baisse de Patronel'),
  (263, 17, 'Saint-Léger-les-Mélèzes - La Bâtie-Neuve'),
  (264, 17, 'La Bâtie-Neuve - Saint-Léger-les-Mélèzes'),
  (265, 17, 'Moissière'),
  (266, 17, 'Ancelle'),
  (267, 17, 'Ravin de Coste Belle'),
  (268, 17, 'Les Borels'),
  (274, 28, 'Santes Creus'),
  (275, 28, 'Valldossera'),
  (276, 28, 'Campdasens'),
  (277, 28, 'Pontils'),
  (278, 28, 'Montagut'),
  (279, 28, 'Aiguamúrcia'),
  (285, 28, 'Alforja'),
  (286, 28, 'Les Irles'),
  (287, 28, 'L''Argentera'),
  (288, 28, 'Les Voltes'),
  (289, 28, 'Montclar'),
  (290, 28, 'Botareli'),
  (296, 24, 'Gravia'),
  (297, 24, 'Prosilio'),
  (298, 24, 'Mariolata'),
  (299, 24, 'Karoutes'),
  (300, 24, 'Viniani'),
  (301, 24, 'Delphi'),
  (302, 24, 'Eptalofos'),
  (303, 24, 'Lilea'),
  (304, 24, 'Parnassós'),
  (305, 24, 'Bauxites'),
  (306, 24, 'Drosochori'),
  (307, 24, 'Amfissa'),
  (308, 25, 'Lake Mikawa'),
  (309, 25, 'Kudarisawa'),
  (310, 25, 'Oninotaira'),
  (311, 25, 'Okuwacho'),
  (312, 25, 'Habu Dam'),
  (313, 25, 'Habucho'),
  (314, 25, 'Nenoue Plateau'),
  (315, 25, 'Tegano'),
  (316, 25, 'Higashino'),
  (317, 25, 'Hokono Lake'),
  (318, 25, 'Nenoue Highlands'),
  (319, 25, 'Nakatsugawa'),
  (320, 8, 'Otepää'),
  (321, 8, 'Rebaste'),
  (322, 8, 'Nüpli'),
  (323, 8, 'Truuta'),
  (324, 8, 'Koigu'),
  (325, 8, 'Kooraste'),
  (326, 8, 'Elva'),
  (327, 8, 'Metsalaane'),
  (328, 8, 'Vahessaare'),
  (329, 8, 'Külaaseme'),
  (330, 8, 'Vissi'),
  (331, 8, 'Vellavere'),
  (332, 26, 'Malewa'),
  (333, 26, 'Tarambete'),
  (334, 26, 'Moi North'),
  (335, 26, 'Marula'),
  (336, 26, 'Wileli'),
  (337, 26, 'Kingono'),
  (338, 26, 'Soysambu'),
  (339, 26, 'Mbaruk'),
  (340, 26, 'Sugunoi'),
  (341, 26, 'Nakuru'),
  (342, 26, 'Kanyawa'),
  (343, 26, 'Kanyawa - Nakura'),
  (344, 16, 'Bliznec'),
  (345, 16, 'Trakošćan'),
  (346, 16, 'Vrbno'),
  (347, 16, 'Zagorska Sela'),
  (348, 16, 'Kumrovec'),
  (349, 16, 'Grdanjci'),
  (350, 16, 'Stojdraga'),
  (351, 16, 'Mali Lipovec'),
  (352, 16, 'Hartje'),
  (353, 16, 'Kostanjevac'),
  (354, 16, 'Krašić'),
  (355, 16, 'Petruš Vrh'),
  (356, 27, 'Oakleigh'),
  (357, 27, 'Doctors Hill'),
  (358, 27, 'Mangapai'),
  (359, 27, 'Brynderwyn'),
  (360, 27, 'Taipuha'),
  (361, 27, 'Mareretu'),
  (362, 27, 'Waiwera'),
  (363, 27, 'Tahekeroa'),
  (364, 27, 'Noakes Hill'),
  (365, 27, 'Orewa'),
  (366, 27, 'Tahekeroa - Orewa'),
  (367, 27, 'Makarau'),
  (368, 29, 'Rouské'),
  (369, 29, 'Lukoveček'),
  (370, 29, 'Raztoka'),
  (371, 29, 'Žabárna'),
  (372, 29, 'Provodovice'),
  (373, 29, 'Chvalčov'),
  (374, 29, 'Vítová'),
  (375, 29, 'Brusné'),
  (376, 29, 'Libosváry'),
  (377, 29, 'Rusava'),
  (378, 29, 'Osíčko'),
  (379, 29, 'Příkazy'),
  (380, 30, 'Vecpils'),
  (381, 30, 'Kaģene'),
  (382, 30, 'Mazilmāja'),
  (383, 30, 'Ķirsīts'),
  (384, 30, 'Baznīca'),
  (385, 30, 'Stroķacs'),
  (386, 30, 'Podnieki'),
  (387, 30, 'Dinsdurbe'),
  (388, 30, 'Kalvene'),
  (389, 30, 'Cērpi'),
  (390, 30, 'Krote'),
  (391, 30, 'Kalēti'),
  (392, 31, 'Swietajno'),
  (393, 31, 'Jelonek'),
  (394, 31, 'Gajrowskie'),
  (395, 31, 'Pietrasze'),
  (396, 31, 'Dybowo'),
  (397, 31, 'Chełchy'),
  (398, 31, 'Mikolajki'),
  (399, 31, 'Zawada'),
  (400, 31, 'Gmina Mragowo'),
  (401, 31, 'Czerwonki'),
  (402, 31, 'Kosewo'),
  (403, 31, 'Probark'),
  (440, 17, 'Briançonnet - Entrevaux'),
  (441, 17, 'Entrevaux - Briançonnet'),
  (442, 17, 'Les Vénières'),
  (443, 17, 'Parbiou'),
  (444, 17, 'Le champ'),
  (445, 17, 'Pertus'),
  (446, 6, 'Fafe'),
  (447, 6, 'Vila Pouca'),
  (448, 6, 'Barbosa'),
  (449, 6, 'Passos'),
  (450, 6, 'Moreira do Rei'),
  (451, 6, 'Ruivães');

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

CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  name TEXT
);

CREATE TABLE IF NOT EXISTS user_logins (
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  user_id TEXT REFERENCES users(id),
);

CREATE SEQUENCE IF NOT EXISTS session_id_sequence START 1;

CREATE TABLE IF NOT EXISTS sessions (
  id                        UINTEGER PRIMARY KEY DEFAULT nextval('session_id_sequence'),
  started_at                TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  user_id                   TEXT REFERENCES users(id),
  game_mode                 USMALLINT,
  location_id               USMALLINT,
  route_id                  USMALLINT,
  stage_length              DOUBLE,
  stage_result_status       USMALLINT REFERENCES stage_result_states(id),
  stage_result_time         FLOAT,
  stage_result_time_penalty FLOAT,
  stage_shakedown           BOOLEAN,
  vehicle_class_id          USMALLINT,
  vehicle_id                USMALLINT,
  vehicle_manufacturer_id   USMALLINT,
);

COMMENT ON COLUMN sessions.game_mode IS 'Game mode unique identifier. See "game_mode" table.';
COMMENT ON COLUMN sessions.location_id IS 'Location unique identifier. See "locations" table.';
COMMENT ON COLUMN sessions.route_id IS 'Route unique identifier. See "routes" table.';
COMMENT ON COLUMN sessions.stage_length IS 'Total length of current stage. [metre]';
COMMENT ON COLUMN sessions.stage_result_status IS 'Unique identifier for stage result status (not finished, finished, disqualified etc.). See "stage_result_status table"';
COMMENT ON COLUMN sessions.stage_result_time IS 'Result time. Has a small delay to update after crossing finish line, but guaranteed correct by the time the race telemetry session has ended. Does not take into account retirement or disqualification. [second]';
COMMENT ON COLUMN sessions.stage_result_time_penalty IS 'Total time penalty gained [second]';
COMMENT ON COLUMN sessions.stage_shakedown IS 'Is in a shakedown event.';
COMMENT ON COLUMN sessions.vehicle_class_id IS 'Vehicle class unique identifier. See "vehicle_classes" table.';
COMMENT ON COLUMN sessions.vehicle_id IS 'Vehicle unique identifier. See "vehicles" table.';
COMMENT ON COLUMN sessions.vehicle_manufacturer_id IS 'Vehicle manufacturer unique identifier. See "vehicle_manufacturers" table.';


CREATE TABLE IF NOT EXISTS telemetry (
  timestamp                    TIMESTAMP,
  stage_current_distance       DOUBLE,
  stage_current_time           FLOAT,
  stage_previous_split_time    FLOAT,
  stage_progress               FLOAT,
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
);

COMMENT ON COLUMN telemetry.stage_current_distance IS 'Distance reached on current stage. [metre]';
COMMENT ON COLUMN telemetry.stage_current_time IS 'Time spent on current stage. [second]';
COMMENT ON COLUMN telemetry.stage_previous_split_time IS 'Split time of previous sector. Value unspecified if no sector completed. [second]';
COMMENT ON COLUMN telemetry.stage_progress IS 'Percentage of stage progress 0 to 1 (during race). Value is unspecified before start line and after finish line.';
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
