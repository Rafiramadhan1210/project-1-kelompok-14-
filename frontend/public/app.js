// Toggle Profile Dropdown
const profileToggle = document.getElementById('profile-toggle');
const profileDropdown = document.getElementById('profile-dropdown');
const mobileMenuToggle = document.getElementById('mobile-menu-toggle');
const mobileMenu = document.getElementById('mobile-menu');

if (profileToggle) {
    profileToggle.addEventListener('click', (e) => {
        e.stopPropagation();
        profileDropdown.classList.toggle('hidden');
    });
}

document.addEventListener('click', () => {
    if (profileDropdown) profileDropdown.classList.add('hidden');
});

if (profileDropdown) {
    profileDropdown.addEventListener('click', (e) => {
        e.stopPropagation();
    });
}

// Toggle Mobile Menu
mobileMenuToggle.addEventListener('click', () => {
    mobileMenu.classList.toggle('hidden');
});

// Cek status login & sesuaikan tampilan navbar
async function checkLoginStatus() {
    const loginLink = document.getElementById('login-link');
    const profileMenu = document.getElementById('profile-menu');
    const profileName = document.getElementById('profile-name');
    const profileDropdownName = document.getElementById('profile-dropdown-name');
    const profileDropdownEmail = document.getElementById('profile-dropdown-email');

    if (!loginLink || !profileMenu) return;

    try {
        const res = await fetch('/api/me', { credentials: 'include' });
        if (res.ok) {
            const result = await res.json();
            const user = result.user || {};
            if (profileName) profileName.textContent = user.nama || 'Akun';
            if (profileDropdownName) profileDropdownName.textContent = user.nama || 'Pengguna';
            if (profileDropdownEmail) profileDropdownEmail.textContent = user.email || '';
            profileMenu.classList.remove('hidden');
            loginLink.classList.add('hidden');
        } else {
            profileMenu.classList.add('hidden');
            loginLink.classList.remove('hidden');
        }
    } catch (err) {
        profileMenu.classList.add('hidden');
        loginLink.classList.remove('hidden');
    }
}
checkLoginStatus();

// Logout Button
const logoutBtn = document.getElementById('logout-btn');
if (logoutBtn) {
    logoutBtn.addEventListener('click', async () => {
        try {
            await fetch('/logout', { method: 'POST', credentials: 'include' });
        } catch (err) {
            // tetap arahkan ke login walau request gagal
        }
        window.location.href = 'login.html';
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
    <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden hover:shadow-lg transition hover:-translate-y-1" data-kategori="${kategori}" data-id="${item._id}">
        <div class="relative h-48 overflow-hidden">
            <img src="${item.gambar || 'https://images.unsplash.com/photo-1507525428034-b723cf961d3e'}" alt="${item.nama}" class="w-full h-full object-cover">
            <span class="absolute bottom-4 left-4 text-xs font-bold text-white bg-black/50 backdrop-blur px-3 py-1 rounded-full">${item.kategori || 'Wisata'}</span>
            <button class="wishlist-btn absolute top-3 right-3 w-9 h-9 rounded-full bg-white/90 flex items-center justify-center hover:bg-white transition" data-id="${item._id}">
                <i class="fa-regular fa-heart text-gray-600"></i>
            </button>
        </div>
        <div class="p-5 cursor-pointer" data-open-detail="${item._id}">
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

// Wishlist: load status awal & pasang event listener tombol hati
let myWishlistIds = new Set();

async function loadMyWishlistIds() {
    try {
        const res = await fetch('/api/my-wishlist', { credentials: 'include' });
        if (!res.ok) return; // belum login, biarkan semua hati kosong
        const result = await res.json();
        myWishlistIds = new Set((result.data || []).map(d => d._id));
        refreshWishlistIcons();
    } catch (err) {
        // diam saja kalau gagal, hati tetap kosong
    }
}

function refreshWishlistIcons() {
    document.querySelectorAll('.wishlist-btn').forEach(btn => {
        const icon = btn.querySelector('i');
        if (myWishlistIds.has(btn.dataset.id)) {
            icon.classList.remove('fa-regular');
            icon.classList.add('fa-solid', 'text-red-500');
        } else {
            icon.classList.remove('fa-solid', 'text-red-500');
            icon.classList.add('fa-regular', 'text-gray-600');
        }
    });
}

document.addEventListener('click', async (e) => {
    const btn = e.target.closest('.wishlist-btn');
    if (!btn) return;
    e.stopPropagation();

    const destinationId = btn.dataset.id;
    try {
        const res = await fetch('/api/wishlist/toggle', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ destination_id: destinationId })
        });

        if (res.status === 401) {
            alert('Silakan login dulu untuk menyimpan wishlist.');
            return;
        }

        const result = await res.json();
        if (result.wishlisted) {
            myWishlistIds.add(destinationId);
        } else {
            myWishlistIds.delete(destinationId);
        }
        refreshWishlistIcons();
    } catch (err) {
        console.error('Gagal update wishlist:', err);
    }
});

loadMyWishlistIds();

// Klik body kartu (bukan tombol hati/booking) untuk buka halaman detail
document.addEventListener('click', (e) => {
    if (e.target.closest('.wishlist-btn')) return; // jangan buka detail kalau klik tombol hati
    if (e.target.closest('button')) return; // jangan buka detail kalau klik tombol Booking di kartu

    const detailTrigger = e.target.closest('[data-open-detail]');
    if (detailTrigger) {
        const id = detailTrigger.dataset.openDetail;
        window.location.href = `destinasi-detail.html?id=${id}`;
    }
});

async function loadKategori() {
    const response = await fetch('/api/kategori');
    const data = await response.json();
    
    const container = document.getElementById('kategori-container');
    container.innerHTML = data.map(k => `
        <button class="kategori-item" onclick="filterDestinasi('${k.nama}')">
            ${k.nama}
        </button>
    `).join('');
}

