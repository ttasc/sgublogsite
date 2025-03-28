import { showError } from '../utils.js';

export default function initUsers() {
    // User-specific logic
    const avatarInput = document.getElementById('avatarUpload');
    avatarInput?.addEventListener('change', function() {
        const reader = new FileReader();
        reader.onload = (e) => {
            document.getElementById('avatarPreview').src = e.target.result;
        };
        reader.readAsDataURL(this.files[0]);
    });

    // User form submission
    document.getElementById('userForm')?.addEventListener('submit', (e) => {
        handleUserFormSubmit(e, 'PUT', window.location.href, '/admin/users');
    });

    document.getElementById('userNewForm')?.addEventListener('submit', (e) => {
        handleUserFormSubmit(e, 'POST', '/admin/users', '/admin/users');
    });

    document.getElementById('profileForm')?.addEventListener('submit', (e) => {
        handleUserFormSubmit(e, 'PUT', window.location.href, window.location.href);
    });
}

async function handleUserFormSubmit(e, method, url, successUrl) {
    e.preventDefault();
    const formData = new FormData(e.target);

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
        showError(error.message || 'Network error. Please try again.');
    }
}
