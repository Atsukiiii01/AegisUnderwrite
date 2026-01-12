let timer;
const TIMEOUT = 60000;

function resetTimer() {
    clearTimeout(timer);
    timer = setTimeout(() => {
        alert("Session expired. Please login again.");
        window.location.href = "/logout";
    }, TIMEOUT);
}

["click", "mousemove", "keydown"].forEach(e =>
    document.addEventListener(e, resetTimer)
);

resetTimer();
