async function loadDashboard() {
    try {
        const response = await fetch('/history.json');
        if (!response.ok) throw new Error("No data found");
        
        const data = await response.json();
        const sessions = Array.isArray(data) ? data : [data];

        updateSummary(sessions);
        updateTable(sessions);
    } catch (err) {
        console.error(err);
        document.getElementById('summary').innerHTML = `<p>Ready to start tracking!</p>`;
    }
}

function updateSummary(sessions) {
    const totals = {};
    sessions.forEach(s => {
        totals[s.category] = (totals[s.category] || 0) + s.duration;
    });

    const summaryDiv = document.getElementById('summary');
    summaryDiv.innerHTML = '';

    for (const [cat, dur] of Object.entries(totals)) {
        const mins = Math.floor(dur / 1e9 / 60);
        const card = document.createElement('div');
        card.className = 'stat-card';
        card.innerHTML = `
            <h3>${cat}</h3>
            <p>${mins}m</p>
        `;
        summaryDiv.appendChild(card);
    }
}

function updateTable(sessions) {
    const tbody = document.querySelector('#historyTable tbody');
    tbody.innerHTML = '';

    sessions.reverse().forEach(s => {
        const date = new Date(s.start_time).toLocaleDateString();
        const mins = Math.floor(s.duration / 1e9 / 60);
        const secs = Math.floor((s.duration / 1e9) % 60);

        const row = `
            <tr>
                <td>${date}</td>
                <td><strong>${s.category}</strong></td>
                <td>${mins}m ${secs}s</td>
            </tr>
        `;
        tbody.insertAdjacentHTML('beforeend', row);
    });
}

loadDashboard();
