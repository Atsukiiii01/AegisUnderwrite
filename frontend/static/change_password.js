const token = localStorage.getItem("token");

if (!token) {
    window.location.href = "/auth";
}

document.getElementById("change-password-form")
.addEventListener("submit", async function (e) {
    e.preventDefault();

    const new_password = document.getElementById("new-password").value;

    const res = await fetch("/auth/change-password", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + token
        },
        body: JSON.stringify({ new_password })
    });

    const msg = document.getElementById("msg");

    if (!res.ok) {
        msg.innerHTML = "<p style='color:red'>Password update failed</p>";
        return;
    }

    localStorage.setItem("must_change_password", "false");
    msg.innerHTML = "<p style='color:green'>Password updated. Redirecting...</p>";

    setTimeout(() => {
        window.location.href = "/";
    }, 1200);
});
