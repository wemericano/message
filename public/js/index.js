document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('signupForm');
    const messageDiv = document.getElementById('message');
    
    form.addEventListener('submit', function(e) {
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
        
        fetch('/api/signup', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => response.json())
        .then(data => {
            messageDiv.classList.add('show');
            if (data.success) {
                messageDiv.textContent = '회원가입 성공!';
                messageDiv.style.color = '#4caf50';
                messageDiv.style.backgroundColor = '#e8f5e9';
            } else {
                messageDiv.textContent = '회원가입 실패: ' + data.message;
                messageDiv.style.color = '#f44336';
                messageDiv.style.backgroundColor = '#ffebee';
            }
        })
        .catch(error => {
            messageDiv.classList.add('show');
            messageDiv.textContent = '에러: ' + error.message;
            messageDiv.style.color = '#f44336';
            messageDiv.style.backgroundColor = '#ffebee';
        });
    });
});

