# Avito-app
# Тестовое задание для стажёра Backend
# Сервис динамического сегментирования пользователей
### Для запуска приложения:

```
make build && make run
```
## Requests (Linux)
#### Privilege creation:

`curl -XPOST http://localhost:8000/priv -d '{"privilege_title": "AVITO_VOICE_MESSAGES"}'`

#### Attach privileges to specific user:

`curl -XPOST http://localhost:8000/priv/user/add -d '{"user_id":1, "add_privilege": ["AVITO_VOICE_MESSAGES", "AVITO_PERFORMANCE_VAS", "AVITO_DISCOUNT_30"]}'`

#### Remove privileges from specific user:

`curl -XPOST http://localhost:8000/priv/user/remove -d '{"user_id":1, "add_privilege": ["AVITO_VOICE_MESSAGES", "AVITO_DISCOUNT_30"]}'`

#### Get users' privileges:

`curl -XGET http://localhost:8000/priv/user`

#### Get existing privileges:

`curl -XGET http://localhost:8000/priv`

# Тестовое задание для стажёра Backend
# Сервис динамического сегментирования пользователей

### Проблема:

В Авито часто проводятся различные эксперименты — тесты новых продуктов, тесты интерфейса, скидочные и многие другие.
На архитектурном комитете приняли решение централизовать работу с проводимыми экспериментами и вынести этот функционал в отдельный сервис.

### Задача:

Требуется реализовать сервис, хранящий пользователя и сегменты, в которых он состоит (создание, изменение, удаление сегментов, а также добавление и удаление пользователей в сегмент)

**Сценарии использования:**
Хотим провести несколько экспериментов и протестировать новый функционал Авито:
- Голосовые сообщения в чатах
- Новые услуги продвижения
- Скидка 30% на услуги продвижения
- Скидка 50% на услуги продвижения

> Кто из пользователей в какой эксперимент попадет будет решать большой отдел аналитики, а мы лишь дадим им возможность для таких тестов.

Допустим аналитики создали сегменты:
- AVITO_VOICE_MESSAGES
- AVITO_PERFORMANCE_VAS
- AVITO_DISCOUNT_30
- AVITO_DISCOUNT_50

и добавили созданные сегменты нескольким пользователям:

| Пользователь | Сегменты которым он принадлежит |
| --- | --- |
| 1000 | [AVITO_VOICE_MESSAGES, AVITO_PERFORMANCE_VAS, AVITO_DISCOUNT_30, …] |
| 1002 | [AVITO_VOICE_MESSAGES, AVITO_DISCOUNT_50, …] |
| 1004 | нет сегментов |

> Формат хранения данных в базе данных не ограничен - можно выбрать любой удобный

**Получили следующие данные:**

- Пользователь 1000 состоит в 3 сегментах: AVITO_VOICE_MESSAGES, AVITO_PERFORMANCE_VAS, AVITO_DISCOUNT_30

- Пользователь 1002 состоит в 2 сегментах: AVITO_VOICE_MESSAGES, AVITO_DISCOUNT_50

- Пользователь 1004 не состоит ни в одном из сегментов 

Теперь мы хотим через API сервиса по user_id получать список сегментов в которых он состоит.

## Требования и детали по заданию

**Технические требования:**

1. Сервис должен предоставлять HTTP API с форматом JSON как при отправке запроса, так и при получении результата.
2. Язык разработки: Golang.
3. Фреймворки и библиотеки можно использовать любые.
4. Реляционная СУБД: MySQL или PostgreSQL.
5. Использование docker и docker-compose для поднятия и развертывания dev-среды.
6. Весь код должен быть выложен на Github/Gitlab с Readme файлом с инструкцией по запуску и примерами запросов/ответов (можно просто описать в Readme методы, можно через Postman, можно в Readme curl запросы скопировать, и так далее).
7. Если есть потребность в асинхронных сценариях, то использование любых систем очередей - допускается.
8. При возникновении вопросов по ТЗ оставляем принятие решения за кандидатом (в таком случае в Readme файле к проекту должен быть указан список вопросов, с которыми кандидат столкнулся и каким образом он их решил).
9. Разработка интерфейса в браузере НЕ ТРЕБУЕТСЯ. Взаимодействие с API предполагается посредством запросов из кода другого сервиса. Для тестирования можно использовать любой удобный инструмент. Например: в терминале через curl или Postman.

**Будет плюсом:**
1. Покрытие кода тестами.
2. [Swagger](https://swagger.io/solutions/api-design/) файл для вашего API.

**Основное задание (минимум):**

1. Метод создания сегмента. Принимает slug (название) сегмента. 
2. Метод удаления сегмента. Принимает slug (название) сегмента. 
3. Метод добавления пользователя в сегмент. Принимает список slug (названий) сегментов которые нужно добавить пользователю, список slug (названий) сегментов которые нужно удалить у пользователя, id пользователя.
4. Метод получения активных сегментов пользователя. Принимает на вход id пользователя.

**Детали по заданию:**

1. По умолчанию сервис не содержит в себе никаких данных о сегментах и пользователях (пустая табличка в БД). Данные появляются при создании сегментов и добавлении их пользователям.
2. Валидацию данных и обработку ошибок оставляем на ваше усмотрение.
3. Список полей к методам не фиксированный. Перечислен лишь необходимый минимум. В рамках выполнения доп. заданий возможны дополнительные поля.
4. Механизм миграции не нужен. Достаточно предоставить конечный SQL файл с созданием всех необходимых таблиц в БД.
5. Сегменты пользователя очень важны - на этих данных в будущем строится аналитика о том насколько продукт востребован. Поэтому нужно следить за тем чтобы сегменты не терялись (не перетирались) и не добавлялись лишним пользователям.
6. Мы можем добавлять сегменты пользователю динамически (он уже состоит в нескольких, мы можем добавить еще несколько, не перетирая существующие).
7. В методе получения сегментов пользователя мы должны получить АКТУАЛЬНУЮ информацию о сегментах пользователя с задержкой не более 1 минуты после добавления сегмента.


