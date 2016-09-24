var ResourceAwareness = {};

(function() {
    "use strict";

    ResourceAwareness.map = new google.maps.Map($('#map')[0], {
	    center: {
            lat: 39.9522, 
            lng: -75.1635
        },
	    zoom: 16
    });

    ResourceAwareness.layers = new google.maps.MVCObject();
    ResourceAwareness.layers.setValues({
        food: ResourceAwareness.map,
        legal: ResourceAwareness.map,
        shelter: ResourceAwareness.map,
        medical: ResourceAwareness.map,
        dental: ResourceAwareness.map
    });

    ResourceAwareness.setCenterGeoLocation = function() {
        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(function(position) {
                var pos = {
                    lat: position.coords.latitude,
                    lng: position.coords.longitude
                };
        
                ResourceAwareness.map.setCenter(pos);
            }, function() {
                console.log('Geolocation could not determine your location');
            });
        } else {
            // Browser doesn't support Geolocation
            console.log('Your browser does not support geolocation or you chose not to share your location');
        }
    }

    ResourceAwareness.getResources = function () {
        $.ajax({
            url: "/resources",
            success: function(data) {
                ResourceAwareness.resources = data;
            },
            async: false
        });
    }
    
    ResourceAwareness.bindMapControls = function() {
        var infowindow =  new google.maps.InfoWindow({
            content: ""
        });

        ResourceAwareness.resources.features.forEach(function(feature) {
            var latLng = new google.maps.LatLng(
                feature.geometry.coordinates[1], //longitude
                feature.geometry.coordinates[0] //latitude
            );
            
            // Create our markers
            var marker = new google.maps.Marker({
                position: latLng,
                map: ResourceAwareness.map,
                icon: {
                    path: google.maps.SymbolPath.BACKWARD_CLOSED_ARROW, 
                    strokeColor: "black",
                    strokeWeight: 1,
                    scale: 8,
                    fillOpacity: 1
                },
                title: feature.properties.OrganizationName
            });

            switch (feature.properties.Category) {
                case "Medical":
                    marker.icon.fillColor = "red";
                    setUniqueCategoryList(feature.properties.Category);
                    marker.bindTo('map', ResourceAwareness.layers, 'medical');
                    break;
                case "Emergency Shelter":
                    marker.icon.fillColor = "blue";
                    setUniqueCategoryList(feature.properties.Category);
                    marker.bindTo('map', ResourceAwareness.layers, 'shelter');
                    break;
                case "Legal":
                    marker.icon.fillColor = "green";
                    setUniqueCategoryList(feature.properties.Category);
                    marker.bindTo('map', ResourceAwareness.layers, 'legal');
                    break;
                case "Emergency Food":
                    marker.icon.fillColor = "yellow";
                    setUniqueCategoryList(feature.properties.Category);
                    marker.bindTo('map', ResourceAwareness.layers, 'food');
                    break;
                default:
                    marker.icon.fillColor = "purple";
                    setUniqueCategoryList(feature.properties.Category);
                    break;
            }


			var content = '<div>' +
			    '<h3>'+feature.properties.OrganizationName+'</h3>'+
			    '<p>'+feature.properties.Description+'</p>'+
			    '<p>'+feature.properties.Address+'</p>';
			if(!!feature.properties.Days){
			    content+= '<p>Days: '+feature.properties.Days+'</p>';
			}
			if(!!feature.properties.TimeOpen){
			    content+= '<p>Open: '+feature.properties.TimeOpen+'</p>';
			}
			if(!!feature.properties.TimeClose){
			    content+= '<p>Close: '+feature.properties.TimeClose+'</p>';
			}
			if(!!feature.properties.PeopleServed){
			    content+= '<p>People Served: '+feature.properties.PeopleServed+'</p>';
			}
			if(!!feature.properties.PhoneNumber){
			    content+= '<p>Phone: '+feature.properties.PhoneNumber+'</p>';
			}
			content+= '</div>';

            bindInfoWindow(marker, ResourceAwareness.map, infowindow, content);
        });
    }

    function bindInfoWindow(marker, map, infowindow, content) {
        marker.addListener('click', function() {
            infowindow.setContent(content);
            infowindow.open(map, this);
        });
    }

    function setUniqueCategoryList (category) {
        if (typeof ResourceAwareness.categoryList === "undefined") {
            ResourceAwareness.categoryList = [];
        } else {
            if (ResourceAwareness.categoryList.indexOf(category) === -1) {
                ResourceAwareness.categoryList.push(category);
            }
        }
    }

    function setToggleButtons() {
        $("#cbMedical").on("change", function() {
            if (ResourceAwareness.layers.medical === null) {
                ResourceAwareness.layers.set("medical", ResourceAwareness.map);
            } else {
                ResourceAwareness.layers.set("medical", null);
            }
        });

        $("#cbShelter").on("change", function() {
            if (ResourceAwareness.layers.shelter === null) {
                ResourceAwareness.layers.set("shelter", ResourceAwareness.map);
            } else {
                ResourceAwareness.layers.set("shelter", null);
            }
        });

        $("#cbLegal").on("change", function() {
            if (ResourceAwareness.layers.legal === null) {
                ResourceAwareness.layers.set("legal", ResourceAwareness.map);
            } else {
                ResourceAwareness.layers.set("legal", null);
            }
        });

        $("#cbFood").on("change", function() {
            if (ResourceAwareness.layers.food === null) {
                ResourceAwareness.layers.set("food", ResourceAwareness.map);
            } else {
                ResourceAwareness.layers.set("food", null);
            }
        });
    }

    ResourceAwareness.getResources();
    ResourceAwareness.bindMapControls();
    ResourceAwareness.setCenterGeoLocation();
    setToggleButtons();
})(window, document, google);