import { slugify, showError, showDialog, closeDialog } from '../utils.js';

let pendingMove = null;
let isAutoSlug = true;
let currentCategory = null;

export default function initCategories() {
    // Initialize Sortable
    initSortable();

    // Slug generation
    document.getElementById('categoryName').addEventListener('input', (e) => {
        if (isAutoSlug) {
            document.getElementById('categorySlug').value = slugify(e.target.value);
        }
    });
    // Turn off auto slug generation when typing
    document.getElementById('categorySlug').addEventListener('input', function() { isAutoSlug = false });
    document.getElementById('slugRefreshBtn').addEventListener('click', refreshSlug);

    document.getElementById('confirmMoveBtn').addEventListener('click', confirmMove);
    document.getElementById('closeMoveDialog').addEventListener('click', closeMoveDialog);

    document.getElementById('addRootCategoryBtn').addEventListener('click', () => showCategoryModal(null));

    // Form submission
    document.getElementById('categoryForm').addEventListener('submit', async (e) => {
        handleCategoryFormSubmit(e);
    });

    document.addEventListener('click', (e) => {
        const addSubBtn = e.target.closest('#addSubCategoryBtn');
        if (addSubBtn) {
            showCategoryModal(addSubBtn.dataset.parentId);
            return;
        }

        const editBtn = e.target.closest('#editCategoryBtn');
        if (editBtn) {
            showEditModal(JSON.parse(editBtn.dataset.category));
            return;
        }

        const closeBtn = e.target.closest('#closeModal');
        if (closeBtn) {
            closeCategoryModal();
            return;
        }

        const toggleBtn = e.target.closest('#toggleChildren');
        if (toggleBtn) {
            toggleChildren(toggleBtn);
            return;
        }
    });
}

function confirmMove() {
    if (!pendingMove) return;
    const { item, from, to, oldIndex } = pendingMove;
    if (from !== to) {
        to.removeChild(item);
    }
    const referenceNode = from.children[oldIndex] || null;
    from.insertBefore(item, referenceNode);
    handleMoveRequest(item.dataset.id, getNewParentId(to));
    closeMoveDialog();
}

function getNewParentId(targetList) {
    const parentItem = targetList.closest('.category-item');
    return parentItem ? parseInt(parentItem.dataset.id) : null;
}

function closeMoveDialog() {
    closeDialog('moveDialog');
    pendingMove = null;
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
    document.getElementById('categoryId').value = category.id;
    document.getElementById('categoryName').value = category.name;
    document.getElementById('categorySlug').value = category.slug;
    document.getElementById('modalTitle').textContent = 'Edit Category';
    document.getElementById('categoryModal').style.display = 'block';
    isAutoSlug = false;
}

function closeCategoryModal() {
    document.getElementById('categoryModal').style.display = 'none';
    document.getElementById('categoryForm').reset();
}

function refreshSlug() {
    const name = document.getElementById('categoryName').value;
    const newSlug = slugify(name);
    document.getElementById('categorySlug').value = newSlug;
    isAutoSlug = true;
}

async function handleCategoryFormSubmit(e) {
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
}

function initSortable() {
    const lists = document.querySelectorAll('.category-list');

    lists.forEach(list => {
        Sortable.create(list, {
            group: "categories",
            handle: ".drag-handle",
            animation: 150,
            onEnd: async (evt) => {
                // Lưu thông tin move
                pendingMove = {
                    item: evt.item,
                    from: evt.from,
                    to: evt.to,
                    oldIndex: evt.oldIndex,
                    newIndex: evt.newIndex
                };

                // Hiển thị dialog xác nhận
                showDialog('moveDialog');
            }
        });
    });
}

async function handleMoveRequest(categoryId, newParentId) {
    try {
        const response = await fetch(`/admin/categories/${categoryId}/move`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                new_parent_id: newParentId
            })
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

