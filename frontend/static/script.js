const token = localStorage.getItem("token");
const mustChange = localStorage.getItem("must_change_password");

if (!token) {
    window.location.href = "/auth";
}

if (mustChange === "true") {
    window.location.href = "/change-password";
}
