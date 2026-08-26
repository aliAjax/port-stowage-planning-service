async function api(path, method, body) {
  const opts = { method: method || 'GET', headers: {} };
  if (body) { opts.headers['content-type'] = 'application/json'; opts.body = JSON.stringify(body); }
  const res = await fetch(path, opts);
  if (!res.ok) throw new Error((await res.json()).error || res.statusText);
  return res.status === 204 ? null : res.json();
}
async function refresh() {
  try {
    const plans = await api('/api/v1/plans');
    const list = document.getElementById('voyage-list');
    list.innerHTML = '';
    plans.forEach(p => { const li = document.createElement('li'); li.textContent = p.id + ' [' + p.state + '] ' + p.decisions.length + ' 箱'; list.appendChild(li); });
    document.getElementById('status').textContent = '服务正常';
    document.getElementById('status').className = 'ok';
  } catch (e) {
    document.getElementById('status').textContent = '服务异常: ' + e.message;
    document.getElementById('status').className = 'err';
  }
}
document.getElementById('voyage-form').addEventListener('submit', async (ev) => {
  ev.preventDefault();
  const f = new FormData(ev.target);
  const eta = new Date(Date.now() + 3600e3).toISOString();
  const etd = new Date(Date.now() + 86400e3).toISOString();
  await api('/api/v1/voyages', 'POST', { id: f.get('id'), vessel_id: f.get('vessel_id'), port_id: f.get('port_id'), eta, etd });
  refresh();
});
document.getElementById('container-form').addEventListener('submit', async (ev) => {
  ev.preventDefault();
  const f = new FormData(ev.target);
  const c = { id: f.get('id'), voyage_id: f.get('voyage_id'), iso_size: f.get('iso_size'), weight_kg: Number(f.get('weight_kg')), destination: f.get('destination'), priority: Number(f.get('priority')), on_deck: true };
  try { await api('/api/v1/containers', 'POST', c); document.getElementById('container-msg').textContent = '已注册'; }
  catch (e) { document.getElementById('container-msg').textContent = e.message; }
  refresh();
});
document.getElementById('plan-form').addEventListener('submit', async (ev) => {
  ev.preventDefault();
  const f = new FormData(ev.target);
  try {
    const out = await api('/api/v1/plans/solve', 'POST', { id: f.get('id'), voyage_id: f.get('voyage_id') });
    document.getElementById('plan-output').textContent = JSON.stringify(out, null, 2);
  } catch (e) { document.getElementById('plan-output').textContent = e.message; }
  refresh();
});
refresh();
