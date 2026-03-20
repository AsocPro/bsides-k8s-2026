import { writable, get } from 'svelte/store';

function createStateStore() {
  return writable({ loading: false, error: null, data: null });
}

export const rbacState = createStateStore();
export const policyState = createStateStore();
export const netpolState = createStateStore();

const storeMap = {
  rbac: rbacState,
  policy: policyState,
  netpol: netpolState,
};

// Transform goss JSON output into per-demo structured data
function transformGossResults(demo, gossOutput) {
  const results = {};
  for (const result of gossOutput.results) {
    results[result['resource-id']] = result.successful;
  }

  switch (demo) {
    case 'rbac':
      return transformRBAC(results);
    case 'policy':
      return transformPolicy(results);
    case 'netpol':
      return transformNetPol(results);
    default:
      return results;
  }
}

function transformRBAC(results) {
  const resources = ['pods', 'roles', 'rolebindings'];
  const verbs = ['get', 'create', 'delete'];
  const permissions = {};

  for (const resource of resources) {
    permissions[resource] = {};
    for (const verb of verbs) {
      const key = `can-i-${verb}-${resource}`;
      permissions[resource][verb] = results[key] ?? false;
    }
  }

  return { permissions };
}

function transformPolicy(results) {
  return {
    policies: {
      'require-non-root': results['policy-require-non-root'] ?? false,
      'require-approved-registries': results['policy-require-approved-registries'] ?? false,
    },
    probes: {
      'root-blocked': results['probe-root-container-blocked'] ?? false,
      'nonroot-allowed': results['probe-nonroot-container-allowed'] ?? false,
    },
  };
}

function transformNetPol(results) {
  return {
    connections: {
      'frontend->api': results['conn-frontend-to-api'] ?? false,
      'frontend->database': results['conn-frontend-to-database'] ?? false,
      'api->database': results['conn-api-to-database'] ?? false,
    },
  };
}

const pollingTimers = {};

export function startPolling(demo, ms = 1000) {
  stopPolling(demo);
  checkState(demo);
  pollingTimers[demo] = setInterval(() => checkState(demo), ms);
}

export function stopPolling(demo) {
  if (pollingTimers[demo]) {
    clearInterval(pollingTimers[demo]);
    delete pollingTimers[demo];
  }
}

export async function checkState(demo) {
  const store = storeMap[demo];
  if (!store) return;

  store.set({ loading: true, error: null, data: get(store).data });

  try {
    const res = await fetch(`/api/state/${demo}`, { method: 'POST' });
    if (!res.ok) {
      throw new Error(`State check failed: ${res.statusText}`);
    }
    const gossOutput = await res.json();
    const data = transformGossResults(demo, gossOutput);
    store.set({ loading: false, error: null, data });
  } catch (err) {
    store.set({ loading: false, error: err.message, data: get(store).data });
  }
}
