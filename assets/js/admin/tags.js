import { slugify, showError, showDialog, debounce } from '../utils.js';

let isAutoSlug = true;

export default function initTags() {
    // Event Listeners
    document.getElementById('tagSearch').addEventListener('input', debounce(searchTags, 300));
    document.getElementById('selectAll').addEventListener('change', toggleSelectAll);
    document.getElementById('bulkDeleteBtn').addEventListener('click', handleBulkDelete);
    document.getElementById('tagForm').addEventListener('submit', handleFormSubmit);
    document.getElementById('slugRefreshBtn').addEventListener('click', refreshSlug);
    document.getElementById('addTagBtn').addEventListener('click', showTagModal);

    // Delegate events for dynamic content
    document.addEventListener('click', (e) => {
        const editBtn = e.target.closest('.edit-btn');
        if (editBtn) {
            openEditModal(JSON.parse(editBtn.dataset.tag));
        }

        const closeBtn = e.target.closest('#closeModal');
        if (closeBtn) {
            closeTagModal();
            return;
        }
    });

    // Auto-slug generation
    document.getElementById('tag_Name').addEventListener('input', (e) => {
        if (isAutoSlug) {
            document.getElementById('tagSlug').value = slugify(e.target.value);
        }
    });

    document.getElementById('tagSlug').addEventListener('input', () => {
        isAutoSlug = false;
    });

    // Initialize empty state
    updateEmptyState();
}

function updateEmptyState() {
    const isEmpty = document.querySelectorAll('#tagsContainer tr').length === 0;
    document.getElementById('emptyState').style.display = isEmpty ? 'block' : 'none';
}

function searchTags(e) {
    const query = e.target.value.toLowerCase();
    document.querySelectorAll('#tagsContainer tr').forEach(row => {
        const name = row.children[1].textContent.toLowerCase();
        const slug = row.children[2].textContent.toLowerCase();
        row.style.display = (name.includes(query) || slug.includes(query)) ? '' : 'none';
    });
}

function toggleSelectAll(e) {
    const checkboxes = document.querySelectorAll('.row-checkbox');
    checkboxes.forEach(checkbox => checkbox.checked = e.target.checked);
}

async function openEditModal(tag) {
    document.getElementById('tagId').value = tag.tag_id;
    document.getElementById('tag_Name').value = tag.name;
    document.getElementById('tagSlug').value = tag.slug;
    document.getElementById('modalTitle').textContent = 'Edit Tag';
    document.getElementById('tagModal').style.display = 'block';
}

function showTagModal() {
    document.getElementById('tagForm').reset();
    document.getElementById('tagModal').style.display = 'block';
    document.getElementById('modalTitle').textContent = 'New Tag';
    document.getElementById('tagId').value = '';
}

function closeTagModal() {
    document.getElementById('tagModal').style.display = 'none';
    document.getElementById('tagForm').reset();
    document.getElementById('tagId').value = '';
}

async function handleFormSubmit(e) {
    e.preventDefault();

    const tagData = {
        tag_id: parseInt(document.getElementById('tagId').value),
        name:   document.getElementById('tag_Name').value,
        slug:   document.getElementById('tagSlug').value
    };

    try {
        const method = tagData.id ? 'PUT' : 'POST';
        const url = tagData.id ? `/admin/tags/${tagData.id}` : '/admin/tags';

        const response = await fetch(url, {
            method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(tagData)
        });

        if (!response.ok) {
            const error = await response.json();
            showError(error.message);
        } else{
            window.location.reload();
        }
    } catch (error) {
        console.error(error);
        showError('Network error. Please try again.');
    }
}
function handleBulkDelete() {
    const tagIds = Array.from(
        document.querySelectorAll('.row-checkbox')
    ).filter(
        checkbox => checkbox.checked
    ).map(
        checkbox => parseInt(checkbox.dataset.id)
    );
    showDialog();
    document.getElementById('confirmDeleteBtn').onclick = async () => {
        try {
            const response = await fetch('/admin/tags/bulk', {
                method: 'DELETE',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ ids: tagIds })
            });

            if (!response.ok) {
                const error = await response.json();
                showError(error.message);
            } else{
                window.location.reload();
            }
        } catch (error) {
            console.error(error);
            showError('Network error. Please try again.');
        }
    };
}

function refreshSlug() {
    const name = document.getElementById('tag_Name').value;
    document.getElementById('tagSlug').value = slugify(name);
    isAutoSlug = true;
}
