const token = localStorage.getItem("token");
const role = localStorage.getItem("role");

if (!token || role !== "ADMIN") {
    alert("Admin access required");
    window.location.href = "/auth";
}

document.getElementById("create-user-form")
.addEventListener("submit", async function (e) {
    e.preventDefault();

    const username = document.getElementById("username").value;
    const roleValue = document.getElementById("role").value;

    const res = await fetch("/admin/create-user", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + token
        },
        body: JSON.stringify({ username, role: roleValue })
    });

    const result = document.getElementById("result");

    if (!res.ok) {
        const err = await res.json();
        result.innerHTML = `<p style="color:red">${err.detail}</p>`;
        return;
    }

    const data = await res.json();
    result.innerHTML = `
        <p><b>User created</b></p>
        <p>Username: ${data.username}</p>
        <p>Temp Password: ${data.temporary_password}</p>
    `;
});
