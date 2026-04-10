<script>
  import { onMount } from 'svelte';
  import { gsap } from 'gsap';
  import Slide from '../lib/Slide.svelte';

  const points = [
    {
      title: 'Roles & ClusterRoles',
      desc: 'Define permissions on resources',
      yaml: `kind: Role
metadata:
  name: pod-reader
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list"]`,
    },
    {
      title: 'RoleBindings',
      desc: 'Associate roles with users or service accounts',
      yaml: `kind: RoleBinding
metadata:
  name: read-pods
roleRef:
  kind: Role
  name: pod-reader
subjects:
- kind: ServiceAccount
  name: demo-user`,
    },
    {
      title: 'Service Accounts',
      desc: 'Identity for processes running in pods',
      yaml: `kind: ServiceAccount
metadata:
  name: demo-user
  namespace: demo`,
    },
    {
      title: 'Principle of Least Privilege',
      desc: 'Grant only what\'s needed — no wildcards, no cluster-admin',
      yaml: null,
    },
  ];

  let listEl;

  onMount(() => {
    if (listEl) {
      gsap.from(listEl.querySelectorAll('.rbac-item'), {
        opacity: 0,
        x: -20,
        duration: 0.4,
        stagger: 0.15,
        delay: 0.2,
        ease: 'power2.out',
      });
    }
  });
</script>

<Slide>
  <h2 class="text-4xl font-bold text-white mb-2">RBAC</h2>
  <p class="text-neutral-500 mb-6">Making sure only the right people make changes</p>

  <div bind:this={listEl} class="flex-1 min-h-0 grid grid-cols-2 gap-3">
    {#each points as point}
      <div class="rbac-item flex flex-col gap-1 p-3 rounded-xl border border-cyan-500/20 bg-cyan-500/5">
        <div>
          <h3 class="text-base font-semibold text-cyan-400">{point.title}</h3>
          <p class="text-neutral-400 text-xs">{point.desc}</p>
        </div>
        {#if point.yaml}
          <pre class="yaml-snippet text-neutral-300 bg-neutral-900/80 rounded-lg p-2 overflow-x-auto font-mono mt-auto"><code>{point.yaml}</code></pre>
        {/if}
      </div>
    {/each}
  </div>
</Slide>

<style>
  .yaml-snippet {
    font-size: 0.65rem;
    line-height: 1.4;
  }
</style>
