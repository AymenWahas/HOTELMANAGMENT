# نظام إدارة الفندق

واجهة برمجية REST مبنية بلغة Go لإدارة عمليات الفندق، وتستخدم PostgreSQL لحفظ البيانات. يتبع المشروع بنية نظيفة تفصل طبقات التطبيق:

```text
معالج HTTP ← حالات الاستخدام ← واجهة المستودع ← مستودع PostgreSQL
```

**اللغات:** [العربية](README.md) | [English](README.en.md)

## العمليات المتاحة

- إدارة الغرف وحالتها
- ملفات النزلاء
- الحجوزات والتحقق من تعارض المواعيد وتسجيل الوصول والمغادرة والإلغاء
- قوائم الوصول والمغادرة في الاستقبال
- مهام النظافة وتجهيز الغرف
- فواتير النزلاء والمدفوعات
- الموظفون والأقسام
- الصيانة وطلبات الإصلاح
- المخزون وحركاته
- لوحة متابعة وتقارير الإشغال والإيرادات
- سجل تدقيق للتغييرات

تستخدم قاعدة البيانات جداول مترابطة وقيودًا للتحقق من صحة البيانات. تُنفّذ بعض العمليات، مثل حركات المخزون وتغييرات الحجز، ضمن معاملات PostgreSQL.

## المتطلبات

- Go 1.26.5 أو أحدث
- Docker Compose
- PostgreSQL 16 (يوفره ملف Docker Compose)

## تشغيل المشروع

من مجلد المشروع:

1. شغّل PostgreSQL:

   ```sh
   docker compose up -d postgres
   ```

2. اضبط متغيرات البيئة:

   ```sh
   export DATABASE_URL='postgres://hotel:hotel@localhost:5433/hotel_management?sslmode=disable'
   export API_TOKEN='ضع-هنا-رمزًا-سريًا-طويلًا'
   export APP_PORT=8081
   export HOTEL_TIMEZONE='Asia/Riyadh'
   ```

3. شغّل الخادم:

   ```sh
   go run ./cmd/api
   ```

عند بدء التشغيل، يطبق الخادم ترحيلات قاعدة البيانات الموجودة في `migrations/001_initial_schema.sql` ويسجل إصدار الترحيل.

افتح [http://localhost:8081](http://localhost:8081) في المتصفح، ثم أدخل قيمة `API_TOKEN`. للتحقق من جاهزية الخادم، افتح `http://localhost:8081/health`.

## الواجهة ونقاط API

الواجهة العربية المضمّنة تعرض لوحة المتابعة والاستقبال والحجوزات والغرف والنزلاء والمدفوعات. وتتصل بنقاط API على الخادم نفسه:

| الشاشة | نقاط API المستخدمة |
| --- | --- |
| الرئيسية | `GET /dashboard`، `GET /reports/occupancy`، `GET /reports/revenue`، `GET /reception/arrivals`، `GET /rooms` |
| الاستقبال | `GET /reception/arrivals`، `GET /reception/departures`، وإجراءات الوصول والمغادرة والإلغاء للحجز |
| الحجوزات | `GET /bookings`، `POST /bookings`، وإجراءات الوصول والمغادرة والإلغاء |
| الغرف | `GET /rooms`، `POST /rooms`، `PATCH /rooms/{id}` |
| النزلاء | `GET /guests`، `POST /guests` |
| المدفوعات | `GET /payments`، `GET /folios`، `POST /folios`، `POST /payments`، `PATCH /payments/{id}` |

كل مسارات `/api/v1` تتطلب ترويسة المصادقة:

```text
Authorization: Bearer <API_TOKEN>
```

## مرجع API

تستخدم الطلبات والاستجابات JSON. يعيد سجل المورد بالشكل التالي:

```json
{
  "id": "UUID",
  "kind": "guest",
  "data": {},
  "created_at": "2026-10-19T12:00:00Z",
  "updated_at": "2026-10-19T12:00:00Z"
}
```

| المورد | المسار |
| --- | --- |
| الغرف | `/api/v1/rooms` |
| النزلاء | `/api/v1/guests` |
| الحجوزات | `/api/v1/bookings` أو `/api/v1/reservations` |
| الاستقبال | `/api/v1/reception/arrivals` و`/api/v1/reception/departures` |
| النظافة | `/api/v1/housekeeping` |
| الفواتير | `/api/v1/folios` |
| المدفوعات | `/api/v1/payments` |
| الموظفون | `/api/v1/employees` |
| الصيانة | `/api/v1/maintenance` |
| عناصر المخزون | `/api/v1/inventory/items` |
| حركات المخزون | `/api/v1/inventory/movements` |

تدعم الموارد عادةً `GET` لعرض القائمة، و`POST` للإضافة، و`GET /{id}` لعرض سجل واحد، و`PATCH /{id}` للتعديل. يدعم معظمها `DELETE /{id}` للأرشفة المنطقية. تقبل القوائم `limit` (الافتراضي 50 والحد الأعلى 200) و`offset`.

### إجراءات الحجز

- `POST /api/v1/bookings/{id}/check-in`
- `POST /api/v1/bookings/{id}/check-out`
- `POST /api/v1/bookings/{id}/cancel`

الإجراءات نفسها متاحة تحت `/api/v1/reservations/{id}/...` و`/api/v1/reception/{id}/...`. تعرض قوائم الاستقبال حجوزات اليوم افتراضيًا، ويمكن اختيار تاريخ عبر `?date=YYYY-MM-DD`.

### التقارير

- `GET /api/v1/dashboard`
- `GET /api/v1/reports/occupancy`
- `GET /api/v1/reports/revenue?from=2026-10-01&to=2026-10-31`

### أمثلة الطلبات

إضافة غرفة:

```sh
curl -H "Authorization: Bearer $API_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"number":"101","room_type":"standard","floor":1,"rate":350}' \
  http://localhost:8081/api/v1/rooms
```

إنشاء حجز:

```json
{
  "room_id": "ROOM_UUID",
  "guest_id": "GUEST_UUID",
  "check_in": "2026-10-21",
  "check_out": "2026-10-24",
  "adults": 2,
  "rate": 350,
  "source": "direct"
}
```

حقول النزيل المطلوبة هي `first_name` و`last_name`. وتحتاج الدفعة إلى `folio_id` و`amount` و`method`. تحتاج حركة المخزون إلى `item_id` و`quantity` و`movement_type`، وتقبل أنواع الحركة `in` و`out` و`restock` و`use`.

## بنية المشروع

- `internal/domain`: أنواع السجلات والتقارير
- `internal/handler`: مسارات HTTP وJSON والمصادقة
- `internal/usecase`: التحقق من المدخلات وقواعد العمليات
- `internal/repository`: واجهة التخزين التي تستخدمها حالات الاستخدام
- `internal/postgres`: استعلامات PostgreSQL والمعاملات وسجل التدقيق
- `internal/webui`: ملفات الواجهة المضمّنة في الخادم
- `migrations`: مخطط قاعدة البيانات
- `cmd/api`: ربط الاعتماديات وتشغيل الخادم
# HOTELMANAGMENT
