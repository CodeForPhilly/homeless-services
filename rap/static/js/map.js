//This page is based on the tutorials from google maps listed below

//google has some advice for mobile once we get past the basics of map display
//https://developers.google.com/maps/documentation/javascript/basics

//this lovely tutorial tells us how to add/remove markers from the map
//https://developers.google.com/maps/documentation/javascript/markers


(function(){
    "use strict";

    var map = new google.maps.Map(document.getElementById('map'), {
	center: {lat: 39.9522, lng: -75.1635},
	zoom: 16
    });
    
    // Set the center of the map based on the browser's location
    function setCenterGeolocation(){
	if (navigator.geolocation) {
	    navigator.geolocation.getCurrentPosition(function(position) {
		//console.log(position);

		var pos = {
		    lat: position.coords.latitude,
		    lng: position.coords.longitude
		};
		
		map.setCenter(pos);
	    }, function() {
		console.log('Geolocation could not determine your location');
	    });
	} else {
	    // Browser doesn't support Geolocation
	    console.log('Your browser does not support geolocation or you chose not to share your location');
	}
    }

    //Searches for RAP markers that relate to the search from the form on the page
    //returns all markers when the form is empty
    function searchForResources(){
	//This is convenient but we can't create the info windows easily this way
	//map.data.loadGeoJson("/resources");
	
	//manual way of doing the request and adding the marker/info windows
	var httpRequest = new XMLHttpRequest();
	//httpRequest.setRequestHeader('Accept', 'application/json');
	httpRequest.onreadystatechange = function(){
	    if (httpRequest.readyState === XMLHttpRequest.DONE) {
		// everything is good, the response is received
		if (httpRequest.status === 200) {
		    //console.log(httpRequest.responseText);
		    var resources = JSON.parse(httpRequest.responseText);

		    var infowindow =  new google.maps.InfoWindow({
			content: ""
		    });
		    
		    for(var i=0;i< resources.features.length;i++){
			var res = resources.features[i];

			var latLng = new google.maps.LatLng(
			    res.geometry.coordinates[1],  //longitude
			    res.geometry.coordinates[0]   //latitude
			);
			
			// Creating a marker and putting it on the map
			var marker = new google.maps.Marker({
			    position: latLng,
			    map: map,
			    icon: {
				path: google.maps.SymbolPath.BACKWARD_CLOSED_ARROW,
				fillColor: "yellow", //Emergency Food
				strokeColor: "black",
				strokeWeight: 1,
				scale: 8,
				fillOpacity: 1
			    },
			    title: res.properties.OrganizationName
			});

			switch(res.properties.Category){
			case "Medical":
			    marker.icon.fillColor= "red";
			    break;
			case "Emergency Shelter":
			    marker.icon.fillColor= "blue";
			    break;
			case "Legal":
			    marker.icon.fillColor= "green";
			    break;
			}
			    

			var content = '<div>' +
			    '<h3>'+res.properties.OrganizationName+'</h3>'+
			    '<p>'+res.properties.Description+'</p>'+
			    '<p>'+res.properties.Address+'</p>';
			if(!!res.properties.Days){
			    content+= '<p>Days: '+res.properties.Days+'</p>';
			}
			if(!!res.properties.TimeOpen){
			    content+= '<p>Open: '+res.properties.TimeOpen+'</p>';
			}
			if(!!res.properties.TimeClose){
			    content+= '<p>Close: '+res.properties.TimeClose+'</p>';
			}
			if(!!res.properties.PeopleServed){
			    content+= '<p>People Served: '+res.properties.PeopleServed+'</p>';
			}
			if(!!res.properties.PhoneNumber){
			    content+= '<p>Phone: '+res.properties.PhoneNumber+'</p>';
			}
			content+= '</div>';
			
			bindInfoWindow(marker, map, infowindow, content);
		    }
		} else {
		    // there was a problem with the request,
		    // for example the response may contain a 404 (Not Found)
		    // or 500 (Internal Server Error) response code
		    console.log(httpRequest.status);
		}
	    } else {
		// still not ready
		console.log("loading resources...");
	    }
	};
	
	httpRequest.open('GET', '/resources', true);
	httpRequest.send();
	
    }

    function bindInfoWindow(marker, map, infowindow, description) {
	marker.addListener('click', function() {
            infowindow.setContent(description);
            infowindow.open(map, this);
	});
    }
    
    //var legend = document.getElementById('legend');
    //legend.addEventListener('click', setResourceTypeSelection);

    //initialize
    //setResourceTypeSelection();
    setCenterGeolocation();
    searchForResources();

    //don't run this till window, document, and google are ready
})(window,document,google);
