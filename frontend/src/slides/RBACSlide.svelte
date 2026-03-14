<script>
  import Slide from '../lib/Slide.svelte';
  import SplitLayout from '../lib/SplitLayout.svelte';
  import AnimatedList from '../lib/AnimatedList.svelte';

  const points = [
    'Roles & ClusterRoles — define permissions on resources',
    'RoleBindings — associate roles with users or service accounts',
    'Service Accounts — identity for processes running in pods',
    'Principle of Least Privilege — grant only what\'s needed',
  ];
</script>

<Slide>
  <h2 class="text-4xl font-bold text-white mb-2">RBAC</h2>
  <p class="text-neutral-500 mb-6">Making sure only the right people make changes</p>

  <div class="flex-1 min-h-0">
    <AnimatedList items={points} />
  </div>
</Slide>
