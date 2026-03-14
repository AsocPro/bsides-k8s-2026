import { writable, derived } from 'svelte/store';

export const slides = writable([]);
export const currentIndex = writable(0);

export const currentSlide = derived(
  [slides, currentIndex],
  ([$slides, $currentIndex]) => $slides[$currentIndex] || null
);

export const totalSlides = derived(slides, ($slides) => $slides.length);

export function next() {
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
