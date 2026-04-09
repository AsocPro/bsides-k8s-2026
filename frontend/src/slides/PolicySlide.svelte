<script>
  import Slide from '../lib/Slide.svelte';
  import AnimatedList from '../lib/AnimatedList.svelte';

  const policies = [
    'Require all containers to run as non-root',
    'Block privileged containers',
    'Enforce resource limits (CPU/memory) on all pods',
    'Mandate images come from approved registries only',
    'Prevent use of latest tag — require explicit versioning',
  ];
</script>

<Slide>
  <h2 class="text-4xl font-bold text-white mb-2">Policy Agents</h2>
  <p class="text-neutral-500 mb-6">Making sure the changes are reasonable</p>

  <div class="mb-6 space-y-3">
    <p class="text-neutral-300 leading-relaxed">
      RBAC controls <span class="text-cyan-400 font-semibold">who</span> can make changes.
    </p>
    <p class="text-neutral-300 leading-relaxed">
      Policy agents control <span class="text-amber-400 font-semibold">what</span> changes are allowed.
    </p>
    <p class="text-neutral-400 text-sm italic pl-8 mt-1">
      Enforced via admission webhooks before resources hit the cluster.
    </p>
  </div>

  <h3 class="text-lg font-semibold text-emerald-400 mb-4">Policies You Can Enforce</h3>
  <AnimatedList items={policies} />
</Slide>
