document.getElementById("limit").onchange = function() {
	var limit = document.getElementById("limit").value;
	var page = 1;
	if(document.location.search != "") {
		const params = new URLSearchParams(document.location.search);
		page = params.get("page");
	}

	window.location.href = `/entries?page=${page}&limit=${limit}`;
}

document.getElementById("page").onchange = function() {
	var page = document.getElementById("page").value;
	if(page == "") {
		page = 1;
	}
	var limit = 10;
	if(document.location.search != "") {
		const params = new URLSearchParams(document.location.search);
		limit = params.get("limit");
	}
	window.location.href = `/entries?page=${page}&limit=${limit}`;
};

document.getElementById("filter").onchange = function() {
	var page = document.getElementById("page").value;
	if(page == "") {
		page = 1
	}

	var limit = document.getElementById("limit").value;
	if(limit == "") {
		limit = 10
	}

	var filter = document.getElementById("filter").value;
	window.location.href = `/entries?page=${page}&limit=${limit}&filter=${filter}`;
}

document.getElementById("acc-limit").onchange = function() {
	var limit = document.getElementById("acc-limit").value;
	var page = 1;
	if(document.location.search != "") {
		const params = new URLSearchParams(document.location.search);
		page = params.get("page");
	}

	window.location.href = `/entries?page=${page}&limit=${limit}`;
}

document.getElementById("acc-page").onchange = function() {
	var page = document.getElementById("acc-page").value;
	if(page == "") {
		page = 1;
	}
	var limit = 10;
	if(document.location.search != "") {
		const params = new URLSearchParams(document.location.search);
		limit = params.get("limit");
	}
	window.location.href = `/entries?page=${page}&limit=${limit}`;
};

document.getElementById("acc-filter").onchange = function() {
	var page = document.getElementById("acc-page").value;
	if(page == "") {
		page = 1
	}

	var limit = document.getElementById("acc-limit").value;
	if(limit == "") {
		limit = 10
	}

	var filter = document.getElementById("acc-filter").value;
	window.location.href = `/entries?page=${page}&limit=${limit}&filter=${filter}`;
}

document.getElementById("globalSearch").addEventListener("keyup", function(event) {
	event.preventDefault();
	if (event.keyCode === 13) {
		var searchValue = document.getElementById("globalSearch").value;
		window.location.href = "/search?query=" + encodeURIComponent(searchValue);
	}
});  

$( function() {
	$( "#tabs" ).tabs();
} );



