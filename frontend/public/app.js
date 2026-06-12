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

// Filter Kategori Klik Handler
let currentFilter = 'Semua';
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
            return; // Biarkan browser membuka tentang.html secara langsung
        }
        // Menangani scroll halus ke section id jika menggunakan # href
        if (link.getAttribute('href').startsWith('#')) {
            e.preventDefault();
            const targetId = link.getAttribute('href').substring(1);
            const targetSection = document.getElementById(targetId);
            if (targetSection) {
                targetSection.scrollIntoView({ behavior: 'smooth' });
            }
        }
        console.log('Menu clicked:', link.dataset.menu);
    });
});

// Filter Destinasi Function (Sembunyikan / Tampilkan Card di Halaman)
function filterDestinasi(filter) {
    const items = document.querySelectorAll('[data-kategori]');
    items.forEach(item => {
        if (filter === 'Semua' || item.dataset.kategori === filter) {
            item.classList.remove('hidden');
        } else {
            item.classList.add('hidden');
        }
    });
}

// Load Destinasi dari API Backend
fetch('/button')
    .then(res => res.json())
    .then(result => {
        const list = document.getElementById('destinasi-list');
        list.innerHTML = ''; // Kosongkan penampung data awal

        result.data.forEach(item => {
            // Menggunakan huruf kecil sesuai struktur database MongoDB Atlas
            const kategori = item.kategori || 'Semua';

            list.innerHTML += `
                <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden hover:shadow-lg transition hover:-translate-y-1" data-kategori="${kategori}">
                    <div class="relative h-48 overflow-hidden">
                        <img src="${item.gambar || 'https://images.unsplash.com/photo-1507525428034-b723cf961d3e'}" alt="${item.nama}" class="w-full h-full object-cover">
                        <span class="absolute bottom-4 left-4 text-xs font-bold text-white bg-black/50 backdrop-blur px-3 py-1 rounded-full">${item.kategori || 'Wisata'}</span>
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