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
        currentFilter = btn.dataset.filter;
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
        if (filter === 'semua' || item.dataset.kategori === filter) {
            item.classList.remove('hidden');
        } else {
            item.classList.add('hidden');
        }
    });
}

// Load Destinasi dari API
fetch('/button')
    .then(res => res.json())
    .then(result => {
        const list = document.getElementById('destinasi-list');
        result.data.forEach(item => {
            const kategori = item.kategori?.toLowerCase() || 'semua';
            list.innerHTML += `
                <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden hover:shadow-lg transition hover:-translate-y-1" data-kategori="${kategori}">
                    <div class="relative h-48 bg-gradient-to-br from-blue-400 to-blue-600 flex items-end justify-start p-4">
                        <span class="text-xs font-bold text-white bg-black/30 backdrop-blur px-3 py-1 rounded-full">${item.kategori || 'Wisata'}</span>
                    </div>
                    <div class="p-5">
                        <h3 class="text-xl font-bold text-gray-800">${item.nama}</h3>
                        <p class="text-sm text-gray-500 mt-2 line-clamp-2">${item.deskripsi}</p>
                        <div class="mt-4 flex items-center gap-2 text-sm text-gray-600 mb-4">
                            <i class="fa-solid fa-location-dot text-red-500"></i>
                            <span>${item.lokasi || 'Lokasi tidak tersedia'}</span>
                        </div>
                        <div class="mt-6 flex justify-between items-center">
                            <div>
                                <span class="text-orange-500 font-bold text-lg">
                                    Rp ${item.harga ? item.harga.toLocaleString('id-ID') : '0'}
                                </span>
                                <p class="text-xs text-gray-400">/per tiket</p>
                            </div>
                            <button class="bg-blue-600 text-white px-5 py-2 rounded-lg font-semibold hover:bg-blue-700 transition flex items-center gap-2">
                                <i class="fa-solid fa-ticket"></i>Booking
                            </button>
                        </div>
                    </div>
                </div>
            `;
        });
    })
    .catch(err => console.error('Error load destinasi:', err));
