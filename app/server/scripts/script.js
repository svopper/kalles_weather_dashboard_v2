'use strict';

document.getElementById("date_picker").addEventListener("keydown", function (e) {
  e.preventDefault();
});

document.getElementById("date_picker").addEventListener("input", function (e) {
  window.location.href = "/?date=" + e.target.value;
});

var map = L.map('map').setView([55.673043, 12.564782], 13);
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '© OpenStreetMap'
}).addTo(map);

  