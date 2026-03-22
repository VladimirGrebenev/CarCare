#!/usr/bin/env node
// generate-icons.js
// Автоматизация генерации maskable-иконок для PWA CarCare
// Стили: glassmorphism, minimalism; Темы: car, fuel, fines
// Использует svg-to-img для конвертации SVG в PNG

const fs = require('fs');
const path = require('path');
const { createCanvas, loadImage } = require('canvas');

const ICONS = [
  { name: 'car', color: '#1976d2', emoji: '🚗' },
  { name: 'fuel', color: '#43a047', emoji: '⛽' },
  { name: 'fines', color: '#e53935', emoji: '💸' }
];

const STYLES = ['glassmorphism', 'minimalism'];
const SIZES = [192, 512];
const OUT_DIR = path.join(__dirname, 'public', 'icons');

function glassBg(ctx, size, color) {
  // Glassmorphism: полупрозрачный круглый фон
  ctx.save();
  ctx.globalAlpha = 0.7;
  ctx.beginPath();
  ctx.arc(size/2, size/2, size/2-8, 0, 2*Math.PI);
  ctx.fillStyle = color;
  ctx.shadowColor = '#fff';
  ctx.shadowBlur = 16;
  ctx.fill();
  ctx.restore();
}

function minimalBg(ctx, size, color) {
  // Minimalism: простой круглый фон
  ctx.save();
  ctx.beginPath();
  ctx.arc(size/2, size/2, size/2-8, 0, 2*Math.PI);
  ctx.fillStyle = color;
  ctx.fill();
  ctx.restore();
}

function drawEmoji(ctx, emoji, size) {
  ctx.font = `${size*0.6}px serif`;
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';
  ctx.fillText(emoji, size/2, size/2);
}

function generateIcon({ name, color, emoji }, style, size) {
  const canvas = createCanvas(size, size);
  const ctx = canvas.getContext('2d');
  if (style === 'glassmorphism') glassBg(ctx, size, color);
  else minimalBg(ctx, size, color);
  drawEmoji(ctx, emoji, size);
  return canvas;
}

function saveIcon(canvas, outPath) {
  const out = fs.createWriteStream(outPath);
  const stream = canvas.createPNGStream();
  stream.pipe(out);
}

function main() {
  if (!fs.existsSync(OUT_DIR)) fs.mkdirSync(OUT_DIR, { recursive: true });
  for (const icon of ICONS) {
    for (const style of STYLES) {
      for (const size of SIZES) {
        const canvas = generateIcon(icon, style, size);
        const file = `${style}-maskable-${icon.name}-${size}x${size}.png`;
        saveIcon(canvas, path.join(OUT_DIR, file));
      }
    }
  }
  console.log('Icons generated in', OUT_DIR);
}

if (require.main === module) main();
