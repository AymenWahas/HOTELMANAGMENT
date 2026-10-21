(() => {
  const form = document.querySelector('#booking-form');
  const checkIn = document.querySelector('#check-in');
  const checkOut = document.querySelector('#check-out');
  const availability = document.querySelector('#availability');
  const guestFields = document.querySelector('#guest-fields');
  const message = document.querySelector('#booking-message');
  const submitButton = document.querySelector('#submit-booking');
  const galleryCards = Array.from(document.querySelectorAll('.gallery-card'));
  const galleryImages = galleryCards.map((card) => card.querySelector('img'));
  const photoDialog = document.querySelector('#photo-dialog');
  const dialogImage = document.querySelector('#dialog-image');
  const dialogTitle = document.querySelector('#dialog-title');
  const dialogCount = document.querySelector('#dialog-count');
  let activePhoto = 0;
  let selectedRoom = null;
  const localDate = (date) => {
    const offset = date.getTimezoneOffset() * 60000;
    return new Date(date.getTime() - offset).toISOString().slice(0, 10);
  };
  const today = localDate(new Date());
  checkIn.min = today;
  checkOut.min = today;
  const year = document.querySelector('#year');
  if (year) year.textContent = new Date().getFullYear();

  function showPhoto(index) {
    activePhoto = (index + galleryImages.length) % galleryImages.length;
    const source = galleryImages[activePhoto];
    dialogImage.src = source.src;
    dialogImage.alt = source.alt;
    dialogTitle.textContent = galleryCards[activePhoto].querySelector('strong').textContent;
    dialogCount.textContent = String(activePhoto + 1).padStart(2, '0') + ' / ' + String(galleryImages.length).padStart(2, '0');
  }
  galleryCards.forEach((card, index) => {
    card.addEventListener('click', () => {
      showPhoto(index);
      photoDialog.showModal();
    });
    card.addEventListener('keydown', (event) => {
      if (event.key === 'Enter' || event.key === ' ') {
        event.preventDefault();
        showPhoto(index);
        photoDialog.showModal();
      }
    });
  });
  document.querySelector('.photo-close').addEventListener('click', () => photoDialog.close());
  document.querySelector('.photo-prev').addEventListener('click', () => showPhoto(activePhoto - 1));
  document.querySelector('.photo-next').addEventListener('click', () => showPhoto(activePhoto + 1));
  photoDialog.addEventListener('click', (event) => {
    if (event.target === photoDialog) photoDialog.close();
  });
  photoDialog.addEventListener('keydown', (event) => {
    if (event.key === 'ArrowLeft') showPhoto(activePhoto + 1);
    if (event.key === 'ArrowRight') showPhoto(activePhoto - 1);
  });

  function clearSelection() {
    selectedRoom = null;
    guestFields.hidden = true;
    submitButton.disabled = true;
  }
  function showError(text) {
    message.className = 'booking-message error';
    message.textContent = text;
  }
  checkIn.addEventListener('change', () => {
    const nextDay = new Date(checkIn.value + 'T12:00:00');
    nextDay.setDate(nextDay.getDate() + 1);
    checkOut.min = localDate(nextDay);
    if (checkOut.value && checkOut.value <= checkIn.value) checkOut.value = '';
    availability.replaceChildren();
    message.textContent = '';
    clearSelection();
  });
  checkOut.addEventListener('change', () => {
    availability.replaceChildren();
    message.textContent = '';
    clearSelection();
  });

  document.querySelector('#search-rooms').addEventListener('click', async () => {
    message.className = 'booking-message';
    message.textContent = '';
    availability.replaceChildren();
    clearSelection();
    if (!checkIn.value || !checkOut.value || checkOut.value <= checkIn.value) {
      showError('اختر تاريخ وصول ومغادرة صحيحين أولًا.');
      return;
    }
    const status = document.createElement('p');
    status.className = 'availability-status';
    status.textContent = 'جارٍ البحث عن الغرف…';
    availability.append(status);
    try {
      const params = new URLSearchParams({ check_in: checkIn.value, check_out: checkOut.value });
      const response = await fetch('/api/v1/public/rooms?' + params.toString());
      const result = await response.json();
      if (!response.ok) throw new Error(result.error || 'تعذر البحث عن الغرف.');
      availability.replaceChildren();
      if (!result.items || !result.items.length) {
        const empty = document.createElement('p');
        empty.className = 'availability-status';
        empty.textContent = 'لا توجد غرف متاحة لهذه التواريخ. جرّب تواريخ أخرى.';
        availability.append(empty);
        return;
      }
      const title = document.createElement('p');
      title.className = 'availability-title';
      title.textContent = 'اختر الغرفة المناسبة';
      availability.append(title);
      result.items.forEach((room, index) => {
        const card = document.createElement('article');
        card.className = 'available-room';
        const thumbnail = document.createElement('img');
        thumbnail.className = 'available-room-image';
        thumbnail.src = galleryImages[index % galleryImages.length].src;
        thumbnail.alt = galleryImages[index % galleryImages.length].alt;
        thumbnail.loading = 'lazy';
        const info = document.createElement('div');
        info.className = 'available-room-info';
        const name = document.createElement('strong');
        name.textContent = room.room_type || 'غرفة';
        const details = document.createElement('small');
        details.textContent = 'سعر الليلة';
        const price = document.createElement('span');
        price.className = 'room-price';
        price.textContent = Number(room.rate || 0).toLocaleString('ar-SA') + ' ر.س';
        const choose = document.createElement('button');
        choose.type = 'button';
        choose.textContent = 'اختيار هذه الغرفة';
        choose.addEventListener('click', () => {
          availability.querySelectorAll('.available-room').forEach((item) => item.classList.remove('selected'));
          card.classList.add('selected');
          availability.querySelectorAll('.available-room button').forEach((button) => {
            button.textContent = 'اختيار هذه الغرفة';
            button.setAttribute('aria-pressed', 'false');
          });
          choose.textContent = 'تم اختيار الغرفة ✓';
          choose.setAttribute('aria-pressed', 'true');
          selectedRoom = room;
          guestFields.hidden = false;
          submitButton.disabled = false;
          guestFields.scrollIntoView({ behavior: 'smooth', block: 'center' });
        });
        info.append(name, details);
        card.append(thumbnail, info, price, choose);
        availability.append(card);
      });
    } catch (error) {
      availability.replaceChildren();
      showError(error.message);
    }
  });

  form.addEventListener('submit', async (event) => {
    event.preventDefault();
    if (!selectedRoom) return;
    message.className = 'booking-message';
    message.textContent = 'جارٍ إرسال طلبك…';
    submitButton.disabled = true;
    const data = new FormData(form);
    const parts = String(data.get('name') || '').trim().split(/\s+/);
    const payload = {
      first_name: parts.shift() || '',
      last_name: parts.join(' '),
      email: String(data.get('email') || '').trim(),
      phone: String(data.get('phone') || '').trim(),
      room_id: selectedRoom.id,
      check_in: String(data.get('check_in')),
      check_out: String(data.get('check_out')),
      adults: Number(data.get('adults')),
      children: Number(data.get('children')),
      website: String(data.get('website') || ''),
    };
    try {
      const response = await fetch('/api/v1/public/booking-requests', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });
      const result = await response.json();
      if (!response.ok) throw new Error(result.error || 'تعذر إرسال طلب الحجز.');
      form.reset();
      checkIn.min = today;
      checkOut.min = today;
      availability.replaceChildren();
      clearSelection();
      message.textContent = 'تم استلام طلب الحجز بنجاح. رقم الحجز ' + String(result.id).slice(0, 8).toUpperCase() + '. سيتم مراجعة الطلب من قبل فريق الاستقبال/المدير قبل التأكيد.';
      message.scrollIntoView({ behavior: 'smooth', block: 'center' });
    } catch (error) {
      showError(error.message);
      submitButton.disabled = false;
    }
  });
})();
