const postcss = require('postcss');
const tailwindcss = require('tailwindcss');
const fs = require('fs');
const path = require('path');

const inputPath = path.join(__dirname, 'public', 'input.css');
const outputPath = path.join(__dirname, 'public', 'output.css');
const configPath = path.join(__dirname, 'tailwind.config.js');

const css = fs.readFileSync(inputPath, 'utf8');

postcss([tailwindcss(configPath)])
    .process(css, { from: inputPath, to: outputPath })
    .then((result) => {
        fs.writeFileSync(outputPath, result.css);
        console.log('Build CSS berhasil! Disimpan ke', outputPath);
    })
    .catch((err) => {
        console.error('Build CSS gagal:', err);
        process.exit(1);
    });