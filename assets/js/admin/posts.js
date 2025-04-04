import { showError, showDialog, closeDialog, toggleSelectAll } from "../utils.js";

export default function initPosts() {

    document.getElementById('select-all-posts').addEventListener('change', toggleSelectAll('.post-checkbox'));
    document.getElementById('bulkDeletePostsBtn')?.addEventListener('click', handleBulkDelete);

    // Date filter validation
    document.querySelectorAll('.filters input').forEach(input => {
        input.addEventListener('change', function() {
            const start = document.getElementById('startDate').value;
            const end = document.getElementById('endDate').value;
            if (start && end && new Date(start) > new Date(end)) {
                alert('End date must be after start date');
                this.value = '';
            }
        });
    });

    // Thêm sự kiện giữ lại các filter khi gửi request
    document.querySelectorAll('.filter-select').forEach(select => {
        select.addEventListener('change', function() {
            const hxGet = this.getAttribute('hx-get');
            const params = new URLSearchParams(window.location.search);

            // Cập nhật tham số filter
            params.set(this.name, this.value);

            // Giữ lại các tham số khác (page, search, date)
            this.setAttribute('hx-get', hxGet + '?' + params.toString());
        });
    });
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
