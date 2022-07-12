'use strict'

var map = L.map('map').setView([55.673043, 12.564782], 13);
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '© OpenStreetMap'
}).addTo(map);


function drawStations(stations) {
    stations.forEach(element => {
        L.marker([element.Coordinates.Lat, element.Coordinates.Lng], {
            title: "Yeet"
        }).bindPopup(`<h2>${element.Name}</h2><ul><li><strong>ID:</strong> ${element.StationID}</li></ul>`).addTo(map);
    });
    console.log(map);
}