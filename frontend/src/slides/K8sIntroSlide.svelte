<script>
  import Slide from '../lib/Slide.svelte';
  import AnimatedList from '../lib/AnimatedList.svelte';

  const concepts = [
    'Pods — smallest deployable unit, one or more containers sharing network/storage',
    'Deployments — declare replica count and rolling update strategy',
    'Services — stable networking abstraction routing traffic to pods',
    'Namespaces — logical isolation boundaries for resources and policies',
  ];

  const securityFeatures = [
    'All changes go through the API server — fully auditable',
    'Secrets management is a first-class concept',
    'Resource isolation via namespaces, quotas, and security contexts',
    'Security posture defined in code — reviewable and version controlled',
  ];
</script>

<Slide>
  <h2 class="text-4xl font-bold text-white mb-8">Kubernetes in 60 Seconds</h2>

  <div class="grid grid-cols-2 gap-12 flex-1 min-h-0">
    <div>
      <h3 class="text-lg font-semibold text-cyan-400 mb-4">Core Concepts</h3>
      <AnimatedList items={concepts} />
    </div>
    <div>
      <h3 class="text-lg font-semibold text-emerald-400 mb-4">Built-in Security</h3>
      <AnimatedList items={securityFeatures} delay={0.6} />
    </div>
  </div>
</Slide>
