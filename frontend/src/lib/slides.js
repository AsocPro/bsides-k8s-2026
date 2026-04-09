import { writable, derived } from 'svelte/store';

export const slides = writable([]);
export const currentIndex = writable(0);

export const currentSlide = derived(
  [slides, currentIndex],
  ([$slides, $currentIndex]) => $slides[$currentIndex] || null
);

export const totalSlides = derived(slides, ($slides) => $slides.length);

// Sub-step support: slides can register a callback that consumes "next"
// presses before the slide advances. Return true to consume the press.
let _substepHandler = null;

export function registerSubsteps(handler) {
  _substepHandler = handler;
  return () => { _substepHandler = null; };
}

export function next() {
  if (_substepHandler && _substepHandler()) return;
  currentIndex.update((i) => {
    const total = getTotalSync();
    return Math.min(i + 1, total - 1);
  });
}

export function prev() {
  currentIndex.update((i) => Math.max(i - 1, 0));
}

export function goTo(index) {
  currentIndex.set(index);
}

// Internal helper — reads total from the store synchronously
let _total = 0;
function getTotalSync() {
  return _total;
}

// Keep _total in sync
totalSlides.subscribe((v) => (_total = v));
