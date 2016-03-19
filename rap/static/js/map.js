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
	map.data.loadGeoJson("/resources");

	/*
//manual way of doing the 
	var httpRequest = new XMLHttpRequest();
	//httpRequest.setRequestHeader('Accept', 'application/json');
	httpRequest.onreadystatechange = function(){
	    if (httpRequest.readyState === XMLHttpRequest.DONE) {
		// everything is good, the response is received
		if (httpRequest.status === 200) {
		    console.log(httpRequest.responseText);
		    
		} else {
		    // there was a problem with the request,
		    // for example the response may contain a 404 (Not Found)
		    // or 500 (Internal Server Error) response code
		    console.log(httpRequest.status);
		}
	    } else {
		// still not ready
		console.log("loading geojson...");
	    }
	};
	
	httpRequest.open('GET', '/resources', true);
	httpRequest.send();
	*/
    }

    function setResourceTypeSelection(){
	var resourceTypes = legend.getElementsByTagName('li');
	//console.log(resourceTypes);
	
	//we will loop over each li and set its background color depend on its checkbox's checked state
	for(var li of resourceTypes){
	    var checkbox = li.getElementsByTagName('input')[0];
	    
	    if(checkbox.checked){
		li.getElementsByTagName('label')[0]
		    .classList.add('selected');
	    } else{
		li.getElementsByTagName('label')[0]
		    .classList.remove('selected');
	    }
	}
    }

    var legend = document.getElementById('legend');
    legend.addEventListener('click', setResourceTypeSelection);

    //initialize
    setResourceTypeSelection();
    setCenterGeolocation();
    searchForResources();

    //don't run this till window, document, and google are ready
})(window,document,google);
