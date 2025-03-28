import { closeError } from './utils.js';

export function initCommon() {
    document.getElementById("closeError")?.addEventListener("click", () => {
        closeError();
    });
}

export default function initIndex() {
    // Home navigation menu toggle
    document.getElementById("menuToggle").addEventListener("click", () => {
        const navLinks = document.querySelector('.nav-links');
        navLinks.classList.toggle('active');
    });
}
