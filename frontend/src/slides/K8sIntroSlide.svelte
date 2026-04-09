<script>
  import Slide from '../lib/Slide.svelte';
  import { onMount } from 'svelte';
  import { gsap } from 'gsap';
  import { registerSubsteps } from '../lib/slides.js';

  let diagramEl;
  let textEl;
  let phase = 0;
  const totalPhases = 4; // namespace, deployment, pod, service

  function revealPhase(p) {
    if (!diagramEl || !textEl) return;

    if (p === 1) {
      // Namespaces
      gsap.to(textEl.querySelector('.concept-ns'), { opacity: 1, x: 0, duration: 0.5, ease: 'power2.out' });
      gsap.fromTo(diagramEl.querySelectorAll('.phase-ns'), { scale: 0.9 }, { scale: 1, duration: 0.5, ease: 'power2.out' });
      gsap.to(diagramEl.querySelectorAll('.phase-ns'), { opacity: 1, duration: 0.5, ease: 'power2.out' });
    } else if (p === 2) {
      // Deployments
      gsap.to(textEl.querySelector('.concept-deploy'), { opacity: 1, x: 0, duration: 0.5, ease: 'power2.out' });
      gsap.fromTo(diagramEl.querySelectorAll('.phase-deploy'), { y: -15 }, { y: 0, duration: 0.4, ease: 'power2.out' });
      gsap.to(diagramEl.querySelectorAll('.phase-deploy'), { opacity: 1, duration: 0.4, ease: 'power2.out' });
    } else if (p === 3) {
      // Pods
      gsap.to(textEl.querySelector('.concept-pod'), { opacity: 1, x: 0, duration: 0.5, ease: 'power2.out' });
      gsap.fromTo(diagramEl.querySelectorAll('.phase-pod'), { scale: 0 }, { scale: 1, duration: 0.3, stagger: 0.06, ease: 'back.out(1.7)' });
      gsap.to(diagramEl.querySelectorAll('.phase-pod'), { opacity: 1, duration: 0.3, stagger: 0.06, ease: 'power2.out' });
    } else if (p === 4) {
      // Services + connections + ports
      gsap.to(textEl.querySelector('.concept-svc'), { opacity: 1, x: 0, duration: 0.5, ease: 'power2.out' });
      gsap.fromTo(diagramEl.querySelectorAll('.phase-svc'), { x: 20 }, { x: 0, duration: 0.4, ease: 'power2.out' });
      gsap.to(diagramEl.querySelectorAll('.phase-svc'), { opacity: 1, duration: 0.4, ease: 'power2.out' });
      gsap.to(diagramEl.querySelectorAll('.phase-conn'), { opacity: 1, duration: 0.4, stagger: 0.05, delay: 0.15, ease: 'power2.out' });
      gsap.to(diagramEl.querySelectorAll('.phase-port'), { opacity: 1, duration: 0.3, delay: 0.3, ease: 'back.out(2)' });
    }
  }

  onMount(() => {
    // Hide all animated elements initially
    gsap.set(diagramEl.querySelectorAll('.phase-ns, .phase-deploy, .phase-pod, .phase-svc, .phase-conn, .phase-port'), { opacity: 0 });
    gsap.set(textEl.querySelectorAll('.concept'), { opacity: 0, x: -20 });

    // Legend is visible from the start — animate it in gently
    gsap.from(diagramEl.querySelector('.legend'), { opacity: 0, duration: 0.6, delay: 0.2, ease: 'power2.out' });

    // Register sub-step handler: consume right-arrow until all phases shown
    const unregister = registerSubsteps(() => {
      if (phase < totalPhases) {
        phase++;
        revealPhase(phase);
        return true; // consumed
      }
      return false; // let slide engine advance
    });

    return unregister;
  });
</script>

<Slide>
  <h2 class="text-4xl font-bold text-white mb-6">Kubernetes Core Concepts</h2>

  <div class="grid grid-cols-2 gap-8 flex-1 min-h-0">
    <!-- Left: Text descriptions, revealed one by one -->
    <div bind:this={textEl} class="flex flex-col justify-center space-y-6">
      <div class="concept concept-ns flex items-start gap-4">
        <div class="shrink-0 w-4 h-4 mt-1 rounded-sm border-2 border-dashed" style="border-color: rgb(6,182,212); background: rgba(6,182,212,0.08);"></div>
        <div>
          <h3 class="text-lg font-semibold text-cyan-400">Namespaces</h3>
          <p class="text-neutral-400 text-sm mt-1">Logical isolation boundaries for resources and policies. Think of them as folders that keep teams and environments separated.</p>
        </div>
      </div>

      <div class="concept concept-deploy flex items-start gap-4">
        <div class="shrink-0 w-4 h-4 mt-1 rounded" style="border: 2px solid rgb(139,92,246); background: rgba(139,92,246,0.12);"></div>
        <div>
          <h3 class="text-lg font-semibold text-violet-400">Deployments</h3>
          <p class="text-neutral-400 text-sm mt-1">Declare your desired replica count and update strategy. Kubernetes ensures reality matches your intent.</p>
        </div>
      </div>

      <div class="concept concept-pod flex items-start gap-4">
        <div class="shrink-0 w-4 h-4 mt-1 rounded" style="border: 2px solid rgb(34,211,238); background: rgba(34,211,238,0.15);"></div>
        <div>
          <h3 class="text-lg font-semibold text-cyan-300">Pods</h3>
          <p class="text-neutral-400 text-sm mt-1">The smallest deployable unit — one or more containers sharing network and storage. Ephemeral by design.</p>
        </div>
      </div>

      <div class="concept concept-svc flex items-start gap-4">
        <div class="shrink-0 w-4 h-4 mt-1 rounded-full" style="border: 2px solid rgb(16,185,129); background: rgba(16,185,129,0.12);"></div>
        <div>
          <h3 class="text-lg font-semibold text-emerald-400">Services</h3>
          <p class="text-neutral-400 text-sm mt-1">Stable networking abstraction that routes traffic to pods. Expose ports internally or externally.</p>
        </div>
      </div>
    </div>

    <!-- Right: Visual diagram -->
    <div bind:this={diagramEl} class="flex-1 min-h-0 flex items-center">
      <svg viewBox="0 0 340 340" class="w-full h-full" xmlns="http://www.w3.org/2000/svg">
        <!-- Namespace: frontend -->
        <rect class="phase-ns" x="10" y="10" width="250" height="150" rx="8"
              fill="rgba(6,182,212,0.08)" stroke="rgb(6,182,212)" stroke-width="1.5" stroke-dasharray="6 3" />
        <text class="phase-ns" x="20" y="30" fill="rgb(6,182,212)" font-size="11" font-weight="600">namespace: frontend</text>

        <!-- Deployment in frontend ns -->
        <rect class="phase-deploy" x="20" y="42" width="135" height="55" rx="6"
              fill="rgba(139,92,246,0.12)" stroke="rgb(139,92,246)" stroke-width="1.2" />
        <text class="phase-deploy" x="30" y="57" fill="rgb(167,139,250)" font-size="9" font-weight="500">Deployment</text>
        <text class="phase-deploy" x="30" y="69" fill="rgb(196,181,253)" font-size="8">replicas: 3</text>

        <!-- Pods inside frontend deployment -->
        <rect class="phase-pod" x="27" y="76" width="34" height="16" rx="4"
              fill="rgba(34,211,238,0.15)" stroke="rgb(34,211,238)" stroke-width="1" />
        <text class="phase-pod" x="34" y="87" fill="rgb(34,211,238)" font-size="7">Pod</text>

        <rect class="phase-pod" x="66" y="76" width="34" height="16" rx="4"
              fill="rgba(34,211,238,0.15)" stroke="rgb(34,211,238)" stroke-width="1" />
        <text class="phase-pod" x="73" y="87" fill="rgb(34,211,238)" font-size="7">Pod</text>

        <rect class="phase-pod" x="105" y="76" width="34" height="16" rx="4"
              fill="rgba(34,211,238,0.15)" stroke="rgb(34,211,238)" stroke-width="1" />
        <text class="phase-pod" x="112" y="87" fill="rgb(34,211,238)" font-size="7">Pod</text>

        <!-- Service in frontend ns -->
        <rect class="phase-svc" x="170" y="50" width="55" height="50" rx="25"
              fill="rgba(16,185,129,0.12)" stroke="rgb(16,185,129)" stroke-width="1.2" />
        <text class="phase-svc" x="178" y="73" fill="rgb(52,211,153)" font-size="9" font-weight="500">Service</text>
        <text class="phase-svc" x="187" y="87" fill="rgb(110,231,183)" font-size="8">:80</text>

        <!-- Connection lines: frontend service -> pods -->
        <line class="phase-conn" x1="170" y1="75" x2="139" y2="84" stroke="rgb(16,185,129)" stroke-width="0.8" stroke-dasharray="3 2" />
        <line class="phase-conn" x1="170" y1="75" x2="100" y2="84" stroke="rgb(16,185,129)" stroke-width="0.8" stroke-dasharray="3 2" />
        <line class="phase-conn" x1="170" y1="75" x2="61" y2="84" stroke="rgb(16,185,129)" stroke-width="0.8" stroke-dasharray="3 2" />

        <!-- Port exposure arrow -->
        <line class="phase-port" x1="225" y1="75" x2="248" y2="75" stroke="rgb(251,191,36)" stroke-width="1.5" />
        <polygon class="phase-port" points="248,71 256,75 248,79" fill="rgb(251,191,36)" />
        <text class="phase-port" x="258" y="78" fill="rgb(251,191,36)" font-size="9" font-weight="600">:80</text>

        <!-- Namespace: backend -->
        <rect class="phase-ns" x="10" y="180" width="250" height="150" rx="8"
              fill="rgba(139,92,246,0.06)" stroke="rgb(139,92,246)" stroke-width="1.5" stroke-dasharray="6 3" />
        <text class="phase-ns" x="20" y="200" fill="rgb(139,92,246)" font-size="11" font-weight="600">namespace: backend</text>

        <!-- Deployment in backend ns -->
        <rect class="phase-deploy" x="20" y="210" width="135" height="55" rx="6"
              fill="rgba(139,92,246,0.12)" stroke="rgb(139,92,246)" stroke-width="1.2" />
        <text class="phase-deploy" x="30" y="225" fill="rgb(167,139,250)" font-size="9" font-weight="500">Deployment</text>
        <text class="phase-deploy" x="30" y="237" fill="rgb(196,181,253)" font-size="8">replicas: 2</text>

        <!-- Pods inside backend deployment -->
        <rect class="phase-pod" x="35" y="244" width="34" height="16" rx="4"
              fill="rgba(34,211,238,0.15)" stroke="rgb(34,211,238)" stroke-width="1" />
        <text class="phase-pod" x="42" y="255" fill="rgb(34,211,238)" font-size="7">Pod</text>

        <rect class="phase-pod" x="78" y="244" width="34" height="16" rx="4"
              fill="rgba(34,211,238,0.15)" stroke="rgb(34,211,238)" stroke-width="1" />
        <text class="phase-pod" x="85" y="255" fill="rgb(34,211,238)" font-size="7">Pod</text>

        <!-- Standalone DB pod -->
        <rect class="phase-pod" x="25" y="275" width="75" height="20" rx="4"
              fill="rgba(244,114,182,0.12)" stroke="rgb(244,114,182)" stroke-width="1" />
        <text class="phase-pod" x="33" y="288" fill="rgb(244,114,182)" font-size="8">Pod: postgres</text>

        <!-- Service in backend ns -->
        <rect class="phase-svc" x="170" y="218" width="55" height="50" rx="25"
              fill="rgba(16,185,129,0.12)" stroke="rgb(16,185,129)" stroke-width="1.2" />
        <text class="phase-svc" x="175" y="241" fill="rgb(52,211,153)" font-size="9" font-weight="500">Service</text>
        <text class="phase-svc" x="182" y="255" fill="rgb(110,231,183)" font-size="8">:8080</text>

        <!-- DB Service -->
        <rect class="phase-svc" x="115" y="273" width="55" height="24" rx="12"
              fill="rgba(16,185,129,0.12)" stroke="rgb(16,185,129)" stroke-width="1" />
        <text class="phase-svc" x="120" y="289" fill="rgb(110,231,183)" font-size="8">Svc :5432</text>

        <!-- Connection lines: backend service -> pods -->
        <line class="phase-conn" x1="170" y1="243" x2="112" y2="252" stroke="rgb(16,185,129)" stroke-width="0.8" stroke-dasharray="3 2" />
        <line class="phase-conn" x1="170" y1="243" x2="69" y2="252" stroke="rgb(16,185,129)" stroke-width="0.8" stroke-dasharray="3 2" />
        <line class="phase-conn" x1="115" y1="285" x2="100" y2="285" stroke="rgb(16,185,129)" stroke-width="0.8" stroke-dasharray="3 2" />

        <!-- Port exposure from backend service -->
        <line class="phase-port" x1="225" y1="243" x2="248" y2="243" stroke="rgb(251,191,36)" stroke-width="1.5" />
        <polygon class="phase-port" points="248,239 256,243 248,247" fill="rgb(251,191,36)" />
        <text class="phase-port" x="258" y="246" fill="rgb(251,191,36)" font-size="9" font-weight="600">:8080</text>

        <!-- Namespace boundary annotation -->
        <path class="phase-ns" d="M130,160 L130,180" stroke="rgb(251,146,60)" stroke-width="1" stroke-dasharray="4 2" fill="none" />
        <text class="phase-ns" x="140" y="174" fill="rgb(251,146,60)" font-size="8" font-style="italic">namespace boundary</text>

        <!-- Legend (always visible) -->
        <g class="legend" transform="translate(280, 10)">
          <text fill="rgb(100,116,139)" font-size="9" font-weight="600" y="0">Legend</text>
          <rect x="0" y="10" width="12" height="8" rx="2" fill="rgba(6,182,212,0.08)" stroke="rgb(6,182,212)" stroke-width="0.8" stroke-dasharray="3 2" />
          <text x="16" y="17" fill="rgb(148,163,184)" font-size="8">Namespace</text>
          <rect x="0" y="24" width="12" height="8" rx="2" fill="rgba(139,92,246,0.12)" stroke="rgb(139,92,246)" stroke-width="0.8" />
          <text x="16" y="31" fill="rgb(148,163,184)" font-size="8">Deployment</text>
          <rect x="0" y="38" width="12" height="8" rx="2" fill="rgba(34,211,238,0.15)" stroke="rgb(34,211,238)" stroke-width="0.8" />
          <text x="16" y="45" fill="rgb(148,163,184)" font-size="8">Pod</text>
          <rect x="0" y="52" width="12" height="8" rx="6" fill="rgba(16,185,129,0.12)" stroke="rgb(16,185,129)" stroke-width="0.8" />
          <text x="16" y="59" fill="rgb(148,163,184)" font-size="8">Service</text>
          <line x1="0" y1="70" x2="12" y2="70" stroke="rgb(251,191,36)" stroke-width="1.2" />
          <polygon points="12,67.5 17,70 12,72.5" fill="rgb(251,191,36)" />
          <text x="20" y="73" fill="rgb(148,163,184)" font-size="8">Port</text>
        </g>
      </svg>
    </div>
  </div>
</Slide>
