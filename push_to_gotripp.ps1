# Script untuk mengunggah backend dan frontend secara terpisah ke organisasi GoTripp (branch: dev-shany)

$PSScriptRoot = Split-Path -Parent -Path $MyInvocation.MyCommand.Definition

# 1. Proses untuk Backend
Write-Host "--------------------------------------------------" -ForegroundColor Yellow
Write-Host "Memulai proses unggah BACKEND..." -ForegroundColor Cyan
Write-Host "--------------------------------------------------" -ForegroundColor Yellow

cd "$PSScriptRoot\backend"

# Inisialisasi git jika belum ada
if (-not (Test-Path ".git")) {
    git init
}

git add .
git commit -m "Initial commit backend GoTrip"
git branch -M dev-shany

# Hapus remote origin lama jika ada, lalu tambahkan yang baru
git remote remove origin 2>$null
git remote add origin https://github.com/GoTripp/backend.git

Write-Host "Mengirim file backend ke https://github.com/GoTripp/backend (branch: dev-shany)..." -ForegroundColor Cyan
git push -u origin dev-shany

# 2. Proses untuk Frontend
Write-Host "`n--------------------------------------------------" -ForegroundColor Yellow
Write-Host "Memulai proses unggah FRONTEND..." -ForegroundColor Cyan
Write-Host "--------------------------------------------------" -ForegroundColor Yellow

cd "$PSScriptRoot\frontend"

# Inisialisasi git jika belum ada
if (-not (Test-Path ".git")) {
    git init
}

git add .
git commit -m "Initial commit frontend GoTrip"
git branch -M dev-shany

# Hapus remote origin lama jika ada, lalu tambahkan yang baru
git remote remove origin 2>$null
git remote add origin https://github.com/GoTripp/frontend.git

Write-Host "Mengirim file frontend ke https://github.com/GoTripp/frontend (branch: dev-shany)..." -ForegroundColor Cyan
git push -u origin dev-shany

Write-Host "`nProses Selesai!" -ForegroundColor Green
Read-Host "Tekan Enter untuk keluar..."
