function filterResources() {
	// Declare variables
	var input, filter, table, tr, td, i, txtValue;
	input = document.getElementById("filterText");
	filter = input.value.toUpperCase();
	table = document.getElementById("resourceTable");
	tr = table.getElementsByTagName("tr");

	// Loop through all table rows, and hide those who don't match the search query
	for (i = 0; i < tr.length; i++) {
		td = tr[i].getElementsByTagName("td")[0];
		if (td) {
			txtValue = td.textContent || td.innerText;
			if (txtValue.toUpperCase().indexOf(filter) > -1) {
				tr[i].style.display = "";
			} else {
				tr[i].style.display = "none";
			}
		}
	}
}

function filterAccessions() {
	
    var input = document.getElementById("filterText");
    var filter = input.value.toUpperCase();
    var table = document.getElementById("accessionsTable");
    var tr = table.getElementsByTagName("tr");
    
	for (var i = 0; i < tr.length; i++) {
		var td0 = tr[i].getElementsByTagName("td")[0];
		var td1 = tr[i].getElementsByTagName("td")[1];
		var td2 = tr[i].getElementsByTagName("td")[2];
		if (td0 || td1 || td2) {
			txtValue1 = td0.textContent || td0.innerText;
			txtValue2 = td1.textContent || td1.innerText;
			txtValue3 = td2.textContent || td2.innerText;
			if ( txtValue1.toUpperCase().indexOf(filter) > -1 || txtValue2.toUpperCase().indexOf(filter) > -1 || txtValue3.toUpperCase().indexOf(filter) > -1) {
				tr[i].style.display = "";
			} else {
				tr[i].style.display = "none";
			}
		}
	}
}

function filterResourceCols() {
	
    var input = document.getElementById("filterText");
    var filter = input.value.toUpperCase();
    var table = document.getElementById("resourceTable");
    var tr = table.getElementsByTagName("tr");
    
	for (var i = 0; i < tr.length; i++) {
		var td0 = tr[i].getElementsByTagName("td")[0];
		var td1 = tr[i].getElementsByTagName("td")[1];
		if (td0 || td1) {
			txtValue1 = td0.textContent || td0.innerText;
			txtValue2 = td1.textContent || td1.innerText;
			if ( txtValue1.toUpperCase().indexOf(filter) > -1 || txtValue2.toUpperCase().indexOf(filter) > -1 ) {
				tr[i].style.display = "";
			} else {
				tr[i].style.display = "none";
			}
		}
	}
}

function togglePassword() {
	var passwd = document.getElementById("password_1");
	if (passwd.type === "password") {
	  passwd.type = "text";
	} else {
	  passwd.type = "password";
	}
}