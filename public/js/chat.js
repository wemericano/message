// 검색 입력 이벤트
document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.querySelector('.search-input');
    const friendsContent = document.querySelector('.friends-content');
    const notificationBell = document.getElementById('notificationBell');
    const notificationBadge = document.getElementById('notificationBadge');
    let searchTimeout;
    let ws = null;
    let currentUserID = null; // 로그인한 사용자 ID (나중에 세션에서 가져와야 함)
    let unreadCount = 0;

    searchInput.addEventListener('input', function(e) {
        const searchTerm = e.target.value.trim();
        
        // 디바운싱: 300ms 후에 검색 실행
        clearTimeout(searchTimeout);
        
        if (searchTerm === '') {
            // 검색어가 비어있으면 기본 화면 표시
            showDefaultView();
            return;
        }
        
        searchTimeout = setTimeout(() => {
            searchUsers(searchTerm);
        }, 300);
    });
    
    function searchUsers(name) {
        fetch(`/api/search?name=${encodeURIComponent(name)}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                const contentType = response.headers.get("content-type");
                if (!contentType || !contentType.includes("application/json")) {
                    throw new Error("Response is not JSON");
                }
                return response.json();
            })
            .then(data => {
                if (data.success) {
                    displaySearchResults(data.users);
                } else {
                    showNoResults();
                }
            })
            .catch(error => {
                console.error('### Search error:', error);
                showNoResults();
            });
    }
    
    function displaySearchResults(users) {
        if (users.length === 0) {
            showNoResults();
            return;
        }
        
        friendsContent.innerHTML = '';
        
        users.forEach(user => {
            const friendItem = document.createElement('div');
            friendItem.className = 'friend-item';
            
            const firstLetter = user.name.charAt(0).toUpperCase();
            
            friendItem.innerHTML = `
                <div class="friend-avatar">${firstLetter}</div>
                <div class="friend-info">
                    <div class="friend-name">${escapeHtml(user.name)}</div>
                    <div class="friend-last-message">${escapeHtml(user.username)}</div>
                </div>
            `;
            
            // 클릭 시 친구 요청
            friendItem.addEventListener('click', function() {
                sendFriendRequest(user.id);
            });
            
            friendsContent.appendChild(friendItem);
        });
    }
    
    function showNoResults() {
        friendsContent.innerHTML = `
            <div class="no-data">
                <div class="no-data-icon">🔍</div>
                <div>검색 결과가 없습니다</div>
            </div>
        `;
    }
    
    function showDefaultView() {
        friendsContent.innerHTML = `
            <div class="no-data">
                <div class="no-data-icon">👥</div>
                <div>친구 목록이 없습니다</div>
            </div>
        `;
    }
    
    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
    
    // 친구 요청 전송
    function sendFriendRequest(toUserID) {
        // 임시로 fromUserID를 1로 설정 (나중에 세션에서 가져와야 함)
        const fromUserID = currentUserID || 1;
        
        fetch('/api/friend/request', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                from_user_id: fromUserID,
                to_user_id: toUserID
            })
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                alert('친구 요청을 보냈습니다');
            } else {
                alert(data.message || '친구 요청 전송에 실패했습니다');
            }
        })
        .catch(error => {
            console.error('### Friend request error:', error);
            alert('친구 요청 전송에 실패했습니다');
        });
    }
    
    // WebSocket 연결
    function connectWebSocket(userID) {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;
        
        ws = new WebSocket(wsUrl);
        
        ws.onopen = function() {
            console.log('### WebSocket connected');
            // 사용자 ID 등록
            if (userID) {
                const registerMsg = {
                    type: 'register',
                    user_id: String(userID)
                };
                console.log('### Registering user ID:', userID);
                ws.send(JSON.stringify(registerMsg));
            } else {
                console.error('### No userID provided for WebSocket registration');
            }
        };
        
        ws.onmessage = function(event) {
            console.log('### WebSocket message received:', event.data);
            try {
                const data = JSON.parse(event.data);
                console.log('### Parsed message:', data);
                if (data.type === 'notification') {
                    console.log('### Notification received:', data.message);
                    handleNotification(data.message);
                } else {
                    console.log('### Unknown message type:', data.type);
                }
            } catch (e) {
                console.error('### WebSocket message parse error:', e, event.data);
            }
        };
        
        ws.onerror = function(error) {
            console.error('### WebSocket error:', error);
        };
        
        ws.onclose = function() {
            console.log('### WebSocket disconnected');
            // 재연결 시도
            setTimeout(() => connectWebSocket(userID), 3000);
        };
    }
    
    // 알림 처리
    function handleNotification(message) {
        unreadCount++;
        updateNotificationBadge();
        shakeBell();
    }
    
    // 알림 배지 업데이트
    function updateNotificationBadge() {
        if (unreadCount > 0) {
            notificationBadge.textContent = unreadCount > 99 ? '99+' : unreadCount;
            notificationBadge.style.display = 'flex';
        } else {
            notificationBadge.style.display = 'none';
        }
    }
    
    // 종 모양 흔들기 애니메이션
    function shakeBell() {
        notificationBell.classList.add('shake');
        setTimeout(() => {
            notificationBell.classList.remove('shake');
        }, 500);
    }
    
    // 종 모양 클릭 시 알림 확인 (나중에 알림 목록 표시)
    notificationBell.addEventListener('click', function() {
        // 알림 확인 처리 (나중에 구현)
        unreadCount = 0;
        updateNotificationBadge();
    });
    
    // WebSocket 연결 시작 (임시로 userID 1 사용)
    // 나중에 로그인 세션에서 가져와야 함
    currentUserID = 1; // 임시
    connectWebSocket(currentUserID);
});

