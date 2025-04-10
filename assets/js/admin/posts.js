import { showError, showDialog, closeDialog, toggleSelectAll } from "../utils.js";

export default function initPosts() {

    document.getElementById('select-all-posts').addEventListener('change', toggleSelectAll('.post-checkbox'));
    document.getElementById('bulkDeletePostsBtn')?.addEventListener('click', handleBulkDelete);

}

async function handleBulkDelete() {
    const selectedPosts = Array.from(
        document.querySelectorAll('.post-checkbox:checked')
    ).map(checkbox => parseInt(checkbox.dataset.id));

    if (selectedPosts.length === 0) {
        showError("Please select at least one post.");
        return;
    }
    showDialog();
    document.getElementById('confirmDeleteBtn').onclick = async () => {
        try {
            const response = await fetch('/admin/posts/bulk', {
                method: 'DELETE',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ ids: selectedPosts }),
            });

            if (!response.ok) {
                const error = await response.json();
                showError(error.message || "Failed to delete posts.");
            } else {
                window.location.reload();
            }
        } catch (error) {
            console.error(error);
            showError('Network error. Please try again.');
        } finally {
            closeDialog();
        }
    };
}
