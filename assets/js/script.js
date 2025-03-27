let isAutoSlug = true;
var deleteId = null;
var currentCategory = null;

function toggleMenu() {
    const navLinks = document.querySelector('.nav-links');
    navLinks.classList.toggle('active');
}

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

function toggleChildren(element) {
    const parent = element.closest('.category-item');
    parent.classList.toggle('collapsed');
    element.querySelector('.toggle-icon').classList.toggle('fa-chevron-down');
    element.querySelector('.toggle-icon').classList.toggle('fa-chevron-right');
}

function showCategoryModal(parentId) {
    currentCategory = null;
    document.getElementById('parentId').value = parentId || '';
    document.getElementById('modalTitle').textContent = parentId ? 'Add Subcategory' : 'New Category';
    document.getElementById('categoryModal').style.display = 'block';
}

function showEditModal(category) {
    currentCategory = category;
    document.getElementById('categoryId').value = category.ID;
    document.getElementById('categoryName').value = category.Name;
    document.getElementById('categorySlug').value = category.Slug;
    document.getElementById('modalTitle').textContent = 'Edit Category';
    document.getElementById('categoryModal').style.display = 'block';
    isAutoSlug = false;
}

function closeCategoryModal() {
    document.getElementById('categoryModal').style.display = 'none';
    document.getElementById('categoryForm').reset();
}

function confirmDelete(id) {
    deleteId = id;
    document.getElementById('deleteDialog').style.display = 'block';
    document.body.insertAdjacentHTML('beforeend', '<div class="dialog-backdrop"></div>');
}

function closeDialog() {
    document.getElementById('deleteDialog').style.display = 'none';
    document.querySelector('.dialog-backdrop').remove();
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
        console.error('Submission Error:', error);
        showError(error.message || 'Network error. Please try again.');
    }
}

async function handleDelete(urlPrefix, deleteId) {
    if (deleteId) {
        try {
            const response = await fetch(urlPrefix + '/' + deleteId, {
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
}

function slugify(string) {
    const a = 'àáäâãåăæąçćčđďèéěėëêęğǵḧìíïîįłḿǹńňñòóöôœøṕŕřßşśšșťțùúüûǘůűūųẃẍÿýźžż·/_,:;'
    const b = 'aaaaaaaaacccddeeeeeeegghiiiiilmnnnnooooooprrsssssttuuuuuuuuuwxyyzzz------'
    const p = new RegExp(a.split('').join('|'), 'g')
    return string.toString().toLowerCase()
        .replace(/á|à|ả|ạ|ã|ă|ắ|ằ|ẳ|ẵ|ặ|â|ấ|ầ|ẩ|ẫ|ậ/gi, 'a')
        .replace(/é|è|ẻ|ẽ|ẹ|ê|ế|ề|ể|ễ|ệ/gi, 'e')
        .replace(/i|í|ì|ỉ|ĩ|ị/gi, 'i')
        .replace(/ó|ò|ỏ|õ|ọ|ô|ố|ồ|ổ|ỗ|ộ|ơ|ớ|ờ|ở|ỡ|ợ/gi, 'o')
        .replace(/ú|ù|ủ|ũ|ụ|ư|ứ|ừ|ử|ữ|ự/gi, 'u')
        .replace(/ý|ỳ|ỷ|ỹ|ỵ/gi, 'y')
        .replace(/đ/gi, 'd')
        .replace(/\s+/g, '-')
        .replace(p, c => b.charAt(a.indexOf(c)))
        .replace(/&/g, '-and-')
        .replace(/[^\w\-]+/g, '')
        .replace(/\-\-+/g, '-')
        .replace(/^-+/, '')
        .replace(/-+$/, '')
}

function refreshSlug() {
    const name = document.getElementById('categoryName').value;
    const newSlug = slugify(name);
    document.getElementById('categorySlug').value = newSlug;
    isAutoSlug = true;
}

document.getElementById('categoryForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const categoryData = {
        id:         parseInt(document.getElementById('categoryId').value),
        parent_id:  parseInt(document.getElementById('parentId').value),
        name:       document.getElementById('categoryName').value,
        slug:       document.getElementById('categorySlug').value
    };

    try {
        const url = categoryData.id ? `/admin/categories/${categoryData.id}` : '/admin/categories';
        const method = categoryData.id ? 'PUT' : 'POST';

        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(categoryData)
        });

        if (!response.ok) {
            const result = await response.json();
            showError(result.message || 'An error occurred');
            return;
        }

        window.location.reload();
    } catch (error) {
        showError('Error saving category');
    }
});

document.addEventListener('click', function(e) {
    const navLinks = document.querySelector('.nav-links');
    if (!e.target.closest('.navbar')) {
        navLinks.classList.remove('active');
    }
});

document.getElementById('avatarUpload')?.addEventListener('change', handleAvatarPreview);

document.getElementById('userForm')?.addEventListener('submit', (e) => {
    handleUserFormSubmit(e, 'PUT', window.location.href, '/admin/users');
});

document.getElementById('userNewForm')?.addEventListener('submit', (e) => {
    handleUserFormSubmit(e, 'POST', '/admin/users', '/admin/users');
});

document.getElementById('profileForm')?.addEventListener('submit', (e) => {
    handleUserFormSubmit(e, 'PUT', window.location.href, window.location.href);
});

document.getElementById('confirmDeleteBtn').addEventListener('click', () => {
    handleDelete(window.location.href, deleteId);
});

document.getElementById('categoryName').addEventListener('input', function(e) {
    if (isAutoSlug) {
        const slug = slugify(e.target.value);
        document.getElementById('categorySlug').value = slug;
    }
});

document.getElementById('categorySlug').addEventListener('input', function() {
    isAutoSlug = false;
});
