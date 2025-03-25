function toggleMenu() {
    const navLinks = document.querySelector('.nav-links');
    navLinks.classList.toggle('active');
}

document.addEventListener('click', function(e) {
    const navLinks = document.querySelector('.nav-links');
    if (!e.target.closest('.navbar')) {
        navLinks.classList.remove('active');
    }
});

function togglePassword() {
    const passwordInput = document.getElementById('passwordInput');
    const toggleIcon = document.querySelector('.password-toggle i');
    const isPassword = passwordInput.type === 'password';

    passwordInput.type = isPassword ? 'text' : 'password';
    toggleIcon.classList.replace(isPassword ? 'fa-eye' : 'fa-eye-slash', isPassword ? 'fa-eye-slash' : 'fa-eye');
}

function handleAvatarPreview() {
    const reader = new FileReader();
    reader.onload = (e) => {
        document.getElementById('avatarPreview').src = e.target.result;
    };
    reader.readAsDataURL(this.files[0]);
}

function showError(message) {
    const popup = document.getElementById('errorPopup');
    document.getElementById('errorMessage').textContent = message;
    popup.classList.add('show');
    setTimeout(() => popup.classList.remove('show'), 5000);
}

function closeError() {
    document.getElementById('errorPopup').classList.remove('show');
}

async function handleFormSubmit(e, method, url, successUrl) {
    e.preventDefault();
    const formData = new FormData(e.target);

    // Append avatar file if exists
    const avatarFile = document.getElementById('avatarUpload').files[0];
    if (avatarFile) formData.append('avatar', avatarFile);

    try {
        const response = await fetch(url, {
            method,
            body: formData
        });

        if (!response.ok) {
            const result = await response.json();
            throw new Error(result.message || 'An error occurred');
        } else {
            window.location.href = successUrl;
        }
    } catch (error) {
        console.error('Submission Error:', error);
        showError(error.message || 'Network error. Please try again.');
    }
}

document.getElementById('avatarUpload')?.addEventListener('change', handleAvatarPreview);

document.getElementById('userForm')?.addEventListener('submit', (e) => {
    handleFormSubmit(e, 'PUT', window.location.href, '/admin/users');
});

document.getElementById('userNewForm')?.addEventListener('submit', (e) => {
    handleFormSubmit(e, 'POST', '/admin/users', '/admin/users');
});

document.getElementById('profileForm')?.addEventListener('submit', (e) => {
    handleFormSubmit(e, 'PUT', window.location.href, window.location.href);
});

let deleteUserId = null;
function confirmDelete(userId) {
    deleteUserId = userId;
    document.getElementById('deleteDialog').style.display = 'block';
    document.body.insertAdjacentHTML('beforeend', '<div class="dialog-backdrop"></div>');
}
function closeDialog() {
    document.getElementById('deleteDialog').style.display = 'none';
    document.querySelector('.dialog-backdrop').remove();
}
document.getElementById('confirmDeleteBtn').addEventListener('click', async function() {
    if (deleteUserId) {
        try {
            const response = await fetch(`/admin/users/${deleteUserId}`, {
                method: 'DELETE'
            });

            if (!response.ok) {
                const result = await response.json();
                showError(result.message || 'An error occurred');
                return;
            }

            window.location.reload();
        } catch (error) {
            showError('Network error. Please try again.');
        }
    }
    closeDialog();
});

