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
