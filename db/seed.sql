-- CarCare seed data for testing
-- Применяется при первом запуске контейнера (docker-entrypoint-initdb.d)
-- Повторное применение безопасно благодаря ON CONFLICT DO NOTHING
-- Для пересоздания: docker-compose down -v && docker-compose up --build

-- ============================================================
-- ПОЛЬЗОВАТЕЛЬ
-- email: test@mail.ru
-- password: test@carcare  (bcrypt cost=10)
-- ============================================================
INSERT INTO users (id, email, name, role, password_hash)
VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'test@mail.ru',
    'Тест Пользователь',
    'user',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'
)
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- АВТОМОБИЛИ (3 шт.)
-- ============================================================
INSERT INTO cars (id, user_id, brand, model, year, vin, plate)
VALUES
    (
        'b0000000-0000-0000-0000-000000000001',
        'a0000000-0000-0000-0000-000000000001',
        'Toyota',
        'Camry',
        2020,
        'XTA21230002345678',
        'А123БВ77'
    ),
    (
        'b0000000-0000-0000-0000-000000000002',
        'a0000000-0000-0000-0000-000000000001',
        'Volkswagen',
        'Tiguan',
        2019,
        'WVG1234567P123456',
        'В456ГД99'
    ),
    (
        'b0000000-0000-0000-0000-000000000003',
        'a0000000-0000-0000-0000-000000000001',
        'Lada',
        'Vesta',
        2022,
        'XTA21230001111111',
        'Е789ЖЗ50'
    )
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- ЗАПРАВКИ (100 на каждую машину = 300 итого)
-- Объём: 30–70 л, цена: 55–75 руб/л, тип: petrol/diesel
-- Даты: случайные за 2024-01-01 — 2026-03-28 (821 день)
-- ============================================================

-- Toyota Camry (petrol)
INSERT INTO fuel_events (id, car_id, volume, price, type, date)
SELECT
    gen_random_uuid(),
    'b0000000-0000-0000-0000-000000000001',
    ROUND((30 + MOD(i * 37 + 13, 41))::numeric, 2),               -- 30..70 л
    ROUND((55 + MOD(i * 17 + 7,  21))::numeric, 2),               -- 55..75 руб/л
    CASE WHEN MOD(i, 10) < 1 THEN 'diesel' ELSE 'petrol' END,
    DATE '2024-01-01' + (MOD(i * 83 + 29, 821)) * INTERVAL '1 day'
FROM generate_series(1, 100) AS s(i)
WHERE NOT EXISTS (
    SELECT 1 FROM fuel_events WHERE car_id = 'b0000000-0000-0000-0000-000000000001' LIMIT 1
);

-- Volkswagen Tiguan (diesel преимущественно)
INSERT INTO fuel_events (id, car_id, volume, price, type, date)
SELECT
    gen_random_uuid(),
    'b0000000-0000-0000-0000-000000000002',
    ROUND((35 + MOD(i * 41 + 11, 36))::numeric, 2),               -- 35..70 л
    ROUND((60 + MOD(i * 23 + 3,  16))::numeric, 2),               -- 60..75 руб/л
    CASE WHEN MOD(i, 10) < 7 THEN 'diesel' ELSE 'petrol' END,
    DATE '2024-01-01' + (MOD(i * 97 + 47, 821)) * INTERVAL '1 day'
FROM generate_series(1, 100) AS s(i)
WHERE NOT EXISTS (
    SELECT 1 FROM fuel_events WHERE car_id = 'b0000000-0000-0000-0000-000000000002' LIMIT 1
);

-- Lada Vesta (petrol)
INSERT INTO fuel_events (id, car_id, volume, price, type, date)
SELECT
    gen_random_uuid(),
    'b0000000-0000-0000-0000-000000000003',
    ROUND((30 + MOD(i * 29 + 17, 31))::numeric, 2),               -- 30..60 л
    ROUND((55 + MOD(i * 13 + 5,  16))::numeric, 2),               -- 55..70 руб/л
    'petrol',
    DATE '2024-01-01' + (MOD(i * 71 + 61, 821)) * INTERVAL '1 day'
FROM generate_series(1, 100) AS s(i)
WHERE NOT EXISTS (
    SELECT 1 FROM fuel_events WHERE car_id = 'b0000000-0000-0000-0000-000000000003' LIMIT 1
);

-- ============================================================
-- ТЕХНИЧЕСКОЕ ОБСЛУЖИВАНИЕ (100 на каждую машину = 300 итого)
-- Стоимость: 2000–25000 руб
-- Даты: случайные за 2024-01-01 — 2026-03-28
-- ============================================================

-- Справочник типов ТО (10 вариантов, выбираем по MOD)
-- 0: Замена масла и фильтра
-- 1: Замена тормозных колодок
-- 2: Замена шин
-- 3: Плановое ТО
-- 4: Замена воздушного фильтра
-- 5: Диагностика подвески
-- 6: Замена ремня ГРМ
-- 7: Промывка форсунок
-- 8: Замена свечей зажигания
-- 9: Кузовной ремонт

-- Toyota Camry
INSERT INTO maintenance_events (id, car_id, type, date, cost, description)
SELECT
    gen_random_uuid(),
    'b0000000-0000-0000-0000-000000000001',
    CASE MOD(i, 10)
        WHEN 0 THEN 'Замена масла и фильтра'
        WHEN 1 THEN 'Замена тормозных колодок'
        WHEN 2 THEN 'Замена шин'
        WHEN 3 THEN 'Плановое ТО'
        WHEN 4 THEN 'Замена воздушного фильтра'
        WHEN 5 THEN 'Диагностика подвески'
        WHEN 6 THEN 'Замена ремня ГРМ'
        WHEN 7 THEN 'Промывка форсунок'
        WHEN 8 THEN 'Замена свечей зажигания'
        ELSE        'Кузовной ремонт'
    END,
    DATE '2024-01-01' + (MOD(i * 79 + 31, 821)) * INTERVAL '1 day',
    ROUND((2000 + MOD(i * 1103 + 57, 23001))::numeric, 2),         -- 2000..25000
    CASE MOD(i, 10)
        WHEN 0 THEN 'Синтетическое масло 5W-40, масляный и воздушный фильтр'
        WHEN 1 THEN 'Передние и задние колодки Brembo'
        WHEN 2 THEN 'Комплект летних шин Michelin 215/55R17'
        WHEN 3 THEN 'Плановое ТО по регламенту 60 000 км'
        WHEN 4 THEN 'Замена салонного и воздушного фильтра'
        WHEN 5 THEN 'Проверка и регулировка развал-схождения'
        WHEN 6 THEN 'Замена ремня ГРМ с роликами и помпой'
        WHEN 7 THEN 'Ультразвуковая промывка форсунок'
        WHEN 8 THEN 'Иридиевые свечи NGK'
        ELSE        'Устранение вмятины на крыле'
    END
FROM generate_series(1, 100) AS s(i)
WHERE NOT EXISTS (
    SELECT 1 FROM maintenance_events WHERE car_id = 'b0000000-0000-0000-0000-000000000001' LIMIT 1
);

-- Volkswagen Tiguan
INSERT INTO maintenance_events (id, car_id, type, date, cost, description)
SELECT
    gen_random_uuid(),
    'b0000000-0000-0000-0000-000000000002',
    CASE MOD(i, 10)
        WHEN 0 THEN 'Замена масла и фильтра'
        WHEN 1 THEN 'Замена тормозных колодок'
        WHEN 2 THEN 'Замена шин'
        WHEN 3 THEN 'Плановое ТО'
        WHEN 4 THEN 'Замена воздушного фильтра'
        WHEN 5 THEN 'Диагностика подвески'
        WHEN 6 THEN 'Замена ремня ГРМ'
        WHEN 7 THEN 'Промывка форсунок'
        WHEN 8 THEN 'Замена свечей зажигания'
        ELSE        'Кузовной ремонт'
    END,
    DATE '2024-01-01' + (MOD(i * 89 + 43, 821)) * INTERVAL '1 day',
    ROUND((2000 + MOD(i * 997 + 83, 23001))::numeric, 2),
    CASE MOD(i, 10)
        WHEN 0 THEN 'Дизельное моторное масло 5W-30 VW 507.00'
        WHEN 1 THEN 'Тормозные колодки TRW, замена дисков'
        WHEN 2 THEN 'Зимние шины Continental 235/50R18'
        WHEN 3 THEN 'Плановое ТО по регламенту 45 000 км'
        WHEN 4 THEN 'Замена топливного и воздушного фильтра'
        WHEN 5 THEN 'Замена передних амортизаторов Bilstein'
        WHEN 6 THEN 'Замена цепи ГРМ с направляющими'
        WHEN 7 THEN 'Чистка форсунок Common Rail'
        WHEN 8 THEN 'Свечи накала Bosch'
        ELSE        'Покраска бампера'
    END
FROM generate_series(1, 100) AS s(i)
WHERE NOT EXISTS (
    SELECT 1 FROM maintenance_events WHERE car_id = 'b0000000-0000-0000-0000-000000000002' LIMIT 1
);

-- Lada Vesta
INSERT INTO maintenance_events (id, car_id, type, date, cost, description)
SELECT
    gen_random_uuid(),
    'b0000000-0000-0000-0000-000000000003',
    CASE MOD(i, 10)
        WHEN 0 THEN 'Замена масла и фильтра'
        WHEN 1 THEN 'Замена тормозных колодок'
        WHEN 2 THEN 'Замена шин'
        WHEN 3 THEN 'Плановое ТО'
        WHEN 4 THEN 'Замена воздушного фильтра'
        WHEN 5 THEN 'Диагностика подвески'
        WHEN 6 THEN 'Замена ремня ГРМ'
        WHEN 7 THEN 'Промывка форсунок'
        WHEN 8 THEN 'Замена свечей зажигания'
        ELSE        'Кузовной ремонт'
    END,
    DATE '2024-01-01' + (MOD(i * 67 + 53, 821)) * INTERVAL '1 day',
    ROUND((2000 + MOD(i * 877 + 37, 23001))::numeric, 2),
    CASE MOD(i, 10)
        WHEN 0 THEN 'Полусинтетика 5W-40, фильтры оригинал LADA'
        WHEN 1 THEN 'Передние колодки LADA, регулировка ручника'
        WHEN 2 THEN 'Всесезонная резина Nokian 195/55R15'
        WHEN 3 THEN 'Плановое ТО 30 000 км, регулировка клапанов'
        WHEN 4 THEN 'Фильтр воздушный, фильтр салона'
        WHEN 5 THEN 'Замена передних шаровых опор и рулевых наконечников'
        WHEN 6 THEN 'Замена ремня ГРМ и помпы'
        WHEN 7 THEN 'Промывка инжектора, чистка дроссельной заслонки'
        WHEN 8 THEN 'Свечи зажигания BRISK'
        ELSE        'Устранение сколов и царапин на капоте'
    END
FROM generate_series(1, 100) AS s(i)
WHERE NOT EXISTS (
    SELECT 1 FROM maintenance_events WHERE car_id = 'b0000000-0000-0000-0000-000000000003' LIMIT 1
);

-- ============================================================
-- ШТРАФЫ (100 на каждую машину = 300 итого)
-- Суммы: 500–5000 руб, статусы: paid/unpaid
-- Даты: случайные за 2024-01-01 — 2026-03-28
-- ============================================================

-- Справочник типов нарушений (10 вариантов)
-- 0: Превышение скорости 20-40 км/ч
-- 1: Превышение скорости 40-60 км/ч
-- 2: Проезд на красный свет
-- 3: Нарушение правил парковки
-- 4: Непристёгнутый ремень безопасности
-- 5: Разговор по телефону за рулём
-- 6: Неправильный обгон
-- 7: Выезд на встречную полосу
-- 8: Нарушение разметки
-- 9: Несоблюдение дистанции

-- Toyota Camry
INSERT INTO fines (id, car_id, amount, type, date, status, description)
SELECT
    gen_random_uuid(),
    'b0000000-0000-0000-0000-000000000001',
    ROUND((500 + MOD(i * 457 + 23, 4501))::numeric, 2),            -- 500..5000
    CASE MOD(i, 10)
        WHEN 0 THEN 'Превышение скорости 20-40 км/ч'
        WHEN 1 THEN 'Превышение скорости 40-60 км/ч'
        WHEN 2 THEN 'Проезд на красный свет'
        WHEN 3 THEN 'Нарушение правил парковки'
        WHEN 4 THEN 'Непристёгнутый ремень безопасности'
        WHEN 5 THEN 'Разговор по телефону за рулём'
        WHEN 6 THEN 'Неправильный обгон'
        WHEN 7 THEN 'Выезд на встречную полосу'
        WHEN 8 THEN 'Нарушение разметки'
        ELSE        'Несоблюдение дистанции'
    END,
    DATE '2024-01-01' + (MOD(i * 113 + 19, 821)) * INTERVAL '1 day',
    CASE WHEN MOD(i * 7 + 3, 10) < 6 THEN 'paid' ELSE 'unpaid' END,  -- ~60% оплачено
    CASE MOD(i, 10)
        WHEN 0 THEN 'Зафиксировано камерой автофиксации на Ленинском пр-те'
        WHEN 1 THEN 'Зафиксировано камерой ГИБДД на трассе М4'
        WHEN 2 THEN 'Нарушение ПДД на регулируемом перекрёстке'
        WHEN 3 THEN 'Автомобиль припаркован под знаком «Стоянка запрещена»'
        WHEN 4 THEN 'Инспектор ГИБДД, водитель не пристёгнут'
        WHEN 5 THEN 'Использование мобильного телефона без гарнитуры'
        WHEN 6 THEN 'Обгон в зоне сплошной разметки'
        WHEN 7 THEN 'Выезд на встречку на загородной трассе'
        WHEN 8 THEN 'Пересечение сплошной разделительной полосы'
        ELSE        'Резкое торможение, несоблюдение дистанции'
    END
FROM generate_series(1, 100) AS s(i)
WHERE NOT EXISTS (
    SELECT 1 FROM fines WHERE car_id = 'b0000000-0000-0000-0000-000000000001' LIMIT 1
);

-- Volkswagen Tiguan
INSERT INTO fines (id, car_id, amount, type, date, status, description)
SELECT
    gen_random_uuid(),
    'b0000000-0000-0000-0000-000000000002',
    ROUND((500 + MOD(i * 541 + 37, 4501))::numeric, 2),
    CASE MOD(i, 10)
        WHEN 0 THEN 'Превышение скорости 20-40 км/ч'
        WHEN 1 THEN 'Превышение скорости 40-60 км/ч'
        WHEN 2 THEN 'Проезд на красный свет'
        WHEN 3 THEN 'Нарушение правил парковки'
        WHEN 4 THEN 'Непристёгнутый ремень безопасности'
        WHEN 5 THEN 'Разговор по телефону за рулём'
        WHEN 6 THEN 'Неправильный обгон'
        WHEN 7 THEN 'Выезд на встречную полосу'
        WHEN 8 THEN 'Нарушение разметки'
        ELSE        'Несоблюдение дистанции'
    END,
    DATE '2024-01-01' + (MOD(i * 131 + 59, 821)) * INTERVAL '1 day',
    CASE WHEN MOD(i * 11 + 7, 10) < 5 THEN 'paid' ELSE 'unpaid' END,  -- ~50% оплачено
    CASE MOD(i, 10)
        WHEN 0 THEN 'Камера фиксации на Московском шоссе'
        WHEN 1 THEN 'Радар ГИБДД на трассе М11'
        WHEN 2 THEN 'Перекрёсток Садового кольца'
        WHEN 3 THEN 'Парковка на тротуаре в центре города'
        WHEN 4 THEN 'Остановка ГИБДД на патрульном посту'
        WHEN 5 THEN 'Использование смартфона во время движения'
        WHEN 6 THEN 'Обгон под знаком «Обгон запрещён»'
        WHEN 7 THEN 'Выезд на полосу встречного движения'
        WHEN 8 THEN 'Перестроение через двойную сплошную'
        ELSE        'ДТП по причине несоблюдения дистанции'
    END
FROM generate_series(1, 100) AS s(i)
WHERE NOT EXISTS (
    SELECT 1 FROM fines WHERE car_id = 'b0000000-0000-0000-0000-000000000002' LIMIT 1
);

-- Lada Vesta
INSERT INTO fines (id, car_id, amount, type, date, status, description)
SELECT
    gen_random_uuid(),
    'b0000000-0000-0000-0000-000000000003',
    ROUND((500 + MOD(i * 389 + 41, 4501))::numeric, 2),
    CASE MOD(i, 10)
        WHEN 0 THEN 'Превышение скорости 20-40 км/ч'
        WHEN 1 THEN 'Превышение скорости 40-60 км/ч'
        WHEN 2 THEN 'Проезд на красный свет'
        WHEN 3 THEN 'Нарушение правил парковки'
        WHEN 4 THEN 'Непристёгнутый ремень безопасности'
        WHEN 5 THEN 'Разговор по телефону за рулём'
        WHEN 6 THEN 'Неправильный обгон'
        WHEN 7 THEN 'Выезд на встречную полосу'
        WHEN 8 THEN 'Нарушение разметки'
        ELSE        'Несоблюдение дистанции'
    END,
    DATE '2024-01-01' + (MOD(i * 107 + 73, 821)) * INTERVAL '1 day',
    CASE WHEN MOD(i * 3 + 1, 10) < 7 THEN 'paid' ELSE 'unpaid' END,  -- ~70% оплачено
    CASE MOD(i, 10)
        WHEN 0 THEN 'Камера видеофиксации на ул. Гагарина'
        WHEN 1 THEN 'Измерение скорости патрульным автомобилем'
        WHEN 2 THEN 'Инспектор ГИБДД, нарушение на светофоре'
        WHEN 3 THEN 'Парковка в зоне действия знака 3.27'
        WHEN 4 THEN 'Проверка на контрольном посту'
        WHEN 5 THEN 'Телефон в руке на перекрёстке'
        WHEN 6 THEN 'Обгон в зоне ограниченной видимости'
        WHEN 7 THEN 'Выезд на встречную полосу через прерывистую, ставшую сплошной'
        WHEN 8 THEN 'Нарушение горизонтальной разметки'
        ELSE        'Резкое торможение без причины'
    END
FROM generate_series(1, 100) AS s(i)
WHERE NOT EXISTS (
    SELECT 1 FROM fines WHERE car_id = 'b0000000-0000-0000-0000-000000000003' LIMIT 1
);
