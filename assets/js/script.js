import { initAdminCommon } from './admin/common.js';
import { initCommon } from './index.js';

document.addEventListener('DOMContentLoaded', initCommon);

const currentPath = window.location.pathname;

if (currentPath.includes('/admin')) {
    // Initialize Admin common functionality
    document.addEventListener('DOMContentLoaded', initAdminCommon);

    // Dynamic module loading for HTMX
    document.body.addEventListener('htmx:afterSwap', () => {
        const path = window.location.pathname;

        if (path.includes('/users')) {
            import('./admin/users.js').then(module => module.default());
        }
        else if (path.includes('/categories')) {
            import('./admin/categories.js').then(module => module.default());
        }
        else if (path.includes('/tags')) {
            import('./admin/tags.js').then(module => module.default());
        }
        else if (path.includes('/posts')) {
            import('./admin/posts.js').then(module => module.default());
        }
        else if (path.includes('/images')) {
            import('./admin/images.js').then(module => module.default());
        }

        initCommon();
        initAdminCommon();
    });

    // Initial page load handling
    if (currentPath.includes('/users')) {
        import('./admin/users.js').then(module => module.default());
    }
    else if (currentPath.includes('/categories')) {
        import('./admin/categories.js').then(module => module.default());
    }
    else if (currentPath.includes('/tags')) {
        import('./admin/tags.js').then(module => module.default());
    }
    else if (currentPath.includes('/posts')) {
        import('./admin/posts.js').then(module => module.default());
    }
    else if (currentPath.includes('/images')) {
        import('./admin/images.js').then(module => module.default());
    }
}
else {
    import('./index.js').then(module => module.default());
}
