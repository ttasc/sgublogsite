import { showError, showDialog, closeDialog } from '../utils.js';

let deleteId = null;

export function initAdminCommon() {
    // Delete handling
    document.addEventListener('click', (e) => {
        const element = e.target.closest('#deleteBtn');
        if (element) {
            deleteId = element.dataset.id;
            showDialog();
        }
    });
    document.getElementById('confirmDeleteBtn')?.addEventListener('click', async () => {
        handleDelete(window.location.href);
    });

    document.getElementById('closeDialog')?.addEventListener('click', closeDialog);
}

async function handleDelete(urlPrefix) {
    if (!deleteId) return;

    try {
        const response = await fetch(urlPrefix + '/' + deleteId, { method: 'DELETE' });
        if (!response.ok) {
            const result = await response.json();
            showError(result.message || 'An error occurred');
            return;
        }
        window.location.reload();
    } catch (error) {
        showError('Network error. Please try again.');
    }
    closeDialog();
}
