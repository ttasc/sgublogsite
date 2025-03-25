const urlParams = new URLSearchParams(window.location.search);
if (urlParams.get('error') === 'auth_failed') {
    const errorDiv = document.getElementById('errorMessage');
    errorDiv.textContent = 'Thông tin đăng nhập không đúng';
    errorDiv.style.display = 'block';

    window.history.replaceState({}, document.title, "/login");
}
