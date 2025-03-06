function togglePassword() {
	var passwd = document.getElementById("password_1");
	if (passwd.type === "password") {
	  passwd.type = "text";
	} else {
	  passwd.type = "password";
	}
} 