# Staywise | نظام إدارة الفندق

نسخة جاهزة للتسليم للعميل، وتعمل باستخدام Go + PostgreSQL + Docker.

## نظرة عامة
- واجهة عميل للحجز
- لوحة إدارة داخلية للفندق
- REST API لإدارة الحجوزات والغرف والنزلاء والمدفوعات
- قاعدة بيانات PostgreSQL
- تشغيل موحد عبر Docker للمراجعة والتسليم

## التشغيل السريع
```bash
cd /home/ayman-wahas/HotelManagement
docker compose up --build -d
```

ثم افتح:
- http://localhost:8089/
- http://localhost:8089/admin/
- http://localhost:8089/health

## إعدادات البيئة
```env
APP_PORT=8089
DATABASE_URL=postgres://hotel:hotel@localhost:5433/hotel_management?sslmode=disable
API_TOKEN=replace-this-with-a-long-random-secret
HOTEL_TIMEZONE=Asia/Riyadh
```

## الرمز الافتراضي للتجربة
```text
demo-token-12345
```

## الأوامر الأساسية
```bash
# التشغيل
cd /home/ayman-wahas/HotelManagement
docker compose up --build -d

# الإيقاف
cd /home/ayman-wahas/HotelManagement
docker compose down

# إعادة التشغيل
cd /home/ayman-wahas/HotelManagement
docker compose restart

# فحص الصحة
curl http://localhost:8089/health
```

## ملاحظات التسليم
- المشروع يعمل بنجاح في Docker مع PostgreSQL على المنفذ 5433
- API يعمل على المنفذ 8089
- النسخة الحالية مناسبة للمراجعة والعمل الداخلي للعميل
- يوصى عند التسليم الحقيقي استبدال `API_TOKEN` بقيمة قوية وفي بيئة خادم حقيقية

## التحقق النهائي
تم التحقق فعليًا أن المشروع يعمل بنجاح:
- `go test ./...` -> نجح
- `docker compose up --build -d` -> نجح
- `curl http://localhost:8089/health` -> أعاد `{"status":"ok"}`

## مزيد من الوثائق
- [FINAL_HANDOFF.md](FINAL_HANDOFF.md)
- [docs/QUICKSTART.md](docs/QUICKSTART.md)

## بنية المشروع
- `cmd/api` – نقطة تشغيل التطبيق
- `internal/handler` – HTTP API وواجهة الإدارة
- `internal/usecase` – منطق الأعمال
- `internal/postgres` – طبقة قاعدة البيانات
- `internal/webui` – الواجهة المدمجة
- `migrations` – ملفات ترحيل قاعدة البيانات
