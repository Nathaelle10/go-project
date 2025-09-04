const ctx = document.getElementById('cpu-chart').getContext('2d');
const now = new Date().toLocaleTimeString();
let cpuChart;

// Navigation 
document.querySelectorAll('.page').forEach(p => p.style.display = 'none');
document.getElementById('cpu').style.display = 'block'; // Par défaut

document.querySelectorAll('sl-tab').forEach(tab => {
  tab.addEventListener('click', () => {
    const target = tab.getAttribute('panel');
    document.querySelectorAll('.page').forEach(p => p.style.display = 'none');
    document.getElementById(target).style.display = 'block';
  });
});



async function fetchStatus() {
  try {
    const res = await fetch("http://localhost:9090/api/status");
    const data = await res.json();

    const disk = data.disk || {};
    const memory = data.memory || {};
    const host = data.host || {};
    const load = data.load || {};
    const network = data.network || [];

    document.getElementById("card-cpu").innerHTML = `
      <div class="label">🧠 CPU</div>
      Utilisation: ${data.cpu_usage?.toFixed(2) || "N/A"}%<br>
      Modèle: ${data.cpu_model || "N/A"}<br>
      Cœurs: ${data.cpu_cores || "N/A"}<br>
      Charge 1/5/15 min: ${load.load1?.toFixed(2) || "?"}, ${load.load5?.toFixed(2) || "?"}, ${load.load15?.toFixed(2) || "?"}
    `;

    document.getElementById("card-memory").innerHTML = `
      <div class="label">💾 Mémoire</div>
      Utilisation: ${memory.used && memory.total ? ((memory.used / memory.total) * 100).toFixed(1) : "?"}%<br>
      RAM: ${memory.total ? (memory.total / 1e9).toFixed(2) : "?"} GB
    `;

    document.getElementById("card-disk").innerHTML = `
      <div class="label">📁 Disque</div>
      Utilisation: ${disk.usedPercent?.toFixed(1) || "?"}%<br>
      Taille: ${disk.total ? (disk.total / 1e9).toFixed(2) : "?"} GB
    `;

    document.getElementById("card-net").innerHTML = `
      <div class="label">🌐 Réseau</div>
      Interfaces: ${network.length}<br>
      Connexions actives: ${data.connections ?? "?"}
    `;

    document.getElementById("card-host").innerHTML = `
      <div class="label">🖥️ Système</div>
      OS: ${host.platform || "?"} ${host.platformVersion || ""}<br>
      Uptime: ${(data.uptime / 3600).toFixed(1) || "?"} heures
    `;

    updateChart(data);

  } catch (err) {
    console.error("Erreur API collector:", err);
  }
}

async function scanIP(ip) {
  if (!ip) return alert("Veuillez entrer une IP valide");
  try {
    const res = await fetch(`http://localhost:9090/api/scan?ip=${ip}`);
    if (!res.ok) throw new Error(`Erreur HTTP ${res.status}`);
    const data = await res.json();

    const portsList = data.openPorts.length
      ? data.openPorts.map(port => `<li>Port ${port}</li>`).join("")
      : "<li>Aucun port ouvert détecté</li>";

    document.getElementById("scan-results").innerHTML = `<ul>${portsList}</ul>`;

  } catch (err) {
    document.getElementById("scan-results").innerHTML = `<div>Erreur de scan : ${err.message}</div>`;
  }
}

function updateChart(data) {
  const usage = data.cpu_usage ?? 0;
  const now = new Date().toLocaleTimeString();

  if (!cpuChart) {
    const ctx = document.getElementById('cpu-chart').getContext('2d');
    cpuChart = new Chart(ctx, {
      type: 'line',
      data: {
        labels: [now],
        datasets: [{
          label: 'Utilisation CPU (%)',
          data: [usage],
          fill: false,
          borderColor: 'rgb(75, 192, 192)',
          tension: 0.1
        }]
      }
    });
  } else {
    cpuChart.data.labels.push(now);
    cpuChart.data.datasets[0].data.push(usage);
    if (cpuChart.data.labels.length > 10) {
      cpuChart.data.labels.shift();
      cpuChart.data.datasets[0].data.shift();
    }
    cpuChart.update();
  }
}


setInterval(fetchStatus, 3000);
fetchStatus();
