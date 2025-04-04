import { showError, showDialog, closeDialog } from '../utils.js';

export default function initImages() {

    document.getElementById('bulkDeleteImagesBtn')?.addEventListener('click', handleBulkDelete);
    document.getElementById('uploadImageBtn')?.addEventListener('click', () => showDialog('uploadDialog'));
    document.getElementById('uploadForm')?.addEventListener('submit', handleFormUpload);

    document.addEventListener('click', (e) => {
        const closeBtn = e.target.closest('#closeUploadDialog');
        if (closeBtn) {
            closeDialog('uploadDialog');
            return;
        }
    });
}

async function handleFormUpload(e) {
    e.preventDefault();
    const formData = new FormData(e.target);

    try {
        const response = await fetch('/admin/images/upload', {
            method: 'POST',
            body: formData
        });

        if (!response.ok){
            const error = await response.json();
            showError(error.message || "Failed to upload images.");
        } else {
            window.location.reload();
        }
    } catch (error) {
        console.error(error);
        showError('Network error. Please try again.');
    } finally {
        closeDialog('uploadDialog');
    }
}

async function handleBulkDelete() {
    const selectedImages = Array.from(
        document.querySelectorAll('.image-checkbox-input:checked')
    ).map(checkbox => checkbox.dataset.id);

    if (selectedImages.length === 0) {
        showError("Please select at least one image to delete.");
        return;
    }

    showDialog();

    document.getElementById('confirmDeleteBtn').onclick = async () => {
        try {
            const response = await fetch('/admin/images/bulk', {
                method: 'DELETE',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ ids: selectedImages })
            });

            if (!response.ok){
                const error = await response.json();
                showError(error.message || "Failed to delete images.");
            } else {
                window.location.reload();
            }
        } catch (error) {
            console.error(error);
            showError('Network error. Please try again.');
        } finally {
            closeDialog('dialog');
        }
    };
}
