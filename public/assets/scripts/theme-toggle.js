(function () {
    const toggleBtn = document.getElementById('theme-toggle');
    const lightIcon = toggleBtn.querySelector('.theme-icon-light');
    const darkIcon = toggleBtn.querySelector('.theme-icon-dark');
    const html = document.documentElement;

    function updateIcons() {
        if (html.classList.contains('dark')) {
            lightIcon.classList.remove('hidden');
            darkIcon.classList.add('hidden');
        } else {
            lightIcon.classList.add('hidden');
            darkIcon.classList.remove('hidden');
        }
    }

    // Initial icon state
    updateIcons();

    if (toggleBtn) {
        toggleBtn.addEventListener('click', () => {
            if (html.classList.contains('dark')) {
                html.classList.remove('dark');
                localStorage.setItem('theme', 'light');
            } else {
                html.classList.add('dark');
                localStorage.setItem('theme', 'dark');
            }
            updateIcons();
        });
    }
})();
