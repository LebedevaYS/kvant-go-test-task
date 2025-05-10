# Go Test Task

REST API для управления пользователями и заказами

## Функционал
- Регистрация и авторизация (JWT)
- CRUD для пользователей
- Создание и просмотр заказов
- Пагинация и фильтрация

## Запуск
```bash
docker-compose up --build
```
## Проверка работы
```bash
curl http://localhost:8080/health
```
## Добавить пользователя:
```bash
Invoke-RestMethod -Uri "http://localhost:8080/users" -Method Post `
  -Headers @{"Content-Type"="application/json"} `
  -Body '{"name":"PowerShell User", "email":"ps@test.com", "password":"password123", "age":28}'
```
## Добавить массив пользователей
```bash
$users = @(
    @{name="Alice Smith"; email="alice@example.com"; age=25; password="password123"},
    @{name="Bob Johnson"; email="bob@example.com"; age=30; password="password123"},
    @{name="Charlie Brown"; email="charlie@example.com"; age=22; password="password123"},
    @{name="Diana Prince"; email="diana@example.com"; age=28; password="password123"},
    @{name="Eve Wilson"; email="eve@example.com"; age=35; password="password123"}
)
foreach ($user in $users) {
    $json = $user | ConvertTo-Json
    $response = Invoke-RestMethod -Uri "http://localhost:8080/users" -Method Post `
        -Headers @{"Content-Type"="application/json"} `
        -Body $json
    Write-Output "Created user: $($response.email)"
}
```
## Авторизоваться и получить токен:
```bash
$login = @{
    email = "bob@example.com"
    password = "password123"
} | ConvertTo-Json
$response = Invoke-RestMethod -Uri "http://localhost:8080/auth/login" -Method Post -Body $login -ContentType "application/json"
$token = $response.token
$headers = @{
    "Authorization" = "Bearer $token"
}
```
## Получить весь список пользователей:
```bash
$users = Invoke-RestMethod -Uri "http://localhost:8080/users" -Method Get -Headers $headers
Write-Output "Список пользователей:"
$users | ConvertTo-Json -Depth 3
```
## Запрос с параметрами:
```bash
$params = @{
    page = 2
    limit = 5
    min_age = 20
    max_age = 30
}
$query = ($params.GetEnumerator() | ForEach-Object { "$($_.Key)=$($_.Value)" }) -join '&'
$response = Invoke-RestMethod -Uri "http://localhost:8080/users?$query" -Method Get -Headers $headers
$response | ConvertTo-Json -Depth 5
```
## Получение пользователя по ID:
```bash
$userId = 3
$user = Invoke-RestMethod -Uri "http://localhost:8080/users/$userId" -Method Get -Headers $headers
Write-Output "Данные пользователя с ID $userId :"
$user | ConvertTo-Json -Depth 3
# Проверка на несуществующего пользователя
Invoke-RestMethod -Uri "http://localhost:8080/users/999999" -Method Get -Headers $headers -StatusCodeVariable statusCode
Write-Output "Статус код для несуществующего пользователя: $statusCode"
```
## Обновление пользователя:
```bash
$userId = 1  
$updateData = @{
    name = "Updated Name"
    email = "updated.email@example.com"
    age = 35
} | ConvertTo-Json
try {
    $updatedUser = Invoke-RestMethod -Uri "http://localhost:8080/users/$userId" -Method Put -Headers $headers -Body $updateData -ContentType "application/json"
    $updatedUser | ConvertTo-Json -Depth 3
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    Write-Output "Ошибка обновления. Статус код: $statusCode"
    Write-Output "Сообщение: $($_.Exception.Message)"
}
```
## Удаление пользователя:
```bash
$userId = 7
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/users/$userId" -Method Delete -Headers $headers
    Write-Output "Пользователь удален (статус 204)"
} catch {
    if ($_.Exception.Response) {
        $statusCode = $_.Exception.Response.StatusCode.value__
        Write-Output "Статус код: $statusCode"
    } else {
        Write-Output "Ошибка: $($_.Exception.Message)"
    }
}
```
## Добавление нового заказа:
```bash
$users = Invoke-RestMethod -Uri "http://localhost:8080/users" -Method Get -Headers $headers
$userId = $users.users[0].id
$orderData = @{
    product = " Lenovo ThinkPad"
    quantity = 1
    price = 1200.50
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/users/$userId/orders" -Method Post `
        -Headers $headers `
        -Body $orderData `
        -ContentType "application/json"
    
    Write-Output "Заказ успешно создан:"
    $response | Format-List
} catch {
    $stream = $_.Exception.Response.GetResponseStream()
    $reader = New-Object System.IO.StreamReader($stream)
    $errorBody = $reader.ReadToEnd()
    $reader.Close()
    
    Write-Output "Ошибка: $($_.Exception.Response.StatusCode.value__)"
    Write-Output $errorBody
}
```
## Получение списка заказов:
```bash
try {
    $orders = Invoke-RestMethod -Uri "http://localhost:8080/users/$userId/orders" -Method Get -Headers $headers
    
    Write-Output "Список заказов пользователя $userId :"
    $orders | Format-Table -AutoSize
} catch {
    $stream = $_.Exception.Response.GetResponseStream()
    $reader = New-Object System.IO.StreamReader($stream)
    $errorBody = $reader.ReadToEnd()
    $reader.Close()
    
    Write-Output "Ошибка: $($_.Exception.Response.StatusCode.value__)"
    Write-Output $errorBody
}
```