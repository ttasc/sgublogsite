import { initAdminCommon } from './admin/common.js';
import { initCommon } from './index.js';

document.addEventListener('DOMContentLoaded', initCommon);

const currentPath = window.location.pathname;

if (currentPath.includes('/admin')) {
    // Initialize Admin common functionality
    document.addEventListener('DOMContentLoaded', initAdminCommon);

    // Dynamic module loading for HTMX
    document.body.addEventListener('htmx:afterSwap', ({ detail }) => {
        if (detail.target.id === 'content') {
            const path = window.location.pathname;

            if (path.includes('/categories')) {
                import('./admin/categories.js').then(module => module.default());
            }
            else if (path.includes('/users')) {
                import('./admin/users.js').then(module => module.default());
            }

            initAdminCommon();
        }
    });

    // Initial page load handling
    if (currentPath.includes('/categories')) {
        import('./admin/categories.js').then(module => module.default());
    }
    else if (currentPath.includes('/users')) {
        import('./admin/users.js').then(module => module.default());
    }
}
else {
    import('./index.js').then(module => module.default());
}
