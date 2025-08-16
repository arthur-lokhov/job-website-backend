-- Включаем расширение для генерации UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Очистка старых данных для идемпотентности скрипта
TRUNCATE TABLE applications, vacancies, locations, levels, departments RESTART IDENTITY CASCADE;

-- Вставка тестовых данных
DO $$
DECLARE
    dept_id_0 UUID := 'a1b2c3d4-e5f6-7890-1234-567890abcdef';
    dept_id_1 UUID := 'b2c3d4e5-f6a7-8901-2345-67890abcdef0';
    dept_id_2 UUID := 'c3d4e5f6-a7b8-9012-3456-7890abcdef01';
    dept_id_3 UUID := 'd4e5f6a7-b8c9-0123-4567-890abcdef012';
    dept_id_4 UUID := 'e5f6a7b8-c9d0-1234-5678-90abcdef0123';
    dept_id_5 UUID := 'f6a7b8c9-d0e1-2345-6789-0abcdef01234';
    dept_id_6 UUID := 'a7b8c9d0-e1f2-3456-7890-abcdef012345';
    dept_id_7 UUID := 'b8c9d0e1-f2a3-4567-8901-bcdef0123456';
    dept_id_8 UUID := 'c9d0e1f2-a3b4-5678-9012-cdef01234567';
    dept_id_9 UUID := 'd0e1f2a3-b4c5-6789-0123-def012345678';

    level_id_0 UUID := '11111111-1111-1111-1111-111111111111';
    level_id_1 UUID := '22222222-2222-2222-2222-222222222222';
    level_id_2 UUID := '33333333-3333-3333-3333-333333333333';
    level_id_3 UUID := '44444444-4444-4444-4444-444444444444';
    level_id_4 UUID := '55555555-5555-5555-5555-555555555555';

    loc_id_0 UUID := 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';
    loc_id_1 UUID := 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb';
    loc_id_2 UUID := 'cccccccc-cccc-cccc-cccc-cccccccccccc';
    loc_id_3 UUID := 'dddddddd-dddd-dddd-dddd-dddddddddddd';
BEGIN
    -- Departments
    INSERT INTO departments (id) VALUES
        (dept_id_0), (dept_id_1), (dept_id_2), (dept_id_3), (dept_id_4),
        (dept_id_5), (dept_id_6), (dept_id_7), (dept_id_8), (dept_id_9);

    -- Levels
    INSERT INTO levels (id, name, priority) VALUES
        (level_id_0, 'Intern', 0),
        (level_id_1, 'Junior', 1),
        (level_id_2, 'Middle', 2),
        (level_id_3, 'Senior', 3),
        (level_id_4, 'Lead', 4);

    -- Locations
    INSERT INTO locations (id, name, priority) VALUES
        (loc_id_0, 'Удаленно', 0),
        (loc_id_1, 'Офис Москва', 1),
        (loc_id_2, 'Офис Санкт-Петербург', 2),
        (loc_id_3, 'Гибрид', 3);

    -- Vacancies
    INSERT INTO vacancies (id, name, department_id, level_id, location_id, info, application_form, is_active, important, priority) VALUES
        (uuid_generate_v4(), 'Backend Developer (Go)', dept_id_0, level_id_3, loc_id_0, 'Мы ищем крутого специалиста на позицию Backend Developer (Go).', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 1),
        (uuid_generate_v4(), 'Frontend Developer (React)', dept_id_1, level_id_2, loc_id_1, 'Нужен Middle React разработчик в команду продукта.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "portfolioLink": {"type": "string"}}}', true, false, 2),
        (uuid_generate_v4(), 'DevOps Engineer (Kubernetes)', dept_id_2, level_id_3, loc_id_0, 'Ищем Senior DevOps для поддержки и развития инфраструктуры.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 3),
        (uuid_generate_v4(), 'QA Engineer (Automation)', dept_id_3, level_id_2, loc_id_2, 'Автоматизатор тестирования со знанием Python и Selenium.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 4),
        (uuid_generate_v4(), 'Product Manager', dept_id_4, level_id_4, loc_id_1, 'Лидер продуктового направления в финтех стартап.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 5),
        (uuid_generate_v4(), 'Data Scientist (Python)', dept_id_5, level_id_3, loc_id_0, 'Исследователь данных для построения моделей машинного обучения.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "githubLink": {"type": "string"}}}', false, false, 6),
        (uuid_generate_v4(), 'UI/UX Designer', dept_id_1, level_id_2, loc_id_3, 'Дизайнер интерфейсов для мобильного приложения.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "portfolioLink": {"type": "string"}}}', true, false, 7),
        (uuid_generate_v4(), 'Fullstack Developer (Node.js/React)', dept_id_6, level_id_2, loc_id_1, 'Разработчик полного цикла для работы над CRM системой.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 8),
        (uuid_generate_v4(), 'Mobile Developer (iOS)', dept_id_7, level_id_3, loc_id_0, 'Senior iOS разработчик для создания нового банковского приложения.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 9),
        (uuid_generate_v4(), 'Junior Java Developer', dept_id_8, level_id_1, loc_id_2, 'Начинающий Java-разработчик в enterprise проект.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 10),
        (uuid_generate_v4(), 'Project Manager', dept_id_9, level_id_3, loc_id_1, 'Опытный менеджер проектов для ведения команды из 10 человек.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 11),
        (uuid_generate_v4(), 'Data Analyst', dept_id_5, level_id_2, loc_id_3, 'Аналитик данных для работы с большими данными.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', false, false, 12),
        (uuid_generate_v4(), 'Lead Frontend Developer (Vue.js)', dept_id_1, level_id_4, loc_id_0, 'Ведущий frontend разработчик для создания UI кита.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "portfolioLink": {"type": "string"}}}', true, true, 13),
        (uuid_generate_v4(), 'System Administrator (Linux)', dept_id_2, level_id_2, loc_id_2, 'Системный администратор для поддержки наших серверов.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 14),
        (uuid_generate_v4(), 'Intern Python Developer', dept_id_0, level_id_0, loc_id_1, 'Стажер в команду бэкенда, готовы обучать.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 15),
        (uuid_generate_v4(), 'Android Developer (Kotlin)', dept_id_7, level_id_2, loc_id_3, 'Разработчик под Android для развития существующего приложения.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 16),
        (uuid_generate_v4(), 'Senior C++ Developer', dept_id_8, level_id_3, loc_id_0, 'С++ разработчик для работы над высоконагруженными системами.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 17),
        (uuid_generate_v4(), 'Marketing Manager', dept_id_4, level_id_2, loc_id_1, 'Менеджер по маркетингу для продвижения нового продукта.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 18),
        (uuid_generate_v4(), 'Junior QA (Manual)', dept_id_3, level_id_1, loc_id_2, 'Младший специалист по ручному тестированию.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 19),
        (uuid_generate_v4(), 'Senior Frontend Developer (Angular)', dept_id_1, level_id_3, loc_id_0, 'Опытный Angular разработчик для работы над SPA.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "portfolioLink": {"type": "string"}}}', false, true, 20),
        (uuid_generate_v4(), 'Backend Developer (Python)', dept_id_0, level_id_2, loc_id_1, 'Python/Django разработчик в команду разработки.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 21),
        (uuid_generate_v4(), 'Lead DevOps', dept_id_2, level_id_4, loc_id_0, 'Лид DevOps команды.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 22),
        (uuid_generate_v4(), 'HR Manager', dept_id_9, level_id_2, loc_id_1, 'Менеджер по персоналу.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 23),
        (uuid_generate_v4(), 'BI Analyst', dept_id_5, level_id_3, loc_id_3, 'BI аналитик для построения отчетности.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 24),
        (uuid_generate_v4(), 'Junior Golang Developer', dept_id_0, level_id_1, loc_id_2, 'Ищем начинающего Go разработчика.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 25),
        (uuid_generate_v4(), 'Technical Writer', dept_id_3, level_id_2, loc_id_0, 'Технический писатель для документирования нашего API.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 26),
        (uuid_generate_v4(), 'Machine Learning Engineer', dept_id_5, level_id_3, loc_id_0, 'Инженер машинного обучения.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "githubLink": {"type": "string"}}}', true, true, 27),
        (uuid_generate_v4(), 'React Native Developer', dept_id_7, level_id_2, loc_id_3, 'Разработчик мобильных приложений на React Native.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 28),
        (uuid_generate_v4(), 'Senior Java Developer', dept_id_8, level_id_3, loc_id_1, 'Ведущий Java-разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 29),
        (uuid_generate_v4(), 'Intern Frontend Developer', dept_id_1, level_id_0, loc_id_1, 'Стажер в команду frontend разработки.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 30),
        (uuid_generate_v4(), 'DBA (PostgreSQL)', dept_id_2, level_id_3, loc_id_0, 'Администратор баз данных PostgreSQL.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 31),
        (uuid_generate_v4(), 'Scrum Master', dept_id_9, level_id_3, loc_id_1, 'Scrum мастер для нескольких команд.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 32),
        (uuid_generate_v4(), 'Security Engineer', dept_id_2, level_id_3, loc_id_0, 'Инженер по информационной безопасности.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 33),
        (uuid_generate_v4(), 'Python Team Lead', dept_id_0, level_id_4, loc_id_1, 'Тимлид команды Python разработки.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 34),
        (uuid_generate_v4(), 'Vue.js Developer', dept_id_1, level_id_2, loc_id_2, 'Разработчик на Vue.js.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "portfolioLink": {"type": "string"}}}', true, false, 35),
        (uuid_generate_v4(), 'Senior QA Automation (Java)', dept_id_3, level_id_3, loc_id_0, 'Ведущий инженер по автоматизации тестирования на Java.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 36),
        (uuid_generate_v4(), 'Business Analyst', dept_id_4, level_id_2, loc_id_1, 'Бизнес-аналитик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 37),
        (uuid_generate_v4(), 'Senior Backend Developer (Node.js)', dept_id_6, level_id_3, loc_id_0, 'Senior Node.js разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 38),
        (uuid_generate_v4(), 'Swift Developer', dept_id_7, level_id_2, loc_id_3, 'Разработчик на Swift.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 39),
        (uuid_generate_v4(), 'Junior C# Developer', dept_id_8, level_id_1, loc_id_1, 'Младший .NET разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 40),
        (uuid_generate_v4(), 'Recruiter', dept_id_9, level_id_2, loc_id_1, 'Специалист по подбору персонала.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 41),
        (uuid_generate_v4(), 'Data Engineer', dept_id_5, level_id_3, loc_id_0, 'Инженер данных.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 42),
        (uuid_generate_v4(), 'Lead UI/UX Designer', dept_id_1, level_id_4, loc_id_1, 'Ведущий дизайнер интерфейсов.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "portfolioLink": {"type": "string"}}}', true, true, 43),
        (uuid_generate_v4(), 'Linux Administrator', dept_id_2, level_id_3, loc_id_2, 'Администратор Linux.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 44),
        (uuid_generate_v4(), 'Intern QA', dept_id_3, level_id_0, loc_id_1, 'Стажер в отдел тестирования.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 45),
        (uuid_generate_v4(), 'Flutter Developer', dept_id_7, level_id_2, loc_id_0, 'Разработчик на Flutter.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 46),
        (uuid_generate_v4(), 'Senior Scala Developer', dept_id_8, level_id_3, loc_id_0, 'Senior Scala разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 47),
        (uuid_generate_v4(), 'Office Manager', dept_id_9, level_id_1, loc_id_1, 'Офис-менеджер.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 48),
        (uuid_generate_v4(), 'Content Manager', dept_id_4, level_id_2, loc_id_3, 'Контент-менеджер.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 49),
        (uuid_generate_v4(), 'Middle Go Developer', dept_id_0, level_id_2, loc_id_2, 'Go разработчик уровня middle.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 50),
        (uuid_generate_v4(), 'JavaScript Developer', dept_id_1, level_id_2, loc_id_0, 'JavaScript разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "portfolioLink": {"type": "string"}}}', true, false, 51),
        (uuid_generate_v4(), 'Cloud Engineer (AWS)', dept_id_2, level_id_3, loc_id_0, 'Облачный инженер AWS.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 52),
        (uuid_generate_v4(), 'Game Tester', dept_id_3, level_id_1, loc_id_3, 'Тестировщик игр.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 53),
        (uuid_generate_v4(), 'Head of Product', dept_id_4, level_id_4, loc_id_1, 'Руководитель продукта.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 54),
        (uuid_generate_v4(), 'Senior Data Scientist', dept_id_5, level_id_3, loc_id_0, 'Ведущий специалист по данным.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "githubLink": {"type": "string"}}}', true, true, 55),
        (uuid_generate_v4(), 'Graphic Designer', dept_id_1, level_id_2, loc_id_1, 'Графический дизайнер.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "portfolioLink": {"type": "string"}}}', true, false, 56),
        (uuid_generate_v4(), 'PHP Developer', dept_id_6, level_id_2, loc_id_2, 'PHP разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 57),
        (uuid_generate_v4(), 'Senior iOS Developer', dept_id_7, level_id_3, loc_id_0, 'Ведущий iOS разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 58),
        (uuid_generate_v4(), 'Ruby on Rails Developer', dept_id_8, level_id_2, loc_id_0, 'Разработчик Ruby on Rails.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 59),
        (uuid_generate_v4(), 'IT Director', dept_id_9, level_id_4, loc_id_1, 'IT-директор.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 60),
        (uuid_generate_v4(), 'Big Data Engineer', dept_id_5, level_id_3, loc_id_0, 'Инженер Big Data.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 61),
        (uuid_generate_v4(), 'Frontend Team Lead', dept_id_1, level_id_4, loc_id_1, 'Тимлид frontend команды.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "portfolioLink": {"type": "string"}}}', true, true, 62),
        (uuid_generate_v4(), 'Network Engineer', dept_id_2, level_id_3, loc_id_2, 'Сетевой инженер.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 63),
        (uuid_generate_v4(), 'QA Team Lead', dept_id_3, level_id_4, loc_id_1, 'Тимлид команды тестирования.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 64),
        (uuid_generate_v4(), 'Product Analyst', dept_id_4, level_id_2, loc_id_3, 'Продуктовый аналитик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 65),
        (uuid_generate_v4(), 'Senior Python Developer', dept_id_0, level_id_3, loc_id_0, 'Ведущий Python разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 66),
        (uuid_generate_v4(), 'Angular Developer', dept_id_1, level_id_2, loc_id_1, 'Angular разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "portfolioLink": {"type": "string"}}}', true, false, 67),
        (uuid_generate_v4(), 'Cloud Architect (Azure)', dept_id_2, level_id_4, loc_id_0, 'Облачный архитектор Azure.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 68),
        (uuid_generate_v4(), 'Manual QA Engineer', dept_id_3, level_id_2, loc_id_2, 'Инженер по ручному тестированию.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 69),
        (uuid_generate_v4(), 'Marketing Director', dept_id_4, level_id_4, loc_id_1, 'Директор по маркетингу.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 70),
        (uuid_generate_v4(), 'ML Team Lead', dept_id_5, level_id_4, loc_id_0, 'Тимлид команды машинного обучения.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "githubLink": {"type": "string"}}}', true, true, 71),
        (uuid_generate_v4(), 'Motion Designer', dept_id_1, level_id_2, loc_id_3, 'Моушн-дизайнер.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "portfolioLink": {"type": "string"}}}', true, false, 72),
        (uuid_generate_v4(), 'Senior PHP Developer', dept_id_6, level_id_3, loc_id_0, 'Ведущий PHP разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 73),
        (uuid_generate_v4(), 'Senior Android Developer', dept_id_7, level_id_3, loc_id_1, 'Ведущий Android разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 74),
        (uuid_generate_v4(), 'Elixir Developer', dept_id_8, level_id_2, loc_id_0, 'Elixir разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 75),
        (uuid_generate_v4(), 'Chief Technology Officer (CTO)', dept_id_9, level_id_4, loc_id_1, 'Технический директор.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 76),
        (uuid_generate_v4(), 'Data Warehouse Engineer', dept_id_5, level_id_3, loc_id_2, 'Инженер DWH.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 77),
        (uuid_generate_v4(), 'UI Designer', dept_id_1, level_id_2, loc_id_1, 'UI дизайнер.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "portfolioLink": {"type": "string"}}}', true, false, 78),
        (uuid_generate_v4(), 'Site Reliability Engineer (SRE)', dept_id_2, level_id_3, loc_id_0, 'SRE инженер.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 79),
        (uuid_generate_v4(), 'Performance QA Engineer', dept_id_3, level_id_3, loc_id_0, 'Инженер по нагрузочному тестированию.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 80),
        (uuid_generate_v4(), 'Agile Coach', dept_id_9, level_id_3, loc_id_1, 'Agile коуч.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 81),
        (uuid_generate_v4(), 'Lead Golang Developer', dept_id_0, level_id_4, loc_id_0, 'Ведущий Go разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 82),
        (uuid_generate_v4(), 'TypeScript Developer', dept_id_1, level_id_2, loc_id_3, 'TypeScript разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "portfolioLink": {"type": "string"}}}', true, false, 83),
        (uuid_generate_v4(), 'DevOps (GCP)', dept_id_2, level_id_3, loc_id_0, 'DevOps инженер в GCP.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 84),
        (uuid_generate_v4(), 'Mobile QA Engineer', dept_id_3, level_id_2, loc_id_2, 'Тестировщик мобильных приложений.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 85),
        (uuid_generate_v4(), 'SEO Specialist', dept_id_4, level_id_2, loc_id_1, 'SEO специалист.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 86),
        (uuid_generate_v4(), 'Computer Vision Engineer', dept_id_5, level_id_3, loc_id_0, 'Инженер по компьютерному зрению.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "githubLink": {"type": "string"}}}', true, true, 87),
        (uuid_generate_v4(), 'UX Researcher', dept_id_1, level_id_2, loc_id_3, 'UX исследователь.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "portfolioLink": {"type": "string"}}}', true, false, 88),
        (uuid_generate_v4(), 'Symfony Developer', dept_id_6, level_id_2, loc_id_0, 'Разработчик на Symfony.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 89),
        (uuid_generate_v4(), 'Cross-platform Developer (Xamarin)', dept_id_7, level_id_2, loc_id_1, 'Кросс-платформенный разработчик Xamarin.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 90),
        (uuid_generate_v4(), 'Rust Developer', dept_id_8, level_id_2, loc_id_0, 'Rust разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 91),
        (uuid_generate_v4(), 'Technical Support Engineer', dept_id_9, level_id_2, loc_id_2, 'Инженер технической поддержки.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 92),
        (uuid_generate_v4(), 'NLP Engineer', dept_id_5, level_id_3, loc_id_0, 'NLP инженер.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "githubLink": {"type": "string"}}}', true, true, 93),
        (uuid_generate_v4(), 'Illustrator', dept_id_1, level_id_2, loc_id_3, 'Иллюстратор.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}, "portfolioLink": {"type": "string"}}}', true, false, 94),
        (uuid_generate_v4(), 'Information Security Analyst', dept_id_2, level_id_2, loc_id_1, 'Аналитик по информационной безопасности.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 95),
        (uuid_generate_v4(), 'Automation Engineer', dept_id_3, level_id_2, loc_id_0, 'Инженер по автоматизации.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 96),
        (uuid_generate_v4(), 'PR Manager', dept_id_4, level_id_2, loc_id_1, 'PR менеджер.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 97),
        (uuid_generate_v4(), 'Backend Team Lead (Java)', dept_id_8, level_id_4, loc_id_1, 'Тимлид backend команды (Java).', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 98),
        (uuid_generate_v4(), 'Laravel Developer', dept_id_6, level_id_2, loc_id_2, 'Laravel разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 99),
        (uuid_generate_v4(), 'Unity Developer', dept_id_7, level_id_2, loc_id_0, 'Unity разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, false, 100),
        (uuid_generate_v4(), 'Haskell Developer', dept_id_8, level_id_3, loc_id_0, 'Haskell разработчик.', '{"formSchema": {"fullName": {"type": "string"}, "email": {"type": "string"}}}', true, true, 101);
END $$;