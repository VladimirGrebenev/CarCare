# Интеграционные адаптеры (Hexagonal Architecture)

## Госуслуги
- **Файл:** internal/adapter/gosuslugi/gosuslugi_adapter.go
- **Тесты:** internal/adapter/gosuslugi/gosuslugi_adapter_test.go
- **Описание:** Проверка штрафов по номеру автомобиля через внешний сервис Госуслуги.

## Карты
- **Файл:** internal/adapter/maps/maps_adapter.go
- **Тесты:** internal/adapter/maps/maps_adapter_test.go
- **Описание:** Получение маршрута между точками через внешний картографический сервис.

## Оплата
- **Файл:** internal/adapter/payment/payment_adapter.go
- **Тесты:** internal/adapter/payment/payment_adapter_test.go
- **Описание:** Проведение оплаты штрафа через внешний платёжный шлюз.

## Мокирование
- Для тестирования интеграций используйте стандартные Go-моки (интерфейсы, подмены).
- Примеры тестов см. в соответствующих *_test.go файлах.
