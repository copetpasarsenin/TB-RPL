let currentToken = null;

// Saat tombol generate token diklik
document.getElementById('generateTokenBtn').addEventListener('click', async () => {
    try {
        const response = await fetch('/auth/token', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ user_id: 1, username: 'admin_dashboard', role: 'admin' })
        });
        const data = await response.json();
        
        if (data.status === 'success') {
            currentToken = data.data.token;
            const badge = document.getElementById('tokenStatus');
            badge.innerHTML = '✅ Token Active';
            badge.className = 'status-badge online';
            
            // Langsung otomatis refresh log setelah dapet token
            fetchLogs();
        }
    } catch (error) {
        console.error('Error fetching token:', error);
        alert('Gagal mengambil token. Pastikan server API jalan.');
    }
});

// Tombol refresh manual
document.getElementById('refreshLogsBtn').addEventListener('click', () => {
    if (!currentToken) {
        alert("Silakan Generate Akses Token terlebih dahulu!");
        return;
    }
    fetchLogs();
});

// Cek status server (Health Check)
async function checkHealth() {
    try {
        const response = await fetch('/health');
        const data = await response.json();
        const el = document.getElementById('healthStatus');
        
        if(data.status === 'success') {
            el.textContent = '✅ Online';
            el.style.color = '#10b981';
        } else {
            el.textContent = '❌ Offline';
            el.style.color = '#ef4444';
        }
    } catch (error) {
        const el = document.getElementById('healthStatus');
        el.textContent = '❌ Offline';
        el.style.color = '#ef4444';
    }
}

// Ambil data log dari /integrator/logging
async function fetchLogs() {
    if (!currentToken) return;
    
    try {
        // Menggunakan endpoint spesifik sesuai spreadsheet dosen
        const response = await fetch('/integrator/logging', {
            headers: { 'Authorization': `Bearer ${currentToken}` }
        });
        const data = await response.json();
        
        if (data.status === 'success' && data.data) {
            const logs = data.data;
            document.getElementById('totalRequests').textContent = logs.length;
            
            const tbody = document.querySelector('#logsTable tbody');
            tbody.innerHTML = '';
            
            if(logs.length === 0) {
                tbody.innerHTML = '<tr><td colspan="5" class="empty-state">Belum ada request yang tercatat hari ini.</td></tr>';
                return;
            }
            
            // Menampilkan 15 log terbaru
            logs.slice(0, 15).forEach(log => {
                const tr = document.createElement('tr');
                
                const methodRaw = log.method.toLowerCase();
                const methodClass = methodRaw === 'get' ? 'get' : methodRaw === 'post' ? 'post' : 'options';
                const statusClass = log.status_code >= 200 && log.status_code < 300 ? 'success' : 'error';
                
                // Format waktu
                const date = new Date(log.created_at).toLocaleString('id-ID');
                
                tr.innerHTML = `
                    <td>${date}</td>
                    <td><span class="badge ${methodClass}">${log.method}</span></td>
                    <td style="font-family: monospace;">${log.path}</td>
                    <td><span class="badge ${statusClass}">${log.status_code}</span></td>
                    <td><span style="opacity:0.7">${log.source_ip}</span></td>
                `;
                tbody.appendChild(tr);
            });
        }
    } catch (error) {
        console.error('Error fetching logs:', error);
    }
}

// Jalankan check health saat pertama kali dibuka
checkHealth();
// Update health check setiap 10 detik
setInterval(checkHealth, 10000);
