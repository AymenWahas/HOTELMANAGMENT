(() => {
  const API = '/api/v1';
  const TOKEN_KEY = 'staywise.api-token';
  const state = { token: localStorage.getItem(TOKEN_KEY) || '', page: 'overview', receptionTab: 'arrivals', date: '', cache: {}, search: '' };
  const translations = {
    'Your token was rejected. Check it and try again.':'رمز الدخول غير صحيح. تحقق منه وحاول مرة أخرى.',
    'Enter your hotel API access token to open the workspace.':'أدخل رمز الدخول لفتح لوحة إدارة الفندق.',
    'Your hotel API access token to open the workspace.':'أدخل رمز الدخول لفتح لوحة إدارة الفندق.',
    'Your token stays in this browser and is sent only to this hotel server.':'يُحفظ الرمز في هذا المتصفح ويُرسل إلى خادم الفندق فقط.',
    'No arrivals scheduled for today':'لا توجد حالات وصول مجدولة اليوم',
    'No bookings yet. Create one to get started.':'لا توجد حجوزات حتى الآن. أنشئ حجزًا للبدء.',
    'No rooms added yet. Add your room inventory.':'لم تتم إضافة غرف بعد. أضف غرف الفندق.',
    'No guest profiles yet. Add a guest or create a booking.':'لا توجد ملفات للنزلاء. أضف نزيلًا أو أنشئ حجزًا.',
    'No payments recorded yet. Create a guest folio to begin.':'لا توجد مدفوعات مسجلة. أنشئ فاتورة نزيل للبدء.',
    'Guest profile unavailable':'ملف النزيل غير متاح',
    'Search payment reference or guest':'ابحث برقم الدفعة أو اسم النزيل',
    'Search guest, room, or status':'ابحث عن نزيل أو غرفة أو حالة',
    'Search room number or type':'ابحث عن رقم الغرفة أو نوعها',
    'Search name, email, or phone':'ابحث بالاسم أو البريد أو الهاتف',
    'Guest folio created.':'تم إنشاء فاتورة النزيل.',
    'Booking created.':'تم إنشاء الحجز.',
    'Booking request approved.':'تمت الموافقة على طلب الحجز.',
    'Approve this booking request?':'هل تريد الموافقة على طلب الحجز؟',
    'Approve booking request':'الموافقة على طلب الحجز',
    'Decline booking request':'رفض طلب الحجز',
    'Guest profile created.':'تم إنشاء ملف النزيل.',
    'Payment recorded.':'تم تسجيل الدفعة.',
    'Room status updated.':'تم تحديث حالة الغرفة.',
    'Payment marked as refunded.':'تم تسجيل استرداد الدفعة.',
    'A housekeeping task was opened.':'وتم إنشاء مهمة لتنظيف الغرفة.',
    'Add a guest before creating a booking.':'أضف نزيلًا قبل إنشاء الحجز.',
    'Add a room before creating a booking.':'أضف غرفة قبل إنشاء الحجز.',
    'Add a guest before opening a folio.':'أضف نزيلًا قبل إنشاء الفاتورة.',
    'Create an open folio before recording a payment.':'أنشئ فاتورة مفتوحة قبل تسجيل الدفعة.',
    'Check-out must be after check-in.':'يجب أن يكون تاريخ المغادرة بعد تاريخ الوصول.',
    'Room availability is checked when the booking is saved.':'سيتم التحقق من توفر الغرفة عند حفظ الحجز.',
    'Create the bill that payments will be posted against.':'أنشئ الفاتورة التي ستُسجّل عليها المدفوعات.',
    'Post a payment to an open guest folio.':'سجّل دفعة على فاتورة نزيل مفتوحة.',
    'Choose the current operational status for this room.':'اختر الحالة التشغيلية الحالية للغرفة.',
    'Add this room to your hotel inventory.':'أضف هذه الغرفة إلى قائمة غرف الفندق.',
    'Create a guest profile for bookings and reception.':'أنشئ ملف نزيل لاستخدامه في الحجوزات والاستقبال.',
    'Add a stay for a guest. Check-out must be after check-in.':'أضف إقامة للنزيل. يجب أن يكون تاريخ المغادرة بعد الوصول.',
    'Make every arrival feel like they’re already home.':'اجعل كل نزيل يشعر وكأنه في منزله منذ لحظة وصوله.',
    'Welcome your guests and prepare their rooms.':'رحّب بالنزلاء وجهّز غرفهم.',
    'Visit reception to check in arrivals.':'انتقل إلى الاستقبال لتسجيل وصول النزلاء.',
    'Keep the important handoffs close.':'تابع المهام والتسليمات المهمة.',
    'Rooms needing attention are marked in room status.':'تظهر الغرف التي تحتاج إلى متابعة ضمن حالة الغرف.',
    'Track and assign follow-up in operations.':'تابع أعمال الصيانة ووزّع مهام المتابعة.',
    'Completed payments by day.':'المدفوعات المكتملة يوميًا.',
    'Guest folios and recorded guest payments.':'فواتير النزلاء والمدفوعات المسجلة.',
    'Room availability and housekeeping status.':'توفر الغرف وحالة النظافة.',
    'Guest profiles, contact details, and guest notes.':'ملفات النزلاء وبيانات التواصل والملاحظات.',
    'Manage stays from reservation to checkout.':'أدر الإقامات من الحجز حتى المغادرة.',
    'Arrivals, departures, and the next guest handoff.':'تابع الوصول والمغادرة وتسليم الغرف.',
    'A clear view of the day and the people checking in.':'نظرة واضحة على يومك والنزلاء المتوقع وصولهم.',
    'No booking attached':'بدون حجز مرتبط',
    'Travel agency':'وكالة سفر',
    'Bank transfer':'تحويل بنكي',
    'Hotel workspace':'مساحة عمل الفندق',
    'Guest care, in motion':'متابعة النزلاء',
    'Today’s arrivals':'حالات الوصول اليوم',
    'Good morning. Here’s your day.':'صباح الخير، إليك ملخص يومك.',
    'Today at a glance':'ملخص اليوم',
    'TODAY AT A GLANCE':'ملخص اليوم',
    'HOTEL OPERATIONS':'إدارة الفندق',
    'HOTEL OPERATIONS':'إدارة الفندق',
    'My property':'منشأتي الفندقية',
    'WORKSPACE':'القائمة الرئيسية',
    'Overview':'الرئيسية',
    'Reception':'الاستقبال',
    'Bookings':'الحجوزات',
    'Rooms':'الغرف',
    'Guests':'النزلاء',
    'Payments':'المدفوعات',
    'Hotel':'الفندق',
    'Staywise home':'الصفحة الرئيسية',
    'Main navigation':'القائمة الرئيسية',
    'Your hotel, at a glance':'ملخص الفندق',
    'Arrivals':'القادمون',
    'Departures':'المغادرون',
    'New booking':'حجز جديد',
    'Add guest':'إضافة نزيل',
    'Add room':'إضافة غرفة',
    'Record payment':'تسجيل دفعة',
    'New folio':'فاتورة جديدة',
    'Open workspace':'دخول إلى لوحة الفندق',
    'API access token':'رمز الدخول',
    'Paste your access token':'أدخل رمز الدخول',
    'Good to have you':'سعداء بعودتك',
    'back.':'إلى الفندق',
    'staywise':'Staywise',
    'Every stay,':'كل إقامة،',
    'in good hands.':'في أيدٍ أمينة.',
    'Front desk · Rooms · Guest care':'الاستقبال · الغرف · رعاية النزلاء',
    'Need a hand?':'هل تحتاج إلى مساعدة؟',
    'Ask your administrator for access':'تواصل مع مسؤول النظام لمنحك صلاحية الدخول',
    'Hotel Manager':'مدير الفندق',
    'Sign out':'تسجيل الخروج',
    'Refresh data':'تحديث البيانات',
    'Toggle navigation':'فتح القائمة',
    'Open reception':'فتح الاستقبال',
    'View reception':'عرض الاستقبال',
    'All rooms':'كل الغرف',
    'Room status':'حالة الغرفة',
    'Live view of your property.':'حالة الفندق الآن.',
    'Occupancy':'الإشغال',
    'Revenue this week':'إيرادات هذا الأسبوع',
    'Collected today':'المحصّل اليوم',
    'Revenue today':'إيرادات اليوم',
    'Rooms in house':'الغرف المشغولة',
    'Available rooms':'الغرف المتاحة',
    'Arrivals today':'القادمون اليوم',
    'Payments recorded':'المدفوعات المسجلة',
    'Completed payments':'المدفوعات المكتملة',
    'Open folios':'الفواتير المفتوحة',
    'Refunded':'المبالغ المستردة',
    'all time':'الإجمالي',
    'collected':'تم تحصيلها',
    'awaiting settlement':'بانتظار السداد',
    'recorded':'مسجل',
    'ready to assign':'جاهزة للتخصيص',
    'expected today':'متوقع اليوم',
    'payments recorded':'مدفوعات مسجلة',
    'arriving':'قادمون',
    'checking out':'مغادرون',
    'room tasks':'مهام الغرف',
    'guests are expected today':'نزيلًا متوقع وصولهم اليوم',
    'housekeeping tasks to follow':'مهام تنظيف تحتاج إلى متابعة',
    'open maintenance requests':'طلبات صيانة مفتوحة',
    'adult':'بالغ',
    'adults':'بالغون',
    'child':'طفل',
    'children':'أطفال',
    'Room ':'الغرفة ',
    'Floor ':'الطابق ',
    ' / night':' / الليلة',
    'Rate / night':'السعر / الليلة',
    'RATE / NIGHT':'السعر / الليلة',
    'Check in':'تسجيل الوصول',
    'Check out':'تسجيل المغادرة',
    'Cancel booking':'إلغاء الحجز',
    'Cancel':'إلغاء',
    'Refund':'استرداد',
    'Update status':'تحديث الحالة',
    'Update room status':'تحديث حالة الغرفة',
    'Search':'بحث',
    'BOOKING':'الحجز',
    'GUEST':'النزيل',
    'ROOM':'الغرفة',
    'STAY DATES':'تواريخ الإقامة',
    'STATUS':'الحالة',
    'EMAIL':'البريد الإلكتروني',
    'PHONE':'الهاتف',
    'DOCUMENT':'الوثيقة',
    'ADDED':'تاريخ الإضافة',
    'PAYMENT':'الدفعة',
    'DATE':'التاريخ',
    'METHOD':'طريقة الدفع',
    'AMOUNT':'المبلغ',
    'total bookings':'حجزًا إجمالًا',
    'guest profiles':'ملف نزيل',
    'rooms':'غرف',
    'expected':'متوقع',
    'departing':'مغادرون',
    'Arrival':'الوصول',
    'Departure':'المغادرة',
    'available':'متاحة',
    'occupied':'مشغولة',
    'needs cleaning':'تحتاج إلى تنظيف',
    'out of order':'خارج الخدمة',
    'checked in':'تم تسجيل الوصول',
    'checked out':'تم تسجيل المغادرة',
    'confirmed':'مؤكد',
    'pending':'قيد المراجعة',
    'completed':'مكتمل',
    'cancelled':'ملغي',
    'refunded':'مسترد',
    'failed':'فشل',
    'open':'مفتوح',
    'done':'منجز',
    'active':'نشط',
    'inactive':'غير نشط',
    'in progress':'قيد التنفيذ',
    'on leave':'في إجازة',
    'void':'ملغى',
    'No arrivals for':'لا يوجد وصول بتاريخ',
    'No departures for':'لا توجد مغادرات بتاريخ',
    'No arrivals scheduled for':'لا توجد حالات وصول مجدولة بتاريخ',
    'Loading hotel operations…':'جارٍ تحميل لوحة الفندق…',
    'Try again':'حاول مرة أخرى',
    'Saving…':'جارٍ الحفظ…',
    'Save guest':'حفظ بيانات النزيل',
    'Save room':'حفظ الغرفة',
    'Create booking':'إنشاء الحجز',
    'Create folio':'إنشاء الفاتورة',
    'Save payment':'حفظ الدفعة',
    'Save status':'حفظ الحالة',
    'Add a guest':'إضافة نزيل',
    'Create a booking':'إنشاء حجز',
    'Add a room':'إضافة غرفة',
    'Open a guest folio':'إنشاء فاتورة نزيل',
    'Record a payment':'تسجيل دفعة',
    'Cancel this booking?':'هل تريد إلغاء هذا الحجز؟',
    'Mark this payment as refunded?':'هل تريد تسجيل استرداد هذه الدفعة؟',
    'Close':'إغلاق',
    'Cancel':'إلغاء',
    'Saving':'جارٍ الحفظ',
    'First name':'الاسم الأول',
    'Last name':'اسم العائلة',
    'Email':'البريد الإلكتروني',
    'Phone':'رقم الهاتف',
    'Nationality':'الجنسية',
    'Document number':'رقم الهوية أو جواز السفر',
    'Guest notes':'ملاحظات النزيل',
    'Guest':'النزيل',
    'Room':'الغرفة',
    'Check-in':'تاريخ الوصول',
    'Check-out':'تاريخ المغادرة',
    'Adults':'عدد البالغين',
    'Children':'عدد الأطفال',
    'Nightly rate (SAR)':'سعر الليلة (ريال)',
    'Booking source':'مصدر الحجز',
    'Special requests':'طلبات خاصة',
    'Room number':'رقم الغرفة',
    'Room type':'نوع الغرفة',
    'Floor':'الطابق',
    'Initial status':'الحالة الأولية',
    'Notes':'ملاحظات',
    'Booking':'الحجز',
    'Folio total (SAR)':'إجمالي الفاتورة (ريال)',
    'Currency':'العملة',
    'Guest folio':'فاتورة النزيل',
    'Amount (SAR)':'المبلغ (ريال)',
    'Payment method':'طريقة الدفع',
    'Reference':'رقم المرجع',
    'Room status':'حالة الغرفة',
    'Saudi':'سعودي',
    'Available':'متاحة',
    'Needs cleaning':'تحتاج إلى تنظيف',
    'Out of order':'خارج الخدمة',
    'Direct':'مباشر',
    'Phone':'الهاتف',
    'Website':'الموقع الإلكتروني',
    'Other':'أخرى',
    'Cash':'نقدًا',
    'Card':'بطاقة',
    'Online':'إلكتروني',
    'e.g. Maya':'مثال: منى',
    'e.g. Hassan':'مثال: حسن',
    'Passport or ID':'جواز سفر أو هوية',
    'Preferences or special requirements':'التفضيلات أو المتطلبات الخاصة',
    'Deluxe king':'غرفة ديلوكس بسرير كبير',
    'Room features or notes':'مواصفات الغرفة أو ملاحظات',
    'Anything we should prepare for this stay?':'هل هناك شيء نجهزه لهذه الإقامة؟',
    'Room charges, dining, or other services':'رسوم الغرفة أو الطعام أو خدمات أخرى',
    'Receipt or transaction ID':'رقم الإيصال أو العملية',
    'Access token rejected.':'رمز الدخول غير مقبول.',
    'Request failed':'فشل الطلب',
    'Failed to fetch':'تعذر الاتصال بالخادم. تحقق من الاتصال ثم حاول مجددًا.',
    'internal server error':'تعذر حفظ البيانات بسبب خطأ في الخادم. حاول مرة أخرى.',
    'value already exists: guests_email_unique':'يوجد نزيل مسجل بهذا البريد الإلكتروني. استخدم بريدًا مختلفًا أو اترك الحقل فارغًا.'
  };
  const translationEntries = Object.entries(translations).sort((a,b) => b[0].length-a[0].length);
  const translationPattern = new RegExp(translationEntries.map(([source]) => source.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')).join('|'), 'g');
  function translateText(value) {
    return String(value).replace(translationPattern, source => translations[source]);
  }
  function localizeDOM(root) {
    const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT);
    let node;
    while ((node = walker.nextNode())) {
      const translated = translateText(node.nodeValue);
      if (translated !== node.nodeValue) node.nodeValue = translated;
    }
    if (root.nodeType === Node.ELEMENT_NODE) {
      for (const name of ['placeholder','title','aria-label']) {
        if (root.hasAttribute(name)) {
          const value = root.getAttribute(name);
          const translated = translateText(value);
          if (translated !== value) root.setAttribute(name, translated);
        }
      }
    }
    root.querySelectorAll?.('[placeholder],[title],[aria-label]').forEach(element => {
      for (const name of ['placeholder','title','aria-label']) {
        if (element.hasAttribute(name)) {
          const value = element.getAttribute(name);
          const translated = translateText(value);
          if (translated !== value) element.setAttribute(name, translated);
        }
      }
    });
  }
  localizeDOM(document.body);
  const $ = (selector, root = document) => root.querySelector(selector);
  const escapeHTML = (value = '') => String(value).replace(/[&<>"']/g, char => ({ '&':'&amp;', '<':'&lt;', '>':'&gt;', '"':'&quot;', "'":'&#39;' })[char]);
  const today = () => new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Riyadh' }).format(new Date());
  const money = amount => new Intl.NumberFormat('ar-SA', { style:'currency', currency:'SAR', maximumFractionDigits:0 }).format(Number(amount || 0));
  const dateText = value => value ? new Intl.DateTimeFormat('ar-SA', { month:'short', day:'numeric', year:'numeric', timeZone:'UTC' }).format(new Date(`${String(value).slice(0,10)}T12:00:00Z`)) : '—';
  const initials = name => String(name || 'Guest').trim().split(/\s+/).slice(0,2).map(part => part[0] || '').join('').toUpperCase();
  const unwrap = value => value?.data && typeof value.data === 'object' ? value.data : {};
  const fullGuestName = (value, fallback = 'Guest') => {
    if (value === null || value === undefined) return fallback;
    if (typeof value === 'object') {
      const first = String(value.first_name || value.firstName || '').trim();
      const last = String(value.last_name || value.lastName || '').trim();
      const full = [first, last].filter(Boolean).join(' ');
      return full || fallback;
    }
    const text = String(value).trim();
    return text || fallback;
  };
  const resolveGuestName = (guestId, guestMap, fallback = 'Guest profile unavailable') => {
    const value = guestId && guestMap ? guestMap.get(guestId) : undefined;
    if (value) return fullGuestName(value, fallback);
    return fallback;
  };

  async function request(path, options = {}) {
    const headers = { 'Content-Type':'application/json', 'Authorization':`Bearer ${state.token}`, ...(options.headers || {}) };
    const response = await fetch(`${API}${path}`, { ...options, headers });
    if (response.status === 401) { showLogin('Your token was rejected. Check it and try again.'); throw new Error('Access token rejected.'); }
    if (!response.ok) {
      const body = await response.json().catch(() => ({}));
      throw new Error(body.error || `Request failed (${response.status})`);
    }
    return response.status === 204 ? null : response.json();
  }
  const list = async path => (await request(path)).items || [];
  const post = (path, data) => request(path, { method:'POST', body:JSON.stringify(data) });
  const patch = (path, data) => request(path, { method:'PATCH', body:JSON.stringify(data) });

  function showLogin(error = '') {
    $('#app-shell').hidden = true;
    $('#login-screen').hidden = false;
    $('#login-error').textContent = translateText(error);
    $('#api-token').value = state.token;
  }
  function showApp() {
    $('#login-screen').hidden = true;
    $('#app-shell').hidden = false;
    $('#today-label').textContent = new Intl.DateTimeFormat('ar-SA', { weekday:'short', month:'short', day:'numeric', timeZone:'Asia/Riyadh' }).format(new Date());
    navigate(state.page);
  }
  async function signIn(token) {
    state.token = token.trim();
    try {
      await request('/dashboard');
      localStorage.setItem(TOKEN_KEY, state.token);
      showApp();
    } catch (error) {
      if (error.message !== 'Access token rejected.') $('#login-error').textContent = translateText(error.message);
    }
  }

  const titles = {
    overview:['Your hotel, at a glance','A clear view of the day and the people checking in.'],
    reception:['Reception','Arrivals, departures, and the next guest handoff.'],
    bookings:['Bookings','Manage stays from reservation to checkout.'],
    rooms:['Rooms','Room availability and housekeeping status.'],
    guests:['Guests','Guest profiles, contact details, and guest notes.'],
    payments:['Payments','Guest folios and recorded guest payments.']
  };
  function navigate(page) {
    if (!titles[page]) return;
    state.page = page; state.search = '';
    document.querySelectorAll('.nav-item').forEach(item => item.classList.toggle('active', item.dataset.page === page));
    $('#page-crumb').textContent = translateText(page.charAt(0).toUpperCase() + page.slice(1));
    $('#sidebar').classList.remove('open');
    renderPage();
  }
  function pageHeading(page, action = '') {
    const [title, subtitle] = titles[page];
    return `<div class="page-heading"><div><p class="eyebrow">${page === 'overview' ? 'TODAY AT A GLANCE' : 'HOTEL OPERATIONS'}</p><h1>${title}</h1><p>${subtitle}</p></div>${action}</div>`;
  }
  function stat(label, value, icon, note = '') {
    return `<article class="stat-card"><div class="stat-top"><span>${label}</span><span class="stat-icon">${icon}</span></div><div class="stat-value">${value}<span class="stat-note">${note}</span></div></article>`;
  }
  function statusBadge(status) {
    const value = String(status || 'available').toLowerCase();
    const cls = ['available','checked_in','completed','done','active'].includes(value) ? 'green' : ['needs_cleaning','confirmed','pending','in_progress','open','on_leave'].includes(value) ? 'amber' : ['out_of_order','cancelled','failed','inactive','refunded','void'].includes(value) ? 'red' : 'gray';
    return `<span class="badge badge-${cls}">${escapeHTML(value.replaceAll('_',' '))}</span>`;
  }
  function person(name, detail = '') {
    return `<div class="table-person"><span class="guest-initial">${escapeHTML(initials(name))}</span><span><strong>${escapeHTML(name || 'Guest')}</strong>${detail ? `<small>${escapeHTML(detail)}</small>` : ''}</span></div>`;
  }
  function empty(message, icon = '⌑') { return `<div class="empty-state"><span class="empty-icon">${icon}</span>${escapeHTML(message)}</div>`; }
  function loading() { $('#page-content').innerHTML = '<div class="page-loading"><div class="load-wrap"><span class="spinner"></span>Loading hotel operations…</div></div>'; localizeDOM($('#page-content')); }
  function failPage(error) { $('#page-content').innerHTML = `<div class="inline-error">${escapeHTML(error.message)} <button class="text-link" data-action="refresh">Try again</button></div>`; localizeDOM($('#page-content')); }

  async function renderPage() {
    loading();
    try {
      if (state.page === 'overview') await overview();
      if (state.page === 'reception') await reception();
      if (state.page === 'bookings') await bookings();
      if (state.page === 'rooms') await rooms();
      if (state.page === 'guests') await guests();
      if (state.page === 'payments') await payments();
      localizeDOM($('#page-content'));
    } catch (error) { failPage(error); }
  }

  async function overview() {
    const [dashboard, occupancy, arrivals, roomsList] = await Promise.all([
      request('/dashboard'), request('/reports/occupancy'), list(`/reception/arrivals?date=${today()}`), list('/rooms?limit=200')
    ]);
    state.cache.arrivals = arrivals; state.cache.rooms = roomsList;
    const arrivalRows = arrivals.slice(0,5).map(entry => `<div class="arrival-row"><div class="guest-cell"><span class="guest-initial">${escapeHTML(initials(entry.guest_name))}</span><span><strong>${escapeHTML(entry.guest_name)}</strong><small>${entry.adults} adult${entry.adults === 1 ? '' : 's'}${entry.children ? ` · ${entry.children} child${entry.children === 1 ? '' : 'ren'}` : ''}</small></span></div><div class="arrival-detail"><strong>Room ${escapeHTML(entry.room_number)}</strong>${dateText(entry.check_in)} – ${dateText(entry.check_out)}</div><span class="room-label">${escapeHTML(entry.room_number)}</span><button class="button button-small button-secondary" data-action="check-in" data-id="${escapeHTML(entry.reservation_id)}">Check in</button></div>`).join('');
    const roomStates = [
      ['Available', occupancy.available, 'dot'], ['Occupied', occupancy.occupied, 'dot'],
      ['Needs cleaning', occupancy.needs_cleaning, 'dot-amber'], ['Out of order', occupancy.out_of_order, 'dot-red']
    ];
    const bars = await revenueBars();
    $('#page-content').innerHTML = `${pageHeading('overview')}<section class="welcome-banner"><div class="welcome-copy"><p class="eyebrow">${new Intl.DateTimeFormat('ar-SA', { weekday:'long', timeZone:'Asia/Riyadh' }).format(new Date())} · ${today()}</p><h1>Good morning. Here’s your day.</h1><p>Make every arrival feel like they’re already home.</p><div class="welcome-meta"><span><strong>${dashboard.arrivals_today}</strong> arriving</span><span><strong>${dashboard.departures_today}</strong> checking out</span><span><strong>${dashboard.open_housekeeping_tasks}</strong> room tasks</span></div></div><button class="button banner-action" data-page="reception">Open reception <span>→</span></button></section><section class="stats-grid">${stat('Rooms in house',`${dashboard.rooms_occupied}<span class="stat-note"> / ${dashboard.rooms_total}</span>`,'▦',`${occupancy.occupancy_percent.toFixed(0)}% occupied`)}${stat('Available rooms',dashboard.rooms_available,'⌂','ready to assign')}${stat('Arrivals today',dashboard.arrivals_today,'↘','expected today')}${stat('Revenue today',money(dashboard.revenue_today),'◉','payments recorded')}</section><section class="content-grid"><article class="panel"><div class="panel-heading"><div><h2>Today’s arrivals</h2><p>Welcome your guests and prepare their rooms.</p></div><button class="text-link" data-page="reception">View reception →</button></div><div class="arrivals-list">${arrivalRows || empty('No arrivals scheduled for today','⌑')}</div></article><article class="panel"><div class="panel-heading"><div><h2>Room status</h2><p>Live view of your property.</p></div><button class="text-link" data-page="rooms">All rooms →</button></div><div class="room-summary">${roomStates.map(([name,count,dot]) => `<div class="room-summary-line"><i class="dot ${dot}"></i><span>${name}</span><strong>${count}</strong></div>`).join('')}</div><div class="progress-track"><span style="width:${Math.min(100,occupancy.occupancy_percent)}%"></span></div><div class="occupancy-caption"><span>Occupancy</span><strong>${occupancy.occupancy_percent.toFixed(0)}%</strong></div></article></section><section class="bottom-grid"><article class="panel"><div class="panel-heading"><div><h2>Guest care, in motion</h2><p>Keep the important handoffs close.</p></div></div><div class="activity-item"><span class="activity-mark">↘</span><span><strong>${dashboard.arrivals_today} guests are expected today</strong><small>Visit reception to check in arrivals.</small></span></div><div class="activity-item"><span class="activity-mark">✳</span><span><strong>${dashboard.open_housekeeping_tasks} housekeeping tasks to follow</strong><small>Rooms needing attention are marked in room status.</small></span></div><div class="activity-item"><span class="activity-mark">⚒</span><span><strong>${dashboard.open_maintenance_requests} open maintenance requests</strong><small>Track and assign follow-up in operations.</small></span></div></article><article class="panel"><div class="panel-heading"><div><h2>Revenue this week</h2><p>Completed payments by day.</p></div></div><div class="revenue-big">${money(dashboard.revenue_today)}</div><div class="revenue-sub">Collected today</div><div class="revenue-bars">${bars}</div></article></section>`;
  }

  async function revenueBars() {
    const end = today(); const startDate = new Date(`${end}T12:00:00+03:00`); startDate.setDate(startDate.getDate()-6);
    const start = new Intl.DateTimeFormat('en-CA',{timeZone:'Asia/Riyadh'}).format(startDate);
    const data = await request(`/reports/revenue?from=${start}&to=${end}`);
    const amounts = data.map(row => Number(row.revenue)); const max = Math.max(...amounts,1);
    const lookup = new Map(data.map(row=>[row.date,Number(row.revenue)]));
    return Array.from({length:7},(_,index)=>{ const date = new Date(`${start}T12:00:00Z`); date.setUTCDate(date.getUTCDate()+index); const key = date.toISOString().slice(0,10); const height=Math.max(8,((lookup.get(key)||0)/max)*58); return `<i title="${escapeHTML(key)} · ${escapeHTML(money(lookup.get(key)||0))}" style="height:${height}px"></i>`; }).join('');
  }

  async function reception() {
    const date = state.date || today(); const action = state.receptionTab;
    const entries = await list(`/reception/${action}?date=${encodeURIComponent(date)}`);
    const cards = entries.map(entry => {
      const guestName = fullGuestName(entry.guest_name || entry.guestName || 'Guest', 'Guest');
      return `<article class="reception-card"><div class="guest-cell"><span class="guest-initial">${escapeHTML(initials(guestName))}</span><span><strong>${escapeHTML(guestName)}</strong><small>${entry.adults} adult${entry.adults===1?'':'s'}${entry.children?` · ${entry.children} child${entry.children===1?'':'ren'}`:''}</small></span></div><div class="reception-meta"><strong>Room ${escapeHTML(entry.room_number)}</strong><small>Room ID ${escapeHTML(entry.room_id.slice(0,8))}</small></div><div class="reception-meta"><strong>${action==='arrivals'?'Arrival':'Departure'}</strong><small>${dateText(action==='arrivals'?entry.check_in:entry.check_out)} · ${escapeHTML(entry.status.replaceAll('_',' '))}</small></div><div class="reception-actions">${action==='arrivals'?`<button class="button button-small button-primary" data-action="check-in" data-id="${escapeHTML(entry.reservation_id)}">Check in</button><button class="button button-small button-secondary" data-action="cancel-booking" data-id="${escapeHTML(entry.reservation_id)}">Cancel</button>`:`<button class="button button-small button-primary" data-action="check-out" data-id="${escapeHTML(entry.reservation_id)}">Check out</button>`}</div></article>`;
    }).join('');
    $('#page-content').innerHTML = `${pageHeading('reception',`<button class="button button-primary" data-action="new-booking">＋ New booking</button>`)}<div class="toolbar"><div class="date-tabs"><button class="date-tab ${action==='arrivals'?'active':''}" data-reception-tab="arrivals">Arrivals</button><button class="date-tab ${action==='departures'?'active':''}" data-reception-tab="departures">Departures</button></div><input type="date" class="date-input" id="queue-date" value="${escapeHTML(date)}"><span class="badge badge-gray">${entries.length} ${action==='arrivals'?'expected':'departing'}</span></div><div class="reception-list">${cards || empty(`No ${action} for ${dateText(date)}.`, '⌑')}</div>`;
  }

  async function bookings() {
    const [records, guestsRecords, roomRecords] = await Promise.all([list('/bookings?limit=200'),list('/guests?limit=200'),list('/rooms?limit=200')]);
    state.cache.guests = guestsRecords; state.cache.rooms = roomRecords; state.cache.bookings = records;
    const guestMap = new Map(guestsRecords.map(row => [row.id, fullGuestName(unwrap(row), 'Guest profile unavailable')]));
    const roomMap = new Map(roomRecords.map(row=>[row.id,unwrap(row).number]));
    const rows = records.filter(record => {
      const d=unwrap(record); const text=`${resolveGuestName(d.guest_id, guestMap, d.guest_name || '')} ${roomMap.get(d.room_id)||''} ${d.status||''}`.toLowerCase(); return text.includes(state.search.toLowerCase());
    }).map(record=>{const d=unwrap(record);const guestName=resolveGuestName(d.guest_id, guestMap, d.guest_name || 'Guest profile unavailable');const actions=d.status==='pending'?`<button class="icon-action" title="Approve booking request" data-action="confirm-booking" data-id="${record.id}">✓</button><button class="icon-action" title="Decline booking request" data-action="cancel-booking" data-id="${record.id}">×</button>`:d.status==='confirmed'?`<button class="icon-action" title="Check in" data-action="check-in" data-id="${record.id}">↘</button><button class="icon-action" title="Cancel booking" data-action="cancel-booking" data-id="${record.id}">×</button>`:d.status==='checked_in'?`<button class="icon-action" title="Check out" data-action="check-out" data-id="${record.id}">↗</button>`:'—';return `<tr><td><strong>#${escapeHTML(record.id.slice(0,8).toUpperCase())}</strong></td><td>${person(guestName)}</td><td>Room ${escapeHTML(roomMap.get(d.room_id)||'—')}</td><td>${dateText(d.check_in)} – ${dateText(d.check_out)}</td><td>${statusBadge(d.status)}</td><td>${money(d.rate||0)}</td><td><div class="row-actions">${actions}</div></td></tr>`;}).join('');
    $('#page-content').innerHTML=`${pageHeading('bookings',`<button class="button button-primary" data-action="new-booking">＋ New booking</button>`)}<div class="toolbar"><div class="search-box"><span>⌕</span><input id="search-input" placeholder="Search guest, room, or status" value="${escapeHTML(state.search)}"></div><span class="badge badge-gray">${records.length} total bookings</span></div><div class="table-wrap"><table class="data-table"><thead><tr><th>BOOKING</th><th>GUEST</th><th>ROOM</th><th>STAY DATES</th><th>STATUS</th><th>RATE / NIGHT</th><th></th></tr></thead><tbody>${rows || `<tr><td colspan="7">${empty('No bookings yet. Create one to get started.','▤')}</td></tr>`}</tbody></table></div>`;
  }

  async function rooms() {
    const records=await list('/rooms?limit=200'); state.cache.rooms=records;
    const html=records.filter(row=>`${unwrap(row).number} ${unwrap(row).room_type} ${unwrap(row).status}`.toLowerCase().includes(state.search.toLowerCase())).map(record=>{const d=unwrap(record);return `<article class="room-card"><div class="room-card-top"><span class="room-number">${escapeHTML(d.number)}</span>${statusBadge(d.status)}</div><p>${escapeHTML(d.room_type)} · Floor ${escapeHTML(d.floor||1)}</p><div class="room-card-bottom"><span class="room-rate"><strong>${money(d.rate)}</strong> / night</span><button class="text-link" data-action="room-status" data-id="${record.id}" data-status="${escapeHTML(d.status)}">Update status</button></div></article>`;}).join('');
    $('#page-content').innerHTML=`${pageHeading('rooms',`<button class="button button-primary" data-action="new-room">＋ Add room</button>`)}<div class="toolbar"><div class="search-box"><span>⌕</span><input id="search-input" placeholder="Search room number or type" value="${escapeHTML(state.search)}"></div><span class="badge badge-gray">${records.length} rooms</span></div><div class="room-grid">${html||`<div class="panel list-empty">${empty('No rooms added yet. Add your room inventory.','▦')}</div>`}</div>`;
  }

  async function guests() {
    const records=await list('/guests?limit=200'); state.cache.guests=records;
    const rows=records.filter(row=>{const d=unwrap(row);return `${d.first_name} ${d.last_name} ${d.email||''} ${d.phone||''}`.toLowerCase().includes(state.search.toLowerCase());}).map(record=>{const d=unwrap(record);return `<tr><td>${person(`${d.first_name||''} ${d.last_name||''}`.trim(),d.nationality||'Guest')}</td><td>${escapeHTML(d.email||'—')}</td><td>${escapeHTML(d.phone||'—')}</td><td>${escapeHTML(d.document_number||'—')}</td><td>${dateText(record.created_at.slice(0,10))}</td><td><button class="button button-small button-secondary" data-action="guest-booking" data-id="${record.id}">New booking</button></td></tr>`;}).join('');
    $('#page-content').innerHTML=`${pageHeading('guests',`<button class="button button-primary" data-action="new-guest">＋ Add guest</button>`)}<div class="toolbar"><div class="search-box"><span>⌕</span><input id="search-input" placeholder="Search name, email, or phone" value="${escapeHTML(state.search)}"></div><span class="badge badge-gray">${records.length} guest profiles</span></div><div class="table-wrap"><table class="data-table"><thead><tr><th>GUEST</th><th>EMAIL</th><th>PHONE</th><th>DOCUMENT</th><th>ADDED</th><th></th></tr></thead><tbody>${rows||`<tr><td colspan="6">${empty('No guest profiles yet. Add a guest or create a booking.','♙')}</td></tr>`}</tbody></table></div>`;
  }

  async function payments() {
    const [paymentsList,foliosList,guestsList]=await Promise.all([list('/payments?limit=200'),list('/folios?limit=200'),list('/guests?limit=200')]);
    state.cache.folios=foliosList;state.cache.guests=guestsList;state.cache.payments=paymentsList;
    const guestMap=new Map(guestsList.map(row=>[row.id,fullGuestName(unwrap(row), 'Guest')]));
    const folioMap=new Map(foliosList.map(row=>[row.id,row]));
    const rows=paymentsList.filter(record=>{const d=unwrap(record);const folio=folioMap.get(d.folio_id);const f=unwrap(folio||{});return `${record.id} ${guestMap.get(f.guest_id)||''} ${d.method||''} ${d.reference||''} ${d.status||''}`.toLowerCase().includes(state.search.toLowerCase());}).map(record=>{const d=unwrap(record);const folio=folioMap.get(d.folio_id);const f=unwrap(folio||{});return `<tr><td><strong>#${escapeHTML(record.id.slice(0,8).toUpperCase())}</strong></td><td>${escapeHTML(guestMap.get(f.guest_id)||'Guest')}</td><td>${dateText(record.created_at.slice(0,10))}</td><td>${escapeHTML(d.method||'—').replaceAll('_',' ')}</td><td><strong>${money(d.amount)}</strong></td><td>${statusBadge(d.status)}</td><td>${d.status==='completed'?`<button class="button button-small button-secondary" data-action="refund-payment" data-id="${record.id}">Refund</button>`:'—'}</td></tr>`;}).join('');
    $('#page-content').innerHTML=`${pageHeading('payments',`<button class="button button-primary" data-action="new-payment">＋ Record payment</button>`)}<section class="stats-grid">${stat('Payments recorded',paymentsList.length,'◉','all time')}${stat('Completed payments',money(paymentsList.filter(p=>unwrap(p).status==='completed').reduce((sum,p)=>sum+Number(unwrap(p).amount||0),0)),'↗','collected')}${stat('Open folios',foliosList.filter(f=>unwrap(f).status==='open').length,'▤','awaiting settlement')}${stat('Refunded',paymentsList.filter(p=>unwrap(p).status==='refunded').length,'↩','recorded')}</section><div class="toolbar"><div class="search-box"><span>⌕</span><input id="search-input" placeholder="Search payment reference or guest" value="${escapeHTML(state.search)}"></div><button class="button button-secondary" data-action="new-folio">＋ New folio</button></div><div class="table-wrap"><table class="data-table"><thead><tr><th>PAYMENT</th><th>GUEST</th><th>DATE</th><th>METHOD</th><th>AMOUNT</th><th>STATUS</th><th></th></tr></thead><tbody>${rows||`<tr><td colspan="7">${empty('No payments recorded yet. Create a guest folio to begin.','◉')}</td></tr>`}</tbody></table></div>`;
  }

  function openModal(title, description, fieldsHTML, submitLabel, onSubmit) {
    const root=$('#modal-root'); root.innerHTML=`<div class="modal-backdrop" data-close-modal><section class="modal-card" role="dialog" aria-modal="true" aria-labelledby="modal-title"><header class="modal-header"><div><h2 id="modal-title">${escapeHTML(title)}</h2><p>${escapeHTML(description)}</p></div><button class="modal-close" type="button" data-close-modal aria-label="Close">×</button></header><form class="modal-form" id="record-form"><div class="form-grid">${fieldsHTML}</div><p id="modal-error" class="form-error" aria-live="polite"></p><div class="form-actions"><button type="button" class="button button-secondary" data-close-modal>Cancel</button><button type="submit" class="button button-primary" id="modal-submit">${escapeHTML(submitLabel)}</button></div></form></section></div>`;
    root._submit=onSubmit; root.querySelector('input,select,textarea')?.focus();
    localizeDOM(root);
  }
  const field=(label,name,{type='text',placeholder='',required=false,options=[],value='',wide=false,step=''}={})=>`<div class="field ${wide?'field-wide':''}"><label for="f-${name}">${escapeHTML(label)}${required?' *':''}</label>${type==='select'?`<select id="f-${name}" name="${name}" ${required?'required':''}>${options.map(([key,text])=>`<option value="${escapeHTML(key)}" ${String(key)===String(value)?'selected':''}>${escapeHTML(text)}</option>`).join('')}</select>`:type==='textarea'?`<textarea id="f-${name}" name="${name}" placeholder="${escapeHTML(placeholder)}" ${required?'required':''}>${escapeHTML(value)}</textarea>`:`<input id="f-${name}" name="${name}" type="${type}" placeholder="${escapeHTML(placeholder)}" value="${escapeHTML(value)}" ${required?'required':''} ${step?`step="${step}"`:''}>`}</div>`;
  function openGuestForm() {
    const html=field('First name','first_name',{required:true,placeholder:'e.g. Maya'})+field('Last name','last_name',{required:true,placeholder:'e.g. Hassan'})+field('Email','email',{type:'email',placeholder:'maya@example.com'})+field('Phone','phone',{type:'tel',placeholder:'+966 5X XXX XXXX'})+field('Nationality','nationality',{placeholder:'Saudi'})+field('Document number','document_number',{placeholder:'Passport or ID'})+field('Guest notes','notes',{type:'textarea',wide:true,placeholder:'Preferences or special requirements'});
    openModal('Add a guest','Create a guest profile for bookings and reception.',html,'Save guest',async values=>{await post('/guests',values);toast('Guest profile created.');await renderPage();});
  }
  async function openBookingForm(guestId='') {
    const [guestsList,roomsList]=await Promise.all([list('/guests?limit=200'),list('/rooms?limit=200')]);
    state.cache.guests=guestsList;state.cache.rooms=roomsList;
    const activeRooms=roomsList.filter(room=>unwrap(room).status!=='out_of_order');
    if (!guestsList.length) { toast('Add a guest before creating a booking.',true); navigate('guests');openGuestForm();return; }
    if (!activeRooms.length) { toast('Add a room before creating a booking.',true);navigate('rooms');return; }
    const guestOpts=guestsList.map(row=>[row.id,`${unwrap(row).first_name} ${unwrap(row).last_name}`]);
    const roomOpts=activeRooms.map(row=>[row.id,`Room ${unwrap(row).number} · ${unwrap(row).room_type} · ${money(unwrap(row).rate)}/night`]);
    const html=field('Guest','guest_id',{type:'select',options:guestOpts,required:true,value:guestId,wide:true})+field('Room','room_id',{type:'select',options:roomOpts,required:true,wide:true})+field('Check-in','check_in',{type:'date',required:true,value:today()})+field('Check-out','check_out',{type:'date',required:true})+field('Adults','adults',{type:'number',required:true,value:1})+field('Children','children',{type:'number',value:0})+field('Nightly rate (SAR)','rate',{type:'number',step:'0.01',value:unwrap(activeRooms[0]).rate||0})+field('Booking source','source',{type:'select',value:'direct',options:[['direct','Direct'],['phone','Phone'],['website','Website'],['agency','Travel agency'],['other','Other']]})+field('Special requests','special_requests',{type:'textarea',wide:true,placeholder:'Anything we should prepare for this stay?'})+`<p class="form-hint field-wide">Room availability is checked when the booking is saved.</p>`;
    openModal('Create a booking','Add a stay for a guest. Check-out must be after check-in.',html,'Create booking',async values=>{values.adults=Number(values.adults||1);values.children=Number(values.children||0);values.rate=Number(values.rate||0);await post('/bookings',values);toast('Booking created.');await renderPage();});
  }
  function openRoomForm() {
    const html=field('Room number','number',{required:true,placeholder:'e.g. 204'})+field('Room type','room_type',{required:true,placeholder:'Deluxe king'})+field('Floor','floor',{type:'number',value:1})+field('Nightly rate (SAR)','rate',{type:'number',step:'0.01',value:0})+field('Initial status','status',{type:'select',value:'available',options:[['available','Available'],['needs_cleaning','Needs cleaning'],['out_of_order','Out of order']]})+field('Notes','notes',{type:'textarea',wide:true,placeholder:'Room features or notes'});
    openModal('Add a room','Add this room to your hotel inventory.',html,'Save room',async values=>{values.floor=Number(values.floor||1);values.rate=Number(values.rate||0);await post('/rooms',values);toast('Room added.');await renderPage();});
  }
  async function openFolioForm() {
    const [guestsList,bookingsList]=await Promise.all([list('/guests?limit=200'),list('/bookings?limit=200')]);
    if (!guestsList.length) {toast('Add a guest before opening a folio.',true);navigate('guests');openGuestForm();return;}
    const guestOpts=guestsList.map(row=>[row.id,`${unwrap(row).first_name} ${unwrap(row).last_name}`]);
    const bookingOpts=[['','No booking attached'],...bookingsList.map(row=>[row.id,`#${row.id.slice(0,8)} · ${unwrap(row).status}`])];
    const html=field('Guest','guest_id',{type:'select',options:guestOpts,required:true,wide:true})+field('Booking','reservation_id',{type:'select',options:bookingOpts,wide:true})+field('Folio total (SAR)','total_amount',{type:'number',step:'0.01',value:0,required:true})+field('Currency','currency',{value:'SAR',required:true})+field('Notes','notes',{type:'textarea',wide:true,placeholder:'Room charges, dining, or other services'});
    openModal('Open a guest folio','Create the bill that payments will be posted against.',html,'Create folio',async values=>{if(!values.reservation_id)delete values.reservation_id;values.total_amount=Number(values.total_amount||0);await post('/folios',values);toast('Guest folio created.');await renderPage();});
  }
  async function openPaymentForm() {
    const [foliosList,guestsList]=await Promise.all([list('/folios?limit=200'),list('/guests?limit=200')]);
    state.cache.folios=foliosList;state.cache.guests=guestsList;
    const guestMap=new Map(guestsList.map(g=>[g.id,`${unwrap(g).first_name} ${unwrap(g).last_name}`]));
    const openFolios=foliosList.filter(f=>unwrap(f).status==='open');
    if(!openFolios.length){toast('Create an open folio before recording a payment.',true);navigate('payments');await openFolioForm();return;}
    const options=openFolios.map(row=>{const d=unwrap(row);return[row.id,`${guestMap.get(d.guest_id)||'Guest'} · ${money(d.total_amount)} · ${d.currency||'SAR'}`]});
    const html=field('Guest folio','folio_id',{type:'select',options,required:true,wide:true})+field('Amount (SAR)','amount',{type:'number',step:'0.01',required:true,placeholder:'0.00'})+field('Payment method','method',{type:'select',required:true,options:[['cash','Cash'],['card','Card'],['bank_transfer','Bank transfer'],['online','Online'],['other','Other']]})+field('Reference','reference',{placeholder:'Receipt or transaction ID'});
    openModal('Record a payment','Post a payment to an open guest folio.',html,'Save payment',async values=>{values.amount=Number(values.amount);await post('/payments',values);toast('Payment recorded.');await renderPage();});
  }

  function toast(message,error=false) {
    const root=$('#toast-root');const node=document.createElement('div');node.className=`toast${error?' error':''}`;node.textContent=translateText(message);root.append(node);setTimeout(()=>node.remove(),3500);
  }
  async function runAction(button) {
    const id=button.dataset.id;const action=button.dataset.action;
    if(action==='check-in'||action==='check-out'){
      const endpoint=state.page==='reception'?`/reception/${id}/${action}`:`/bookings/${id}/${action}`;
      await post(endpoint,{});toast(action==='check-in'?'Guest checked in.':'Guest checked out. A housekeeping task was opened.');await renderPage();return;
    }
    if(action==='confirm-booking'){
      if(!window.confirm('Approve this booking request?'))return;await post(`/bookings/${id}/confirm`,{});toast('Booking request approved.');await renderPage();return;
    }
    if(action==='cancel-booking'){
      if(!window.confirm('Cancel this booking?'))return;await post(`/bookings/${id}/cancel`,{});toast('Booking cancelled.');await renderPage();return;
    }
    if(action==='new-booking'){await openBookingForm();return;}
    if(action==='new-guest'){openGuestForm();return;}
    if(action==='new-room'){openRoomForm();return;}
    if(action==='new-payment'){await openPaymentForm();return;}
    if(action==='new-folio'){await openFolioForm();return;}
    if(action==='guest-booking'){await openBookingForm(id);return;}
    if(action==='room-status'){
      const current=button.dataset.status;
      const options=[['available','Available'],['occupied','Occupied'],['needs_cleaning','Needs cleaning'],['out_of_order','Out of order']];
      openModal('Update room status','Choose the current operational status for this room.',field('Room status','status',{type:'select',value:current,options,required:true,wide:true}),'Save status',async values=>{await patch(`/rooms/${id}`,values);toast('Room status updated.');await renderPage();});return;
    }
    if(action==='refund-payment'){
      if(!window.confirm('Mark this payment as refunded?'))return;await patch(`/payments/${id}`,{status:'refunded'});toast('Payment marked as refunded.');await renderPage();return;
    }
    if(action==='refresh')await renderPage();
  }

  document.addEventListener('submit',async event=>{
    if(event.target.id==='login-form'){
      event.preventDefault();await signIn(new FormData(event.target).get('token'));return;
    }
    if(event.target.id==='record-form'){
      event.preventDefault();const submit=$('#modal-submit');submit.disabled=true;submit.textContent=translateText('Saving…');$('#modal-error').textContent='';
      const values=Object.fromEntries(new FormData(event.target).entries());
      try{await $('#modal-root')._submit(values);$('#modal-root').innerHTML='';}
      catch(error){$('#modal-error').textContent=translateText(error.message);submit.disabled=false;submit.textContent=translateText('Try again');}
    }
  });
  document.addEventListener('click',async event=>{
    const nav=event.target.closest('[data-page]');if(nav){navigate(nav.dataset.page);return;}
    const tab=event.target.closest('[data-reception-tab]');if(tab){state.receptionTab=tab.dataset.receptionTab;await reception().then(()=>localizeDOM($('#page-content'))).catch(failPage);return;}
    if(event.target.closest('#menu-toggle')){$('#sidebar').classList.toggle('open');return;}
    if(event.target.closest('#refresh-button')){await renderPage();return;}
    if(event.target.closest('#sign-out')){localStorage.removeItem(TOKEN_KEY);state.token='';showLogin();return;}
    const close=event.target.closest('[data-close-modal]');
    if(close){
      if(close.classList.contains('modal-backdrop')&&event.target!==close)return;
      $('#modal-root').innerHTML='';return;
    }
    const button=event.target.closest('[data-action]');if(button){try{await runAction(button);}catch(error){toast(error.message,true);}}
  });
  document.addEventListener('input',event=>{
    if(event.target.id==='search-input'){
      state.search=event.target.value.toLowerCase();
      document.querySelectorAll('.data-table tbody tr').forEach(row=>{row.hidden=!row.textContent.toLowerCase().includes(state.search);});
      document.querySelectorAll('.room-card').forEach(card=>{card.hidden=!card.textContent.toLowerCase().includes(state.search);});
    }
  });
  document.addEventListener('change',event=>{
    if(event.target.id==='queue-date'){state.date=event.target.value;reception().catch(failPage);}
  });
  document.addEventListener('keydown',event=>{if(event.key==='Escape')$('#modal-root').innerHTML='';});

  if(state.token) showApp(); else showLogin();
})();
