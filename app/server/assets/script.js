'use strict';

document.getElementById("date_picker").addEventListener("keydown", function (e) {
  e.preventDefault();
});

document.getElementById("date_picker").addEventListener("input", function (e) {
  window.location.href = "/?date=" + e.target.value;
});

