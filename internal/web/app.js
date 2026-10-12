const refresh = document.querySelector('#refresh');
async function check(path, element, healthy) {
  element.textContent = 'Checking…';
  try {
    const response = await fetch(path, { cache: 'no-store', signal: AbortSignal.timeout(5000) });
    element.textContent = response.ok ? healthy : 'Unavailable';
  } catch {
    element.textContent = 'Unable to connect';
  }
}
async function refreshStatus() {
  refresh.disabled = true;
  await Promise.all([
    check('/healthz', document.querySelector('#api-status'), 'Online'),
    check('/readyz', document.querySelector('#db-status'), 'Connected'),
  ]);
  refresh.disabled = false;
}
refresh.addEventListener('click', refreshStatus);
refreshStatus();
