// Toggle Profile Dropdown
const profileToggle = document.getElementById('profile-toggle');
const profileDropdown = document.getElementById('profile-dropdown');
const mobileMenuToggle = document.getElementById('mobile-menu-toggle');
const mobileMenu = document.getElementById('mobile-menu');

profileToggle.addEventListener('click', (e) => {
    e.stopPropagation();
    profileDropdown.classList.toggle('hidden');
});

document.addEventListener('click', () => {
    profileDropdown.classList.add('hidden');
});

profileDropdown.addEventListener('click', (e) => {
    e.stopPropagation();
});

// Toggle Mobile Menu
mobileMenuToggle.addEventListener('click', () => {
    mobileMenu.classList.toggle('hidden');
});

// Logout Button
const logoutBtn = document.getElementById('logout-btn');
if (logoutBtn) {
    logoutBtn.addEventListener('click', () => {
        alert('Logging out...');
        // Implementasi logout ke backend
    });
}

// Notification Button
const notificationBtn = document.getElementById('notification-btn');
if (notificationBtn) {
    notificationBtn.addEventListener('click', () => {
        alert('Notifikasi: Belum ada notifikasi baru');
    });
}

// Wishlist Button
const wishlistBtn = document.getElementById('wishlist-btn');
if (wishlistBtn) {
    wishlistBtn.addEventListener('click', () => {
        alert('Anda memiliki 3 item di wishlist');
    });
}

// Filter Kategori
let currentFilter = 'semua';
document.querySelectorAll('[data-filter]').forEach(btn => {
    btn.addEventListener('click', () => {
        document.querySelectorAll('[data-filter]').forEach(b => {
            b.classList.remove('bg-blue-600', 'text-white');
            b.classList.add('bg-gray-100', 'text-gray-700');
        });
        btn.classList.remove('bg-gray-100', 'text-gray-700');
        btn.classList.add('bg-blue-600', 'text-white');
        currentFilter = (btn.dataset.filter || 'semua').toLowerCase();
        console.log('Filter dipilih:', currentFilter);
        filterDestinasi(currentFilter);
    });
});

// Menu Navigation
document.querySelectorAll('[data-menu]').forEach(link => {
    link.addEventListener('click', (e) => {
        if (link.dataset.menu === 'tentang') {
            return; // biarkan browser buka tentang.html
        }

        e.preventDefault();
        console.log('Menu clicked:', link.dataset.menu);
    });
});

// Filter Destinasi Function
function filterDestinasi(filter) {
    const items = document.querySelectorAll('[data-kategori]');
    items.forEach(item => {
        const kategori = (item.dataset.kategori || 'semua').toLowerCase();
        if (filter === 'semua' || kategori === filter) {
            item.classList.remove('hidden');
        } else {
            item.classList.add('hidden');
        }
    });
}

function escapeHTML(value) {
    return String(value ?? '').replace(/[&<>"']/g, (char) => ({
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#39;',
    }[char]));
}

function getImageUrl(gambar) {
    const url = String(gambar || '').trim();
    if (!url) return '';
    if (/^(https?:)?\/\//i.test(url) || url.startsWith('data:') || url.startsWith('/')) {
        return url;
    }
    return `/${url.replace(/^\.?\//, '')}`;
}

function cssUrl(value) {
    return String(value || '').replace(/\\/g, '\\\\').replace(/"/g, '\\"');
}

// Load Destinasi dari API
fetch('/button')
    .then(res => res.json())
    .then(result => {
        const list = document.getElementById('destinasi-list');
        result.data.forEach(item => {
            const kategori = item.kategori?.toLowerCase() || 'semua';
            const destinationId = item._id || item.id || '';
            const imageUrl = getImageUrl(item.gambar);
            const imageStyle = imageUrl
                ? `background-image: url("${cssUrl(imageUrl)}"); background-size: cover; background-position: center;`
                : '';
            const detailUrl = destinationId
                ? `destinasi-detail.html?id=${encodeURIComponent(destinationId)}`
                : 'destinasi-detail.html';
            list.innerHTML += `
                <div class="destination-card bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden hover:shadow-lg transition hover:-translate-y-1" data-kategori="${kategori}">
                    <div class="card-img-wrapper relative h-48 bg-gradient-to-br from-blue-400 to-blue-600 flex items-end justify-start p-4" style="${imageStyle}">
                        ${imageUrl ? `<img src="${escapeHTML(imageUrl)}" alt="${escapeHTML(item.nama || 'Destinasi')}" loading="eager" decoding="async" style="position:absolute;inset:0;width:100%;height:100%;object-fit:cover;display:block;">` : ''}
                        <div style="position:absolute;inset:0;background:linear-gradient(to top,rgba(0,0,0,0.55),rgba(0,0,0,0));"></div>
                        <span class="relative text-xs font-bold text-white bg-black/30 backdrop-blur px-3 py-1 rounded-full" style="z-index:1;">${escapeHTML(item.kategori || 'Wisata')}</span>
                    </div>
                    <div class="p-5">
                        <h3 class="text-xl font-bold text-gray-800">${escapeHTML(item.nama)}</h3>
                        <p class="text-sm text-gray-500 mt-2 line-clamp-2">${escapeHTML(item.deskripsi)}</p>
                        <div class="mt-4 flex items-center gap-2 text-sm text-gray-600 mb-4">
                            <i class="fa-solid fa-location-dot text-red-500"></i>
                            <span>${escapeHTML(item.lokasi || 'Lokasi tidak tersedia')}</span>
                        </div>
                        <div class="mt-6 flex justify-between items-center">
                            <div>
                                <span class="text-orange-500 font-bold text-lg">
                                    Rp ${item.harga ? item.harga.toLocaleString('id-ID') : '0'}
                                </span>
                                <p class="text-xs text-gray-400">/per tiket</p>
                            </div>
                            <a href="${detailUrl}" class="booking-btn bg-blue-600 text-white px-5 py-2 rounded-lg font-semibold hover:bg-blue-700 transition flex items-center gap-2">
                                <i class="fa-solid fa-ticket"></i>Booking
                            </a>
                        </div>
                    </div>
                </div>
            `;
        });
    })
    .catch(err => console.error('Error load destinasi:', err));
