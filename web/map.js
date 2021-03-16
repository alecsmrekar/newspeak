function initMap() {
    var myLatlng = new google.maps.LatLng(-25.363882,131.044922);

    const map = new google.maps.Map(document.getElementById("map-canvas"),{
        zoom: 3,
        center: myLatlng,
    });
    var marker = new google.maps.Marker({
        position: myLatlng,
        map: map,
        draggable:true,
        title:"Drag me!"
    });

    let circle = new google.maps.Circle({
        strokeColor: "#ff0000",
        strokeOpacity: 0.8,
        strokeWeight: 2,
        fillColor: "#FF0000",
        fillOpacity: 0.5,
        map,
        center: myLatlng,
        radius: 500000,
    });

    marker.addListener("dragend", () => {
        circle.setCenter(marker.getPosition());
        circle.setRadius(circle.getRadius()*1.2);
    });
}
