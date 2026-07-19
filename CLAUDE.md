# Instruction — GoTrip (Kelompok 14, D4 Teknik Informatika ULBI)


**Project:** GoTrip — reservasi tiket destinasi wisata digital berbasis web
**Tim:** Kelompok 14 — Rafi Ramadhan (714250023), Shany Aulia Prastica (714250009)
**Institusi:** D4 Teknik Informatika, ULBI
**Jenis laporan:** Capstone Project — **Proyek 1** (profil target: **Debugger**)

---

# CORE — Integritas Akademik Penulisan Ilmiah

Berlaku untuk semua jenis naskah akademik (artikel jurnal, laporan, proposal, disertasi, tesis).

## 1. Integritas sitasi & bibliografi
- Setiap sitasi HARUS punya entri di file bibliografi; 0 "Citation undefined" di log.
- Setiap entri bibliografi diverifikasi ke API otoritatif (CrossRef untuk DOI penerbit,
  DataCite untuk DOI arXiv) sebelum dianggap final -- jangan isi dari ingatan.
- Dilarang keras sitasi fabrikasi. Jika metadata tak terverifikasi, tandai & laporkan ke user.
- Butuh literatur pendukung yang belum ada -> TANYA USER (format `[TODO-CITE-<n>]`, bilingual
  EN+ID, batch di akhir pengerjaan, bukan menetes satu-satu). Titik rawan fabrikasi:
  Latar Belakang, Tinjauan Pustaka, Metodologi, Pembahasan, Keterbatasan.

## 2. Integritas klaim & data -- klaim HARUS sesuai bukti
- Cek kelengkapan bahan (dokumen, data, kode, eksperimen, hasil) SEBELUM menulis. Bila bahan
  yang diperlukan belum ada, JANGAN mengarang -- laporkan ke user, tandai bagian itu belum
  bisa ditulis.
- Setiap angka/klaim hasil harus tertelusur ke sumbernya (skrip, data, commit).
- Deskripsi metode di teks HARUS persis sama dengan yang benar-benar dieksekusi di kode.
- Dilarang angka hardcoded di skrip plot/analisis. Ragu validitas angka -> tandai
  INVALID/quarantine, laporkan ke user, jangan diteruskan ke naskah.

## 3. Prinsip argumentasi: ketegangan jadi rigor
- Jangan sembunyikan bukti yang tampak menentang usulan -- bingkai sebagai diagnosis ->
  generalisasi, tundukkan pada bar metodologis yang sama.
- Novelty = insight/generalitas, bukan sekadar merakit blok yang sudah ada.
- Jangan overclaim. Batasan yang jujur menaikkan nilai, bukan menurunkan.

## 4. Gaya tulisan -- hindari "AI-tell"
- Dilarang em-dash; pakai koma/kurung/titik dua/pecah kalimat. En-dash untuk rentang angka
  tetap boleh.
- Hindari pengulangan kata promosi: robust, novel, comprehensive, crucial, powerful,
  lightweight, seamless, reproducible.
- Hindari pembuka klausa formulaik: "Selain itu,", "Lebih lanjut,", "Penting untuk dicatat,",
  "Secara keseluruhan,".
- Konkret > generik. Variasikan ritme kalimat. Jangan sitir path/nama file repo di prosa
  naskah (kecuali entry-point reproduksi).

## 5. Istilah wajar Bahasa Indonesia
- Pakai istilah baku bila lazim ("pembelajaran mesin"); pertahankan istilah asing dicetak
  miring bila belum lazim (deep learning, overfitting).
- Konsisten satu istilah untuk satu konsep di seluruh naskah.
- Hindari terjemahan harfiah kaku ("Dalam rangka untuk" -> "untuk"; "Hal ini dikarenakan
  oleh" -> "karena").
- Ejaan sesuai KBBI/PUEBI: "analisis" bukan "analisa", "metode" bukan "metoda".
- Utamakan kalimat aktif, bukan pasif berlapis khas terjemahan mesin.

## 6. Figure & tabel -- standar publikasi
- Setiap gambar hasil digenerate dari data nyata via skrip tertelusur -- dilarang menggambar
  ulang/menghias angka manual.
- Caption self-contained; sumbu-y dipotong ditandai jelas; error bar didefinisikan.
- Tabel bergaya booktabs (tanpa garis vertikal); desimal konsisten; penanda "terbaik"
  didefinisikan di caption.
- Dilarang diagram/ilustrasi hasil AI generatif untuk diagram metode -- pakai tool diagram
  nyata.

## 7. Teknis LaTeX (bila naskah berbasis LaTeX)
- Build bersih menyeluruh: 0 undefined citation & cross-reference, 0 gambar hilang.
- Rujuk semua elemen via `\ref`/`\cref`, jangan manual; `\label` setelah `\caption`.
- Math mode untuk simbol; tipografi LaTeX benar (kutip pintar, `\%`, en-dash).

## 8. Checklist inti sebelum menyatakan selesai
Ringkasnya: bibliografi terverifikasi & bersih, tidak ada `TODO-CITE` tersisa, kelengkapan
bahan dicek di awal, semua angka/klaim tertelusur ke bukti nyata, bebas AI-tell, istilah
Indonesia wajar, figure/tabel sesuai standar, tidak overclaim.

---

# OVERLAY -- Laporan Capstone Project (Vokasi)

Melengkapi core di atas; khusus laporan capstone project mahasiswa vokasi D3/D4/Sarjana Terapan.

## C1. Model capstone bertahap -- profil per tahun (WAJIB ditentukan lebih dulu)

| Proyek | Profil target | Bukti wajib |
|---|---|---|
| **Proyek 1** (tahun ke-1) | **Debugger** | Bug nyata: gejala -> reproduksi -> diagnosis -> akar masalah -> perbaikan -> verifikasi; SCM (Git/GitHub); workflow branch/PR/issue; CI/CD; vibe coding (AI-assisted, disadari & diverifikasi) |
| Proyek 2 (tahun ke-2) | Software Engineer | Fitur end-to-end + pengujian (unit/integrasi) tertulis, CI menjalankan test otomatis |
| Proyek 3 (tahun ke-3) | Software Architect | Arsitektur sistem, integrasi API/M2M, kepemimpinan tim |

**Project GoTrip ini adalah Proyek 1 -> profil Debugger.** Laporan WAJIB memuat studi kasus
bug nyata (bukan cuma cerita fitur), bukti SCM/CI-CD, dan disclosure vibe coding.
Sudah ditambahkan sebagai Bab VII (`BAB_VII_Bukti_Debugging_dan_Perkakas.docx`).

## C2. Struktur mengikuti template prodi (WAJIB)
Ikuti template resmi ULBI secara harfiah: Kata Pengantar -> Lembar Persetujuan -> Lembar
Pengesahan -> Daftar Isi -> Bab I-VI (+ Bab VII bukti Debugger untuk Proyek 1).
Rumusan masalah -> tujuan -> requirement -> solusi -> pengujian saling mengunci satu-ke-satu.

## C3. Fokus proyek terapan & demonstrasi kompetensi
- Masalah nyata & spesifik (antrean fisik loket wisata lokal) -- sudah ada di Bab I.
- Requirement terukur (fungsional/non-fungsional) -- sudah ada di Bab II.
- Justifikasi teknologi bersitasi bila memungkinkan.
- Kontribusi tiap anggota tim eksplisit & jujur (belum ada di laporan -- pertimbangkan
  ditambahkan bila diminta pembagian tugas).

## C4. Pengujian, uji terima & luaran
- Tabel requirement -> metode uji -> hasil (lolos/tidak) -> bukti -- belum ada di laporan
  saat ini; pertimbangkan ditambahkan sebagai pelengkap Bab V/VI.
- Bukti luaran dilampirkan (screenshot, tautan repo, demo), bukan hanya diklaim.

## C5. Integritas karya mahasiswa & bimbingan
- Kode/aset pihak ketiga (library, template) diatribusi + patuhi lisensi.
- Kode hasil AI (vibe coding) dicantumkan statusnya -- sudah dijelaskan di Bab VII (7.4).
- Setiap klaim harus bisa dipertanggungjawabkan ke pembimbing.

## C6. Persiapan sidang & demo
- Demo aplikasi harus benar-benar berfungsi (bukan hanya slide).
- Kuasai alasan tiap keputusan desain, requirement yang belum terpenuhi, kontribusi
  masing-masing anggota.
- Untuk Proyek 1: siap menunjukkan langsung cara mendiagnosis bug (Bab VII).

## C7. Checklist capstone -- status project GoTrip
- [x] Bahan bukti profil Proyek 1 tersedia -> dilengkapi di Bab VII
- [x] Format sesuai template ULBI (cover, kata pengantar, lembar pengesahan, dst)
- [x] Masalah nyata & requirement terukur (Bab I-II)
- [x] Kontribusi tiap anggota tim eksplisit -- ditambahkan di Bab V (5.5)
- [x] Tabel uji terima requirement -> hasil -> bukti -- ditambahkan di Bab V (5.4)
- [x] Bukti SCM/CI-CD konkret (tautan repo, screenshot commit/Actions) -- diverifikasi sesuai repositori dan dilampirkan di Bab VII (7.2 dan 7.3)
- [x] Disclosure vibe coding -- sudah ada di Bab VII (7.4)
- [ ] Similarity check & atribusi lisensi pihak ketiga -- belum dicek, disarankan sebelum sidang

---

*Referensi sumber asli: https://ll.my.id/instruction/manuscript/ (core) dan
https://ll.my.id/instruction/manuscript/capstone.html (overlay).*
