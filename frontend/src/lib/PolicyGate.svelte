<script>
  import { gsap } from 'gsap';

  let { data } = $props();

  let badgeEls = {};
  let prevValues = {};
  let hasAnimatedIn = false;

  $effect(() => {
    if (!data) return;

    // First time: stagger-animate all rows in
    if (!hasAnimatedIn) {
      hasAnimatedIn = true;
      // Defer so bind:this has resolved
      requestAnimationFrame(() => {
        Object.values(badgeEls).forEach((el, i) => {
          if (el) {
            gsap.fromTo(el.closest('[data-row]'),
              { opacity: 0, x: -10 },
              { opacity: 1, x: 0, duration: 0.3, delay: i * 0.1, ease: 'power2.out' }
            );
          }
        });
      });
    }

    // On subsequent updates: only animate badges whose value changed
    const allEntries = {
      ...data.policies,
      ...data.probes,
    };

    for (const [key, val] of Object.entries(allEntries)) {
      if (prevValues[key] !== undefined && prevValues[key] !== val && badgeEls[key]) {
        gsap.fromTo(badgeEls[key],
          { scale: 1.3 },
          { scale: 1, duration: 0.3, ease: 'back.out(2)' }
        );
      }
      prevValues[key] = val;
    }
  });

  const policyLabels = {
    'require-non-root': 'Require Non-Root',
    'require-approved-registries': 'Approved Registries',
  };

  const probeLabels = {
    'root-blocked': 'Root container blocked',
    'nonroot-allowed': 'Non-root container allowed',
  };
</script>

<div class="space-y-1.5">
  <div class="text-[10px] font-mono text-amber-400/60 uppercase tracking-wider">Policies</div>
  {#each Object.entries(data?.policies ?? {}) as [key, active]}
    <div data-row class="flex items-center justify-between px-2 py-1 rounded bg-neutral-800/50">
      <span class="text-xs text-neutral-300 font-mono">{policyLabels[key] || key}</span>
      <span
        bind:this={badgeEls[key]}
        class="text-[10px] font-mono px-1.5 py-0.5 rounded {active ? 'bg-emerald-500/20 text-emerald-400' : 'bg-red-500/20 text-red-400'}"
      >
        {active ? 'active' : 'missing'}
      </span>
    </div>
  {/each}

  <div class="text-[10px] font-mono text-amber-400/60 uppercase tracking-wider mt-2">Probes</div>
  {#each Object.entries(data?.probes ?? {}) as [key, passed]}
    <div data-row class="flex items-center justify-between px-2 py-1 rounded bg-neutral-800/50">
      <span class="text-xs text-neutral-300 font-mono">{probeLabels[key] || key}</span>
      <span
        bind:this={badgeEls[key]}
        class="text-[10px] font-mono px-1.5 py-0.5 rounded {passed ? 'bg-emerald-500/20 text-emerald-400' : 'bg-red-500/20 text-red-400'}"
      >
        {passed ? 'pass' : 'fail'}
      </span>
    </div>
  {/each}
</div>
