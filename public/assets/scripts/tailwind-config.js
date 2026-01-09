// External Tailwind configuration + small theme initializer
// This script mirrors the inline tailwind.config used in the template.
// It also initializes the initial dark/light class on <html> using
// localStorage (if set) or system preference.
//
// Usage: included after the Tailwind CDN script:
// <script src="https://cdn.tailwindcss.com?plugins=forms,container-queries"></script>
// <script src="/assets/scripts/tailwind-config.js"></script>

(function () {
  // Initialize theme (dark class) early to avoid flash
  try {
    var storedTheme = localStorage.getItem('theme');
    var prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
    if (storedTheme === 'dark' || (storedTheme === null && prefersDark)) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  } catch (e) {
    // If any error occurs (e.g. localStorage blocked), don't break the page.
    // Just let the default stylesheet behavior apply.
  }

  // Assign Tailwind config
  if (typeof tailwind !== 'undefined') {
    tailwind.config = tailwind.config || {};
    tailwind.config.darkMode = 'class';
    tailwind.config.theme = tailwind.config.theme || {};
    tailwind.config.theme.extend = Object.assign({}, tailwind.config.theme.extend || {}, {
      colors: {
        primary: '#ee2b8c',
        'primary-hover': '#d41b76',
        'background-light': '#f8f6f7',
        'background-dark': '#221019',
        'text-main': '#181114',
        'text-muted': '#896175'
      },
      fontFamily: {
        display: ['Spline Sans', 'sans-serif']
      },
      borderRadius: {
        DEFAULT: '1rem',
        lg: '1.5rem',
        xl: '2rem',
        '2xl': '3rem',
        full: '9999px'
      },
      boxShadow: {
        soft: '0 4px 20px -2px rgba(238, 43, 140, 0.1)',
        hover: '0 10px 25px -5px rgba(238, 43, 140, 0.15)'
      }
    });
  } else {
    // If tailwind isn't loaded yet, expose a global config holder so the CDN can pick it up later.
    // tailwind CDN checks for `tailwind.config` on the window, so set it there.
    window.tailwind = window.tailwind || {};
    window.tailwind.config = window.tailwind.config || {
      darkMode: 'class',
      theme: {
        extend: {
          colors: {
            primary: '#ee2b8c',
            'primary-hover': '#d41b76',
            'background-light': '#f8f6f7',
            'background-dark': '#221019',
            'text-main': '#181114',
            'text-muted': '#896175'
          },
          fontFamily: {
            display: ['Spline Sans', 'sans-serif']
          },
          borderRadius: {
            DEFAULT: '1rem',
            lg: '1.5rem',
            xl: '2rem',
            '2xl': '3rem',
            full: '9999px'
          },
          boxShadow: {
            soft: '0 4px 20px -2px rgba(238, 43, 140, 0.1)',
            hover: '0 10px 25px -5px rgba(238, 43, 140, 0.15)'
          }
        }
      }
    };
  }
})();
