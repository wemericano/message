document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginFormElement');
    const signupForm = document.getElementById('signupFormElement');
    const messageDiv = document.getElementById('message');
    
    // 로그인 폼 처리
    loginForm.addEventListener('submit', function(e) {
        e.preventDefault();
        
        const username = document.getElementById('loginUsername').value;
        const password = document.getElementById('loginPassword').value;
        
        const data = {
            username: username,
            password: password
        };
        
        // 로딩 UI 표시
        showLoading('Authenticating...');
        
        fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                // userID를 localStorage에 저장
                if (data.user_id) {
                    localStorage.setItem('userID', data.user_id);
                }
                showLoading('Login successful! Redirecting...');
                setTimeout(() => {
                    window.location.href = '/html/chat.html';
                }, 1500);
            } else {
                hideLoading();
                messageDiv.classList.add('show');
                messageDiv.textContent = '로그인 실패: ' + data.message;
                messageDiv.style.color = '#f44336';
                messageDiv.style.backgroundColor = '#ffebee';
            }
        })
        .catch(error => {
            hideLoading();
            messageDiv.classList.add('show');
            messageDiv.textContent = '에러: ' + error.message;
            messageDiv.style.color = '#f44336';
            messageDiv.style.backgroundColor = '#ffebee';
        });
    });
    
    // 회원가입 폼 처리
    signupForm.addEventListener('submit', function(e) {
        e.preventDefault();
        
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        const name = document.getElementById('name').value;
        const email = document.getElementById('email').value;
        
        const data = {
            username: username,
            password: password,
            name: name,
            email: email
        };
        
        // 로딩 UI 표시
        showLoading('Creating account...');
        
        fetch('/api/signup', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showLoading('Account created! Redirecting...');
                setTimeout(() => {
                    window.location.href = '/html/chat.html';
                }, 1500);
            } else {
                hideLoading();
                messageDiv.classList.add('show');
                messageDiv.textContent = '회원가입 실패: ' + data.message;
                messageDiv.style.color = '#f44336';
                messageDiv.style.backgroundColor = '#ffebee';
            }
        })
        .catch(error => {
            hideLoading();
            messageDiv.classList.add('show');
            messageDiv.textContent = '에러: ' + error.message;
            messageDiv.style.color = '#f44336';
            messageDiv.style.backgroundColor = '#ffebee';
        });
    });
});

// 로딩 UI 함수
function showLoading(text) {
    const overlay = document.getElementById('loadingOverlay');
    overlay.classList.add('show');
}

function hideLoading() {
    const overlay = document.getElementById('loadingOverlay');
    overlay.classList.remove('show');
}

// 탭 전환 함수
function showLogin() {
    document.getElementById('loginForm').style.display = 'block';
    document.getElementById('signupForm').style.display = 'none';
    document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
    event.target.classList.add('active');
    document.getElementById('message').classList.remove('show');
}

function showSignup() {
    document.getElementById('loginForm').style.display = 'none';
    document.getElementById('signupForm').style.display = 'block';
    document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
    event.target.classList.add('active');
    document.getElementById('message').classList.remove('show');
}

