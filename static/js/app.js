// ============================================
// SistemInfoGo — JavaScript Interaktif
// ============================================

document.addEventListener('DOMContentLoaded', function () {
    // === Mobile Navigation Toggle ===
    const navToggle = document.getElementById('navToggle');
    const navLinks = document.querySelector('.nav-links');

    if (navToggle) {
        navToggle.addEventListener('click', function () {
            navLinks.classList.toggle('active');

            // Animasi hamburger ke X
            const spans = navToggle.querySelectorAll('span');
            navToggle.classList.toggle('open');

            if (navToggle.classList.contains('open')) {
                spans[0].style.transform = 'rotate(45deg) translate(5px, 5px)';
                spans[1].style.opacity = '0';
                spans[2].style.transform =
                    'rotate(-45deg) translate(5px, -5px)';
            } else {
                spans[0].style.transform = 'none';
                spans[1].style.opacity = '1';
                spans[2].style.transform = 'none';
            }
        });
    }

    // === Navbar Scroll Effect ===
    const navbar = document.querySelector('.navbar');
    let lastScrollY = window.scrollY;

    window.addEventListener('scroll', function () {
        if (window.scrollY > 50) {
            navbar.style.background = 'rgba(10, 10, 15, 0.95)';
            navbar.style.boxShadow = '0 4px 20px rgba(0, 0, 0, 0.3)';
        } else {
            navbar.style.background = 'rgba(10, 10, 15, 0.85)';
            navbar.style.boxShadow = 'none';
        }
        lastScrollY = window.scrollY;
    });

    // === Smooth Scroll untuk anchor links ===
    document.querySelectorAll('a[href^="#"]').forEach(function (anchor) {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start',
                });
            }
        });
    });

    // === Flash/Alert Auto Dismiss ===
    const alerts = document.querySelectorAll('.alert');
    alerts.forEach(function (alert) {
        setTimeout(function () {
            alert.style.transition = 'opacity 0.5s ease, transform 0.5s ease';
            alert.style.opacity = '0';
            alert.style.transform = 'translateY(-10px)';
            setTimeout(function () {
                alert.remove();
            }, 500);
        }, 5000);
    });

    // === Form Validation Visual Feedback ===
    const formInputs = document.querySelectorAll('.form-input');
    formInputs.forEach(function (input) {
        input.addEventListener('blur', function () {
            if (this.value.trim() !== '' && this.checkValidity()) {
                this.style.borderColor = 'rgba(6, 214, 160, 0.5)';
            } else if (this.value.trim() !== '' && !this.checkValidity()) {
                this.style.borderColor = 'rgba(239, 68, 68, 0.5)';
            } else {
                this.style.borderColor = '';
            }
        });

        input.addEventListener('focus', function () {
            this.style.borderColor = '';
        });
    });

    // === Password Match Validation ===
    const password = document.getElementById('password');
    const confirmPassword = document.getElementById('confirm_password');

    if (password && confirmPassword) {
        confirmPassword.addEventListener('input', function () {
            if (this.value !== password.value) {
                this.style.borderColor = 'rgba(239, 68, 68, 0.5)';
            } else {
                this.style.borderColor = 'rgba(6, 214, 160, 0.5)';
            }
        });
    }

    // === Password Toggle Visibility ===
    const passwordToggle = document.getElementById('passwordToggle');
    const passwordInput = document.getElementById('password');
    const eyeIcon = document.getElementById('eyeIcon');

    if (passwordToggle && passwordInput) {
        passwordToggle.addEventListener('click', function () {
            if (passwordInput.type === 'password') {
                passwordInput.type = 'text';
                eyeIcon.textContent = '🙈';
            } else {
                passwordInput.type = 'password';
                eyeIcon.textContent = '👁️';
            }
        });
    }

    // === Form Loading State ===
    const loginForm = document.getElementById('loginForm');
    const loginBtn = document.getElementById('loginBtn');

    if (loginForm && loginBtn) {
        loginForm.addEventListener('submit', function () {
            loginBtn.classList.add('btn-loading');
        });
    }

    // === Intersection Observer untuk animasi scroll ===
    const observerOptions = {
        threshold: 0.1,
        rootMargin: '0px 0px -50px 0px',
    };

    const observer = new IntersectionObserver(function (entries) {
        entries.forEach(function (entry) {
            if (entry.isIntersecting) {
                entry.target.style.opacity = '1';
                entry.target.style.transform = 'translateY(0)';
                observer.unobserve(entry.target);
            }
        });
    }, observerOptions);

    // Observe elements yang akan di-animasikan
    document
        .querySelectorAll('.stat-card, .article-card')
        .forEach(function (el) {
            observer.observe(el);
        });
});
