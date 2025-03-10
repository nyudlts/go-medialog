function filterAccessions() {
    var input = document.getElementById("filterText");
    var filter = input.value.toUpperCase();
    var table = document.getElementById("accessionsTable");
    var tr = table.getElementsByTagName("tr");
    
	// Loop through all table rows, and hide those who don't match the search query
	for (i = 0; i < tr.length; i++) {
		td = tr[i].getElementsByTagName("td")[2];
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
